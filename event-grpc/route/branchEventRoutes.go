package route

import (
	"init.bulbul.tv/bulbul-backend/event-grpc/api/branchevent"
)

var BranchEventRoutes = Routes{
	Route{
		"Event",
		"POST",
		"/",
		branchevent.HandleBranchEvent,
	},
}
