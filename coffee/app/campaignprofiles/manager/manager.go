package manager

import (
	"coffee/app/campaignprofiles/dao"
	"coffee/app/campaignprofiles/domain"
	discoverydomain "coffee/app/discovery/domain"
	discoverymanager "coffee/app/discovery/manager"
	beatservice "coffee/client/beat"
	"coffee/constants"
	"errors"

	"coffee/core/appcontext"
	coredomain "coffee/core/domain"

	"context"
	"strconv"
	"strings"

	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

type CampaignProfileManager struct {
	dao           *dao.CampaignProfileDao
	searchManager *discoverymanager.SearchManager
}

func CreateManager(ctx context.Context) *CampaignProfileManager {
	campaignDao := dao.CreateCampaignProfileDao(ctx)
	searchManager := discoverymanager.CreateSearchManager(ctx)
	manager := &CampaignProfileManager{
		dao:           campaignDao,
		searchManager: searchManager,
	}
	return manager
}

func (m *CampaignProfileManager) UpsertGccCampaignProfile(ctx context.Context, platform string, cpInput coredomain.CampaignProfileInput) (coredomain.CampaignProfileEntry, error) {
	var adminDetails *coredomain.AdminDetails
	var userDetails *coredomain.UserDetails
	var cpEntry coredomain.CampaignProfileEntry
	var err error
	gccCpEntity, _ := m.dao.FindCPByAccountId(ctx, platform, cpInput.AccountId)
	if cpInput.AccountId == nil {
		return cpEntry, errors.New("accountID missing in request")
	}
	if gccCpEntity == nil {
		var err error
		gccCpEntity = &dao.CampaignProfileEntity{
			Platform:          string(constants.GCCPlatform),
			PlatformAccountId: *cpInput.AccountId,
			Enabled:           true,
		}
		gccCpEntity, err = m.dao.Create(ctx, gccCpEntity)
		if err != nil {
			log.Error("error creating CP account for GCC platform", err)
			return cpEntry, err
		}
	}
	adminDetails, userDetails, *gccCpEntity = domain.MergeUserAdminDetails(cpInput, gccCpEntity)
	gccCpEntity = domain.SetHasEmailAndHasPhoneFieldsInCampaignEntity(ctx, gccCpEntity, string(constants.GCCPlatform), nil)
	gccCpEntity, err = m.dao.Update(ctx, *gccCpEntity, gccCpEntity.ID)
	if err != nil {
		log.Error("Error updating gcc account: ", err)
		return cpEntry, err
	}
	cpEntry = coredomain.CampaignProfileEntry{
		Id:                gccCpEntity.ID,
		Platform:          string(constants.GCCPlatform),
		PlatformAccountId: *cpInput.AccountId,
		AdminDetails:      adminDetails,
		UserDetails:       userDetails,
		UpdatedBy:         cpInput.UpdatedBy,
		AccountId:         cpInput.AccountId,
		OnGCC:             cpInput.OnGCC,
		OnGCCApp:          cpInput.OnGCCApp,
		HasEmail:          gccCpEntity.HasEmail,
		HasPhone:          gccCpEntity.HasPhone,
	}
	return cpEntry, nil
}

func (m *CampaignProfileManager) FindProfileInDiscovery(ctx context.Context, social *coredomain.SocialAccount, linkedSource string) (*coredomain.Profile, error) {

	if social.Platform == string(constants.InstagramPlatform) {
		*social.Handle = strings.ToLower(*social.Handle)
	}
	profile, err := m.searchManager.FindProfileByPlatformHandle(ctx, social.Platform, *social.Handle, linkedSource, false)
	if err != nil {
		return nil, err
	}
	if profile.GccCode == nil {
		log.Error("Profile_id is missing in profile ", err)
		return nil, err
	}
	return profile, nil
}

func (m *CampaignProfileManager) UpsertInstaYtCampaignProfile(ctx context.Context, profile *coredomain.Profile, cpInput coredomain.CampaignProfileInput) (coredomain.CampaignProfileEntry, *coredomain.Profile, *coredomain.AdminDetails, *coredomain.UserDetails, error) {

	var adminDetails *coredomain.AdminDetails
	var userDetails *coredomain.UserDetails
	var cpEntry coredomain.CampaignProfileEntry
	var err error
	var gender, dob *string
	var phone, location, languages []string

	cpEntity, _ := m.dao.FindByPlatformAndPlatformId(ctx, profile.Platform, profile.PlatformCode)

	if cpEntity == nil {
		platformAccountId, _ := strconv.ParseInt(profile.PlatformCode, 10, 64)
		cpEntity = &dao.CampaignProfileEntity{
			Platform:          profile.Platform,
			PlatformAccountId: platformAccountId,
			Enabled:           true,
		}
		cpEntity, err = m.dao.Create(ctx, cpEntity)
		if err != nil {
			log.Error(err)
		}
	}
	adminDetails, userDetails, *cpEntity = domain.MergeUserAdminDetails(cpInput, cpEntity)
	cpEntity = domain.SetHasEmailAndHasPhoneFieldsInCampaignEntity(ctx, cpEntity, profile.Platform, profile)
	if cpInput.OnGCCApp != nil && !*cpInput.OnGCCApp {
		cpEntity.GCCUserAccountID = nil
		flag := false
		cpEntity.OnGCCApp = &flag
	}
	// Can Move This To Separate Function
	if userDetails.Gender != nil {
		gender = userDetails.Gender
	} else if adminDetails.Gender != nil {
		gender = adminDetails.Gender
	} else if profile.Gender != nil {
		gender = profile.Gender
	}
	if userDetails.Phone != nil {
		phone = append(phone, *userDetails.Phone)
	}
	if adminDetails.Phone != nil {
		phone = append(phone, *adminDetails.Phone)
	}
	if userDetails.Dob != nil {
		dob = userDetails.Dob
	} else if adminDetails.Dob != nil {
		dob = adminDetails.Dob
	}

	if userDetails.Location != nil {
		location = append(location, *userDetails.Location...)
	}
	if adminDetails.Location != nil {
		location = append(location, *adminDetails.Location...)
	}
	if profile.LocationList != nil && len(*profile.LocationList) > 0 {
		location = append(location, *profile.LocationList...)
	}
	// Will ADD location List From Profile i.e IA/YA

	if userDetails.Languages != nil {
		languages = append(languages, *userDetails.Languages...)
	}
	if adminDetails.Languages != nil {
		languages = append(languages, *adminDetails.Languages...)
	}
	if profile.Languages != nil && len(*profile.Languages) > 0 {
		languages = append(languages, *profile.Languages...)
	}

	cpEntity.Gender = gender
	cpEntity.Phone = pq.StringArray(phone)
	cpEntity.Dob = dob
	cpEntity.Location = pq.StringArray(location)
	cpEntity.Languages = pq.StringArray(languages)

	cpEntity, err = m.dao.Update(ctx, *cpEntity, cpEntity.ID)
	if err != nil {
		log.Error("Error updating adminDetails in CP account: ", err)
		return cpEntry, nil, adminDetails, userDetails, err
	}

	cpEntry = domain.CreateInstaYtCampaignProfileResult(cpEntity, adminDetails, userDetails, profile)
	return cpEntry, profile, adminDetails, userDetails, nil
}

func (m *CampaignProfileManager) updateSocialLinkage(ctx context.Context, instagramAccount *coredomain.Profile, youtubeAccount *coredomain.Profile, linkedSource string, searchPhraseAdminIA string, keywordsAdminIA []string, searchPhraseAdminYA string, keywordsAdminYA []string) error {

	var err error
	if instagramAccount != nil && youtubeAccount != nil {
		err = m.UnLinkingInstagramAccount(ctx, instagramAccount, linkedSource, keywordsAdminIA, searchPhraseAdminIA)
		if err != nil {
			log.Error("Error while linking-unlinking instagram account: ", err)
			return err
		}
		err = m.UnLinkingYoutubeAccount(ctx, youtubeAccount, linkedSource, keywordsAdminYA, searchPhraseAdminYA)
		if err != nil {
			log.Error("Error while linking-unlinking youtube account: ", err)
			return err
		}
		_, err = m.searchManager.UpdateMultipleFieldsYA(ctx, *instagramAccount.GccCode, instagramAccount.Handle, youtubeAccount.PlatformCode, nil, nil)
		if err != nil {
			log.Error("Error while linking youtube account: ", err)
			return err
		}
		_, err = m.searchManager.UpdateMultipleFieldsIA(ctx, youtubeAccount.Handle, instagramAccount.PlatformCode, nil, nil)
		if err != nil {
			log.Error("Error while linking instagram account: ", err)
			return err
		}
	} else if instagramAccount != nil {
		err = m.UnLinkingInstagramAccount(ctx, instagramAccount, linkedSource, keywordsAdminIA, searchPhraseAdminIA)
		if err != nil {
			log.Error("Error while linking-unlinking instagram account: ", err)
			return err
		}
	} else if youtubeAccount != nil {
		err = m.UnLinkingYoutubeAccount(ctx, youtubeAccount, linkedSource, keywordsAdminYA, searchPhraseAdminYA)
		if err != nil {
			log.Error("Error while linking-unlinking youtube account: ", err)
			return err
		}
	}
	return nil
}

func (m *CampaignProfileManager) GenerateKeywordsAdminAndSearchPhraseAdminForIaYa(ctx context.Context, instagramAccount *coredomain.Profile, youtubeAccount *coredomain.Profile, adminDetails *coredomain.AdminDetails, userDetails *coredomain.UserDetails) (string, []string) {
	if instagramAccount != nil {
		searchPhraseAdminIA, keywordsAdminIA := domain.CreateKeywordsAndSearchPhraseAdmin(adminDetails, userDetails, string(constants.InstagramPlatform), instagramAccount)
		return searchPhraseAdminIA, keywordsAdminIA
	}

	if youtubeAccount != nil {
		searchPhraseAdminYA, keywordsAdminYA := domain.CreateKeywordsAndSearchPhraseAdmin(adminDetails, userDetails, string(constants.YoutubePlatform), youtubeAccount)
		return searchPhraseAdminYA, keywordsAdminYA
	}
	return "", nil
}

func (m *CampaignProfileManager) sanitizeInput(ctx context.Context, cpInput coredomain.CampaignProfileInput) coredomain.CampaignProfileInput {
	if cpInput.AdminDetails != nil {
		cpInput.AdminDetails = domain.CapitalizeAdminDetailsCityStateCountry(*cpInput.AdminDetails)
		cpInput.AdminDetails.Location = domain.SetLocation(cpInput.AdminDetails)
	}
	if cpInput.OnGCCApp != nil && *cpInput.OnGCCApp {
		cpInput.OnGCC = cpInput.OnGCCApp
	}
	if cpInput.UserDetails != nil && cpInput.UserDetails.Location != nil {
		cpInput.UserDetails.Location = domain.CapitalizeUserDetailsLocation(cpInput.UserDetails.Location)
	}
	return cpInput
}

func (m *CampaignProfileManager) UpsertProfileAdminDetails(ctx context.Context, profileCode string, cpInput coredomain.CampaignProfileInput, linkedSource string) ([]coredomain.CampaignProfileEntry, error) {
	var err error
	var cpEntries []coredomain.CampaignProfileEntry
	var adminDetails *coredomain.AdminDetails
	var userDetails *coredomain.UserDetails
	var iaProfile *coredomain.Profile
	var yaProfile *coredomain.Profile
	var searchPhraseAdminIA, searchPhraseAdminYA string
	var keywordsAdminIA, keywordsAdminYA []string

	cpInput, err = m.enrichHandleInSocials(ctx, cpInput, profileCode)
	if err != nil {
		return cpEntries, err
	}
	if cpInput.AccountId != nil {
		cpInput = m.enrichGccPlatformInSocials(ctx, cpInput, profileCode)
	}

	cpInput = m.sanitizeInput(ctx, cpInput)

	for _, social := range cpInput.Socials {
		var cpEntry coredomain.CampaignProfileEntry
		var account *coredomain.Profile
		if social.Platform == string(constants.GCCPlatform) || social.Platform == string(constants.InstagramPlatform) || social.Platform == string(constants.YoutubePlatform) {
			switch social.Platform {
			case string(constants.GCCPlatform):
				cpEntry, err = m.UpsertGccCampaignProfile(ctx, social.Platform, cpInput)
				if err != nil {
					log.Error("Error creating GCC Campaign: ", err)
					return cpEntries, err
				}
				cpEntries = append(cpEntries, cpEntry)

			case string(constants.InstagramPlatform), string(constants.YoutubePlatform):
				profile, err := m.FindProfileInDiscovery(ctx, social, linkedSource)
				if err != nil {
					return cpEntries, err
				}
				cpEntry, account, adminDetails, userDetails, err = m.UpsertInstaYtCampaignProfile(ctx, profile, cpInput)
				if err != nil {
					log.Error("Error creating instagram-youtube campaign: ", err)
					return cpEntries, err
				}
				if social.Platform == string(constants.InstagramPlatform) {
					iaProfile = account
					searchPhraseAdminIA, keywordsAdminIA = m.GenerateKeywordsAdminAndSearchPhraseAdminForIaYa(ctx, iaProfile, nil, adminDetails, userDetails)
				} else {
					yaProfile = account
					searchPhraseAdminYA, keywordsAdminYA = m.GenerateKeywordsAdminAndSearchPhraseAdminForIaYa(ctx, nil, yaProfile, adminDetails, userDetails)
				}
				cpEntries = append(cpEntries, cpEntry)
			}
		}
	}
	err = m.updateSocialLinkage(ctx, iaProfile, yaProfile, linkedSource, searchPhraseAdminIA, keywordsAdminIA, searchPhraseAdminYA, keywordsAdminYA)
	if err != nil {
		log.Error("Error while linking-unlinking instagram and youtube accounts: ", err)
		return cpEntries, err
	}
	cpEntries = m.fillResultWithLinkedSocials(ctx, cpInput, cpEntries, linkedSource)
	return cpEntries, nil
}

func (m *CampaignProfileManager) FindCampaignProfileByPlatformHandle(ctx context.Context, platform string, platformHandle string, linkedSource string, fullRefresh bool) (*coredomain.CampaignProfileEntry, error) {
	var cpEntry *coredomain.CampaignProfileEntry
	if linkedSource == "" {
		linkedSource = string(constants.GCCPlatform)
	}
	if platform == string(constants.InstagramPlatform) {
		platformHandle = strings.ToLower(platformHandle)
	}
	if platform != string(constants.GCCPlatform) && platform != string(constants.InstagramPlatform) && platform != string(constants.YoutubePlatform) {
		err := errors.New("platform is not GCC or INSTAGRAM or YOUTUBE")
		return cpEntry, err
	}
	if platform == string(constants.YoutubePlatform) && strings.HasPrefix(platformHandle, "@") {
		beatClient := beatservice.New(context.TODO())
		ytProfile, err := beatClient.FindYoutubeChannelIdByHandle(platformHandle)
		if err != nil {
			return nil, err
		}
		if ytProfile.YoutubeChannel.ChannelId != "" {
			platformHandle = ytProfile.YoutubeChannel.ChannelId
		} else {
			return nil, errors.New("unable to find channel id for handle")
		}
	}
	profile, err := m.searchManager.FindProfileByPlatformHandle(ctx, platform, platformHandle, linkedSource, false)
	if err != nil {
		return cpEntry, err
	}
	campaignEntity, err := m.dao.FindByPlatformAndPlatformId(ctx, platform, profile.PlatformCode)
	if err != nil {
		return nil, err
	}
	if err == nil && campaignEntity == nil {
		campaignEntity, err = m.CreateCampaignProfileForNewHandles(ctx, profile)
		if err != nil {
			return cpEntry, err
		}
	}
	cpEntry, _ = domain.ToCampaignEntry(ctx, campaignEntity, profile)
	if fullRefresh {
		cpEntry, err = m.UpdateCampaignProfileResultForFullRefresh(ctx, cpEntry, platform, platformHandle, profile)
		if err != nil {
			log.Debug("error getting profile from beat", err)
		}

	}
	return cpEntry, nil
}

func (m *CampaignProfileManager) UnLinkingInstagramAccount(ctx context.Context, instagramAccount *coredomain.Profile, linkedSource string, keywordsAdmin []string, searchPhraseAdminIA string) error {
	var err error
	_, err = m.searchManager.UpdateMultipleFieldsIA(ctx, new(string), instagramAccount.PlatformCode, &keywordsAdmin, &searchPhraseAdminIA)
	if err != nil {
		log.Error("Error while updating channelId in IA: ", err)
		return err
	}
	ytAccouts, _ := m.searchManager.FindYoutubeAccountsByProfileId(ctx, *instagramAccount.GccCode, linkedSource) // check this call.
	if ytAccouts != nil {
		for _, ytAccount := range *ytAccouts {
			newGccProfileID := "YA_" + ytAccount.PlatformCode
			_, err = m.searchManager.UpdateMultipleFieldsYA(ctx, newGccProfileID, new(string), ytAccount.PlatformCode, nil, nil)
			if err != nil {
				log.Error("Error while unlinking youtube account: ", err)
				return err
			}
		}
	}
	return nil
}

func (m *CampaignProfileManager) UnLinkingYoutubeAccount(ctx context.Context, youtubeAccount *coredomain.Profile, linkedSource string, keywordsAdminYA []string, searchPhraseAdminYA string) error {
	var err error
	currentYoutubeGccProfileId := youtubeAccount.GccCode

	newYoutubeGccProfileId := "YA_" + youtubeAccount.PlatformCode
	_, err = m.searchManager.UpdateMultipleFieldsYA(ctx, newYoutubeGccProfileId, new(string), youtubeAccount.PlatformCode, &keywordsAdminYA, &searchPhraseAdminYA)
	if err != nil {
		log.Error("Error while unlinking youtube account: ", err)
		return err
	}

	profile, _ := m.searchManager.FindByProfileId(ctx, *currentYoutubeGccProfileId, linkedSource)
	if profile != nil && profile.Platform == string(constants.InstagramPlatform) && profile.PlatformCode != "" {
		_, err = m.searchManager.UpdateMultipleFieldsIA(ctx, new(string), profile.PlatformCode, nil, nil)
		if err != nil {
			log.Error("Error while updating channelId in IA: ", err)
			return err
		}
	}
	return nil
}

/*
on CM, in pop-up edit only profile_id will be sent in request. this method is used to fill socials when profile_id starts with IA, YA
*/
func (m *CampaignProfileManager) enrichSocialsFromProfileCode(ctx context.Context, profileCode string, entry coredomain.CampaignProfileInput) (coredomain.CampaignProfileInput, error) {
	isIA := strings.HasPrefix(profileCode, "IA")
	isYA := strings.HasPrefix(profileCode, "YA")

	var socials []*coredomain.SocialAccount
	if isIA {
		// search in both IA and YA tables.
		iaAccount, err := m.searchManager.FindByProfileId(ctx, profileCode, string(constants.GCCPlatform))
		if err != nil {
			log.Error("Error during filling socials. IA not found for given profileID: ", err)
			return entry, err
		}
		linkedSocials := iaAccount.LinkedSocials
		for _, social := range linkedSocials {
			if social.Platform == string(constants.InstagramPlatform) && social.Source == string(constants.GCCPlatform) {
				instaCompactProfile := &coredomain.SocialAccount{
					Platform: social.Platform,
					Handle:   social.Handle,
				}
				socials = append(socials, instaCompactProfile)
			}
			if social.Platform == string(constants.YoutubePlatform) && social.Source == string(constants.GCCPlatform) {
				youtubeCompactProfile := &coredomain.SocialAccount{
					Platform: social.Platform,
					Handle:   social.Handle,
				}
				socials = append(socials, youtubeCompactProfile)
			}
		}
	} else if isYA {
		// search in only YA table.
		yaAccount, err := m.searchManager.FindByProfileId(ctx, profileCode, string(constants.GCCPlatform))
		if err != nil {
			log.Error("Error during filling socials. YA not found for given profileID: ", err)
			return entry, err
		}

		linkedSocials := yaAccount.LinkedSocials
		for _, social := range linkedSocials {
			if social.Platform == string(constants.YoutubePlatform) && social.Source == string(constants.GCCPlatform) {
				youtubeCompactProfile := &coredomain.SocialAccount{
					Platform: social.Platform,
					Handle:   social.Handle,
				}
				socials = append(socials, youtubeCompactProfile)
			}
		}
	}
	entry.Socials = socials
	return entry, nil
}

func (m *CampaignProfileManager) FindById(ctx context.Context, id int64, socials bool) (*coredomain.CampaignProfileEntry, error) {
	var result *coredomain.CampaignProfileEntry

	campaignEntity, err := m.dao.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	if !socials || (campaignEntity.Platform != string(constants.InstagramPlatform) && campaignEntity.Platform != string(constants.YoutubePlatform)) {
		result, _ = domain.ToCampaignEntry(ctx, campaignEntity, nil)
		return result, nil
	}
	profile, err := m.searchManager.FindByPlatformProfileId(ctx, campaignEntity.Platform, campaignEntity.PlatformAccountId, string(constants.GCCPlatform), false)
	if err != nil {
		log.Error("IA/YA is not found for cp", campaignEntity.ID, err)
		return nil, err
	}
	result, _ = domain.ToCampaignEntry(ctx, campaignEntity, profile)
	return result, nil
}

/*
go-gateway requires socials in compactProfile format,
so filling linkedSocials in results. this will be used by go-gateway to fill socials data.
*/
func (m *CampaignProfileManager) fillResultWithLinkedSocials(ctx context.Context, entry coredomain.CampaignProfileInput, results []coredomain.CampaignProfileEntry, linkedSource string) []coredomain.CampaignProfileEntry {

	for i := range results {
		if results[i].Platform == string(constants.InstagramPlatform) || results[i].Platform == string(constants.YoutubePlatform) {
			profile, _ := m.searchManager.FindProfileByPlatformHandle(ctx, results[i].Platform, *results[i].Handle, linkedSource, false)
			if profile != nil {
				for j := range profile.LinkedSocials {
					if profile.LinkedSocials[j].Platform == results[i].Platform && profile.LinkedSocials[j].Source == string(constants.GCCPlatform) && strconv.FormatInt(profile.LinkedSocials[j].Code, 10) == profile.PlatformCode {
						results[i].SocialDetails = profile.LinkedSocials[j]
					}
				}
			}
		}
	}
	return results
}

func (m *CampaignProfileManager) UpdateById(ctx context.Context, id int64, cpInput coredomain.CampaignProfileInput, linkedSource string) ([]coredomain.CampaignProfileEntry, error) {
	var result coredomain.CampaignProfileEntry
	var results []coredomain.CampaignProfileEntry
	var updatedCampaignEntity dao.CampaignProfileEntity
	var adminDetails *coredomain.AdminDetails
	var userDetails *coredomain.UserDetails
	var profile *coredomain.Profile
	var gender, dob *string
	var phone, location, languages []string

	campaignEntity, err := m.dao.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	cpInput = m.sanitizeInput(ctx, cpInput)
	adminDetails, userDetails, updatedCampaignEntity = domain.MergeUserAdminDetails(cpInput, campaignEntity)
	if campaignEntity.Platform == string(constants.InstagramPlatform) || campaignEntity.Platform == string(constants.YoutubePlatform) {
		profile, err = m.searchManager.FindByPlatformProfileId(ctx, campaignEntity.Platform, campaignEntity.PlatformAccountId, linkedSource, false)
		if err != nil {
			log.Error("profile not found error ", err)
			return nil, err
		}
	}
	updatedCampaignEntity = *domain.SetHasEmailAndHasPhoneFieldsInCampaignEntity(ctx, &updatedCampaignEntity, campaignEntity.Platform, profile)
	// Can Move This To Separate Function
	if userDetails.Gender != nil {
		gender = userDetails.Gender
	} else if adminDetails.Gender != nil {
		gender = adminDetails.Gender
	} else if profile.Gender != nil {
		gender = profile.Gender
	}
	if userDetails.Phone != nil {
		phone = append(phone, *userDetails.Phone)
	}
	if adminDetails.Phone != nil {
		phone = append(phone, *adminDetails.Phone)
	}
	if userDetails.Dob != nil {
		dob = userDetails.Dob
	} else if adminDetails.Dob != nil {
		dob = adminDetails.Dob
	}

	if userDetails.Location != nil {
		location = append(location, *userDetails.Location...)
	}
	if adminDetails.Location != nil {
		location = append(location, *adminDetails.Location...)
	}
	if profile.LocationList != nil && len(*profile.LocationList) > 0 {
		location = append(location, *profile.LocationList...)
	}
	// Will ADD location List From Profile i.e IA/YA

	if userDetails.Languages != nil {
		languages = append(languages, *userDetails.Languages...)
	}
	if adminDetails.Languages != nil {
		languages = append(languages, *adminDetails.Languages...)
	}
	if profile.Languages != nil && len(*profile.Languages) > 0 {
		languages = append(languages, *profile.Languages...)
	}

	campaignEntity.Gender = gender
	campaignEntity.Phone = pq.StringArray(phone)
	campaignEntity.Dob = dob
	campaignEntity.Location = pq.StringArray(location)
	campaignEntity.Languages = pq.StringArray(languages)

	campaignEntity, err = m.dao.Update(ctx, updatedCampaignEntity, campaignEntity.ID)
	if err != nil {
		log.Error("Error updating adminDetails in CP account: ", err)
		return nil, err
	}

	if profile != nil {
		if profile.Platform == string(constants.InstagramPlatform) {
			searchPhraseAdminIA, keywordsAdminIA := m.GenerateKeywordsAdminAndSearchPhraseAdminForIaYa(ctx, profile, nil, adminDetails, userDetails)
			_, err := m.searchManager.UpdateMultipleFieldsIA(ctx, nil, profile.PlatformCode, &keywordsAdminIA, &searchPhraseAdminIA)
			if err != nil {
				log.Error("Error while updating keywords admin and search phrase admin in IA: ", err)
				return results, err
			}
		} else if profile != nil && profile.Platform == string(constants.YoutubePlatform) {
			searchPhraseAdminYA, keywordsAdminYA := m.GenerateKeywordsAdminAndSearchPhraseAdminForIaYa(ctx, nil, profile, adminDetails, userDetails)
			_, err := m.searchManager.UpdateMultipleFieldsYA(ctx, *profile.GccCode, nil, profile.PlatformCode, &keywordsAdminYA, &searchPhraseAdminYA)
			if err != nil {
				log.Error("Error while updating keywords admin and search phrase admin in YA: ", err)
				return results, err
			}
		}
		result = domain.CreateInstaYtCampaignProfileResult(campaignEntity, adminDetails, userDetails, profile)
		result = *domain.EnrichCampaignProfileFromSocialProfile(&result, profile)
	} else {
		campaignEntry, _ := domain.ToCampaignEntry(ctx, campaignEntity, nil)
		result = *campaignEntry
	}
	results = append(results, result)
	return results, nil
}

func (m *CampaignProfileManager) Search(ctx context.Context, query coredomain.SearchQuery, sortBy string, sortDir string, page int, size int, fetchSocials bool) ([]coredomain.CampaignProfileEntry, int64, error) {
	var entries []coredomain.CampaignProfileEntry
	filterMap := domain.CampaignProfileSearchFilterMap
	var newSearchQuery coredomain.SearchQuery
	var profileCollectionItemJoin, profileCollectionJoin bool

	for _, filter := range query.Filters {
		var newFilter coredomain.SearchFilter

		field := filter.Field
		value, ok := filterMap[field]
		if ok {
			filter.Field = value
		} else {
			filter.Field = "campaign_profiles." + field
		}

		newFilter.Field = filter.Field
		newFilter.FilterType = filter.FilterType
		newFilter.Value = filter.Value
		newSearchQuery.Filters = append(newSearchQuery.Filters, newFilter)
		if strings.Contains(value, "profile_collection_item.") {
			profileCollectionItemJoin = true
		}

		if strings.Contains(value, "profile_collection.") {
			profileCollectionItemJoin = true
			profileCollectionJoin = true
		}
	}
	var joinClauseStruct []coredomain.JoinClauses

	if profileCollectionItemJoin {
		joinClauseStruct = append(joinClauseStruct, coredomain.JoinClauses{
			JoinClause: "JOIN profile_collection_item ON campaign_profiles.platform_account_id = profile_collection_item.platform_account_code AND profile_collection_item.platform = campaign_profiles.platform",
		})
	}
	if profileCollectionJoin {
		joinClauseStruct = append(joinClauseStruct, coredomain.JoinClauses{
			JoinClause: "JOIN profile_collection ON profile_collection_item.profile_collection_id = profile_collection.id AND profile_collection_item.platform = campaign_profiles.platform",
		})
	}
	entities, filteredCount, err := m.dao.SearchJoins(ctx, newSearchQuery, sortBy, sortDir, page, size, joinClauseStruct)
	if err != nil {
		return entries, 0, err
	}
	var instaAccountCode, ytAccountCode []int64
	var profiles *[]coredomain.Profile
	if !fetchSocials {
		for i := range entities {
			entry, _ := domain.ToCampaignEntry(ctx, &entities[i], nil)
			entries = append(entries, *entry)
		}
		return entries, filteredCount, err
	}
	for i := range entities {
		if entities[i].Platform == string(constants.InstagramPlatform) {
			instaAccountCode = append(instaAccountCode, entities[i].PlatformAccountId)
		}
		if entities[i].Platform == string(constants.YoutubePlatform) {
			ytAccountCode = append(ytAccountCode, entities[i].PlatformAccountId)
		}

	}

	instaMap := make(map[string]*coredomain.Profile)
	ytMap := make(map[string]*coredomain.Profile)
	if len(instaAccountCode) > 0 {
		profiles, err = m.searchManager.FindByPlatformProfileIds(ctx, string(constants.InstagramPlatform), instaAccountCode, string(constants.GCCPlatform))
		if err == nil {
			for _, profile := range *profiles {
				instaMap[profile.PlatformCode] = &profile
			}
		}
	}

	if len(ytAccountCode) > 0 {
		profiles, err = m.searchManager.FindByPlatformProfileIds(ctx, string(constants.YoutubePlatform), ytAccountCode, string(constants.GCCPlatform))
		if err == nil {
			for _, profile := range *profiles {
				ytMap[profile.PlatformCode] = &profile
			}
		}
	}

	for i := range entities {
		var profile *coredomain.Profile
		if entities[i].Platform == string(constants.InstagramPlatform) {
			platformAccountId := strconv.FormatInt(entities[i].PlatformAccountId, 10)
			profile = instaMap[platformAccountId]
		}
		if entities[i].Platform == string(constants.YoutubePlatform) {
			platformAccountId := strconv.FormatInt(entities[i].PlatformAccountId, 10)
			profile = ytMap[platformAccountId]
		}
		entry, _ := domain.ToCampaignEntry(ctx, &entities[i], profile)
		entries = append(entries, *entry)
	}

	return entries, filteredCount, err
}

func (m *CampaignProfileManager) FindCampaignProfileByPlatformAccountCode(ctx context.Context, platform string, platformAccountCode int64, linkedSource string) (*coredomain.CampaignProfileEntry, error) {
	var cpEntry *coredomain.CampaignProfileEntry
	log.Debug("FindCampaignProfileByPlatformAccountCode Manager")
	if platform != string(constants.GCCPlatform) && platform != string(constants.InstagramPlatform) && platform != string(constants.YoutubePlatform) {
		err := errors.New("platform is not GCC or INSTAGRAM or YOUTUBE")
		return cpEntry, err
	}
	profile, err := m.searchManager.FindByPlatformProfileId(ctx, platform, platformAccountCode, linkedSource, false)
	if err != nil {
		return cpEntry, err
	}
	if profile == nil {
		return cpEntry, errors.New("invalid social account")
	}
	campaignEntity, err := m.dao.FindByPlatformAndPlatformId(ctx, platform, profile.PlatformCode)
	if err != nil {
		return nil, err
	}
	if err == nil && campaignEntity == nil {
		// entity doesnt exist, create in db now
		campaignEntity, err = m.CreateCampaignProfileForNewHandles(ctx, profile)
		if err != nil {
			return cpEntry, err
		}
	}
	campaignEntry, _ := domain.ToCampaignEntry(ctx, campaignEntity, profile)
	cpEntry = domain.EnrichCampaignProfileFromSocialProfile(campaignEntry, profile)
	return cpEntry, nil
}

func (m *CampaignProfileManager) enrichHandleInSocials(ctx context.Context, cpInput coredomain.CampaignProfileInput, profileCode string) (coredomain.CampaignProfileInput, error) {
	var err error
	instagramLinked := false
	ytLinked := false
	for i := range cpInput.Socials {
		social := cpInput.Socials[i]
		if social.CampaignProfileId != nil && *social.CampaignProfileId > 0 {
			cp, err := m.FindById(ctx, *social.CampaignProfileId, true)
			if err != nil {
				return coredomain.CampaignProfileInput{}, err
			}
			cpInput.Socials[i].Platform = cp.Platform
			if cp.Platform == string(constants.InstagramPlatform) {
				if instagramLinked {
					return coredomain.CampaignProfileInput{}, errors.New("cannot link multiple instagram profiles")
				}
				instagramLinked = true
			}
			if cp.Platform == string(constants.YoutubePlatform) {
				if ytLinked {
					return coredomain.CampaignProfileInput{}, errors.New("cannot link multiple youtube profiles")
				}
				ytLinked = true
			}
			cpInput.Socials[i].Handle = cp.Handle
		}
	}
	if len(cpInput.Socials) == 0 && profileCode != "" {
		cpInput, err = m.enrichSocialsFromProfileCode(ctx, profileCode, cpInput)
		if err != nil {
			log.Error("Error while filling socials for given profile_id: ", err)
			return cpInput, err
		}
	}
	// var gccSocials *coredomain.SocialAccount
	// if cpInput.AccountId != nil {
	// 	gccSocials = &coredomain.SocialAccount{
	// 		Platform: string(constants.GCCPlatform),
	// 	}
	// 	cpInput.Socials = append(cpInput.Socials, gccSocials)
	// }
	return cpInput, nil
}

func (m *CampaignProfileManager) enrichGccPlatformInSocials(ctx context.Context, cpInput coredomain.CampaignProfileInput, profileCode string) coredomain.CampaignProfileInput {
	gccSocials := &coredomain.SocialAccount{
		Platform: string(constants.GCCPlatform),
	}
	cpInput.Socials = append(cpInput.Socials, gccSocials)
	return cpInput
}

func (m *CampaignProfileManager) UpdateCampaignProfileResultForFullRefresh(ctx context.Context, campaignProfileEntry *coredomain.CampaignProfileEntry, platform string, handle string, profile *coredomain.Profile) (*coredomain.CampaignProfileEntry, error) {
	profileFromBeat, err := m.FetchProfileFromBeat(ctx, platform, handle)
	if err != nil {
		return campaignProfileEntry, nil
	}
	if platform == string(constants.InstagramPlatform) {
		campaignProfileEntry = domain.FillSocialDetailsInFullRefreshResultForIA(campaignProfileEntry, profileFromBeat)
		_ = m.UpdateInstagramAccountFromBeatResponse(ctx, profile, profileFromBeat)
	} else if platform == string(constants.YoutubePlatform) {
		campaignProfileEntry = domain.FillSocialDetailsInFullRefreshResultForYA(campaignProfileEntry, profileFromBeat)
		_ = m.UpdateYoutubeAccountFromBeatResponse(ctx, profile, profileFromBeat)
	}
	return campaignProfileEntry, nil
}

func (m *CampaignProfileManager) FetchProfileFromBeat(ctx context.Context, platform string, handle string) (*beatservice.CompleteProfileResponse, error) {
	beatClient := beatservice.New(ctx)
	apiData, err := beatClient.FetchCompleteProfileByHandle(platform, handle)
	if err != nil {
		log.Error("error fetching data from beat api: ", err)
		return nil, err
	}
	if apiData.Status.Type != "SUCCESS" {
		return nil, errors.New("profile not found")
	}
	return apiData, nil
}

func (m *CampaignProfileManager) UpdateInstagramAccountFromBeatResponse(ctx context.Context, profile *coredomain.Profile, beatProfile *beatservice.CompleteProfileResponse) error {
	var err error
	instagramEntity := domain.CreateInstagramEntityUsingBeatResponse(profile, beatProfile)
	_, err = m.searchManager.UpdateInstagramForFullRefresh(ctx, profile.PlatformCode, instagramEntity)
	if err != nil {
		log.Error("error while updating instagram account in full refresh")
		return err
	}
	return nil
}

func (m *CampaignProfileManager) UpdateYoutubeAccountFromBeatResponse(ctx context.Context, profile *coredomain.Profile, beatProfile *beatservice.CompleteProfileResponse) error {
	var err error
	YoutubeEntity := domain.CreateYoutubeEntityUsingBeatResponse(profile, beatProfile)
	_, err = m.searchManager.UpdateYoutubeForFullRefresh(ctx, profile.PlatformCode, YoutubeEntity)
	if err != nil {
		log.Error("error while updating instagram account in full refresh")
		return err
	}
	return nil
}

func (m *CampaignProfileManager) CreateCampaignProfileForNewHandles(ctx context.Context, profile *coredomain.Profile) (*dao.CampaignProfileEntity, error) {
	var err error
	falseFlag := false
	updatedBy := "SYSTEM"
	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	if appCtx.AccountId != nil {
		updatedBy = strconv.Itoa(int(*appCtx.AccountId))
	}
	cpInput := &coredomain.CampaignProfileInput{
		OnGCC:     &falseFlag,
		OnGCCApp:  &falseFlag,
		UpdatedBy: &updatedBy,
	}
	cpEntity, _ := domain.ToCampaignEntity(cpInput)
	cpEntity.PlatformAccountId, _ = strconv.ParseInt(profile.PlatformCode, 10, 64)
	cpEntity.Platform = profile.Platform
	cpEntity = domain.SetHasEmailAndHasPhoneFieldsInCampaignEntity(ctx, cpEntity, profile.Platform, profile)
	cpEntity, err = m.dao.Create(ctx, cpEntity)
	if err != nil {
		return nil, err
	}

	if profile.Platform == string(constants.InstagramPlatform) {
		searchPhraseAdmin, keywordsAdmin := m.GenerateKeywordsAdminAndSearchPhraseAdminForIaYa(ctx, profile, nil, nil, nil)
		_, err = m.searchManager.UpdateMultipleFieldsIA(ctx, new(string), profile.PlatformCode, &keywordsAdmin, &searchPhraseAdmin)
		if err != nil {
			log.Error("Error while updating keywordsAdmin and searchPhraseAdmin in IA: ", err)
			return nil, err
		}
	}
	if profile.Platform == string(constants.YoutubePlatform) {
		searchPhraseAdmin, keywordsAdmin := m.GenerateKeywordsAdminAndSearchPhraseAdminForIaYa(ctx, nil, profile, nil, nil)
		_, err = m.searchManager.UpdateMultipleFieldsYA(ctx, *profile.GccCode, new(string), profile.PlatformCode, &keywordsAdmin, &searchPhraseAdmin)
		if err != nil {
			log.Error("Error while updating keywordsAdmin and searchPhraseAdmin in YA: ", err)
			return nil, err
		}
	}
	return cpEntity, nil
}

func (m *CampaignProfileManager) RefreshSocialProfileById(ctx context.Context, id int64) (*coredomain.CampaignProfileEntry, error) {
	campaigns, err := m.FindById(ctx, id, true)
	if err != nil {
		return nil, err
	}
	switch campaigns.Platform {
	case "GCC":
		return campaigns, nil
	case "INSTAGRAM":
		platformCode := campaigns.SocialDetails.Code
		platformCodeStr := strconv.FormatInt(platformCode, 10)
		beatClient := beatservice.New(ctx)
		apiData, err := beatClient.FindSocialProfileByProfileId(campaigns.Platform, *campaigns.SocialDetails.IgId)
		if err != nil {
			log.Error("error fetching data from beat api: ", err)
			return nil, err
		}
		if apiData.Status.Type == "ERROR" {
			return nil, errors.New("profile not found")
		}
		profile, err := m.searchManager.FindByPlatformProfileId(ctx, campaigns.Platform, platformCode, "GCC", false)
		if err != nil {
			return nil, errors.New("profile not found")
		}
		instagramEntity := domain.CreateInstagramEntityUsingBeatResponse(profile, apiData)
		m.searchManager.UpdateInstagramForFullRefresh(ctx, platformCodeStr, instagramEntity)
	case "YOUTUBE":
		platformCode := campaigns.SocialDetails.Code
		platformCodeStr := strconv.FormatInt(platformCode, 10)
		beatClient := beatservice.New(ctx)
		apiData, err := beatClient.FindSocialProfileByProfileId(campaigns.Platform, *campaigns.SocialDetails.Handle)
		if err != nil {
			log.Error("error fetching data from beat api: ", err)
			return nil, err
		}
		if apiData.Status.Type == "ERROR" {
			return nil, errors.New("profile not found")
		}
		profile, err := m.searchManager.FindByPlatformProfileId(ctx, campaigns.Platform, platformCode, "GCC", false)
		if err != nil {
			return nil, errors.New("profile not found")
		}
		YoutubeEntity := domain.CreateYoutubeEntityUsingBeatResponse(profile, apiData)
		m.searchManager.UpdateYoutubeForFullRefresh(ctx, platformCodeStr, YoutubeEntity)
	}
	campaignProfile, _ := m.FindById(ctx, id, true)
	return campaignProfile, nil
}

func (m *CampaignProfileManager) CreatorInsightsById(ctx context.Context, id int64, token string, platform string, userId string) (*coredomain.CampaignProfileEntry, error) {
	if platform != string(constants.InstagramPlatform) {
		err := errors.New("platform is not supported")
		return nil, err
	}
	cpEntry, err := m.FindById(ctx, id, true)
	if err != nil {
		return nil, err
	}
	profile, err := m.searchManager.FindByPlatformProfileId(ctx, "INSTAGRAM", cpEntry.PlatformAccountId, "SAAS", true) // confirm source and campaignProfileJoin
	if err != nil {
		return nil, err
	}
	beatClient := beatservice.New(ctx)
	apiData, err := beatClient.FindCreatorInsightsByHandle(*profile.Handle, token, userId)
	if err != nil {
		return nil, err
	}
	if apiData.Status.Type == "ERROR" {
		err := "error from beat: " + apiData.Status.Message
		return nil, errors.New(err)
	}
	cpEntry = domain.FillSocialDetailsInFullRefreshResultForIA(cpEntry, apiData)
	cpEntry.SocialDetails.FollowersGrowth7d = profile.Metrics.FollowersGrowth7d
	cpEntry.SocialDetails.FollowersGrowth30d = profile.Metrics.FollowersGrowth30d
	cpEntry.SocialDetails.UpdatedAt = strconv.FormatInt(profile.UpdatedAt, 10)
	return cpEntry, nil
}

func (m *CampaignProfileManager) CreatorAudienceInsightsDataById(ctx context.Context, id int64, token string, platform string, userId string) (*domain.CampaignProfileAudienceInsightsEntry, error) {
	if platform != string(constants.InstagramPlatform) {
		err := errors.New("platform is not supported")
		return nil, err
	}
	cpEntry, err := m.FindById(ctx, id, false)
	if err != nil {
		return nil, err
	}
	profile, err := m.searchManager.FindByPlatformProfileId(ctx, "INSTAGRAM", cpEntry.PlatformAccountId, "SAAS", true) // confirm source and campaignProfileJoin
	if err != nil {
		return nil, err
	}
	beatClient := beatservice.New(ctx)
	apiData, err := beatClient.FindCreatorAudienceInsightsByHandle(*profile.Handle, token, userId)
	if err != nil {
		return nil, err
	}
	if apiData.Status.Type == "ERROR" {
		err := "error from beat: " + apiData.Status.Message
		return nil, errors.New(err)
	}
	cityWiseData, normalizedData, audienceGenderSplit := convertCityAndAudienceDataToPercentage(apiData)
	return &domain.CampaignProfileAudienceInsightsEntry{
		Id:                cpEntry.Id,
		Platform:          "INSTAGRAM",
		PlatformAccountId: cpEntry.PlatformAccountId,
		Handle:            profile.Handle,
		UpdatedAt:         strconv.FormatInt(profile.UpdatedAt, 10),
		AudienceData: discoverydomain.SocialProfileAudienceInfoEntry{
			CityWiseAudienceLocation: cityWiseData,
			AudienceGenderSplit:      audienceGenderSplit,
			AudienceAgeGenderSplit: discoverydomain.AudienceAgeGender{
				Male:   normalizedData["male"],
				Female: normalizedData["female"],
			},
		},
	}, nil
}

func convertCityAndAudienceDataToPercentage(apiData *beatservice.CreatorInsight) (map[string]*float64, map[string]map[string]*float64, map[string]*float64) {
	cityWiseData := make(map[string]*float64)
	cityTotal := int64(0)
	otherCitiesTotal := int64(0)
	for _, count := range apiData.Profile.AudienceData.City {
		cityTotal += count
	}

	for city, count := range apiData.Profile.AudienceData.City {
		if len(cityWiseData) < 5 {
			value := (float64(count) / float64(cityTotal)) * 100
			cityWiseData[city] = &value
		} else {
			otherCitiesTotal += count
		}
	}

	if otherCitiesTotal > 0 {
		value := (float64(otherCitiesTotal) / float64(cityTotal)) * 100
		cityWiseData["others"] = &value
	}

	normalizedData := make(map[string]map[string]*float64)

	total := make(map[string]int)
	totalAge := 0
	for gender, ageGroup := range apiData.Profile.AudienceData.GenderAge {
		ageGroupMap, _ := ageGroup.(map[string]interface{})
		for _, value := range ageGroupMap {
			age, _ := value.(float64)
			total[gender] += int(age)
		}
	}
	if maleAge, ok := total["male"]; ok {
		totalAge += maleAge
	}
	if femaleAge, ok := total["female"]; ok {
		totalAge += femaleAge
	}

	audienceGenderSplit := make(map[string]*float64)
	if maleAge, ok := total["male"]; ok {
		male_per := (float64(maleAge) / float64(totalAge)) * 100
		audienceGenderSplit["male_per"] = &male_per
	}
	if femaleAge, ok := total["female"]; ok {
		female_per := (float64(femaleAge) / float64(totalAge)) * 100
		audienceGenderSplit["female_per"] = &female_per
	}

	for gender, ageGroup := range apiData.Profile.AudienceData.GenderAge {
		if gender == "male" || gender == "female" {
			normalizedAgeGroup := make(map[string]*float64)
			ageGroupMap, _ := ageGroup.(map[string]interface{})
			for age, value := range ageGroupMap {
				ageValue, _ := value.(float64)
				value := (float64(ageValue) / float64(totalAge)) * 100
				normalizedAgeGroup[age] = &value
			}
			normalizedData[gender] = normalizedAgeGroup
		}
	}
	return cityWiseData, normalizedData, audienceGenderSplit
}
