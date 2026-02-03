package rest

import (
	"coffee/core/domain"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"golang.org/x/exp/slices"

	"github.com/go-chi/render"

	"github.com/go-chi/chi/v5"
)

func PathParam(r *http.Request, param string) string {
	return chi.URLParam(r, param)
}

func ParseBody[T any](w http.ResponseWriter, r *http.Request, body *T) bool {
	rawRequestBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return true
	}
	err = json.Unmarshal(rawRequestBody, body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return true
	}
	return false
}

func RenderError(w http.ResponseWriter, r *http.Request, err error) {
	errorResponse := domain.BaseResponse{Status: domain.Status{
		Status:  "ERROR",
		Message: err.Error(),
	}}
	render.JSON(w, r, errorResponse)
}

func ParseSearchParameters(w http.ResponseWriter, r *http.Request) (string, string, int, int, domain.SearchQuery, bool) {
	sortBy := r.URL.Query().Get("sortBy")
	sortDir := r.URL.Query().Get("sortDir")
	pageStr := r.URL.Query().Get("cursor")
	sizeStr := r.URL.Query().Get("size")
	page := 1
	size := 50
	if sortBy == "" {
		sortBy = "id"
	}
	if sortDir == "" {
		sortDir = "DESC"
	}
	var err error
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			RenderError(w, r, err)
		}
		if page >= 1000 {
			page = 1
		}
	}
	if sizeStr != "" {
		size, err = strconv.Atoi(sizeStr)
		if err != nil {
			RenderError(w, r, err)
		}
		if size >= 100 {
			size = 50
		}
	}
	var query domain.SearchQuery
	if r.Body == http.NoBody {
		return sortBy, sortDir, page, size, query, false
	}
	err = json.NewDecoder(r.Body).Decode(&query)
	if err != nil {
		RenderError(w, r, err)
		return "", "", 0, 0, domain.SearchQuery{}, true
	}
	return sortBy, sortDir, page, size, query, false
}

func ParseCollectionProfilesParameters(w http.ResponseWriter, r *http.Request) (string, string, int, int, bool) {
	sortBy := r.URL.Query().Get("sortBy")
	sortDir := r.URL.Query().Get("sortDir")
	pageStr := r.URL.Query().Get("cursor")
	sizeStr := r.URL.Query().Get("size")

	page := 1
	size := 50
	if sortBy == "" {
		sortBy = "id"
	}
	if sortDir == "" {
		sortDir = "DESC"
	}
	var err error
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			RenderError(w, r, err)
		}
		if page >= 1000 {
			page = 1
		}
	}
	if sizeStr != "" {
		size, err = strconv.Atoi(sizeStr)
		if err != nil {
			RenderError(w, r, err)
		}
		if size >= 100 {
			size = 50
		}
	}

	return sortBy, sortDir, page, size, false
}

func ParseActivityParameters(w http.ResponseWriter, r *http.Request) (string, string, int, int, bool) {
	sortBy := r.URL.Query().Get("sortBy")
	sortDir := r.URL.Query().Get("sortDir")
	pageStr := r.URL.Query().Get("cursor")
	sizeStr := r.URL.Query().Get("size")
	//activityType := r.URL.Query().Get("activityType")

	page := 1
	size := 50
	if sortBy == "" {
		sortBy = "id"
	}
	if sortDir == "" {
		sortDir = "DESC"
	}
	var err error
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			RenderError(w, r, err)
		}
		if page >= 1000 {
			page = 1
		}
	}
	if sizeStr != "" {
		size, err = strconv.Atoi(sizeStr)
		if err != nil {
			RenderError(w, r, err)
		}
		if size >= 100 {
			size = 50
		}
	}

	return sortBy, sortDir, page, size, false
}

func GetFilterValueForKey(filters []domain.SearchFilter, key string) *string {
	if len(filters) > 0 {
		for _, filter := range filters {
			if filter.Field == key {
				return &filter.Value
			}
		}
	}
	return nil
}

func RemoveKeysFromQuery(query domain.SearchQuery, keys ...string) domain.SearchQuery {
	var newFilters []domain.SearchFilter
	if len(query.Filters) > 0 {
		for _, filter := range query.Filters {
			if !slices.Contains(keys, filter.Field) {
				newFilters = append(newFilters, filter)
			}
		}
	}
	query.Filters = newFilters
	return query
}

func GetFilterValueForKeyOperation(filters []domain.SearchFilter, key string, filterType string) *string {
	if len(filters) > 0 {
		for _, filter := range filters {
			if filter.Field == key && filter.FilterType == filterType {
				return &filter.Value
			}
		}
	}
	return nil
}

func GetFilterForKey(filters []domain.SearchFilter, key string) *domain.SearchFilter {
	if len(filters) > 0 {
		for _, filter := range filters {
			if filter.Field == key {
				return &filter
			}
		}
	}
	return nil
}

func ParseGenderAgeFilters(searchQuery domain.SearchQuery) domain.SearchQuery {
	genderMap := map[string]string{"MALE": "audience_gender.male_per", "FEMALE": "audience_gender.female_per"}
	ageMap := map[string]string{"13-17": "audience_age.13-17", "18-24": "audience_age.18-24", "25-34": "audience_age.25-34", "35-44": "audience_age.35-44", "45-54": "audience_age.45-54", "55-64": "audience_age.55-64", "65+": "audience_age.65+"}

	for filtterKey := range searchQuery.Filters {
		if searchQuery.Filters[filtterKey].Field == "audience_gender" {
			filterValue := searchQuery.Filters[filtterKey].Value
			if val, ok := genderMap[filterValue]; ok {
				searchQuery.Filters[filtterKey].FilterType = "GTE"
				searchQuery.Filters[filtterKey].Field = val
				searchQuery.Filters[filtterKey].Value = "40"
			}
		}
		if searchQuery.Filters[filtterKey].Field == "audience_age" {
			filterValue := searchQuery.Filters[filtterKey].Value
			if val, ok := ageMap[filterValue]; ok {
				searchQuery.Filters[filtterKey].FilterType = "GTE"
				searchQuery.Filters[filtterKey].Field = val
				searchQuery.Filters[filtterKey].Value = "10"
			}
		}
	}
	return searchQuery
}
