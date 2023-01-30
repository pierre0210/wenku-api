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

func getVolume(aidNum int, vidNum int) (int, volumeResponse) {
	volumeList := wenku.GetVolumeList(aidNum)
	vidNum = volumeList[(vidNum - 1)].Vid

	req, _ := http.NewRequest("GET", fmt.Sprintf("https://dl.wenku8.com/pack.php?aid=%d&vid=%d", aidNum, vidNum), nil)
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Println(err.Error())
		return 500, volumeResponse{Message: err.Error()}
	} else if res.StatusCode == 404 {
		log.Println("Volume not found.")
		return 404, volumeResponse{Message: "Volume not found."}
	} else if res.StatusCode != 200 {
		log.Println("Other problem.")
		return res.StatusCode, volumeResponse{Message: "Other problem."}
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err.Error())
		return 500, volumeResponse{Message: err.Error()}
	}
	return 200, volumeResponse{Message: "Volume found.", Volume: string(body)}
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
	statusCode, res := getVolume(aidNum, vidNum)
	c.JSON(statusCode, res)
}

func HandleGetChapter(c *gin.Context) {

}
