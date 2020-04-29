package exporter

import (
	"fmt"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/ShotaKitazawa/minecraft-bot/pkg/domain"
	"github.com/ShotaKitazawa/minecraft-bot/pkg/mock"
	"github.com/ShotaKitazawa/minecraft-bot/pkg/sharedmem"
)

var (
	loggerForTest = logrus.New()
	userForTest   = domain.User{
		Name:    mock.MockUserNameValue,
		Health:  mock.MockUserHealthValue,
		XpLevel: mock.MockUserXpValue,
		Position: domain.Position{
			X: mock.MockUserPosXValue,
			Y: mock.MockUserPosYValue,
			Z: mock.MockUserPosZValue,
		},
	}
)

func NewCollectorForTest(m sharedmem.SharedMem) Collector {
	c, err := New(m, loggerForTest)
	if err != nil {
		panic(err)
	}

	return c
}

func NewStreamMetricForTest() (chan<- prometheus.Metric, <-chan prometheus.Metric) {
	stream := make(chan prometheus.Metric)
	var receiveChan chan<- prometheus.Metric
	receiveChan = stream
	var sendChan <-chan prometheus.Metric
	sendChan = stream
	return receiveChan, sendChan
}

func NewStreamDescForTest() (chan<- *prometheus.Desc, <-chan *prometheus.Desc) {
	stream := make(chan *prometheus.Desc)
	var receiveChan chan<- *prometheus.Desc
	receiveChan = stream
	var sendChan <-chan *prometheus.Desc
	sendChan = stream
	return receiveChan, sendChan
}

func TestExporter(t *testing.T) {

	t.Run(`normal`, func(t *testing.T) {

		t.Run(`Collect()`, func(t *testing.T) {

			t.Run(`exist data in SharedMem (AllUsers:1,LoginUser:0)`, func(t *testing.T) {
				c := NewCollectorForTest(&mock.SharedmemMockValid{Data: &domain.Entity{
					AllUsers: []domain.User{userForTest},
				}})
				receiveChan, sendChan := NewStreamMetricForTest()
				go c.Collect(receiveChan)

				var metric prometheus.Metric
				for i := 0; i < len(c.descriptors); i++ {
					metric = <-sendChan
					assert.Equal(t, metric.Desc(), c.descriptors[i])
					// TODO: conpare metric.value and expected-value
				}
				// check to no value get
				go func() {
					<-sendChan
					t.Error(fmt.Errorf(`sharedmem has no value, but channel get value`))
				}()
				time.Sleep(time.Second)
			})

			t.Run(`exist data in SharedMem (AllUsers:1,LoginUser:1)`, func(t *testing.T) {
				c := NewCollectorForTest(&mock.SharedmemMockValid{Data: &domain.Entity{
					AllUsers:   []domain.User{userForTest},
					LoginUsers: []domain.User{userForTest},
				}})
				receiveChan, sendChan := NewStreamMetricForTest()

				go c.Collect(receiveChan)

				var metric prometheus.Metric
				for i := 0; i < len(c.descriptors); i++ {
					metric = <-sendChan
					assert.Equal(t, metric.Desc(), c.descriptors[i])
					// TODO: conpare metric.value and expected-value
				}
				// check to no value get
				go func() {
					<-sendChan
					t.Error(fmt.Errorf(`sharedmem has no value, but channel get value`))
				}()
				time.Sleep(time.Second)
			})
			t.Run(`no exist data in SharedMem`, func(t *testing.T) {
				c := NewCollectorForTest(&mock.SharedmemMockValid{})
				receiveChan, sendChan := NewStreamMetricForTest()
				c.Collect(receiveChan)

				// check to no value get
				go func() {
					<-sendChan
					t.Error(fmt.Errorf(`sharedmem has no value, but channel get value`))
				}()
				time.Sleep(time.Second)
			})
		})

		t.Run(`Describe()`, func(t *testing.T) {

			t.Run(`exist data in SharedMem (AllUsers:1,LoginUser:1)`, func(t *testing.T) {
				c := NewCollectorForTest(&mock.SharedmemMockValid{Data: &domain.Entity{
					AllUsers:   []domain.User{userForTest},
					LoginUsers: []domain.User{userForTest},
				}})
				receiveChan, sendChan := NewStreamDescForTest()

				go c.Describe(receiveChan)

				var desc *prometheus.Desc
				for i := 0; i < len(c.descriptors); i++ {
					desc = <-sendChan
					assert.Equal(t, desc, c.descriptors[i])
				}

				// check to no value get
				go func() {
					<-sendChan
					t.Error(fmt.Errorf(`sharedmem has no value, but channel get value`))
				}()
				time.Sleep(time.Second)
			})

		})
	})

	t.Run(`abnormal`, func(t *testing.T) {
		c := NewCollectorForTest(&mock.SharedmemMockInvalid{})
		receiveChan, sendChan := NewStreamMetricForTest()
		c.Collect(receiveChan)

		// check to no value get
		go func() {
			<-sendChan
			t.Error(fmt.Errorf(`sharedmem has no value, but channel get value`))
		}()
		time.Sleep(time.Second)
	})

}
