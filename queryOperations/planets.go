package queryoperations

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FilterParam struct {
    Eq    interface{} `json:"eq,omitempty"`
    Neq   interface{} `json:"neq,omitempty"`
    Gt    interface{} `json:"gt,omitempty"`
    Gte   interface{} `json:"gte,omitempty"`
    Lt    interface{} `json:"lt,omitempty"`
    Lte   interface{} `json:"lte,omitempty"`
    Like  string      `json:"like,omitempty"`
    In    []string    `json:"in,omitempty"`
    NotIn []string    `json:"notin,omitempty"`
    Or    []FilterParam `json:"or,omitempty"` // future scope
}

type QueryParams struct {
    Sort   string `form:"sort"`
    Filters map[string]FilterParam `form:"-"`
    Page   int    `form:"page"`
    Limit  int    `form:"limit"`
}

func (q *QueryParams) BindQuery(c *gin.Context) error {
    // Bind the simple fields
    if err := c.ShouldBindQuery(q); err != nil {
        return err
    }

    // Parse the filters
    q.Filters = make(map[string]FilterParam)
    filterQuery := c.QueryMap("filter")
    for field, value := range filterQuery {
        var filterParam FilterParam
        if err := json.Unmarshal([]byte(value), &filterParam); err != nil {
            return fmt.Errorf("invalid filter for field %s: %v", field, err)
        }
        q.Filters[field] = filterParam
    }

    return nil
}

func Filter(db *gorm.DB, params *QueryParams, allowedFilters *map[string]string) *gorm.DB {
	for field, filter := range params.Filters {
        if dataType, allowed := (*allowedFilters)[field]; allowed {
            if filter.Eq != nil {
                db = db.Where(fmt.Sprintf("%s = ?", field), filter.Eq)
            }
            if filter.Neq != nil {
                db = db.Where(fmt.Sprintf("%s != ?", field), filter.Neq)
            }
            if filter.Gt != nil {
                db = db.Where(fmt.Sprintf("%s > ?", field), filter.Gt)
            }
            if filter.Gte != nil {
                db = db.Where(fmt.Sprintf("%s >= ?", field), filter.Gte)
            }
            if filter.Lt != nil {
                db = db.Where(fmt.Sprintf("%s < ?", field), filter.Lt)
            }
            if filter.Lte != nil {
                db = db.Where(fmt.Sprintf("%s <= ?", field), filter.Lte)
            }
            if filter.Like != "" && dataType == "string" {
                db = db.Where(fmt.Sprintf("%s LIKE ?", field), "%"+filter.Like+"%")
            }
            if len(filter.In) > 0 {
                db = db.Where(fmt.Sprintf("%s IN (?)", field), filter.In)
            }
            if len(filter.NotIn) > 0 {
                db = db.Where(fmt.Sprintf("%s NOT IN (?)", field), filter.NotIn)
            }
        }
    }
    return db
}

func Sort(db *gorm.DB, params *QueryParams) *gorm.DB {
    if params.Sort != "" {
        db = db.Order(params.Sort)
    }
    return db
}

func Paginate(db *gorm.DB, params *QueryParams) *gorm.DB {
    if params.Page > 0 && params.Limit > 0 {
		offset := (params.Page - 1) * params.Limit
		return db.Limit(params.Limit).Offset(offset)
	}
	return db
}

func Apply(db *gorm.DB, params *QueryParams, allowedFilters *map[string]string) *gorm.DB {
	db = Filter(db, params, allowedFilters)
	db = Sort(db, params)
	db = Paginate(db, params)
	return db
}