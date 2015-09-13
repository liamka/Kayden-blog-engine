package main

import (
	"net/http"
	"fmt"
	"time"
	"io"
	"io/ioutil"
	"os"
)


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
// Delete uploaded file
////////////////////////////////
func uploadsDelete(w http.ResponseWriter, r *http.Request) {
	// Access
	access := CheckCookies(r)
	if !access {
		http.Redirect(w, r, "/kayden/login", 302)
	}

	// Get file name
	file := r.URL.Path[len("/upload/upload/delete/"):]

	err := os.Remove("./uploads/"+file)
	if err != nil {
        fmt.Println(err)
        return
    }

    http.Redirect(w, r, "/kayden/upload", 302)
}

////////////////////////////////
// All drafts
////////////////////////////////
func drafts(w http.ResponseWriter, r *http.Request) {
	// Access
	access := CheckCookies(r)
	if !access {
		http.Redirect(w, r, "/kayden/login", 302)
	}

	// Query
	rows, _ := db.Query("SELECT * FROM kayden_blog_drafts ORDER BY id DESC LIMIT 100000")
	defer rows.Close()
	posts := []post{}
	for rows.Next() {
		p := post{}
		rows.Scan(&p.Id, &p.Title, &p.Body, &p.Tags, &p.Time)

		posts = append(posts, p)
	}
	t.ExecuteTemplate(w, "header", v)
	t.ExecuteTemplate(w, "drafts", posts)
}

////////////////////////////////
// New draft
////////////////////////////////
func newDraft(w http.ResponseWriter, r *http.Request) {
	// Access
	access := CheckCookies(r)
	if !access {
		http.Redirect(w, r, "/kayden/login", 302)
	}

	t.ExecuteTemplate(w, "header", v)
	t.ExecuteTemplate(w, "draft_new", nil)
}

////////////////////////////////
// Save draft
////////////////////////////////
func saveDraft(w http.ResponseWriter, r *http.Request) {
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
	stmt, _ := db.Prepare("INSERT kayden_blog_drafts SET title=?,body=?,tags=?,time=?")
	stmt.Exec(title, body, tags, time)
	// Redirect
	http.Redirect(w, r, "/kayden/drafts", 302)
}

////////////////////////////////
// Edit draft
////////////////////////////////
func editDraft(w http.ResponseWriter, r *http.Request) {
	// Access
	access := CheckCookies(r)
	if !access {
		http.Redirect(w, r, "/kayden/login", 302)
	}

	// Get single
	id := r.URL.Path[len("/kayden/drafts/edit/"):]
	stmt, _ := db.Prepare("SELECT * FROM kayden_blog_drafts where id = ? LIMIT 1")
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
	t.ExecuteTemplate(w, "draft_edit", posts)
}

////////////////////////////////
// Update draft
////////////////////////////////
func updateDraft(w http.ResponseWriter, r *http.Request) {
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
	stmt, _ := db.Prepare("update kayden_blog_drafts SET title=?,body=?,tags=? where id=?")
	stmt.Exec(title, body, tags, id)
	// Redirect
	http.Redirect(w, r, "/kayden/drafts", 302)
}

////////////////////////////////
// Delete draft
////////////////////////////////
func deleteDraft(w http.ResponseWriter, r *http.Request) {
	// Access
	access := CheckCookies(r)
	if !access {
		http.Redirect(w, r, "/kayden/login", 302)
	}

	// Get single
	id := r.URL.Path[len("/kayden/drafts/delete/"):]
	stmt, _ := db.Prepare("DELETE from kayden_blog_drafts where id=?")
	stmt.Exec(id)
	http.Redirect(w, r, "/kayden/drafts", 302)
}

////////////////////////////////
// Publish draft
////////////////////////////////
func publishDraft(w http.ResponseWriter, r *http.Request) {
	// Access
	access := CheckCookies(r)
	if !access {
		http.Redirect(w, r, "/kayden/login", 302)
	}

	id := r.URL.Path[len("/kayden/drafts/publish/"):]

	stmt, _ := db.Prepare("SELECT * FROM kayden_blog_drafts where id = ? LIMIT 1")
	rows, _ := stmt.Query(id)
	defer rows.Close()
	posts := []post{}
	for rows.Next() {
		p := post{}
		rows.Scan(&p.Id, &p.Title, &p.Body, &p.Tags, &p.Time)
		p.Time = timeX(p.Time)
		posts = append(posts, p)
	}

	stmta, _ := db.Prepare("INSERT kayden_blog_posts SET title=?,body=?,tags=?,time=?")
	stmta.Exec(posts[0].Title, posts[0].Body, posts[0].Tags, posts[0].Time)

	stmtd, _ := db.Prepare("DELETE from kayden_blog_drafts where id=?")
	stmtd.Exec(id)

	http.Redirect(w, r, "/kayden", 302)
}