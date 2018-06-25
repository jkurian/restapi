package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"

	"database/sql"

	_ "github.com/lib/pq"
)

// DB connection
const (
	dbUser     = "jerrykurian"
	dbPassword = "jerrykurian"
	dbName     = "restapi_test"
)

var dbinfo = fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", dbUser, dbPassword, dbName)
var db, err = sql.Open("postgres", dbinfo)

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

//@TODO - Look at all / update queries
//Handlers
func booksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		w.Header().Set("Content-Type", "application/json")

		var books []Book

		rows, err := db.Query("SELECT b.id, b.isbn, b.title, a.firstname, a.lastname FROM books as b INNER JOIN authors a on b.author = a.id;")
		fmt.Print(rows)
		if err != nil {
			log.Fatal(err.Error())
		}
		for rows.Next() {
			var book Book
			var author Author
			rows.Scan(&book.ID, &book.Isbn, &book.Title, &author.Firstname, &author.Lastname)
			book.Author = &author
			books = append(books, book)
		}
		errEncode := json.NewEncoder(w).Encode(books)
		if errEncode != nil {
			log.Fatal(errEncode.Error())
		}
		break
	case "POST":
		w.Header().Set("Content-Type", "application/json")

		var book Book
		_ = json.NewDecoder(r.Body).Decode(&book)
		var authorID int
		//@TODO - Handle error here / update query
		db.QueryRow("INSERT INTO authors (firstname, lastname) VALUES ($1, $2) ON CONFLICT(firstname, lastname) DO UPDATE SET firstname=EXCLUDED.firstname RETURNING id;", book.Author.Firstname, book.Author.Lastname).Scan(&authorID)

		fmt.Print(authorID)
		_, errInsertBook := db.Exec("INSERT INTO books (isbn, title, author) VALUES ($1, $2, $3);", book.Isbn, book.Title, authorID)
		if errInsertBook != nil {
			log.Fatal(errInsertBook.Error())
		}

		errEncode := json.NewEncoder(w).Encode(book)
		if errEncode != nil {
			log.Fatal(errEncode.Error())
		}
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
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		if err != nil {
			log.Fatal(err.Error())
		}

		row := db.QueryRow("SELECT b.id, b.isbn, b.title, a.firstname, a.lastname FROM books as b INNER JOIN authors a on b.author = a.id WHERE b.id = $1;", params["id"])

		var book Book
		var author Author

		err := row.Scan(&book.ID, &book.Isbn, &book.Title, &author.Firstname, &author.Lastname)
		if err != nil {
			log.Fatal(err.Error())
		}

		book.Author = &author
		errWrite := json.NewEncoder(w).Encode(book)

		if errWrite != nil {
			log.Fatal(errWrite.Error())
			json.NewEncoder(w).Encode(&Book{})
		}
		break
	case "PUT":
		w.Header().Set("Content-Type", "application/jsom")
		params := mux.Vars(r)

		row := db.QueryRow("SELECT b.id, b.isbn, b.title, a.firstname, a.lastname FROM books as b WHERE b.id = $1;", params["id"])

		var book Book
		var authorID int
		err := row.Scan(&book.ID, &book.Isbn, &book.Title, &authorID)
		if err != nil {
			log.Fatal(err.Error())
		}

		_, errInsertBook := db.Exec("UPDATE books (isbn, title, author) VALUES ($1, $2, $3) WHERE id = $4;", book.Isbn, book.Title, authorID, params["id"])
		if errInsertBook != nil {
			log.Fatal(errInsertBook.Error())
			json.NewEncoder(w).Encode(&Book{})
		}

		errWrite := json.NewEncoder(w).Encode(book)
		if errWrite != nil {
			log.Fatal(errWrite.Error())
		}
		break
	case "DELETE":
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)

		_, err := db.Query("DELETE FROM books WHERE id = $1;", params["id"])
		if err != nil {
			log.Fatal(err.Error())
		}

		//@TODO - Send back proper response on delete
		json.NewEncoder(w).Encode(&Book{})
		break
	default:
		//@TODO - Send a proper error request
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&Book{})
	}
}

func main() {
	if err != nil {
		log.Fatal(err)
	}

	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	r := mux.NewRouter()

	//Route hanlders / End points
	r.HandleFunc("/api/books", booksHandler).Methods("GET", "POST")
	r.HandleFunc("/api/book/{id}", bookHandler).Methods("GET", "PUT", "DELETE")

	srv := &http.Server{
		Addr: "0.0.0.0:8000",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}
