package domain

import (
	"coffee/app/campaignprofiles/dao"
	discoverydao "coffee/app/discovery/dao"
	beatservice "coffee/client/beat"
	"coffee/constants"
	"coffee/core/appcontext"
	coredomain "coffee/core/domain"
	"coffee/helpers"
	"context"
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/imdario/mergo"
	log "github.com/sirupsen/logrus"
)

func ConvertCreatorProgramsToString(adminDetails *coredomain.AdminDetails) *[]string {
	creatorPrograms := &[]string{}
	if adminDetails.CreatorPrograms != nil {
		for _, program := range *adminDetails.CreatorPrograms {
			cp := strconv.Itoa(constants.CreatorProgramsNameToId[program.Tag]) + "_" + program.Level
			*creatorPrograms = append(*creatorPrograms, cp)
		}
		return creatorPrograms
	}
	return nil
}

func ConvertCreatorCohortsToString(adminDetails *coredomain.AdminDetails) *[]string {
	creatorCohorts := &[]string{}
	if adminDetails.CreatorCohorts != nil {
		for _, cohort := range *adminDetails.CreatorCohorts {
			cp := strconv.Itoa(constants.CreatorCohortsNameToId[cohort.Tag]) + "_" + cohort.Level
			*creatorCohorts = append(*creatorCohorts, cp)
		}
		return creatorCohorts
	}
	return nil
}

func convertCampaignCategoryIdsToString(CampaignCategoryIds *[]int64) *[]string {
	campaignCategoryIdStrings := &[]string{}
	for _, val := range *CampaignCategoryIds {
		id := strconv.FormatInt(val, 10)
		*campaignCategoryIdStrings = append(*campaignCategoryIdStrings, id)
	}
	return campaignCategoryIdStrings
}

func ToAdminDetailsEntity(adminDetails *coredomain.AdminDetails) dao.AdminDetailsEntity {
	var whatsappOptIn, isBlacklisted bool

	if adminDetails.WhatsappOptIn != nil {
		whatsappOptIn = *adminDetails.WhatsappOptIn
	}
	if adminDetails.IsBlacklisted != nil {
		isBlacklisted = *adminDetails.IsBlacklisted
	}

	return dao.AdminDetailsEntity{
		Name:              adminDetails.Name,
		Email:             adminDetails.Email,
		Phone:             adminDetails.Phone,
		SecondaryPhone:    adminDetails.SecondaryPhone,
		Gender:            adminDetails.Gender,
		Languages:         adminDetails.Languages,
		City:              adminDetails.City,
		State:             adminDetails.State,
		Country:           adminDetails.Country,
		Dob:               adminDetails.Dob,
		Bio:               adminDetails.Bio,
		IsBlacklisted:     isBlacklisted,
		BlacklistedBy:     adminDetails.BlacklistedBy,
		BlacklistedReason: adminDetails.BlacklistedReason,
		Location:          adminDetails.Location,
		WhatsappOptIn:     whatsappOptIn,
		CountryCode:       adminDetails.CountryCode,
	}
}

func ToUserDetailsEntity(userDetails *coredomain.UserDetails) dao.UserDetailsEntity {
	var whatsappOptIn, ekycPending, instantGratificationInvited, amazonStoreLinkVerified bool
	if userDetails.WhatsappOptIn != nil {
		whatsappOptIn = *userDetails.WhatsappOptIn
	}
	if userDetails.EKYCPending != nil {
		ekycPending = *userDetails.EKYCPending
	}
	if userDetails.InstantGratificationInvited != nil {
		instantGratificationInvited = *userDetails.InstantGratificationInvited
	}
	if userDetails.AmazonStoreLinkVerified != nil {
		amazonStoreLinkVerified = *userDetails.AmazonStoreLinkVerified
	}
	return dao.UserDetailsEntity{
		Email:                       userDetails.Email,
		Phone:                       userDetails.Phone,
		SecondaryPhone:              userDetails.SecondaryPhone,
		Name:                        userDetails.Name,
		Dob:                         userDetails.Dob,
		Gender:                      userDetails.Gender,
		WhatsappOptIn:               whatsappOptIn,
		Languages:                   userDetails.Languages,
		Location:                    userDetails.Location,
		NotificationToken:           userDetails.NotificationToken,
		Bio:                         userDetails.Bio,
		MemberId:                    userDetails.MemberId,
		ReferenceCode:               userDetails.ReferenceCode,
		EKYCPending:                 ekycPending,
		WebEngageUserId:             userDetails.WebEngageUserId,
		InstantGratificationInvited: instantGratificationInvited,
		AmazonStoreLinkVerified:     amazonStoreLinkVerified,
	}
}

func ToCampaignEntity(cpInput *coredomain.CampaignProfileInput) (*dao.CampaignProfileEntity, error) {
	var adminDetailsString, userDetailsString *string
	var adminDetailsEntity dao.AdminDetailsEntity
	var userDetailsEntity dao.UserDetailsEntity

	if cpInput.AdminDetails != nil {
		adminDetailsEntity = ToAdminDetailsEntity(cpInput.AdminDetails)
		adminDetailsEntity = convertAdminDetailsCreatorProgramsCohortsAndCategoryIdsToString(cpInput.AdminDetails, adminDetailsEntity)
	}
	adminDetailsString = getJSONString(adminDetailsEntity)

	if cpInput.UserDetails != nil {
		userDetailsEntity = ToUserDetailsEntity(cpInput.UserDetails)
		userDetailsEntity = convertUserDetailsCategoryIdsToString(cpInput.UserDetails, userDetailsEntity)
	}
	userDetailsString = getJSONString(userDetailsEntity)

	return &dao.CampaignProfileEntity{
		GCCUserAccountID: cpInput.AccountId,
		OnGCC:            cpInput.OnGCC,
		OnGCCApp:         cpInput.OnGCCApp,
		AdminDetails:     adminDetailsString,
		UserDetails:      userDetailsString,
		Enabled:          true,
		UpdatedBy:        cpInput.UpdatedBy,
	}, nil
}

func convertAdminDetailsCreatorProgramsCohortsAndCategoryIdsToString(entryAdminDetails *coredomain.AdminDetails, entityAdminDetails dao.AdminDetailsEntity) dao.AdminDetailsEntity {
	if entryAdminDetails.Languages != nil && len(*entryAdminDetails.Languages) == 0 {
		entityAdminDetails.Languages = nil
	}
	if entryAdminDetails.CampaignCategoryIds != nil && len(*entryAdminDetails.CampaignCategoryIds) == 0 {
		entityAdminDetails.CampaignCategoryIds = nil
	}
	if entryAdminDetails.CreatorPrograms != nil {
		entityAdminDetails.CreatorPrograms = ConvertCreatorProgramsToString(entryAdminDetails)
	}
	if entryAdminDetails.CreatorCohorts != nil {
		entityAdminDetails.CreatorCohorts = ConvertCreatorCohortsToString(entryAdminDetails)
	}
	if entryAdminDetails.CampaignCategoryIds != nil {
		entityAdminDetails.CampaignCategoryIds = convertCampaignCategoryIdsToString(entryAdminDetails.CampaignCategoryIds)
	}
	return entityAdminDetails
}

func convertUserDetailsCategoryIdsToString(entryUserDetails *coredomain.UserDetails, entityUserDetails dao.UserDetailsEntity) dao.UserDetailsEntity {
	if entryUserDetails.Languages != nil && len(*entryUserDetails.Languages) == 0 {
		entityUserDetails.Languages = nil
	}
	if entryUserDetails.CampaignCategoryIds != nil && len(*entryUserDetails.CampaignCategoryIds) == 0 {
		entityUserDetails.CampaignCategoryIds = nil
	}
	if entryUserDetails.Location != nil && len(*entryUserDetails.Location) == 0 {
		entityUserDetails.Location = nil
	}
	if entryUserDetails.CampaignCategoryIds != nil {
		entityUserDetails.CampaignCategoryIds = convertCampaignCategoryIdsToString(entryUserDetails.CampaignCategoryIds)
	}
	return entityUserDetails
}

func ConvertStringToAdminDetails(adminDetailsString *string) coredomain.AdminDetails {
	var adminDetails coredomain.AdminDetails
	var creatorStructs coredomain.CreatorStructs
	json.Unmarshal([]byte(*adminDetailsString), &adminDetails)
	json.Unmarshal([]byte(*adminDetailsString), &creatorStructs)
	adminDetails = ConvertStringToCreatorPrograms(creatorStructs, adminDetails)
	adminDetails = ConvertStringToCreatorCohorts(creatorStructs, adminDetails)
	adminDetails.CampaignCategoryIds = convertStringToCampaignCategoryIds(creatorStructs)
	return adminDetails
}

func ConvertStringToUserDetails(userDetailsSting *string) coredomain.UserDetails {
	var userDetails coredomain.UserDetails
	var creatorStructs coredomain.CreatorStructs
	json.Unmarshal([]byte(*userDetailsSting), &userDetails)
	json.Unmarshal([]byte(*userDetailsSting), &creatorStructs)
	userDetails.CampaignCategoryIds = convertStringToCampaignCategoryIds(creatorStructs)
	return userDetails
}

func ToCampaignEntry(ctx context.Context, entity *dao.CampaignProfileEntity, profile *coredomain.Profile) (*coredomain.CampaignProfileEntry, error) {
	var adminDetails coredomain.AdminDetails
	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)

	if entity.AdminDetails != nil {
		adminDetails = ConvertStringToAdminDetails(entity.AdminDetails)
	}
	var userDetails coredomain.UserDetails
	if entity.UserDetails != nil {
		userDetails = ConvertStringToUserDetails(entity.UserDetails)
	}

	// Encoding CampaignProfile Contact Info For Non Admin
	if appCtx.PartnerId == nil || *appCtx.PartnerId != -1 {
		adminDetails, userDetails = hideContactInfo(adminDetails, userDetails)
	}
	// If Plan Is paid and cp Contact Info is nil Then Replace With Profile Info
	if appCtx.PlanType != nil && (*appCtx.PlanType == constants.PaidPlan || *appCtx.PlanType == constants.SaasPlan) {
		adminDetails = enrichContactInfoFromSocials(adminDetails, profile)

	}
	campaignProfileEntry := coredomain.CampaignProfileEntry{
		Id:                entity.ID,
		Platform:          entity.Platform,
		PlatformAccountId: entity.PlatformAccountId,
		OnGCC:             entity.OnGCC,
		OnGCCApp:          entity.OnGCCApp,
		AdminDetails:      &adminDetails,
		UserDetails:       &userDetails,
		UpdatedBy:         entity.UpdatedBy,
		AccountId:         entity.GCCUserAccountID,
		HasEmail:          entity.HasEmail,
		HasPhone:          entity.HasPhone,
	}
	if profile != nil {
		campaignProfileEntry = *EnrichCampaignProfileFromSocialProfile(&campaignProfileEntry, profile)
	}
	return &campaignProfileEntry, nil
}

func hideContactInfo(adminDetails coredomain.AdminDetails, userDetails coredomain.UserDetails) (coredomain.AdminDetails, coredomain.UserDetails) {
	if adminDetails.Email != nil {
		encodedEmail := helpers.EncodeEmail(*adminDetails.Email)
		adminDetails.Email = &encodedEmail
	}
	if adminDetails.Phone != nil {
		encodedPhone := helpers.EncodePhone(*adminDetails.Phone)
		adminDetails.Phone = &encodedPhone
	}
	if adminDetails.SecondaryPhone != nil {
		encodedPhone := helpers.EncodePhone(*adminDetails.SecondaryPhone)
		adminDetails.SecondaryPhone = &encodedPhone
	}

	if userDetails.Email != nil {
		encodedEmail := helpers.EncodeEmail(*userDetails.Email)
		userDetails.Email = &encodedEmail
	}
	if userDetails.Phone != nil {
		encodedPhone := helpers.EncodePhone(*userDetails.Phone)
		userDetails.Phone = &encodedPhone
	}
	if userDetails.SecondaryPhone != nil {
		encodedPhone := helpers.EncodePhone(*userDetails.SecondaryPhone)
		userDetails.SecondaryPhone = &encodedPhone
	}
	return adminDetails, userDetails
}

func EnrichCampaignProfileFromSocialProfile(entry *coredomain.CampaignProfileEntry, profile *coredomain.Profile) *coredomain.CampaignProfileEntry {
	entry.Platform = profile.Platform
	entry.Handle = profile.Handle
	for j := range profile.LinkedSocials {
		if profile.LinkedSocials[j].Platform == entry.Platform && profile.LinkedSocials[j].Source == string(constants.GCCPlatform) && strconv.FormatInt(profile.LinkedSocials[j].Code, 10) == profile.PlatformCode {
			// entry.Socials = append(entry.Socials, profile.LinkedSocials[j])
			entry.SocialDetails = profile.LinkedSocials[j]
		}
	}
	return entry
}
func enrichContactInfoFromSocials(adminDetails coredomain.AdminDetails, profile *coredomain.Profile) coredomain.AdminDetails {
	if adminDetails.Email == nil && profile != nil && profile.Email != nil {
		adminDetails.Email = profile.Email
	}
	if adminDetails.Phone == nil && profile != nil && profile.Email != nil {
		adminDetails.Phone = profile.Phone
	}
	return adminDetails

}

func ConvertStringToCreatorPrograms(creatorStructs coredomain.CreatorStructs, adminDetails coredomain.AdminDetails) coredomain.AdminDetails {
	if adminDetails.CreatorPrograms == nil {
		return adminDetails
	}
	programs := make([]coredomain.CreatorProgram, 0)
	adminDetails.CreatorPrograms = &programs
	if creatorStructs.CreatorProgram != nil {
		for _, program := range *creatorStructs.CreatorProgram {
			var cp coredomain.CreatorProgram
			parts := strings.Split(program, "_")
			if parts[0] != "" && parts[1] != "" {
				tag, _ := strconv.Atoi(parts[0])
				level := parts[1]

				cp.Tag = constants.CreatorProgramsIdToName[tag]
				cp.Level = level
				*adminDetails.CreatorPrograms = append(*adminDetails.CreatorPrograms, cp)
			} else {
				log.Error("creatorPrograms value is not correct")
			}
		}
	}
	return adminDetails
}

func ConvertStringToCreatorCohorts(creatorStructs coredomain.CreatorStructs, adminDetails coredomain.AdminDetails) coredomain.AdminDetails {
	if adminDetails.CreatorCohorts == nil {
		return adminDetails
	}
	cohorts := make([]coredomain.CreatorCohort, 0)
	adminDetails.CreatorCohorts = &cohorts
	if creatorStructs.CreatorCohorts != nil {
		for _, cohort := range *creatorStructs.CreatorCohorts {
			var cc coredomain.CreatorCohort
			parts := strings.Split(cohort, "_")
			if parts[0] != "" && parts[1] != "" {
				tag, _ := strconv.Atoi(parts[0])
				level := parts[1]

				cc.Tag = constants.CreatorCohortsIdToName[tag]
				cc.Level = level
				*adminDetails.CreatorCohorts = append(*adminDetails.CreatorCohorts, cc)
			} else {
				log.Error("creatorCohorts value is not correct")
			}
		}
	}
	return adminDetails
}

func CreateKeywordsAndSearchPhraseAdmin(adminDetails *coredomain.AdminDetails, userDetails *coredomain.UserDetails, platform string, profile *coredomain.Profile) (string, []string) {
	var keywords_admin []string
	var search_phrase_admin string

	if adminDetails != nil {
		if adminDetails.Name != nil && *adminDetails.Name != "" {
			keywords_admin = append(keywords_admin, strings.Fields(strings.ToLower(*adminDetails.Name))...)
			search_phrase_admin += strings.ToLower(*adminDetails.Name) + " "
		}
		if adminDetails.Email != nil && *adminDetails.Email != "" {
			keywords_admin = append(keywords_admin, strings.ToLower(*adminDetails.Email))
			search_phrase_admin += strings.ToLower(*adminDetails.Email) + " "
		}
		if adminDetails.Phone != nil && *adminDetails.Phone != "" {
			keywords_admin = append(keywords_admin, strings.ToLower(*adminDetails.Phone))
			search_phrase_admin += strings.ToLower(*adminDetails.Phone) + " "
		}
		if adminDetails.SecondaryPhone != nil && *adminDetails.SecondaryPhone != "" {
			keywords_admin = append(keywords_admin, strings.ToLower(*adminDetails.SecondaryPhone))
			search_phrase_admin += strings.ToLower(*adminDetails.SecondaryPhone) + " "
		}
	}
	if userDetails != nil {
		if userDetails.Name != nil && *userDetails.Name != "" {
			keywords_admin = append(keywords_admin, strings.Fields(strings.ToLower(*userDetails.Name))...)
			search_phrase_admin += strings.ToLower(*userDetails.Name) + " "
		}
		if userDetails.Email != nil && *userDetails.Email != "" {
			keywords_admin = append(keywords_admin, strings.ToLower(*userDetails.Email))
			search_phrase_admin += strings.ToLower(*userDetails.Email) + " "
		}
		if userDetails.Phone != nil && *userDetails.Phone != "" {
			keywords_admin = append(keywords_admin, strings.ToLower(*userDetails.Phone))
			search_phrase_admin += strings.ToLower(*userDetails.Phone) + " "
		}
		if userDetails.SecondaryPhone != nil && *userDetails.SecondaryPhone != "" {
			keywords_admin = append(keywords_admin, strings.ToLower(*userDetails.SecondaryPhone))
			search_phrase_admin += strings.ToLower(*userDetails.SecondaryPhone) + " "
		}
	}
	if platform == string(constants.InstagramPlatform) {
		search_phrase_admin, keywords_admin = CreateKeywordsAndSearchPhraseAdminForIA(keywords_admin, search_phrase_admin, profile)
		search_phrase_admin = strings.TrimSpace(search_phrase_admin)
		return search_phrase_admin, keywords_admin
	} else if platform == string(constants.YoutubePlatform) {
		search_phrase_admin, keywords_admin = CreateKeywordsAndSearchPhraseAdminForYA(keywords_admin, search_phrase_admin, profile)
		search_phrase_admin = strings.TrimSpace(search_phrase_admin)
		return search_phrase_admin, keywords_admin
	}
	return search_phrase_admin, keywords_admin
}

func CreateKeywordsAndSearchPhraseAdminForIA(keywords_admin []string, search_phrase_admin string, profile *coredomain.Profile) (string, []string) {
	keywords_admin_insta := keywords_admin
	if profile.Name != nil && *profile.Name != "" {
		keywords_admin_insta = append(keywords_admin_insta, strings.Fields(strings.ToLower(*profile.Name))...)
		search_phrase_admin += strings.ToLower(*profile.Name) + " "
	}
	if profile.Handle != nil && *profile.Handle != "" {
		keywords_admin_insta = append(keywords_admin_insta, strings.ToLower(*profile.Handle))
		search_phrase_admin += strings.ToLower(*profile.Handle) + " "
	}
	return search_phrase_admin, keywords_admin_insta
}

func CreateKeywordsAndSearchPhraseAdminForYA(keywords_admin []string, search_phrase_admin string, profile *coredomain.Profile) (string, []string) {
	keywords_admin_youtube := keywords_admin
	if profile.Name != nil && *profile.Name != "" {
		keywords_admin_youtube = append(keywords_admin_youtube, strings.Fields(strings.ToLower(*profile.Name))...)
		search_phrase_admin += strings.ToLower(*profile.Name) + " "
	}
	if profile.Username != nil && *profile.Username != "" {
		keywords_admin_youtube = append(keywords_admin_youtube, strings.ToLower(*profile.Username))
		search_phrase_admin += strings.ToLower(*profile.Username) + " "
	}
	if profile.Handle != nil {
		keywords_admin_youtube = append(keywords_admin_youtube, *profile.Handle)
		search_phrase_admin += *profile.Handle + " "
	}
	return search_phrase_admin, keywords_admin_youtube
}

func MergeUserAdminDetails(cpInput coredomain.CampaignProfileInput, cpEntity *dao.CampaignProfileEntity) (*coredomain.AdminDetails, *coredomain.UserDetails, dao.CampaignProfileEntity) {

	newAdminDetailsString, entryAdminDetails := MergeAdminDetails(cpInput, cpEntity)
	newUserDetailsString, entryUserDetails := MergeUserDetails(cpInput, cpEntity)

	if cpInput.OnGCC == nil {
		cpInput.OnGCC = cpEntity.OnGCC
	}
	if cpInput.OnGCCApp == nil {
		cpInput.OnGCCApp = cpEntity.OnGCCApp
	}

	gccUserAccountId := cpEntity.GCCUserAccountID
	if cpInput.AccountId != nil {
		gccUserAccountId = cpInput.AccountId
	}

	updatedCpEntity := dao.CampaignProfileEntity{
		ID:                cpEntity.ID,
		UpdatedBy:         cpInput.UpdatedBy,
		Platform:          cpEntity.Platform,
		PlatformAccountId: cpEntity.PlatformAccountId,
		AdminDetails:      newAdminDetailsString,
		UserDetails:       newUserDetailsString,
		OnGCC:             cpInput.OnGCC,
		OnGCCApp:          cpInput.OnGCCApp,
		GCCUserAccountID:  gccUserAccountId,
		HasEmail:          cpEntity.HasEmail,
		HasPhone:          cpEntity.HasPhone,
		Enabled:           cpEntity.Enabled,
		CreatedAt:         cpEntity.CreatedAt,
		UpdatedAt:         cpEntity.UpdatedAt,
	}

	return entryAdminDetails, entryUserDetails, updatedCpEntity
}

type blankTransformer struct{}

func (t blankTransformer) Transformer(typ reflect.Type) func(dst, src reflect.Value) error {
	switch typ.Kind() {
	case reflect.String:
		return func(dst, src reflect.Value) error {
			if dst.Kind() == reflect.Ptr && dst.IsNil() && !src.IsNil() {
				dst.SetString(src.String())
			} else if dst.Kind() == reflect.Ptr && dst.IsNil() && src.IsNil() {
				return nil
			} else {
				dst.SetString(dst.String())
			}
			return nil
		}
	case reflect.Int, reflect.Int64:
		return func(dst, src reflect.Value) error {
			if src.Int() == 0 {
				dst.SetInt(0)
			} else {
				dst.SetInt(src.Int())
			}
			return nil
		}
	case reflect.Ptr:
		return func(dst, src reflect.Value) error {
			if src.Kind() == reflect.Ptr && dst.Kind() == reflect.Ptr && dst.IsNil() && !src.IsNil() {
				dst.Elem().SetString(src.Elem().String())
			} else if src.Kind() == reflect.Ptr && dst.Kind() == reflect.Ptr && dst.IsNil() && src.IsNil() {
				return nil
			} else {
				if dst.Elem().Kind() == reflect.Slice || dst.Elem().Kind() == reflect.Bool {
					dst.Elem().Set(dst.Elem())
				} else {
					dst.Elem().SetString(dst.Elem().String())
				}
			}
			return nil
		}
	case reflect.Bool:
		return func(dst, _ reflect.Value) error {
			dst.SetBool(dst.Bool())
			return nil
		}
	default:
		return nil
	}
}

func MergeAdminDetails(cpInput coredomain.CampaignProfileInput, campaign *dao.CampaignProfileEntity) (*string, *coredomain.AdminDetails) {
	var tempEntryAdminDetails coredomain.AdminDetails
	encoded, _ := json.Marshal(cpInput.AdminDetails)
	json.Unmarshal(encoded, &tempEntryAdminDetails)

	var updatedCampaignAdminDetails coredomain.AdminDetails
	var entryAdminDetails *coredomain.AdminDetails

	if campaign.AdminDetails != nil {
		updatedCampaignAdminDetails = ConvertStringToAdminDetails(campaign.AdminDetails)
	}
	if cpInput.AdminDetails != nil {
		encoded, _ := json.Marshal(cpInput.AdminDetails)
		json.Unmarshal(encoded, &entryAdminDetails)
	} else {
		encoded, _ := json.Marshal(updatedCampaignAdminDetails)
		json.Unmarshal(encoded, &entryAdminDetails)
	}

	mergo.Merge(entryAdminDetails, updatedCampaignAdminDetails, mergo.WithOverride, mergo.WithTransformers(blankTransformer{}))

	entryAdminDetails = MergeAdminDetailArrays(tempEntryAdminDetails, cpInput, entryAdminDetails, updatedCampaignAdminDetails)
	entryAdminDetails.Location = SetLocation(entryAdminDetails)

	adminDetailsEntity := ToAdminDetailsEntity(entryAdminDetails)

	if entryAdminDetails.CreatorPrograms != nil && len(*entryAdminDetails.CreatorPrograms) != 0 {
		adminDetailsEntity.CreatorPrograms = ConvertCreatorProgramsToString(entryAdminDetails)
	}

	if entryAdminDetails.CreatorCohorts != nil && len(*entryAdminDetails.CreatorCohorts) != 0 {
		adminDetailsEntity.CreatorCohorts = ConvertCreatorCohortsToString(entryAdminDetails)
	}
	if entryAdminDetails.CampaignCategoryIds != nil {
		adminDetailsEntity.CampaignCategoryIds = convertCampaignCategoryIdsToString(entryAdminDetails.CampaignCategoryIds)
	}

	jsonData, _ := json.Marshal(adminDetailsEntity)
	jsonDataInString := string(jsonData)
	adminDetailsString := &jsonDataInString

	if *adminDetailsString == "{}" {
		return nil, entryAdminDetails
	}
	return adminDetailsString, entryAdminDetails
}

func MergeAdminDetailArrays(tempEntryAdminDetails coredomain.AdminDetails, cpInput coredomain.CampaignProfileInput, entryAdminDetails *coredomain.AdminDetails, updatedCampaignAdminDetails coredomain.AdminDetails) *coredomain.AdminDetails {
	if cpInput.AdminDetails != nil {
		if tempEntryAdminDetails.Languages != nil {
			if len(*tempEntryAdminDetails.Languages) == 0 {
				entryAdminDetails.Languages = nil
			} else {
				entryAdminDetails.Languages = tempEntryAdminDetails.Languages
			}
		} else {
			entryAdminDetails.Languages = updatedCampaignAdminDetails.Languages
		}

		if tempEntryAdminDetails.CampaignCategoryIds != nil {
			if len(*tempEntryAdminDetails.CampaignCategoryIds) == 0 {
				entryAdminDetails.CampaignCategoryIds = nil
			} else {
				entryAdminDetails.CampaignCategoryIds = tempEntryAdminDetails.CampaignCategoryIds
			}
		} else {
			entryAdminDetails.CampaignCategoryIds = updatedCampaignAdminDetails.CampaignCategoryIds
		}

		if tempEntryAdminDetails.CreatorPrograms != nil {
			if len(*tempEntryAdminDetails.CreatorPrograms) == 0 {
				entryAdminDetails.CreatorPrograms = nil
			} else {
				entryAdminDetails.CreatorPrograms = tempEntryAdminDetails.CreatorPrograms
			}
		} else {
			if updatedCampaignAdminDetails.CreatorPrograms == nil {
				entryAdminDetails.CreatorPrograms = nil
			} else {
				entryAdminDetails.CreatorPrograms = updatedCampaignAdminDetails.CreatorPrograms
			}
		}

		if tempEntryAdminDetails.CreatorCohorts != nil {
			if len(*tempEntryAdminDetails.CreatorCohorts) == 0 {
				entryAdminDetails.CreatorCohorts = nil
			} else {
				entryAdminDetails.CreatorCohorts = tempEntryAdminDetails.CreatorCohorts
			}
		} else {
			if updatedCampaignAdminDetails.CreatorCohorts == nil {
				entryAdminDetails.CreatorCohorts = nil
			} else {
				entryAdminDetails.CreatorCohorts = updatedCampaignAdminDetails.CreatorCohorts
			}
		}

		if entryAdminDetails.IsBlacklisted != nil && !*entryAdminDetails.IsBlacklisted {
			entryAdminDetails.BlacklistedBy = nil
			entryAdminDetails.BlacklistedReason = nil
		}
		falseFlag := false
		if entryAdminDetails.WhatsappOptIn == nil {
			entryAdminDetails.WhatsappOptIn = &falseFlag
		}
		if entryAdminDetails.IsBlacklisted == nil {
			entryAdminDetails.IsBlacklisted = &falseFlag
		}
	}
	return entryAdminDetails
}

func MergeUserDetails(cpInput coredomain.CampaignProfileInput, campaign *dao.CampaignProfileEntity) (*string, *coredomain.UserDetails) {
	var tempEntryUserDetails coredomain.UserDetails
	encoded, _ := json.Marshal(cpInput.UserDetails)
	json.Unmarshal(encoded, &tempEntryUserDetails)
	var entryUserDetails *coredomain.UserDetails
	var updatedCampaignUserDetails coredomain.UserDetails
	if campaign.UserDetails != nil {
		updatedCampaignUserDetails = ConvertStringToUserDetails(campaign.UserDetails)
	}
	if cpInput.UserDetails != nil {
		encoded, _ := json.Marshal(cpInput.UserDetails)
		json.Unmarshal(encoded, &entryUserDetails)
	} else {
		encoded, _ := json.Marshal(updatedCampaignUserDetails)
		json.Unmarshal(encoded, &entryUserDetails)
	}
	mergo.Merge(entryUserDetails, updatedCampaignUserDetails, mergo.WithOverride, mergo.WithTransformers(blankTransformer{}))

	entryUserDetails = MergeUserDetailArrays(cpInput, tempEntryUserDetails, entryUserDetails, updatedCampaignUserDetails)
	userDetailsEntity := ToUserDetailsEntity(entryUserDetails)
	if entryUserDetails.CampaignCategoryIds != nil {
		userDetailsEntity.CampaignCategoryIds = convertCampaignCategoryIdsToString(entryUserDetails.CampaignCategoryIds)
	}
	jsonData, _ := json.Marshal(userDetailsEntity)
	jsonDataInString := string(jsonData)
	userDetailsString := &jsonDataInString
	if *userDetailsString == "{}" {
		return nil, entryUserDetails
	}
	return userDetailsString, entryUserDetails
}

func MergeUserDetailArrays(cpInput coredomain.CampaignProfileInput, tempEntryUserDetails coredomain.UserDetails, entryUserDetails *coredomain.UserDetails, updatedCampaignUserDetails coredomain.UserDetails) *coredomain.UserDetails {
	if cpInput.UserDetails != nil {
		if tempEntryUserDetails.Languages != nil {
			if len(*tempEntryUserDetails.Languages) == 0 {
				entryUserDetails.Languages = nil
			} else {
				entryUserDetails.Languages = tempEntryUserDetails.Languages
			}
		} else {
			entryUserDetails.Languages = updatedCampaignUserDetails.Languages
		}
		if tempEntryUserDetails.CampaignCategoryIds != nil {
			if len(*tempEntryUserDetails.CampaignCategoryIds) == 0 {
				entryUserDetails.CampaignCategoryIds = nil
			} else {
				entryUserDetails.CampaignCategoryIds = tempEntryUserDetails.CampaignCategoryIds
			}
		} else {
			entryUserDetails.CampaignCategoryIds = updatedCampaignUserDetails.CampaignCategoryIds
		}
		if tempEntryUserDetails.Location != nil {
			if len(*tempEntryUserDetails.Location) == 0 {
				entryUserDetails.Location = nil
			} else {
				entryUserDetails.Location = tempEntryUserDetails.Location
			}
		} else {
			entryUserDetails.Location = updatedCampaignUserDetails.Location
		}
		falseFlag := false
		if entryUserDetails.InstantGratificationInvited == nil {
			entryUserDetails.InstantGratificationInvited = &falseFlag
		}
		if entryUserDetails.AmazonStoreLinkVerified == nil {
			entryUserDetails.AmazonStoreLinkVerified = &falseFlag
		}
		if entryUserDetails.WhatsappOptIn == nil {
			entryUserDetails.WhatsappOptIn = &falseFlag
		}
		if entryUserDetails.EKYCPending == nil {
			entryUserDetails.EKYCPending = &falseFlag
		}
	}
	return entryUserDetails
}

func languageToCode(language string) string {
	if value, ok := constants.LanguageCodes[language]; ok {
		return value
	} else {
		return language
	}
}

func CodeToLanguage(language string) string {
	if value, ok := constants.LanguageNames[language]; ok {
		return value
	} else {
		return language
	}
}

func SetLocation(adminDetails *coredomain.AdminDetails) *[]string {
	var adminLocations []string
	if adminDetails.City != nil && *adminDetails.City != "" {
		city := strings.Title(*adminDetails.City)
		adminLocations = append(adminLocations, city)
	}
	if adminDetails.State != nil && *adminDetails.State != "" {
		state := strings.Title(*adminDetails.State)
		adminLocations = append(adminLocations, state)
	}
	if adminDetails.Country != nil && *adminDetails.Country != "" {
		country := strings.Title(*adminDetails.Country)
		adminLocations = append(adminLocations, country)
	}
	if len(adminLocations) == 0 {
		return nil
	}
	return &adminLocations
}

func TransformLanguagesToCode(languages *[]string) *[]string {
	if languages == nil {
		return nil
	}
	newLanguages := []string{}
	for _, language := range *languages {
		newLanguages = append(newLanguages, languageToCode(language))
	}
	return &newLanguages
}

func TransformCodeToLanguages(languages *[]string) *[]string {
	if languages == nil {
		return nil
	}
	newLanguages := []string{}
	if languages != nil {
		for _, language := range *languages {
			newLanguages = append(newLanguages, CodeToLanguage(language))
		}
	}
	return &newLanguages
}

func CreateInstaYtCampaignProfileResult(campaignEntity *dao.CampaignProfileEntity, adminDetails *coredomain.AdminDetails, userDetails *coredomain.UserDetails, profile *coredomain.Profile) coredomain.CampaignProfileEntry {
	var adminLanguages *[]string
	var userLanguages *[]string
	var adminDetailsResult *coredomain.AdminDetails
	var userDetailsResult *coredomain.UserDetails
	if adminDetails != nil {
		if adminDetails.Languages != nil {
			adminLanguages = TransformCodeToLanguages(adminDetails.Languages)
		}
		adminDetailsResult = ToAdminDetailsEntry(adminDetails, adminLanguages)
	}
	if userDetails != nil {
		if userDetails.Languages != nil {
			userLanguages = TransformCodeToLanguages(userDetails.Languages)
		}
		userDetailsResult = ToUserDetailsEntry(userDetails, userLanguages)
	}

	return coredomain.CampaignProfileEntry{
		Id:                campaignEntity.ID,
		Platform:          profile.Platform,
		PlatformAccountId: campaignEntity.PlatformAccountId,
		Handle:            profile.Handle,
		AccountId:         campaignEntity.GCCUserAccountID,
		OnGCC:             campaignEntity.OnGCC,
		OnGCCApp:          campaignEntity.OnGCCApp,
		UserDetails:       userDetailsResult,
		AdminDetails:      adminDetailsResult,
		// Socials:           profile.LinkedSocials,
		UpdatedBy: campaignEntity.UpdatedBy,
		HasEmail:  campaignEntity.HasEmail,
		HasPhone:  campaignEntity.HasPhone,
	}
}

func ToAdminDetailsEntry(adminDetails *coredomain.AdminDetails, adminLanguages *[]string) *coredomain.AdminDetails {
	return &coredomain.AdminDetails{
		Name:                adminDetails.Name,
		Email:               adminDetails.Email,
		Phone:               adminDetails.Phone,
		SecondaryPhone:      adminDetails.SecondaryPhone,
		Gender:              adminDetails.Gender,
		Languages:           adminLanguages,
		City:                adminDetails.City,
		State:               adminDetails.State,
		Country:             adminDetails.Country,
		Dob:                 adminDetails.Dob,
		Bio:                 adminDetails.Bio,
		CampaignCategoryIds: adminDetails.CampaignCategoryIds,
		CreatorPrograms:     adminDetails.CreatorPrograms,
		IsBlacklisted:       adminDetails.IsBlacklisted,
		BlacklistedBy:       adminDetails.BlacklistedBy,
		BlacklistedReason:   adminDetails.BlacklistedReason,
		WhatsappOptIn:       adminDetails.WhatsappOptIn,
		CountryCode:         adminDetails.CountryCode,
		CreatorCohorts:      adminDetails.CreatorCohorts,
		Location:            adminDetails.Location,
	}
}

func ToUserDetailsEntry(userDetails *coredomain.UserDetails, userLanguages *[]string) *coredomain.UserDetails {
	return &coredomain.UserDetails{
		Email:                       userDetails.Email,
		Phone:                       userDetails.Phone,
		SecondaryPhone:              userDetails.SecondaryPhone,
		Name:                        userDetails.Name,
		Dob:                         userDetails.Dob,
		Gender:                      userDetails.Gender,
		WhatsappOptIn:               userDetails.WhatsappOptIn,
		CampaignCategoryIds:         userDetails.CampaignCategoryIds,
		Languages:                   userLanguages,
		Location:                    userDetails.Location,
		NotificationToken:           userDetails.NotificationToken,
		Bio:                         userDetails.Bio,
		MemberId:                    userDetails.MemberId,
		ReferenceCode:               userDetails.ReferenceCode,
		EKYCPending:                 userDetails.EKYCPending,
		WebEngageUserId:             userDetails.WebEngageUserId,
		InstantGratificationInvited: userDetails.InstantGratificationInvited,
		AmazonStoreLinkVerified:     userDetails.AmazonStoreLinkVerified,
	}
}

func convertStringToCampaignCategoryIds(creatorStructs coredomain.CreatorStructs) *[]int64 {
	var detailsIntArray []int64
	if creatorStructs.CampaignCategoryIds != nil {
		for _, id := range *creatorStructs.CampaignCategoryIds {
			num := helpers.ParseToInt64(id)
			if num != nil {
				detailsIntArray = append(detailsIntArray, *num)
			}
		}
		return &detailsIntArray
	}
	return nil
}

func CapitalizeUserDetailsLocation(locations *[]string) *[]string {
	var locs []string
	for _, location := range *locations {
		value := strings.Title(location)
		locs = append(locs, value)
	}
	return &locs
}

func CapitalizeAdminDetailsCityStateCountry(adminDetails coredomain.AdminDetails) *coredomain.AdminDetails {
	if adminDetails.City != nil && *adminDetails.City != "" {
		city := strings.Title(*adminDetails.City)
		adminDetails.City = &city
	}
	if adminDetails.State != nil && *adminDetails.State != "" {
		state := strings.Title(*adminDetails.State)
		adminDetails.State = &state
	}
	if adminDetails.Country != nil && *adminDetails.Country != "" {
		country := strings.Title(*adminDetails.Country)
		adminDetails.Country = &country
	}
	return &adminDetails
}

func FillSocialDetailsInFullRefreshResultForIA(cpEntry *coredomain.CampaignProfileEntry, apiData *beatservice.CompleteProfileResponse) *coredomain.CampaignProfileEntry {
	socialDetails := cpEntry.SocialDetails
	var thumbnail *string
	var engagementRate, imageEngagementRate, reelsEngagementRate float64
	var followers, following, storyReach, imageReach, reelsReach, avgReelsPlay30d, uploads *int64
	var avgReach, avgLikes, imageAvgLikes, reelsAvgLikes, avgComments, imageAvgComments, reelsAvgComments *float64
	var recentPosts []coredomain.SocialProfilePostsEntry
	var isPrivate *bool
	var updatedAt, igId string
	thumbnail = socialDetails.PlatformThumbnail
	if apiData.CompleteProfile.ProfilePicUrl != "" {
		thumbnail = &apiData.CompleteProfile.ProfilePicUrl
	}
	following = socialDetails.Following
	if apiData.CompleteProfile.Following != 0 {
		following = &apiData.CompleteProfile.Following
	}
	followers = socialDetails.Followers
	if apiData.CompleteProfile.Followers != 0 {
		followers = &apiData.CompleteProfile.Followers
	}
	engagementRate = socialDetails.EngagementRatePercenatge
	if apiData.CompleteProfile.PostMetrics.All.AverageEngagementRate != nil {
		engagementRate = *apiData.CompleteProfile.PostMetrics.All.AverageEngagementRate
	}
	imageEngagementRate = socialDetails.ImageEngagementRatePercenatge
	if apiData.CompleteProfile.PostMetrics.Static.AverageEngagementRate != nil {
		imageEngagementRate = *apiData.CompleteProfile.PostMetrics.Static.AverageEngagementRate
	}
	reelsEngagementRate = socialDetails.ReelsEngagementRatePercenatge
	if apiData.CompleteProfile.PostMetrics.Reels.AverageEngagementRate != nil {
		reelsEngagementRate = *apiData.CompleteProfile.PostMetrics.Reels.AverageEngagementRate
	}
	avgLikes = socialDetails.AvgLikes
	if apiData.CompleteProfile.PostMetrics.All.AverageLikes != nil {
		avgLikes = apiData.CompleteProfile.PostMetrics.All.AverageLikes
	}
	imageAvgLikes = socialDetails.ImageAvgLikes
	if apiData.CompleteProfile.PostMetrics.Static.AverageLikes != nil {
		imageAvgLikes = apiData.CompleteProfile.PostMetrics.Static.AverageLikes
	}
	reelsAvgLikes = socialDetails.ReelsAvgLikes
	if apiData.CompleteProfile.PostMetrics.Reels.AverageLikes != nil {
		reelsAvgLikes = apiData.CompleteProfile.PostMetrics.Reels.AverageLikes
	}
	avgComments = socialDetails.AvgComments
	if apiData.CompleteProfile.PostMetrics.All.AverageComments != nil {
		avgComments = apiData.CompleteProfile.PostMetrics.All.AverageComments
	}
	imageAvgComments = socialDetails.ImageAvgComments
	if apiData.CompleteProfile.PostMetrics.Static.AverageComments != nil {
		imageAvgComments = apiData.CompleteProfile.PostMetrics.Static.AverageComments
	}
	reelsAvgComments = socialDetails.ReelsAvgComments
	if apiData.CompleteProfile.PostMetrics.Reels.AverageComments != nil {
		reelsAvgComments = apiData.CompleteProfile.PostMetrics.Reels.AverageComments
	}
	avgReach = socialDetails.AvgReach
	if apiData.CompleteProfile.PostMetrics.All.AverageReach != nil {
		avgReach = apiData.CompleteProfile.PostMetrics.All.AverageReach
	}
	storyReach = socialDetails.StoryReach
	if apiData.CompleteProfile.PostMetrics.Story.AverageReach != nil {
		temp := int64(*apiData.CompleteProfile.PostMetrics.Story.AverageReach)
		storyReach = &temp
	}
	imageReach = socialDetails.ImageReach
	if apiData.CompleteProfile.PostMetrics.Static.AverageReach != nil {
		temp := int64(*apiData.CompleteProfile.PostMetrics.Static.AverageReach)
		imageReach = &temp
	}
	reelsReach = socialDetails.ReelsReach
	if apiData.CompleteProfile.PostMetrics.Reels.AverageReach != nil {
		temp := int64(*apiData.CompleteProfile.PostMetrics.Reels.AverageReach)
		reelsReach = &temp
	}

	if socialDetails.AvgReelsPlayCount != nil {
		avgReelsPlay30dint := int64(*socialDetails.AvgReelsPlayCount)
		avgReelsPlay30d = &avgReelsPlay30dint
	}
	if apiData.CompleteProfile.PostMetrics.Reels.AveragePlays != nil {
		temp := int64(*apiData.CompleteProfile.PostMetrics.Reels.AveragePlays)
		avgReelsPlay30d = &temp
	}
	if apiData.CompleteProfile.IsPrivate != nil {
		isPrivate = apiData.CompleteProfile.IsPrivate
	}
	uploads = socialDetails.Uploads
	if apiData.CompleteProfile.Uploads != nil {
		uploads = apiData.CompleteProfile.Uploads
	}
	updatedAt = socialDetails.UpdatedAt
	if apiData.CompleteProfile.UpdatedAt != "" {
		updatedAt = apiData.CompleteProfile.UpdatedAt
	}
	recentPosts = TransformDataToPosts(apiData.CompleteProfile.RecentPosts, socialDetails.Platform, *socialDetails.Handle)
	igId = *socialDetails.IgId
	if apiData.CompleteProfile.ProfileId != "" {
		igId = apiData.CompleteProfile.ProfileId
	}
	cpEntry.SocialDetails = &coredomain.SocialAccount{
		IgId:                          &igId,
		CampaignProfileId:             socialDetails.CampaignProfileId,
		Code:                          socialDetails.Code,
		ProfileCode:                   socialDetails.ProfileCode,
		GCCProfileCode:                socialDetails.GCCProfileCode,
		Name:                          socialDetails.Name,
		Platform:                      socialDetails.Platform,
		Handle:                        socialDetails.Handle,
		Username:                      socialDetails.Username,
		PlatformThumbnail:             thumbnail,
		EngagementRatePercenatge:      engagementRate,
		Followers:                     followers,
		Following:                     following,
		AvgLikes:                      avgLikes,
		AvgComments:                   avgComments,
		Uploads:                       uploads,
		AvgVideoViews30d:              socialDetails.AvgVideoViews30d,
		AvgReelsPlay30d:               avgReelsPlay30d,
		AvgReach:                      avgReach,
		StoryReach:                    storyReach,
		ImageReach:                    imageReach,
		ReelsReach:                    reelsReach,
		Views:                         socialDetails.Views,
		IsVerified:                    socialDetails.IsVerified,
		Source:                        socialDetails.Source,
		AccountType:                   socialDetails.AccountType,
		ProfileLink:                   socialDetails.ProfileLink,
		ImageEngagementRatePercenatge: imageEngagementRate,
		ReelsEngagementRatePercenatge: reelsEngagementRate,
		ImageAvgLikes:                 imageAvgLikes,
		ReelsAvgLikes:                 reelsAvgLikes,
		ImageAvgComments:              imageAvgComments,
		ReelsAvgComments:              reelsAvgComments,
		RecentPosts:                   recentPosts,
		UpdatedAt:                     updatedAt,
		AudienceCityAvailable:         socialDetails.AudienceCityAvailable,
		AudienceGenderAvailable:       socialDetails.AudienceGenderAvailable,
		IsPrivate:                     isPrivate,
	}
	return cpEntry
}
func TransformDataToPosts(recentPosts []beatservice.Posts, platform string, handle string) []coredomain.SocialProfilePostsEntry {
	var socialProfilePostsEntry coredomain.SocialProfilePostsEntry
	var posts []coredomain.SocialProfilePostsEntry
	for _, post := range recentPosts {
		var publishTime time.Time
		var publishTimestamp int64
		publishTime, err := time.Parse("2006-01-02T15:04:05", post.PublishTime)
		if err != nil {
			log.Println(err)
		}
		publishTimestamp = publishTime.Unix()
		socialProfilePostsEntry = coredomain.SocialProfilePostsEntry{
			Platform:       platform,
			Handle:         handle,
			PostId:         post.PostId,
			PostLink:       post.PostUrl,
			PostType:       post.PostType,
			Thumbnail:      post.Dimensions.ThumbnailUrl,
			LikesCount:     post.Metrics.Likes,
			CommentsCount:  post.Dimensions.Comments,
			ViewsCount:     post.Metrics.Reach,
			EngagementRate: post.Metrics.EngagementRate,
			PublishedAt:    publishTimestamp,
			Caption:        post.Dimensions.Caption,
			PlayCount:      post.Metrics.PlayCount,
		}
		posts = append(posts, socialProfilePostsEntry)
	}
	return posts
}
func FillSocialDetailsInFullRefreshResultForYA(result *coredomain.CampaignProfileEntry, apiData *beatservice.CompleteProfileResponse) *coredomain.CampaignProfileEntry {
	var thumbnail *string
	var views, followers, uploads *int64
	socialDetails := result.SocialDetails
	thumbnail = socialDetails.PlatformThumbnail
	if apiData.CompleteProfile.YoutubeThumbnail != "" {
		thumbnail = &apiData.CompleteProfile.YoutubeThumbnail
	}
	views = socialDetails.Views
	if apiData.CompleteProfile.YoutubeViews != 0 {
		views = &apiData.CompleteProfile.YoutubeViews
	}
	followers = socialDetails.Followers
	if apiData.CompleteProfile.YoutubeSubscribers != 0 {
		followers = &apiData.CompleteProfile.YoutubeSubscribers
	}
	uploads = socialDetails.Uploads
	if apiData.CompleteProfile.YoutubeUploads != 0 {
		uploads = &apiData.CompleteProfile.YoutubeUploads
	}
	result.SocialDetails = &coredomain.SocialAccount{
		CampaignProfileId:             socialDetails.CampaignProfileId,
		Code:                          socialDetails.Code,
		ProfileCode:                   socialDetails.ProfileCode,
		GCCProfileCode:                socialDetails.GCCProfileCode,
		Name:                          socialDetails.Name,
		Platform:                      socialDetails.Platform,
		Handle:                        socialDetails.Handle,
		Username:                      socialDetails.Username,
		PlatformThumbnail:             thumbnail,
		Followers:                     followers,
		Uploads:                       uploads,
		Views:                         views,
		Source:                        socialDetails.Source,
		AvgReach:                      socialDetails.AvgReach,
		EngagementRatePercenatge:      socialDetails.EngagementRatePercenatge,
		Following:                     socialDetails.Following,
		AvgLikes:                      socialDetails.AvgLikes,
		AvgComments:                   socialDetails.AvgComments,
		AvgVideoViews30d:              socialDetails.AvgVideoViews30d,
		AvgReelsPlay30d:               socialDetails.AvgReelsPlay30d,
		StoryReach:                    socialDetails.StoryReach,
		ImageReach:                    socialDetails.ImageReach,
		ReelsReach:                    socialDetails.ReelsReach,
		IsVerified:                    socialDetails.IsVerified,
		AccountType:                   socialDetails.AccountType,
		ProfileLink:                   socialDetails.ProfileLink,
		ImageEngagementRatePercenatge: socialDetails.ImageEngagementRatePercenatge,
		ReelsEngagementRatePercenatge: socialDetails.ReelsEngagementRatePercenatge,
		ImageAvgLikes:                 socialDetails.ImageAvgLikes,
		ReelsAvgLikes:                 socialDetails.ReelsAvgLikes,
		ImageAvgComments:              socialDetails.ImageAvgComments,
		ReelsAvgComments:              socialDetails.ReelsAvgComments,
		RecentPosts:                   socialDetails.RecentPosts,
		AudienceCityAvailable:         socialDetails.AudienceCityAvailable,
		AudienceGenderAvailable:       socialDetails.AudienceGenderAvailable,
		IsPrivate:                     socialDetails.IsPrivate,
	}
	return result
}
func CreateInstagramEntityUsingBeatResponse(profile *coredomain.Profile, apiData *beatservice.CompleteProfileResponse) discoverydao.InstagramAccountEntity {
	var thumbnail, name, handle, igId, label, bio *string
	var engagementRate float64
	var followers, following, storyReach, imageReach, reelsReach, uploads *int64
	var avgLikes, avgComments, avgReelsPlayCount, avgReach *float64
	var isPrivate *bool
	var id int64

	thumbnail = profile.Thumbnail
	if apiData.CompleteProfile.ProfilePicUrl != "" {
		thumbnail = &apiData.CompleteProfile.ProfilePicUrl
	}
	name = profile.Name
	if apiData.CompleteProfile.FullName != "" {
		name = &apiData.CompleteProfile.FullName
	}
	handle = profile.Handle
	if apiData.CompleteProfile.Handle != "" {
		handle = &apiData.CompleteProfile.Handle
	}
	igId = profile.IgId
	if apiData.CompleteProfile.ProfileId != "" {
		igId = &apiData.CompleteProfile.ProfileId
	}
	engagementRate = profile.Metrics.EngagementRatePercenatge / 100
	if apiData.CompleteProfile.PostMetrics.All.AverageEngagementRate != nil {
		engagementRate = (*apiData.CompleteProfile.PostMetrics.All.AverageEngagementRate) / 100
	}
	followers = profile.Metrics.Followers
	label = profile.ProfileType
	if apiData.CompleteProfile.Followers != 0 {
		followers = &apiData.CompleteProfile.Followers
		var instalabel string
		if *followers < 1000 {
			instalabel = "Nano"
		} else if *followers < 75000 {
			instalabel = "Micro"
		} else if *followers < 1000000 {
			instalabel = "Macro"
		} else {
			instalabel = "Mega"
		}
		label = &instalabel
	}
	following = profile.Metrics.Following
	if apiData.CompleteProfile.Following != 0 {
		following = &apiData.CompleteProfile.Following
	}
	storyReach = profile.Metrics.StoryReach
	if apiData.CompleteProfile.PostMetrics.Story.AverageReach != nil {
		temp := int64(*apiData.CompleteProfile.PostMetrics.Story.AverageReach)
		storyReach = &temp
	}
	imageReach = profile.Metrics.ImageReach
	if apiData.CompleteProfile.PostMetrics.Static.AverageReach != nil {
		temp := int64(*apiData.CompleteProfile.PostMetrics.Static.AverageReach)
		imageReach = &temp
	}
	reelsReach = profile.Metrics.ReelsReach
	if apiData.CompleteProfile.PostMetrics.Reels.AverageReach != nil {
		temp := int64(*apiData.CompleteProfile.PostMetrics.Reels.AverageReach)
		reelsReach = &temp
	}
	avgLikes = profile.Metrics.AvgLikes
	if apiData.CompleteProfile.PostMetrics.All.AverageLikes != nil {
		avgLikes = apiData.CompleteProfile.PostMetrics.All.AverageLikes
	}
	avgComments = profile.Metrics.AvgComments
	if apiData.CompleteProfile.PostMetrics.All.AverageComments != nil {
		avgComments = apiData.CompleteProfile.PostMetrics.All.AverageComments
	}
	avgReelsPlayCount = profile.Metrics.AvgReelsPlayCount
	if apiData.CompleteProfile.PostMetrics.Reels.AveragePlays != nil {
		avgReelsPlayCount = apiData.CompleteProfile.PostMetrics.Reels.AveragePlays
	}
	if apiData.CompleteProfile.IsPrivate != nil {
		isPrivate = apiData.CompleteProfile.IsPrivate
	}

	avgReach = profile.Metrics.AvgReach
	if apiData.CompleteProfile.PostMetrics.All.AverageReach != nil {
		avgReach = apiData.CompleteProfile.PostMetrics.All.AverageReach
	}
	bio = profile.Description
	if apiData.CompleteProfile.Biography != "" {
		bio = &apiData.CompleteProfile.Biography
	}

	uploads = profile.Metrics.Uploads
	if apiData.CompleteProfile.Uploads != nil {
		uploads = apiData.CompleteProfile.Uploads
	}

	id, _ = strconv.ParseInt(profile.PlatformCode, 10, 64)
	return discoverydao.InstagramAccountEntity{
		ID:                id,
		ProfileId:         profile.Code,
		Name:              name,
		Handle:            handle,
		IgID:              igId,
		Thumbnail:         thumbnail,
		Bio:               bio,
		Followers:         followers,
		Following:         following,
		EngagementRate:    engagementRate,
		AvgLikes:          avgLikes,
		AvgReach:          avgReach,
		ReelsReach:        reelsReach,
		StoryReach:        storyReach,
		ImageReach:        imageReach,
		AvgComments:       avgComments,
		AvgReelsPlayCount: avgReelsPlayCount,
		Label:             label,
		IsPrivate:         isPrivate,
		PostCount:         uploads,
		AvgViews:          profile.Metrics.AvgViews,
	}
}
func CreateYoutubeEntityUsingBeatResponse(profile *coredomain.Profile, apiData *beatservice.CompleteProfileResponse) discoverydao.YoutubeAccountEntity {
	var thumbnail, channelId, title, label *string
	var followers, uploadsCount, viewsCount *int64
	var id int64
	thumbnail = profile.Thumbnail
	if apiData.CompleteProfile.YoutubeThumbnail != "" {
		thumbnail = &apiData.CompleteProfile.YoutubeThumbnail
	}
	title = profile.Name
	if apiData.CompleteProfile.YoutubeTitle != "" {
		title = &apiData.CompleteProfile.YoutubeTitle
	}
	channelId = profile.Handle
	if apiData.CompleteProfile.YoutubeChannelId != "" {
		channelId = &apiData.CompleteProfile.YoutubeChannelId
	}
	followers = profile.Metrics.Followers
	label = profile.ProfileType
	if apiData.CompleteProfile.YoutubeSubscribers != 0 {
		followers = &apiData.CompleteProfile.YoutubeSubscribers
		var ytlabel string
		if *followers < 1000 {
			ytlabel = "Nano"
		} else if *followers < 75000 {
			ytlabel = "Micro"
		} else if *followers < 1000000 {
			ytlabel = "Macro"
		} else {
			ytlabel = "Mega"
		}
		label = &ytlabel
	}
	uploadsCount = profile.Metrics.Uploads
	if apiData.CompleteProfile.YoutubeUploads != 0 {
		uploadsCount = &apiData.CompleteProfile.YoutubeUploads
	}
	viewsCount = profile.Metrics.Views
	if apiData.CompleteProfile.YoutubeViews != 0 {
		viewsCount = &apiData.CompleteProfile.YoutubeViews
	}
	id, _ = strconv.ParseInt(profile.PlatformCode, 10, 64)
	return discoverydao.YoutubeAccountEntity{
		ID:           id,
		ProfileId:    profile.Code,
		GccProfileId: profile.GccCode,
		Username:     profile.Username,
		Description:  profile.Description,
		VideoReach:   profile.Metrics.VideoReach,
		AvgViews:     profile.Metrics.AvgViews,
		ChannelId:    channelId,
		Title:        title,
		Thumbnail:    thumbnail,
		Label:        label,
		UploadsCount: uploadsCount,
		Followers:    followers,
		ViewsCount:   viewsCount,
	}
}

func SetHasEmailAndHasPhoneFieldsInCampaignEntity(ctx context.Context, campaign *dao.CampaignProfileEntity, platform string, profile *coredomain.Profile) *dao.CampaignProfileEntity {
	entry, _ := ToCampaignEntry(ctx, campaign, profile)
	hasEmail := false
	hasPhone := false
	if entry.HasEmail != nil {
		hasEmail = *entry.HasEmail
	}
	if entry.HasPhone != nil {
		hasPhone = *entry.HasPhone
	}

	if (entry.AdminDetails != nil && entry.AdminDetails.Email != nil && *entry.AdminDetails.Email != "") ||
		(entry.UserDetails != nil && entry.UserDetails.Email != nil && *entry.UserDetails.Email != "") {
		hasEmail = true
	}
	if (entry.AdminDetails != nil && entry.AdminDetails.Phone != nil && *entry.AdminDetails.Phone != "") ||
		(entry.AdminDetails != nil && entry.AdminDetails.SecondaryPhone != nil && *entry.AdminDetails.SecondaryPhone != "") ||
		(entry.UserDetails != nil && entry.UserDetails.Phone != nil && *entry.UserDetails.Phone != "") ||
		(entry.UserDetails != nil && entry.UserDetails.SecondaryPhone != nil && *entry.UserDetails.SecondaryPhone != "") {
		hasPhone = true
	}
	if platform == string(constants.InstagramPlatform) || platform == string(constants.YoutubePlatform) {
		hasEmail = hasEmail || (profile != nil && profile.Email != nil && *profile.Email != "")
		hasPhone = hasPhone || (profile != nil && profile.Phone != nil && *profile.Phone != "")
	}
	campaign.HasEmail = &hasEmail
	campaign.HasPhone = &hasPhone
	return campaign
}

func getJSONString(data interface{}) *string {
	jsonData, _ := json.Marshal(data)
	jsonDataInString := string(jsonData)
	return &jsonDataInString
}
