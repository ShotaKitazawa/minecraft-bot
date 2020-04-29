package mock

import (
	"errors"

	"github.com/ShotaKitazawa/minecraft-bot/pkg/rcon"
)

type RconClientMockValid struct {
	LoginUsernames       []string
	WhitelistedUsernames []string
}

func (r *RconClientMockValid) List() ([]string, error) {
	return r.LoginUsernames, nil
}
func (r *RconClientMockValid) WhitelistAdd(username string) error {
	r.WhitelistedUsernames = append(r.WhitelistedUsernames, username)
	return nil
}
func (r *RconClientMockValid) WhitelistRemove(username string) error {
	var matched bool
	for idx, whitelistedUsername := range r.WhitelistedUsernames {
		if username == whitelistedUsername {
			matched = true
			if idx == len(r.WhitelistedUsernames)-1 {
				r.WhitelistedUsernames = r.WhitelistedUsernames[0:idx]
			} else {
				r.WhitelistedUsernames = append(r.WhitelistedUsernames[0:idx], r.WhitelistedUsernames[idx+1:]...)
			}
		}
	}
	if !matched {
		return errors.New(``)
	}
	return nil
}
func (r *RconClientMockValid) WhitelistList() ([]string, error) {
	return r.WhitelistedUsernames, nil
}
func (r *RconClientMockValid) DataGetEntity(string) (*rcon.User, error) {
	return &rcon.User{}, nil
}
func (r *RconClientMockValid) Title(string) ([]string, error) {
	return r.LoginUsernames, nil
}

type RconClientMockInvalid struct {
}

func (r *RconClientMockInvalid) List() ([]string, error)                  { return nil, errors.New(``) }
func (r *RconClientMockInvalid) WhitelistAdd(string) error                { return errors.New(``) }
func (r *RconClientMockInvalid) WhitelistRemove(string) error             { return errors.New(``) }
func (r *RconClientMockInvalid) WhitelistList() ([]string, error)         { return nil, errors.New(``) }
func (r *RconClientMockInvalid) DataGetEntity(string) (*rcon.User, error) { return nil, errors.New(``) }
func (r *RconClientMockInvalid) Title(string) ([]string, error)           { return nil, errors.New(``) }
