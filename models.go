package main

import "github.com/AryanBhatDev/CeleryClone/internal/database"




type User struct{
	Name string `json:"name"`
	Email string `json:"email"`
}

func databaseUserToUser(dbUser database.User) User{
	return User{
		Name: dbUser.Name,
		Email: dbUser.Email,
	}
}