package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

//Book struct (model)
type Book struct {
	ID     string  `json:"id"`
	Isbn   string  `json:"isbn"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`
}

//Author struct
type Author struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

// Init books var as slice book struct
var books []Book

//Handlers
func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {

}

func createBook(w http.ResponseWriter, r *http.Request) {

}

func updateBook(w http.ResponseWriter, r *http.Request) {

}

func deleteBook(w http.ResponseWriter, r *http.Request) {

}
func main() {
	r := mux.NewRouter()

	//Mock data @todo - implement DB
	books = append(books, Book{ID: "1", Isbn: "12345", Title: "Harry Potter and the Goblet of Fire", Author: &Author{Firstname: "Joanne", Lastname: "Rowling"}})
	books = append(books, Book{ID: "2", Isbn: "44382", Title: "Macbeth", Author: &Author{Firstname: "William", Lastname: "Shakespeare"}})

	//Route hanlders / End points
	r.HandleFunc("/api/books", getBooks).Methods("GET")
	r.HandleFunc("/api/book/{id}", getBook).Methods("GET")
	r.HandleFunc("/api/books", createBook).Methods("POST")
	r.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", r))
}
