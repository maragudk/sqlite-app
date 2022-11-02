package sql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/maragudk/sqlite-app/model"
)

// GetTOC of all articles with no content.
func (d *Database) GetTOC(ctx context.Context) (as []model.Article, err error) {
	err = d.DB.SelectContext(ctx, &as, `select id, title from articles order by created desc`)
	return
}

// CreateArticle with title and content, ignoring any ID or timestamps.
func (d *Database) CreateArticle(ctx context.Context, a model.Article) error {
	_, err := d.DB.NamedExecContext(ctx, `insert into articles (title, content) values (:title, :content)`, a)
	return err
}

// GetArticle by ID, returning nil if no such ID exists.
func (d *Database) GetArticle(ctx context.Context, id int) (*model.Article, error) {
	var a model.Article
	if err := d.DB.GetContext(ctx, &a, `select * from articles where id = ?`, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &a, nil
}
