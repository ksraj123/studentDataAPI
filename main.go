package main

import (
	//core package
	"encoding/json"
	"fmt"
	"log"       // core package - logs erros etc
	"math/rand" // corepackage
	"net/http"  // core package - used to work with http
	"strconv"

	"github.com/gorilla/mux"
)

// Student Struct (Model)
type Student struct {
	Roll   string  `json:"roll"`
	Name   string  `json:"name"`
	Branch string  `json:"branch"`
	Parent *Parent `json:"parent"`
}

// Parent Struct
type Parent struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

// Local datastore save all students inside a slice named students which contain vlues of type Student
var students []Student

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
	students = append(students, student)
	json.NewEncoder(w).Encode(student)
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
			fmt.Println(item)
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
			fmt.Println(item)
			break
		}
	}
	json.NewEncoder(w).Encode(students)
}

func main() {
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
