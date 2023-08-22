package main

type CreatePersonPayload struct {
	Name      string   `json:"nome" validate:"required,max=100"`
	Nickname  string   `json:"apelido" validate:"required,max=32"`
	Birthdate string   `json:"nascimento" validate:"datetime=2006-01-02"`
	Stack     []string `json:"stack" validate:"required,dive,required,max=32"`
}
