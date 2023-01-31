package novel

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pierre0210/wenku-api/internal/wenku"
)

type volumeResponse struct {
	Message string `json:"message"`
	Content string `json:"content"`
}

type chapterResponse struct {
	Message string        `json:"message"`
	Content wenku.Chapter `json:"content"`
}

type indexResponse struct {
	Message string     `json:"message"`
	Content novelIndex `json:"content"`
}

func HandleGetVolume(c *gin.Context) {
	aid := c.Param("aid")
	vid := c.Param("vid")

	aidNum, aidErr := strconv.Atoi(aid)
	vidNum, vidErr := strconv.Atoi(vid)
	if aidErr != nil || vidErr != nil {
		log.Println("Invalid params data type.")
		c.JSON(400, volumeResponse{Message: "Invalid params data type."})
		return
	}
	statusCode, volumeRes, _ := getVolume(aidNum, vidNum)
	c.JSON(statusCode, volumeRes)
}

func HandleGetChapter(c *gin.Context) {
	aid := c.Param("aid")
	vid := c.Param("vid")
	cid := c.Param("cid")
	aidNum, aidErr := strconv.Atoi(aid)
	vidNum, vidErr := strconv.Atoi(vid)
	cidNum, cidErr := strconv.Atoi(cid)
	if aidErr != nil || vidErr != nil || cidErr != nil {
		log.Println("Invalid params data type.")
		c.JSON(400, chapterResponse{Message: "Invalid params data type."})
		return
	}
	statusCode, chapterRes := getChapter(aidNum, vidNum, cidNum)
	c.JSON(statusCode, chapterRes)
}

func HandleGetIndex(c *gin.Context) {
	aid := c.Param("aid")
	aidNum, aidErr := strconv.Atoi(aid)
	if aidErr != nil {
		log.Println("Invalid params data type.")
	}
	statusCode, indexRes := getIndex(aidNum)
	c.JSON(statusCode, indexRes)
}
