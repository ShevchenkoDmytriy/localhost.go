package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Article struct {
	Id       uint
	Title    string
	Anons    string
	Fulltext string
}

var posts = []Article{}

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("template/index.html", "template/header.html", "template/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/golang")
	if err != nil {
		panic(err.Error)
	}

	defer db.Close()

	send, err := db.Query("SELECT *FROM `articles`")
	if err != nil {
		panic(err.Error)
	}
	posts = []Article{}
	for send.Next() {
		var post Article
		err = send.Scan(&post.Id, &post.Title, &post.Anons, &post.Fulltext)
		if err != nil {
			panic(err.Error)
		}

		posts = append(posts, post)
	}
	t.ExecuteTemplate(w, "index", posts)
}
func create(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("template/create.html", "template/header.html", "template/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	t.ExecuteTemplate(w, "create", nil)
}
func SaveArticle(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	anons := r.FormValue("anons")
	fulltext := r.FormValue("FullText")

	if title == "" || anons == "" || fulltext == "" {
		fmt.Fprintf(w, "Try again")
	} else {
		db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/golang")
		if err != nil {
			panic(err.Error)
		}

		defer db.Close()

		insert, err := db.Query(fmt.Sprintf("INSERT INTO `articles`(`title`,`anons`,`fulltext`) VALUES('%s','%s','%s')", title, anons, fulltext))
		if err != nil {
			panic(err.Error)
		}

		defer insert.Close()

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

}
func HandlePage() {
	http.HandleFunc("/", index)
	http.HandleFunc("/create", create)
	http.HandleFunc("/SaveArticle", SaveArticle)
	http.ListenAndServe(":8080", nil)
}
func main() {
	HandlePage()
}
