package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Student Struct (Model)
type Student struct {
	Roll   string  `json:"roll" bson:"roll,omitempty"`
	Name   string  `json:"name" bson:"name,omitempty"`
	Branch string  `json:"branch" bson:"branch,omitempty"`
	Parent *Parent `json:"parent" bson:"parent,omitempty"`
}

// Parent Struct
type Parent struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var studentColl *mongo.Collection

// Get all Students
func getStudents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	var results []*Student
	cur, err := studentColl.Find(ctx, bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	for cur.Next(ctx) {
		var elem Student
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, &elem)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	cur.Close(ctx)
	json.NewEncoder(w).Encode(results)
}

// Get Single Student
func getStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Get params
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	var result Student
	err := studentColl.FindOne(ctx, bson.M{"roll": params["roll"]}).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(result)
	// Handle error when wrong roll number provided
}

// create a new Student
func createStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var student Student
	json.NewDecoder(r.Body).Decode(&student)
	student.Roll = strconv.Itoa(rand.Intn(1000000))
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	insertResult, err := studentColl.InsertOne(ctx, student)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(insertResult)
}

// Udate a Student
func updateStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Get params
	var student Student
	json.NewDecoder(r.Body).Decode(&student)
	student.Roll = params["roll"] // Roll is internal and cannot be changed
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	_, err := studentColl.DeleteMany(ctx, bson.M{"roll": params["roll"]})
	if err != nil {
		log.Fatal(err)
	}
	insertResult, err := studentColl.InsertOne(ctx, student)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(insertResult)
}

// Delete a Student
func deleteStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Get params
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	deleteResult, err := studentColl.DeleteMany(ctx, bson.M{"roll": params["roll"]})
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(deleteResult)
	// Handle error when wrong roll number provided
}

func main() {
	// Connecting to database

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
	studentColl = client.Database("test").Collection("students")

	// Init router
	r := mux.NewRouter()
	r.HandleFunc("/api/students", getStudents).Methods("GET")
	r.HandleFunc("/api/students/{roll}", getStudent).Methods("GET")
	r.HandleFunc("/api/students", createStudent).Methods("POST")
	r.HandleFunc("/api/students/{roll}", updateStudent).Methods("PUT")
	r.HandleFunc("/api/students/{roll}", deleteStudent).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8001", r))
}
