package html

import (
	"fmt"
	"html/template"
	"net/url"
	"regexp"
	"strings"

	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents-heroicons/solid"
	. "github.com/maragudk/gomponents/html"

	"github.com/maragudk/litefs-app/model"
)

func HomePage(articles []model.Article, search, region string) g.Node {
	return Page(PageProps{Title: "Home", Description: "All articles."},

		H1(g.Text("Articles")),

		P(g.Rawf(`Served to you from the <strong>%v</strong> region.`, region)),

		g.If(len(articles) == 0 && search == "",
			P(g.Raw(`No articles yet. <a href="/new">Create a new one</a>.`)),
		),

		g.If(len(articles) > 0 || search != "",
			Div(
				P(g.Raw(`These are my articles. <a href="/new">Create a new one</a>.`)),

				FormEl(Action("/"), Method("get"), Class("flex items-center w-full"),
					Label(For("search"), Class("sr-only"), g.Text("Search")),
					Div(Class("relative rounded-md shadow-sm flex-grow"),
						Div(Class("absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none"),
							solid.Search(Class("h-5 w-5 text-gray-400")),
						),
						Input(Type("text"), Name("search"), ID("search"), Value(search), TabIndex("1"),
							Class("focus:ring-gray-500 focus:border-gray-500 block w-full pl-10 text-sm border-gray-300 rounded-md")),
					),
				),

				g.If(len(articles) == 0,
					P(g.Raw(`No search results for `), Mark(g.Text(search)), g.Raw(`.`)),
				),

				g.If(len(articles) > 0,
					Ul(
						g.Group(g.Map(articles, func(a model.Article) g.Node {
							return Li(
								P(A(Href(fmt.Sprintf("/articles?id=%v&search=%v", a.ID, url.QueryEscape(search))), highlight(a.Title))),
								g.If(search != "",
									P(highlight(a.Content)),
								),
							)
						})),
					),
				),
			),
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

		FormEl(Method("post"), Action("/new"),
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

func ArticlePage(a model.Article, search string) g.Node {
	return Page(PageProps{Title: a.Title, Description: ""},
		H1(highlight(a.Title)),

		P(
			g.Textf(`Published %v`, formatDate(a.Created)),
			g.If(a.Updated.T.After(a.Created.T), g.Textf(`, last updated %v`, formatDate(a.Updated))),
			g.Text(`.`),
		),

		g.Group(g.Map(strings.Split(a.Content, "\n\n"), func(p string) g.Node {
			return P(highlight(p))
		})),

		P(A(Href("/?search="+url.QueryEscape(search)), g.Text("Back to front."))),
	)
}

func formatDate(t model.Time) string {
	return t.T.Format("Monday January 2 2006 at 15:04")
}

// highlightMatcher matches highlighted text between unit separators, greedily.
// See https://regex101.com/r/cWSkkZ/latest
var highlightMatcher = regexp.MustCompile(`␟(.*?)␟`)

// highlight escapes the given string for HTML output, and then highlights substrings with the HTML mark tag,
// using highlightMatcher.
func highlight(s string) g.Node {
	s = template.HTMLEscapeString(s)
	s = highlightMatcher.ReplaceAllString(s, "<mark>$1</mark>")
	return g.Raw(s)
}
