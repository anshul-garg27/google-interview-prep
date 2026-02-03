package route

import (
	"init.bulbul.tv/bulbul-backend/saas-gateway/api/heartbeat"
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
