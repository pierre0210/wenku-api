package novel

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pierre0210/wenku-api/internal/wenku"
)

type volumeResponse struct {
	Message string `json:"message"`
	Volume  string `json:"volume"`
}

func HandleGetVolume(c *gin.Context) {
	var response volumeResponse
	aid := c.Param("aid")
	vid := c.Param("vid")

	aidNum, aidErr := strconv.Atoi(aid)
	vidNum, vidErr := strconv.Atoi(vid)
	if aidErr != nil || vidErr != nil {
		log.Println("Invalid params data type.")
		c.JSON(400, volumeResponse{Message: "Invalid params data type."})
		return
	}
	volumeList := wenku.GetVolumeList(aidNum)
	vidNum = volumeList[(vidNum - 1)].Vid

	req, _ := http.NewRequest("GET", fmt.Sprintf("https://dl.wenku8.com/pack.php?aid=%d&vid=%d", aidNum, vidNum), nil)
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Println(err.Error())
		c.JSON(500, volumeResponse{Message: err.Error()})
		return
	} else if res.StatusCode == 404 {
		log.Println("Volume not found.")
		c.JSON(404, volumeResponse{Message: "Volume not found."})
		return
	} else if res.StatusCode != 200 {
		log.Println("Other problem.")
		c.JSON(res.StatusCode, volumeResponse{Message: "Other problem."})
		return
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err.Error())
		c.JSON(500, volumeResponse{Message: err.Error()})
		return
	}
	response.Message = "Volume found."
	response.Volume = string(body)
	c.JSON(200, response)
}
