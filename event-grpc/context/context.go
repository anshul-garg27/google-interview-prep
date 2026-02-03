package context

import (
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"init.bulbul.tv/bulbul-backend/event-grpc/client/entry"
	"init.bulbul.tv/bulbul-backend/event-grpc/config"
	"init.bulbul.tv/bulbul-backend/event-grpc/generator"
	"init.bulbul.tv/bulbul-backend/event-grpc/header"
	"init.bulbul.tv/bulbul-backend/event-grpc/locale/locale"
	"net/http"
	"os"
	"strconv"
)

var XHEADERS = []string{
	header.RequestID,
	header.Authorization,
	header.ApolloOpName,
	header.Channel,
	header.ClientID,
	header.ClientType,
	header.Device,
	header.BBDevice,
	header.Os,
	header.UserRT,
	header.VersionName,
	header.Locale,
}

type Context struct {
	*gin.Context

	Logger zerolog.Logger

	// all x-bb- headers in one place
	XHeader http.Header

	// device authorization or user authorized header if user exists
	DeviceAuthorization string

	// user authorization header
	UserAuthorization string

	// user authorization header
	ERPUserAuthorization string

	// client authorization header
	ClientAuthorization string

	// client
	Client *entry.Client

	// user client account
	UserClientAccount *entry.UserClientAccount

	// clientToken
	ClientToken *string

	// userToken
	UserToken *string

	Config config.Config

	Locale locale.Locale

	// fixme add feature gates
}

func DisallowBots(c *gin.Context) {
	if c != nil && c.Writer != nil && c.Writer.Header() != nil {
		c.Writer.Header().Set(header.RobotsTag, "noindex")
	}
}

func New(c *gin.Context, config config.Config) *Context {
	reqId := generator.SetNXRequestIdOnContext(c)
	DisallowBots(c)
	ctx := &Context{
		Context: c,
		Logger:  zerolog.New(os.Stdout).Level(zerolog.Level(config.LOG_LEVEL)).With().Str(header.RequestID, reqId).Logger(),
		Config:  config,
		XHeader: make(http.Header),
		Locale:  locale.LocaleEnglish,
	}
	ctx.XHeader.Set(header.Locale, string(ctx.Locale))
	if c != nil {
		ctx.PickXHeadersFromReq()
	}
	return ctx
}

func (c *Context) PickXHeadersFromReq() *Context {

	// if request is not there we should pick request id from keys
	if c.Keys != nil && c.Keys[header.RequestID] != nil {
		c.XHeader.Set(header.RequestID, c.Keys[header.RequestID].(string))
	}

	if c.Context.Request == nil {
		return c
	}

	// whitelist headers to pick
	keys := []string{header.RequestID, header.ForwardedFor, header.UID, header.Timestamp, header.Os, header.ClientID, header.ClientType, header.Version, header.VersionName, header.Latitude, header.Longitude, header.Location, header.Device, header.BBDevice, header.Channel, header.Country, header.Locale, header.ApolloOpName, header.ApolloOpID, header.UserRT, header.AdvertisingID, header.PpId}
	validLocales := map[locale.Locale]bool{locale.LocaleEnglish: true, locale.LocaleHindi: true, locale.LocaleBengali: true, locale.LocaleTamil: true, locale.LocaleTelgu: true}

	for _, key := range keys {
		// canonical done for a reason
		// capture only non empty headers
		if c.Context.Request.Header.Get(http.CanonicalHeaderKey(key)) != "" {

			// read ABTestInfo header as RemoteConfig, rest all remains same
			if key == header.ABTestInfo {
				c.XHeader[http.CanonicalHeaderKey(header.RemoteConfig)], _ = c.Context.Request.Header[http.CanonicalHeaderKey(key)]
			} else {
				c.XHeader[http.CanonicalHeaderKey(key)], _ = c.Context.Request.Header[http.CanonicalHeaderKey(key)]
			}

			// overwrite accept-lang header if not among valid locales
			if key == header.Locale {
				if _, validLocale := validLocales[locale.Locale(c.Context.Request.Header.Get(http.CanonicalHeaderKey(key)))]; !validLocale {
					c.XHeader.Set(key, string(locale.LocaleEnglish))
				}
			}

		}
	}

	// set locale in context and in Xheaders if not set
	if _, ok := c.XHeader[http.CanonicalHeaderKey(header.Locale)]; ok {
		c.Locale = locale.Locale(c.XHeader.Get(header.Locale))
	} else {
		c.Locale = locale.LocaleEnglish
		c.XHeader.Set(header.Locale, string(c.Locale))
	}

	return c
}

func (c *Context) GenerateAuthorization() {
	c.GenerateClientAuthorization()
	c.GenerateUserAuthorization()
}

func (c *Context) GenerateUserAuthorization() {
	if c.UserClientAccount != nil && c.Client != nil && c.Client.AppType == "ERP" {
		c.ERPUserAuthorization = "Basic " + base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d%s%s%s%d%s%d", c.UserClientAccount.UserId, "~", getName(*c.UserClientAccount), "~", 1, ":", c.UserClientAccount.UserId)))
		c.UserAuthorization = c.ERPUserAuthorization
		c.DeviceAuthorization = c.UserAuthorization
	} else if c.UserClientAccount != nil {
		c.UserAuthorization = "Basic " + base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d%s%s%s%d%s%d", c.UserClientAccount.UserId, "~", getName(*c.UserClientAccount), "~", 1, ":", c.UserClientAccount.UserId)))
		c.DeviceAuthorization = c.UserAuthorization
	} else if c.XHeader.Get(header.Device) != "" {
		c.DeviceAuthorization = "Basic " + base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s%s%s%s%d%s%s", c.XHeader.Get(header.Device), "~", c.XHeader.Get(header.Device), "~", 0, ":", c.XHeader.Get(header.Device))))
	}
}

func getName(userClientAccount entry.UserClientAccount) string {
	var name string
	if userClientAccount.Name == nil {
		name = strconv.Itoa(userClientAccount.UserId)
	} else {
		name = *userClientAccount.Name
	}
	return name
}

func (c *Context) GenerateClientAuthorization() {
	if c.Client != nil {
		c.XHeader.Set(header.ClientID, c.Client.Id)
		c.XHeader.Set(header.ClientType, c.Client.AppType)

		c.ClientAuthorization = "Basic " + base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s%s%s%s%s", c.Client.Id, "~", c.Client.Id, ":", c.Client.Id)))
	}
}

func (c *Context) GenerateClientAuthorizationWithClientId(ClientId string, ClientAppType string) {
	if ClientId != "" {
		c.SetClient(&entry.Client{Id: ClientId, AppType: ClientAppType})
	}
}

func (c *Context) SetUserClientAccount(userClientAccount *entry.UserClientAccount) {
	c.UserClientAccount = userClientAccount
	c.GenerateUserAuthorization()
}

func (c *Context) SetChannel(channel string) {
	c.XHeader.Set(header.Channel, channel)
}

func (c *Context) SetClient(client *entry.Client) {
	c.Client = client
	c.GenerateClientAuthorization()
}
