package localmem

import (
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/ShotaKitazawa/minecraft-bot/pkg/domain"
)

var (
	loggerForTest = logrus.New()
)

func NewLocalMemForTest() *SharedMem {
	m, err := New(loggerForTest)
	if err != nil {
		panic(err)
	}
	return m
}

func TestLocalMem(t *testing.T) {
	t.Run(`store & read domain.Entity`, func(t *testing.T) {

		t.Run(`AsyncWriteEntity && SyncReadEntity`, func(t *testing.T) {
			m := NewLocalMemForTest()
			testData := domain.Entity{AllUsers: []domain.User{{Name: `test`}}}
			m.AsyncWriteEntityToSharedMem(testData)
			time.Sleep(time.Millisecond)
			readData, err := m.SyncReadEntityFromSharedMem()
			assert.Nil(t, err)
			assert.Equal(t, testData, readData)
		})

		t.Run(`SyncReadEntity (data is nil)`, func(t *testing.T) {
			m := NewLocalMemForTest()
			_, err := m.SyncReadEntityFromSharedMem()
			assert.NotNil(t, err)
		})

	})
	t.Run(`publish & subscribe domain.Message`, func(t *testing.T) {
		t.Run(`AsyncPublishMessage && SyncSubscribeMessage`, func(t *testing.T) {
			m := NewLocalMemForTest()
			testData := domain.Message{
				UserID: `hoge`,
				Msg:    `fuga`,
			}
			m.AsyncPublishMessage(testData)
			time.Sleep(time.Millisecond)
			readData, err := m.SyncSubscribeMessage()
			assert.Nil(t, err)
			assert.Equal(t, testData, readData)
		})
		t.Run(`SyncSubscribeMessage & AsyncPublishMessage`, func(t *testing.T) {
			m := NewLocalMemForTest()
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
}
