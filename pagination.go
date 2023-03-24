package dmp

// This needs to be properly tested

import (
	"fmt"
	"math"
	"regexp"
	"strings"

	"gorm.io/gorm"
)

type pagination interface {
	Run(count int64, pagination *Pagination) func(db *gorm.DB) *gorm.DB
}

type Pagination struct {
	Page          int32
	PageSize      int32
	Total         int32
	TotalPages    int32
	OrderByColumn string
	Sort          string
	Override      bool
}

func (p *Pagination) Run(count int64, pagination *Pagination) func(db *gorm.DB) *gorm.DB {
	p = pagination
	return func(db *gorm.DB) *gorm.DB {
		if pagination.Override {
			return db
		}

		pageSize := p.GetPageSize()
		pagination.Total = int32(count)
		totalPages := int32(math.Ceil(float64(count) / float64(pageSize)))
		pagination.TotalPages = totalPages

		return db.Offset(p.GetOffset()).Limit(pageSize).Order(p.GetOrderByColumn())
	}
}

// GetOffset return current offset for the pagination struct
func (p *Pagination) GetOffset() int {
	// kept original line bellow since I'm not sure where
	// p.GetTotal() comes from...
	// return int((p.GetPage() - 1) * p.GetTotal())

	return int((p.GetPage() - 1) * int(p.Total))

}

func (p *Pagination) GetPageSize() int {
	if p.PageSize == 0 {
		p.PageSize = 10
	}
	return int(p.PageSize)
}

func (p *Pagination) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}
	return int(p.Page)
}

func (p *Pagination) GetOrderByColumn() string {
	if p.OrderByColumn == "" {
		p.OrderByColumn = "id"
	}
	return fmt.Sprintf("%s %s", strings.ToLower(p.toSnakeCase((p.OrderByColumn))), strings.ToLower(p.Sort))
}

// ToSnakeCase format string to snake case
func (p *Pagination) toSnakeCase(str string) string {
	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
