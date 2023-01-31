package wenku

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
	"github.com/pierre0210/wenku-api/internal/util"
)

type Volume struct {
	Title       string    `json:"title"`
	Vid         int       `json:"vid"`
	ChapterList []Chapter `json:"chapterList"`
}

type Chapter struct {
	Title   string `json:"title"`
	Cid     int    `json:"cid"`
	Content string `json:"content"`
}

func GetVolumeList(aid int) (string, []Volume) {
	var volumeList []Volume
	var novelTitle string
	c := colly.NewCollector()
	c.OnHTML("#title", func(h *colly.HTMLElement) {
		titleByte, _ := util.GbkToUtf8([]byte(h.Text))
		novelTitle = string(titleByte)
	})
	c.OnHTML("td", func(h *colly.HTMLElement) {
		if h.DOM.HasClass("vcss") {
			var vol Volume
			vid, err := strconv.Atoi(h.Attr("vid"))
			if err != nil {
				log.Println(err.Error())
				return
			}
			titleByte, _ := util.GbkToUtf8([]byte(h.Text))
			vol.Title = string(titleByte)
			vol.Vid = vid
			volumeList = append(volumeList, vol)
		} else if h.DOM.HasClass("ccss") && h.ChildAttr("a", "href") != "" {
			var ch Chapter
			titleByte, _ := util.GbkToUtf8([]byte(h.ChildText("a")))
			ch.Title = string(titleByte)
			ch.Cid, _ = strconv.Atoi(strings.Split(h.ChildAttr("a", "href"), "&cid=")[1])
			volumeList[len(volumeList)-1].ChapterList = append(volumeList[len(volumeList)-1].ChapterList, ch)
		}
	})

	c.Visit(fmt.Sprintf("https://www.wenku8.net/modules/article/reader.php?aid=%d", aid))

	return novelTitle, volumeList
}
