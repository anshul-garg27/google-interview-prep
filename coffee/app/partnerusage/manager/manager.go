package manager

import (
	"coffee/app/partnerusage/dao"
	"coffee/app/partnerusage/domain"
	partnerclient "coffee/client/partner"
	"coffee/publishers"

	"github.com/google/uuid"

	"coffee/constants"
	"coffee/core/appcontext"
	"coffee/core/persistence/redis"
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	discoverymanager "coffee/app/discovery/manager"

	log "github.com/sirupsen/logrus"
)

type PartnerUsageManager struct {
	partnerUsageDao        *dao.PartnerUsageDao
	partnerProfileTrackDao *dao.PartnerProfileTrackDao
	activityTrackerDao     *dao.ActivityTrackerDao
	searchManager          *discoverymanager.SearchManager
}

func CreatePartnerUsageManager(ctx context.Context) *PartnerUsageManager {
	partnerUsageDao := dao.CreatePartnerUsageDao(ctx)
	partnerProfileTrackDao := dao.CreatePartnerProfileTrackDao(ctx)
	activityTrackerDao := dao.CreateActivityTrackerDao(ctx)
	searchManager := discoverymanager.CreateSearchManager(ctx)

	manager := &PartnerUsageManager{
		partnerUsageDao:        partnerUsageDao,
		partnerProfileTrackDao: partnerProfileTrackDao,
		activityTrackerDao:     activityTrackerDao,
		searchManager:          searchManager,
	}
	return manager
}

func (m *PartnerUsageManager) FindByPartnerId(ctx context.Context, partnerId int64) (*domain.PartnerInput, error) {
	var partnerContract *domain.PartnerInput
	var err error
	redisKey := "partnercontract-" + strconv.FormatInt(partnerId, 10)
	partnerContract, err = GetPartnerContractFromCache(redisKey)

	if err != nil {
		PartnerClient := partnerclient.New(ctx)
		partnerResponse, err := PartnerClient.GetPartnerContract(partnerId)
		if len(partnerResponse.Partners) < 1 {
			err = errors.New("contract missing")
		}
		if err != nil {
			return nil, err
		}
		partnerContract = &partnerResponse.Partners[0]
		expiration := 12 * time.Hour
		setPartnerContractInRedis(redisKey, partnerContract, expiration)
	}
	partnerUsage, err := m.EnrichPartnerWithConsumption(ctx, partnerId, partnerContract)
	if err != nil {
		return nil, err
	}
	return partnerUsage, nil
}

func (m *PartnerUsageManager) UpdatePartnerUsage(ctx context.Context, partnerId int64, partnerContract *domain.PartnerInput) error {
	if len(partnerContract.Contracts) == 0 {
		return errors.New("contract missing")
	}
	contract := partnerContract.Contracts[0]
	if len(contract.Configurations) == 0 {
		return errors.New("configurations missing")
	}

	startDate := time.Unix(0, contract.StartTime*int64(time.Millisecond))
	endDate := time.Unix(0, contract.EndTime*int64(time.Millisecond))
	for _, configuration := range contract.Configurations {
		namespace := configuration.Namespace
		if len(configuration.Configuration) > 0 {
			for _, config := range configuration.Configuration {
				if config.Key == "USAGE" && config.Type == "COUNTER" {
					entity := makeEntity(partnerId, namespace, startDate, endDate, &config)
					err := m.upSertPartnerUsage(ctx, entity)
					if err != nil {
						return err
					}
				}
			}
		}
	}
	partnerCacheKey := "partnercontract-" + strconv.FormatInt(partnerId, 10)
	err := redis.DeleteValueByKey(partnerCacheKey)
	if err != nil {
		log.Error(err)
	}
	return nil
}

func (m *PartnerUsageManager) upSertPartnerUsage(ctx context.Context, entity *dao.PartnerUsageEntity) error {
	partnerExist, err := m.partnerUsageDao.CheckUsage(ctx, entity.PartnerId, entity.Namespace)

	if err != nil && err.Error() != "record not found" {
		return err
	}
	if err != nil && err.Error() == "record not found" {
		m.partnerUsageDao.Create(ctx, entity)
	} else {
		entity.Consumed = partnerExist.Consumed
		m.partnerUsageDao.Update(ctx, partnerExist.ID, entity)
	}
	return nil
}

func (m *PartnerUsageManager) EnrichPartnerWithConsumption(ctx context.Context, partnerId int64, partnerUsage *domain.PartnerInput) (*domain.PartnerInput, error) {
	partnerConsuption, err := m.partnerUsageDao.FindConsuptionByPartnerId(ctx, partnerId)
	if err != nil {
		return nil, err
	}
	consumptionMap := make(map[string]dao.PartnerUsageEntity)
	for _, v := range *partnerConsuption {
		consumptionMap[v.Namespace] = v
	}
	if len(partnerUsage.Contracts) == 0 {
		return nil, errors.New("contract missing")
	}

	contract := partnerUsage.Contracts[0]
	if len(contract.Configurations) == 0 {
		return nil, errors.New("configurations missing")
	}
	contractConfigurations := contract.Configurations
	for moduleKey := range contractConfigurations {
		moduleName := contractConfigurations[moduleKey].Namespace
		EnrichModuleUsage(contractConfigurations, moduleKey, consumptionMap, moduleName)
	}
	return partnerUsage, nil
}

func EnrichModuleUsage(contractConfigurations []domain.ContractNamespace, moduleKey int, consumptionMap map[string]dao.PartnerUsageEntity, moduleName string) {
	if len(contractConfigurations[moduleKey].Configuration) > 0 {
		moduleConfig := contractConfigurations[moduleKey].Configuration
		for moduleConfigKey := range moduleConfig {
			EnrichUsageCounter(moduleConfig, moduleConfigKey, consumptionMap, moduleName)
		}
	}
}

func EnrichUsageCounter(moduleConfig []domain.ContractConfiguration, moduleConfigKey int, consumptionMap map[string]dao.PartnerUsageEntity, moduleName string) {
	if moduleConfig[moduleConfigKey].Key == "USAGE" && moduleConfig[moduleConfigKey].Type == "COUNTER" {
		if val, ok := consumptionMap[moduleName]; ok {
			consumeDetails := val
			moduleConfig[moduleConfigKey].Consumed = *consumeDetails.Consumed
		}
	}
}

func (m *PartnerUsageManager) CheckUsage(ctx context.Context, partnerId int64, moduleName constants.Module) (*dao.PartnerUsageEntity, error) {
	module := string(moduleName)
	partnerUsageEntity, err := m.partnerUsageDao.CheckUsage(ctx, partnerId, module)
	if err != nil {
		return nil, errors.New("partner contract not found")
	}

	return partnerUsageEntity, nil
}

func (m *PartnerUsageManager) IncrementPartnerModuleUsage(ctx context.Context, partnerUsageEntity dao.PartnerUsageEntity, incrementCount int64) (*domain.PartnerUsageEntry, error) {
	consumed := *partnerUsageEntity.Consumed + incrementCount
	partnerUsageEntity.Consumed = &consumed
	updatedEntity, err := m.partnerUsageDao.Update(ctx, partnerUsageEntity.ID, &partnerUsageEntity)
	if err != nil {
		return nil, err
	}
	return ToPartnerUsageEntry(updatedEntity)
}

func (m *PartnerUsageManager) CreateProfileTracker(ctx context.Context, partnerId int64, key string) (*domain.PartnerProfileTrackEntry, error) {
	var err error
	profileTrackEEntity := &dao.PartnerProfileTrackEntity{
		PartnerId: partnerId,
		Key:       key,
	}
	profileTrackEEntity, err = m.partnerProfileTrackDao.Create(ctx, profileTrackEEntity)
	if err != nil {
		return nil, err
	}
	trackEntry, _ := ToProfileTrackEntry(profileTrackEEntity)
	return trackEntry, nil
}
func (m *PartnerUsageManager) FindProfilePageByKey(ctx context.Context, partnerId int64, key string) (*domain.PartnerProfileTrackEntry, error) {
	trackEntity, err := m.partnerProfileTrackDao.FindByProfileKey(ctx, partnerId, key)
	if err != nil {
		return nil, err
	}
	trackEntry, _ := ToProfileTrackEntry(trackEntity)
	return trackEntry, nil

}

func (m *PartnerUsageManager) DeletePartnerContract(ctx context.Context, partnerId int64) error {
	err := m.partnerUsageDao.DeletePartnerUsage(ctx, partnerId)
	if err != nil {
		return err
	}
	partnerCacheKey := "partnercontract-" + strconv.FormatInt(partnerId, 10)
	return redis.DeleteValueByKey(partnerCacheKey)
}

func GetPartnerContractFromCache(key string) (*domain.PartnerInput, error) {
	partnerusage := &domain.PartnerInput{}
	partnerContractString, err := redis.GetValueByKey(key)
	if err != nil {
		return nil, err
	}
	partnerContractJson := strings.ReplaceAll(partnerContractString, "\\\"", "\"")
	err = json.Unmarshal([]byte(partnerContractJson), partnerusage)
	if err != nil {
		return nil, err
	}
	return partnerusage, nil
}

func setPartnerContractInRedis(key string, partnerusage *domain.PartnerInput, expiryTime time.Duration) {
	redisClient := redis.Redis()
	marshaledData, err := json.Marshal(partnerusage)
	if err != nil {
		log.Debug("error in partner contract marshling for cache", partnerusage)
	}
	err = redisClient.Set(key, string(marshaledData), expiryTime).Err()
	if err != nil {
		log.Debug("error in setting partner contract in cache", marshaledData)
	}
}

func (m *PartnerUsageManager) LogPaidActivity(ctx context.Context, moduleName constants.Module, activityName constants.Activity, incrementCount int64, partnerId int64, accountId int64, activityId string, activityMeta string, profilePageTrackKey string, planType string) (*dao.AcitivityTrackerEntity, error) {
	return m.logActivity(ctx, moduleName, activityName, incrementCount, partnerId, accountId, activityId, activityMeta, true, profilePageTrackKey, planType)
}

func (m *PartnerUsageManager) LogFreeActivity(ctx context.Context, moduleName constants.Module, activityName constants.Activity, incrementCount int64, partnerId int64, accountId int64, activityId string, activityMeta string, profilePageTrackKey string, planType string) (*dao.AcitivityTrackerEntity, error) {
	return m.logActivity(ctx, moduleName, activityName, incrementCount, partnerId, accountId, activityId, activityMeta, false, profilePageTrackKey, planType)
}

func (m *PartnerUsageManager) logActivity(ctx context.Context, moduleName constants.Module, activityName constants.Activity, incrementCount int64, partnerId int64, accountId int64, activityId string, activityMeta string, checkUsage bool, profilePageTrackKey string, planType string) (*dao.AcitivityTrackerEntity, error) {
	// Check Usage Call if Consumed Return Error

	if planType == string(constants.PaidPlan) || planType == string(constants.SaasPlan) {
		if constants.PartnerLimitModules[string(moduleName)] {
			partnerUsageEntity, err := m.CheckUsage(ctx, partnerId, moduleName)
			if checkUsage {
				if partnerUsageEntity != nil && *partnerUsageEntity.Consumed >= partnerUsageEntity.Limit {
					err = errors.New("limit consumed")
				}
				if err != nil {
					return nil, err
				}
			}
			_, err = m.IncrementPartnerModuleUsage(ctx, *partnerUsageEntity, incrementCount)
			if err != nil {
				return nil, err
			}
			if moduleName == constants.ProfilePage && profilePageTrackKey != "" {
				_, err := m.CreateProfileTracker(ctx, partnerId, profilePageTrackKey)
				if err != nil {
					return nil, err
				}
			}
		}
	}
	var entity *dao.AcitivityTrackerEntity
	var err error
	if constants.TrackedModule[string(activityName)] {
		entity, err = m.UpsertActivity(ctx, moduleName, partnerId, accountId, activityId, activityMeta)
		if err != nil {
			return nil, err
		}
	}
	var activityIdForLogEvent *int64
	if entity != nil {
		activityIdForLogEvent = &entity.Id
	}
	err = m.createPartnerActivityLogEvent(ctx, moduleName, activityName, incrementCount, partnerId, accountId, activityIdForLogEvent, activityMeta)
	if err != nil {
		log.Error("error while publishing activity to queue")
	}
	return entity, nil
}

func (m *PartnerUsageManager) UpsertActivity(ctx context.Context, moduleName constants.Module, partnerId int64, accountId int64, activityId string, activityMeta string) (*dao.AcitivityTrackerEntity, error) {
	enabled := true
	var entity *dao.AcitivityTrackerEntity
	var err error
	if activityId != "" {
		id, _ := strconv.ParseInt(activityId, 10, 64)

		// Fill The Parameters To Be Updated
		activity := dao.AcitivityTrackerEntity{
			Id:        id,
			PartnerId: partnerId,
			AccountId: accountId,
			Type:      string(moduleName),
			Meta:      &activityMeta,
			Enabled:   &enabled,
		}
		entity, err = m.activityTrackerDao.Update(ctx, id, &activity)
		if err != nil {
			return nil, err
		}
	} else {
		activity := &dao.AcitivityTrackerEntity{
			PartnerId: partnerId,
			AccountId: accountId,
			Type:      string(moduleName),
			Meta:      &activityMeta,
			Enabled:   &enabled,
		}
		entity, err = m.activityTrackerDao.Create(ctx, activity)
		if err != nil {
			return nil, err
		}
	}
	return entity, nil
}

func (m *PartnerUsageManager) createPartnerActivityLogEvent(ctx context.Context, moduleName constants.Module, activityName constants.Activity, incrementCount int64, partnerId int64, accountId int64, activityId *int64, activityMeta string) error {
	messagePayload := m.createPayloadForClickhouse(moduleName, activityName, incrementCount, partnerId, accountId, activityId, activityMeta)
	jsonBytes, _ := json.Marshal(messagePayload)
	topic := "coffee.dx___activity_tracker_rk"
	err := publishers.PublishMessage(jsonBytes, topic)
	return err
}

func (m *PartnerUsageManager) createPayloadForClickhouse(moduleName constants.Module, activityName constants.Activity, incrementCount int64, partnerId int64, accountId int64, activityId *int64, activityMeta string) map[string]interface{} {
	messagePayload := make(map[string]interface{})
	messagePayload["moduleName"] = string(moduleName)
	messagePayload["activityName"] = string(activityName)
	messagePayload["incrementCount"] = incrementCount
	messagePayload["partnerId"] = partnerId
	messagePayload["accountId"] = accountId
	messagePayload["activityId"] = activityId
	messagePayload["activityMeta"] = activityMeta
	uuidObj := uuid.New()
	uuidStr := uuidObj.String()
	messagePayload["eventId"] = uuidStr
	currentTime := time.Now()
	messagePayload["eventTimestamp"] = currentTime.Format("2006-01-02 15:04:05")
	return messagePayload
}
func (m *PartnerUsageManager) UnlockProfile(ctx context.Context, partnerID *int64, profileCode string, platform string) error {
	profile, err := m.searchManager.FindByProfileId(ctx, profileCode, string(constants.SaasLinkedProfile))
	if err != nil {
		return err
	}
	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	for _, socialAccount := range profile.LinkedSocials {
		if socialAccount.Platform == platform {
			platformProfileId := socialAccount.Code
			profilePageTrackKey := "profile-" + platform + "-" + strconv.FormatInt(socialAccount.Code, 10)
			trackEntry, err := m.FindProfilePageByKey(ctx, *partnerID, profilePageTrackKey)
			activityMeta := makeDiscoveryAudienceActivityMeta(platform, platformProfileId)
			if err != nil && trackEntry == nil {
				_, err := m.LogPaidActivity(ctx, constants.ProfilePage, constants.ProfilePageActivity, int64(1), *partnerID, *appCtx.AccountId, "", activityMeta, profilePageTrackKey, string(*appCtx.PlanType))
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func makeDiscoveryAudienceActivityMeta(platform string, platformProfileId int64) string {
	var requestMetaStr string
	var requestMeta = make(map[string]interface{})
	requestMeta["platform"] = platform
	requestMeta["platformProfileId"] = platformProfileId
	requestMetaJson, err := json.Marshal(requestMeta)
	if err != nil {
		log.Error(err)
		return requestMetaStr
	}
	requestMetaStr = string(requestMetaJson)
	return requestMetaStr
}
