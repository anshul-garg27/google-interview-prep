package appcontext

import (
	"coffee/constants"
	"coffee/core/event"
	"coffee/core/persistence"
	"coffee/helpers"
	"context"
	"encoding/base64"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-resty/resty/v2"

	log "github.com/sirupsen/logrus"
)

const (
	Authorization       string = "Authorization"
	ForwardedFor        string = "x-forwarded-for"
	UID                 string = "x-bb-uid"
	Timestamp           string = "x-bb-timestamp"
	Os                  string = "x-bb-os"
	ClientID            string = "x-bb-clientid"
	RequestID           string = "x-bb-requestid"
	AccountID           string = "x-bb-account-id"
	PartnerID           string = "x-bb-partner-id"
	ClientType          string = "x-bb-clienttype"
	Version             string = "x-bb-version"
	VersionName         string = "x-bb-version-name"
	Device              string = "x-bb-deviceid"
	BBDevice            string = "x-bb-custom-deviceid"
	Channel             string = "x-bb-channelid"
	Locale              string = "accept-language"
	UserRT              string = "x-bb-user-rt"
	AdvertisingID       string = "x-bb-advertising-id"
	XRetryCount         string = "x-retry-count"
	PlanType            string = "x-bb-plan-type"
	ClientAuthorization string = "x-bb-client-authorization"
)

type RequestContext struct {
	Ctx             context.Context
	Mutex           sync.Mutex
	Properties      map[string]interface{}
	PartnerId       *int64
	AccountId       *int64
	UserId          *int64
	UserName        *string
	Authorization   string
	DeviceId        *string
	IsLoggedIn      bool
	IsPremiumMember bool
	IsBasicMember   bool
	Session         persistence.Session
	CHSession       persistence.Session
	PlanType        *constants.PlanType
}

func CreateApplicationContext(ctx context.Context, r *http.Request) (RequestContext, error) {
	appCtx := RequestContext{Ctx: ctx}
	properties := make(map[string]interface{})
	appCtx.Properties = properties
	if r != nil {
		for name, values := range r.Header {
			for _, value := range values {
				switch strings.ToLower(name) {
				case strings.ToLower(PartnerID):
					appCtx.PartnerId = helpers.ParseToInt64(value)
				case strings.ToLower(PlanType):
					if value == string(constants.FreePlan) {
						PlanType := constants.FreePlan
						appCtx.PlanType = &PlanType

					} else if value == string(constants.PaidPlan) {
						PlanType := constants.PaidPlan
						appCtx.PlanType = &PlanType
					} else if value == string(constants.SaasPlan) {
						PlanType := constants.SaasPlan
						appCtx.PlanType = &PlanType
					} else {
						log.Error("invalid plan type")
						PlanType := constants.PlanType(value)
						appCtx.PlanType = &PlanType
					}

				case strings.ToLower(Authorization):
					appCtx.Authorization = value
					var err error
					appCtx, err = PopulateUserInfoFromAuthToken(appCtx, value)
					if err != nil {
						return appCtx, err
					}
				}
			}
		}
	}
	return appCtx, nil
}

func CreateApplicationContextPubSub(m *message.Message) (*message.Message, error) {
	ctx := m.Context()
	appCtx := RequestContext{Ctx: ctx}
	properties := make(map[string]interface{})
	appCtx.Properties = properties
	if m != nil {
		for name, value := range m.Metadata {
			log.Println("Adding to context - ", name, value)
			switch strings.ToLower(name) {
			case strings.ToLower(PartnerID):
				appCtx.PartnerId = helpers.ParseToInt64(value)

			case strings.ToLower(Authorization):
				appCtx.Authorization = value
				var err error
				appCtx, err = PopulateUserInfoFromAuthToken(appCtx, value)
				if err != nil {
					return m, err
				}
			case strings.ToLower(PlanType):
				if value == string(constants.FreePlan) {
					PlanType := constants.FreePlan
					appCtx.PlanType = &PlanType

				} else if value == string(constants.PaidPlan) {
					PlanType := constants.PaidPlan
					appCtx.PlanType = &PlanType
				} else {
					log.Error("invalid plan type")
					PlanType := constants.PlanType(value)
					appCtx.PlanType = &PlanType
				}
			}
		}
	}
	// log.Println("Initialized PubSub Context - %+v", appCtx)
	newCtx := context.WithValue(ctx, constants.AppContextKey, &appCtx)
	m.SetContext(newCtx)
	return m, nil
}

func (c *RequestContext) SetValue(key string, value interface{}) {
	c.Properties[key] = value
}

func (c *RequestContext) GetValue(key string) interface{} {
	return c.Properties[key]
}

func (c *RequestContext) SetSession(ctx context.Context, session persistence.Session) {
	c.Session = session
	c.Session.PerformAfterCommit(ctx, c.applyEvents)
}

func (c *RequestContext) GetSession() persistence.Session {
	return c.Session
}
func (c *RequestContext) applyEvents(ctx context.Context) {
	log.Debug("applying events")
	appContext := ctx.Value(constants.AppContextKey).(*RequestContext)
	eventStore := appContext.GetValue("event_store")
	if eventStore == nil {
		eventStore = []event.Event{}
	}
	_store := eventStore.([]event.Event)
	for i := range _store {
		_store[i].Apply()
	}
}

func (c *RequestContext) PopulateHeadersForHttpReq(req *resty.Request) {
	if c.PartnerId != nil {
		req.Header.Add(PartnerID, strconv.FormatInt(*c.PartnerId, 10))
	}
	if c.PlanType != nil {
		req.Header.Add(PlanType, string(*c.PlanType))
	}
	if c.AccountId != nil {
		req.Header.Add(AccountID, strconv.FormatInt(*c.AccountId, 10))
	}
	req.Header.Add(Authorization, c.Authorization)
}

func PopulateUserInfoFromAuthToken(appCtx RequestContext, token string) (RequestContext, error) {
	tokenParts := strings.Split(token, " ")
	decodedAuthToken, _ := base64.StdEncoding.DecodeString(tokenParts[1])
	if decodedAuthToken == nil {
		return appCtx, errors.New("UnAuthorized Access")
	}
	namePassword := strings.Split(string(decodedAuthToken), ":")
	if len(namePassword) < 2 {
		return appCtx, errors.New("UnAuthorized Access")
	}
	userIdNameOperationArr := strings.Split(namePassword[0], "~")
	appCtx.IsLoggedIn = false
	if len(userIdNameOperationArr) >= 3 {
		loggedInFlag, err1 := strconv.Atoi(userIdNameOperationArr[2])
		if err1 == nil && loggedInFlag == 1 {
			appCtx.IsLoggedIn = true
		}
	}
	if !appCtx.IsLoggedIn {
		appCtx.DeviceId = &userIdNameOperationArr[0]
	} else {
		userId, _ := strconv.ParseInt(userIdNameOperationArr[0], 10, 64)
		if userId != 0 {
			appCtx.UserId = &userId
		}
		if len(userIdNameOperationArr) > 1 {
			appCtx.UserName = &userIdNameOperationArr[1]
		}
		if len(userIdNameOperationArr) > 3 {
			userAccountId, _ := strconv.ParseInt(userIdNameOperationArr[3], 10, 64)
			appCtx.AccountId = &userAccountId
		}
	}
	return appCtx, nil
}

func (c *RequestContext) SetCHSession(ctx context.Context, session persistence.Session) {
	c.CHSession = session
}

func (c *RequestContext) GetCHSession() persistence.Session {
	return c.CHSession
}
