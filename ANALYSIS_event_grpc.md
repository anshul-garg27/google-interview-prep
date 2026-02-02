# ULTRA-DEEP ANALYSIS: EVENT-GRPC PROJECT

## PROJECT OVERVIEW

| Attribute | Value |
|-----------|-------|
| **Project Name** | Event-gRPC |
| **Purpose** | High-throughput event ingestion & distribution system for real-time analytics |
| **Architecture** | gRPC Server + RabbitMQ Message Broker + Multi-Database Sinks |
| **Language** | Go 1.14 |
| **Project Size** | 11MB, 10,000+ LOC |
| **Ports** | 8017 (gRPC), 8019 (HTTP/Gin) |

---

## 1. COMPLETE DIRECTORY STRUCTURE

```
event-grpc/
├── main.go (573 lines)              # Entry point with 26 consumer configurations
├── go.mod / go.sum                   # 30+ direct dependencies
├── .gitlab-ci.yml                    # CI/CD pipeline
├── .env.*                            # Environment configs (local, stage, production)
│
├── proto/                            # Protocol Buffers
│   ├── eventservice.proto (688 lines) # 60+ event types defined
│   ├── healthcheck.proto             # gRPC health service
│   └── go/bulbulgrpc/                # Generated Go code
│
├── eventworker/                      # Main gRPC event worker pool
├── brancheventworker/                # Branch.io events
├── vidoolyeventworker/               # Vidooly analytics events
├── webengageeventworker/             # WebEngage marketing events
├── graphyeventworker/                # Graphy platform events
├── shopifyeventworker/               # Shopify e-commerce events
│
├── sinker/ (20+ files)               # Event persistence layer
│   ├── eventsinker.go (38KB)         # Main ClickHouse sink
│   ├── clickeventsinker.go (78KB)    # Click tracking (600+ lines)
│   ├── brancheventsinker.go          # Branch attribution
│   ├── webengageeventsinker.go       # WebEngage API sync
│   └── parser/                       # Event parsers
│
├── model/ (19 files)                 # Database models
│   ├── event.go                      # Core event model
│   ├── branchevent.go (49 fields)    # Attribution tracking
│   └── clickevent.go                 # Click analytics
│
├── rabbit/rabbit.go (293 lines)      # RabbitMQ client
├── clickhouse/clickhouse.go          # ClickHouse connection pool
├── pg/pg.go                          # PostgreSQL connection pool
├── cache/                            # Redis + Ristretto caching
│
├── config/config.go (283 lines)      # 40+ configuration fields
├── router/router.go                  # Gin HTTP setup
├── middleware/                       # Request logging, auth, context
├── context/context.go (199 lines)    # Gateway context management
├── client/                           # External API clients
└── scripts/start.sh                  # Graceful deployment script
```

---

## 2. ARCHITECTURE DEEP DIVE

### System Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────────────┐
│                        CLIENT APPLICATIONS                               │
│              (Mobile Apps, Web Apps, Backend Services)                  │
└────────────────────────────────┬────────────────────────────────────────┘
                                 │ gRPC (Protobuf)
                                 ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                     EVENT-GRPC SERVICE                                   │
├─────────────────────────────────────────────────────────────────────────┤
│  ┌─────────────────────┐    ┌─────────────────────────────────────┐    │
│  │  gRPC Server        │    │  HTTP Server (Gin)                  │    │
│  │  Port: 8017         │    │  Port: 8019                         │    │
│  │  - Dispatch RPC     │    │  - /heartbeat (health)              │    │
│  │  - Health Check     │    │  - /metrics (Prometheus)            │    │
│  └──────────┬──────────┘    │  - /vidooly/event                   │    │
│             │               │  - /branch/event                    │    │
│             │               │  - /webengage/event                 │    │
│             │               │  - /shopify/event                   │    │
│             │               └─────────────────────────────────────┘    │
│             ▼                                                           │
│  ┌──────────────────────────────────────────────────────────────┐      │
│  │              WORKER POOL SYSTEM (6 Pools)                    │      │
│  ├──────────────────────────────────────────────────────────────┤      │
│  │  Event Worker Pool      │  Configurable size per pool        │      │
│  │  Branch Worker Pool     │  Buffered channels (1000+ capacity)│      │
│  │  WebEngage Worker Pool  │  Safe goroutine execution          │      │
│  │  Vidooly Worker Pool    │  Panic recovery                    │      │
│  │  Graphy Worker Pool     │  Event grouping by type            │      │
│  │  Shopify Worker Pool    │  Timestamp correction              │      │
│  └──────────────────────────────────────────────────────────────┘      │
└────────────────────────────────┬────────────────────────────────────────┘
                                 │ Publish (JSON)
                                 ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                     RABBITMQ MESSAGE BROKER                              │
├─────────────────────────────────────────────────────────────────────────┤
│  EXCHANGES:                    │  QUEUES (26 Configured):               │
│  - grpc_event.tx              │  - grpc_clickhouse_event_q (2 workers) │
│  - grpc_event.dx              │  - clickhouse_click_event_q (2)        │
│  - grpc_event_error.dx        │  - app_init_event_q (2)                │
│  - branch_event.tx            │  - branch_event_q (2)                  │
│  - webengage_event.dx         │  - webengage_event_q (3)               │
│  - identity.dx                │  - webengage_ch_event_q (5)            │
│  - beat.dx                    │  - post_log_events_q (20)              │
│  - coffee.dx                  │  - profile_log_events_q (2)            │
│  - shopify_event.dx           │  - sentiment_log_events_q (2)          │
│  - affiliate.dx               │  - scrape_request_log_events_q (2)     │
│  - bigboss                    │  - activity_tracker_q (2)              │
│                               │  ... and 15 more queues                │
├─────────────────────────────────────────────────────────────────────────┤
│  FEATURES:                                                               │
│  - Durable queues             - Retry logic (max 2 attempts)            │
│  - Error queues               - Dead letter routing                     │
│  - Prefetch QoS (1)           - Auto-reconnect on failure              │
└────────────────────────────────┬────────────────────────────────────────┘
                                 │ Consume
                                 ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                        SINKER LAYER (20+ Sinkers)                       │
├─────────────────────────────────────────────────────────────────────────┤
│  EVENT SINKERS:                │  BUFFERED SINKERS:                     │
│  - SinkEventToClickhouse       │  - TraceLogEventsSinker (batch)       │
│  - SinkErrorEventToClickhouse  │  - AffiliateOrdersSinker (batch)      │
│  - SinkClickEventToClickhouse  │  - BigBossVotesSinker (batch)         │
│  - SinkAppInitEvent            │  - PostLogEventsSinker (batch)        │
│  - SinkLaunchReferEvent        │  - ProfileLogEventsSinker (batch)     │
│  - SinkBranchEvent             │  - SentimentLogEventsSinker (batch)   │
│  - SinkGraphyEvent             │  - ScrapeLogEventsSinker (batch)      │
│  - SinkShopifyEvent            │  - OrderLogEventsSinker (batch)       │
│  - SinkWebengageToClickhouse   │  - ActivityTrackerSinker (batch)      │
│  - SinkWebengageToAPI          │                                        │
└────────────────────────────────┬────────────────────────────────────────┘
                                 │
        ┌────────────────────────┼────────────────────────┐
        ▼                        ▼                        ▼
┌───────────────────┐  ┌───────────────────┐  ┌───────────────────────┐
│    CLICKHOUSE     │  │    POSTGRESQL     │  │    EXTERNAL APIs      │
│   (Analytics DB)  │  │   (Transactional) │  │                       │
├───────────────────┤  ├───────────────────┤  ├───────────────────────┤
│  Tables:          │  │  Tables:          │  │  - WebEngage API      │
│  - event          │  │  - referral_event │  │  - Identity Service   │
│  - error_event    │  │  - user_account   │  │  - Social Stream API  │
│  - branch_event   │  │                   │  │                       │
│  - click_event    │  │  Features:        │  │  Features:            │
│  - trace_log      │  │  - Schema-based   │  │  - Resty HTTP client  │
│  - webengage_event│  │  - Pool: 10/20    │  │  - Retry logic        │
│  - graphy_event   │  │  - Transactions   │  │  - Error handling     │
│  - + 10 more      │  │                   │  │                       │
├───────────────────┤  └───────────────────┘  └───────────────────────┘
│  Features:        │
│  - Multi-DB       │
│  - Auto-reconnect │
│  - GORM ORM       │
│  - Batch inserts  │
└───────────────────┘
```

---

## 3. PROTOCOL BUFFERS & gRPC SERVICE

### eventservice.proto (688 lines)

```protobuf
syntax = "proto3";
option go_package = "proto/go/bulbulgrpc";
option java_package = "com.bulbul.grpc.event";

service EventService {
    rpc dispatch(Events) returns (Response) {}
}

message Events {
    Header header = 1;
    repeated Event events = 2;
}

message Header {
    string sessionId = 1;
    string bbDeviceId = 2;
    string userId = 3;
    string deviceId = 4;
    string androidAdvertisingId = 5;
    string clientId = 6;
    string channel = 7;
    string os = 8;
    string clientType = 9;
    string appLanguage = 10;
    string merchantId = 11;
    string ppId = 12;
    string appVersion = 13;
    string currentURL = 14;
    string utmReferrer = 15;
    string utmPlatform = 16;
    string utmSource = 17;
    string utmMedium = 18;
    string utmCampaign = 19;
    // ... 5 more fields (24 total)
}
```

### Event Types Defined (60+)

**User Flow Events:**
- `LaunchEvent`, `LaunchReferEvent`
- `PageOpenedEvent`, `PageLoadedEvent`, `PageLoadFailedEvent`

**Widget Events:**
- `WidgetViewEvent`, `WidgetCtaClickedEvent`
- `WidgetElementViewEvent`, `WidgetElementClickedEvent`

**Commerce Events:**
- `AddToCartEvent`, `GoToPaymentsEvent`
- `InitiatePurchaseEvent`, `CompletePurchaseEvent`, `PurchaseFailedEvent`

**Streaming Events:**
- `StreamEnterEvent`, `ChooseProductEvent`
- `SocialStreamEnterEvent`, `SocialStreamActionEvent`

**Review Events:**
- `InitiateReview`, `EnterReviewScreen`, `SubmitReview`
- `ViewAllReviews`, `ExpandReviewImage`

**Authentication Events:**
- `PhoneVerificationInitiateEvent`
- `PhoneNumberEntered`, `OTPVerifiedEvent`

**System Events:**
- `AppStartEvent`, `NotificationActionEvent`
- `SessionIdChangeEvent`, `TestEvent`

---

## 4. WORKER POOL IMPLEMENTATION

### Architecture Pattern

```go
// Singleton pattern with lazy initialization
var (
    eventWrapperChannel chan bulbulgrpc.Events
    channelInit         sync.Once
)

func GetChannel(config config.Config) chan bulbulgrpc.Events {
    channelInit.Do(func() {
        // Create buffered channel
        eventWrapperChannel = make(chan bulbulgrpc.Events,
            config.EVENT_WORKER_POOL_CONFIG.EVENT_BUFFERED_CHANNEL_SIZE)

        // Initialize worker pool
        initWorkerPool(config,
            config.EVENT_WORKER_POOL_CONFIG.EVENT_WORKER_POOL_SIZE,
            eventWrapperChannel)
    })
    return eventWrapperChannel
}

func initWorkerPool(config config.Config, poolSize int,
    eventChannel <-chan bulbulgrpc.Events) {
    for i := 0; i < poolSize; i++ {
        safego.GoNoCtx(func() {
            worker(config, eventChannel)
        })
    }
}
```

### Worker Processing Logic

```go
func worker(config config.Config, eventChannel <-chan bulbulgrpc.Events) {
    for e := range eventChannel {
        rabbitConn := rabbit.Rabbit(config)

        // Group events by type
        eventsGrouped := make(map[string][]*bulbulgrpc.Event)
        for _, evt := range e.Events {
            eventName := fmt.Sprintf("%v", reflect.TypeOf(evt.GetEventOf()))
            eventName = strings.ReplaceAll(eventName, "*bulbulgrpc.Event_", "")

            // Timestamp correction (future timestamps set to now)
            et := time.Unix(transformerToInt64(evt.Timestamp)/1000, 0)
            if et.Unix() > time.Now().Unix() {
                evt.Timestamp = strconv.FormatInt(time.Now().Unix()*1000, 10)
            }

            eventsGrouped[eventName] = append(eventsGrouped[eventName], evt)
        }

        // Publish grouped events to RabbitMQ
        for eventName, events := range eventsGrouped {
            e := bulbulgrpc.Events{Header: header, Events: events}
            if b, err := protojson.Marshal(&e); err == nil {
                rabbitConn.Publish("grpc_event.tx", eventName, b, map[string]interface{}{})
            }
        }
    }
}
```

### Worker Pool Configurations

| Worker Pool | Channel Size | Pool Size | Exchange |
|-------------|--------------|-----------|----------|
| Event Worker | 1000 | Configurable | grpc_event.tx |
| Branch Worker | 1000 | Configurable | branch_event.tx |
| WebEngage Worker | 1000 | Configurable | webengage_event.dx |
| Vidooly Worker | 1000 | Configurable | vidooly_event.dx |
| Graphy Worker | 1000 | Configurable | graphy_event.tx |
| Shopify Worker | 1000 | Configurable | shopify_event.dx |

---

## 5. RABBITMQ INTEGRATION

### Connection Management

```go
var (
    singletonRabbit  *RabbitConnection
    rabbitCloseError chan *amqp.Error
    rabbitInit       sync.Once
)

func Rabbit(config config.Config) *RabbitConnection {
    rabbitInit.Do(func() {
        rabbitConnected := make(chan bool)
        safego.GoNoCtx(func() {
            rabbitConnector(config, rabbitConnected)
        })
        select {
        case <-rabbitConnected:
        case <-time.After(5 * time.Second):
        }
    })
    return singletonRabbit
}
```

### Consumer Configuration (26 Consumers in main.go)

```go
// Example: High-volume buffered consumer
traceLogChan := make(chan interface{}, 10000)
traceLogEventConsumerCfg := rabbit.RabbitConsumerConfig{
    QueueName:            "trace_log",
    Exchange:             "identity.dx",
    RoutingKey:           "trace_log",
    RetryOnError:         true,
    ErrorExchange:        &errorExchange,
    ErrorRoutingKey:      &errorRoutingKey,
    ConsumerCount:        2,
    BufferChan:           traceLogChan,
    BufferedConsumerFunc: sinker.BufferTraceLogEvent,
}
rabbit.Rabbit(config).InitConsumer(traceLogEventConsumerCfg)
go sinker.TraceLogEventsSinker(traceLogChan)  // Batch processor
```

### Retry Logic Implementation

```go
func (rabbit *RabbitConnection) Consume(cfg RabbitConsumerConfig) error {
    for msg := range msgs {
        listenerResponse := consumerFunc(msg)

        if listenerResponse {
            msg.Ack(false)
        } else if !retryOnError {
            publishToErrorQueue(msg)
            msg.Ack(false)
        } else {
            retryCount := getRetryCount(msg.Headers)
            if retryCount >= 2 {
                publishToErrorQueue(msg)
                msg.Ack(false)
            } else {
                // Republish with incremented retry count
                msg.Headers["x-retry-count"] = retryCount + 1
                republish(msg)
                msg.Ack(false)
            }
        }
    }
}
```

### All 26 Consumer Queues

| Queue | Workers | Purpose |
|-------|---------|---------|
| grpc_clickhouse_event_q | 2 | Main events → ClickHouse |
| grpc_clickhouse_event_error_q | 2 | Error events → ClickHouse |
| clickhouse_click_event_q | 2 | Click tracking |
| app_init_event_q | 2 | App initialization |
| launch_refer_event_q | 2 | Launch referrals |
| create_user_account_q | 2 | User account creation |
| event.account_profile_complete | 2 | Profile completion |
| ab_assignments | 2 | A/B test assignments |
| branch_event_q | 2 | Branch.io events |
| graphy_event_q | 2 | Graphy events |
| trace_log | 2 | Trace logs (buffered) |
| affiliate_orders_event_q | 2 | Affiliate orders (buffered) |
| bigboss_votes_log_q | 2 | BigBoss votes (buffered) |
| post_log_events_q | 20 | Social post logs (buffered) |
| sentiment_log_events_q | 2 | Sentiment analysis (buffered) |
| post_activity_log_events_q | 2 | Post activity (buffered) |
| profile_log_events_q | 2 | Profile logs (buffered) |
| profile_relationship_log_events_q | 2 | Relationships (buffered) |
| scrape_request_log_events_q | 2 | Scrape requests (buffered) |
| order_log_events_q | 2 | Order logs (buffered) |
| shopify_events_q | 2 | Shopify events (buffered) |
| webengage_event_q | 3 | WebEngage → API |
| webengage_ch_event_q | 5 | WebEngage → ClickHouse |
| webengage_user_event_q | 5 | WebEngage user events |
| post_log_events_q_bkp | 5 | Backup post logs |
| activity_tracker_q | 2 | Partner activity (buffered) |

---

## 6. DATABASE LAYER

### ClickHouse Connection (Multi-DB Support)

```go
var (
    singletonClickhouseMap map[string]*gorm.DB
    clickhouseInit         sync.Once
)

func Clickhouse(config config.Config, dbName *string) *gorm.DB {
    clickhouseInit.Do(func() {
        connectToClickhouse(config)
        go clickhouseConnectionCron(config)  // Reconnect every 1 second
    })

    if dbName != nil {
        if db, ok := singletonClickhouseMap[*dbName]; ok {
            return db
        }
    }
    return singletonClickhouseMap[config.CLICKHOUSE_DB_NAME]
}

func connectToClickhouse(config config.Config) {
    for _, dbName := range config.CLICKHOUSE_DB_NAMES {
        dsn := "tcp://" + host + ":" + port +
               "?database=" + dbName +
               "&read_timeout=10&write_timeout=20"

        db, _ := gorm.Open(clickhouse.New(clickhouse.Config{
            DSN: dsn,
            DisableDatetimePrecision: true,
        }), &gorm.Config{})

        singletonClickhouseMap[dbName] = db
    }
}
```

### PostgreSQL Connection (Schema-Based)

```go
func PG(config config.Config) *gorm.DB {
    pgInit.Do(func() {
        connectToPG(config)
        go pgConnectionCron(config)
    })
    return singletonPg
}

func connectToPG(config config.Config) {
    singletonPg, _ = gorm.Open("postgres",
        "host=" + host + " port=" + port +
        " user=" + user + " dbname=" + dbname +
        " password=" + password)

    singletonPg.DB().SetMaxIdleConns(10)
    singletonPg.DB().SetMaxOpenConns(20)

    // Schema-based table naming
    gorm.DefaultTableNameHandler = func(db *gorm.DB, tableName string) string {
        return config.POSTGRES_EVENT_SCHEMA + "." + tableName
    }
}
```

### ClickHouse Tables Created

| Table | Purpose | Key Fields |
|-------|---------|------------|
| event | Main app events | event_id, event_name, event_params (JSONB) |
| error_event | Failed event tracking | All event fields + error info |
| click_event | Click tracking | 600+ lines of click-specific data |
| branch_event | Deep linking | 49 fields for attribution |
| webengage_event | Marketing events | userId, eventData (JSONB) |
| trace_log | Request tracing | hostName, serviceName, timeTaken |
| graphy_event | Graphy platform | Platform-specific events |
| app_init_event | App startup | Device, session metrics |
| launch_refer_event | Launch referrals | Attribution data |
| ab_assignment | A/B test data | Variant assignments |
| entity_metric | Generic metrics | Flexible metric storage |
| post_log_event | Social posts | Post activity data |
| sentiment_log | Sentiment analysis | Comment sentiment scores |
| profile_log | Profile activity | Profile metrics over time |
| profile_relationship_log | Relationships | Follower/following data |
| order_log | E-commerce orders | Order transaction data |
| shopify_event | Shopify events | Commerce event data |
| affiliate_order | Affiliate tracking | Affiliate transaction data |

---

## 7. DATA MODELS

### Core Event Model

```go
type Event struct {
    EventId         string    `gorm:"column:event_id"`
    EventName       string    `gorm:"column:event_name"`
    EventTimestamp  time.Time `gorm:"column:event_timestamp"`
    InsertTimestamp time.Time `gorm:"column:insert_timestamp"`

    // Header fields
    SessionId       string    `gorm:"column:session_id"`
    BbDeviceId      int64     `gorm:"column:bb_device_id"`
    UserId          int64     `gorm:"column:user_id"`
    DeviceId        string    `gorm:"column:device_id"`
    ClientId        string    `gorm:"column:client_id"`
    Channel         string    `gorm:"column:channel"`
    Os              string    `gorm:"column:os"`
    ClientType      string    `gorm:"column:client_type"`
    AppLanguage     string    `gorm:"column:app_language"`
    AppVersion      string    `gorm:"column:app_version"`

    // UTM parameters
    UtmReferrer     string    `gorm:"column:utm_referrer"`
    UtmPlatform     string    `gorm:"column:utm_platform"`
    UtmSource       string    `gorm:"column:utm_source"`
    UtmMedium       string    `gorm:"column:utm_medium"`
    UtmCampaign     string    `gorm:"column:utm_campaign"`

    // Dynamic event data
    EventParams     JSONB     `gorm:"column:event_params;type:String"`
}
```

### Branch Event Model (49 Fields)

```go
type BranchEvent struct {
    // Core identifiers
    Id                    string `gorm:"column:id"`
    Name                  string `gorm:"column:name"`
    Timestamp             int64  `gorm:"column:timestamp"`

    // Attribution data
    Attributed            bool   `gorm:"column:attributed"`
    DeepLinked            bool   `gorm:"column:deep_linked"`
    ExistingUser          bool   `gorm:"column:existing_user"`
    HasClicked            bool   `gorm:"column:has_clicked"`
    HasApp                bool   `gorm:"column:has_app"`

    // Timing metrics
    SecondsFromInstall    int64  `gorm:"column:seconds_from_install"`
    SecondsFromLastOpen   int64  `gorm:"column:seconds_from_last_attributed_touch_timestamp"`

    // Geo data
    GeoCountryEn          string `gorm:"column:geo_country_en"`
    GeoRegionEn           string `gorm:"column:geo_region_en"`
    GeoCityEn             string `gorm:"column:geo_city_en"`

    // Campaign data
    Campaign              string `gorm:"column:campaign"`
    Channel               string `gorm:"column:channel"`
    Feature               string `gorm:"column:feature"`
    Tags                  string `gorm:"column:tags"`

    // Complex nested data (JSONB)
    EventData             JSONB  `gorm:"column:event_data;type:String"`
    UserData              JSONB  `gorm:"column:user_data;type:String"`
    LastAttributedTouchData JSONB `gorm:"column:last_attributed_touch_data;type:String"`

    // ... 30 more fields for complete attribution tracking
}
```

---

## 8. SINKER IMPLEMENTATIONS

### Main Event Sinker

```go
func SinkEventToClickhouse(delivery amqp.Delivery) bool {
    grpcEvent := &bulbulgrpc.Events{}
    if err := protojson.Unmarshal(delivery.Body, grpcEvent); err != nil {
        return false
    }

    if grpcEvent.Events == nil || len(grpcEvent.Events) == 0 {
        return false
    }

    events := []model.Event{}
    for _, e := range grpcEvent.Events {
        event, err := TransformToEventModel(grpcEvent, e)
        if err != nil {
            continue
        }
        events = append(events, event)
    }

    if len(events) > 0 {
        result := clickhouse.Clickhouse(config.New(), nil).Create(events)
        return result != nil && result.Error == nil
    }
    return false
}
```

### Buffered Sinker Pattern (High-Volume)

```go
// Buffer function - adds to channel
func BufferTraceLogEvent(delivery amqp.Delivery, c chan interface{}) bool {
    var traceLog map[string]interface{}
    if err := json.Unmarshal(delivery.Body, &traceLog); err == nil {
        c <- traceLog
        return true
    }
    return false
}

// Batch processor - runs in separate goroutine
func TraceLogEventsSinker(c chan interface{}) {
    ticker := time.NewTicker(5 * time.Second)
    batch := []model.TraceLogEvent{}

    for {
        select {
        case event := <-c:
            traceLog := parseTraceLog(event)
            batch = append(batch, traceLog)

            if len(batch) >= 1000 {
                flushBatch(batch)
                batch = []model.TraceLogEvent{}
            }

        case <-ticker.C:
            if len(batch) > 0 {
                flushBatch(batch)
                batch = []model.TraceLogEvent{}
            }
        }
    }
}
```

### All Sinker Functions

| Sinker | Type | Destination |
|--------|------|-------------|
| SinkEventToClickhouse | Direct | ClickHouse event table |
| SinkErrorEventToClickhouse | Direct | ClickHouse error_event |
| SinkClickEventToClickhouse | Direct | ClickHouse click_event |
| SinkAppInitEvent | Direct | ClickHouse app_init_event |
| SinkLaunchReferEvent | Direct | ClickHouse + PostgreSQL |
| SinkBranchEvent | Direct | ClickHouse branch_event |
| SinkGraphyEvent | Direct | ClickHouse graphy_event |
| SinkWebengageToClickhouse | Direct | ClickHouse webengage_event |
| SinkWebengageToAPI | HTTP | WebEngage REST API |
| SinkShopifyEvent | Buffered | ClickHouse shopify_event |
| TraceLogEventsSinker | Buffered | ClickHouse trace_log |
| AffiliateOrdersSinker | Buffered | ClickHouse affiliate_order |
| BigBossVotesSinker | Buffered | ClickHouse bigboss_votes |
| PostLogEventsSinker | Buffered | ClickHouse post_log_event |
| SentimentLogEventsSinker | Buffered | ClickHouse sentiment_log |
| ProfileLogEventsSinker | Buffered | ClickHouse profile_log |
| ScrapeLogEventsSinker | Buffered | ClickHouse scrape_request_log |
| OrderLogEventsSinker | Buffered | ClickHouse order_log |
| ActivityTrackerSinker | Buffered | ClickHouse activity_tracker |

---

## 9. ERROR HANDLING & RESILIENCE

### Safe Goroutine Execution

```go
func GoNoCtx(f func()) {
    go func() {
        defer func() {
            if panicMessage := recover(); panicMessage != nil {
                stack := debug.Stack()
                log.Printf("RECOVERED FROM PANIC: %v\nSTACK: %s",
                    panicMessage, stack)
            }
        }()
        f()
    }()
}
```

### Connection Auto-Recovery

```go
// ClickHouse reconnect cron
func clickhouseConnectionCron(config config.Config) {
    ticker := time.NewTicker(1 * time.Second)
    for range ticker.C {
        for dbName, db := range singletonClickhouseMap {
            if db == nil || db.Error != nil {
                reconnect(dbName)
            }
        }
    }
}

// RabbitMQ reconnect on close
func rabbitConnector(config config.Config, connected chan bool) {
    for {
        connectRabbit(config)
        connected <- true

        // Wait for close notification
        <-rabbitCloseError

        // Reconnect with backoff
        time.Sleep(2 * time.Second)
    }
}
```

### Error Flow

```
Event Processing
    ↓
Success? → ACK & Done
    ↓ (Failure)
Retry 1 → Republish with x-retry-count=1
    ↓ (Failure)
Retry 2 → Republish with x-retry-count=2
    ↓ (Failure)
Route to Error Exchange → Dead Letter Queue
```

---

## 10. CI/CD PIPELINE

### GitLab CI Configuration

```yaml
stages:
  - build
  - deploy
  - publish

build:
  stage: build
  variables:
    GOOS: linux
    GOARCH: amd64
  script:
    - protoc --go_out=. --go-grpc_out=. proto/eventservice.proto
    - env GOOS=$GOOS GOARCH=$GOARCH go build
  artifacts:
    paths:
      - event-grpc
      - .env*
      - scripts/start.sh
    expire_in: 1 week

deploy_staging:
  stage: deploy
  variables:
    ENV: STAGE
    GOGC: 250
  script:
    - /bin/bash scripts/start.sh
  tags:
    - deployer-stage-search
  only:
    - master
    - dev

deploy_production:
  stage: deploy
  when: manual
  only:
    - master

publish_web:
  stage: publish
  script:
    - protoc proto/eventservice.proto
        --js_out=import_style=commonjs:./proto/web/src
        --grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:./proto/web/src/
    - npm publish --registry http://artifactory.bulbul.tv:4873/
  when: manual

publish_java:
  stage: publish
  script:
    - mvn clean install -U && mvn deploy
  when: manual
```

### Graceful Deployment Script

```bash
#!/bin/bash
ulimit -n 100000

PID=$(ps aux | grep event-grpc | grep -v grep | awk '{print $2}')

if [ ! -z "$PID" ]; then
    # Remove from load balancer
    curl -vXPUT http://localhost:$GIN_PORT/heartbeat/?beat=false
    sleep 15

    # Kill existing process
    kill -9 $PID
fi

sleep 10

# Start new process
ENV=$CI_ENVIRONMENT_NAME ./event-grpc >> "logs/out.log" 2>&1 &

# Symlink for log access
ln -s $(pwd)/logs /bulbul/services/event-grpc/logs

# Wait for startup
sleep 20

# Add back to load balancer
curl -vXPUT http://localhost:$GIN_PORT/heartbeat/?beat=true
```

---

## 11. CONFIGURATION MANAGEMENT

### Environment-Based Config (40+ Fields)

```go
type Config struct {
    // Server
    Port    string  // 8017 (gRPC)
    GinPort string  // 8019 (HTTP)
    Env     string  // PRODUCTION, STAGE, LOCAL

    // External APIs
    IDENTITY_URL      string
    SOCIAL_STREAM_URL string
    WEBENGAGE_API_KEY string
    WEBENGAGE_URL     string
    WEBENGAGE_USER_URL string
    WEBENGAGE_BRAND_API_KEY string
    WEBENGAGE_BRAND_URL string

    // OAuth Client IDs
    CUSTOMER_APP_CLIENT_ID string
    HOST_APP_CLIENT_ID     string
    GCC_HOST_APP_CLIENT_ID string
    GCC_BRAND_CLIENT_ID    string

    // Redis
    RedisClusterAddresses  []string
    REDIS_CLUSTER_PASSWORD string

    // Worker Pools (6 configurations)
    EVENT_WORKER_POOL_CONFIG          *EVENT_WORKER_POOL_CONFIG
    BRANCH_EVENT_WORKER_POOL_CONFIG   *EVENT_WORKER_POOL_CONFIG
    WEBENGAGE_EVENT_WORKER_POOL_CONFIG *EVENT_WORKER_POOL_CONFIG
    VIDOOLY_EVENT_WORKER_POOL_CONFIG  *EVENT_WORKER_POOL_CONFIG
    GRAPHY_EVENT_WORKER_POOL_CONFIG   *EVENT_WORKER_POOL_CONFIG
    SHOPIFY_EVENT_WORKER_POOL_CONFIG  *EVENT_WORKER_POOL_CONFIG

    // Databases
    CLICKHOUSE_CONNECTION_CONFIG *CLICKHOUSE_CONNECTION_CONFIG
    POSTGRES_CONNECTION_CONFIG   *POSTGRES_CONNECTION_CONFIG
    RABBIT_CONNECTION_CONFIG     *RABBIT_CONNECTION_CONFIG

    // Security
    HMAC_SECRET []byte
    LOG_LEVEL   int
}

type EVENT_WORKER_POOL_CONFIG struct {
    EVENT_WORKER_POOL_SIZE      int
    EVENT_BUFFERED_CHANNEL_SIZE int
}

type CLICKHOUSE_CONNECTION_CONFIG struct {
    CLICKHOUSE_HOST     string
    CLICKHOUSE_PORT     string
    CLICKHOUSE_USER     string
    CLICKHOUSE_DB_NAMES []string  // Comma-separated
    CLICKHOUSE_PASSWORD string
}

type RABBIT_CONNECTION_CONFIG struct {
    RABBIT_USER      string
    RABBIT_PASSWORD  string
    RABBIT_HOST      string
    RABBIT_PORT      string
    RABBIT_HEARTBEAT int  // milliseconds
    RABBIT_VHOST     string
}
```

---

## 12. OBSERVABILITY

### Prometheus Metrics

```go
router.GET("/metrics", gin.WrapH(promhttp.Handler()))
```

### Sentry Integration

```go
sentry.Init(sentry.ClientOptions{
    Dsn:              "http://xxx@172.31.14.149:9000/26",
    Environment:      config.Env,
    AttachStacktrace: true,
})

router.Use(sentrygin.New(sentrygin.Options{Repanic: true}))
```

### Request Logging

```go
func RequestLogger(config config.Config) gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()

        // Capture request/response bodies
        buf, _ := ioutil.ReadAll(c.Request.Body)
        c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(buf))

        c.Next()

        // Log with configurable verbosity
        if gc.Config.LOG_LEVEL < 3 {
            gc.Logger.Error().Msg(fmt.Sprintf(
                "%s - %s - [%s] \"%s %s\" %d %s\n%s",
                c.Request.Header.Get(header.RequestID),
                c.ClientIP(),
                time.Now().Format(time.RFC1123),
                c.Request.Method, path,
                c.Writer.Status(),
                time.Since(start),
                body,  // Request body (verbose mode)
            ))
        }
    }
}
```

### Health Check Endpoint

```go
// GET /heartbeat/ - Returns health status
// PUT /heartbeat/?beat=false - Remove from LB
// PUT /heartbeat/?beat=true - Add to LB

func (s *healthserver) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest)
    (*grpc_health_v1.HealthCheckResponse, error) {
    if heartbeat.BEAT {
        return &grpc_health_v1.HealthCheckResponse{
            Status: grpc_health_v1.HealthCheckResponse_SERVING,
        }, nil
    }
    return &grpc_health_v1.HealthCheckResponse{
        Status: grpc_health_v1.HealthCheckResponse_NOT_SERVING,
    }, nil
}
```

---

## 13. KEY METRICS & STATISTICS

| Metric | Value |
|--------|-------|
| **Total Lines of Code** | 10,000+ |
| **Go Source Files** | 90+ |
| **Test Files** | 8 |
| **Proto Definitions** | 688 lines |
| **Event Types** | 60+ |
| **Database Models** | 19 |
| **Sinker Functions** | 20+ |
| **RabbitMQ Consumers** | 26 configured |
| **Total Consumer Workers** | 70+ |
| **Worker Pools** | 6 |
| **Buffered Channel Capacity** | 100K+ combined |
| **Direct Dependencies** | 30+ |
| **Exchanges** | 11 |
| **ClickHouse Tables** | 18+ |

---

## 14. SKILLS DEMONSTRATED

### Technical Skills

| Skill | Evidence |
|-------|----------|
| **gRPC & Protocol Buffers** | 688-line proto with 60+ event types, Java/JS client publishing |
| **Go Concurrency** | Worker pools, channels, sync.Once, goroutines, panic recovery |
| **Message Queues** | RabbitMQ with 26 queues, retry logic, dead letter routing |
| **Database Design** | ClickHouse (multi-DB), PostgreSQL (schema-based), connection pooling |
| **Distributed Systems** | Event-driven architecture, fault tolerance, auto-recovery |
| **API Design** | gRPC + REST hybrid, health checks, metrics endpoints |
| **DevOps** | GitLab CI/CD, graceful deployments, zero-downtime, client library publishing |
| **Observability** | Prometheus, Sentry, structured logging, request tracing |

### Architecture Patterns

1. **Worker Pool Pattern** - Configurable concurrency with buffered channels
2. **Publish-Subscribe** - RabbitMQ exchange-based routing
3. **Singleton Pattern** - Connection pools with lazy initialization
4. **Circuit Breaker** - Auto-reconnect on database/queue failures
5. **Buffered Sinker** - Batch processing for high-volume events
6. **Gateway Pattern** - Unified entry point with middleware pipeline

---

## 15. INTERVIEW TALKING POINTS

### System Design Questions

**"Design a real-time event processing system"**
- Built gRPC server handling high-throughput events
- 6 Worker pools with configurable concurrency
- RabbitMQ for reliable message delivery with retry logic (max 2 retries)
- Multi-destination routing (ClickHouse, PostgreSQL, External APIs)
- Buffered sinkers for high-volume batch processing

**"How do you ensure reliability in distributed systems?"**
- Auto-reconnect on connection failures (1-second cron)
- Message retry with dead letter queues
- Safe goroutine execution with panic recovery
- Health checks for load balancer integration
- Graceful deployment with zero-downtime

**"Explain your experience with message queues"**
- 26 RabbitMQ consumer configurations
- 11 exchanges for event type distribution
- Durable queues with prefetch QoS
- Error queue pattern for debugging
- Buffered consumption for high-volume events

### Behavioral Questions

**"Tell me about a complex system you built"**
- Event-gRPC: 10K+ LOC, 60+ event types, 26 message queues
- Handles real-time analytics for mobile/web applications
- Multi-database strategy (ClickHouse for analytics, PostgreSQL for transactions)
- Client library publishing (Java via Maven, JavaScript via NPM)

**"How do you handle scale?"**
- Configurable worker pool sizes
- Buffered channels (100K+ capacity)
- Batch inserts to databases
- Horizontal scaling via stateless design
- 70+ concurrent consumer workers

---

*Generated through comprehensive source code analysis of the event-grpc project.*
