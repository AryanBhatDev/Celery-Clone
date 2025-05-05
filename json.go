package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)






func respondWithJson(w http.ResponseWriter, code int , payload interface{}){
	data, err := json.Marshal(payload)

	if err != nil{
		error := fmt.Sprintf("Error marshaling data %v and err:%s",payload,err)
		fmt.Println(error)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-type","application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func respondWithError(w http.ResponseWriter, code int , msg string){
	if code > 499{
		error := fmt.Sprintf("Internal error with code %d and err:%v",code, msg)
		fmt.Println(error)
	}

	type Error struct{
		Error string `json:"error"`
	}

	respondWithJson(w, 400, Error{
		Error: msg,
	})
}