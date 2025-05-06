package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/AryanBhatDev/CeleryClone/internal/types"
	"github.com/asaskevich/govalidator"

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


	if !govalidator.StringLength(params.Name, "3", "20") {
		respondWithError(w,400,fmt.Sprintf("Name should be minimum 3 characters and 20 characters:%v",err))
		return
	}
	if !govalidator.IsEmail(params.Email) {
		respondWithError(w,400,fmt.Sprintf("Invalid email format:%v",err))
		return
	}
	if len(params.Password) < 8 {
		respondWithError(w,400,fmt.Sprintf("Password should be more than 8 characters long:%v",err))
		return
	}

	if err != nil{
		respondWithError(w,401,fmt.Sprintf("Incorrect paramss:%v",err))
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

	taskId := task.Payload.ID.String()

	taskJson, err := json.Marshal(task)

	if err != nil{
		respondWithError(w, 400, fmt.Sprintf("Error while creating task: %v",err))
		return
	}

	ctx := context.Background()

	err = apiCfg.Redis.LPush(ctx,"user_signup_queue",taskJson).Err()

	if err != nil{
		respondWithError(w, 500, fmt.Sprintf("Failed to push to queue: %v",err))
		return
	}
	apiCfg.Redis.Set(ctx,fmt.Sprintf("task_status:%s",taskId),"pending",0)

	type SignupResponse struct{
		Message string `json:"message"`
		TaskId string   `json:"task_id"`	
	}

	respondWithJson(w, 201, SignupResponse{
		Message: "Signed up",
		TaskId: taskId,
	})
}

func (apiCfg *apiConfig) handlerTaskStatus(w http.ResponseWriter, r *http.Request){
	ctx := context.Background()

	taskId := r.URL.Query().Get("task_id")

	if taskId == "" {
		respondWithError(w, 400, "Missing task_id in query parameters")
		return
	}

	status, err := apiCfg.Redis.Get(ctx,fmt.Sprintf("task_status:%s",taskId)).Result()

	if err != nil {
		respondWithError(w, 404, fmt.Sprintf("Error while fetching task status: %v",err))
		return
	}


	result := types.UserPayload{}

	if status == "completed"{
		jsonStr, err := apiCfg.Redis.Get(ctx,fmt.Sprintf("task_result:%s",taskId)).Result()

		if err != nil {
			respondWithError(w, 404, fmt.Sprintf("Error while fetching task result: %v",err))
			return
		}
		
		json.Unmarshal([]byte(jsonStr), &result)

		respondWithJson(w, 200,databaseUserToUser(result))
	}else{
		respondWithJson(w,200,fmt.Sprintf("status:%v",status))
	}
	
}