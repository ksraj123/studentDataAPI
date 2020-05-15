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

// Local datastore save all students inside a slice named students which contain vlues of type Student
var students []Student

var studentColl *mongo.Collection

// Get all Students
func getStudents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}

// Get Single Student
func getStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Get params
	for _, item := range students {
		if item.Roll == params["roll"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Student{}) // return empty student when not found
}

// create a new Student
func createStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var student Student
	_ = json.NewDecoder(r.Body).Decode(&student)
	student.Roll = strconv.Itoa(rand.Intn(1000000))
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	insertResult, err := studentColl.InsertOne(ctx, student)
	if err != nil {
		log.Fatal(err)
	}

	// students = append(students, student)
	json.NewEncoder(w).Encode(insertResult)
}

// Udate a Student
func updateStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Get params
	var student Student
	_ = json.NewDecoder(r.Body).Decode(&student)
	for index, item := range students {
		if item.Roll == params["roll"] {
			student.Roll = params["roll"] // Roll is internal and cannot be changed
			students = append(append(students[:index], student), students[index+1:]...)
			log.Println(item)
			break
		}
	}
	json.NewEncoder(w).Encode(students)
}

// Delete a Student
func deleteStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Get params
	for index, item := range students {
		if item.Roll == params["roll"] {
			students = append(students[:index], students[index+1:]...)
			log.Println(item)
			break
		}
	}
	json.NewEncoder(w).Encode(students)
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
	// fmt.Printf("type of ctx variable = %T\n", ctx)

	// insertResult, err := studentColl.InsertOne(ctx, Student{Roll: "1", Branch: "123456", Name: "Student 1",
	// 	Parent: &Parent{Firstname: "saurabh", Lastname: "raj"}})
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("Inserted a single document: ", insertResult.InsertedID)

	// Init router
	r := mux.NewRouter()

	students = append(students, Student{Roll: "1", Branch: "123456", Name: "Student 1",
		Parent: &Parent{Firstname: "saurabh", Lastname: "raj"}})

	students = append(students, Student{Roll: "2", Branch: "789456", Name: "Student 2",
		Parent: &Parent{Firstname: "rahul", Lastname: "singh"}})

	// Route Handler - will establish endpoints for our APIs
	// route then the function to be executed when that route with given method is hit
	r.HandleFunc("/api/students", getStudents).Methods("GET")
	r.HandleFunc("/api/students/{roll}", getStudent).Methods("GET")
	r.HandleFunc("/api/students", createStudent).Methods("POST")
	r.HandleFunc("/api/students/{roll}", updateStudent).Methods("PUT")
	r.HandleFunc("/api/students/{roll}", deleteStudent).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8001", r))
}
