package sql

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/maragudk/sqlite-app/model"
)

// GetTOC of all articles with no content.
func (d *Database) GetTOC(ctx context.Context) ([]model.Article, error) {
	var as []model.Article
	err := d.DB.SelectContext(ctx, &as, `select id, title from articles order by created desc`)
	return as, err
}

// CreateArticle with title and content, ignoring any ID or timestamps.
func (d *Database) CreateArticle(ctx context.Context, a model.Article) error {
	a.Title = strings.ReplaceAll(a.Title, "␟", "")
	a.Content = strings.ReplaceAll(a.Content, "␟", "")
	_, err := d.DB.NamedExecContext(ctx, `insert into articles (title, content) values (:title, :content)`, a)
	return err
}

// GetArticle by ID, returning nil if no such ID exists.
// If search is not empty, highlight the given search query in the title and content.
func (d *Database) GetArticle(ctx context.Context, id int, search string) (*model.Article, error) {
	var a model.Article

	query := `select * from articles where id = ?`
	var args []any
	args = append(args, id)

	if search != "" {
		query = `
			select
				a.id,
				highlight(articles_fts, 0, '␟', '␟') title,
				highlight(articles_fts, 1, '␟', '␟') content,
				a.created,
				a.updated
			from articles a
				join articles_fts af on (af.rowid = a.id)
			where id = ? and articles_fts match ?`
		args = append(args, escapeSearch(search))
	}

	if err := d.DB.GetContext(ctx, &a, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &a, nil
}

// SearchArticles with the given search query. Matches in titles are highlighted with the unit separator character ␟.
// Matches in content return a snippet of the content, also highlighted with the unit separator character ␟.
// Results are ordered by the internal rank of fts5.
// See https://www.sqlite.org/fts5.html
func (d *Database) SearchArticles(ctx context.Context, search string) ([]model.Article, error) {
	var as []model.Article
	query := `
		select
			a.id,
			highlight(articles_fts, 0, '␟', '␟') title,
			snippet(articles_fts, 1, '␟', '␟', '', 8) content,
			a.created,
			a.updated
		from articles a
			join articles_fts af on (af.rowid = a.id)
		where articles_fts match ?
		order by rank`
	err := d.DB.SelectContext(ctx, &as, query, escapeSearch(search))
	return as, err
}

func escapeSearch(s string) string {
	s = strings.ReplaceAll(s, `"`, `""`)
	return `"` + s + `"`
}
