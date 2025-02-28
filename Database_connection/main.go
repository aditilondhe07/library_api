package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/aditilonde07/libraryManagement/ent"
	_ "github.com/lib/pq"
)

func main() {
	client, err := ent.Open("postgres", "host=localhost dbname=library user=your_username password=your_password sslmode=disable")
	if err != nil {
		log.Fatalf("failed connecting to postgres: %v", err)
	}
	defer client.Close()

	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	http.HandleFunc("/books", createBook(client))
	http.HandleFunc("/users", createUser(client))
	http.HandleFunc("/racks", createRack(client))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func createBook(client *ent.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var b ent.Book
		if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		newBook, err := client.Book.Create().
			SetTitle(b.Title).
			SetAuthor(b.Author).
			SetIsbn(b.ISBN).
			Save(context.Background())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(newBook)
	}
}

func createUser(client *ent.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u ent.User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		newUser, err := client.User.Create().
			SetName(u.Name).
			SetEmail(u.Email).
			Save(context.Background())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(newUser)
	}
}

func createRack(client *ent.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var r ent.Rack
		if err := json.NewDecoder(r.Body).Decode(&r); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		newRack, err := client.Rack.Create().
			SetLocation(r.Location).
			Save(context.Background())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(newRack)
	}
}
