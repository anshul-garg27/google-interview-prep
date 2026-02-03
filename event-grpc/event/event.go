package event

import (
	"context"
	"errors"
	"init.bulbul.tv/bulbul-backend/event-grpc/config"
	"init.bulbul.tv/bulbul-backend/event-grpc/eventworker"
	"init.bulbul.tv/bulbul-backend/event-grpc/proto/go/bulbulgrpc"
	"init.bulbul.tv/bulbul-backend/event-grpc/util"
	"time"
)

func DispatchEvent(ctx context.Context, events bulbulgrpc.Events) error {
	gCtx := util.GatewayContextFromGinContext(nil, config.New())
	//md,_:=metadata.FromIncomingContext(ctx)

	incompleteEvent := false
	processedEvent := false

	if events.Header == nil || events.Header.DeviceId == "" {
		//log.Printf("Incomplete event with context: %v %v", events, md)
		incompleteEvent = true
	}

	//log.Printf("Complete event with context: %v %v", events, md)

	if eventChannel := eventworker.GetChannel(gCtx.Config); eventChannel == nil {
		gCtx.Logger.Error().Msgf("Event channel nil: %v", events)
	} else {
		select {
		case <-time.After(1 * time.Second):
			gCtx.Logger.Log().Msgf("Error publishing event: %v", events)
		case eventChannel <- events:
			processedEvent = true
			//return nil
		}
	}

	if incompleteEvent {
		return errors.New("Incomplete event")
	} else if processedEvent {
		return nil
	} else {
		return errors.New("Error publishing event")
	}
}

