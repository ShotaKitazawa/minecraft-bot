package mock

import (
	"errors"

	"github.com/ShotaKitazawa/minecraft-bot/pkg/domain"
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
func (m *SharedmemMockValid) SyncSubscribeMessage() (domain.Message, error) {
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
func (m *SharedmemMockInvalid) SyncSubscribeMessage() (domain.Message, error) {
	return domain.Message{}, errors.New(``)
}
