package localmem

import (
	"fmt"
	"sync"

	"github.com/ShotaKitazawa/minecraft-bot/pkg/domain"
	"github.com/ShotaKitazawa/minecraft-bot/pkg/sharedmem"
	"github.com/sirupsen/logrus"
)

var (
	mu sync.Mutex
)

type SharedMem struct {
	logger               *logrus.Logger
	sendStreamEntity     chan<- domain.Entity
	receiveStreamEntity  <-chan domain.Entity
	sendStreamMessage    chan<- domain.Message
	receiveStreamMessage <-chan domain.Message
	publishMessage       chan<- domain.Message
	subscribeMessage     <-chan domain.Message
	sharedMemoryEntity   *domain.Entity
}

func New(logger *logrus.Logger) (*SharedMem, error) {
	streamEntity := make(chan domain.Entity)
	streamQueue := make(chan domain.Message)
	streamPubSubMsg := make(chan domain.Message)
	m := &SharedMem{
		logger:               logger,
		sendStreamEntity:     streamEntity,
		receiveStreamEntity:  streamEntity,
		sendStreamMessage:    streamQueue,
		receiveStreamMessage: streamQueue,
		publishMessage:       streamPubSubMsg,
		subscribeMessage:     streamPubSubMsg,
	}
	go m.receiveFromChannelAndWriteSharedMem()
	return m, nil
}

func (m *SharedMem) receiveFromChannelAndWriteSharedMem() error {
	for {
		select {
		case d := <-m.receiveStreamEntity:
			mu.Lock()
			m.sharedMemoryEntity = &d
			mu.Unlock()
		case d := <-m.receiveStreamMessage:
			go func() {
				m.sendStreamMessage <- d
			}()
		}
	}
	// return nil
}

func (m *SharedMem) SyncReadEntity() (domain.Entity, error) {
	mu.Lock()
	result := m.sharedMemoryEntity
	mu.Unlock()
	if result == nil {
		return domain.Entity{}, fmt.Errorf("no such data")
	}
	return *result, nil
}

func (m *SharedMem) AsyncWriteEntity(data domain.Entity) error {
	m.sendStreamEntity <- data
	return nil
}

func (m *SharedMem) AsyncPublishMessage(data domain.Message) error {
	m.sendStreamMessage <- data
	return nil
}

// subscriber instance

type Subscriber struct {
	logger               *logrus.Logger
	receiveStreamMessage <-chan domain.Message
}

func (m *SharedMem) NewSubscriber() (sharedmem.Subscriber, error) {
	return &Subscriber{
		logger:               m.logger,
		receiveStreamMessage: m.receiveStreamMessage,
	}, nil
}

func (sub *Subscriber) SyncSubscribeMessage() (domain.Message, error) {
	result := <-sub.receiveStreamMessage
	return result, nil
}
