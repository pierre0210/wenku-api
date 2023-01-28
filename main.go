package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pierre0210/wenku-reverseproxy/internal/wenku"
)

type Novel struct {
	Aid     int
	Volume  int
	Chapter int
}

func modifyRequest(r *http.Request) *http.Request {
	var novel Novel
	err := json.NewDecoder(r.Body).Decode(&novel)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	aid := novel.Aid
	volumeList := wenku.GetVolumeList(aid)
	vid := volumeList[(novel.Volume - 1)].Vid
	req, err := http.NewRequest("GET", fmt.Sprintf("https://dl.wenku8.com/pack.php?aid=%d&vid=%d", aid, vid), nil)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return req
}

func handler(w http.ResponseWriter, r *http.Request) {
	req := modifyRequest(r)
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		w.Write([]byte(err.Error()))
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err.Error())
		w.Write([]byte(err.Error()))
	}
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
