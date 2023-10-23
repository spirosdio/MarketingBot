package main

import (
	"encoding/json"
	"net/http"
	"parameters"
	"userdb"
)

func RunClientServer() {
	http.HandleFunc("/webhook/callback", handleWebhook)
	http.HandleFunc("/get_stats", getStats)
	http.ListenAndServe(":6000", nil)
}

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Extract user and handle the received answer
	user := userdb.GetUser(data["user_id"].(int))
	handleAnswerReceived(user, data)

	// Respond to the webhook
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func handleAnswerReceived(user *userdb.User, data map[string]interface{}) {
	// Implement the logic to handle the answer received
}

func sendMessagetoMockServer(data map[string]interface{}) {
	// Implement the logic to send message to mock server
}

func getStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(parameters.Stats)
}
