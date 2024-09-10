package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type acceptedVals struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type Users struct {
	users map[int]User
}

type User struct {
	id       int
	username string
	email    string
}
type returnVals struct {
	Error    string       `json:"error,omitempty"`
	Id       int          `json:"id,omitempty"`
	Body     string       `json:"body,omitempty"`
	Username string       `json:"username,omitempty"`
	Email    string       `json:"email,omitempty"`
	Token    string       `json:"token,omitempty"`
	Users    map[int]User `json:"users,omitempty"`
}

func (apiCfg *apiConfig) handlerUsers(w http.ResponseWriter, r *http.Request) {
	// Unmarshal body data and return params
	params, err := getParams(r, w)
	if err != nil {
		log.Printf("Error: %v\n", err)
	}

	switch r.Method {
	// POST create a user
	case http.MethodPost:
		_, err := apiCfg.DB.Exec("INSERT INTO users (username, password, email) VALUES ($1, $2, $3)", params.Username, params.Password, params.Email)
		if err != nil {
			fmt.Printf("Error coming back: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		payload := &returnVals{
			Username: params.Username,
			Email:    params.Email,
		}
		respondWithJSON(w, 200, payload)

	case http.MethodGet:
		rows, err := apiCfg.DB.Query("SELECT id, username, email FROM users")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		userMap := make(map[int]User)
		users := &Users{
			users: userMap,
		}
		for rows.Next() {
			rows.Next()
			fmt.Printf("\nNext Row\n")
			user := User{}
			rows.Scan(&user.id, &user.username, &user.email)
			fmt.Printf("username: %v \nemail: %v\n", user.username, user.email)
			users.users[user.id] = user
			fmt.Printf("Users: %v\n", users.users)
		}
		fmt.Printf("Users completed: %v\n", users.users)
		payload := &returnVals{
			Users: users.users,
		}
		fmt.Printf("Payload: %v\n", payload)

		respondWithJSON(w, http.StatusOK, payload)
	}
}

// func (apiCfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
// 	// Unmarshal body data and return params
// 	params, err := getParams(r, w)
// 	if err != nil {
// 		log.Printf("Error: %v\n", err)
// 	}

// 	switch r.Method {
// 	// POST login a user
// 	case http.MethodPost:
// 		user, id, token, err := apiCfg.DB.LoginUser(params.Email, params.Password, jwtSecret, params.ExpiresInSeconds)
// 		if err != nil {
// 			respondWithError(w, 401, "Unauthorized")
// 		}

// 		payload := &returnVals{
// 			Id:    id,
// 			Email: user.Email,
// 			Token: token,
// 		}

// 		respondWithJSON(w, 200, payload)
// 	}
// }

func getParams(r *http.Request, w http.ResponseWriter) (acceptedVals, error) {
	decoder := json.NewDecoder(r.Body)
	params := acceptedVals{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Error decoding params")
		return params, err
	}

	return params, nil
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(code)
	w.Write(data)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	w.Header().Add("Content-Type", "application/json")
	respBody := returnVals{
		Error: msg,
	}

	data, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		return
	}

	w.WriteHeader(code)
	w.Write(data)
}
