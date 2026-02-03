package route

import (
	"init.bulbul.tv/bulbul-backend/event-grpc/api/webengageevent"
	"init.bulbul.tv/bulbul-backend/event-grpc/api/webengageuserevent"
)

var WebengageEventRoutes = Routes{
	Route{
		"Event",
		"POST",
		"/",
		webengageevent.HandleWebengageEvent,
	}, Route{
		"Event",
		"POST",
		"/users",
		webengageuserevent.HandleWebengageUserEvent,
	},
}
