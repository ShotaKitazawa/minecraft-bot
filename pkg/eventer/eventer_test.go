package eventer

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/ShotaKitazawa/minecraft-bot/pkg/domain"
	"github.com/ShotaKitazawa/minecraft-bot/pkg/mock"
	"github.com/ShotaKitazawa/minecraft-bot/pkg/rcon"
	"github.com/ShotaKitazawa/minecraft-bot/pkg/sharedmem"
)

var (
	loggerForTest            = logrus.New()
	minecraftHostnameForTest = `test`
)

func newForTest(m sharedmem.SharedMem, rcon rcon.RconClient) *Eventer {
	return &Eventer{
		MinecraftHostname: minecraftHostnameForTest,
		sharedMem:         m,
		rcon:              rcon,
		Logger:            loggerForTest,
	}
}

func TestEventer(t *testing.T) {

	t.Run(`New()`, func(t *testing.T) {
		sharedMemValid := &mock.SharedmemMockValid{}
		rconValid := &mock.RconClientMockValid{}
		e, err := New(minecraftHostnameForTest, sharedMemValid, rconValid, loggerForTest)
		assert.Nil(t, err)
		assert.Equal(t, e, &Eventer{
			MinecraftHostname: minecraftHostnameForTest,
			sharedMem:         sharedMemValid,
			rcon:              rconValid,
			Logger:            loggerForTest,
		})
	})
	t.Run(`job()`, func(t *testing.T) {
		t.Run(`normal`, func(t *testing.T) {
			t.Run(`numOfAllLoginUser:0,numOfCurrentLoginUser:0,numOfPreviousLoginUser:0`, func(t *testing.T) {
				e := newForTest(
					&mock.SharedmemMockValid{},
					&mock.RconClientMockValid{},
				)
				assert.Nil(t, e.job())
			})
			t.Run(`numOfAllLoginUser:1,numOfCurrentLoginUser:0,numOfPreviousLoginUser:0`, func(t *testing.T) {
				e := newForTest(
					&mock.SharedmemMockValid{Data: &domain.Entity{
						AllUsers: []domain.User{{Name: `test`}},
					}},
					&mock.RconClientMockValid{},
				)
				assert.Nil(t, e.job())
			})
			t.Run(`numOfAllLoginUser:1,numOfCurrentLoginUser:1,numOfPreviousLoginUser:0`, func(t *testing.T) {
				e := newForTest(
					&mock.SharedmemMockValid{Data: &domain.Entity{
						AllUsers: []domain.User{{Name: `test`}},
					}},
					&mock.RconClientMockValid{LoginUsernames: []string{`test`}},
				)
				assert.Nil(t, e.job())
			})
			t.Run(`numOfAllLoginUser:1,numOfCurrentLoginUser:0,numOfPreviousLoginUser:1`, func(t *testing.T) {
				e := newForTest(
					&mock.SharedmemMockValid{Data: &domain.Entity{
						AllUsers:   []domain.User{{Name: `test`}},
						LoginUsers: []domain.User{{Name: `test`}},
					}},
					&mock.RconClientMockValid{},
				)
				assert.Nil(t, e.job())
			})
			t.Run(`numOfAllLoginUser:1,numOfCurrentLoginUser:1,numOfPreviousLoginUser:1`, func(t *testing.T) {
				e := newForTest(
					&mock.SharedmemMockValid{Data: &domain.Entity{
						AllUsers:   []domain.User{{Name: `test`}},
						LoginUsers: []domain.User{{Name: `test`}},
					}},
					&mock.RconClientMockValid{LoginUsernames: []string{`test`}},
				)
				assert.Nil(t, e.job())
			})
		})
		t.Run(`abnormal(SharedMem)`, func(t *testing.T) {
			e := newForTest(
				&mock.SharedmemMockInvalid{},
				&mock.RconClientMockValid{},
			)
			assert.NotNil(t, e.job())
		})
		t.Run(`abnormal(RconClient)`, func(t *testing.T) {
			e := newForTest(
				&mock.SharedmemMockValid{},
				&mock.RconClientMockInvalid{},
			)
			assert.NotNil(t, e.job())
		})
	})
}
