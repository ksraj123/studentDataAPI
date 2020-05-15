package main

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
