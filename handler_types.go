package main

type User struct {
	Email string `json:"email"`
	ID    int    `json:"id"`
}

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}
