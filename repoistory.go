package dmp

import (
	"fmt"

	"github.com/GoLangWebSDK/dmp/database"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository[T any] struct {
	DB             *database.Database
	Model          T
	stmt           *gorm.Statement
	query          *gorm.DB
	filter         Filter
	pagination     pagination
	deletedAtQuery string
}

func NewRepository[T any](model T) *Repository[T] {
	return &Repository[T]{
		DB:             database.DBManager,
		Model:          model,
		stmt:           &gorm.Statement{DB: database.DBManager.Engine},
		query:          database.DBManager.Engine.Session(&gorm.Session{}),
		filter:         &DefaultFilter{},
		pagination:     &Pagination{},
		deletedAtQuery: "%s.deleted_at IS NULL",
	}
}

func (repo *Repository[T]) Add(model T) error {
	return repo.DB.Engine.Create(&model).Error
}

func (repo *Repository[T]) GetAll() ([]T, error) {
	results := []T{}
	err := repo.query.Find(&results).Error
	// not sure if this is needed, I'm trying to reset repo.query
	// for next call by creating new query session after the finisher
	// GetAll is called
	repo.query = repo.DB.Engine.Session(&gorm.Session{})
	return results, err
}

func (repo *Repository[T]) Get(ID uint32) (T, error) {
	err := repo.DB.Engine.Where("id = ?", ID).First(&repo.Model).Error
	return repo.Model, err
}

func (repo *Repository[T]) Update(ID uint32, model T) error {
	if ID == 0 {
		return fmt.Errorf("Missing ID")
	}
	err := repo.DB.Engine.First(&repo.Model, ID).Error

	if err != nil {
		return err
	}

	err = repo.DB.Engine.Model(&model).Where("id = ?", ID).Updates(model).Error

	return err
}

func (repo *Repository[T]) Delete(ID uint32) error {
	if ID == 0 {
		return fmt.Errorf("Missing ID")
	}
	err := repo.DB.Engine.First(&repo.Model, ID).Error

	if err != nil {
		return err
	}

	err = repo.DB.Engine.Delete(&repo.Model, ID).Error
	return err
}

func (repo *Repository[T]) Filter(options *FilterOptions) *Repository[T] {
	repo.query = repo.query.Model(&repo.Model).
		Preload(clause.Associations).
		Scopes(repo.filter.Run(options, repo.Model))
	return repo
}

func (repo *Repository[T]) SetFilter(filter Filter) *Repository[T] {
	repo.filter = filter
	return repo
}

func (repo *Repository[T]) Paginate(p *Pagination) *Repository[T] {
	var count int64
	repo.query = repo.query.Count(&count).Scopes(repo.pagination.Run(count, p))
	return repo
}
