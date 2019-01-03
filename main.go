package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/gorilla/mux"
)

func main() {
	db, err := NewFruitsRepository("postgres", "host=localhost port=5432 user=postgres dbname=postgres password=mysecretpassword sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.AutoMigrate(&Fruit{})

	fruitController := NewFruitController(db)

	r := mux.NewRouter()
	r.HandleFunc("/api/fruits", fruitController.List).Methods("GET")
	r.HandleFunc("/api/fruits/{id:[0-9]+}", fruitController.Show).Methods("GET")
	r.HandleFunc("/api/fruits", fruitController.Create).Methods("POST")
	r.HandleFunc("/api/fruits/{id:[0-9]+}", fruitController.Update).Methods("PUT")
	r.HandleFunc("/api/fruits/{id:[0-9]+}", fruitController.Delete).Methods("DELETE")

	http.Handle("/", r)
	r.Use(loggingMiddleware)
	fmt.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("\n[%s] %q %q",
			time.Now().Format("02/Jan/2006:15:04:05 -0700"),
			r.Method,
			r.RequestURI)
		next.ServeHTTP(w, r)

	})
}
