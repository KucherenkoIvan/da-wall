package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"
)

const PORT = 8080

type Post struct {
	text        string
	create_date time.Time
}

var posts []Post = []Post{{"init", time.Date(2000, 2, 1, 12, 30, 0, 0, time.UTC)}}

func renderPosts(w http.ResponseWriter) {
	component := posts_list(posts)
	component.Render(context.Background(), w)
}

func main() {
	// Serve static from /public dir
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	http.HandleFunc("/posts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			renderPosts(w)
			return
		} else if r.Method != http.MethodPost {
			w.WriteHeader(404)
			return
		}

		r.ParseForm()

		postText := r.FormValue("post_text")
		if postText != "" {
			post := Post{postText, time.Now()}
			posts = append(posts, post)

			sort.Slice(posts[:], func(i, j int) bool {
				return !posts[i].create_date.Before(posts[j].create_date)
			})
		}

		renderPosts(w)
	})

	// Root
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		component := root()
		component.Render(context.Background(), w)
	})

	// Start server
	fmt.Printf("Starting on http://localhost:%d\n", PORT)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(PORT), nil))
}
