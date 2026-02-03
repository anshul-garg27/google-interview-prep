package route

import (
	"init.bulbul.tv/bulbul-backend/event-grpc/api/heartbeat"
)

var HeartbeatRoutes = Routes{
	Route{
		"Get heartbeat",
		"GET",
		"/",
		heartbeat.Beat,
	},
	Route{
		"Modify heartbeat",
		"PUT",
		"/",
		heartbeat.ModifyBeat,
	},
}
