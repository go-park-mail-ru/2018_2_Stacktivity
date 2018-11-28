package game

import "github.com/prometheus/client_golang/prometheus"

var (
	RoomCountMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rooms_total",
			Help: "Total number of rooms on service. Has label 'type' = {'single', 'mult'}",
		},
		[]string{"type"},
	)

	PlayersPendingRoomMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "players_pending",
			Help: "Total number of players waiting for the room",
		},
		[]string{"type"},
	)

	PlayersLeftGameMetric = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "players_left",
			Help: "Total number of players who left the game unfinished",
		},
	)

	labelTypeSingle = prometheus.Labels{"type": "single"}
	labelTypeMult   = prometheus.Labels{"type": "mult"}
)
