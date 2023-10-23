package main

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
