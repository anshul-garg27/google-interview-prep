package event

import (
	"github.com/joho/godotenv"
	"init.bulbul.tv/bulbul-backend/event-grpc/config"
	"init.bulbul.tv/bulbul-backend/event-grpc/proto/go/bulbulgrpc"
	"init.bulbul.tv/bulbul-backend/event-grpc/util"
	"testing"
	"time"
)

func TestDispatchEvent(t *testing.T) {
	godotenv.Load(".env")

	gCtx := util.GatewayContextFromGinContext(nil, config.New())

	ste :=	bulbulgrpc.Event_StreamEnterEvent{
		StreamEnterEvent: &bulbulgrpc.StreamEnterEvent{
			StreamId:  123,
		},
	}
	event := bulbulgrpc.Event{
		Id:        "fvdvd",
		Timestamp: time.Now().String(),
		EventOf: &ste,
	}

	events := bulbulgrpc.Events{
		Header: &bulbulgrpc.Header{
			BbDeviceId: 123,
		},
		Events: []*bulbulgrpc.Event{&event},
	}

	DispatchEvent(gCtx, events)

	time.Sleep(time.Hour)

}