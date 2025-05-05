package types

import (
	"time"
	"github.com/google/uuid"
)

type CreateUserTask struct {
	TaskType string      `json:"task_type"`
	Payload  UserPayload `json:"payload"`
}

type UserPayload struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
}
