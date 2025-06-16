package main

import (
	"crypto/subtle"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"

	blackfriday "github.com/russross/blackfriday/v2"
)

// simple in-memory store
var (
	posts  = make(map[int]Post)
	nextID = 1
	mu     sync.Mutex
)

type Post struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	expected := os.Getenv("BLOG_TOKEN")
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if len(expected) == 0 || !checkToken(token, expected) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

func checkToken(header, token string) bool {
	const prefix = "Bearer "
	if len(header) <= len(prefix) || header[:len(prefix)] != prefix {
		return false
	}
	provided := header[len(prefix):]
	return subtle.ConstantTimeCompare([]byte(provided), []byte(token)) == 1
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response{Message: "Hello from Go backend"})
}

type response struct {
	Message string `json:"message"`
}

func createPostFromJSON(w http.ResponseWriter, r *http.Request) {
	var p Post
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	mu.Lock()
	p.ID = nextID
	nextID++
	posts[p.ID] = p
	mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func createPostFromNotebook(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Title    string `json:"title"`
		Notebook string `json:"notebook"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	content, err := parseNotebook([]byte(payload.Notebook))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	mu.Lock()
	p := Post{ID: nextID, Title: payload.Title, Content: content}
	nextID++
	posts[p.ID] = p
	mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func listPosts(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	list := make([]Post, 0, len(posts))
	for _, p := range posts {
		list = append(list, p)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func parseNotebook(data []byte) (string, error) {
	var nb struct {
		Cells []struct {
			CellType string   `json:"cell_type"`
			Source   []string `json:"source"`
		} `json:"cells"`
	}
	if err := json.Unmarshal(data, &nb); err != nil {
		return "", err
	}
	var out []byte
	for _, c := range nb.Cells {
		text := ""
		for _, l := range c.Source {
			text += l
		}
		switch c.CellType {
		case "markdown":
			out = append(out, blackfriday.Run([]byte(text))...)
		case "code":
			out = append(out, []byte("<pre><code>")...)
			out = append(out, []byte(text)...)
			out = append(out, []byte("</code></pre>")...)
		}
	}
	return string(out), nil
}

func main() {
	http.HandleFunc("/api/hello", helloHandler)
	http.HandleFunc("/api/posts", authMiddleware(createPostFromJSON))
	http.HandleFunc("/api/posts/ipynb", authMiddleware(createPostFromNotebook))
	http.HandleFunc("/api/posts/list", listPosts)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
