package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/AryanBhatDev/CeleryClone/internal/database"
	"github.com/google/uuid"
)


func (apiCfg *apiConfig)handlerCreateUser(w http.ResponseWriter, r *http.Request){
	type parameters struct{
		Name string `json:"name"`
		Email string `json:"email"`
		Password string `json:"password"`
	}

	params := parameters{}

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&params)

	if err != nil{
		respondWithError(w,400,fmt.Sprintf("Error while decoding request body:%v",err))
		return 
	}

	user, err := apiCfg.DB.CreateUser(
		r.Context(),
		database.CreateUserParams{
			ID: uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name: params.Name,
			Email: params.Email,
			Password: params.Password,
		},
	)

	if err != nil{
		respondWithError(w, 400, fmt.Sprintf("Error while creating user: %v",err))
		return
	}

	respondWithJson(w, 201, databaseUserToUser(user))
}