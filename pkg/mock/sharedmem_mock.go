package mock

import (
	"errors"

	"github.com/ShotaKitazawa/minecraft-bot/pkg/domain"
	"github.com/ShotaKitazawa/minecraft-bot/pkg/sharedmem"
)

type SharedmemMockValid struct {
	Data *domain.Entity
}

func (m *SharedmemMockValid) SyncReadEntity() (domain.Entity, error) {
	if m.Data == nil {
		return domain.Entity{}, errors.New(``)
	}
	return *m.Data, nil
}

func (m *SharedmemMockValid) AsyncWriteEntity(data domain.Entity) error {
	m.Data = &data
	return nil
}

func (m *SharedmemMockValid) AsyncPublishMessage(data domain.Message) error {
	return nil
}

func (m *SharedmemMockValid) NewSubscriber() (sharedmem.Subscriber, error) {
	return nil, errors.New(``)
}

type SubscriberMockValid struct{}

func (sub *SubscriberMockValid) SyncSubscribeMessage() (domain.Message, error) {
	return domain.Message{
		UserID: MockUserNameValue,
		Msg:    `test`,
	}, nil
}

type SharedmemMockInvalid struct {
}

func (m *SharedmemMockInvalid) SyncReadEntity() (domain.Entity, error) {
	return domain.Entity{}, errors.New(``)
}

func (m *SharedmemMockInvalid) AsyncWriteEntity(data domain.Entity) error {
	return errors.New(``)
}
func (m *SharedmemMockInvalid) AsyncPublishMessage(data domain.Message) error {
	return errors.New(``)
}
func (m *SharedmemMockInvalid) NewSubscriber() (sharedmem.Subscriber, error) {
	return nil, errors.New(``)
}

type SubscriberMockInvalid struct{}

func (m *SubscriberMockInvalid) SyncSubscribeMessage() (domain.Message, error) {
	return domain.Message{}, errors.New(``)
}
