package main

import (
	"net/http"
	"text/template"
	_ "github.com/go-sql-driver/mysql"
	"github.com/liamka/Superior"
	"database/sql"
	"./models"
)

var (
	db *sql.DB
	config models.Config
	v models.Vars
	t = template.Must(template.ParseGlob("views/*"))
)

type post struct {
	Id 			int
	Title      	string
	Body       	string
	Tags       	string
	Time       	string
}

////////////////////////////////
// Index post
////////////////////////////////
func index(w http.ResponseWriter, r *http.Request) {
	// Query
	rows, _ := db.Query("SELECT * FROM kayden_blog_posts ORDER BY id DESC LIMIT 10")
	defer rows.Close()
	posts := []post{}
	for rows.Next() {
		p := post{}
		rows.Scan(&p.Id, &p.Title, &p.Body, &p.Tags, &p.Time)
		v.Title = config.Title
		p.Body = ConvertMarkdownToHtml(p.Body)
		p.Time = timeX(p.Time)
		p.Tags = tagsX(p.Tags,"true")
		posts = append(posts, p)
	}
	t.ExecuteTemplate(w, "header", v)
	t.ExecuteTemplate(w, "index", posts)
}

////////////////////////////////
// Single post
////////////////////////////////
func single(w http.ResponseWriter, r *http.Request) {
	// Get single
	id := r.URL.Path[len("/note/"):]
	stmt, _ := db.Prepare("SELECT * FROM kayden_blog_posts where id = ? LIMIT 1")
	rows, _ := stmt.Query(id)
	defer rows.Close()
	posts := []post{}
	for rows.Next() {
		p := post{}
		rows.Scan(&p.Id, &p.Title, &p.Body, &p.Tags, &p.Time)
		v.Title = config.Title + ": " + p.Title
		p.Body = ConvertMarkdownToHtml(p.Body)
		p.Time = timeX(p.Time)
		p.Tags = tagsX(p.Tags,"true")
		posts = append(posts, p)
	}
	// 404 
	if len(posts) == 0 {
		http.Redirect(w, r, "/404", 302)
	}
	t.ExecuteTemplate(w, "header", v)
	t.ExecuteTemplate(w, "index", posts)
}
////////////////////////////////
// All posts
////////////////////////////////
func allPosts(w http.ResponseWriter, r *http.Request) {
	// Query
	rows, _ := db.Query("SELECT * FROM kayden_blog_posts ORDER BY id DESC LIMIT 100000")
	defer rows.Close()
	posts := []post{}
	for rows.Next() {
		p := post{}
		rows.Scan(&p.Id, &p.Title, &p.Body, &p.Tags, &p.Time)
		p.Body = ConvertMarkdownToHtml(p.Body)
		p.Time = timeX(p.Time)
		p.Tags = tagsX(p.Tags,"false")
		posts = append(posts, p)
	}
	t.ExecuteTemplate(w, "header", v)
	t.ExecuteTemplate(w, "all", posts)
}
////////////////////////////////
// Tags
////////////////////////////////
func tagPosts(w http.ResponseWriter, r *http.Request) {
	// Get tag
	tag := r.URL.Path[len("/tag/"):]
	stmt, _ := db.Prepare("SELECT * FROM kayden_blog_posts where tags LIKE ? ORDER BY id DESC LIMIT 1000000")
	rows, _ := stmt.Query("%"+tag+"%")
	defer rows.Close()
	posts := []post{}
	for rows.Next() {
		p := post{}
		rows.Scan(&p.Id, &p.Title, &p.Body, &p.Tags, &p.Time)
		p.Body = ConvertMarkdownToHtml(p.Body)
		p.Time = timeX(p.Time)
		p.Tags = tagsX(p.Tags,"true")
		posts = append(posts, p)
	}
	t.ExecuteTemplate(w, "header", v)
	t.ExecuteTemplate(w, "index", posts)
}

////////////////////////////////
// RSS
////////////////////////////////
func rssPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/xml")
	// Query
	rows, _ := db.Query("SELECT * FROM kayden_blog_posts ORDER BY id DESC LIMIT 10")
	defer rows.Close()

	type rssPost struct {
		Id 			int
		Title      	string
		Body       	string
		Tags       	string
		Time       	string
		URI       	string
	}

	posts := []rssPost{}
	k := 0
	var lastUpdate string
	for rows.Next() {
		p := rssPost{}
		rows.Scan(&p.Id, &p.Title, &p.Body, &p.Tags, &p.Time)
			if k == 0 {
				lastUpdate = timeRFC(p.Time)
			}
		p.Body = ConvertMarkdownToHtml(p.Body)
		p.Time = timeRFC(p.Time)
		p.Tags = tagsX(p.Tags,"false")
		p.URI = config.URI
		posts = append(posts, p)
		k++
	}

	type Rss struct {
		Title   	string
		Subtitle 	string
		URI 		string
		Description string
		LastBuild 	string
	}

	rss := Rss{
		Title:   		config.Title,
		URI:  			config.URI,
		Description:   	config.Description,
		LastBuild:   	lastUpdate,
	}

	t.ExecuteTemplate(w, "rss_header", rss)
	t.ExecuteTemplate(w, "rss", posts)
}

////////////////////////////////
// 404
////////////////////////////////
func notfound(w http.ResponseWriter, r *http.Request) {
	t.ExecuteTemplate(w, "header", v)
	t.ExecuteTemplate(w, "404", nil)
}

////////////////////////////////
// MAIN
////////////////////////////////
func main() {
	// Sample
	Superior.Print("Kayden love it!", "normal", "green")
	// Conf
	config = models.Conf()
	v = models.Values(config)
	// Open mysql
	db, _ = sql.Open("mysql", config.Mysql)
	// Invoke folders
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))
	// Routes
	http.HandleFunc("/", index)
	http.HandleFunc("/note/", single)
	http.HandleFunc("/all/", allPosts)
	http.HandleFunc("/tag/", tagPosts)
	http.HandleFunc("/rss/", rssPosts)
	http.HandleFunc("/404/", notfound)
	http.HandleFunc("/kayden/login", login)
	http.HandleFunc("/kayden/", kayden)
	http.HandleFunc("/kayden/new/", newPost)
	http.HandleFunc("/kayden/new/save", savePost)
	http.HandleFunc("/kayden/edit/", editPost)
	http.HandleFunc("/kayden/edit/save", updatePost)
	http.HandleFunc("/kayden/delete/", deletePost)
	http.HandleFunc("/kayden/upload", uploads)
	http.HandleFunc("/kayden/upload/delete/", uploadsDelete)
	http.HandleFunc("/kayden/drafts/", drafts)
	http.HandleFunc("/kayden/drafts/new/", newDraft)
	http.HandleFunc("/kayden/drafts/new/save", saveDraft)
	http.HandleFunc("/kayden/drafts/edit/", editDraft)
	http.HandleFunc("/kayden/drafts/edit/save", updateDraft)
	http.HandleFunc("/kayden/drafts/delete/", deleteDraft)
	http.HandleFunc("/kayden/drafts/publish/", publishDraft)
	
	// Get port
	http.ListenAndServe(":3000", nil)
}