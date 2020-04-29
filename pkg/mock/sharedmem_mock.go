package mock

import (
	"errors"

	"github.com/ShotaKitazawa/minecraft-bot/pkg/domain"
)

type SharedmemMockValid struct {
	Data *domain.Entity
}

func (m *SharedmemMockValid) SyncReadEntityFromSharedMem() (domain.Entity, error) {
	if m.Data == nil {
		return domain.Entity{}, errors.New(``)
	}
	return *m.Data, nil
}

func (m *SharedmemMockValid) AsyncWriteEntityToSharedMem(data domain.Entity) error {
	m.Data = &data
	return nil
}

type SharedmemMockInvalid struct {
}

func (m *SharedmemMockInvalid) SyncReadEntityFromSharedMem() (domain.Entity, error) {
	return domain.Entity{}, errors.New(``)
}

func (m *SharedmemMockInvalid) AsyncWriteEntityToSharedMem(data domain.Entity) error {
	return errors.New(``)
}
