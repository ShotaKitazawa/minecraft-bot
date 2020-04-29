package exporter

import (
	"github.com/ShotaKitazawa/minecraft-bot/pkg/sharedmem"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

type Collector struct {
	descriptors []*prometheus.Desc
	sharedmem   sharedmem.SharedMem
	Logger      *logrus.Logger
}

func New(m sharedmem.SharedMem, l *logrus.Logger) (Collector, error) {
	descriptors := []*prometheus.Desc{
		prometheus.NewDesc(
			"minecraft_user_info",
			"Minecraft Login Users",
			[]string{"username"},
			nil,
		),
		prometheus.NewDesc(
			"minecraft_health_gauge",
			"Minecraft User's Health",
			[]string{"username"},
			nil,
		),
		prometheus.NewDesc(
			"minecraft_pos_x_gauge",
			"Minecraft User's Position of X axis",
			[]string{"username"},
			nil,
		),
		prometheus.NewDesc(
			"minecraft_pos_y_gauge",
			"Minecraft User's Position of Y axis",
			[]string{"username"},
			nil,
		),
		prometheus.NewDesc(
			"minecraft_pos_z_gauge",
			"Minecraft User's Position of Z axis",
			[]string{"username"},
			nil,
		),
		prometheus.NewDesc(
			"minecraft_xp_level_gauge",
			"Minecraft User's Xp Level",
			[]string{"username"},
			nil,
		),
	}

	return Collector{
		descriptors: descriptors,
		sharedmem:   m,
		Logger:      l,
	}, nil
}

func (c Collector) Collect(ch chan<- prometheus.Metric) {
	describeUserInfo := c.descriptors[0]
	describeHealthGauge := c.descriptors[1]
	describePosXGauge := c.descriptors[2]
	describePosYGauge := c.descriptors[3]
	describePosZGauge := c.descriptors[4]
	describeXpLevelGauge := c.descriptors[5]

	data, err := c.sharedmem.SyncReadEntityFromSharedMem()
	if err != nil {
		c.Logger.Warn(err)
		return
	}

	for _, user := range data.AllUsers {

		// check if user log in
		userIsLoggingin := 0
		for _, loginUser := range data.LoginUsers {
			if user.Name == loginUser.Name {
				userIsLoggingin = 1
			}
		}

		// metrics
		ch <- prometheus.MustNewConstMetric(
			describeUserInfo,
			prometheus.GaugeValue,
			float64(userIsLoggingin),
			user.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			describeHealthGauge,
			prometheus.GaugeValue,
			float64(user.Health),
			user.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			describePosXGauge,
			prometheus.GaugeValue,
			float64(user.Position.X),
			user.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			describePosYGauge,
			prometheus.GaugeValue,
			float64(user.Position.Y),
			user.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			describePosZGauge,
			prometheus.GaugeValue,
			float64(user.Position.Z),
			user.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			describeXpLevelGauge,
			prometheus.GaugeValue,
			float64(user.XpLevel),
			user.Name,
		)
	}
}

func (c Collector) Describe(ch chan<- *prometheus.Desc) {
	for _, describe := range c.descriptors {
		ch <- describe
	}
}
