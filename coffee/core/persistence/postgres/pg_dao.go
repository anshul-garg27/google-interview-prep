package postgres

import (
	"coffee/constants"
	"coffee/core/appcontext"
	"coffee/core/domain"
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"gorm.io/gorm/clause"

	"gorm.io/gorm"
)

type PgDao[EN domain.Entity, I domain.ID] struct {
	ctx context.Context
}

func NewPgDao[EN domain.Entity, I domain.ID](ctx context.Context) *PgDao[EN, I] {
	return &PgDao[EN, I]{
		ctx: ctx,
	}
}

func (r *PgDao[EN, I]) SetCtx(ctx context.Context) {
	r.ctx = ctx
}

func (r *PgDao[EN, I]) GetSession(ctx context.Context) *gorm.DB {
	reqContext := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	session := reqContext.GetSession()
	reqContext.Mutex.Lock()
	if session == nil {
		// Create new txn if one already does not exist
		session = InitializePgSession(ctx)
		reqContext.SetSession(ctx, session)
	}
	reqContext.Mutex.Unlock()
	return session.(PgSession).db
}

func (r *PgDao[EN, I]) FindById(ctx context.Context, id I) (*EN, error) {
	var entity EN
	idStr := fmt.Sprintf("%v", id)
	db := r.GetSession(ctx)
	tx := db.WithContext(ctx).Model(&entity).Where("id = ? AND enabled = ?", idStr, true).First(&entity)
	if tx.Error != nil && tx.Error.Error() == "record not found" {
		return nil, fmt.Errorf("invalid Record Id - %v", id)
	}
	return &entity, tx.Error
}

func (r *PgDao[EN, I]) FindByIds(ctx context.Context, ids []I) ([]EN, error) {
	var entity EN
	db := r.GetSession(ctx)
	entities := []EN{}

	tx := db.WithContext(ctx).Model(&entity).Where("id IN ? AND enabled = ?", ids, true).Scan(&entities)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return entities, tx.Error
}

func (r PgDao[EN, I]) Create(ctx context.Context, entity *EN) (*EN, error) {
	if entity == nil {
		return nil, errors.New("null entity passed to create")
	}
	db := r.GetSession(ctx)
	tx := db.WithContext(ctx).Model(&entity).Create(entity)
	return entity, tx.Error
}

func (r PgDao[EN, I]) Update(ctx context.Context, id I, entity *EN) (*EN, error) {
	if entity == nil {
		return nil, errors.New("null entity passed to update")
	}
	var err error
	db := r.GetSession(ctx)
	db = db.Model(&entity).Clauses(clause.Returning{}).Where("id = ?", id).Updates(&entity)
	if db.Error != nil {
		return nil, err
	} else if db.RowsAffected == 0 {
		return nil, fmt.Errorf("invalid Record Id - %v", id)

	}
	return entity, nil
}

func (r PgDao[EN, I]) Search(ctx context.Context, query domain.SearchQuery, sortBy string, sortDir string, page int, size int) ([]EN, int64, error) {
	db := r.GetSession(ctx)
	var entities []EN
	filteredReq := db.Model(&entities)
	for i := range query.Filters {
		filter := query.Filters[i]
		if filter.Value != "" {
			r.addPredicateForSearch(filteredReq, filter)
		}
	}
	if strings.ToLower(sortDir) == "desc" {
		filteredReq.Order(fmt.Sprintf("%s %s NULLS LAST, %s %s NULLS LAST", sortBy, sortDir, "id", "ASC"))
	} else {
		filteredReq.Order(fmt.Sprintf("%s %s, %s %s", sortBy, sortDir, "id", "asc"))

	}
	var filteredCount int64
	//rand.Seed(time.Now().UnixNano())
	//filteredCount = int64(rand.Intn(200000-10+1) + 10)
	filteredReq.Count(&filteredCount)
	paginatedReq := filteredReq.Limit(size)
	paginatedReq = paginatedReq.Offset((page - 1) * size)
	paginatedReq.Find(&entities)
	if len(entities) < size {
		filteredCount = int64(len(entities))
	}
	if filteredReq.Error != nil {
		return nil, 0, filteredReq.Error
	}
	return entities, filteredCount, nil
}

func (r PgDao[EN, I]) SearchJoins(ctx context.Context, query domain.SearchQuery, sortBy string, sortDir string, page int, size int, joinTables []domain.JoinClauses) ([]EN, int64, error) {
	db := r.GetSession(ctx)
	var entities []EN
	filteredReq := db.Model(&entities)

	for key := range joinTables {
		if joinTables[key].PreloadVariable != "" {
			preloadVariable := joinTables[key].PreloadVariable

			preloadCondition := joinTables[key].PreloadCondition
			preloadValue := joinTables[key].PreloadValue

			if preloadCondition == "" || preloadValue == "" {
				filteredReq = filteredReq.Preload(preloadVariable)
			} else {
				filteredReq = filteredReq.Preload(preloadVariable, preloadCondition, preloadValue)
			}

		}
		joinClause := joinTables[key].JoinClause
		filteredReq = filteredReq.Joins(joinClause)
	}
	for i := range query.Filters {
		filter := query.Filters[i]
		if filter.Value != "" {
			r.AddPredicateForSearchJoins(filteredReq, filter)
		}
	}
	// filteredReq = filteredReq.Distinct()
	if strings.ToLower(sortDir) == "desc" {
		filteredReq.Order(fmt.Sprintf("%s %s NULLS LAST", sortBy, sortDir))
	} else {
		filteredReq.Order(fmt.Sprintf("%s %s", sortBy, sortDir))
	}

	var filteredCount int64
	//rand.Seed(time.Now().UnixNano())
	//filteredCount = int64(rand.Intn(200000-10+1) + 10)
	filteredReq.Count(&filteredCount)

	paginatedReq := filteredReq.Limit(size)
	paginatedReq = paginatedReq.Offset((page - 1) * size)
	paginatedReq.Find(&entities)
	if len(entities) < size {
		filteredCount = int64(len(entities))
	}
	if filteredReq.Error != nil {
		return nil, 0, filteredReq.Error
	}
	return entities, filteredCount, nil
}

func (r PgDao[EN, I]) AddPredicateForSearchJoins(req *gorm.DB, filter domain.SearchFilter) *gorm.DB {
	if strings.Count(filter.Field, ".") >= 2 && !strings.Contains(filter.Field, ",") && filter.FilterType != "ARRAY" && filter.FilterType != "ARRAY_IN" {
		filter.Field = makeJsonbFilterFieldForSearchJoins(filter.Field, filter.FilterType)
	}

	if filter.FilterType == "EQ" {
		if strings.Contains(filter.Field, ",") {
			filterType := "="
			req = makeOrFilter(req, filter, filterType)
		} else {
			req = req.Where(fmt.Sprintf(filter.Field)+" = ?", filter.Value)
		}

	} else if filter.FilterType == "NE" {
		req = req.Where(fmt.Sprintf(filter.Field)+" != ?", filter.Value)
	} else if filter.FilterType == "LIKE" {
		if strings.Contains(filter.Field, ",") {
			filterType := "ILIKE"
			filter.Value = "%" + filter.Value + "%"
			req = makeOrFilter(req, filter, filterType)
		} else {
			req = req.Where(fmt.Sprintf(filter.Field)+" ILIKE ?", "%"+filter.Value+"%")
		}
	} else if filter.FilterType == "GTE" {
		if strings.Contains(filter.Field, ",") {
			filterType := ">="
			req = makeOrFilter(req, filter, filterType)
		} else {
			value, err := strconv.Atoi(filter.Value)
			if err == nil {
				if strings.Contains(filter.Field, "->") {
					req = req.Where("(("+fmt.Sprintf("%s)::numeric >= ?", filter.Field)+")", value)
				} else {
					req = req.Where(fmt.Sprintf(filter.Field)+" >= ?", value)
				}

			} else {
				// for handling date type like '2023-07-01'
				req = req.Where(fmt.Sprintf(filter.Field)+" >= ?", filter.Value)
			}
		}
	} else if filter.FilterType == "LTE" {
		if strings.Contains(filter.Field, ",") {
			filterType := "<="
			req = makeOrFilter(req, filter, filterType)
		} else {
			value, err := strconv.Atoi(filter.Value)
			if err == nil {
				if strings.Contains(filter.Field, "->") {
					req = req.Where("(("+fmt.Sprintf("%s)::numeric <= ?", filter.Field)+")", value)
				} else {
					req = req.Where(fmt.Sprintf(filter.Field)+" <= ?", value)
				}
			} else {
				// for handling date type like '2023-07-01'
				req = req.Where(fmt.Sprintf(filter.Field)+" <= ?", filter.Value)
			}
		}

	} else if filter.FilterType == "GT" {
		if strings.Contains(filter.Field, ",") {
			filterType := ">"
			req = makeOrFilter(req, filter, filterType)
		} else {
			value, err := strconv.Atoi(filter.Value)
			if err == nil {
				if strings.Contains(filter.Field, "->") {
					req = req.Where("(("+fmt.Sprintf("%s)::numeric > ?", filter.Field)+")", value)
				} else {
					req = req.Where(fmt.Sprintf(filter.Field)+" > ?", value)
				}
			} else {
				// for handling date type like '2023-07-01'
				req = req.Where(fmt.Sprintf(filter.Field)+" > ?", filter.Value)
			}
		}

	} else if filter.FilterType == "LT" {
		if strings.Contains(filter.Field, ",") {
			filterType := "<"
			req = makeOrFilter(req, filter, filterType)
		} else {
			value, err := strconv.Atoi(filter.Value)
			if err == nil {
				if strings.Contains(filter.Field, "->") {
					req = req.Where("(("+fmt.Sprintf("%s)::numeric < ?", filter.Field)+")", value)
				} else {
					req = req.Where(fmt.Sprintf(filter.Field)+" < ?", value)
				}

			} else {
				// for handling date type like '2023-07-01'
				req = req.Where(fmt.Sprintf(filter.Field)+" < ?", filter.Value)
			}
		}
	} else if filter.FilterType == "IN" {
		if strings.Contains(filter.Field, ",") {
			filterType := "IN"
			req = makeOrFilter(req, filter, filterType)
		} else {
			if len(strings.Split(filter.Value, ",")) == 1 {
				req = req.Where(fmt.Sprintf(filter.Field)+" = ?", filter.Value)
			} else {
				req = req.Where(fmt.Sprintf(filter.Field)+" IN ?", strings.Split(filter.Value, ","))

			}
		}
	} else if filter.FilterType == "NOTIN" {
		req = req.Not(fmt.Sprintf(filter.Field), strings.Split(filter.Value, ","))
	} else if filter.FilterType == "BOOL" {
		value, err := strconv.ParseBool(filter.Value)
		if err == nil {
			req = req.Where(fmt.Sprintf(filter.Field)+" = ?", value)
		}
	} else if filter.FilterType == "ARRAY" || filter.FilterType == "ARRAY_IN" {
		if strings.Contains(filter.Field, ",") {
			filterType := "ARRAY_IN"
			req = makeOrFilter(req, filter, filterType)
		} else if strings.Count(filter.Field, ".") >= 2 {
			filter.Field = makeJsonbFilterFieldForSearchJoins(filter.Field, filter.FilterType)
			filter.Value = makeInFilterforJsonBColumn(filter.Value)
			req = req.Where(fmt.Sprintf("%s @>  %s", filter.Field, filter.Value))
		} else {
			query, valuesInterface := makeArrayFilter(filter.Field, filter.Value)
			req = req.Where(query, valuesInterface...)
		}

	}
	return req
}
func makeInFilterforJsonBColumn(filterValue string) string {
	values := strings.Split(filterValue, ",")

	jsonbStrings := make([]string, len(values))
	for i, v := range values {
		jsonbStrings[i] = fmt.Sprintf(`'["%s"]'::jsonb`, v)
	}
	// where admin_details->'campaignCategories' @> ANY(array['["Health"]'::jsonb, '["Lifestyle"]'::jsonb])
	filterValue = fmt.Sprintf("ANY(array[%s])", strings.Join(jsonbStrings, ", "))
	return filterValue
}

func makeOrFilter(req *gorm.DB, filter domain.SearchFilter, filterType string) *gorm.DB {
	var conditions []string
	var values, value []interface{}
	var condition string
	fields := strings.Split(filter.Field, ",")
	if filterType == "ARRAY_IN" {
		for _, field := range fields {
			if strings.Count(field, ".") >= 2 {
				filter.Field = makeJsonbFilterFieldForSearchJoins(field, filterType)
				filterValue := makeInFilterforJsonBColumn(filter.Value)
				conditions = append(conditions, fmt.Sprintf("%s @>  %s", filter.Field, filterValue))
			} else {
				condition, value = makeArrayFilter(field, filter.Value)
				conditions = append(conditions, condition)
				for i := range value {
					values = append(values, value[i])
				}
			}
		}
		query := strings.Join(conditions, " OR ")
		req = req.Where(query, values...)
	} else if filterType == "IN" {
		for _, field := range fields {
			if strings.Count(field, ".") >= 2 {
				field = makeJsonbFilterFieldForSearchJoins(field, filterType)
			}
			conditions = append(conditions, fmt.Sprintf("%s %s ?", field, filterType))
			values = append(values, strings.Split(filter.Value, ","))
		}
		query := strings.Join(conditions, " OR ")
		req = req.Where(query, values...)
	} else {
		for _, field := range fields {
			if strings.Count(field, ".") >= 2 {
				field = makeJsonbFilterFieldForSearchJoins(field, filterType)
			}
			conditions = append(conditions, fmt.Sprintf("%s %s ?", field, filterType))
			values = append(values, filter.Value)
		}
		query := strings.Join(conditions, " OR ")
		req = req.Where(query, values...)
	}
	return req
}

func makeArrayFilter(fieldName string, filedValue string) (string, []interface{}) {
	filedValueSlice := strings.Split(filedValue, ",")
	fieldValueInterface := make([]interface{}, len(filedValueSlice))
	for i, v := range filedValueSlice {
		fieldValueInterface[i] = interface{}(v)
	}
	var placeholders []string
	for range fieldValueInterface {
		placeholders = append(placeholders, "?")
	}
	query := fmt.Sprintf("%s && ARRAY[%s]::text[]", fieldName, strings.Join(placeholders, ", "))
	return query, fieldValueInterface
}

func makeJsonbFilterFieldForSearchJoins(field string, filterType string) string {
	newSlice := strings.Split(field, ".")
	field = ""
	for key, value := range newSlice {
		if key == 0 {
			field = field + value + "."
		} else if key == 1 {
			if len(newSlice) == 3 && filterType != "ARRAY_IN" && filterType != "ARRAY" {
				field = field + value + "->>'"
			} else {
				field = field + value + "->'"
			}

		} else if key == len(newSlice)-1 {
			field = field + value + "'"
		} else {
			field = field + value + "'->'"
		}
	}
	return field
}

func (r PgDao[EN, I]) addPredicateForSearch(req *gorm.DB, filter domain.SearchFilter) *gorm.DB {
	if strings.Contains(filter.Field, ".") && filter.FilterType != "ARRAY" && filter.FilterType != "ARRAY_IN" {
		filter.Field = makeJsonbFilterFieldForSearch(filter.Field, filter.FilterType)
	}
	if filter.FilterType == "EQ" {
		req = req.Where(fmt.Sprintf(filter.Field)+" = ?", filter.Value)
	} else if filter.FilterType == "NE" {
		req = req.Where(fmt.Sprintf(filter.Field)+" != ?", filter.Value)
	} else if filter.FilterType == "LIKE" {
		req = req.Where(fmt.Sprintf(filter.Field)+" ILIKE ?", "%"+filter.Value+"%")
	} else if filter.FilterType == "GTE" {
		value, err := strconv.Atoi(filter.Value)
		if err == nil {
			if strings.Contains(filter.Field, "->") {
				req = req.Where("(("+fmt.Sprintf("%s)::numeric >= ?", filter.Field)+")", value)
			} else {
				req = req.Where(fmt.Sprintf(filter.Field)+" >= ?", value)
			}
		} else {
			// for handling date type like '2023-07-01'
			req = req.Where(fmt.Sprintf(filter.Field)+" >= ?", filter.Value)
		}
	} else if filter.FilterType == "LTE" {
		value, err := strconv.Atoi(filter.Value)
		if err == nil {
			if strings.Contains(filter.Field, "->") {
				req = req.Where("(("+fmt.Sprintf("%s)::numeric <= ?", filter.Field)+")", value)
			} else {
				req = req.Where(fmt.Sprintf(filter.Field)+" <= ?", value)
			}
		} else {
			// for handling date type like '2023-07-01'
			req = req.Where(fmt.Sprintf(filter.Field)+" <= ?", filter.Value)
		}
	} else if filter.FilterType == "GT" {
		value, err := strconv.Atoi(filter.Value)
		if err == nil {
			if strings.Contains(filter.Field, "->") {
				req = req.Where("(("+fmt.Sprintf("%s)::numeric > ?", filter.Field)+")", value)
			} else {
				req = req.Where(fmt.Sprintf(filter.Field)+" > ?", value)
			}
		} else {
			// for handling date type like '2023-07-01'
			req = req.Where(fmt.Sprintf(filter.Field)+" > ?", filter.Value)
		}
	} else if filter.FilterType == "LT" {
		value, err := strconv.Atoi(filter.Value)
		if err == nil {
			if strings.Contains(filter.Field, "->") {
				req = req.Where("(("+fmt.Sprintf("%s)::numeric < ?", filter.Field)+")", value)
			} else {
				req = req.Where(fmt.Sprintf(filter.Field)+" < ?", value)
			}
		} else {
			// for handling date type like '2023-07-01'
			req = req.Where(fmt.Sprintf(filter.Field)+" < ?", filter.Value)

		}
	} else if filter.FilterType == "IN" {
		if len(strings.Split(filter.Value, ",")) == 1 {
			req = req.Where(fmt.Sprintf(filter.Field)+" = ?", filter.Value)
		} else {
			req = req.Where(fmt.Sprintf(filter.Field)+" IN ?", strings.Split(filter.Value, ","))
		}
	} else if filter.FilterType == "BOOL" {
		value, err := strconv.ParseBool(filter.Value)
		if err == nil {
			req = req.Where(fmt.Sprintf(filter.Field)+" = ?", value)
		}
	} else if filter.FilterType == "ARRAY" || filter.FilterType == "ARRAY_IN" {
		if strings.Contains(filter.Field, ".") {
			filter.Field = makeJsonbFilterFieldForSearch(filter.Field, filter.FilterType)
			filter.Value = makeInFilterforJsonBColumn(filter.Value)
			req = req.Where(fmt.Sprintf("%s @>  %s", filter.Field, filter.Value))
		} else {
			query, valuesInterface := makeArrayFilter(filter.Field, filter.Value)
			req = req.Where(query, valuesInterface...)
		}
	}
	return req
}

func makeJsonbFilterFieldForSearch(field string, filterType string) string {
	newSlice := strings.Split(field, ".")
	field = ""
	for key, value := range newSlice {
		if key == 0 {
			if len(newSlice) == 2 && filterType != "ARRAY_IN" && filterType != "ARRAY" {
				field = field + value + "->>'"
			} else {
				field = field + value + "->'"
			}
		} else if key == len(newSlice)-1 {
			field = field + value + "'"
		} else {
			field = field + value + "'->'"
		}
	}
	return field
}
