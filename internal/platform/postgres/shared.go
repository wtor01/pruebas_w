package postgres

import (
	"gorm.io/gorm"
	"strings"
)

type Paginate struct {
	Limit  int
	Offset int
}

func NewPaginate(limit int, offset *int) Paginate {
	o := 0
	if offset != nil {
		o = *offset
	}
	return Paginate{Limit: limit, Offset: o}
}

func WithPaginate(paginate Paginate) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(paginate.Limit).Offset(paginate.Offset)
	}
}

func FetchAndCount(db *gorm.DB) (*gorm.DB, int64) {
	var count int64
	result := db.Limit(-1).Offset(-1).Count(&count)

	return result, count
}

func WithLike(query string, keys ...string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(query) <= 0 {
			return db
		}
		var strQueries []string
		var values []interface{}

		for _, k := range keys {
			strQueries = append(strQueries, k+" ILIKE ?")
			values = append(values, "%"+query+"%")
		}

		return db.Where(strings.Join(strQueries, " OR "), values...)
	}
}

func WithOrder(validColumns map[string]struct{}, order map[string]string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for key, sort := range order {
			_, ok := validColumns[key]
			if ok {
				db.Order(key + " " + sort)
			}

		}
		return db
	}
}

func WithValues(validColumns map[string]struct{}, values map[string]string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for key, value := range values {
			_, ok := validColumns[key]
			if ok {
				db.Where(key, value)
			}

		}
		return db
	}
}
