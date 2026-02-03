package manager

import (
	"coffee/app/discovery/dao"
	"coffee/app/discovery/domain"
	coredomain "coffee/core/domain"
	"context"
	"strings"
)

type LocationManager struct {
	locationsDao *dao.LocationsDao
}

func CreateLocationManager(ctx context.Context) *LocationManager {
	locationsDao := dao.CreateLocationsDao(ctx)
	managerLocation := &LocationManager{
		locationsDao: locationsDao,
	}
	return managerLocation
}

func (m *LocationManager) SearchLocations(ctx context.Context, query coredomain.SearchQuery, sortBy string, sortDir string, page int, size int) (*[]domain.LocationsEntry, int64, error) {
	entities, filteredCount, err := m.locationsDao.Search(ctx, query, sortBy, sortDir, page, size)
	var locations []domain.LocationsEntry
	for _, entity := range entities {
		locationEntry, _ := ToLocationsEntry(entity)
		locations = append(locations, *locationEntry)
	}
	return &locations, filteredCount, err
}

func (m *LocationManager) SearchLocationsNew(ctx context.Context, searchQuery coredomain.SearchQuery, size int) ([]domain.LocationsEntry, error) {
	var locations []string
	var err error
	table_name := "mv_location_master"
	var name string
	filters := searchQuery.Filters
	for i := range filters {
		name = filters[i].Value
	}
	locations, err = m.locationsDao.FindLocationList(ctx, name, table_name, size)
	var locationsArray []domain.LocationsEntry
	for i := range locations {
		locationEntry := domain.LocationsEntry{}
		if strings.Contains(locations[i], "city_") {
			name := strings.Replace(locations[i], "city_", "", -1)
			if name != "" {
				locationEntry.Name = name
				locationEntry.FullName = locationEntry.Name + ", India"
				locationEntry.Type = "city"
				locationsArray = append(locationsArray, locationEntry)
			}
		}
		if strings.Contains(locations[i], "state_") {
			name := strings.Replace(locations[i], "state_", "", -1)
			if name != "" {
				locationEntry.Name = name
				locationEntry.FullName = locationEntry.Name + ", India"
				locationEntry.Type = "state"
				locationsArray = append(locationsArray, locationEntry)
			}
		}
		if strings.Contains(locations[i], "country_") {
			name := strings.Replace(locations[i], "country_", "", -1)
			if name != "" {
				locationEntry.Name = name
				locationEntry.FullName = "India"
				locationEntry.Type = "country"
				locationsArray = append(locationsArray, locationEntry)
			}
		}
	}
	return locationsArray, err
}

func (m *LocationManager) SearchGccLocations(ctx context.Context, name string) ([]string, error) {
	var locations []string
	var err error
	table_name := "mv_location_master_tokens"

	locations, err = m.locationsDao.FindLocationList(ctx, name, table_name, 10)
	return locations, err
}
