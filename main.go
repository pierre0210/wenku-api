package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/pierre0210/wenku-reverseproxy/internal/wenku"
)

type Novel struct {
	Aid     int
	Volume  int
	Chapter int
}

func modifyRequest(r *http.Request) *http.Request {
	var novel Novel
	//err := json.NewDecoder(r.Body).Decode(&novel)
	novel.Aid, _ = strconv.Atoi(r.URL.Query().Get("aid"))
	novel.Volume, _ = strconv.Atoi(r.URL.Query().Get("vid"))

	volumeList := wenku.GetVolumeList(novel.Aid)
	if novel.Volume-1 >= len(volumeList) {
		fmt.Println("Index out of range.")
		return nil
	}
	vid := volumeList[(novel.Volume - 1)].Vid
	req, err := http.NewRequest("GET", fmt.Sprintf("https://dl.wenku8.com/pack.php?aid=%d&vid=%d", novel.Aid, vid), nil)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return req
}

func handler(w http.ResponseWriter, r *http.Request) {
	if !r.URL.Query().Has("aid") || !r.URL.Query().Has("vid") {
		w.Write([]byte("Home page."))
		return
	}
	req := modifyRequest(r)
	if req == nil {
		w.Write([]byte("Some thing went wrong. No content."))
		return
	}
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		w.Write([]byte(err.Error()))
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err.Error())
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.Write(body)
}

func main() {
	port := flag.Int("p", 5000, "Port")
	flag.Parse()

	addr := fmt.Sprintf(":%d", *port)

	http.HandleFunc("/", handler)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		panic(err)
	}
}
