package dmp

import (
	"fmt"

	"gorm.io/gorm"
)

//
////
// Not sure about this defitnion for Filter interface, an adapter pattern might be used better
// but the main problem is a predifend FitlerOptions struct, which is not
// flexible enough, and it is not possible to add custom filter options.
// The best solution would be to use a map[string]interface{} <-- this was provided by AI, can we use that
// insted of a predifend type for filter options?
// In anyway I would think about replacing FilerOptions with a map[string]interface{} or use
// generics if possible so that bot the type and options can be defined by the developer in
// the repository layer and/or the rest/grpc layer.
//
////
//

type Filter interface {
	Run(options *FilterOptions, model interface{}) func(db *gorm.DB) *gorm.DB
}

type FilterPair struct {
	Key   string
	Value string
}

type FilterOptions struct {
	Filters []FilterPair
}

//
////
// DefaultFilter is the default filter implementation of the Filter interface
// it is automatically select as defualt printer in the repository.go, and it can be easily replaced
// by creating a new struct that implements the Filter interface and then passing it to the repository
// using Repository.AddFilter(filter *Filter) method
////
//

type DefaultFilter struct {
	*FilterOptions
	FieldString string
}

func (filter *DefaultFilter) Run(options *FilterOptions, model interface{}) func(db *gorm.DB) *gorm.DB {
	// activeFilters will contain all valid key:value pairs that were sent in filterOptions
	activeFilters := make(map[string]interface{})
	filter.FilterOptions = options
	filter.FieldString = "%s.%s"

	return func(db *gorm.DB) *gorm.DB {
		stmt := &gorm.Statement{DB: db}

		if filter.FilterOptions == nil {
			return db
		}

		if err := stmt.Parse(&model); err != nil {
			return db
		}

		for _, f := range filter.FilterOptions.Filters {
			for _, field := range stmt.Schema.Fields {
				if field.Name == f.Key {
					activeFilters[fmt.Sprintf(filter.FieldString, stmt.Schema.Table, field.DBName)] = f.Value
				}
			}
		}

		return db.Where(activeFilters)
	}
}
