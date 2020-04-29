package sharedmem

import "github.com/ShotaKitazawa/minecraft-bot/pkg/domain"

type SharedMem interface {
	SyncReadEntityFromSharedMem() (domain.Entity, error)
	AsyncWriteEntityToSharedMem(data domain.Entity) error
}
