package route

import "init.bulbul.tv/bulbul-backend/event-grpc/api/vidoolyevent"

var VidoolyEventRoutes = Routes{
	Route{
		"Event",
		"POST",
		"/",
		vidoolyevent.HandleVidoolyEvent,
	},
}
