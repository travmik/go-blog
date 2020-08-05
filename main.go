package main

import (
	"fmt"
	"go-blog/models"
	"html/template"
	"net/http"
)

var posts map[string]*models.Post

func main() {
	fmt.Println("Listening port 3000")

	posts = make(map[string]*models.Post, 0)

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/write", writeHandler)
	http.HandleFunc("/edit", editHandler)
	http.HandleFunc("/delete", deleteHandler)
	http.HandleFunc("/SavePost", savePostHandler)

	http.ListenAndServe("localhost:3000", nil)
}

func indexHandler(writer http.ResponseWriter, request *http.Request) {
	t, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(writer, err.Error())
	}

	fmt.Println(posts)

	t.ExecuteTemplate(writer, "index", posts)
}

func writeHandler(writer http.ResponseWriter, request *http.Request) {
	t, err := template.ParseFiles("templates/write.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(writer, err.Error())
	}
	t.ExecuteTemplate(writer, "write", nil)
}

func editHandler(writer http.ResponseWriter, request *http.Request) {
	t, err := template.ParseFiles("templates/write.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(writer, err.Error())
	}

	id := request.FormValue("id")
	post, found := posts[id]
	if !found {
		http.NotFound(writer, request)
	}

	t.ExecuteTemplate(writer, "write", post)
}

func deleteHandler(writer http.ResponseWriter, request *http.Request) {
	id := request.FormValue("id")
	if id == "" {
		http.NotFound(writer, request)
	}

	delete(posts, id)

	http.Redirect(writer, request, "/", 302)
}

func savePostHandler(writer http.ResponseWriter, request *http.Request) {
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

	http.Redirect(writer, request, "/", 302)
}
