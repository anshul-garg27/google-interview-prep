package manager

import (
	"coffee/app/partnerusage/dao"
	"coffee/app/partnerusage/domain"
	"context"
	"time"
)

func ToPartnerUsageEntity(entry *domain.PartnerUsageEntry) (*dao.PartnerUsageEntity, error) {
	return &dao.PartnerUsageEntity{
		PartnerId:   entry.PartnerId,
		Namespace:   entry.Namespace,
		Key:         entry.Key,
		Limit:       entry.Limit,
		Consumed:    entry.Consumed,
		NextResetOn: entry.NextResetOn,
		StartDate:   entry.StartDate,
		EndDate:     entry.EndDate,
		Enabled:     entry.Enabled,
	}, nil
}

func ToPartnerUsageEntry(entity *dao.PartnerUsageEntity) (*domain.PartnerUsageEntry, error) {
	return &domain.PartnerUsageEntry{
		ID:          entity.ID,
		PartnerId:   entity.PartnerId,
		Namespace:   entity.Namespace,
		Key:         entity.Key,
		Limit:       entity.Limit,
		Consumed:    entity.Consumed,
		NextResetOn: entity.NextResetOn,
		Frequency:   entity.Frequency,
		StartDate:   entity.StartDate,
		EndDate:     entity.EndDate,
		Enabled:     entity.Enabled,
	}, nil
}

// TO DO make this like ToEntity
func makeEntity(id int64, namespace string, startDate time.Time, endDate time.Time, configuration *domain.ContractConfiguration) *dao.PartnerUsageEntity {
	enabled := true
	nextResetOn := getNextResetDate(startDate, configuration.Frequency)

	return &dao.PartnerUsageEntity{
		PartnerId:   id,
		Namespace:   namespace,
		Key:         configuration.Key,
		Limit:       configuration.Value,
		Consumed:    new(int64),
		NextResetOn: nextResetOn,
		Frequency:   configuration.Frequency,
		StartDate:   startDate,
		EndDate:     endDate,
		Enabled:     &enabled,
	}
}

func getNextResetDate(startDateObj time.Time, frequency string) time.Time {
	currentDate := time.Now()
	var nextResetDate time.Time
	if startDateObj.Before(currentDate) {
		if frequency == "daily" {
			nextResetDate = currentDate.AddDate(0, 0, 1)
		} else if frequency == "weekly" {
			duration := currentDate.Sub(startDateObj)
			daysDifference := duration / (24 * time.Hour)
			daysAdded := int(7 - daysDifference%7)
			nextResetDate = currentDate.AddDate(0, 0, daysAdded)
		} else if frequency == "monthly" {
			year := currentDate.Year()
			month := currentDate.Month()
			day := startDateObj.Day()
			nextResetDate = time.Date(year, month, day, 0, 0, 0, 0, currentDate.Location())
			if nextResetDate.Before(currentDate) {
				nextResetDate = nextResetDate.AddDate(0, 1, 0)
			}
		} else if frequency == "quarterly" {
			diff := currentDate.Sub(startDateObj)
			months := int(diff.Hours() / (24 * 30))
			monthsAdded := int(3-months%3) + months
			nextResetDate = startDateObj.AddDate(0, monthsAdded, 0)
		} else if frequency == "yearly" {
			years := currentDate.Year() - startDateObj.Year()
			nextResetDate = startDateObj.AddDate(years+1, 0, 0)
		}
	} else {
		if frequency == "daily" {
			nextResetDate = startDateObj.AddDate(0, 0, 1)
		} else if frequency == "weekly" {
			nextResetDate = startDateObj.AddDate(0, 0, 7)
		} else if frequency == "monthly" {
			nextResetDate = startDateObj.AddDate(0, 1, 0)
		} else if frequency == "quarterly" {
			nextResetDate = startDateObj.AddDate(0, 3, 0)
		} else if frequency == "yearly" {
			nextResetDate = startDateObj.AddDate(1, 0, 0)
		}
	}
	return nextResetDate
}

func ToProfileTrackEntry(entity *dao.PartnerProfileTrackEntity) (*domain.PartnerProfileTrackEntry, error) {
	return &domain.PartnerProfileTrackEntry{
		ID:        entity.ID,
		PartnerId: entity.PartnerId,
		Key:       entity.Key,
		Enabled:   entity.Enabled,
	}, nil

}

func ToActivityEntry(ctx context.Context, entity dao.AcitivityTrackerEntity) (*domain.ActivityInput, error) {
	return &domain.ActivityInput{
		Id:        entity.Id,
		PartnerId: entity.PartnerId,
		Type:      entity.Type,
		UpdatedAt: entity.UpdatedAt.Unix(),
	}, nil
}
