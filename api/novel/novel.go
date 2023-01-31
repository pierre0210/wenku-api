package novel

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

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

func splitVolume(content string, volume wenku.Volume) {
	for index, chapter := range volume.ChapterList {
		r, _ := regexp.Compile(`<div class="chaptertitle"><a name="` + strconv.Itoa(chapter.Cid) + `">[\s\S]+?<span></span></div>`)
		volume.ChapterList[index].Content = r.FindString(content)
		volume.ChapterList[index].Content = strings.ReplaceAll(volume.ChapterList[index].Content, "<br />\r\n<br />", "\r\n")
		volume.ChapterList[index].Content = strings.ReplaceAll(volume.ChapterList[index].Content, "&nbsp;", " ")
	}
}

func getVolume(aidNum int, vidNum int) (int, volumeResponse, wenku.Volume) {
	volumeList := wenku.GetVolumeList(aidNum)
	if len(volumeList) == 0 || vidNum > len(volumeList) {
		return 404, volumeResponse{Message: "Not found."}, wenku.Volume{}
	}
	volume := volumeList[(vidNum - 1)]
	vidNum = volume.Vid

	req, _ := http.NewRequest("GET", fmt.Sprintf("https://dl.wenku8.com/pack.php?aid=%d&vid=%d", aidNum, vidNum), nil)
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Println(err.Error())
		return 500, volumeResponse{Message: err.Error()}, wenku.Volume{}
	} else if res.StatusCode == 404 {
		log.Println("Volume not found.")
		return 404, volumeResponse{Message: "Volume not found."}, wenku.Volume{}
	} else if res.StatusCode != 200 {
		log.Println("Other problem.")
		return res.StatusCode, volumeResponse{Message: "Other problem."}, wenku.Volume{}
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err.Error())
		return 500, volumeResponse{Message: err.Error()}, wenku.Volume{}
	}
	return 200, volumeResponse{Message: "Volume found.", Content: string(body)}, volume
}

func getChapter(aid int, vid int, cid int) (int, chapterResponse) {
	statusCode, res, volume := getVolume(aid, vid)
	if statusCode != 200 {
		log.Printf("%d %s\n", statusCode, res.Message)
		return statusCode, chapterResponse{Message: res.Message}
	}
	splitVolume(res.Content, volume)
	if cid > len(volume.ChapterList) {
		log.Println("Index out of range.")
		return 404, chapterResponse{Message: "Index out of range."}
	}
	return 200, chapterResponse{Message: "Chapter found.", Content: volume.ChapterList[(cid - 1)]}
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
		c.JSON(400, volumeResponse{Message: "Invalid params data type."})
		return
	}
	statusCode, chapterRes := getChapter(aidNum, vidNum, cidNum)
	c.JSON(statusCode, chapterRes)
}
