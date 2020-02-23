package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gorilla/mux"
)

// Book is a book
type Book struct {
	ID     string  `json:"id"`
	Isbn   string  `json:"isbn"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`
}

// Author is an author
type Author struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

var books []Book

func main() {
	r := mux.NewRouter()
	initBooks(&books)
	r.HandleFunc("/", home).Methods("GET")
	r.HandleFunc("/api/book", createBook).Methods("POST")
	r.HandleFunc("/api/books", getBooks).Methods("GET")
	r.HandleFunc("/api/book/{id}", getBook).Methods("GET")
	r.HandleFunc("/api/book/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/api/book/{id}", deleteBook).Methods("DELETE")
	fmt.Println("Server running on http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Books API home page"))
	return
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, book := range books {
		if book.ID == params["id"] {
			json.NewEncoder(w).Encode(book)
			break
		}
	}
	w.Write([]byte("No books found!"))
}

func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	var book Book
	// body, _ := ioutil.ReadAll(r.Body)
	// json.Unmarshal(body, &book)
	json.NewDecoder(r.Body).Decode(&book)
	book.ID = strconv.Itoa(rand.Intn(10000000))
	books = append(books, book)
	json.NewEncoder(w).Encode(book)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(204)
	params := mux.Vars(r)
	for i, book := range books {
		if book.ID == params["id"] {
			books = append(books[:i], books[i+1:]...)
			break
		}
	}
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Matching book from the param ID stored in memory
	var originalBook Book

	// Book values from the payload
	var updatedBook Book
	params := mux.Vars(r)

	// Get the book values from the payload
	json.NewDecoder(r.Body).Decode(&updatedBook)

	for i, book := range books {
		if book.ID == params["id"] {
			originalBook = book
			books = append(books[:i], books[i+1:]...)
			break
		}
	}

	// Get all the fields of the struct from the updatedBook
	updatedBookFields := reflect.TypeOf(updatedBook)

	// Get all the values of the struct from the updatedBook
	updatedBookValues := reflect.ValueOf(updatedBook)

	// Create a pointer to the original book
	origValues := reflect.ValueOf(&originalBook)

	// Loop through each updated book field to set the coresponding field on the original book.
	for i := 0; i < updatedBookFields.NumField(); i++ {
		// Get the value for the current field. Ex: Field(i) == "title", then field == "The Book Title"
		field := updatedBookValues.Field(i)

		// If this value is a string, and is not empty, then the value should be updated
		if field.Kind() == reflect.String && field.String() != "" {
			// Use the pointer to update the matching field in memory on the original matching book
			reflect.Indirect(origValues).Field(i).SetString(field.String())
		}

	}
	books = append(books, originalBook)
	json.NewEncoder(w).Encode(originalBook)
}

func initBooks(books *[]Book) {
	*books = append(*books, Book{
		ID:    "1",
		Isbn:  "123",
		Title: "Book one",
		Author: &Author{
			FirstName: "Foo",
			LastName:  "Bar",
		},
	},
		Book{
			ID:    "2",
			Isbn:  "321",
			Title: "Book two",
			Author: &Author{
				FirstName: "Peter",
				LastName:  "Parker",
			},
		},
	)
}
