package main

type User struct {
	Email       string `json:"email"`
	ID          int    `json:"id"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}

type Chirp struct {
	AuthorId int    `json:"author_id"`
	Body     string `json:"body"`
	ID       int    `json:"id"`
}
