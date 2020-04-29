package eventer

import (
	"time"

	mapset "github.com/deckarep/golang-set"
	"github.com/sirupsen/logrus"

	"github.com/ShotaKitazawa/minecraft-bot/pkg/botplug"
	"github.com/ShotaKitazawa/minecraft-bot/pkg/domain"
	"github.com/ShotaKitazawa/minecraft-bot/pkg/domain/i18n"
	"github.com/ShotaKitazawa/minecraft-bot/pkg/rcon"
	"github.com/ShotaKitazawa/minecraft-bot/pkg/sharedmem"
)

const (
	cronJobInterval = 10
)

type Eventer struct {
	botplug.BotPluginSender

	MinecraftHostname string
	sharedMem         sharedmem.SharedMem
	rcon              rcon.RconClient
	Logger            *logrus.Logger
}

func New(minecraftHostname string, sender botplug.BotPluginSender, m sharedmem.SharedMem, rcon rcon.RconClient, logger *logrus.Logger) (*Eventer, error) {
	return &Eventer{
		BotPluginSender:   sender,
		MinecraftHostname: minecraftHostname,
		sharedMem:         m,
		rcon:              rcon,
		Logger:            logger,
	}, nil
}

func (e *Eventer) Run() error {
	return e.cronjob()
}

func (e *Eventer) cronjob() error {
	if err := e.job(); err != nil {
		e.Logger.Error(err)
	}
	t := time.NewTicker(cronJobInterval * time.Second)
	for {
		select {
		case <-t.C:
			if err := e.job(); err != nil {
				e.Logger.Error(err)
			}
		}
	}
	// t.Stop()
	// return nil
}

func (e *Eventer) job() error {
	var err error

	// get Minecraft metrics by RCON
	currentData, err := e.getMetricsUsingRCON()
	if err != nil {
		return err
	}

	// create currentLoginUsernameSet
	currentLoginUsernameSet := mapset.NewSet()
	for _, loginUser := range currentData.LoginUsers {
		currentLoginUsernameSet.Add(loginUser.Name)
	}

	// get logged in users from SharedMem
	previousData, err := e.sharedMem.SyncReadEntityFromSharedMem()
	if err != nil {
		// write to sharedMem & return
		return e.sharedMem.AsyncWriteEntityToSharedMem(currentData)
	}

	// create previousLoginUsernameSet
	previousLoginUsernameSet := mapset.NewSet()
	for _, previousLoginUser := range previousData.LoginUsers {
		previousLoginUsernameSet.Add(previousLoginUser.Name)
	}

	// store to currentData.AllUsers
	for _, currentUser := range currentData.LoginUsers {
		currentData.AllUsers = append(currentData.AllUsers, currentUser)
	}
	for _, previousUser := range previousData.AllUsers {
		var flag bool
		for _, currentUser := range currentData.LoginUsers {
			if previousUser.Name == currentUser.Name {
				flag = true
			}
		}
		if !flag {
			currentData.AllUsers = append(currentData.AllUsers, previousUser)
		}
	}

	// send to LINE (PUSH notification) if d.LoginUsers != sharedmem.Domain.LoginUsers
	loggingInUsernameSet := currentLoginUsernameSet.Difference(previousLoginUsernameSet)
	if loggingInUsernameSet.Cardinality() != 0 {
		if err := e.BotPluginSender.SendTextMessage(i18n.T.Sprintf(i18n.MessageUsersLogin, loggingInUsernameSet.ToSlice())); err != nil {
			return err
		}
	}
	loggingOutUsernameSet := previousLoginUsernameSet.Difference(currentLoginUsernameSet)
	if loggingOutUsernameSet.Cardinality() != 0 {
		if err := e.BotPluginSender.SendTextMessage(i18n.T.Sprintf(i18n.MessageUsersLogout, loggingOutUsernameSet.ToSlice())); err != nil {
			return err
		}
	}

	// write to sharedMem
	e.sharedMem.AsyncWriteEntityToSharedMem(currentData)

	return nil
}

func (e *Eventer) getMetricsUsingRCON() (domain.Entity, error) {
	var currentData domain.Entity

	currentLoginUsernames, err := e.rcon.List()
	if err != nil {
		return domain.Entity{}, err
	}
	for _, username := range currentLoginUsernames {
		userData, err := e.rcon.DataGetEntity(username)
		if err != nil {
			e.Logger.Warn(`userData is nil`)
			return domain.Entity{}, err
			// TODO: e.rcon.DataGetEntity の return nil, nil をやめる
		} else if userData == nil {
			return domain.Entity{}, nil
		}
		currentLoginUser := domain.User{
			Name:    username,
			Health:  userData.Health,
			XpLevel: userData.XpLevel,
			Position: domain.Position{
				X: userData.X,
				Y: userData.Y,
				Z: userData.Z,
			},
		}
		currentData.LoginUsers = append(currentData.LoginUsers, currentLoginUser)
	}
	currentData.WhitelistUsernames, err = e.rcon.WhitelistList()
	if err != nil {
		return domain.Entity{}, err
	}
	return currentData, nil
}
