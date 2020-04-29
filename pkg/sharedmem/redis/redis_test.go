package redis

import (
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ShotaKitazawa/minecraft-bot/pkg/domain"
	"github.com/alicebob/miniredis"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var (
	loggerForTest = logrus.New()
)

func NewMockRedis(t *testing.T) *SharedMem {
	t.Helper()

	// redisサーバを作る
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	addr := strings.Split(s.Addr(), ":")
	port, err := strconv.Atoi(addr[1])
	if err != nil {
		panic(err)
	}
	m, err := New(loggerForTest, addr[0], port)
	if err != nil {
		panic(err)
	}
	return m
}

func TestRedis(t *testing.T) {

	t.Run(`AsyncWriteEntity && SyncReadEntity`, func(t *testing.T) {
		m := NewMockRedis(t)
		testData := domain.Entity{AllUsers: []domain.User{{Name: `test`}}}
		m.AsyncWriteEntityToSharedMem(testData)
		time.Sleep(time.Second)
		readData, err := m.SyncReadEntityFromSharedMem()
		assert.Nil(t, err)
		assert.Equal(t, testData, readData)
	})

	t.Run(`SyncReadEntity (data is nil)`, func(t *testing.T) {
		m := NewMockRedis(t)
		_, err := m.SyncReadEntityFromSharedMem()
		assert.NotNil(t, err)
	})

}
