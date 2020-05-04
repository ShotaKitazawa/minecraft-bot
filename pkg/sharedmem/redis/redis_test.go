package redis

import (
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/ShotaKitazawa/minecraft-bot/pkg/domain"
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
		m.AsyncWriteEntity(testData)
		time.Sleep(time.Second)
		readData, err := m.SyncReadEntity()
		assert.Nil(t, err)
		assert.Equal(t, testData, readData)
	})

	t.Run(`SyncReadEntity (data is nil)`, func(t *testing.T) {
		m := NewMockRedis(t)
		_, err := m.SyncReadEntity()
		assert.NotNil(t, err)
	})

	/* MEMO: ERR unknown command `PUBLISH` & `SUBSCRIBE`
	t.Run(`publish & subscribe domain.Message`, func(t *testing.T) {
		t.Run(`SyncSubscribeMessage & AsyncPublishMessage`, func(t *testing.T) {
			m := NewMockRedis(t)
			testData := domain.Message{
				UserID: `hoge`,
				Msg:    `fuga`,
			}
			var flag int
			go func() {
				flag = 0
				readData, err := m.SyncSubscribeMessage()
				flag = 1
				assert.Nil(t, err)
				assert.Equal(t, testData, readData)
			}()
			assert.Equal(t, flag, 0)
			m.AsyncPublishMessage(testData)
			time.Sleep(time.Millisecond)
			assert.Equal(t, flag, 1)
		})
	})
	*/
}
