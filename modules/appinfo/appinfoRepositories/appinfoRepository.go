package appinfoRepositories

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/k0msak007/kawaii-shop/modules/appinfo"
)

type IAppinfoRepository interface {
	FindCategory(req *appinfo.CategoryFilter) ([]*appinfo.Category, error)
}

type appinfoRepository struct {
	db *sqlx.DB
}

func AppinfoRepository(db *sqlx.DB) IAppinfoRepository {
	return &appinfoRepository{
		db: db,
	}
}

func (r *appinfoRepository) FindCategory(req *appinfo.CategoryFilter) ([]*appinfo.Category, error) {
	query := `
		SELECT
			"id",
			"title"
		FROM "categories"
	`

	filterValues := make([]any, 0)
	if req.Title != "" {
		query += `
			WHERE LOWER("title") LIKE $1
		`

		filterValues = append(filterValues, "%"+strings.ToLower(req.Title)+"%")
	}
	query += ";"

	category := make([]*appinfo.Category, 0)
	if err := r.db.Select(&appinfo.Category{}, query, filterValues...); err != nil {
		return nil, fmt.Errorf("select category failed: %v", err)
	}

	return category, nil
}
