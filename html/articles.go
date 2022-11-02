package html

import (
	"fmt"
	"strings"

	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"

	"github.com/maragudk/sqlite-app/model"
)

func HomePage(articles []model.Article) g.Node {
	return Page(PageProps{Title: "Home", Description: "All articles."},

		H1(g.Text("Articles")),

		P(g.Raw(`These are my articles. <a href="/new">Create a new one</a>.`)),

		Ul(
			g.Group(g.Map(articles, func(a model.Article) g.Node {
				return Li(A(Href(fmt.Sprintf("/articles?id=%v", a.ID)), g.Text(a.Title)))
			})),
		),
	)
}

const (
	textInputStyle = "block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
	buttonStyle    = "inline-flex items-center rounded-md border border-transparent bg-indigo-600 px-6 py-3 text-base " +
		"font-medium text-white shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 " +
		"focus:ring-offset-2"
)

func NewArticlePage() g.Node {
	return Page(PageProps{Title: "New article", Description: "Create a new article."},
		H1(g.Text(`New article`)),

		FormEl(Method("post"), Action("/new"), g.Attr("enctype", "multipart/form-data"),
			Div(Class("flex flex-col space-y-4"),

				Label(For("title"), g.Text(`Title`)),
				Input(
					Type("text"), ID("title"), Name("title"), Placeholder("My new article"),
					AutoComplete("off"), Required(), AutoFocus(), Class(textInputStyle),
				),

				Label(For("content"), g.Text(`Content`)),
				Textarea(
					ID("content"), Name("content"), Placeholder("This article is about…"), Required(),
					Class(textInputStyle), Rows("10"),
				),

				Div(Class("mx-auto"),
					Input(Type("submit"), Value("Create"), Class(buttonStyle)),
				),

				P(A(Href("/"), g.Text("Back to front."))),
			),
		),
	)
}

func ArticlePage(a model.Article) g.Node {
	return Page(PageProps{Title: a.Title, Description: ""},
		H1(g.Text(a.Title)),

		P(
			g.Textf(`Published %v`, formatDate(a.Created)),
			g.If(a.Updated.T.After(a.Created.T), g.Textf(`, last updated %v`, formatDate(a.Updated))),
			g.Text(`.`),
		),

		g.Group(g.Map(strings.Split(a.Content, "\n\n"), func(p string) g.Node {
			return P(g.Text(p))
		})),

		P(A(Href("/"), g.Text("Back to front."))),
	)
}

func formatDate(t model.Time) string {
	return t.T.Format("Monday January 2 2006 at 15:04")
}
