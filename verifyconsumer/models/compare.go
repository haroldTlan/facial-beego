package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"github.com/astaxie/beego"
	"verifyconsumer/models/util"
)

type Extract struct {
	Features interface{} `json:"features"`
	TimeUsed int         `json:"time_used"`
}

type Compare struct {
	Confidence float64 `json:"confidence"`
	TimeUsed   int     `json:"time_used"`
}

var (
	extract  = beego.AppConfig.String("locExtract")
	compares = beego.AppConfig.String("locCompare")
)

//Compare
func Comparing(file1, file2 multipart.File) (conf float64, err error) {
	feat1, err := postFile(extract, file1)
	if err != nil {
		util.AddLog(err)
		return
	}
	feat2, err := postFile(extract, file2)
	if err != nil {
		util.AddLog(err)
		return
	}

	compare, err := getRemoteData(compares, feat1, feat2)
	if err != nil {
		util.AddLog(err)
		return
	}
	conf = compare.Confidence //TODO different type
	return
}

//get feature
func getFeat(feat Extract) (f string, err error) {
	if feature, ok := feat.Features.([]interface{}); ok {
		for _, i := range feature {
			f = i.(string)
		}
	} else if feat.Features == nil {
		err = fmt.Errorf("unknown picture")
		util.AddLog(err)
		return
	}
	return
}

//use feature to compare
func getRemoteData(url string, f1, f2 Extract) (conf Compare, err error) {
	feat1, err := getFeat(f1)
	if err != nil {
		util.AddLog(err)
		return
	}

	feat2, err := getFeat(f2)
	if err != nil {
		util.AddLog(err)
		return
	}

	query := map[string]string{
		"feat1": feat1,
		"feat2": feat2,
	}

	jsonString, err := json.Marshal(query)
	if err != nil {
		util.AddLog(err)
		return
	}

	req, err := DateHttpRequest(url, string(jsonString))
	if err != nil {
		util.AddLog(err)
		return
	}

	if err = json.Unmarshal([]byte(req), &conf); err != nil {
		util.AddLog(err)
		return
	}

	return
}

//POST features
func DateHttpRequest(url, query string) (body string, err error) {
	var jsonStr = []byte(query)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		util.AddLog(err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		util.AddLog(err)
		return
	}
	defer resp.Body.Close()

	bodys, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		util.AddLog(err)
		return
	}
	body = string(bodys)

	return
}

//POST file
func postFile(url string, f multipart.File) (body Extract, err error) {
	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	// Add your image file

	defer f.Close()
	fw, err := w.CreateFormFile("image", "image")
	if err != nil {
		util.AddLog(err)
		return
	}
	if _, err = io.Copy(fw, f); err != nil {
		util.AddLog(err)
		return
	}
	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		util.AddLog(err)
		return
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Submit the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		util.AddLog(err)
		return
	}
	// Check the response
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", res.Status)
	}
	bodys, err := ioutil.ReadAll(res.Body)
	if err != nil {
		util.AddLog(err)
		return
	}
	if err = json.Unmarshal(bodys, &body); err != nil {
		util.AddLog(err)
		return
	}
	return
}
