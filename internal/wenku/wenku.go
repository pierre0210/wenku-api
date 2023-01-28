package wenku

import (
	"fmt"
	"strconv"

	"github.com/gocolly/colly"
)

type Volume struct {
	Title string
	Vid   int
}

func GetVolumeList(aid int) []Volume {
	var volumeList []Volume
	c := colly.NewCollector()
	c.OnHTML(".vcss", func(h *colly.HTMLElement) {
		var vol Volume
		vid, err := strconv.Atoi(h.Attr("vid"))
		if err != nil {
			fmt.Printf("Error %s", err)
			return
		}
		vol.Title = h.Text
		vol.Vid = vid
		volumeList = append(volumeList, vol)
	})
	c.Visit(fmt.Sprintf("https://www.wenku8.net/modules/article/reader.php?aid=%d", aid))

	return volumeList
}

func GetCapterList() {

}
