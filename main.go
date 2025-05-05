package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/AryanBhatDev/CeleryClone/internal/database"
	_ "github.com/lib/pq"
)


type apiConfig struct{
	DB *database.Queries
}


func main(){

	err := godotenv.Load()

	if err!= nil{
		log.Fatal("Env file not found")
	}

	GetRedisClient()

	port := os.Getenv("PORT")

	if port == ""{
		log.Fatal("Error getting the port")
	}

	dbUrl := os.Getenv("DATABASE_URL")

	if dbUrl == ""{
		log.Fatal("Error getting db url")
	}

	conn, err := sql.Open("postgres",dbUrl)

	if err != nil{
		log.Fatal("Error connecting to db")
	}

	queries := database.New(conn)


	apiCfg := apiConfig{
		DB : queries,
	}


	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*","http://*"},
		AllowedMethods: []string{"GET","POST","PUT","DELETE"},
		AllowedHeaders:	[]string{"*"},
		MaxAge: 		300,
	}))

	v1Router := chi.NewRouter()

	v1Router.Post("/user/signup",apiCfg.handlerCreateUser)

	router.Mount("/api/v1",v1Router)

	srv := &http.Server{
		Handler: router,
		Addr : ":" + port,
	}

	err = srv.ListenAndServe()

	if err != nil{
		log.Fatal("Error while serving the http client",err)
	}


}