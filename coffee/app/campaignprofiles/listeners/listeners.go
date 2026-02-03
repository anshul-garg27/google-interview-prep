package listeners

import (
	"coffee/app/campaignprofiles/domain"
	"coffee/client/winkl"
	"coffee/core/client"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
	log "github.com/sirupsen/logrus"

	cpmanager "coffee/app/campaignprofiles/manager"
)

func SetupListeners(router *message.Router, subscriber *amqp.Subscriber) {
	router.AddNoPublisherHandler(
		"sync_winkl_api",
		"winkl.dx___sync_winkl_api_q",
		subscriber,
		UpdateWinklApiUsingId,
	)
}

type Client struct {
	*client.BaseClient
}

func UpdateWinklApiUsingId(message *message.Message) error {
	ctx := message.Context()
	cpId := &winkl.SyncWinklApiEntry{}
	err := json.Unmarshal(message.Payload, cpId)

	if cpId.Id == nil {
		return nil
	}
	campaignManager := cpmanager.CreateManager(ctx)
	campaignEntry, err := campaignManager.FindById(ctx, *cpId.Id, true)

	if err != nil {
		log.Error("entry not found on coffee side")
		return nil
	}

	var onGcc, instantGratificationInvited, amazonStoreLinkVerified *bool
	var isBlacklisted, whatsappOptin, ekycPending *bool
	var gender, name, email, blacklistedBy, blacklistedReason, city, state, country, bio, dob, phone, accountId, secondaryPhone, notificationToken, memberId, referenceCode, webengageUserId *string
	var creatorProgramArary *[]winkl.CreatorProgramInput
	var creatorCohortArray *[]winkl.CreatorCohortInput
	var campaignCategoriesArray *[]int64
	var socialAccounts []winkl.SocialAccountInput
	var languages *[]string

	if campaignEntry.AccountId != nil {
		strValue := strconv.FormatInt(*campaignEntry.AccountId, 10)
		accountId = &strValue
	}
	if campaignEntry.AdminDetails != nil {

		adminDetails := *campaignEntry.AdminDetails
		if adminDetails.IsBlacklisted != nil {
			isBlacklisted = adminDetails.IsBlacklisted
		}
		if adminDetails.BlacklistedBy != nil {
			blacklistedBy = adminDetails.BlacklistedBy
		}
		if adminDetails.BlacklistedReason != nil {
			blacklistedReason = adminDetails.BlacklistedReason
		}
		if adminDetails.Name != nil {
			name = adminDetails.Name
		}
		if adminDetails.Email != nil {
			email = adminDetails.Email
		}
		if adminDetails.Phone != nil {
			phone = adminDetails.Phone
		}
		if adminDetails.City != nil {
			city = adminDetails.City
		}
		if adminDetails.State != nil {
			state = adminDetails.State
		}
		if adminDetails.Country != nil {
			country = adminDetails.Country
		}
		if adminDetails.CampaignCategoryIds != nil {
			if len(*adminDetails.CampaignCategoryIds) == 0 {
				emptySlice := make([]int64, 0)
				campaignCategoriesArray = &emptySlice
			} else {
				var temporaryArray []int64
				temporaryArray = append(temporaryArray, *adminDetails.CampaignCategoryIds...)
				campaignCategoriesArray = &temporaryArray
			}
		}
		if adminDetails.Languages != nil {
			languages = domain.TransformCodeToLanguages(adminDetails.Languages)
		}
		if adminDetails.Dob != nil {
			dob = adminDetails.Dob
		}
		if adminDetails.Bio != nil {
			bio = adminDetails.Bio
		}
		if adminDetails.Gender != nil {
			gender = adminDetails.Gender
		}

		if adminDetails.CreatorPrograms != nil {
			if len(*adminDetails.CreatorPrograms) == 0 {
				emptySlice := make([]winkl.CreatorProgramInput, 0)
				creatorProgramArary = &emptySlice
			} else {
				var temporaryArray []winkl.CreatorProgramInput
				for _, program := range *adminDetails.CreatorPrograms {
					creator := winkl.CreatorProgramInput{
						Id:    winkl.TransformCategoryProgram(program.Tag),
						Tag:   program.Tag,
						Level: winkl.TransformCategoryProgramLevelToCamelCase(program.Level),
					}
					temporaryArray = append(temporaryArray, creator)
				}
				creatorProgramArary = &temporaryArray
			}
		}
		if adminDetails.CreatorCohorts != nil {
			if len(*adminDetails.CreatorCohorts) == 0 {
				emptySlice := make([]winkl.CreatorCohortInput, 0)
				creatorCohortArray = &emptySlice
			} else {
				var temporaryArray []winkl.CreatorCohortInput
				for _, cohort := range *adminDetails.CreatorCohorts {
					creator := winkl.CreatorCohortInput{
						Id:    winkl.TransformCohortProgram(cohort.Tag),
						Tag:   cohort.Tag,
						Level: cohort.Level,
					}
					temporaryArray = append(temporaryArray, creator)
				}
				creatorCohortArray = &temporaryArray
			}
		}
		if adminDetails.WhatsappOptIn != nil {
			whatsappOptin = adminDetails.WhatsappOptIn
		}
		if adminDetails.SecondaryPhone != nil {
			secondaryPhone = adminDetails.SecondaryPhone
		}
	}

	// have to fix this part.

	// if campaignEntry.Socials != nil {
	// 	for _, social := range campaignEntry.Socials {

	// 		var platform winkl.SocialNetworkType

	// 		id := strconv.FormatInt(*cpId.Id, 10)
	// 		cpid := id
	// 		if social.Platform == string(constants.InstagramPlatform) {
	// 			platform = winkl.SocialNetworkTypeInstagram
	// 		} else if social.Platform == string(constants.YoutubePlatform) {
	// 			platform = winkl.SocialNetworkTypeYoutube
	// 		} else {
	// 			platform = winkl.SocialNetworkTypeGcc
	// 		}
	// 		var handle string
	// 		if social.Handle != nil {
	// 			handle = *social.Handle
	// 		}
	// 		socialAccount := &winkl.SocialAccountInput{
	// 			Handle:   handle,
	// 			Platform: platform,
	// 			CpId:     cpid,
	// 		}

	// 		socialAccounts = append(socialAccounts, *socialAccount)
	// 	}
	// }

	if campaignEntry.OnGCC != nil {
		onGcc = campaignEntry.OnGCC
	}
	if campaignEntry.UserDetails != nil {
		userDetails := *campaignEntry.UserDetails
		if userDetails.AmazonStoreLinkVerified != nil {
			amazonStoreLinkVerified = userDetails.AmazonStoreLinkVerified
		}
		if userDetails.InstantGratificationInvited != nil {
			instantGratificationInvited = userDetails.InstantGratificationInvited
		}
		if userDetails.NotificationToken != nil {
			notificationToken = userDetails.NotificationToken
		}
		if userDetails.MemberId != nil {
			memberId = userDetails.MemberId
		}
		if userDetails.ReferenceCode != nil {
			referenceCode = userDetails.ReferenceCode
		}
		if userDetails.EKYCPending != nil {
			ekycPending = userDetails.EKYCPending
		}
		if userDetails.WebEngageUserId != nil {
			webengageUserId = userDetails.WebEngageUserId
		}
	}

	winklInput := &winkl.WinklUpdateInfluencerInput{
		IsBlacklisted:               isBlacklisted,
		BlacklistedBy:               blacklistedBy,
		BlacklistedReason:           blacklistedReason,
		Name:                        name,
		Email:                       email,
		PhoneNumber:                 phone,
		City:                        city,
		State:                       state,
		Country:                     country,
		UserCategories:              campaignCategoriesArray,
		Languages:                   languages,
		Dob:                         dob,
		Bio:                         bio,
		Gender:                      gender,
		SocialAccounts:              socialAccounts,
		CreatorPrograms:             creatorProgramArary,
		CreatorCohorts:              creatorCohortArray,
		AccountId:                   accountId,
		WhatsappOptin:               whatsappOptin,
		OnGCC:                       onGcc,
		InstantGratificationInvited: instantGratificationInvited,
		AmazonStoreLinkVerified:     amazonStoreLinkVerified,
		SecondaryPhone:              secondaryPhone,
		NotificationToken:           notificationToken,
		MemberId:                    memberId,
		ReferenceCode:               referenceCode,
		EKYCPending:                 ekycPending,
		WebEngageUserId:             webengageUserId,
	}

	winklRequestBodyJSON, _ := json.Marshal(winklInput)
	log.Debug("winkl input ", string(winklRequestBodyJSON))

	client := &Client{client.New(ctx, "WINKL_URL")}

	client.SetTimeout(1 * time.Minute)
	response, err := client.GenerateRequest().
		SetBody(winklInput).
		Post(client.BaseUrl + "/social_data_service/update")

	log.Debug("winkl response ", response)
	if err != nil {
		err = fmt.Errorf("%s winkl input %s", err, string(winklRequestBodyJSON))
		return nil
	}

	var responseBody map[string]interface{}
	err = json.Unmarshal(response.Body(), &responseBody)
	if err != nil {
		log.Debug("error in Unmarshling response", err)
	}
	status, _ := responseBody["status"].(map[string]interface{})
	statusType, _ := status["type"].(string)

	if statusType == "ERROR" {
		code := strconv.FormatFloat(status["code"].(float64), 'f', -1, 64)
		message := status["message"].(string)
		errorMessage := "Error calling winkl API. Error code: " + code + ". Error message: " + message
		log.Error(errorMessage)
		return nil
	}

	return nil
}
