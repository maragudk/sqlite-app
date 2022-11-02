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

func TestDatabase_GetArticle(t *testing.T) {
	t.Run("gets an article", func(t *testing.T) {
		db := sqltest.CreateDatabase(t)

		err := db.CreateArticle(context.Background(), model.Article{
			Title:   "Foo",
			Content: "Bar",
		})
		require.NoError(t, err)

		a, err := db.GetArticle(context.Background(), 1)
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

		a, err := db.GetArticle(context.Background(), 1)
		require.NoError(t, err)
		require.Nil(t, a)
	})
}
