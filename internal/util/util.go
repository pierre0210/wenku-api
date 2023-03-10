package util

import (
	"bytes"
	"io/ioutil"
	"log"

	"github.com/liuzl/gocc"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
)

var s2tw, _ = gocc.New("s2tw")
var t2tw, _ = gocc.New("t2tw")

func GbkToUtf8(b []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(b), simplifiedchinese.GBK.NewDecoder())
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func Utf8ToBig5(b []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(b), traditionalchinese.Big5.NewEncoder())
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func SimplifiedToTW(content string) (string, error) {
	result, err := s2tw.Convert(content)
	if err != nil {
		log.Println(err)
		return content, err
	}
	return result, nil
}

func TraditionalToTW(content string) (string, error) {
	result, err := t2tw.Convert(content)
	if err != nil {
		log.Println(err)
		return content, err
	}
	return result, nil
}
