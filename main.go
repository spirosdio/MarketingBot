package main

import (
	"parameters"
	"servercomm"
	"sync"
	"tree"
	"userdb"
)

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
