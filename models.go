package main

import (
	"github.com/AryanBhatDev/CeleryClone/internal/types"
	"github.com/google/uuid"
)

type Result struct{
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
}

func databaseUserToUser(user types.UserPayload)Result {
	return Result{
		ID: user.ID,
		Name: user.Name,
		Email: user.Email,
	}
}