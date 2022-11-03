package http

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	g "github.com/maragudk/gomponents"
	ghttp "github.com/maragudk/gomponents/http"

	"github.com/maragudk/sqlite-app/html"
	"github.com/maragudk/sqlite-app/model"
)

type httpError struct {
	Code int
}

func (e httpError) Error() string {
	return http.StatusText(e.Code)
}

func (e httpError) StatusCode() int {
	return e.Code
}

type articlesGetter interface {
	GetTOC(ctx context.Context) ([]model.Article, error)
	SearchArticles(ctx context.Context, search string) ([]model.Article, error)
}

func Home(mux chi.Router, log *log.Logger, db articlesGetter) {
	mux.Get("/", ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (g.Node, error) {
		search := r.URL.Query().Get("search")

		var articles []model.Article
		var err error
		if search != "" {
			articles, err = db.SearchArticles(r.Context(), search)
		} else {
			articles, err = db.GetTOC(r.Context())
		}
		if err != nil {
			log.Println("Error getting/searching articles:", err)
			return html.ErrorPage(), err
		}

		return html.HomePage(articles, search), nil
	}))
}

type articleGetter interface {
	GetArticle(ctx context.Context, id int, search string) (*model.Article, error)
}

func Articles(mux chi.Router, log *log.Logger, db articleGetter) {
	mux.Get("/articles", ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (g.Node, error) {
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			return html.ErrorPage(), httpError{http.StatusBadRequest}
		}

		search := r.URL.Query().Get("search")

		a, err := db.GetArticle(r.Context(), id, search)
		if err != nil {
			log.Println("Error getting article:", err)
			return html.ErrorPage(), err
		}

		if a == nil {
			return html.NotFoundPage(), httpError{http.StatusNotFound}
		}

		return html.ArticlePage(*a, search), nil
	}))
}

type articleCreator interface {
	CreateArticle(ctx context.Context, a model.Article) error
}

func NewArticle(mux chi.Router, log *log.Logger, db articleCreator) {
	mux.Route("/new", func(r chi.Router) {
		r.Get("/", ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (g.Node, error) {
			return html.NewArticlePage(), nil
		}))

		r.Post("/", ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (g.Node, error) {
			a := model.Article{
				Title:   r.PostFormValue("title"),
				Content: normalizeLinebreaks(r.PostFormValue("content")),
			}

			if err := db.CreateArticle(r.Context(), a); err != nil {
				log.Println("Error creating article:", err)
				return html.ErrorPage(), err
			}

			http.Redirect(w, r, "/", http.StatusFound)
			return nil, nil
		}))
	})
}

func normalizeLinebreaks(v string) string {
	return strings.ReplaceAll(v, "\r\n", "\n")
}
