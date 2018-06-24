package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"

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
func booksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(books)
		break
	case "POST":
		w.Header().Set("Content-Type", "application/jsom")
		params := mux.Vars(r)

		for _, item := range books {
			if item.ID == params["id"] {
				json.NewEncoder(w).Encode(item)
				return
			}
		}
		//If book is not found, we reutn an empty book struct
		json.NewEncoder(w).Encode(&Book{})
		break
	default:
		//@TODO - Send a proper error request
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&Book{})
	}
}

func bookHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		w.Header().Set("Content-Type", "application/jsom")
		params := mux.Vars(r)

		for _, item := range books {
			if item.ID == params["id"] {
				json.NewEncoder(w).Encode(item)
				return
			}
		}
		//If book is not found, we reutn an empty book struct
		json.NewEncoder(w).Encode(&Book{})
		break
	case "PUT":
		w.Header().Set("Content-Type", "application/jsom")
		params := mux.Vars(r)
		for index, item := range books {
			if item.ID == params["ID"] {
				//Remove book to update
				books := append(books[:index], books[index+1:]...)
				//Create updated book
				var book Book
				_ = json.NewDecoder(r.Body).Decode(&book)
				book.ID = strconv.Itoa(rand.Intn(10000000)) //Mock ID --> Not safe
				books = append(books, book)
				json.NewEncoder(w).Encode(book)
				return
			}
		}
		json.NewEncoder(w).Encode(books)
		break
	case "DELETE":
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		for index, item := range books {
			if item.ID == params["ID"] {
				//Delete book
				books := append(books[:index], books[index+1:]...)
				json.NewEncoder(w).Encode(books)
				break
			}
		}
		json.NewEncoder(w).Encode(books)
		break
	default:
		//@TODO - Send a proper error request
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&Book{})
	}
}

func main() {
	r := mux.NewRouter()

	//Mock data @todo - implement DB
	books = append(books, Book{ID: "1", Isbn: "12345", Title: "Harry Potter and the Goblet of Fire", Author: &Author{Firstname: "Joanne", Lastname: "Rowling"}})
	books = append(books, Book{ID: "2", Isbn: "44382", Title: "Macbeth", Author: &Author{Firstname: "William", Lastname: "Shakespeare"}})

	//Route hanlders / End points
	r.HandleFunc("/api/books", booksHandler).Methods("GET", "POST")
	r.HandleFunc("/api/book/{id}", bookHandler).Methods("GET", "PUT", "DELETE")

	log.Fatal(http.ListenAndServe(":8000", r))
}
