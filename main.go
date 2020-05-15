package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

var studentColl *mongo.Collection

func main() {
	// Connecting to database
	connectToDb()

	// Init router
	r := mux.NewRouter()
	r.HandleFunc("/api/students", getStudents).Methods("GET")
	r.HandleFunc("/api/students/{roll}", getStudent).Methods("GET")
	r.HandleFunc("/api/students", createStudent).Methods("POST")
	r.HandleFunc("/api/students/{roll}", updateStudent).Methods("PUT")
	r.HandleFunc("/api/students/{roll}", deleteStudent).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8001", r))
}
