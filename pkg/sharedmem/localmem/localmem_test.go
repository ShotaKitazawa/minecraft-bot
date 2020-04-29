package localmem

import (
	"testing"
	"time"

	"github.com/ShotaKitazawa/minecraft-bot/pkg/domain"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
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

	t.Run(`AsyncWriteEntity && SyncReadEntity`, func(t *testing.T) {
		m := NewLocalMemForTest()
		testData := domain.Entity{AllUsers: []domain.User{{Name: `test`}}}
		m.AsyncWriteEntityToSharedMem(testData)
		time.Sleep(time.Second)
		readData, err := m.SyncReadEntityFromSharedMem()
		assert.Nil(t, err)
		assert.Equal(t, testData, readData)
	})

	t.Run(`SyncReadEntity (data is nil)`, func(t *testing.T) {
		m := NewLocalMemForTest()
		_, err := m.SyncReadEntityFromSharedMem()
		assert.NotNil(t, err)
	})

}
