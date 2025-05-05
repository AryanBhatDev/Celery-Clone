package main

import (
	"context"
	"encoding/json"

	"log"
	"time"

	"github.com/AryanBhatDev/CeleryClone/internal/database"
	"github.com/AryanBhatDev/CeleryClone/internal/types"
)


func (apiCfg *apiConfig) signupWorker() {
	ctx := context.Background()
	log.Print("worker started")
	for range 10000{
		result, err := apiCfg.Redis.BRPop(ctx, 0*time.Second, "user_signup_queue").Result()
		if err != nil {
			log.Printf("Error fetching from Redis queue: %v", err)
			continue
		}

		if len(result) < 2 {
			log.Println("Invalid task data received.")
			continue
		}

		task := types.CreateUserTask{}
		err = json.Unmarshal([]byte(result[1]), &task)
		if err != nil {
			log.Printf("Error unmarshalling task: %v", err)
			continue
		}

		if task.TaskType != "create_user" {
			log.Printf("Unknown task type: %s", task.TaskType)
			continue
		}


		_, err = apiCfg.DB.CreateUser(ctx, database.CreateUserParams{
			ID:        task.Payload.ID,
			Name:      task.Payload.Name,
			Email:     task.Payload.Email,
			Password:  task.Payload.Password,
			CreatedAt: task.Payload.CreatedAt,
			UpdatedAt: task.Payload.UpdatedAt,
		})

		
		if err != nil {
			log.Printf("Error inserting user into DB: %v", err)
			continue
		}

		log.Printf("User %s created successfully.", task.Payload.Email)
	}
}