package winkl

import (
	"coffee/core/client"
	coredomain "coffee/core/domain"
	"context"
	"strconv"
)

type Client struct {
	*client.BaseClient
}

func New(ctx context.Context) *Client {
	return &Client{client.New(ctx, "WINKL_URL")}
}

// func (c *Client) UpdateInfluencerData(results []coredomain.CampaignProfileEntry, input coredomain.CampaignProfileEntry) error {
// 	var onGcc, instantGratificationInvited, amazonStoreLinkVerified, isBlacklisted, whatsappOptin, ekycPending *bool
// 	var gender, name, email, blacklistedBy, blacklistedReason, city, state, country, bio, dob, phone, accountId, secondaryPhone, notificationToken, memberId, referenceCode, webengageUserId *string
// 	var creatorProgramArary *[]CreatorProgramInput
// 	var creatorCohortArray *[]CreatorCohortInput
// 	var campaignCategoriesArray *[]int64
// 	var socialAccounts []SocialAccountInput
// 	var languages *[]string

// 	cpIdMap := getCpIdFromResults(results)
// 	log.Debug(cpIdMap)

// 	if input.AccountId != nil {
// 		strValue := strconv.FormatInt(*input.AccountId, 10)
// 		accountId = &strValue
// 	}
// 	if input.AdminDetails != nil {

// 		adminDetails := *input.AdminDetails
// 		if adminDetails.IsBlacklisted != nil {
// 			isBlacklisted = adminDetails.IsBlacklisted
// 		}
// 		if adminDetails.BlacklistedBy != nil {
// 			blacklistedBy = adminDetails.BlacklistedBy
// 		}
// 		if adminDetails.BlacklistedReason != nil {
// 			blacklistedReason = adminDetails.BlacklistedReason
// 		}
// 		if adminDetails.Name != nil {
// 			name = adminDetails.Name
// 		}
// 		if adminDetails.Email != nil {
// 			email = adminDetails.Email
// 		}
// 		if adminDetails.Phone != nil {
// 			phone = adminDetails.Phone
// 		}
// 		if adminDetails.City != nil {
// 			city = adminDetails.City
// 		}
// 		if adminDetails.State != nil {
// 			state = adminDetails.State
// 		}
// 		if adminDetails.Country != nil {
// 			country = adminDetails.Country
// 		}
// 		if adminDetails.CampaignCategoryIds != nil {
// 			if len(*adminDetails.CampaignCategoryIds) == 0 {
// 				emptySlice := make([]int64, 0)
// 				campaignCategoriesArray = &emptySlice
// 			} else {
// 				var temporaryArray []int64
// 				// for _, campaignCategory := range *adminDetails.CampaignCategoryIds {
// 				// 	// num, _ := strconv.Atoi(campaignCategory)
// 				// 	temporaryArray = append(temporaryArray, campaignCategory)
// 				// }
// 				temporaryArray = append(temporaryArray, *adminDetails.CampaignCategoryIds...)
// 				campaignCategoriesArray = &temporaryArray
// 			}
// 		}
// 		if adminDetails.Languages != nil {
// 			languages = cpdao.TransformCodeToLanguages(adminDetails.Languages)
// 		}
// 		if adminDetails.Dob != nil {
// 			dob = adminDetails.Dob
// 		}
// 		if adminDetails.Bio != nil {
// 			bio = adminDetails.Bio
// 		}
// 		if adminDetails.Gender != nil {
// 			gender = adminDetails.Gender
// 		}

// 		if adminDetails.CreatorPrograms != nil {
// 			if len(*adminDetails.CreatorPrograms) == 0 {
// 				emptySlice := make([]CreatorProgramInput, 0)
// 				creatorProgramArary = &emptySlice
// 			} else {
// 				var temporaryArray []CreatorProgramInput
// 				for _, program := range *adminDetails.CreatorPrograms {
// 					creator := CreatorProgramInput{
// 						Id:    TransformCategoryProgram(program.Tag),
// 						Tag:   program.Tag,
// 						Level: TransformCategoryProgramLevelToCamelCase(program.Level),
// 					}
// 					temporaryArray = append(temporaryArray, creator)
// 				}
// 				creatorProgramArary = &temporaryArray
// 			}
// 		}
// 		if adminDetails.CreatorCohorts != nil {
// 			if len(*adminDetails.CreatorCohorts) == 0 {
// 				emptySlice := make([]CreatorCohortInput, 0)
// 				creatorCohortArray = &emptySlice
// 			} else {
// 				var temporaryArray []CreatorCohortInput
// 				for _, cohort := range *adminDetails.CreatorCohorts {
// 					creator := CreatorCohortInput{
// 						Id:    TransformCohortProgram(cohort.Tag),
// 						Tag:   cohort.Tag,
// 						Level: cohort.Level,
// 					}
// 					temporaryArray = append(temporaryArray, creator)
// 				}
// 				creatorCohortArray = &temporaryArray
// 			}
// 		}
// 		if adminDetails.WhatsappOptIn != nil {
// 			whatsappOptin = adminDetails.WhatsappOptIn
// 		}
// 		if adminDetails.SecondaryPhone != nil {
// 			secondaryPhone = adminDetails.SecondaryPhone
// 		}
// 	}

// 	if input.Socials != nil {
// 		for _, social := range input.Socials {

// 			var platform SocialNetworkType
// 			var cpid string
// 			if social.Platform == "INSTAGRAM" {
// 				platform = SocialNetworkTypeInstagram
// 				cpid = cpIdMap[social.Platform]
// 			} else if social.Platform == "YOUTUBE" {
// 				platform = SocialNetworkTypeYoutube
// 				cpid = cpIdMap[social.Platform]
// 			} else {
// 				platform = SocialNetworkTypeGcc
// 				cpid = cpIdMap[social.Platform]
// 			}
// 			var handle string
// 			if social.Handle != nil {
// 				handle = *social.Handle
// 			}
// 			socialAccount := &SocialAccountInput{
// 				Handle:   handle,
// 				Platform: platform,
// 				CpId:     cpid,
// 			}

// 			socialAccounts = append(socialAccounts, *socialAccount)
// 		}
// 	}

// 	if input.OnGCC != nil {
// 		onGcc = input.OnGCC
// 	}
// 	if input.UserDetails != nil {
// 		userDetails := *input.UserDetails
// 		if userDetails.Name != nil {
// 			name = userDetails.Name
// 		}
// 		if userDetails.Email != nil {
// 			email = userDetails.Email
// 		}
// 		if userDetails.Phone != nil {
// 			phone = userDetails.Phone
// 		}
// 		if userDetails.CampaignCategoryIds != nil {
// 			if len(*userDetails.CampaignCategoryIds) == 0 {
// 				emptySlice := make([]int64, 0)
// 				campaignCategoriesArray = &emptySlice
// 			} else {
// 				var temporaryArray []int64
// 				// for _, campaignCategory := range *userDetails.CampaignCategoryIds {
// 				// 	// num, _ := strconv.Atoi(campaignCategory)
// 				// 	temporaryArray = append(temporaryArray, campaignCategory)
// 				// }
// 				temporaryArray = append(temporaryArray, *userDetails.CampaignCategoryIds...)
// 				campaignCategoriesArray = &temporaryArray
// 			}
// 		}
// 		if userDetails.Languages != nil {
// 			languages = cpdao.TransformCodeToLanguages(userDetails.Languages)
// 		}
// 		if userDetails.Dob != nil {
// 			dob = userDetails.Dob
// 		}
// 		if userDetails.Bio != nil {
// 			bio = userDetails.Bio
// 		}
// 		if userDetails.Gender != nil {
// 			gender = userDetails.Gender
// 		}
// 		if userDetails.WhatsappOptIn != nil {
// 			whatsappOptin = userDetails.WhatsappOptIn
// 		}
// 		if userDetails.AmazonStoreLinkVerified != nil {
// 			amazonStoreLinkVerified = userDetails.AmazonStoreLinkVerified
// 		}
// 		if userDetails.InstantGratificationInvited != nil {
// 			instantGratificationInvited = userDetails.InstantGratificationInvited
// 		}
// 		if userDetails.SecondaryPhone != nil {
// 			secondaryPhone = userDetails.SecondaryPhone
// 		}
// 		if userDetails.NotificationToken != nil {
// 			notificationToken = userDetails.NotificationToken
// 		}
// 		if userDetails.MemberId != nil {
// 			memberId = userDetails.MemberId
// 		}
// 		if userDetails.ReferenceCode != nil {
// 			referenceCode = userDetails.ReferenceCode
// 		}
// 		if userDetails.EKYCPending != nil {
// 			ekycPending = userDetails.EKYCPending
// 		}
// 		if userDetails.WebEngageUserId != nil {
// 			webengageUserId = userDetails.WebEngageUserId
// 		}
// 	}

// 	winklInput := &WinklUpdateInfluencerInput{
// 		IsBlacklisted:               isBlacklisted,
// 		BlacklistedBy:               blacklistedBy,
// 		BlacklistedReason:           blacklistedReason,
// 		Name:                        name,
// 		Email:                       email,
// 		PhoneNumber:                 phone,
// 		City:                        city,
// 		State:                       state,
// 		Country:                     country,
// 		UserCategories:              campaignCategoriesArray,
// 		Languages:                   languages,
// 		Dob:                         dob,
// 		Bio:                         bio,
// 		Gender:                      gender,
// 		SocialAccounts:              socialAccounts,
// 		CreatorPrograms:             creatorProgramArary,
// 		CreatorCohorts:              creatorCohortArray,
// 		AccountId:                   accountId,
// 		WhatsappOptin:               whatsappOptin,
// 		OnGCC:                       onGcc,
// 		InstantGratificationInvited: instantGratificationInvited,
// 		AmazonStoreLinkVerified:     amazonStoreLinkVerified,
// 		SecondaryPhone:              secondaryPhone,
// 		NotificationToken:           notificationToken,
// 		MemberId:                    memberId,
// 		ReferenceCode:               referenceCode,
// 		EKYCPending:                 ekycPending,
// 		WebEngageUserId:             webengageUserId,
// 	}

// 	winklRequestBodyJSON, _ := json.Marshal(winklInput)
// 	log.Debug("winkl input ", string(winklRequestBodyJSON))

// 	c.SetTimeout(1 * time.Minute)
// 	startTime := time.Now()
// 	response, err := c.GenerateRequest().
// 		SetBody(winklInput).
// 		Post(c.BaseUrl + "/social_data_service/update")
// 	endTime := time.Now()
// 	duration := endTime.Sub(startTime)
// 	log.Debug("winkl response ", response)
// 	if err != nil {
// 		err = fmt.Errorf("%s winkl input %s. Time taken is %s", err, string(winklRequestBodyJSON), duration)
// 		return err
// 	}

// 	var responseBody map[string]interface{}
// 	err = json.Unmarshal(response.Body(), &responseBody)
// 	if err != nil {
// 		log.Debug("error in Unmarshling response", err)
// 	}
// 	status, _ := responseBody["status"].(map[string]interface{})
// 	statusType, _ := status["type"].(string)

// 	if statusType == "ERROR" {
// 		code := strconv.FormatFloat(status["code"].(float64), 'f', -1, 64)
// 		message := status["message"].(string)
// 		errorMessage := "Error calling winkl API. Error code: " + code + ". Error message: " + message
// 		return errors.New(errorMessage)
// 	}

// 	return nil
// }

func TransformCategoryProgram(creatorProgram string) int {
	categories := map[string]int{
		"BEAUTY":         1,
		"GOOD_PARENTING": 2,
		"GOOD_LIFE":      3,
	}

	value := categories[creatorProgram]
	return value
}

func TransformCategoryProgramLevelToCamelCase(creatorProgram string) string {
	categories := map[string]string{
		"PRO":     "Pro",
		"CREATOR": "Creator",
		"SQUAD":   "Squad",
	}

	value := categories[creatorProgram]
	return value
}

func TransformCohortProgram(creatorCohort string) int {
	cohorts := map[string]int{
		"CONTENT": 1,
		"ORDERS":  2,
		"VIEWS":   3,
	}

	value := cohorts[creatorCohort]
	return value
}

func getCpIdFromResults(results []coredomain.CampaignProfileEntry) map[string]string {

	cpIdMap := make(map[string]string)
	for _, campaignProfile := range results {
		if campaignProfile.Platform == "INSTAGRAM" {
			cpIdMap["INSTAGRAM"] = strconv.FormatInt(campaignProfile.Id, 10)
		} else if campaignProfile.Platform == "YOUTUBE" {
			cpIdMap["YOUTUBE"] = strconv.FormatInt(campaignProfile.Id, 10)
		} else if campaignProfile.Platform == "GCC" {
			cpIdMap["GCC"] = strconv.FormatInt(campaignProfile.Id, 10)
		}
	}

	return cpIdMap
}
