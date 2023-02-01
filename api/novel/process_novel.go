package novel

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/pierre0210/wenku-api/internal/database"
	chapterTable "github.com/pierre0210/wenku-api/internal/database/table/chapter"
	"github.com/pierre0210/wenku-api/internal/util"
	"github.com/pierre0210/wenku-api/internal/wenku"
)

type volumeIndex struct {
	Title       string         `json:"title"`
	Vid         int            `json:"vid"`
	ChapterList []chapterIndex `json:"chapterList"`
}

type chapterIndex struct {
	Title string `json:"title"`
	Cid   int    `json:"cid"`
}

type novelIndex struct {
	Title      string        `json:"title"`
	Author     string        `json:"author"`
	Aid        int           `json:"aid"`
	Cover      string        `json:"cover"`
	VolumeList []volumeIndex `json:"volumeList"`
}

func splitVolume(content string, volume wenku.Volume, aid int, vid int) {
	for index, chapter := range volume.ChapterList {
		r, _ := regexp.Compile(`<div class="chaptertitle"><a name="` + strconv.Itoa(chapter.Cid) + `">[\s\S]+?<span></span></div>`)
		rHtml, _ := regexp.Compile("<[^<>]+>")
		rUrl, _ := regexp.Compile(`https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`)
		volume.ChapterList[index].Content = r.FindString(content)
		volume.ChapterList[index].Content = strings.ReplaceAll(volume.ChapterList[index].Content, "<br />\r\n<br />", "\r\n")
		volume.ChapterList[index].Content = strings.ReplaceAll(volume.ChapterList[index].Content, "&nbsp;", " ")
		volume.ChapterList[index].Content = strings.ReplaceAll(volume.ChapterList[index].Content, "</a></div>", "")
		volume.ChapterList[index].Content = strings.ReplaceAll(volume.ChapterList[index].Content, "</a>", "\r\n")
		volume.ChapterList[index].Content = rHtml.ReplaceAllString(volume.ChapterList[index].Content, "")
		volume.ChapterList[index].Content, _ = util.SimplifyToTW(volume.ChapterList[index].Content)
		volume.ChapterList[index].Urls = rUrl.FindAllString(volume.ChapterList[index].Content, -1)

		//chapterExists, _ := chapterTable.CheckChapter(database.DB, aid, vid, cid)
		contentObj, _ := json.Marshal(volume.ChapterList[index])
		_, err := chapterTable.AddChapter(database.DB, aid, vid, index+1, chapter.Title, string(contentObj))
		if err != nil {
			log.Println(err)
		}
	}
}

func getVolume(aidNum int, vidNum int) (int, volumeResponse, wenku.Volume) {
	_, _, volumeList := wenku.GetVolumeList(aidNum)
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
	chapterExists, _ := chapterTable.CheckChapter(database.DB, aid, vid, cid)
	if chapterExists {
		var chapterObj wenku.Chapter
		_, _, content, _ := chapterTable.GetChapter(database.DB, aid, vid, cid)
		json.Unmarshal([]byte(content), &chapterObj)

		return 200, chapterResponse{Message: "Saved chapter found.", Content: chapterObj}
	}
	statusCode, res, volume := getVolume(aid, vid)
	if statusCode != 200 {
		log.Printf("%d %s\n", statusCode, res.Message)
		return statusCode, chapterResponse{Message: res.Message}
	}
	splitVolume(res.Content, volume, aid, vid)
	if cid > len(volume.ChapterList) {
		log.Println("Index out of range.")
		return 404, chapterResponse{Message: "Index out of range."}
	}

	return 200, chapterResponse{Message: "Chapter found.", Content: volume.ChapterList[(cid - 1)]}
}

func getIndex(aid int) (int, indexResponse) {
	var index novelIndex
	var volumeList []wenku.Volume
	index.Aid = aid
	index.Title, index.Author, volumeList = wenku.GetVolumeList(aid)
	if index.Title == "" || len(volumeList) == 0 {
		return 404, indexResponse{Message: "Not found."}
	}
	index.Cover = fmt.Sprintf("https://img.wenku8.com/image/%d/%d/%ds.jpg", int(math.Floor(float64(aid)/1000)), aid, aid)
	for _, volume := range volumeList {
		index.VolumeList = append(index.VolumeList, volumeIndex{Title: volume.Title, Vid: volume.Vid})
		for _, chapter := range volume.ChapterList {
			index.VolumeList[len(index.VolumeList)-1].ChapterList = append(index.VolumeList[len(index.VolumeList)-1].ChapterList, chapterIndex{Title: chapter.Title, Cid: chapter.Cid})
		}
	}
	return 200, indexResponse{Message: "Index found.", Content: index}
}
