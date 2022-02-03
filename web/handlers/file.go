package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/bll"
	"github.com/yzx9/otodo/web/utils"
)

// Upload file, only support single file now
func PostFileHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		utils.AbortWithJson(c, "invalid file")
		return
	}

	filename, err := bll.UploadFile(file)
	if err != nil {
		utils.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, struct {
		FileName string `json:"file_name"`
	}{filename})
}

// Upload file, only support single file now
func GetFileHandler(c *gin.Context) {
	params := struct{ id string }{}
	err := c.ShouldBind(&params)
	if err != nil {
		utils.AbortWithJson(c, "invalid file")
		return
	}

	userID := utils.MustGetAccessUserID(c)
	filepath, err := bll.GetFilePath(params.id, userID)
	if err != nil {
		utils.AbortWithJson(c, "invalid file")
		return
	}

	c.File(filepath)
}
