package redis

import (
	"encoding/json"
	"strconv"

	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"

	"github.com/ShotaKitazawa/minecraft-bot/pkg/domain"
)

type SharedMem struct {
	logger        *logrus.Logger
	sendStream    chan<- domain.Entity
	receiveStream <-chan domain.Entity
	Conn          redis.Conn
	redisHostname string
}

func New(logger *logrus.Logger, addr string, port int) (*SharedMem, error) {
	stream := make(chan domain.Entity)
	redisHostname := addr + ":" + strconv.Itoa(port)
	c, err := redis.Dial("tcp", redisHostname)
	if err != nil {
		return nil, err
	}
	m := &SharedMem{
		logger:        logger,
		sendStream:    stream,
		receiveStream: stream,
		Conn:          c,
		redisHostname: redisHostname,
	}
	go m.receiveFromChannelAndWriteSharedMem()
	return m, nil
}

func (m *SharedMem) SyncReadEntityFromSharedMem() (domain.Entity, error) {
	data, err := redis.Bytes(m.Conn.Do("GET", "entity"))
	if err != nil {
		m.logger.Error(err)
		return domain.Entity{}, err
	}
	entity := domain.Entity{}
	if err := json.Unmarshal(data, &entity); err != nil {
		m.logger.Error(err)
		return domain.Entity{}, err
	}
	return entity, nil
}

func (m *SharedMem) AsyncWriteEntityToSharedMem(data domain.Entity) error {
	m.sendStream <- data
	return nil
}

func (m *SharedMem) receiveFromChannelAndWriteSharedMem() error {
	for {
		select {
		case d := <-m.receiveStream:
			data, err := json.Marshal(&d)
			if err != nil {
				return err
			}
			_, err = m.Conn.Do("SET", "entity", data)
			if err != nil {
				m.logger.Error(err)
				return err
			}
		}
	}
	// return nil
}

/* TODO
func (m *SharedMem) reconnect() error {
	c, err := redis.Dial("tcp", m.redisHostname)
	if err != nil {
		return err
	}
	m.Conn = c
	return nil
}
*/
