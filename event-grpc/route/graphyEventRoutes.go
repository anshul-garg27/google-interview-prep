package route

import (
	"init.bulbul.tv/bulbul-backend/event-grpc/api/graphyevent"
)

var GraphyEventRoutes = Routes{
	Route{
		"Event",
		"POST",
		"/",
		graphyevent.HandleGraphyEvent,
	},
}
