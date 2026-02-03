package route

import (
	"init.bulbul.tv/bulbul-backend/event-grpc/api/shopifyevent"
)

var ShopifyEventRoutes = Routes{
	Route{
		"Event",
		"POST",
		"/",
		shopifyevent.HandleShopifyEvent,
	},
}
