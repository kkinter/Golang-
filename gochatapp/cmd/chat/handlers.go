package main

import (
	"encoding/json"
	"fmt"
	"gochatapp/internel/redisrepo"
	"log"
	"net/http"
)

type userReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Client   string `json:"client"`
}

type response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Total   int         `json:"total,omitempty"`
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	u := &userReq{}
	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		http.Error(w, "error decoding request obj", http.StatusBadRequest)
		return
	}

	resp := register(u)
	json.NewEncoder(w).Encode(resp)
}

func register(u *userReq) *response {
	resp := &response{Status: true}

	status := redisrepo.IsUserExist(u.Username)
	if status {
		resp.Status = false
		resp.Message = "username already taken. try something else"
		return resp
	}

	err := redisrepo.RegisterNewUser(u.Username, u.Password)
	if err != nil {
		resp.Status = false
		resp.Message = "something went wrong while registering the user. please try again after sometime."
		return resp
	}

	return resp
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	u := &userReq{}
	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		http.Error(w, "error decoidng request object", http.StatusBadRequest)
		return
	}

	resp := login(u)
	json.NewEncoder(w).Encode(resp)
}

func login(u *userReq) *response {
	resp := &response{Status: true}

	err := redisrepo.IsUserAuthentic(u.Username, u.Password)
	if err != nil {
		resp.Status = false
		resp.Message = err.Error()
		return resp
	}

	return resp
}

func verifyContactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	u := &userReq{}
	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		http.Error(w, "error decoidng request object", http.StatusBadRequest)
		return
	}

	res := verifyContact(u.Username)
	json.NewEncoder(w).Encode(res)
}

func verifyContact(username string) *response {
	res := &response{Status: true}

	status := redisrepo.IsUserExist(username)
	if !status {
		res.Status = false
		res.Message = "invalid username"
	}

	return res
}

func chatHistoryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// user1 user2
	u1 := r.URL.Query().Get("u1")
	u2 := r.URL.Query().Get("u2")

	// chat between timerange fromTS toTS
	// where TS is timestamp
	// 0 to positive infinity
	fromTS, toTS := "0", "+inf"

	if r.URL.Query().Get("from-ts") != "" && r.URL.Query().Get("to-ts") != "" {
		fromTS = r.URL.Query().Get("from-ts")
		toTS = r.URL.Query().Get("to-ts")
	}

	res := chatHistory(u1, u2, fromTS, toTS)
	json.NewEncoder(w).Encode(res)
}

func chatHistory(username1, username2, fromTS, toTS string) *response {
	res := &response{}

	fmt.Println(username1, username2)

	if !redisrepo.IsUserExist(username1) || !redisrepo.IsUserExist(username2) {
		res.Message = "incorrect username"
		return res
	}

	chats, err := redisrepo.FetchChatBetween(username1, username2, fromTS, toTS)
	if err != nil {
		log.Println("error in fetch chat between", err)
		res.Message = "unable to fetch chat history. please try again later."
		return res
	}

	res.Status = true
	res.Data = chats
	res.Total = len(chats)
	return res
}

func contactListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	u := r.URL.Query().Get("username")

	res := contactList(u)
	json.NewEncoder(w).Encode(res)
}

func contactList(username string) *response {
	res := &response{}

	if !redisrepo.IsUserExist(username) {
		res.Message = "incorrect username"
		return res
	}

	contactList, err := redisrepo.FetchContactList(username)
	if err != nil {
		log.Println("error in fetch contact list of username: ", username, err)
		res.Message = "unable to fetch contact list. please try again later."
		return res
	}

	res.Status = true
	res.Data = contactList
	res.Total = len(contactList)
	return res
}
