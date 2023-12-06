package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Article struct {
	Id       uint
	Title    string
	Anons    string
	Fulltext string
}

var posts = []Article{}
var showPost = Article{}

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
func AboutPost(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("template/show.html", "template/header.html", "template/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	vars := mux.Vars(r)
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/golang")
	if err != nil {
		panic(err.Error)
	}

	defer db.Close()

	send, err := db.Query(fmt.Sprintf("SELECT *FROM `articles` WHERE `Id`='%s'", vars["Id"]))
	if err != nil {
		panic(err.Error)
	}
	showPost = Article{}
	for send.Next() {
		var post Article
		err = send.Scan(&post.Id, &post.Title, &post.Anons, &post.Fulltext)
		if err != nil {
			panic(err.Error)
		}

		showPost = post
	}
	t.ExecuteTemplate(w, "show", posts)
}
func HandlePage() {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/", index).Methods("GET")
	rtr.HandleFunc("/create", create).Methods("GET")
	rtr.HandleFunc("/SaveArticle", SaveArticle).Methods("POST")
	rtr.HandleFunc("/post/{Id:[0-9]+}", AboutPost).Methods("GET")

	http.Handle("/", rtr)
	http.ListenAndServe(":8080", nil)
}
func main() {
	HandlePage()
}
