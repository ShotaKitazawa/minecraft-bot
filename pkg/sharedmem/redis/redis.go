package redis

import (
	"encoding/json"
	"strconv"

	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"

	"github.com/ShotaKitazawa/minecraft-bot/pkg/domain"
)

var (
	pubsubMsgChannelName = `message`
)

type SharedMem struct {
	logger               *logrus.Logger
	sendStreamEntity     chan<- domain.Entity
	receiveStreamEntity  <-chan domain.Entity
	sendStreamMessage    chan<- domain.Message
	receiveStreamMessage <-chan domain.Message
	Conn                 redis.Conn
	PubMsgConn           redis.PubSubConn
	SubMsgConn           redis.PubSubConn
	redisHostname        string
}

func New(logger *logrus.Logger, addr string, port int) (*SharedMem, error) {
	streamEntity := make(chan domain.Entity)
	streamQueue := make(chan domain.Message)
	redisHostname := addr + ":" + strconv.Itoa(port)
	c, err := redis.Dial("tcp", redisHostname)
	if err != nil {
		return nil, err
	}
	pubconn := redis.PubSubConn{Conn: c}
	clientForSubscribe, err := redis.Dial("tcp", redisHostname)
	if err != nil {
		return nil, err
	}
	subconn := redis.PubSubConn{Conn: clientForSubscribe}
	subconn.Subscribe(pubsubMsgChannelName)
	m := &SharedMem{
		logger:               logger,
		sendStreamEntity:     streamEntity,
		receiveStreamEntity:  streamEntity,
		sendStreamMessage:    streamQueue,
		receiveStreamMessage: streamQueue,
		Conn:                 c,
		PubMsgConn:           pubconn,
		SubMsgConn:           subconn,
		redisHostname:        redisHostname,
	}
	go m.receiveFromChannelAndWriteSharedMem()
	return m, nil
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

func (m *SharedMem) receiveFromChannelAndWriteSharedMem() error {
	for {
		select {
		case d := <-m.receiveStreamEntity:
			data, err := json.Marshal(&d)
			if err != nil {
				return err
			}
			_, err = m.Conn.Do("SET", "entity", data)
			if err != nil {
				m.logger.Error(err)
				return err
			}
		case d := <-m.receiveStreamMessage:
			data, err := json.Marshal(&d)
			if err != nil {
				return err
			}
			_, err = m.Conn.Do("PUBLISH", pubsubMsgChannelName, data)
			if err != nil {
				m.logger.Error(err)
				return err
			}
		}
	}
	// return nil
}

func (m *SharedMem) SyncReadEntity() (domain.Entity, error) {
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

func (m *SharedMem) AsyncWriteEntity(data domain.Entity) error {
	m.sendStreamEntity <- data
	return nil
}

func (m *SharedMem) AsyncPublishMessage(data domain.Message) error {
	m.sendStreamMessage <- data
	return nil
}

func (m *SharedMem) SyncSubscribeMessage() (domain.Message, error) {
	message := domain.Message{}
	switch v := m.SubMsgConn.Receive().(type) {
	case redis.Message:
		if err := json.Unmarshal(v.Data, &message); err != nil {
			m.logger.Error(err)
			return domain.Message{}, err
		}
		// case redis.Subscription:
		// 	break
		// case error:
		// 	return
	}

	return message, nil
}
