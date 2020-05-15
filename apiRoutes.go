package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

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
