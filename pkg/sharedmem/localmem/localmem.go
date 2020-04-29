package localmem

import (
	"fmt"
	"sync"

	"github.com/ShotaKitazawa/minecraft-bot/pkg/domain"
	"github.com/sirupsen/logrus"
)

var (
	mu sync.Mutex
)

type SharedMem struct {
	logger        *logrus.Logger
	sendStream    chan<- domain.Entity
	receiveStream <-chan domain.Entity
	sharedMemory  *domain.Entity
}

func New(logger *logrus.Logger) (*SharedMem, error) {
	stream := make(chan domain.Entity)
	m := &SharedMem{
		logger:        logger,
		sendStream:    stream,
		receiveStream: stream,
	}
	go m.receiveFromChannelAndWriteSharedMem()
	return m, nil
}

func (m *SharedMem) SyncReadEntityFromSharedMem() (domain.Entity, error) {
	mu.Lock()
	result := m.sharedMemory
	mu.Unlock()
	if result == nil {
		return domain.Entity{}, fmt.Errorf("no such data")
	}
	return *result, nil
}

func (m *SharedMem) AsyncWriteEntityToSharedMem(data domain.Entity) error {
	m.sendStream <- data
	return nil
}

func (m *SharedMem) receiveFromChannelAndWriteSharedMem() error {
	for {
		select {
		case d := <-m.receiveStream:
			mu.Lock()
			m.sharedMemory = &d
			mu.Unlock()
		}
	}
	// return nil
}
