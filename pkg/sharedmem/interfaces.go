package sharedmem

import "github.com/ShotaKitazawa/minecraft-bot/pkg/domain"

type SharedMem interface {
	SyncReadEntity() (domain.Entity, error)
	AsyncWriteEntity(data domain.Entity) error
	AsyncPublishMessage(data domain.Message) error
	NewSubscriber() (Subscriber, error)
}

type Subscriber interface {
	SyncSubscribeMessage() (domain.Message, error)
}
