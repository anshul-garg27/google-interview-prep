# ULTRA-DEEP ANALYSIS: SAAS-GATEWAY API GATEWAY

## PROJECT OVERVIEW

| Attribute | Value |
|-----------|-------|
| **Project Name** | SaaS Gateway |
| **Purpose** | API Gateway for Bulbul Creator Platform (goodcreator.co) |
| **Architecture** | Reverse proxy with 7-layer middleware pipeline |
| **Framework** | Gin v1.8.1 (Go) |
| **Port** | 8009 (main), 6069 (pprof debug) |
| **Services Proxied** | 13 microservices |
| **Binary Size** | 23MB |

---

## 1. COMPLETE DIRECTORY STRUCTURE

```
/saas-gateway/
├── .git/                              # Git repository
├── .gitignore                         # Git ignore rules
├── .gitlab-ci.yml                     # CI/CD pipeline (42 lines)
│
├── .env                               # Default environment
├── .env.local                         # Local development config
├── .env.stage                         # Staging environment config
├── .env.production                    # Production environment config
│
├── go.mod                             # Go module definition
├── go.sum                             # Dependency checksums
├── main.go                            # Application entry point
├── saas-gateway                       # Compiled binary (23MB)
│
├── api/                               # API Handlers
│   └── heartbeat/
│       └── heartbeat.go               # Health check endpoints
│
├── cache/                             # Caching Layer
│   ├── redis.go                       # Redis cluster client (singleton)
│   └── ristretto.go                   # In-memory cache (singleton)
│
├── client/                            # HTTP Client Layer
│   ├── client.go                      # Base HTTP client (resty wrapper)
│   │
│   ├── entry/                         # Domain Models
│   │   ├── identity.go                # User, Client, Account entities
│   │   ├── status.go                  # Response status model
│   │   └── asset.go                   # Asset information model
│   │
│   ├── identity/                      # Identity Service Client
│   │   ├── identity.go                # Auth API methods
│   │   └── identity_test.go           # Unit tests
│   │
│   ├── input/                         # Request DTOs
│   │   ├── commons.go                 # Common input types
│   │   ├── filter.go                  # Filter parameters
│   │   └── identity.go                # Auth input types
│   │
│   ├── partner/                       # Partner Service Client
│   │   ├── partner.go                 # Partner API methods
│   │   └── partner_test.go            # Unit tests
│   │
│   └── response/                      # Response DTOs
│       ├── abc.go                     # Generic responses
│       ├── asset.go                   # Asset responses
│       ├── identity.go                # Auth responses
│       ├── partner.go                 # Partner responses
│       └── status.go                  # Status responses
│
├── config/                            # Configuration
│   └── config.go                      # Config struct and loader
│
├── context/                           # Request Context
│   └── context.go                     # Gateway context utilities
│
├── custom/                            # Custom Types
│   └── error.go                       # Custom error definitions
│
├── generator/                         # Utilities
│   └── requestid.go                   # UUID request ID generator
│
├── handler/                           # HTTP Handlers
│   └── saas/
│       └── saas.go                    # Reverse proxy handler
│
├── header/                            # Constants
│   └── header.go                      # HTTP header constants (25+ headers)
│
├── locale/                            # Localization
│   ├── locale/
│   │   └── locale.go                  # Locale utilities
│   └── localeconfig/
│       └── locale.go                  # Locale configuration
│
├── logger/                            # Logging
│   └── logger.go                      # Logger setup (logrus + sentry)
│
├── metrics/                           # Observability
│   └── metrics.go                     # Prometheus metrics collectors
│
├── middleware/                        # Middleware Layer
│   ├── auth.go                        # JWT authentication (~150 lines)
│   ├── requestid.go                   # Request ID injection
│   └── requestlogger.go               # Request/response logging
│
├── route/                             # Route Definitions
│   ├── heartbeatRoutes.go             # Health check routes
│   └── route.go                       # Route types
│
├── router/                            # Router Setup
│   └── router.go                      # Main router with middleware (~93 lines)
│
├── safego/                            # Safe Concurrency
│   └── safego.go                      # Panic-safe goroutine wrapper
│
├── scripts/                           # Deployment
│   └── start.sh                       # Graceful deployment script
│
├── sentry/                            # Error Tracking
│   └── sentry.go                      # Sentry initialization
│
└── util/                              # Utilities
    ├── transformer.go                 # Data transformers
    └── util.go                        # Helper functions
```

---

## 2. TECHNOLOGY STACK

### Core Dependencies
| Category | Technology | Version | Purpose |
|----------|-----------|---------|---------|
| **Web Framework** | Gin | v1.8.1 | HTTP routing & middleware |
| **JWT Auth** | golang-jwt/jwt | v3.2.0 | Token validation |
| **HTTP Client** | go-resty/resty | v2.3.0 | Upstream service calls |
| **In-Memory Cache** | dgraph-io/ristretto | v0.0.3 | Local LFU caching (10M keys) |
| **Distributed Cache** | go-redis | v7.0.0-beta.6 | Redis cluster client |
| **Logging** | logrus + zerolog | v1.9.0 / v1.18.0 | Structured logging |
| **Error Tracking** | getsentry/sentry-go | v0.20.0 | Error monitoring |
| **Metrics** | prometheus/client_golang | v1.11.1 | Observability |
| **Config** | joho/godotenv | v1.3.0 | Environment management |
| **CORS** | rs/cors | v1.7.0 | Cross-origin configuration |
| **UUID** | google/uuid | v1.3.0 | Request ID generation |
| **Config Mgmt** | spf13/viper | v1.15.0 | Advanced configuration |

### Go Module
```go
module init.bulbul.tv/bulbul-backend/saas-gateway
go 1.12
```

---

## 3. GATEWAY ARCHITECTURE

### Request Flow
```
┌─────────────────────────────────────────────────────────────────────┐
│                     SAAS GATEWAY ARCHITECTURE                        │
└─────────────────────────────────────────────────────────────────────┘

INTERNET / MOBILE APP / WEB
            ↓
┌─────────────────────────────────────────────────────────────────────┐
│                    LOAD BALANCER                                     │
│  ┌─────────────────┐  ┌─────────────────┐                          │
│  │   Node 1        │  │   Node 2        │                          │
│  │   cb1-1:8009    │  │   cb2-1:8009    │                          │
│  └─────────────────┘  └─────────────────┘                          │
└─────────────────────────────────────────────────────────────────────┘
            ↓
┌─────────────────────────────────────────────────────────────────────┐
│                    SaaS GATEWAY (Port 8009)                          │
├─────────────────────────────────────────────────────────────────────┤
│ MIDDLEWARE PIPELINE (Executed in Order):                             │
│                                                                      │
│  1. gin.Recovery()           ← Panic recovery                       │
│  2. sentrygin.New()          ← Error tracking                       │
│  3. cors.New()               ← CORS (41 origins)                    │
│  4. GatewayContextMiddleware ← Context injection                    │
│  5. RequestIdMiddleware()    ← UUID generation (x-bb-requestid)     │
│  6. RequestLogger()          ← Request/response logging             │
│  7. AppAuth()                ← JWT + Redis session validation       │
│                                                                      │
├─────────────────────────────────────────────────────────────────────┤
│ ROUTE HANDLERS:                                                      │
│                                                                      │
│  GET  /metrics               ← Prometheus metrics                   │
│  GET  /heartbeat/            ← Health check (200 OK / 410 Gone)     │
│  PUT  /heartbeat/?beat=      ← Health state toggle                  │
│  ANY  /{service}/*           ← Reverse proxy to downstream          │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
            ↓
┌─────────────────────────────────────────────────────────────────────┐
│                    CACHING LAYER                                     │
│                                                                      │
│  ┌────────────────────────┐  ┌────────────────────────┐            │
│  │ LAYER 1: Ristretto     │  │ LAYER 2: Redis Cluster │            │
│  │ (In-Memory LFU)        │  │ (Distributed)          │            │
│  │                        │  │                        │            │
│  │ - 10M key capacity     │  │ - 3-6 nodes            │            │
│  │ - 1GB max memory       │  │ - 100 pool size        │            │
│  │ - Nanosecond lookups   │  │ - session:{id} keys    │            │
│  │ - Per-instance         │  │ - Shared across nodes  │            │
│  └────────────────────────┘  └────────────────────────┘            │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
            ↓
┌─────────────────────────────────────────────────────────────────────┐
│                    EXTERNAL SERVICES                                 │
│                                                                      │
│  ┌────────────────────────┐  ┌────────────────────────┐            │
│  │ Identity Service       │  │ Partner Service        │            │
│  │ /auth/verify/token     │  │ /partner/{id}          │            │
│  │ /auth/login            │  │                        │            │
│  └────────────────────────┘  └────────────────────────┘            │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
            ↓
┌─────────────────────────────────────────────────────────────────────┐
│                    DOWNSTREAM MICROSERVICES                          │
│                                                                      │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │                     SAAS_URL (Primary)                       │   │
│  │  discovery-service      │ leaderboard-service               │   │
│  │  profile-collection     │ post-collection-service           │   │
│  │  activity-service       │ collection-analytics-service      │   │
│  │  genre-insights-service │ content-service                   │   │
│  │  collection-group       │ keyword-collection-service        │   │
│  │  partner-usage-service  │ campaign-profile-service          │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                                                                      │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │                    SAAS_DATA_URL (Data)                      │   │
│  │  social-profile-service  (Uses different backend URL)       │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

### Proxied Services (13 Total)

| Service | Route Pattern | Backend |
|---------|--------------|---------|
| discovery-service | `/discovery-service/*` | SAAS_URL |
| leaderboard-service | `/leaderboard-service/*` | SAAS_URL |
| profile-collection-service | `/profile-collection-service/*` | SAAS_URL |
| activity-service | `/activity-service/*` | SAAS_URL |
| collection-analytics-service | `/collection-analytics-service/*` | SAAS_URL |
| post-collection-service | `/post-collection-service/*` | SAAS_URL |
| genre-insights-service | `/genre-insights-service/*` | SAAS_URL |
| content-service | `/content-service/*` | SAAS_URL |
| **social-profile-service** | `/social-profile-service/*` | **SAAS_DATA_URL** |
| collection-group-service | `/collection-group-service/*` | SAAS_URL |
| keyword-collection-service | `/keyword-collection-service/*` | SAAS_URL |
| partner-usage-service | `/partner-usage-service/*` | SAAS_URL |
| campaign-profile-service | `/campaign-profile-service/*` | SAAS_URL |

---

## 4. MIDDLEWARE PIPELINE

### Execution Order

```go
// router/router.go (Lines 18-93)

func SetupRouter(config config.Config) *gin.Engine {
    router := gin.New()

    // 1. Panic Recovery
    router.Use(gin.Recovery())

    // 2. Sentry Error Tracking
    router.Use(sentrygin.New(sentrygin.Options{Repanic: true}))

    // 3. CORS (41 whitelisted origins)
    router.Use(cors.New(cors.Options{
        AllowedOrigins: []string{
            "http://staging.app.vidooly.com",
            "https://stage.cf-provider.goodcreator.co",
            "https://cf-provider.goodcreator.co",
            "https://www.instagram.com",
            "https://www.youtube.com",
            "https://suite.goodcreator.co",
            "https://goodcreator.co",
            "http://localhost:3000",
            "http://localhost:3001",
            "http://localhost:8298",
            // ... 31 more origins
        },
        AllowedMethods:   []string{"HEAD", "GET", "POST", "OPTIONS", "PUT", "PATCH", "DELETE"},
        AllowedHeaders:   []string{"*"},
        AllowCredentials: true,
    }))

    // 4. Gateway Context Injection
    router.Use(GatewayContextToContextMiddleware(config))

    // 5. Request ID Generation
    router.Use(middleware.RequestIdMiddleware())

    // 6. Request/Response Logging
    router.Use(middleware.RequestLogger(config))

    // Route Groups
    router.GET("/metrics", gin.WrapH(promhttp.Handler()))

    heartbeatRouter := router.Group("/heartbeat")
    groupRoutes(config, heartbeatRouter, route.HeartbeatRoutes)

    // 7. Authentication (per-route)
    router.Any("/discovery-service/*any", middleware.AppAuth(config), saas.ReverseProxy("discovery-service", config))
    // ... 12 more service routes

    return router
}
```

### Middleware Details

#### 1. Panic Recovery
```go
router.Use(gin.Recovery())
// Built-in Gin middleware
// Catches all panics, returns 500 Internal Server Error
// Prevents server crash from individual request failures
```

#### 2. Sentry Error Tracking
```go
router.Use(sentrygin.New(sentrygin.Options{Repanic: true}))
// Captures errors and sends to Sentry
// Repanic: true means panic is re-raised after capture
// Works with Recovery() middleware for complete coverage
```

#### 3. CORS Configuration
```go
cors.New(cors.Options{
    AllowedOrigins:   []string{...}, // 41 origins
    AllowedMethods:   []string{"HEAD", "GET", "POST", "OPTIONS", "PUT", "PATCH", "DELETE"},
    AllowedHeaders:   []string{"*"},
    AllowCredentials: true,
})

// Whitelisted Origins Include:
// - Production: cf-provider.goodcreator.co, suite.goodcreator.co
// - Staging: stage.cf-provider.goodcreator.co
// - Local: localhost:3000, localhost:3001, localhost:8298
// - External: instagram.com, youtube.com
```

#### 4. Gateway Context Middleware
```go
func GatewayContextToContextMiddleware(config config.Config) gin.HandlerFunc {
    return func(c *gin.Context) {
        gc := context.New(c, config)
        c.Set("gatewayContext", gc)
        c.Next()
    }
}

// context.New() does:
// - Generates request ID (UUID v4)
// - Blocks bots (User-Agent filtering)
// - Creates zerolog logger with request ID
// - Extracts whitelisted x-bb-* headers
// - Sets default locale (en)
```

#### 5. Request ID Middleware
```go
// generator/requestid.go
func SetNXRequestIdOnContext(c *gin.Context) {
    if c.Keys == nil {
        c.Keys = make(map[string]interface{})
    }

    if c.Keys[header.RequestID] == nil {
        uuid := uuid.New().String()
        c.Keys[header.RequestID] = uuid

        // Set in request header
        if c.Request != nil && c.Request.Header != nil {
            c.Request.Header.Set(header.RequestID, uuid)
        }

        // Set in response header
        if c.Writer != nil && c.Writer.Header() != nil {
            c.Writer.Header().Set(header.RequestID, uuid)
        }
    }
}
```

#### 6. Request Logger Middleware
```go
// middleware/requestlogger.go
func RequestLogger(config config.Config) gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path
        query := c.Request.URL.RawQuery

        gc := util.GatewayContextFromGinContext(c, config)

        // Capture request body (non-destructive read)
        buf, _ := ioutil.ReadAll(c.Request.Body)
        rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
        rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))
        body := readBody(rdr1)
        c.Request.Body = rdr2

        // Wrap response writer to capture response
        w := &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
        c.Writer = w

        c.Next()

        // Log based on LOG_LEVEL
        gc.Logger.Error().Msg(fmt.Sprintf(
            "%s - %s - [%s] \"%s %s %s %s %d %s %s\"\n%s\n%s\n",
            c.Request.Header.Get(header.RequestID),
            c.ClientIP(),
            time.Now().Format(time.RFC1123),
            c.Request.Method,
            path,
            c.Request.Header.Get(header.ApolloOpName),
            c.Request.Proto,
            c.Writer.Status(),
            "Response time: ", time.Now().Sub(start),
            c.Request.UserAgent(),
            c.Request.Header,
        ))
    }
}
```

#### 7. Authentication Middleware
```go
// middleware/auth.go (Lines 23-151)
func AppAuth(config config.Config) gin.HandlerFunc {
    return func(c *gin.Context) {
        gc := util.GatewayContextFromGinContext(c, config)
        verified := false

        // Step 1: Extract Client ID
        if clientId := gc.GetHeader(header.ClientID); clientId != "" {
            clientType := gc.GetHeader(header.ClientType)
            if clientType == "" {
                clientType = "CUSTOMER"
            }
            gc.GenerateClientAuthorizationWithClientId(clientId, clientType)
        }

        // Step 2: Skip auth for init device endpoints
        apolloOp := gc.Context.GetHeader(header.ApolloOpName)
        if apolloOp == "initDeviceGoMutation" || apolloOp == "initDeviceV2" ||
           apolloOp == "initDeviceV2Mutation" || apolloOp == "initDevice" ||
           strings.Contains(c.Request.URL.Path, "/api/auth/init") {
            verified = true
        }

        // Step 3: JWT Token Validation
        if authorization := gc.Context.GetHeader("Authorization"); authorization != "" {
            splittedAuthHeader := strings.Split(authorization, " ")
            if len(splittedAuthHeader) > 1 && splittedAuthHeader[1] != "" {

                // Parse JWT with HMAC-SHA256
                token, err := jwt.Parse(splittedAuthHeader[1], func(token *jwt.Token) (interface{}, error) {
                    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                        return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
                    }
                    hmac, _ := b64.StdEncoding.DecodeString(config.HMAC_SECRET)
                    return hmac, nil
                })

                if err == nil && token.Valid {
                    claims := token.Claims.(jwt.MapClaims)

                    // Step 4: Redis Session Lookup
                    if sessionId, ok := claims["sid"].(string); ok && sessionId != "" {
                        keyExists := cache.Redis(gc.Config).Exists("session:" + sessionId)
                        if exist, _ := keyExists.Result(); exist == 1 {
                            verified = true
                            // Extract user info from claims
                            gc.UserClientAccount = &entry.UserClientAccount{
                                UserId: claims["uid"].(int),
                                Id:     claims["userAccountId"].(int),
                            }
                        }
                    }

                    // Step 5: Fallback to Identity Service
                    if !verified {
                        verifyTokenResponse, err := identity.New(gc).VerifyToken(
                            &input.Token{Token: splittedAuthHeader[1]},
                            false,
                        )
                        if err == nil && verifyTokenResponse.UserClientAccount != nil {
                            verified = true
                            gc.UserClientAccount = verifyTokenResponse.UserClientAccount
                        }
                    }

                    // Step 6: Partner Plan Validation
                    if verified && gc.UserClientAccount != nil {
                        if gc.UserClientAccount.PartnerProfile != nil {
                            partnerId := strconv.FormatInt(
                                gc.UserClientAccount.PartnerProfile.PartnerId, 10)
                            partnerResponse, _ := partner.New(gc).FindPartnerById(partnerId)

                            for _, contract := range partnerResponse.Partners[0].Contracts {
                                if contract.ContractType == "SAAS" {
                                    currentTime := time.Now()
                                    startTime := time.Unix(0, contract.StartTime*int64(time.Millisecond))
                                    endTime := time.Unix(0, contract.EndTime*int64(time.Millisecond))

                                    if currentTime.After(startTime) && currentTime.Before(endTime) {
                                        gc.UserClientAccount.PartnerProfile.PlanType = contract.Plan
                                    } else {
                                        gc.UserClientAccount.PartnerProfile.PlanType = "FREE"
                                    }
                                }
                            }
                        }
                    }
                }
            }
        }

        // Step 7: Generate Authorization Headers
        gc.GenerateAuthorization()

        if !verified {
            c.AbortWithStatus(http.StatusUnauthorized)
            return
        }

        c.Next()
    }
}
```

---

## 5. CACHING ARCHITECTURE

### Two-Layer Caching Strategy

```
┌─────────────────────────────────────────────────────────────────────┐
│                    TWO-LAYER CACHING                                 │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  REQUEST → [Layer 1: Ristretto] → HIT? → Return                     │
│                      ↓ MISS                                          │
│            [Layer 2: Redis]     → HIT? → Return + Cache L1          │
│                      ↓ MISS                                          │
│            [Identity Service]   → Return + Cache L1 + Cache L2      │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

### Layer 1: Ristretto In-Memory Cache

```go
// cache/ristretto.go
var (
    singletonRistretto *ristretto.Cache
    ristrettoOnce      sync.Once
)

func Ristretto(config config.Config) *ristretto.Cache {
    ristrettoOnce.Do(func() {
        singletonRistretto, _ = ristretto.NewCache(&ristretto.Config{
            NumCounters: 1e7,     // 10,000,000 keys tracked for frequency
            MaxCost:     1 << 30, // 1GB maximum memory
            BufferItems: 64,      // 64 keys per Get buffer (batching)
        })
    })
    return singletonRistretto
}
```

**Characteristics:**
- **Capacity**: 10 million keys
- **Memory**: 1GB maximum
- **Eviction**: LFU (Least Frequently Used)
- **Scope**: Per-instance (not shared)
- **Latency**: Nanoseconds

### Layer 2: Redis Cluster

```go
// cache/redis.go
var (
    singletonRedis *redis.ClusterClient
    redisInit      sync.Once
)

func Redis(config config.Config) *redis.ClusterClient {
    redisInit.Do(func() {
        singletonRedis = redis.NewClusterClient(&redis.ClusterOptions{
            Addrs:    config.RedisClusterAddresses,
            PoolSize: 100,
            Password: config.REDIS_CLUSTER_PASSWORD,
        })
    })
    return singletonRedis
}
```

**Configuration by Environment:**

| Environment | Nodes | Addresses |
|-------------|-------|-----------|
| Production | 3 | bulbul-redis-prod-1:6379, -2:6380, -3:6381 |
| Staging | 3 | bulbul-redis-stage-1:6379, -2:6380, -3:6381 |
| Local | 6 | localhost:30001-30006 |

**Characteristics:**
- **Pool Size**: 100 connections
- **Key Format**: `session:{sessionId}`
- **Scope**: Shared across all instances
- **Latency**: Milliseconds

---

## 6. REVERSE PROXY IMPLEMENTATION

### Proxy Handler

```go
// handler/saas/saas.go
func ReverseProxy(module string, config config.Config) func(c *gin.Context) {
    return func(c *gin.Context) {
        appContext := util.GatewayContextFromGinContext(c, config)

        // Select backend URL
        saasUrl := config.SAAS_URL
        if module == "social-profile-service" {
            saasUrl = config.SAAS_DATA_URL  // Different backend for data service
        }

        // Parse target URL
        remote, err := url.Parse(saasUrl)
        if err != nil {
            panic(err)
        }

        // Extract auth info
        authHeader := appContext.UserAuthorization
        var partnerId string
        planType := "FREE"

        if appContext.UserClientAccount != nil {
            if appContext.UserClientAccount.PartnerProfile != nil {
                partnerId = strconv.Itoa(int(appContext.UserClientAccount.PartnerProfile.PartnerId))
                if appContext.UserClientAccount.PartnerProfile.PlanType != "" {
                    planType = appContext.UserClientAccount.PartnerProfile.PlanType
                }
            }
        }

        // Create reverse proxy
        proxy := httputil.NewSingleHostReverseProxy(remote)

        // Request transformation
        deviceId := appContext.XHeader.Get(header.Device)
        proxy.Director = func(req *http.Request) {
            req.Header = c.Request.Header                    // Copy all headers
            req.Header.Set("Authorization", authHeader)      // Override auth
            req.Header.Set(header.PartnerId, partnerId)      // Add partner ID
            req.Header.Set(header.PlanType, planType)        // Add plan type
            req.Header.Set(header.Device, deviceId)          // Add device ID
            req.Header.Del("accept-encoding")                // Remove encoding
            req.Host = remote.Host
            req.URL.Scheme = remote.Scheme
            req.URL.Host = remote.Host
            req.URL.Path = "/" + module + c.Param("any")     // Reconstruct path
        }

        // Response handling
        proxy.ModifyResponse = responseHandler(appContext, c.Request.Method,
            c.Param("any"), c.Request.Header.Get("origin"), ...)

        proxy.ServeHTTP(c.Writer, c.Request)
    }
}

func responseHandler(appContext *context.Context, method string, path string,
    origin string, setOrigin bool, originalAuth string) func(*http.Response) error {

    return func(resp *http.Response) error {
        resp.Header.Del("Access-Control-Allow-Origin")
        if setOrigin {
            resp.Header.Set("Access-Control-Allow-Origin", origin)
        }

        // Log non-200 responses
        if resp.Status != "200 OK" {
            log.Printf("SAAS-Bad-Response - %s, %s, %v : from : %s",
                method, path, resp.Status, originalAuth)
        }

        // Increment Prometheus counter
        metrics.WinklRequestInc(appContext, method, path, resp.Status)

        return nil
    }
}
```

### Request Enrichment Headers

| Header | Source | Example |
|--------|--------|---------|
| `Authorization` | Generated Basic auth | `Basic dXNlcjE...` |
| `x-bb-partner-id` | JWT claims | `12345` |
| `x-bb-plan-type` | Partner contract | `PRO`, `FREE`, `BUSINESS` |
| `x-bb-device` | Original request | `device-uuid-123` |
| `x-bb-requestid` | Generated UUID | `550e8400-e29b-41d4...` |

---

## 7. HTTP HEADERS SPECIFICATION

### Whitelisted Headers (25+)

```go
// header/header.go
const (
    ForwardedFor   = "x-forwarded-for"      // Original client IP
    UID            = "x-bb-uid"              // User ID
    Timestamp      = "x-bb-timestamp"        // Request timestamp
    Os             = "x-bb-os"               // Operating system
    ClientID       = "x-bb-clientid"         // Client/app ID
    RequestID      = "x-bb-requestid"        // Unique request ID
    ClientType     = "x-bb-clienttype"       // Client type (CUSTOMER, ERP)
    Version        = "x-bb-version"          // App version number
    VersionName    = "x-bb-version-name"     // App version name
    Latitude       = "x-bb-latitude"         // Device latitude
    Longitude      = "x-bb-longitude"        // Device longitude
    Location       = "x-bb-location"         // Location string
    Device         = "x-bb-deviceid"         // Device ID
    BBDevice       = "x-bb-custom-deviceid"  // Custom device ID
    Channel        = "x-bb-channelid"        // Channel/source ID
    Country        = "x-bb-country"          // Country code
    Locale         = "accept-language"       // Language (en, hi, bn, ta, te)
    ApolloOpName   = "x-apollo-operation-name"   // GraphQL operation
    ApolloOpID     = "x-apollo-operation-id"     // GraphQL operation ID
    UserRT         = "x-bb-user-rt"          // User RT
    AdvertisingID  = "x-bb-advertising-id"   // Advertising ID
    ABTestInfo     = "x-ab-test-info"        // A/B test assignments
    PartnerId      = "x-bb-partner-id"       // Partner ID
    PlanType       = "x-bb-plan-type"        // SaaS plan type
    NewUser        = "x-bb-new-user"         // New user flag
    PpId           = "x-bb-pp-id"            // PP ID
)
```

---

## 8. DOMAIN MODELS

### User Authentication Models

```go
// client/entry/identity.go

// Client information
type Client struct {
    Id      string  // Client ID (app identifier)
    AppName string  // Application name
    AppType string  // "CUSTOMER", "ERP", etc.
    Enabled bool    // Client status
}

// User profile
type User struct {
    Id           int
    Uidx         string
    Dob          *int64
    Emails       []UserEmail
    Phones       []UserPhone
    Gender       *Gender  // "MALE", "FEMALE", "OTHERS"
    ReferralCode *string
}

// Main authorization entity
type UserClientAccount struct {
    Id                 int
    ClientId           string
    ClientAppType      string
    UserId             int
    Status             UserClientAccountStatus  // ONBOARDING, ACTIVE, INACTIVE, BLOCKED

    Name               *string
    Bio                *string
    Email              *string
    ProfileImageId     *int
    ProfileImage       *AssetInfo
    CoverImageId       *int
    CoverImage         *AssetInfo

    User               *User
    HostProfile        *HostProfile        // Influencer profile
    CustomerProfile    *CustomerProfile    // Buyer profile
    PartnerProfile     *PartnerUserProfile // SaaS partner profile

    SocialAccounts     []SocialAccount
    KycData            *UserKycData
}

// Partner-specific user profile
type PartnerUserProfile struct {
    PartnerId           int64
    Phone               *string
    PlanType            string  // "FREE", "PRO", "BUSINESS"
    VisitedDemo         bool
    AssignedCampaignIds []int
    WinklBrand          *WinklBrand
}

// Social media account
type SocialAccount struct {
    Platform           SocialNetworkType  // GOOGLE, FB, TWITTER, TIKTOK, YOUTUBE, INSTAGRAM
    Handle             string
    SocialUserId       *string
    AccessToken        *string
    GraphAccessToken   *string
    AccessTokenExpired *bool
    Verified           bool
    Metrics            *SocialAccountMetrics
}
```

### Partner Models

```go
// client/response/partner.go

type PartnerResponse struct {
    Status   entry.Status
    Partners []Partner
}

type Partner struct {
    ID                   int64
    Name                 string
    Status               string
    InventorySyncEnabled bool
    Contracts            []Contract
}

type Contract struct {
    ID               int64   // Contract ID
    PartnerID        int64   // Partner ID
    ContractType     string  // "SAAS"
    Status           string
    OnboardingStatus string
    Plan             string  // Pricing plan (FREE, PRO, BUSINESS)
    StartTime        int64   // Unix milliseconds
    EndTime          int64   // Unix milliseconds
}
```

---

## 9. EXTERNAL SERVICE CLIENTS

### Base HTTP Client

```go
// client/client.go
type BaseClient struct {
    *resty.Client
    Context *gatewayContext.Context
}

func New(ctx *gatewayContext.Context) *BaseClient {
    return &BaseClient{
        resty.New().
            SetDebug(true).
            SetTimeout(15 * time.Second).
            SetLogger(NewLogger(ctx.Logger)).
            OnRequestLog(func(r *resty.RequestLog) error {
                if ctx.Config.LOG_LEVEL >= 3 {
                    r.Body = ""  // Don't log body at high log levels
                }
                return nil
            }).
            OnResponseLog(func(r *resty.ResponseLog) error {
                if ctx.Config.LOG_LEVEL >= 3 {
                    r.Body = ""
                }
                return nil
            }),
        ctx,
    }
}

func (c *BaseClient) GenerateRequest(authorization string) *resty.Request {
    c.SetAllowGetMethodPayload(true)
    req := c.NewRequest()
    req.SetHeader("Authorization", authorization)

    // Propagate X-Headers from context
    if c.Context.XHeader != nil {
        for key, vals := range c.Context.XHeader {
            for _, val := range vals {
                req.Header.Add(key, val)
            }
        }
    }

    return req
}
```

### Identity Service Client

```go
// client/identity/identity.go
type Client struct {
    *client.BaseClient
}

func New(ctx *context.Context) *Client {
    c := &Client{client.New(ctx)}
    c.SetTimeout(time.Duration(1 * time.Minute))  // Longer timeout for auth
    return c
}

func (c *Client) VerifyToken(input *input.Token, skipTokenCache bool) (*response.UserAuthResponse, error) {
    response := &response.UserAuthResponse{}
    _, err := c.GenerateRequest(c.Context.ClientAuthorization).
        SetBody(input).
        SetResult(response).
        SetQueryParam("skipTokenCache", strconv.FormatBool(skipTokenCache)).
        Post(c.Context.Config.IDENTITY_URL + "/auth/verify/token?loadPreferences=true")

    if err != nil {
        return nil, err
    }
    return response, nil
}

func (c *Client) UserAuth(userAuthInput *input.UserAuth) (*response.UserAuthResponse, error) {
    response := &response.UserAuthResponse{}
    _, err := c.GenerateRequest(c.Context.ClientAuthorization).
        SetBody(userAuthInput).
        SetResult(response).
        Post(c.Context.Config.IDENTITY_URL + "/auth/login")

    if err != nil {
        return nil, err
    }
    return response, nil
}
```

### Partner Service Client

```go
// client/partner/partner.go
type Client struct {
    *client.BaseClient
}

func New(ctx *context.Context) *Client {
    c := &Client{client.New(ctx)}
    c.SetTimeout(time.Duration(1 * time.Minute))
    return c
}

func (c *Client) FindPartnerById(partnerId string) (*response.PartnerResponse, error) {
    response := &response.PartnerResponse{}
    _, err := c.GenerateRequest(c.Context.ClientAuthorization).
        SetPathParams(map[string]string{"partnerId": partnerId}).
        SetResult(response).
        Get(c.Context.Config.PARTNER_URL + "/partner/{partnerId}")

    if err != nil {
        return nil, err
    }
    return response, nil
}
```

---

## 10. OBSERVABILITY

### Prometheus Metrics

```go
// metrics/metrics.go
var (
    WinklRequestCounter          *prometheus.CounterVec
    LoadPageEmptyResponseCounter *prometheus.CounterVec
    LoadPageInvalidSlugCounter   *prometheus.CounterVec
    LoadPageErrorCounter         *prometheus.CounterVec
)

func SetupCollectors(config config.Config) {
    collectorInit.Do(func() {
        // Track all proxy requests
        WinklRequestCounter = promauto.NewCounterVec(prometheus.CounterOpts{
            Name: "requests_winkl",
            Help: "Track Requests via Winkl Proxy",
        }, []string{"method", "path", "channel", "clientId", "responseStatus"})

        // Track loadPage calls with 0 results
        LoadPageEmptyResponseCounter = promauto.NewCounterVec(prometheus.CounterOpts{
            Name: "zero_edges_loadPage",
            Help: "Track # of loadPage calls with 0 edges",
        }, []string{"slug", "version", "channel", "pageNo"})

        // Track invalid slug requests
        LoadPageInvalidSlugCounter = promauto.NewCounterVec(prometheus.CounterOpts{
            Name: "invalid_slug_loadPage",
            Help: "Track # of loadPage calls with invalid slug",
        }, []string{"slug", "version", "channel", "pageNo"})

        // Track loadPage errors
        LoadPageErrorCounter = promauto.NewCounterVec(prometheus.CounterOpts{
            Name: "errors_loadPage",
            Help: "Track # of loadPage errors",
        }, []string{"slug", "version", "channel", "pageNo"})
    })
}

func WinklRequestInc(gc *context.Context, method string, path string, responseStatus string) {
    if WinklRequestCounter != nil {
        WinklRequestCounter.With(prometheus.Labels{
            "method":         method,
            "path":           path,
            "responseStatus": responseStatus,
            "channel":        gc.XHeader.Get(header.Channel),
            "clientId":       gc.XHeader.Get(header.ClientID),
        }).Inc()
    }
}
```

**Metrics Endpoint**: `GET /metrics`

### Sentry Error Tracking

```go
// sentry/sentry.go
func Setup(config config.Config) {
    if config.Env != "local" {
        if err := sentry.Init(sentry.ClientOptions{
            Dsn:              config.SENTRY_DSN,
            Environment:      config.Env,
            AttachStacktrace: true,
        }); err != nil {
            log.Println("Sentry init failed")
        }
    }
}

func PatchErrorToSentry(gc *context.Context, e error, err interface{}) {
    gcInterface := make(map[string]interface{})

    if gc != nil {
        gcInterface["headers"] = gc.XHeader
        if gc.UserClientAccount != nil {
            gcInterface["user"] = gc.UserClientAccount
        }
        if gc.Client != nil {
            gcInterface["client"] = gc.Client
        }

        sentry.ConfigureScope(func(scope *sentry.Scope) {
            tags := map[string]string{
                header.RequestID:    gc.XHeader.Get(header.RequestID),
                header.ApolloOpName: gc.XHeader.Get(header.ApolloOpName),
                header.Version:      gc.XHeader.Get(header.Version),
                header.ClientID:     gc.XHeader.Get(header.ClientID),
                header.Channel:      gc.XHeader.Get(header.Channel),
            }
            if gc.UserClientAccount != nil {
                tags["userId"] = strconv.Itoa(gc.UserClientAccount.UserId)
            }
            scope.SetTags(tags)
        })
    }

    if e != nil {
        sentry.CaptureException(e)
    }
    if err != nil {
        sentry.CaptureException(errors.New(fmt.Sprintf("%+v", err)))
    }
}
```

### Health Check Endpoints

```go
// api/heartbeat/heartbeat.go
var beat = false  // Health state flag

// GET /heartbeat/
func Beat(config config.Config) func(c *gin.Context) {
    return func(c *gin.Context) {
        if beat {
            c.Status(http.StatusOK)   // 200 - Healthy
        } else {
            c.Status(http.StatusGone) // 410 - Out of LB
        }
    }
}

// PUT /heartbeat/?beat=true|false
func ModifyBeat(config config.Config) func(c *gin.Context) {
    return func(c *gin.Context) {
        modifyBeat, err := strconv.ParseBool(c.Query("beat"))
        if err != nil {
            c.Status(http.StatusBadRequest)
        } else {
            beat = modifyBeat
            c.Status(http.StatusOK)
        }
    }
}
```

---

## 11. CI/CD & DEPLOYMENT

### GitLab CI Pipeline

```yaml
# .gitlab-ci.yml
variables:
  PORT: 8009

stages:
  - test
  - build
  - deploy_staging
  - deploy_production

test:
  stage: test
  script: echo "Running tests"

build:
  stage: build
  variables:
    GOOS: linux
    GOARCH: amd64
  script:
    - echo "Building SaaS gateway"
    - env GOOS=$GOOS GOARCH=$GOARCH go build
  tags:
    - saas-builder
  artifacts:
    paths:
      - .env, .env.local, .env.stage, .env.production
      - saas-gateway
      - scripts/start.sh
    expire_in: 1 week
  only:
    - master
    - staging

deploy_staging:
  stage: deploy_staging
  variables:
    ENV: STAGE
    GOGC: 250  # Aggressive garbage collection
  script:
    - /bin/bash scripts/start.sh
  tags:
    - deployer-stage-search
  when: manual
  only:
    - master
    - staging

deploy_prod_1:
  stage: deploy_production
  variables:
    ENV: PRODUCTION
    GOGC: 250
  script:
    - /bin/bash scripts/start.sh
  tags:
    - cb1-1  # Node 1
  when: manual
  only:
    - master

deploy_prod_2:
  stage: deploy_production
  variables:
    ENV: PRODUCTION
    GOGC: 250
  script:
    - /bin/bash scripts/start.sh
  tags:
    - cb2-1  # Node 2
  when: manual
  only:
    - master
```

### Deployment Script

```bash
# scripts/start.sh
#!/bin/bash
ulimit -n 100000  # Increase file descriptors

# Check if process already running
PID=$(ps aux | grep saas-gateway | grep -v grep | awk '{print $2}')

if [ -z "$PID" ]; then
    echo "SaaS gateway is not running"
else
    # Graceful shutdown
    echo "Bringing OOLB (Out of Load Balancer)"
    curl -vXPUT http://localhost:$PORT/heartbeat/?beat=false

    # Wait for in-flight requests
    sleep 15

    echo "Killing SaaS Gateway"
    kill -9 $PID
fi

sleep 10

# Start new process
echo "Starting SaaS Gateway"
ENV=$CI_ENVIRONMENT_NAME ./saas-gateway >> "logs/out.log" 2>&1 &

# Wait for startup
sleep 20

# Bring back into load balancer
echo "Bringing in LB (Load Balancer)"
curl -vXPUT http://localhost:$PORT/heartbeat/?beat=true
```

---

## 12. CONFIGURATION

### Environment Variables

```go
// config/config.go
type Config struct {
    Port                   string      // Server port (default: 8009)
    HMAC_SECRET            string      // Base64-encoded JWT signing secret
    Env                    string      // Environment: local, STAGE, PRODUCTION
    SAAS_URL               string      // Main backend URL
    SAAS_DATA_URL          string      // Data service URL
    IDENTITY_URL           string      // Identity service endpoint
    PARTNER_URL            string      // Partner service endpoint
    SENTRY_DSN             string      // Sentry error tracking
    ERP_CLIENT_ID          string      // ERP client identifier
    LOG_LEVEL              int         // 0=verbose, 3=silent
    RedisClusterAddresses  []string    // Redis cluster nodes
    REDIS_CLUSTER_PASSWORD string      // Redis auth password
}
```

### Environment Configuration

| Config | Local | Staging | Production |
|--------|-------|---------|------------|
| SAAS_URL | localhost:7179 | stage-search:7179 | coffee.goodcreator.co |
| IDENTITY_URL | sb2:7000/identity-service/api | sb2:7000/identity-service/api | identityservice.bulbul.tv/identity-service/api |
| PARTNER_URL | sb2:7100/product-service/api | sb2:7100/product-service/api | productservice.bulbul.tv/product-service/api |
| Redis Nodes | 6 (localhost:30001-6) | 3 (stage:6379-6381) | 3 (prod:6379-6381) |
| LOG_LEVEL | 2 | 2 | 2 |

---

## 13. KEY METRICS SUMMARY

| Metric | Value |
|--------|-------|
| **Framework** | Gin v1.8.1 |
| **Go Version** | 1.12 |
| **Binary Size** | 23MB |
| **Services Proxied** | 13 |
| **Middleware Layers** | 7 |
| **CORS Origins** | 41 |
| **Whitelisted Headers** | 25+ |
| **Authentication** | JWT + Redis |
| **Cache Layer 1** | Ristretto (10M keys, 1GB) |
| **Cache Layer 2** | Redis Cluster (3-6 nodes) |
| **Prometheus Metrics** | 4 counters |
| **Production Nodes** | 2 (multi-node) |
| **Port** | 8009 |
| **Debug Port** | 6069 (pprof) |

---

## 14. SKILLS DEMONSTRATED

### API Gateway Engineering
- **Reverse Proxy**: Using Go's `httputil.NewSingleHostReverseProxy`
- **Middleware Pipeline**: 7-layer request processing
- **Request Transformation**: Header enrichment, path rewriting
- **Response Handling**: CORS, metrics, logging

### Authentication & Security
- **JWT Validation**: HMAC-SHA256 signature verification
- **Session Management**: Redis cluster with fallback
- **Plan-Based Authorization**: Contract date validation
- **Multi-tenant Support**: Partner ID and plan type propagation

### Caching Architecture
- **Two-Layer Caching**: Ristretto (local) + Redis (distributed)
- **LFU Eviction**: Ristretto's frequency-based eviction
- **Connection Pooling**: 100 connections to Redis cluster

### Observability
- **Prometheus Metrics**: Request counters with labels
- **Sentry Integration**: Stack traces and breadcrumbs
- **Structured Logging**: Zerolog + Logrus
- **Health Checks**: Graceful LB integration

### DevOps
- **CI/CD**: GitLab pipeline with multi-node deployment
- **Graceful Deployment**: Health state toggling
- **Configuration Management**: Environment-based configs

---

## 15. INTERVIEW TALKING POINTS

### 1. "Tell me about an API gateway you built"
- **Architecture**: 7-layer middleware pipeline with reverse proxy
- **Authentication**: JWT + Redis session caching with Identity service fallback
- **Caching**: Two-layer (Ristretto 1GB + Redis cluster)
- **Scale**: 13 microservices, dual-node production

### 2. "Describe your approach to authentication"
- **JWT**: HMAC-SHA256 with claims extraction
- **Session**: Redis cluster with `session:{id}` keys
- **Fallback**: Identity service API when cache miss
- **Plan Validation**: Contract date range checking

### 3. "How do you handle multi-tenancy?"
- **Partner ID**: Extracted from JWT, enriched to headers
- **Plan Type**: Validated against contract dates
- **Header Propagation**: `x-bb-partner-id`, `x-bb-plan-type`

### 4. "Explain your caching strategy"
- **Layer 1**: Ristretto (10M keys, 1GB, LFU, per-instance)
- **Layer 2**: Redis cluster (3-6 nodes, shared)
- **Pattern**: Check L1 → Check L2 → Fetch source → Populate both

### 5. "How do you ensure zero-downtime deployments?"
- **Health Endpoint**: `/heartbeat/?beat=false` to remove from LB
- **Drain Period**: 15 seconds for in-flight requests
- **Restart**: Kill and restart with new binary
- **Re-enable**: `/heartbeat/?beat=true` after stabilization

---

*Analysis covers the complete SaaS Gateway implementation with 7-layer middleware, two-tier caching, JWT authentication, and multi-node deployment for a creator platform serving 13 microservices.*
