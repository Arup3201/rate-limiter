package handlers

import (
	"encoding/json"
	"net/http"
)

type User struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

var users = []User{
	{
		Id:       "511c1c2c-b6f3-4160-bb8e-6d3d19e899b4",
		Username: "example.user.1",
		Email:    "user1@example.com",
	},
	{
		Id:       "f36fb408-40e9-4749-afb7-38742645cc24",
		Username: "example.user.2",
		Email:    "user2@example.com",
	},
	{
		Id:       "7bb5bc14-ed50-40c0-a066-c4e181d5f31f",
		Username: "example.user.3",
		Email:    "user3@example.com",
	},
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(users)
}
