package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

var (
	UserNodeDict = make(map[*userdb.User]*tree.TreeNode)
	Stats        = map[string]int{
		"users":              0,
		"positive_responses": 0,
		"negative_responses": 0,
		"coupon_reveals":     0,
		"messages_read":      0,
	}
)

func CreateWelcomingMessage(user *userdb.User) map[string]interface{} {
	// Convert Python dictionary to Go map and return
}

func CreateCouponMessage(user *userdb.User) map[string]interface{} {
	// Convert Python dictionary to Go map and return
}

func CreateMediaMessage(user *userdb.User) map[string]interface{} {
	// Convert Python dictionary to Go map and return
}

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
		// Print user details
		fmt.Println(user.Name, user.Surname, user.UserID, user.Age)
	}
}

type Data struct {
	JSON               func(*userdb.User) map[string]interface{}
	ButtonNameOfOrigin string
}

func NewData(jsonFunc func(*userdb.User) map[string]interface{}, buttonNameOfOrigin string) *Data {
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

func (node *TreeNode) Display(level int) {
	// Implement the display logic here
}

func createUsers(db *userdb.UserDatabase) {
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
		servercomm.RunClientServer()
	}()

	db := userdb.NewUserDatabase()
	createUsers(db)

	// Creating data instances
	dataRoot := tree.NewData(parameters.CreateWelcomingMessage, "")
	dataChild1 := tree.NewData(parameters.CreateCouponMessage, "show_coupon")
	dataChild2 := tree.NewData(parameters.CreateMediaMessage, "no_thanks")

	// Creating the tree nodes with data instances
	root := tree.NewTreeNode(dataRoot)
	child1 := tree.NewTreeNode(dataChild1)
	child2 := tree.NewTreeNode(dataChild2)

	// Building the tree structure
	root.AddChild(child1)
	root.AddChild(child2)

	parameters.Stats["users"] = len(db.GetAllUsers())

	for _, user := range db.GetAllUsers() {
		parameters.UserNodeDict[user] = root
		// send message to mock server logic here
		// sendMessagetoMockServer(root.Data.JSON(user))
	}

	wg.Wait()
}
