package wenku

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

type Volume struct {
	Title       string
	Vid         int
	ChapterList []Chapter
}

type Chapter struct {
	Title string
	Cid   int
}

func GetVolumeList(aid int) []Volume {
	var volumeList []Volume
	c := colly.NewCollector()
	c.OnHTML("td", func(h *colly.HTMLElement) {
		if h.DOM.HasClass("vcss") {
			var vol Volume
			vid, err := strconv.Atoi(h.Attr("vid"))
			if err != nil {
				log.Println(err.Error())
				return
			}
			vol.Title = h.Text
			vol.Vid = vid
			volumeList = append(volumeList, vol)
		} else if h.DOM.HasClass("ccss") && h.ChildAttr("a", "href") != "" {
			var ch Chapter
			ch.Title = h.ChildText("a")
			ch.Cid, _ = strconv.Atoi(strings.Split(h.ChildAttr("a", "href"), "&cid=")[1])
			volumeList[len(volumeList)-1].ChapterList = append(volumeList[len(volumeList)-1].ChapterList, ch)
		}
	})

	c.Visit(fmt.Sprintf("https://www.wenku8.net/modules/article/reader.php?aid=%d", aid))

	return volumeList
}
