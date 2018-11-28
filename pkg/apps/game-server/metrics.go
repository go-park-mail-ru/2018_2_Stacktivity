package game_server

import "github.com/prometheus/client_golang/prometheus"

var (
	singleplayerHitsMetric = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "game_api_singleplayer_total_requests",
			Help: "Total number of requests on game api singleplayer with status codes",
		},
		[]string{"code", "method"},
	)

	multiplayerHitsMetric = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "game_api_multiplayer_total_requests",
			Help: "Total number of requests on game api multiplayer with status codes",
		},
		[]string{"code", "method"},
	)

	gameHitsMetric = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "game_api_rooms_total_requests",
			Help: "Total number of requests on game api game rooms with status codes",
		},
		[]string{"code", "method"},
	)
)
