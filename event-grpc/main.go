/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a server for Greeter service.
package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"google.golang.org/grpc"
	"init.bulbul.tv/bulbul-backend/event-grpc/api/heartbeat"
	"init.bulbul.tv/bulbul-backend/event-grpc/config"
	"init.bulbul.tv/bulbul-backend/event-grpc/event"
	grpc_health_v1 "init.bulbul.tv/bulbul-backend/event-grpc/proto"
	"init.bulbul.tv/bulbul-backend/event-grpc/proto/go/bulbulgrpc"
	"init.bulbul.tv/bulbul-backend/event-grpc/rabbit"
	"init.bulbul.tv/bulbul-backend/event-grpc/router"
	"init.bulbul.tv/bulbul-backend/event-grpc/sinker"
)

type server struct {
	bulbulgrpc.UnimplementedEventServiceServer
}

type healthserver struct {
	grpc_health_v1.UnimplementedHealthServer
}

func (s *healthserver) Check(ctx context.Context, request *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	if heartbeat.BEAT {
		return &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}, nil
	} else {
		return &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_NOT_SERVING}, nil
	}
}
func (s *healthserver) Watch(req *grpc_health_v1.HealthCheckRequest, server grpc_health_v1.Health_WatchServer) error {
	if heartbeat.BEAT {
		server.Send(&grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING})
	} else {
		server.Send(&grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_NOT_SERVING})
	}
	return nil
}

func (s *server) Dispatch(ctx context.Context, in *bulbulgrpc.Events) (*bulbulgrpc.Response, error) {
	//log.Printf("Received: %v", in)

	//marshaler := jsonpb.Marshaler{}
	//log.Println(marshaler.MarshalToString(in))

	//md,ok:=metadata.FromIncomingContext(ctx)
	//log.Printf("\nContext: %+v%+v",md,ok)

	if in != nil {
		if err := event.DispatchEvent(ctx, *in); err != nil {
			return &bulbulgrpc.Response{Status: "ERROR"}, err
		}
	}

	return &bulbulgrpc.Response{Status: "SUCCESS"}, nil
}

func main() {

	config := config.New()

	if !strings.EqualFold(config.Env, "PRODUCTION") || strings.EqualFold(os.Getenv("DEBUG"), "1") {
		go func() {
			log.Println(http.ListenAndServe(":6063", nil))
		}()
	}

	go func() {
		lis, err := net.Listen("tcp", ":"+config.Port)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		bulbulgrpc.RegisterEventServiceServer(s, &server{})
		grpc_health_v1.RegisterHealthServer(s, &healthserver{})
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	boolCh := make(chan bool)

	clickhouseEventErrorEx := "grpc_event_error.dx"
	clickhouseEventErrorRK := "error.grpc_clickhouse_event_q"

	rabbitClickhouseEventConsumerCfg := rabbit.RabbitConsumerConfig{
		QueueName:       "grpc_clickhouse_event_q",
		Exchange:        "grpc_event.tx",
		RoutingKey:      "*",
		RetryOnError:    true,
		ErrorExchange:   &clickhouseEventErrorEx,
		ErrorRoutingKey: &clickhouseEventErrorRK,
		ExitCh:          boolCh,
		ConsumerCount:   2,
		ConsumerFunc:    sinker.SinkEventToClickhouse,
	}
	rabbit.Rabbit(config).InitConsumer(rabbitClickhouseEventConsumerCfg)

	clickhouseEventErrorEx1 := "grpc_event_error.dx"
	clickhouseEventErrorRK1 := "error.grpc_clickhouse_event_error_q"

	rabbit.Rabbit(config).InitConsumer(rabbit.RabbitConsumerConfig{
		QueueName:       "grpc_clickhouse_event_error_q",
		Exchange:        "grpc_event_error.dx",
		RoutingKey:      "error.grpc_clickhouse_event_q",
		RetryOnError:    true,
		ErrorExchange:   &clickhouseEventErrorEx1,
		ErrorRoutingKey: &clickhouseEventErrorRK1,
		ExitCh:          boolCh,
		ConsumerCount:   2,
		ConsumerFunc:    sinker.SinkErrorEventToClickhouse,
	})

	clickhouseClickEventErrorEx := "grpc_event.dx"
	clickhouseClickEventErrorRK := "error.clickhouse_click_event_q"

	clickhouseClickEventConsumerCfg := rabbit.RabbitConsumerConfig{
		QueueName:       "clickhouse_click_event_q",
		Exchange:        "grpc_event.dx",
		RoutingKey:      "CLICK",
		RetryOnError:    true,
		ErrorExchange:   &clickhouseClickEventErrorEx,
		ErrorRoutingKey: &clickhouseClickEventErrorRK,
		ExitCh:          boolCh,
		ConsumerCount:   2,
		ConsumerFunc:    sinker.SinkClickEventToClickhouse,
	}
	rabbit.Rabbit(config).InitConsumer(clickhouseClickEventConsumerCfg)

	clickhouseAppInitEventErrorEx := "event_error.dx"
	clickhouseAppInitEventErrorRK := "error.app_init_event_q"

	clickhouseAppInitEventConsumerCfg := rabbit.RabbitConsumerConfig{
		QueueName:       "app_init_event_q",
		Exchange:        "event.dx",
		RoutingKey:      "APP_INIT",
		RetryOnError:    true,
		ErrorExchange:   &clickhouseAppInitEventErrorEx,
		ErrorRoutingKey: &clickhouseAppInitEventErrorRK,
		ExitCh:          boolCh,
		ConsumerCount:   2,
		ConsumerFunc:    sinker.SinkAppInitEventToClickhouse,
	}
	rabbit.Rabbit(config).InitConsumer(clickhouseAppInitEventConsumerCfg)

	clickhouseLaunchReferEventErrorEx := "event_error.dx"
	clickhouseLaunchReferEventErrorRK := "error.launch_refer_event_q"

	clickhouseLaunchReferEventConsumerCfg := rabbit.RabbitConsumerConfig{
		QueueName:       "launch_refer_event_q",
		Exchange:        "event.dx",
		RoutingKey:      "LAUNCH_REFER_EVENT",
		RetryOnError:    true,
		ErrorExchange:   &clickhouseLaunchReferEventErrorEx,
		ErrorRoutingKey: &clickhouseLaunchReferEventErrorRK,
		ExitCh:          boolCh,
		ConsumerCount:   2,
		ConsumerFunc:    sinker.SinkLaunchReferEventToClickhouse,
	}
	rabbit.Rabbit(config).InitConsumer(clickhouseLaunchReferEventConsumerCfg)

	createUserAccountEventErrorEx := "event_error.dx"
	createUserAccountEventErrorRK := "error.create_user_account_q"

	createUserAccountEventConsumerCfg := rabbit.RabbitConsumerConfig{
		QueueName:       "create_user_account_q",
		Exchange:        "identity.dx",
		RoutingKey:      "user.client.account.create",
		RetryOnError:    true,
		ErrorExchange:   &createUserAccountEventErrorEx,
		ErrorRoutingKey: &createUserAccountEventErrorRK,
		ExitCh:          boolCh,
		ConsumerCount:   2,
		ConsumerFunc:    sinker.ProcessUserAccountCreationEventForReferral,
	}
	rabbit.Rabbit(config).InitConsumer(createUserAccountEventConsumerCfg)

	profileCompletionEventErrorEx := "event_error.dx"
	profileCompletionEventErrorRK := "error.create_user_account_q"

	profileCompletionEventConsumerCfg := rabbit.RabbitConsumerConfig{
		QueueName:       "event.account_profile_complete",
		Exchange:        "identity.dx",
		RoutingKey:      "user.client.account.profile.completed",
		RetryOnError:    true,
		ErrorExchange:   &profileCompletionEventErrorEx,
		ErrorRoutingKey: &profileCompletionEventErrorRK,
		ExitCh:          boolCh,
		ConsumerCount:   2,
		ConsumerFunc:    sinker.ProcessUserAccountProfileCompletionEventForReferral,
	}
	rabbit.Rabbit(config).InitConsumer(profileCompletionEventConsumerCfg)

	abAssignmentsErrorEx := "ab.dx"
	abAssignmentsErrorRK := "error.ab_assignments"

	abAssignmentsConsumerCfg := rabbit.RabbitConsumerConfig{
		QueueName:       "ab_assignments",
		Exchange:        "ab.dx",
		RoutingKey:      "ab_assignments",
		RetryOnError:    true,
		ErrorExchange:   &abAssignmentsErrorEx,
		ErrorRoutingKey: &abAssignmentsErrorRK,
		ExitCh:          boolCh,
		ConsumerCount:   2,
		ConsumerFunc:    sinker.SinkABAssignmentsToClickhouse,
	}
	rabbit.Rabbit(config).InitConsumer(abAssignmentsConsumerCfg)

	branchEventErrorEx := "branch_event_error.dx"
	branchEventErrorRK := "error.branch_event_q"

	branchEventConsumerCfg := rabbit.RabbitConsumerConfig{
		QueueName:       "branch_event_q",
		Exchange:        "branch_event.tx",
		RoutingKey:      "*",
		RetryOnError:    true,
		ErrorExchange:   &branchEventErrorEx,
		ErrorRoutingKey: &branchEventErrorRK,
		ExitCh:          boolCh,
		ConsumerCount:   2,
		ConsumerFunc:    sinker.SinkBranchEventToClickhouse,
	}
	rabbit.Rabbit(config).InitConsumer(branchEventConsumerCfg)

	graphyEventErrorEx := "graphy_event_error.dx"
	graphyEventErrorRK := "error.graphy_event_q"

	graphyEventConsumerCfg := rabbit.RabbitConsumerConfig{
		QueueName:       "graphy_event_q",
		Exchange:        "graphy_event.tx",
		RoutingKey:      "*",
		RetryOnError:    true,
		ErrorExchange:   &graphyEventErrorEx,
		ErrorRoutingKey: &graphyEventErrorRK,
		ExitCh:          boolCh,
		ConsumerCount:   2,
		ConsumerFunc:    sinker.SinkGraphyEventToClickhouse,
	}
	rabbit.Rabbit(config).InitConsumer(graphyEventConsumerCfg)

	traceLogEventErrorEx := "trace_log_event_error.dx"
	traceLogEventErrorRK := "error.trace_log_event_q"
	traceLogChan := make(chan interface{}, 10000)

	traceLogEventConsumerCfg := rabbit.RabbitConsumerConfig{
		QueueName:            "trace_log",
		Exchange:             "identity.dx",
		RoutingKey:           "*",
		RetryOnError:         true,
		ErrorExchange:        &traceLogEventErrorEx,
		ErrorRoutingKey:      &traceLogEventErrorRK,
		ExitCh:               boolCh,
		ConsumerCount:        2,
		BufferedConsumerFunc: sinker.BufferTraceLogEvent,
		BufferChan:           traceLogChan,
	}
	rabbit.Rabbit(config).InitConsumer(traceLogEventConsumerCfg)
	go sinker.TraceLogEventsSinker(traceLogChan)

	affiliateOrdersEventErrorEx := "affiliate_orders_event_error.dx"
	affiliateOrdersEventErrorRK := "error.affiliate_orders_event_q"
	affiliateOrdersChan := make(chan interface{}, 10000)

	affiliateOrdersEventConsumerCfg := rabbit.RabbitConsumerConfig{
		QueueName:            "affiliate_orders_event_q",
		Exchange:             "affiliate.dx",
		RoutingKey:           "*",
		RetryOnError:         true,
		ErrorExchange:        &affiliateOrdersEventErrorEx,
		ErrorRoutingKey:      &affiliateOrdersEventErrorRK,
		ExitCh:               boolCh,
		ConsumerCount:        2,
		BufferedConsumerFunc: sinker.BufferAffiliateOrderEvents,
		BufferChan:           affiliateOrdersChan,
	}
	rabbit.Rabbit(config).InitConsumer(affiliateOrdersEventConsumerCfg)
	go sinker.AffiliateOrderEventsSinker(affiliateOrdersChan)

	bigBossVotesEx := "bigboss"
	bigBossVotesErrorRK := "error.big_boss_vote_log_q"
	bigBossVotesChan := make(chan interface{}, 10000)

	bigBossConsumerConfig := rabbit.RabbitConsumerConfig{
		QueueName:            "bigboss_votes_log_q",
		Exchange:             bigBossVotesEx,
		RoutingKey:           "*",
		RetryOnError:         true,
		ErrorExchange:        &bigBossVotesEx,
		ErrorRoutingKey:      &bigBossVotesErrorRK,
		ExitCh:               boolCh,
		ConsumerCount:        2,
		BufferedConsumerFunc: sinker.BufferBigBossVotes,
		BufferChan:           bigBossVotesChan,
	}
	rabbit.Rabbit(config).InitConsumer(bigBossConsumerConfig)
	go sinker.BigBossVotesSinker(bigBossVotesChan)

	postLogEx := "beat.dx"
	postLogRk := "post_log_events"
	postLogExErrorRK := "error.post_log_events_q"
	postLogChan := make(chan interface{}, 10000)

	postLogConsumerConfig := rabbit.RabbitConsumerConfig{
		QueueName:            "post_log_events_q",
		Exchange:             postLogEx,
		RoutingKey:           postLogRk,
		RetryOnError:         true,
		ErrorExchange:        &postLogEx,
		ErrorRoutingKey:      &postLogExErrorRK,
		ExitCh:               boolCh,
		ConsumerCount:        20,
		BufferedConsumerFunc: sinker.BufferPostLogEvents,
		BufferChan:           postLogChan,
	}
	rabbit.Rabbit(config).InitConsumer(postLogConsumerConfig)
	go sinker.PostLogEventsSinker(postLogChan)

	sentimentLogEx := "beat.dx"
	sentimentLogRk := "sentiment_log_events"
	sentimentLogExErrorRK := "error.sentiment_log_events_q"
	sentimentLogChan := make(chan interface{}, 10000)
	sentimentLogConsumerConfig := rabbit.RabbitConsumerConfig{
		QueueName:            "sentiment_log_events_q",
		Exchange:             sentimentLogEx,
		RoutingKey:           sentimentLogRk,
		RetryOnError:         true,
		ErrorExchange:        &sentimentLogEx,
		ErrorRoutingKey:      &sentimentLogExErrorRK,
		ExitCh:               boolCh,
		ConsumerCount:        2,
		BufferedConsumerFunc: sinker.BufferSentimentLogEvents,
		BufferChan:           sentimentLogChan,
	}
	rabbit.Rabbit(config).InitConsumer(sentimentLogConsumerConfig)
	go sinker.SentimentLogEventsSinker(sentimentLogChan)

	postActivityLogEx := "beat.dx"
	postActivityLogRk := "post_activity_log_events"
	postActivityLogExErrorRK := "error.post_activity_log_events_q"
	postActivityLogChan := make(chan interface{}, 10000)
	postActivityLogConsumerConfig := rabbit.RabbitConsumerConfig{
		QueueName:            "post_activity_log_events_q",
		Exchange:             postActivityLogEx,
		RoutingKey:           postActivityLogRk,
		RetryOnError:         true,
		ErrorExchange:        &postActivityLogEx,
		ErrorRoutingKey:      &postActivityLogExErrorRK,
		ExitCh:               boolCh,
		ConsumerCount:        2,
		BufferedConsumerFunc: sinker.BufferPostActivityLogEvents,
		BufferChan:           postActivityLogChan,
	}
	rabbit.Rabbit(config).InitConsumer(postActivityLogConsumerConfig)
	go sinker.PostActivityLogEventsSinker(postActivityLogChan)

	profileLogEx := "beat.dx"
	profileLogRk := "profile_log_events"
	profileLogExErrorRK := "error.profile_log_events_q"
	profileLogChan := make(chan interface{}, 10000)

	profileLogConsumerConfig := rabbit.RabbitConsumerConfig{
		QueueName:            "profile_log_events_q",
		Exchange:             profileLogEx,
		RoutingKey:           profileLogRk,
		RetryOnError:         true,
		ErrorExchange:        &profileLogEx,
		ErrorRoutingKey:      &profileLogExErrorRK,
		ExitCh:               boolCh,
		ConsumerCount:        2,
		BufferedConsumerFunc: sinker.BufferProfileLogEvents,
		BufferChan:           profileLogChan,
	}
	rabbit.Rabbit(config).InitConsumer(profileLogConsumerConfig)
	go sinker.ProfileLogEventsSinker(profileLogChan)

	profileRelationshipLogEx := "beat.dx"
	profileRelationshipLogRk := "profile_relationship_log_events"
	profileRelationshipLogExErrorRK := "error.profile_relationship_log_events"
	profileRelationshipLogChan := make(chan interface{}, 10000)

	profileRelationshipLogConsumerConfig := rabbit.RabbitConsumerConfig{
		QueueName:            "profile_relationship_log_events_q",
		Exchange:             profileRelationshipLogEx,
		RoutingKey:           profileRelationshipLogRk,
		RetryOnError:         true,
		ErrorExchange:        &profileRelationshipLogEx,
		ErrorRoutingKey:      &profileRelationshipLogExErrorRK,
		ExitCh:               boolCh,
		ConsumerCount:        2,
		BufferedConsumerFunc: sinker.BufferProfileRelationshipLogEvents,
		BufferChan:           profileRelationshipLogChan,
	}
	rabbit.Rabbit(config).InitConsumer(profileRelationshipLogConsumerConfig)
	go sinker.ProfileRelationshipLogEventsSinker(profileRelationshipLogChan)

	scrapeRequestLogEx := "beat.dx"
	scrapeRequestLogRk := "scrape_request_log_events"
	scrapeRequestLogExErrorRK := "error.scrape_request_log_events_q"
	scrapeRequestLogChan := make(chan interface{}, 10000)

	scrapeRequestLogConsumerConfig := rabbit.RabbitConsumerConfig{
		QueueName:            "scrape_request_log_events_q",
		Exchange:             scrapeRequestLogEx,
		RoutingKey:           scrapeRequestLogRk,
		RetryOnError:         true,
		ErrorExchange:        &scrapeRequestLogEx,
		ErrorRoutingKey:      &scrapeRequestLogExErrorRK,
		ExitCh:               boolCh,
		ConsumerCount:        2,
		BufferedConsumerFunc: sinker.BufferScrapeRequestLogEvents,
		BufferChan:           scrapeRequestLogChan,
	}
	rabbit.Rabbit(config).InitConsumer(scrapeRequestLogConsumerConfig)
	go sinker.ScrapeRequestLogEventsSinker(scrapeRequestLogChan)

	orderLogEx := "beat.dx"
	orderLogRk := "order_log_events"
	orderLogExErrorRK := "error.order_log_events_q"
	orderLogChan := make(chan interface{}, 10000)

	orderLogConsumerConfig := rabbit.RabbitConsumerConfig{
		QueueName:            "order_log_events_q",
		Exchange:             orderLogEx,
		RoutingKey:           orderLogRk,
		RetryOnError:         true,
		ErrorExchange:        &orderLogRk,
		ErrorRoutingKey:      &orderLogExErrorRK,
		ExitCh:               boolCh,
		ConsumerCount:        2,
		BufferedConsumerFunc: sinker.BufferOrderLogEvents,
		BufferChan:           orderLogChan,
	}
	rabbit.Rabbit(config).InitConsumer(orderLogConsumerConfig)
	go sinker.OrderLogEventsSinker(orderLogChan)

	shopifyEventEx := "shopify_event.dx"
	shopifyEventRk := "shopify_events"
	shopifyExErrorRK := "error.shopify_events_q"
	shopifyEventChan := make(chan interface{}, 10000)

	shopifyEventConsumerConfig := rabbit.RabbitConsumerConfig{
		QueueName:            "shopify_events_q",
		Exchange:             shopifyEventEx,
		RoutingKey:           shopifyEventRk,
		RetryOnError:         true,
		ErrorExchange:        &shopifyEventEx,
		ErrorRoutingKey:      &shopifyExErrorRK,
		ExitCh:               boolCh,
		ConsumerCount:        2,
		BufferedConsumerFunc: sinker.BufferShopifyEvents,
		BufferChan:           shopifyEventChan,
	}
	rabbit.Rabbit(config).InitConsumer(shopifyEventConsumerConfig)
	go sinker.ShopifyEventsSinker(shopifyEventChan)

	webengageEventErrorEx := "webengage_event.dx"
	webengageEventErrorRK := "error.webengage_event_q"
	webengageEventConsumerCfg := rabbit.RabbitConsumerConfig{
		QueueName:       "webengage_event_q",
		Exchange:        "webengage_event.dx",
		RoutingKey:      "EVENT",
		RetryOnError:    true,
		ErrorExchange:   &webengageEventErrorEx,
		ErrorRoutingKey: &webengageEventErrorRK,
		ExitCh:          boolCh,
		ConsumerCount:   3,
		ConsumerFunc:    sinker.SinkWebengageEventToWebengage,
	}
	rabbit.Rabbit(config).InitConsumer(webengageEventConsumerCfg)

	webengageChEventErrorEx := "webengage_event.dx"
	webengageChEventErrorRK := "error.webengage_ch_event_q"
	webengageChEventConsumerCfg := rabbit.RabbitConsumerConfig{
		QueueName:       "webengage_ch_event_q",
		Exchange:        "webengage_event.dx",
		RoutingKey:      "EVENT",
		RetryOnError:    true,
		ErrorExchange:   &webengageChEventErrorEx,
		ErrorRoutingKey: &webengageChEventErrorRK,
		ExitCh:          boolCh,
		ConsumerCount:   5,
		ConsumerFunc:    sinker.SinkWebengageEventToClickhouse,
	}
	rabbit.Rabbit(config).InitConsumer(webengageChEventConsumerCfg)

	webengageUserEventErrorEx := "webengage_event.dx"
	webengageUserEventErrorRK := "error.webengage_user_event_q"
	webengageUserEventConsumerCfg := rabbit.RabbitConsumerConfig{
		QueueName:       "webengage_user_event_q",
		Exchange:        "webengage_event.dx",
		RoutingKey:      "USER_EVENT",
		RetryOnError:    true,
		ErrorExchange:   &webengageUserEventErrorEx,
		ErrorRoutingKey: &webengageUserEventErrorRK,
		ExitCh:          boolCh,
		ConsumerCount:   5,
		ConsumerFunc:    sinker.SinkWebengageUserEventToWebengage,
	}
	rabbit.Rabbit(config).InitConsumer(webengageUserEventConsumerCfg)

	/* TEMP - REMOVE LATER */

	postLogRkBkp := "post_log_events_bkp"
	postLogExErrorRKBkp := "error.post_log_events_q_bkp"
	postLogChanBkp := make(chan interface{}, 10000)

	postLogConsumerConfigBkp := rabbit.RabbitConsumerConfig{
		QueueName:            "post_log_events_q_bkp",
		Exchange:             postLogEx,
		RoutingKey:           postLogRkBkp,
		RetryOnError:         true,
		ErrorExchange:        &postLogEx,
		ErrorRoutingKey:      &postLogExErrorRKBkp,
		ExitCh:               boolCh,
		ConsumerCount:        5,
		BufferedConsumerFunc: sinker.BufferPostLogEvents,
		BufferChan:           postLogChanBkp,
	}
	rabbit.Rabbit(config).InitConsumer(postLogConsumerConfigBkp)
	go sinker.PostLogEventsSinker(postLogChanBkp)

	partnerActivityLogEx := "coffee.dx"
	partnerActivityLogRk := "activity_tracker_rk"
	partnerActivityLogExErrorRK := "error.activity_tracker"
	partnerActivityLogChan := make(chan interface{}, 10000)

	partnerActivityLogConsumerConfig := rabbit.RabbitConsumerConfig{
		QueueName:            "activity_tracker_q",
		Exchange:             partnerActivityLogEx,
		RoutingKey:           partnerActivityLogRk,
		RetryOnError:         true,
		ErrorExchange:        &partnerActivityLogEx,
		ErrorRoutingKey:      &partnerActivityLogExErrorRK,
		ExitCh:               boolCh,
		ConsumerCount:        2,
		BufferedConsumerFunc: sinker.BufferPartnerActivityLogEvents,
		BufferChan:           partnerActivityLogChan,
	}
	rabbit.Rabbit(config).InitConsumer(partnerActivityLogConsumerConfig)
	go sinker.PartnerActivityLogEventsSinker(partnerActivityLogChan)

	err := router.SetupRouter(config).Run(":" + config.GinPort)

	if err != nil {
		log.Fatal("Error starting server")
	}
}
