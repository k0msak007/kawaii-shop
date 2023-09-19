package productsPatterns

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/k0msak007/kawaii-shop/modules/entities"
	"github.com/k0msak007/kawaii-shop/modules/files/filesUsecases"
	"github.com/k0msak007/kawaii-shop/modules/products"
)

type IUpdateProductBuilder interface {
	initTransaction() error
	initQuery()
	updateTitleQuery()
	updateDescriptionQuery()
	updatePriceQuery()
	updateCategory() error
	insertImages() error
	getOldImages() []*entities.Image
	deleteOldImages() error
	closeQuery()
	updateProduct() error
	getQueryFields() []string
	getValues() []any
	getQuery() string
	setQuery(query string)
	getImagesLen() int
	commit() error
}

type updateProductBuilder struct {
	db             *sqlx.DB
	tx             *sqlx.Tx
	req            *products.Product
	filesUsecases  filesUsecases.IFilesUsecase
	query          string
	queryFields    []string
	lastStackIndex int
	values         []any
}

func UpdateProductBuilder(db *sqlx.DB, req *products.Product, filesUsecases filesUsecases.IFilesUsecase) IUpdateProductBuilder {
	return &updateProductBuilder{
		db:            db,
		req:           req,
		filesUsecases: filesUsecases,
		queryFields:   make([]string, 0),
		values:        make([]any, 0),
	}
}

type updateProductEngineer struct {
	builder IUpdateProductBuilder
}

func (b *updateProductBuilder) initTransaction() error {
	tx, err := b.db.BeginTxx(context.Background(), nil)
	if err != nil {
		return err
	}
	b.tx = tx
	return nil
}

func (b *updateProductBuilder) initQuery() {
	b.query += `
		UPDATE "products" SET
	`
}

func (b *updateProductBuilder) updateTitleQuery() {
	if b.req.Title != "" {
		b.values = append(b.values, b.req.Title)
		b.lastStackIndex = len(b.values)

		b.queryFields = append(b.queryFields, fmt.Sprintf(`
			"title" = $%d
		`, b.lastStackIndex))
	}
}

func (b *updateProductBuilder) updateDescriptionQuery() {
	if b.req.Description != "" {
		b.values = append(b.values, b.req.Description)
		b.lastStackIndex = len(b.values)

		b.queryFields = append(b.queryFields, fmt.Sprintf(`
			"description" = $%d
		`, b.lastStackIndex))
	}
}

func (b *updateProductBuilder) updatePriceQuery() {
	if b.req.Price != 0 {
		b.values = append(b.values, b.req.Price)
		b.lastStackIndex = len(b.values)

		b.queryFields = append(b.queryFields, fmt.Sprintf(`
			"price" = $%d
		`, b.lastStackIndex))
	}
}

func (b *updateProductBuilder) updateCategory() error {
	if b.req.Category == nil {
		return nil
	}
	if b.req.Category.Id == 0 {
		return nil
	}

	query := `
		UPDATE "products_categories" SET "category_id" = $1
		WHERE "product_id" = $2;
	`

	if _, err := b.tx.ExecContext(context.Background(), query, b.req.Category.Id, b.req.Id); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("update products_categories failed: %v", err)
	}
	return nil
}

func (b *updateProductBuilder) insertImages() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	query := `
		INSERT INTO "images" (
			"filename",
			"url",
			"product_id"
		) VALUES
	`

	valueStack := make([]any, 0)
	var index int
	for i := range b.req.Image {
		valueStack = append(valueStack,
			b.req.Image[i].FileName,
			b.req.Image[i].Url,
			b.req.Id,
		)

		if i != len(b.req.Image)-1 {
			query += fmt.Sprintf(`($%d, $%d, $%d),`, index+1, index+2, index+3)
		} else {
			query += fmt.Sprintf(`($%d, $%d, $%d);`, index+1, index+2, index+3)
		}
		index += 3
	}

	if _, err := b.tx.ExecContext(
		ctx,
		query,
		valueStack...,
	); err != nil {
		b.tx.Rollback()
		fmt.Printf("insert images failed: %v", err)
	}

	return nil
}

func (b *updateProductBuilder) getOldImages() []*entities.Image {
	return nil
}

func (b *updateProductBuilder) deleteOldImages() error {
	return nil
}

func (b *updateProductBuilder) closeQuery() {
	b.values = append(b.values, b.req.Id)
	b.lastStackIndex = len(b.values)

	b.queryFields = append(b.queryFields, fmt.Sprintf(`
			WHERE "id" = $%d
		`, b.lastStackIndex))
}

func (b *updateProductBuilder) updateProduct() error {
	return nil
}

func (b *updateProductBuilder) getQueryFields() []string {
	return nil
}

func (b *updateProductBuilder) getValues() []any {
	return nil
}

func (b *updateProductBuilder) getQuery() string {
	return ""
}

func (b *updateProductBuilder) setQuery(query string) {}

func (b *updateProductBuilder) getImagesLen() int {
	return 0
}

func (b *updateProductBuilder) commit() error {
	if err := b.tx.Commit(); err != nil {
		return err
	}
	return nil
}

func UpdateProductEngineer(b IUpdateProductBuilder) *updateProductEngineer {
	return &updateProductEngineer{
		builder: b,
	}
}

func (en *updateProductEngineer) UpdateProduct() error {
	return nil
}
