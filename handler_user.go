package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"github.com/AryanBhatDev/CeleryClone/internal/types"
	"github.com/google/uuid"
)


func (apiCfg *apiConfig)handlerPushCreateUser(w http.ResponseWriter, r *http.Request){
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


	_, err = apiCfg.DB.GetUser(r.Context(),params.Email)

	if err == nil{
		respondWithError(w,409,fmt.Sprintf("Email is already taken:%v",err))
		return
	}

	task := types.CreateUserTask{
		TaskType: "create_user",
		Payload: types.UserPayload{
			ID: uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name: params.Name,
			Email: params.Email,
			Password: params.Password,
		},
	}

	taskJson, err := json.Marshal(task)

	if err != nil{
		respondWithError(w, 400, fmt.Sprintf("Error while creating task: %v",err))
		return
	}

	ctx := context.Background()

	err = apiCfg.Redis.LPush(ctx,"user_signup_queue",taskJson).Err()
	log.Println("after push")
	if err != nil{
		respondWithError(w, 500, fmt.Sprintf("Failed to push to queue: %v",err))
		return
	}

	respondWithJson(w, 201, "Signed up")
}