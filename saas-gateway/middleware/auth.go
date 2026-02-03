package middleware

import (
	b64 "encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"init.bulbul.tv/bulbul-backend/saas-gateway/cache"
	"init.bulbul.tv/bulbul-backend/saas-gateway/client/entry"
	"init.bulbul.tv/bulbul-backend/saas-gateway/client/identity"
	"init.bulbul.tv/bulbul-backend/saas-gateway/client/input"
	"init.bulbul.tv/bulbul-backend/saas-gateway/client/partner"
	"init.bulbul.tv/bulbul-backend/saas-gateway/config"
	"init.bulbul.tv/bulbul-backend/saas-gateway/header"
	"init.bulbul.tv/bulbul-backend/saas-gateway/util"
)

func AppAuth(config config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// get gateway context
		gc := util.GatewayContextFromGinContext(c, config)
		verified := false

		// condensed auth flow
		// checks client id
		// authorization header
		// validates jwt
		// check if session id in jwt exists in redis else fallsback to identity

		// check for authorization header
		if clientId := gc.GetHeader(header.ClientID); clientId != "" {
			clientType := gc.GetHeader(header.ClientType)
			if clientType == "" {
				clientType = "CUSTOMER"
			}
			gc.GenerateClientAuthorizationWithClientId(clientId, clientType)

			apolloOp := gc.Context.GetHeader(header.ApolloOpName)
			if apolloOp == "initDeviceGoMutation" || apolloOp == "initDeviceV2" || apolloOp == "initDeviceV2Mutation" || apolloOp == "initDevice" || strings.Contains(c.Request.URL.Path, "/api/auth/init") || strings.Contains(c.Request.URL.Path, "/api/auth/init/v2") {
				verified = true
			} else if authorization := gc.Context.GetHeader("Authorization"); authorization != "" {
				// skip this flow for init device
				splittedAuthHeader := strings.Split(authorization, " ")
				// verify user authorization header only if its not init device call (this is to prevent 2 calls of verify token on init device, one on auth header another on user token)
				if len(splittedAuthHeader) > 1 && splittedAuthHeader[1] != "" {
					if token, err := jwt.Parse(splittedAuthHeader[1], func(token *jwt.Token) (interface{}, error) {
						if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
							return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
						}
						hmac, _ := b64.StdEncoding.DecodeString(config.HMAC_SECRET)
						return hmac, nil
					}); err == nil {
						if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
							if sessionVal, ok := claims["sid"]; ok {
								if sessionId, ok := sessionVal.(string); ok && sessionId != "" {
									keyExists := cache.Redis(gc.Config).Exists("session:" + sessionId)
									if exist, err := keyExists.Result(); err == nil && exist == 1 {
										verified = true
										gc.UserToken = &splittedAuthHeader[1]
										if val, ok := claims["uid"]; ok {
											if userId := util.InterfaceToInt(val); userId != 0 && userId != -1 {
												// update user and client from verify token response
												userClientAccount := &entry.UserClientAccount{UserId: userId, ClientId: clientId}
												if userAccountId, ok := claims["userAccountId"]; ok {
													userClientAccount.Id = int(userAccountId.(float64))
												}
												if userPartnerId, ok := claims["partnerId"]; ok {
													if userClientAccount.PartnerProfile == nil {
														userClientAccount.PartnerProfile = &entry.PartnerUserProfile{}
													}
													userClientAccount.PartnerProfile.PartnerId = int64(userPartnerId.(float64))
													partnerId := strconv.FormatInt(userClientAccount.PartnerProfile.PartnerId, 10)
													partnerResponse, err := partner.New(gc).FindPartnerById(partnerId)
													if err == nil {
														if len(partnerResponse.Partners) > 0 {
															partnerDetail := partnerResponse.Partners[0]
															if len(partnerDetail.Contracts) > 0 {
																for _, contract := range partnerDetail.Contracts {
																	if contract.ContractType == "SAAS" {
																		userClientAccount.PartnerProfile.PlanType = contract.Plan
																	}
																}
															}
														}
													}
												}
												// TBD IMPLEMENT GET PARTNER DETAILS API AND PASS THIS PLAN DETAILS
												gc.SetUserClientAccount(userClientAccount)
											}
										}
									} else {
										verifyTokenResponse, err := identity.New(gc).VerifyToken(&input.Token{Token: splittedAuthHeader[1]}, false)
										if err == nil && verifyTokenResponse != nil && verifyTokenResponse.Status.Type == "SUCCESS" {
											verified = true
											// update user and client from verify token response
											gc.UserToken = &verifyTokenResponse.Token
											userClientAccount := verifyTokenResponse.UserClientAccount
											if userClientAccount != nil && userClientAccount.UserId != -1 && userClientAccount.UserId != 0 {
												partnerId := strconv.FormatInt(userClientAccount.PartnerProfile.PartnerId, 10)
												partnerResponse, err := partner.New(gc).FindPartnerById(partnerId)
												if err == nil {
													if len(partnerResponse.Partners) > 0 {
														partnerDetail := partnerResponse.Partners[0]
														if len(partnerDetail.Contracts) > 0 {
															for _, contract := range partnerDetail.Contracts {
																if contract.ContractType == "SAAS" {
																	currentTime := time.Now()
																	startTime := time.Unix(0, contract.StartTime*int64(time.Millisecond))
																	endTime := time.Unix(0, contract.EndTime*int64(time.Millisecond))
																	if currentTime.After(startTime) && currentTime.Before(endTime) {
																		userClientAccount.PartnerProfile.PlanType = contract.Plan
																	} else {
																		userClientAccount.PartnerProfile.PlanType = "FREE"
																	}

																}
															}
														}
													}
												}
												gc.SetUserClientAccount(userClientAccount)
											}
											gc.Client = verifyTokenResponse.Client
										}
									}
								}
							}
						}
					}
				}
			}
		}

		// irrespective of above calls, will generate device authorization at the least
		// gc.GenerateAuthorization()
		// c.Next()
		//This Verified Flag checks that token is correct.  verified  && gc.UserClientAccount== nil

		if !verified || (gc.UserClientAccount == nil && strings.Contains(c.Request.URL.Path, "/test/test/sample")) {
			gc.AbortWithStatus(http.StatusUnauthorized)
		} else {
			gc.GenerateAuthorization()
			c.Next()
		}
	}
}
