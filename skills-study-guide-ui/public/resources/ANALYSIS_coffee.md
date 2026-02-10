# ULTRA-DEEP ANALYSIS: COFFEE PROJECT

## PROJECT OVERVIEW

| Attribute | Value |
|-----------|-------|
| **Project Name** | Coffee |
| **Purpose** | Multi-tenant SaaS platform for influencer discovery, profile collections, and social media analytics |
| **Architecture** | 4-Layer REST microservice with PostgreSQL + ClickHouse dual database strategy |
| **Language** | Go 1.18 |
| **Total Lines of Code** | ~8,500+ |
| **Port** | 7179 (main), 9292 (debug/metrics) |

---

## 1. COMPLETE DIRECTORY STRUCTURE

```
coffee/
├── main.go                           # Application entry point
├── go.mod / go.sum                    # Dependencies (30+ packages)
├── schema.sql (918 lines)             # Database schema (27+ tables)
├── .gitlab-ci.yml                     # CI/CD pipeline
├── .env.*                             # Environment configs
│
├── server/
│   ├── server.go                      # HTTP server setup with Chi
│   └── middlewares/
│       ├── context.go                 # Application context extraction
│       ├── session.go                 # Transaction management
│       ├── errorhandling.go           # Sentry error integration
│       └── requestinterceptor.go      # Request/response capture
│
├── core/
│   ├── rest/
│   │   ├── api.go                     # Generic API handler interface
│   │   ├── service.go                 # Generic Service layer
│   │   ├── manager.go                 # Generic Manager layer
│   │   └── dao.go                     # Generic DAO interface
│   ├── persistence/
│   │   ├── postgres/db.go             # PostgreSQL connection
│   │   ├── clickhouse/db.go           # ClickHouse connection
│   │   └── redis/redis.go             # Redis cluster client
│   ├── appcontext/
│   │   └── requestcontext.go          # Multi-tenant request context
│   └── domain/
│       └── domain.go                  # Core types (Entry, Entity, Response)
│
├── app/
│   └── app.go                         # Service container & DI
│
├── routes/
│   └── services.go                    # Route registration
│
├── discovery/                         # Profile discovery module
│   ├── api/api.go
│   ├── service/service.go
│   ├── manager/
│   │   ├── manager.go
│   │   ├── searchmanager.go
│   │   ├── timeseriesmanager.go
│   │   ├── hashtagsmanager.go
│   │   └── audiencemanager.go
│   ├── dao/dao.go
│   ├── domain/
│   └── listeners/
│
├── profilecollection/                 # Profile collections
├── postcollection/                    # Post collections
├── leaderboard/                       # Rankings system
├── collectionanalytics/               # Collection analytics
├── genreinsights/                     # Genre insights
├── keywordcollection/                 # Keyword collections
├── campaignprofiles/                  # Campaign profiles
├── collectiongroup/                   # Collection grouping
├── partnerusage/                      # Partner usage tracking
├── content/                           # Content management
│
├── listeners/
│   └── listeners.go                   # Watermill event listeners
├── publishers/
│   └── publisher.go                   # Event publishing
├── amqp/
│   └── config.go                      # AMQP configuration
│
├── client/                            # External service clients
│   ├── partner/                       # Partner service client
│   ├── dam/                           # Digital asset management
│   ├── winkl/                         # Legacy Winkl integration
│   └── beat/                          # Beat service client
│
├── helpers/                           # Utility functions
├── constants/constants.go             # Constants (languages, categories)
├── config/config.go                   # Viper configuration
├── logger/                            # Logrus setup
├── sentry/                            # Sentry error tracking
└── scripts/
    └── start.sh                       # Deployment script
```

---

## 2. FOUR-LAYER REST ARCHITECTURE

### Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────────────┐
│                           HTTP REQUEST                                   │
└────────────────────────────────┬────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                     MIDDLEWARE PIPELINE (10 Layers)                      │
├─────────────────────────────────────────────────────────────────────────┤
│  1. RequestInterceptor     │  Capture request/response for errors       │
│  2. SlowQueryLogger        │  Log slow queries                          │
│  3. SentryMiddleware       │  Error tracking integration                │
│  4. ApplicationContext     │  Extract headers → RequestContext          │
│  5. ServiceSession         │  Transaction + timeout management          │
│  6. chi/Logger             │  HTTP request logging                      │
│  7. chi/RequestID          │  Request ID generation                     │
│  8. chi/RealIP             │  Client IP extraction                      │
│  9. chi-prometheus         │  Metrics (300ms-30s buckets)               │
│ 10. Recoverer              │  Panic recovery                            │
└────────────────────────────────┬────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                         API LAYER (Handler)                              │
├─────────────────────────────────────────────────────────────────────────┤
│  Interface: ApiWrapper                                                   │
│  Methods:                                                                │
│    - AttachRoutes(r *chi.Mux)  // Register HTTP routes                  │
│    - GetPrefix() string         // Service path prefix                  │
│                                                                          │
│  Responsibilities:                                                       │
│    - Parse HTTP request (body, params, headers)                         │
│    - Validate input                                                      │
│    - Call Service layer                                                  │
│    - Render JSON response                                                │
└────────────────────────────────┬────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                        SERVICE LAYER                                     │
├─────────────────────────────────────────────────────────────────────────┤
│  Generic Type: Service[RES Response, EX Entry, EN Entity, I ID]         │
│                                                                          │
│  Methods:                                                                │
│    - Create(ctx, entry) (*RES, error)                                   │
│    - FindById(ctx, id) (*RES, error)                                    │
│    - FindByIds(ctx, ids) ([]RES, error)                                 │
│    - Update(ctx, id, entry) (*RES, error)                               │
│    - Search(ctx, query, sort, page, size) ([]RES, int64, error)         │
│    - InitializeSession(ctx) error                                        │
│    - Close(ctx) error                                                    │
│    - Rollback(ctx) error                                                 │
│                                                                          │
│  Responsibilities:                                                       │
│    - Transaction management                                              │
│    - Orchestrate Manager calls                                           │
│    - Response transformation                                             │
└────────────────────────────────┬────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                        MANAGER LAYER                                     │
├─────────────────────────────────────────────────────────────────────────┤
│  Generic Type: Manager[EX Entry, EN Entity, I ID]                        │
│                                                                          │
│  Methods:                                                                │
│    - toEntity(*EX) (*EN, error)    // Entry → Entity                    │
│    - toEntry(*EN) (*EX, error)     // Entity → Entry                    │
│    - CRUD operations via DAO                                             │
│                                                                          │
│  Responsibilities:                                                       │
│    - Business logic enforcement                                          │
│    - Data transformation                                                 │
│    - Validation rules                                                    │
│    - Domain-specific operations                                          │
└────────────────────────────────┬────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                          DAO LAYER                                       │
├─────────────────────────────────────────────────────────────────────────┤
│  Interface: DaoProvider[EN Entity, I ID]                                 │
│                                                                          │
│  Methods:                                                                │
│    - FindById(ctx, id) (*EN, error)                                     │
│    - FindByIds(ctx, ids) ([]EN, error)                                  │
│    - Create(ctx, entity) (*EN, error)                                   │
│    - Update(ctx, id, entity) (*EN, error)                               │
│    - Search(ctx, query, sort, page, size) ([]EN, int64, error)          │
│    - SearchJoins(ctx, query, sort, page, size, joins) ([]EN, int64)     │
│    - GetSession(ctx) *gorm.DB                                           │
│    - AddPredicateForSearchJoins(req, filter) *gorm.DB                   │
│                                                                          │
│  Implementations:                                                        │
│    - PostgreSQL DAO (OLTP - transactional)                              │
│    - ClickHouse DAO (OLAP - analytics)                                  │
└─────────────────────────────────────────────────────────────────────────┘
```

### Generic Type Implementation

```go
// core/rest/service.go
type Service[RES domain.Response, EX domain.Entry, EN domain.Entity, I domain.ID] struct {
    Manager    *Manager[EX, EN, I]
    Repository DaoProvider[EN, I]
}

// core/rest/manager.go
type Manager[EX domain.Entry, EN domain.Entity, I domain.ID] struct {
    Dao       DaoProvider[EN, I]
    ToEntity  func(*EX) (*EN, error)
    ToEntry   func(*EN) (*EX, error)
}

// core/rest/dao.go
type DaoProvider[EN domain.Entity, I domain.ID] interface {
    FindById(ctx context.Context, id I) (*EN, error)
    Create(ctx context.Context, entity *EN) (*EN, error)
    Update(ctx context.Context, id I, entity *EN) (*EN, error)
    Search(ctx context.Context, query interface{}, sortBy, sortDir string, page, size int) ([]EN, int64, error)
    GetSession(ctx context.Context) *gorm.DB
}
```

---

## 3. DATABASE SCHEMA (27+ Tables)

### schema.sql Structure (918 lines)

#### Social Profile Tables

```sql
-- Instagram profiles with 100+ columns
CREATE TABLE instagram_account (
    id SERIAL PRIMARY KEY,
    profile_id VARCHAR(255) UNIQUE,
    handle VARCHAR(255),
    full_name VARCHAR(255),
    bio TEXT,

    -- Metrics
    followers BIGINT,
    following BIGINT,
    engagement_rate DECIMAL(10,6),
    avg_likes DECIMAL(10,2),
    avg_comments DECIMAL(10,2),
    avg_views DECIMAL(10,2),
    total_posts INTEGER,

    -- Audience Demographics (JSONB)
    audience_gender JSONB,
    audience_age JSONB,
    audience_location JSONB,

    -- Grades
    engagement_rate_grade VARCHAR(10),
    followers_grade VARCHAR(10),

    -- Linked Accounts
    youtube_account_id INTEGER REFERENCES youtube_account(id),
    winkl_profile_id VARCHAR(255),
    gcc_profile_id VARCHAR(255),

    -- Timestamps
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- YouTube channels
CREATE TABLE youtube_account (
    id SERIAL PRIMARY KEY,
    channel_id VARCHAR(255) UNIQUE,
    title VARCHAR(255),
    subscribers BIGINT,
    total_views BIGINT,
    avg_views DECIMAL(10,2),
    plays BIGINT,
    shorts_reach BIGINT,
    -- Similar structure to Instagram
);
```

#### Collection Tables

```sql
-- Profile collections
CREATE TABLE profile_collection (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    partner_id BIGINT NOT NULL,
    account_id BIGINT,
    user_id BIGINT,
    share_id VARCHAR(255) UNIQUE,
    source VARCHAR(50),  -- SAAS, SAAS-AT, GCC_CAMPAIGN
    is_public BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Collection items (profiles in collection)
CREATE TABLE profile_collection_item (
    id SERIAL PRIMARY KEY,
    collection_id INTEGER REFERENCES profile_collection(id),
    profile_id INTEGER,
    platform VARCHAR(50),
    platform_profile_id VARCHAR(255),
    added_by BIGINT,
    custom_data JSONB,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Post collections
CREATE TABLE post_collection (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    partner_id BIGINT NOT NULL,
    account_id BIGINT,
    ingestion_frequency VARCHAR(50),
    show_in_report BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW()
);
```

#### Analytics Tables

```sql
-- Leaderboard rankings
CREATE TABLE leaderboard (
    id SERIAL PRIMARY KEY,
    platform VARCHAR(50),
    profile_id VARCHAR(255),
    month DATE,

    -- Rankings
    followers_rank INTEGER,
    followers_change_rank INTEGER,
    engagement_rate_rank INTEGER,
    views_rank INTEGER,
    plays_rank INTEGER,

    -- Category/Language rankings
    followers_rank_by_cat INTEGER,
    followers_rank_by_lang INTEGER,
    followers_rank_by_cat_lang INTEGER,

    category VARCHAR(100),
    language VARCHAR(50),

    UNIQUE(platform, profile_id, month)
);
CREATE INDEX idx_leaderboard_month ON leaderboard(month);
CREATE INDEX idx_leaderboard_platform ON leaderboard(platform);

-- Time-series profile data
CREATE TABLE social_profile_time_series (
    id SERIAL PRIMARY KEY,
    platform VARCHAR(50),
    platform_profile_id VARCHAR(255),
    date DATE,
    followers BIGINT,
    following BIGINT,
    engagement_rate DECIMAL(10,6),
    posts_count INTEGER
);
CREATE INDEX idx_ts_platform_profile ON social_profile_time_series(platform, platform_profile_id, date);

-- Post metrics summary
CREATE TABLE collection_post_metrics_summary (
    id SERIAL PRIMARY KEY,
    collection_id INTEGER,
    post_id VARCHAR(255),
    platform VARCHAR(50),
    likes BIGINT,
    comments BIGINT,
    views BIGINT,
    shares BIGINT,
    engagement_rate DECIMAL(10,6),
    last_updated TIMESTAMP
);
```

#### Business Tables

```sql
-- Partner usage tracking
CREATE TABLE partner_usage (
    id SERIAL PRIMARY KEY,
    partner_id BIGINT NOT NULL,
    module VARCHAR(100),
    usage_count INTEGER DEFAULT 0,
    limit_count INTEGER,
    plan_type VARCHAR(50),
    valid_from DATE,
    valid_to DATE
);

-- Activity tracking
CREATE TABLE activity_tracker (
    id SERIAL PRIMARY KEY,
    partner_id BIGINT,
    account_id BIGINT,
    user_id BIGINT,
    activity_type VARCHAR(100),
    entity_type VARCHAR(100),
    entity_id VARCHAR(255),
    metadata JSONB,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Profile page access tracking (for paid plans)
CREATE TABLE partner_profile_page_track (
    id SERIAL PRIMARY KEY,
    partner_id BIGINT,
    platform VARCHAR(50),
    platform_profile_id VARCHAR(255),
    access_count INTEGER DEFAULT 1,
    first_accessed TIMESTAMP,
    last_accessed TIMESTAMP,
    UNIQUE(partner_id, platform, platform_profile_id)
);
```

---

## 4. MULTI-TENANT IMPLEMENTATION

### Request Context Structure

```go
// core/appcontext/requestcontext.go
type RequestContext struct {
    Ctx             context.Context
    Mutex           sync.Mutex
    Properties      map[string]interface{}

    // Tenant Identification
    PartnerId       *int64           // Primary tenant ID
    AccountId       *int64           // Sub-tenant/account
    UserId          *int64           // Individual user
    UserName        *string

    // Authentication
    Authorization   string           // Bearer token
    DeviceId        *string
    IsLoggedIn      bool

    // Plan & Access Control
    IsPremiumMember bool
    IsBasicMember   bool
    PlanType        *constants.PlanType  // FREE, SAAS, PAID

    // Database Sessions
    Session         persistence.Session  // PostgreSQL transaction
    CHSession       persistence.Session  // ClickHouse transaction
}
```

### Header Extraction

```go
// server/middlewares/context.go
func ApplicationContext(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx := appcontext.NewRequestContext(r.Context())

        // Extract tenant headers
        if partnerId := r.Header.Get("x-bb-partner-id"); partnerId != "" {
            pid, _ := strconv.ParseInt(partnerId, 10, 64)
            ctx.PartnerId = &pid
        }

        if accountId := r.Header.Get("x-bb-account-id"); accountId != "" {
            aid, _ := strconv.ParseInt(accountId, 10, 64)
            ctx.AccountId = &aid
        }

        // Plan type
        if planType := r.Header.Get("x-bb-plan-type"); planType != "" {
            pt := constants.PlanType(planType)
            ctx.PlanType = &pt
        }

        // Authorization token parsing
        if auth := r.Header.Get("Authorization"); auth != "" {
            // Format: userId~userName~isLoggedIn~accountId:password (base64)
            ctx.ParseAuthorization(auth)
        }

        next.ServeHTTP(w, r.WithContext(ctx.ToContext()))
    })
}
```

### Plan-Based Feature Gating

```go
// constants/constants.go
const (
    FreePlan PlanType = "FREE"    // Limited features
    SaasPlan PlanType = "SAAS"    // Standard SaaS
    PaidPlan PlanType = "PAID"    // Enterprise
)

var PartnerLimitModules = map[string]bool{
    "DISCOVERY":        true,
    "CAMPAIGN_REPORT":  true,
    "PROFILE_PAGE":     true,
    "ACCOUNT_TRACKING": true,
}

// Security checks for free users
func blockProfileDataForFreeUsers(ctx *RequestContext, profile *Profile) {
    if ctx.PlanType != nil && *ctx.PlanType == FreePlan {
        // Nullify premium fields
        profile.Email = nil
        profile.Phone = nil
        profile.AudienceDetails = nil
    }
}
```

---

## 5. MESSAGE QUEUE INTEGRATION (Watermill + AMQP)

### AMQP Configuration

```go
// amqp/config.go
func GetAMQPConfig() amqp.Config {
    return amqp.Config{
        Connection: amqp.ConnectionConfig{
            AmqpURI: viper.GetString("AMQP_URI"),
        },
        Marshaler: amqp.DefaultMarshaler{},
        Exchange: amqp.ExchangeConfig{
            GenerateName: func(topic string) string {
                // Pattern: {exchange}___{queue}
                parts := strings.Split(topic, "___")
                return parts[0]
            },
            Type:    "topic",
            Durable: true,
        },
        Queue: amqp.QueueConfig{
            GenerateName: func(topic string) string {
                parts := strings.Split(topic, "___")
                if len(parts) > 1 {
                    return parts[1]
                }
                return topic
            },
            Durable: true,
        },
        QueueBind: amqp.QueueBindConfig{
            GenerateRoutingKey: func(topic string) string {
                parts := strings.Split(topic, "___")
                if len(parts) > 1 {
                    return parts[1]
                }
                return topic
            },
        },
    }
}
```

### Event Listeners Setup

```go
// listeners/listeners.go
func SetupListeners(container *app.ApplicationContainer) {
    subscriber, _ := amqp.NewSubscriber(GetAMQPConfig(), logger)

    router, _ := message.NewRouter(message.RouterConfig{}, logger)

    // Middleware chain
    router.AddMiddleware(
        middleware.MessageApplicationContext,      // Extract context
        middleware.TransactionSessionHandler,      // Transaction management
        middleware.Retry{
            MaxRetries:      3,
            InitialInterval: 100 * time.Millisecond,
        },
        middleware.Recoverer,                      // Panic recovery
    )

    // Register handlers for each module
    router.AddNoPublisherHandler(
        "discovery_handler",
        "coffee.dx___discovery_events_q",
        subscriber,
        container.Discovery.Listeners.HandleEvent,
    )

    router.AddNoPublisherHandler(
        "profile_collection_handler",
        "coffee.dx___profile_collection_q",
        subscriber,
        container.ProfileCollection.Listeners.HandleEvent,
    )

    // ... more handlers for each module

    router.Run(context.Background())
}
```

### Event Publishing

```go
// publishers/publisher.go
type Publisher struct {
    pub *amqp.Publisher
}

func (p *Publisher) PublishMessage(jsonBytes []byte, topic string) error {
    msg := message.NewMessage(watermill.NewUUID(), jsonBytes)
    return p.pub.Publish(topic, msg)
}

// Usage in service layer
func (s *Service) Create(ctx context.Context, entry *Entry) (*Response, error) {
    // ... create logic ...

    // After-commit callback for event publishing
    session.PerformAfterCommit(ctx, func() {
        eventData, _ := json.Marshal(entry)
        publisher.PublishMessage(eventData, "coffee.dx___profile_collection_q")
    })

    return response, nil
}
```

---

## 6. DUAL DATABASE STRATEGY

### PostgreSQL (OLTP - Transactional)

```go
// core/persistence/postgres/db.go
var (
    db     *gorm.DB
    pgOnce sync.Once
)

func GetDB() *gorm.DB {
    pgOnce.Do(func() {
        dsn := fmt.Sprintf(
            "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Kolkata",
            viper.GetString("PG_HOST"),
            viper.GetString("PG_USER"),
            viper.GetString("PG_PASS"),
            viper.GetString("PG_DB"),
            viper.GetString("PG_PORT"),
        )

        db, _ = gorm.Open(postgres.Open(dsn), &gorm.Config{
            Logger: logger.Default.LogMode(logger.Info),
        })

        sqlDB, _ := db.DB()
        sqlDB.SetMaxIdleConns(viper.GetInt("PG_MAX_IDLE_CONN"))  // 5
        sqlDB.SetMaxOpenConns(viper.GetInt("PG_MAX_OPEN_CONN"))  // 5
    })
    return db
}
```

### ClickHouse (OLAP - Analytics)

```go
// core/persistence/clickhouse/db.go
var (
    chDB     *gorm.DB
    chOnce   sync.Once
)

func GetDB() *gorm.DB {
    chOnce.Do(func() {
        dsn := fmt.Sprintf(
            "tcp://%s:%s?database=%s&username=%s&password=%s&read_timeout=10&write_timeout=20",
            viper.GetString("CH_HOST"),
            viper.GetString("CH_PORT"),
            viper.GetString("CH_DB"),
            viper.GetString("CH_USER"),
            viper.GetString("CH_PASS"),
        )

        chDB, _ = gorm.Open(clickhouse.Open(dsn), &gorm.Config{
            Logger: logger.Default.LogMode(logger.Info),
        })

        sqlDB, _ := chDB.DB()
        sqlDB.SetMaxIdleConns(viper.GetInt("CH_MAX_IDLE_CONN"))  // 1
        sqlDB.SetMaxOpenConns(viper.GetInt("CH_MAX_OPEN_CONN"))  // 1
    })
    return chDB
}
```

### Session Management

```go
// server/middlewares/session.go
func ServiceSessionMiddlewares(serviceWrappers map[ServiceWrapper]ApiWrapper) func(next http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            ctx := appcontext.GetRequestContext(r.Context())

            // Initialize PostgreSQL session
            pgSession := postgres.NewSession()
            ctx.Session = pgSession

            // Initialize ClickHouse session (if needed)
            if needsClickHouse(r.URL.Path) {
                chSession := clickhouse.NewSession()
                ctx.CHSession = chSession
            }

            // Timeout context
            timeoutCtx, cancel := context.WithTimeout(r.Context(),
                time.Duration(viper.GetInt("SERVER_TIMEOUT")) * time.Second)
            defer cancel()

            // Execute handler
            done := make(chan bool)
            go func() {
                next.ServeHTTP(w, r.WithContext(ctx.ToContext()))
                done <- true
            }()

            select {
            case <-done:
                // Commit on success
                pgSession.Commit(ctx)
                if ctx.CHSession != nil {
                    ctx.CHSession.Commit(ctx)
                }
                // Execute after-commit callbacks
                pgSession.ExecuteAfterCommitCallbacks()

            case <-timeoutCtx.Done():
                // Rollback on timeout
                pgSession.Rollback(ctx)
                if ctx.CHSession != nil {
                    ctx.CHSession.Rollback(ctx)
                }
                http.Error(w, "Request timeout", http.StatusGatewayTimeout)
            }
        })
    }
}
```

---

## 7. BUSINESS MODULES (12 Modules)

### Module Structure Pattern

Each module follows the same structure:
```
{module}/
├── api/api.go              # HTTP handlers
├── service/service.go      # Business orchestration
├── manager/manager.go      # Domain logic
├── dao/dao.go              # Data access
├── domain/
│   ├── entry.go            # Input DTOs
│   ├── entity.go           # Database entities
│   └── response.go         # Output DTOs
└── listeners/listeners.go  # Event handlers
```

### 1. Discovery Module

```
Endpoints:
  GET  /discovery-service/api/profile/{profileId}
  GET  /discovery-service/api/profile/{platform}/{platformProfileId}
  GET  /discovery-service/api/profile/{platform}/{platformProfileId}/audience
  POST /discovery-service/api/profile/{platform}/{platformProfileId}/timeseries
  POST /discovery-service/api/profile/{platform}/{platformProfileId}/content
  GET  /discovery-service/api/profile/{platform}/{platformProfileId}/hashtags
  POST /discovery-service/api/profile/{platform}/{platformProfileId}/similar_accounts
  POST /discovery-service/api/profile/{platform}/search
  POST /discovery-service/api/profile/locations

Managers:
  - SearchManager: Profile search with filters
  - TimeSeriesManager: Growth data over time
  - HashtagsManager: Hashtag analytics
  - AudienceManager: Demographic data
  - LocationManager: Location-based search
```

### 2. Profile Collection Module

```
Endpoints:
  POST   /profile-collection-service/api/collection/
  GET    /profile-collection-service/api/collection/{id}
  PUT    /profile-collection-service/api/collection/{id}
  DELETE /profile-collection-service/api/collection/{id}
  POST   /profile-collection-service/api/collection/search
  GET    /profile-collection-service/api/collection/byshareid/{shareId}
  POST   /profile-collection-service/api/collection/{id}/link/renew
  GET    /profile-collection-service/api/collection/recent

  POST   /profile-collection-service/api/collection/{collectionId}/item/
  DELETE /profile-collection-service/api/collection/{collectionId}/item
  PUT    /profile-collection-service/api/collection/{collectionId}/item/bulk
  POST   /profile-collection-service/api/collection/item/search
```

### 3. Leaderboard Module

```
Endpoints:
  POST /leaderboard-service/api/leaderboard/platform/{platform}
  POST /leaderboard-service/api/leaderboard/platform/{platform}/category/{category}
  POST /leaderboard-service/api/leaderboard/platform/{platform}/language/{language}
  POST /leaderboard-service/api/leaderboard/cross-platform
```

### 4-12. Other Modules

| Module | Purpose | Key Features |
|--------|---------|--------------|
| Post Collection | Curated post collections | Ingestion frequency, metrics |
| Collection Analytics | Hashtag & keyword tracking | Post metrics summary |
| Genre Insights | Genre-based analytics | Trending content |
| Keyword Collection | Keyword-based collections | Search optimization |
| Campaign Profiles | Campaign management | Profile associations |
| Collection Group | Grouped collections | Hierarchical organization |
| Partner Usage | Usage tracking | Plan limits, feature access |
| Content | Content management | Media handling |

---

## 8. API RESPONSE FORMAT

### Standard Response Structure

```go
// domain/response.go
type StandardResponse struct {
    Data       interface{} `json:"data"`
    Count      int         `json:"count"`
    TotalCount int64       `json:"totalCount"`
    Status     Status      `json:"status"`
}

type Status struct {
    Status     string `json:"status"`      // SUCCESS, ERROR
    Message    string `json:"message"`
    NextCursor string `json:"nextCursor,omitempty"`
}

// Example response
{
    "data": [...],
    "count": 10,
    "totalCount": 1000,
    "status": {
        "status": "SUCCESS",
        "message": "Records retrieved successfully",
        "nextCursor": "11"
    }
}
```

### Request Headers

```
Authorization: Bearer <base64_encoded_token>
x-bb-partner-id: {int64}
x-bb-account-id: {int64}
x-bb-plan-type: FREE|SAAS|PAID
x-bb-uid: {string}
x-bb-clientid: {string}
x-bb-requestid: {string}
```

---

## 9. CI/CD PIPELINE

### GitLab CI Configuration

```yaml
stages:
  - test
  - build
  - deploy_staging
  - deploy_prod

build:
  stage: build
  variables:
    GOOS: linux
    GOARCH: amd64
  script:
    - go build -o bin/coffee
    - tar -czvf coffee.tar.gz .
  artifacts:
    paths:
      - bin/coffee
      - .env*
      - scripts/
      - coffee.tar.gz
    expire_in: 1 week
  only:
    - master
    - dev

deploy_staging:
  stage: deploy_staging
  script:
    - /bin/bash scripts/start.sh
  tags:
    - deployer-stage-search
  when: manual
  only:
    - master
    - dev

deploy_prod_1:
  stage: deploy_prod
  script:
    - /bin/bash scripts/start.sh
  tags:
    - cb1-1
  when: manual
  only:
    - master

deploy_prod_2:
  stage: deploy_prod
  script:
    - /bin/bash scripts/start.sh
  tags:
    - cb2-1
  when: manual
  only:
    - master
```

---

## 10. OBSERVABILITY

### Prometheus Metrics

```go
// Chi-Prometheus integration
chiprometheus.NewMiddleware(
    "coffee",
    chiprometheus.WithBuckets([]float64{300, 500, 1000, 5000, 10000, 20000, 30000}),
)

// Exposed at /metrics
```

### Sentry Error Tracking

```go
// sentry/sentry.go
func Setup() {
    sentry.Init(sentry.ClientOptions{
        Dsn:              viper.GetString("SENTRY_DSN"),
        Environment:      viper.GetString("ENV"),
        AttachStacktrace: true,
    })
}

// Middleware captures errors with request context
func SentryErrorLoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        hub := sentry.GetHubFromContext(r.Context())
        hub.Scope().SetRequest(r)

        defer func() {
            if err := recover(); err != nil {
                hub.RecoverWithContext(r.Context(), err)
            }
        }()

        next.ServeHTTP(w, r)
    })
}
```

### Logging

```go
// logger/logger.go
func Setup() {
    logrus.SetFormatter(&logrus.JSONFormatter{})

    level, _ := logrus.ParseLevel(viper.GetString("LOG_LEVEL"))
    logrus.SetLevel(level)

    // Sentry hook for errors
    hook, _ := logrus_sentry.NewSentryHook(viper.GetString("SENTRY_DSN"), []logrus.Level{
        logrus.PanicLevel,
        logrus.FatalLevel,
        logrus.ErrorLevel,
    })
    logrus.AddHook(hook)
}
```

---

## 11. KEY METRICS & STATISTICS

| Metric | Value |
|--------|-------|
| **Total Lines of Code** | 8,500+ |
| **Go Source Files** | 100+ |
| **Database Tables** | 27+ |
| **API Endpoints** | 50+ |
| **Business Modules** | 12 |
| **Middleware Layers** | 10 |
| **Go Version** | 1.18 |
| **Schema SQL Lines** | 918 |

---

## 12. SKILLS DEMONSTRATED

### Technical Skills

| Skill | Evidence |
|-------|----------|
| **Go Generics** | Type-safe Service/Manager/DAO with generic constraints |
| **REST API Design** | 4-layer architecture with consistent patterns |
| **Multi-Tenancy** | Partner/Account isolation with plan-based gating |
| **Dual Database** | PostgreSQL (OLTP) + ClickHouse (OLAP) strategy |
| **Message Queues** | Watermill + AMQP for event-driven architecture |
| **Transaction Management** | Request-scoped transactions with auto-commit/rollback |
| **Middleware Pipeline** | 10-layer middleware chain with timeout handling |
| **Observability** | Prometheus, Sentry, structured logging |

### Architecture Patterns

1. **4-Layer REST Architecture** - API → Service → Manager → DAO
2. **Generic Repository Pattern** - Type-safe CRUD with Go generics
3. **Service Container** - Dependency injection container
4. **Event-Driven** - Watermill pub/sub with after-commit callbacks
5. **Multi-Tenant** - Header-based tenant isolation
6. **Feature Flags** - Plan-based feature gating

---

## 13. INTERVIEW TALKING POINTS

### System Design Questions

**"Describe a multi-tenant SaaS architecture you built"**
- 4-layer REST architecture with Go generics for type safety
- Partner/Account/User hierarchy with header-based isolation
- Plan-based feature gating (FREE/SAAS/PAID)
- Activity tracking and usage limits per tenant

**"How do you handle different database needs?"**
- PostgreSQL for OLTP (collections, profiles, relationships)
- ClickHouse for OLAP (time-series, analytics, aggregations)
- Request-scoped transactions with dual session management
- Automatic commit/rollback with timeout handling

**"Explain your event-driven architecture"**
- Watermill framework with AMQP transport
- After-commit callbacks for guaranteed event publishing
- Retry middleware (3 attempts, 100ms interval)
- Separate queues per business module

### Behavioral Questions

**"Tell me about a complex backend system you built"**
- Coffee: 8,500+ LOC, 12 business modules, 50+ endpoints
- Generic REST framework used across all modules
- Multi-tenant with plan-based access control
- Dual database strategy for different workloads

**"How do you ensure code quality and consistency?"**
- Generic types enforce consistent patterns
- Middleware chain ensures security at all endpoints
- Transaction management prevents data inconsistency
- Prometheus metrics for performance monitoring

---

*Generated through comprehensive source code analysis of the coffee project.*
