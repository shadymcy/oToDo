package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/bll"
	"github.com/yzx9/otodo/web/common"
)

// Get todo list folder
func GetTodoListFolderHandler(c *gin.Context) {
	id, err := common.GetRequiredParam(c, "id")
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	userID := common.MustGetAccessUserID(c)
	folder, err := bll.GetTodoListFolder(userID, id)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, folder)
}

// Delete todo list folder
func DeleteTodoListFolderHandler(c *gin.Context) {
	id, err := common.GetRequiredParam(c, "id")
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	userID := common.MustGetAccessUserID(c)
	todo, err := bll.DeleteTodoListFolder(userID, id)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, todo)
}
