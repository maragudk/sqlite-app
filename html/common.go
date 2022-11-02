package html

import (
	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	. "github.com/maragudk/gomponents/html"
)

type PageProps struct {
	Title       string
	Description string
}

func Page(p PageProps, body ...g.Node) g.Node {
	return c.HTML5(c.HTML5Props{
		Title:       p.Title,
		Description: p.Description,
		Language:    "en",
		Head: []g.Node{
			Script(Src("https://cdn.tailwindcss.com?plugins=forms,typography")),
		},
		Body: []g.Node{
			Class("bg-gradient-to-b from-white to-teal-50 min-h-screen bg-no-repeat"),
			Container(true,
				Prose(
					g.Group(body),
				),
			),
		},
	})
}

func Container(padY bool, children ...g.Node) g.Node {
	return Div(
		c.Classes{
			"max-w-7xl mx-auto px-4 sm:px-6 lg:px-8": true,
			"py-4 sm:py-6 lg:py-8":                   padY,
		},
		g.Group(children),
	)
}

func Prose(children ...g.Node) g.Node {
	return Div(Class("prose prose-lg lg:prose-xl xl:prose-2xl prose-indigo"), g.Group(children))
}

func ErrorPage() g.Node {
	return Page(PageProps{Title: "Something went wrong", Description: "Oh no! 😵"},
		H1(g.Text("Something went wrong")),
		P(g.Text("Oh no! 😵")),
		P(A(Href("/"), g.Text("Back to front."))),
	)
}

func NotFoundPage() g.Node {
	return Page(PageProps{Title: "There's nothing here! 💨", Description: "Just the void."},
		H1(g.Text("There's nothing here! 💨")),
		P(A(Href("/"), g.Text("Back to front."))),
	)
}
