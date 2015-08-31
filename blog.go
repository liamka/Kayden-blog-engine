package main

import (
	"fmt"
	"net/http"
	"text/template"
	"time"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"./models"
	"io"
	"io/ioutil"
	"os"
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
	stmt, _ := db.Prepare("SELECT * FROM kayden_blog_posts where tags LIKE ?")
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
// Login page
////////////////////////////////
func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t.ExecuteTemplate(w, "header", v)
		t.ExecuteTemplate(w, "login", nil)
	} else {
		r.ParseForm()
		if r.FormValue("pass") == config.Pass {
			expiration := time.Now().Add(365 * 24 * time.Hour)
			cookie := http.Cookie{Name: config.CookieName, Value: config.Pass, Expires: expiration}
			http.SetCookie(w, &cookie)
			http.Redirect(w, r, "/kayden", 302)
		} else {
			http.Redirect(w, r, "/kayden/login", 302)
		}
	}
}

////////////////////////////////
// Admin all posts
////////////////////////////////
func kayden(w http.ResponseWriter, r *http.Request) {
	// Access
	access := CheckCookies(r)
	if !access {
		http.Redirect(w, r, "/kayden/login", 302)
	}

	// Query
	rows, _ := db.Query("SELECT * FROM kayden_blog_posts ORDER BY id DESC LIMIT 100000")
	defer rows.Close()
	posts := []post{}
	for rows.Next() {
		p := post{}
		rows.Scan(&p.Id, &p.Title, &p.Body, &p.Tags, &p.Time)

		posts = append(posts, p)
	}
	t.ExecuteTemplate(w, "header", v)
	t.ExecuteTemplate(w, "posts", posts)
}

////////////////////////////////
// New post
////////////////////////////////
func newPost(w http.ResponseWriter, r *http.Request) {
	// Access
	access := CheckCookies(r)
	if !access {
		http.Redirect(w, r, "/kayden/login", 302)
	}

	t.ExecuteTemplate(w, "header", v)
	t.ExecuteTemplate(w, "new", nil)
}

////////////////////////////////
// Save post
////////////////////////////////
func savePost(w http.ResponseWriter, r *http.Request) {
	// Access
	access := CheckCookies(r)
	if !access {
		http.Redirect(w, r, "/kayden/login", 302)
	}

	// Catch data
	title := r.FormValue("title")
	body := r.FormValue("body")
	tags := r.FormValue("tags")

	now := time.Now()
	time := now.Unix()

	// Prepare and insert
	stmt, _ := db.Prepare("INSERT kayden_blog_posts SET title=?,body=?,tags=?,time=?")
	stmt.Exec(title, body, tags, time)
	// Redirect
	http.Redirect(w, r, "/", 302)
}

////////////////////////////////
// Edit post
////////////////////////////////
func editPost(w http.ResponseWriter, r *http.Request) {
	// Access
	access := CheckCookies(r)
	if !access {
		http.Redirect(w, r, "/kayden/login", 302)
	}

	// Get single
	id := r.URL.Path[len("/kayden/edit/"):]
	stmt, _ := db.Prepare("SELECT * FROM kayden_blog_posts where id = ? LIMIT 1")
	rows, _ := stmt.Query(id)
	defer rows.Close()
	posts := []post{}
	for rows.Next() {
		p := post{}
		rows.Scan(&p.Id, &p.Title, &p.Body, &p.Tags, &p.Time)
		p.Time = timeX(p.Time)
		posts = append(posts, p)
	}
	t.ExecuteTemplate(w, "header", v)
	t.ExecuteTemplate(w, "edit", posts)
}

////////////////////////////////
// Update post
////////////////////////////////
func updatePost(w http.ResponseWriter, r *http.Request) {
	// Access
	access := CheckCookies(r)
	if !access {
		http.Redirect(w, r, "/kayden/login", 302)
	}

	// Catch data
	id := r.FormValue("id")
	title := r.FormValue("title")
	body := r.FormValue("body")
	tags := r.FormValue("tags")

	// Prepare and insert
	stmt, _ := db.Prepare("update kayden_blog_posts SET title=?,body=?,tags=? where id=?")
	stmt.Exec(title, body, tags, id)
	// Redirect
	http.Redirect(w, r, "/kayden", 302)
}

////////////////////////////////
// Delete post
////////////////////////////////
func deletePost(w http.ResponseWriter, r *http.Request) {
	// Access
	access := CheckCookies(r)
	if !access {
		http.Redirect(w, r, "/kayden/login", 302)
	}

	// Get single
	id := r.URL.Path[len("/kayden/delete/"):]
	stmt, _ := db.Prepare("DELETE from kayden_blog_posts where id=?")
	stmt.Exec(id)
	http.Redirect(w, r, "/kayden", 302)
}

////////////////////////////////
// Upload files
////////////////////////////////
func uploads(w http.ResponseWriter, r *http.Request) {
	type file struct {
		Name string
	}

	if r.Method == "GET" {
		filez, _ := ioutil.ReadDir("./uploads/")
		files := []file{}
		for _, fn := range filez {
			f := file{}
			f.Name = fn.Name()
			files = append(files, f)
		}
		t.ExecuteTemplate(w, "header", v)
        t.ExecuteTemplate(w, "upload", files)
	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		f, err := os.OpenFile("./uploads/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
        if err != nil {
            fmt.Println(err)
            return
        }
        defer f.Close()
        io.Copy(f, file)
        http.Redirect(w, r, "/kayden/upload", 302)
	}
}

////////////////////////////////
// MAIN
////////////////////////////////
func main() {
	// Sample
	log("Kayden love it!")
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
	
	// Get port
	http.ListenAndServe(":3000", nil)
}