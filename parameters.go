package main

import (
	"userdb"
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
