package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

var (
	db           = NewUserDatabase()
	UserNodeDict = make(map[*User]*TreeNode)
	Stats        = map[string]int{
		"users":              0,
		"positive_responses": 0,
		"negative_responses": 0,
		"coupon_reveals":     0,
		"messages_read":      0,
	}
)

func CreateWelcomingMessage(user *User) map[string]interface{} {
	welcomingAttachment := map[string]interface{}{
		"actions": []map[string]interface{}{
			{
				"name": "show_coupon",
				"text": "Yes! Show me coupon",
				"type": "button",
			},
			{
				"name": "no_thanks",
				"text": "No, thanks",
				"type": "button",
			},
		},
	}

	return map[string]interface{}{
		"user_id":      user.UserID,
		"callback_url": "http://localhost:6000/webhook/callback",
		"text":         fmt.Sprintf("Welcome to the demo promotional flow, %s! Are you interested in our coupon promotion?", user.Name),
		"attachments":  welcomingAttachment,
	}
}

func CreateCouponMessage(user *User) map[string]interface{} {
	couponAttachment := map[string]interface{}{
		"actions": []map[string]interface{}{
			{
				"name": "reveal_coupon",
				"text": "Reveal Coupon",
				"type": "button",
				"url":  "https://example.com/coupon/reveal",
			},
		},
	}

	return map[string]interface{}{
		"user_id":      user.UserID,
		"callback_url": "http://localhost:6000/webhook/callback",
		"text":         "Here is our unique promotional coupon! 10% off. Limit 1 per customer.",
		"attachments":  couponAttachment,
	}
}

func CreateMediaMessage(user *User) map[string]interface{} {
	mediaAttachment := map[string]interface{}{
		"image_url": "jpg", // Note: You might want to replace this with the actual JPG URL.
	}

	return map[string]interface{}{
		"user_id":      user.UserID,
		"callback_url": "http://localhost:6000/webhook/callback",
		"text":         "No worries! Have a nice day!",
		"attachments":  mediaAttachment,
	}
}

func RunClientServer() {
	http.HandleFunc("/webhook/callback", handleWebhook)
	http.HandleFunc("/get_stats", getStats)
	err := http.ListenAndServe(":6000", nil)
	if err != nil {
		return
	}
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
	user := db.GetUser(data["user_id"].(int))
	handleAnswerReceived(user, data)

	// Respond to the webhook
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	if err != nil {
		return
	}
}

func handleAnswerReceived(user *User, data map[string]interface{}) {
	log.Printf("Received answer from user %s , %v\n", user.Name, data)

	eventType, ok := data["event_type"].(string)
	if !ok {
		log.Println("Error: event_type missing or not a string.")
		return
	}

	switch eventType {
	case "message_read":
		Stats["messages_read"]++
		log.Printf("User %s %s read the message at %v\n", user.Name, user.Surname, data["interaction_timestamp"])

	case "link_click":
		log.Printf("User %s %s clicked the link!!!!!\n", user.Name, user.Surname)
		Stats["coupon_reveals"]++

	default:
		node, exists := UserNodeDict[user]
		if !exists {
			log.Println("Error: User not found in UserNodeDict.")
			return
		}

		buttonName, ok := data["button_name"].(string)
		if !ok {
			log.Println("Error: button_name missing or not a string.")
			return
		}

		var nextNode *TreeNode
		for _, child := range node.Children {
			if child.Data.ButtonNameOfOrigin == buttonName {
				nextNode = child
				break
			}
		}

		if nextNode == nil {
			log.Println("Error: No matching child node found.")
			return
		}

		if buttonName == "no_thanks" {
			Stats["negative_responses"]++
		} else {
			Stats["positive_responses"]++
		}

		sendMessageMockServer(nextNode.Data.JSON(user))
	}
}

func sendMessageMockServer(data map[string]interface{}) {
	serverURL := "http://localhost:5000/api/message"
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Failed to marshal data: %v\n", err)
		return
	}

	response, err := http.Post(serverURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Failed to send request: %v\n", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		var respData map[string]interface{}
		if err := json.NewDecoder(response.Body).Decode(&respData); err != nil {
			log.Println("Error decoding response:", err)
			return
		}
		log.Printf("Server responded with: %v\n", respData)
	} else {
		log.Printf("Failed to send message. Server responded with status code: %d\n", response.StatusCode)
	}
}
func getStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(Stats)
	if err != nil {
		return
	}
}

type User struct {
	Name    string
	Surname string
	UserID  int
	Age     string
}

type UserDatabase struct {
	users map[int]*User
}

func NewUserDatabase() *UserDatabase {
	return &UserDatabase{
		users: make(map[int]*User),
	}
}

func (db *UserDatabase) CreateUser(name string, surname string, userID int, age string) *User {
	if _, exists := db.users[userID]; !exists {
		newUser := &User{
			Name:    name,
			Surname: surname,
			UserID:  userID,
			Age:     age,
		}
		db.users[userID] = newUser
		return newUser
	}
	return nil
}

func (db *UserDatabase) GetUser(userID int) *User {
	return db.users[userID]
}

func (db *UserDatabase) GetAllUsers() []*User {
	var allUsers []*User
	for _, user := range db.users {
		allUsers = append(allUsers, user)
	}
	return allUsers
}

func (db *UserDatabase) ListUsers() {
	for _, user := range db.users {
		fmt.Println(user.Name, user.Surname, user.UserID, user.Age)
	}
}

type Data struct {
	JSON               func(*User) map[string]interface{}
	ButtonNameOfOrigin string
}

func NewData(jsonFunc func(*User) map[string]interface{}, buttonNameOfOrigin string) *Data {
	return &Data{
		JSON:               jsonFunc,
		ButtonNameOfOrigin: buttonNameOfOrigin,
	}
}

type TreeNode struct {
	Data     *Data
	Children []*TreeNode
}

func NewTreeNode(data *Data) *TreeNode {
	return &TreeNode{
		Data:     data,
		Children: make([]*TreeNode, 0),
	}
}

func (node *TreeNode) AddChild(childNode *TreeNode) {
	node.Children = append(node.Children, childNode)
}

func createUsers(db *UserDatabase) {
	db.CreateUser("spiros", "diochnos", 1, "26")
	db.CreateUser("vaso", "kollia", 2, "27")
	db.CreateUser("angelos", "todri", 3, "28")
}

func main() {
	// Start the client's server in a separate goroutine
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		RunClientServer()
	}()

	createUsers(db)

	// Creating data instances
	dataRoot := NewData(CreateWelcomingMessage, "")
	dataChild1 := NewData(CreateCouponMessage, "show_coupon")
	dataChild2 := NewData(CreateMediaMessage, "no_thanks")

	// Creating the tree nodes with data instances
	root := NewTreeNode(dataRoot)
	child1 := NewTreeNode(dataChild1)
	child2 := NewTreeNode(dataChild2)

	// Building the tree structure
	root.AddChild(child1)
	root.AddChild(child2)

	Stats["users"] = len(db.GetAllUsers())

	for _, user := range db.GetAllUsers() {
		UserNodeDict[user] = root
		// send message to mock server logic here
		// sendMessageMockServer(root.Data.JSON(user))
	}

	wg.Wait()
}
