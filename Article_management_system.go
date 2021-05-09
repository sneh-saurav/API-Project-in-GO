package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Article struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

var ArticleManagement []Article

func getAllArticles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ArticleManagement)
}
func createArticle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var article Article
	_ = json.NewDecoder(r.Body).Decode(&article)
	article.ID = strconv.Itoa(rand.Intn(100))
	ArticleManagement = append(ArticleManagement, article)
	json.NewEncoder(w).Encode(&article)
}
func getArticle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range ArticleManagement {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Article{})
}
func updateArticle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range ArticleManagement {
		if item.ID == params["id"] {
			ArticleManagement = append(ArticleManagement[:index], ArticleManagement[index+1:]...)
			var article Article
			_ = json.NewDecoder(r.Body).Decode(&article)
			article.ID = params["id"]
			ArticleManagement = append(ArticleManagement, article)
			json.NewEncoder(w).Encode(&article)
			return
		}
	}
	json.NewEncoder(w).Encode(ArticleManagement)
}
func deleteArticle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range ArticleManagement {
		if item.ID == params["id"] {
			ArticleManagement = append(ArticleManagement[:index], ArticleManagement[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(ArticleManagement)
}
func main() {
	router := mux.NewRouter()
	ArticleManagement = append(ArticleManagement, Article{ID: "1", Title: "My first article", Body: "This is the content of my first article"})
	router.HandleFunc("/articles", getAllArticles).Methods("GET")
	router.HandleFunc("/articles", createArticle).Methods("POST")
	router.HandleFunc("/articles/{id}", getArticle).Methods("GET")
	router.HandleFunc("/articles/{id}", updateArticle).Methods("PUT")
	router.HandleFunc("/articles/{id}", deleteArticle).Methods("DELETE")
	http.ListenAndServe(":8000", router)
}
