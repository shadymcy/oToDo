package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/bll"
	"github.com/yzx9/otodo/web/common"
)

// Get basic todo lists for current user
func GetCurrentUserHandler(c *gin.Context) {
	userID := common.MustGetAccessUserID(c)
	user, err := bll.GetUser(userID)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

// Get todo lists for current user
func GetCurrentUserTodoListsHandler(c *gin.Context) {
	userID := common.MustGetAccessUserID(c)
	todos, err := bll.GetTodoLists(userID)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, todos)
}

// Get todo list
func GetTodoListHandler(c *gin.Context) {
	id, ok := c.Params.Get("id")
	if !ok {
		common.AbortWithError(c, fmt.Errorf("id required"))
		return
	}

	getTodoListHandler(c, id)
}

// Get basic todo lists for current user
func GetCurrentUserBasicTodoListHandler(c *gin.Context) {
	userID := common.MustGetAccessUserID(c)
	user, err := bll.GetUser(userID)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	getTodoListHandler(c, user.BasicTodoListID)
}

func getTodoListHandler(c *gin.Context, todoListID string) {
	userID := common.MustGetAccessUserID(c)
	todoList, err := bll.GetTodoList(userID, todoListID)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, todoList)
}

// Get todos in todo list
func GetTodoListTodosHandler(c *gin.Context) {
	todoListID, ok := c.Params.Get("id")
	if !ok {
		common.AbortWithError(c, fmt.Errorf("id required"))
		return
	}

	userID := common.MustGetAccessUserID(c)
	todos, err := bll.GetTodos(userID, todoListID)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, todos)
}

// Delete todo list
func DeleteTodoListHandler(c *gin.Context) {
	id, ok := c.Params.Get("id")
	if !ok {
		common.AbortWithError(c, fmt.Errorf("id required"))
		return
	}

	userID := common.MustGetAccessUserID(c)
	todo, err := bll.DeleteTodoList(userID, id)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, todo)
}
