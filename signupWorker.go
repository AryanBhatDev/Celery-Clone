package main

import (
	"context"
	"encoding/json"
	"fmt"

	"log"
	"time"

	"github.com/AryanBhatDev/CeleryClone/internal/database"
	"github.com/AryanBhatDev/CeleryClone/internal/types"
)


func (apiCfg *apiConfig) signupWorker() {
	ctx := context.Background()

	for {
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
		fmt.Println("going to sleep for 20 secs")

		time.Sleep(20*time.Second)

		fmt.Println("awake now")

		hashedPassword := encryptPassword(task.Payload.Password)
		
		user, err := apiCfg.DB.CreateUser(ctx, database.CreateUserParams{
			ID:        task.Payload.ID,
			Name:      task.Payload.Name,
			Email:     task.Payload.Email,
			Password:  hashedPassword,
			CreatedAt: task.Payload.CreatedAt,
			UpdatedAt: task.Payload.UpdatedAt,
		})

		fmt.Println("after db entry")

		if err != nil {
			log.Printf("Error inserting user into DB: %v", err)
			continue
		}

		err = apiCfg.EmailSender.SendWelcomeEmail(user.Email, user.Name)

		if err != nil {
			log.Printf("Failed to send welcome email: %v", err)
		}

		userJson, err := json.Marshal(user)

		if err != nil {
			log.Printf("Error while marshaling : %v", err)
			continue
		}
 
		apiCfg.Redis.Set(ctx, fmt.Sprintf("task_result:%s", task.Payload.ID.String()), userJson, 0)
		apiCfg.Redis.Set(ctx, fmt.Sprintf("task_status:%s", task.Payload.ID.String()), "completed", 0)

	}
}