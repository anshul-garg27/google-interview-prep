package saas

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"init.bulbul.tv/bulbul-backend/saas-gateway/config"
	"init.bulbul.tv/bulbul-backend/saas-gateway/context"
	"init.bulbul.tv/bulbul-backend/saas-gateway/header"
	"init.bulbul.tv/bulbul-backend/saas-gateway/metrics"
	"init.bulbul.tv/bulbul-backend/saas-gateway/util"
)

func responseHandler(appContext *context.Context, method string, path string, origin string, setOrigin bool, originalAuth string) func(*http.Response) error {
	if setOrigin {
		return func(resp *http.Response) (err error) {
			resp.Header.Del("Access-Control-Allow-Origin")
			resp.Header.Set("Access-Control-Allow-Origin", origin)
			if resp.Status != "200 OK" {
				log.Printf("SAAS-Bad-Response - %s, %s, %v : from : %s", method, path, resp.Status, originalAuth)
			}
			metrics.WinklRequestInc(appContext, method, path, resp.Status)
			return nil
		}
	} else {
		return func(resp *http.Response) (err error) {
			resp.Header.Del("Access-Control-Allow-Origin")
			if resp.Status != "200 OK" {
				log.Printf("SAAS-Bad-Response - %s, %s, %v : from : %s", method, path, resp.Status, originalAuth)
			}
			metrics.WinklRequestInc(appContext, method, path, resp.Status)
			return nil
		}
	}
}

func ReverseProxy(module string, config config.Config) func(c *gin.Context) {
	return func(c *gin.Context) {
		appContext := util.GatewayContextFromGinContext(c, config)
		saasUrl := config.SAAS_URL
		if module == "social-profile-service" {
			saasUrl = config.SAAS_DATA_URL
		}
		remote, err := url.Parse(saasUrl)
		if err != nil {
			panic(err)
		}
		authHeader := appContext.UserAuthorization
		var partnerId string
		planType := "FREE"
		if appContext.UserClientAccount != nil {
			if appContext.UserClientAccount.PartnerProfile != nil {
				partnerId = strconv.Itoa(int(appContext.UserClientAccount.PartnerProfile.PartnerId))
			}
			if appContext.UserClientAccount.PartnerProfile != nil {
				if appContext.UserClientAccount.PartnerProfile.PlanType != "" {
					planType = appContext.UserClientAccount.PartnerProfile.PlanType
				}

			}
		}
		proxy := httputil.NewSingleHostReverseProxy(remote)
		deviceId := appContext.XHeader.Get(header.Device)
		proxy.Director = func(req *http.Request) {
			req.Header = c.Request.Header
			req.Header.Set("Authorization", authHeader)
			//req.Header.Set(header.AccountId, accountId)
			req.Header.Set(header.PartnerId, partnerId)
			req.Header.Set(header.PlanType, planType)
			req.Header.Set(header.Device, deviceId)
			req.Header.Del("accept-encoding")
			req.Host = remote.Host
			req.URL.Scheme = remote.Scheme
			req.URL.Host = remote.Host
			req.URL.Path = "/" + module + c.Param("any")
			log.Print(req)
		}

		if c.Request.Method == "OPTIONS" {
			proxy.ModifyResponse = responseHandler(appContext, c.Request.Method, c.Param("any"), c.Request.Header.Get("origin"), true, appContext.Context.GetHeader("Authorization"))
		} else {
			proxy.ModifyResponse = responseHandler(appContext, c.Request.Method, c.Param("any"), c.Request.Header.Get("origin"), false, appContext.Context.GetHeader("Authorization"))
		}
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
