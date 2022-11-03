package sql_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/maragudk/sqlite-app/model"
	"github.com/maragudk/sqlite-app/sqltest"
)

func TestDatabase_GetTOC(t *testing.T) {
	t.Run("gets all articles with only id and title reverse chronological order", func(t *testing.T) {
		db := sqltest.CreateDatabase(t)

		err := db.CreateArticle(context.Background(), model.Article{
			Title:   "Foo",
			Content: "Bar",
		})
		require.NoError(t, err)

		err = db.CreateArticle(context.Background(), model.Article{
			Title:   "Baz",
			Content: "Boo",
		})
		require.NoError(t, err)

		as, err := db.GetTOC(context.Background())
		require.NoError(t, err)
		require.Len(t, as, 2)
		require.Equal(t, 2, as[0].ID)
		require.Equal(t, "Baz", as[0].Title)
		require.Equal(t, "", as[0].Content)

		require.Equal(t, 1, as[1].ID)
	})
}

func TestDatabase_CreateArticle(t *testing.T) {
	t.Run("discards the unit separator character in title and content", func(t *testing.T) {
		db := sqltest.CreateDatabase(t)

		err := db.CreateArticle(context.Background(), model.Article{
			Title:   "Foo␟",
			Content: "Bar␟",
		})
		require.NoError(t, err)

		a, err := db.GetArticle(context.Background(), 1, "")
		require.NoError(t, err)
		require.NotNil(t, a)
		require.Equal(t, "Foo", a.Title)
		require.Equal(t, "Bar", a.Content)
	})
}

func TestDatabase_GetArticle(t *testing.T) {
	t.Run("gets an article", func(t *testing.T) {
		db := sqltest.CreateDatabase(t)

		err := db.CreateArticle(context.Background(), model.Article{
			Title:   "Foo",
			Content: "Bar",
		})
		require.NoError(t, err)

		a, err := db.GetArticle(context.Background(), 1, "")
		require.NoError(t, err)
		require.NotNil(t, a)
		require.Equal(t, 1, a.ID)
		require.Equal(t, "Foo", a.Title)
		require.Equal(t, "Bar", a.Content)
		require.WithinDuration(t, time.Now(), a.Created.T, time.Second)
		require.WithinDuration(t, time.Now(), a.Updated.T, time.Second)
	})

	t.Run("returns nil on no such id", func(t *testing.T) {
		db := sqltest.CreateDatabase(t)

		a, err := db.GetArticle(context.Background(), 1, "")
		require.NoError(t, err)
		require.Nil(t, a)
	})

	t.Run("highlights substrings if search given", func(t *testing.T) {
		db := sqltest.CreateDatabase(t)

		err := db.CreateArticle(context.Background(), model.Article{
			Title:   "The Foo Bar",
			Content: "Foo Bar Foo",
		})
		require.NoError(t, err)

		a, err := db.GetArticle(context.Background(), 1, "foo")
		require.NoError(t, err)
		require.NotNil(t, a)
		require.Equal(t, "The ␟Foo␟ Bar", a.Title)
		require.Equal(t, "␟Foo␟ Bar ␟Foo␟", a.Content)
	})
}

func TestDatabase_SearchArticles(t *testing.T) {
	db := sqltest.CreateDatabase(t)
	err := db.CreateArticle(context.Background(), model.Article{
		Title:   "The Foo is great",
		Content: "I wish that bar was also.",
	})
	require.NoError(t, err)

	err = db.CreateArticle(context.Background(), model.Article{
		Title:   "Bar me up a notch",
		Content: "Boo ya.",
	})
	require.NoError(t, err)

	t.Run("searches article titles and highlights", func(t *testing.T) {
		as, err := db.SearchArticles(context.Background(), "bar")
		require.NoError(t, err)
		require.Len(t, as, 2)

		require.Equal(t, "␟Bar␟ me up a notch", as[0].Title)
		require.Equal(t, "Boo ya.", as[0].Content)

		require.Equal(t, "The Foo is great", as[1].Title)
		require.Equal(t, "I wish that ␟bar␟ was also.", as[1].Content)
	})
}
