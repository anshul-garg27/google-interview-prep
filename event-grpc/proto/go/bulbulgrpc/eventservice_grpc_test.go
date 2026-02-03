package bulbulgrpc

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"
)

const (
	//address     = "https://stageevents.mera.link" // fixme
	address     = "stage-envoy:80"
	//address     = "localhost:50051"
)

func TestEventServiceClient_Dispatch(t *testing.T) {
	wg := sync.WaitGroup{}
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go fire(wg)
	}
	wg.Wait()
}

func fire(wg sync.WaitGroup) {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := NewEventServiceClient(conn)

	//ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	clientDeadline := time.Now().Add(5000 * time.Millisecond)
	ctx, cancel := context.WithDeadline(context.Background(), clientDeadline)

	defer cancel()
	r, err := c.Dispatch(ctx, &Events{
		Header: &Header{
			DeviceId:   "1111",
			BbDeviceId: 1,
			UserId:     1,
		},
		Events: []*Event{
			{
				Id:        uuid.New().String(),
				Timestamp: strconv.FormatInt(time.Now().Unix(), 10),
				EventOf: &Event_TestEvent{
					TestEvent: &TestEvent{
						TestParam1: uuid.New().String(),
						TestParam2: int64(rand.Int()),
					},
				},
			},
			{
				Id:        uuid.New().String(),
				Timestamp: strconv.FormatInt(time.Now().Unix(), 10),
				EventOf: &Event_TestEvent{
					TestEvent: &TestEvent{
						TestParam1: uuid.New().String(),
						TestParam2: int64(rand.Int()),
					},
				},
			},
			{
				Id:        uuid.New().String(),
				Timestamp: strconv.FormatInt(time.Now().Unix(), 10),
				EventOf: &Event_TestEvent{
					TestEvent: &TestEvent{
						TestParam1: uuid.New().String(),
						TestParam2: int64(rand.Int()),
					},
				},
			},
			{
				Id:        uuid.New().String(),
				Timestamp: strconv.FormatInt(time.Now().Unix(), 10),
				EventOf: &Event_TestEvent{
					TestEvent: &TestEvent{
						TestParam1: uuid.New().String(),
						TestParam2: int64(rand.Int()),
					},
				},
			},
			{
				Id:        uuid.New().String(),
				Timestamp: strconv.FormatInt(time.Now().Unix(), 10),
				EventOf: &Event_TestEvent{
					TestEvent: &TestEvent{
						TestParam1: uuid.New().String(),
						TestParam2: int64(rand.Int()),
					},
				},
			},
		},
	})
	if err != nil {
		log.Printf("could not dispatch event: %v", err)
	}
	log.Printf("Dispatched event: %s", r.GetStatus())
	wg.Done()
}
