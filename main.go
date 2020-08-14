package main

import (
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"go-blog/models"
	"net/http"
)

var posts map[string]*models.Post
var counter int

func main() {
	fmt.Println("Listening port 3000")

	posts = make(map[string]*models.Post, 0)
	counter = 0

	m := martini.Classic()

	m.Use(render.Renderer(render.Options{
		Directory: "templates", // Specify what path to load the templates from.
		Layout: "layout", // Specify a layout template. Layouts can call {{ yield }} to render the current template.
		Extensions: []string{".tmpl", ".html"}, // Specify extensions to load for templates.
		Charset: "UTF-8", // Sets encoding for json and html content-types. Default is "UTF-8".
		IndentJSON: true, // Output human readable JSON
	}))

	m.Use(func(r *http.Request) {
		if r.URL.Path == "/write" {
			counter++
		}

	})

	staticOptions := martini.StaticOptions{Prefix: "assets"}
	m.Use(martini.Static("assets", staticOptions))

	m.Get("/", indexHandler)
	m.Get("/write", writeHandler)
	m.Get("/edit/:id", editHandler)
	m.Get("/delete/:id", deleteHandler)
	m.Post("/SavePost", savePostHandler)

	m.Get("/test", func() string {
		return "test"
	})
	m.Run()
}

func indexHandler(rndr render.Render) {
	fmt.Println(counter)

	rndr.HTML(200, "index", posts)
}

func writeHandler(rndr render.Render) {
	rndr.HTML(200, "write", nil)
}

func editHandler(rndr render.Render, params martini.Params) {
	id := params["id"]
	post, found := posts[id]
	if !found {
		rndr.Redirect("/")
		return
	}

	rndr.HTML(200, "write", post)
}

func deleteHandler(rndr render.Render, params martini.Params) {
	id := params["id"]
	if id == "" {
		rndr.Redirect("/")
		return
	}

	delete(posts, id)

	rndr.Redirect("/")
}

func savePostHandler(rndr render.Render, request *http.Request) {
	id := request.FormValue("id")
	title := request.FormValue("title")
	content := request.FormValue("content")

	var post *models.Post
	if id != "" {
		post = posts[id]
		post.Title = title
		post.Content = content
	} else {
		id = GenerateId()
		post := models.NewPost(id, title, content)
		posts[post.Id] = post
	}

	rndr.Redirect("/")
}
