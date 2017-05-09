package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"time"

	models2 "DemoForService/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"github.com/crackcomm/nsqueue/consumer"
	"verifyconsumer/controllers"
	"verifyconsumer/models/util"
)

type VerifyMsg struct {
	CompareMsg map[string]string
	ID         string
	Files      []string
	User       string
}

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

func handleTest(msg *consumer.Message) {
	var result VerifyMsg
	err := msg.ReadJSON(&result)
	fmt.Printf("result:\n%+v", result)
	if err != nil {
		log.Println("ReadJSON error: ", err)
		msg.Success()
		return
	}

	var confidence float64
	confidence, err = locFacePlusPlus(result.Files)

	if confidence < 0.7 {
		if err != nil {
			util.AddLog(fmt.Errorf("First error: ", err))
		}

		//confidence, err = CurlByFacePlusPlus(result.Files)
		//confidence, err = FacePlusPlus(result.Files)
		confidence, err = LinkFace(result.Files)
		//confidence, err = locFacePlusPlus(result.Files)
		if err != nil {
			util.AddLog(fmt.Errorf("Second error: ", err))
			publish <- controllers.NewCompareEvent(result.User, "0")
			msg.Success()
			return
		}

	}

	msg.Success()
	//confidence = math.Trunc(confidence*1e2+0.5) * 1e-2
	con := strconv.FormatFloat(confidence, 'f', -1, 64)
	publish <- controllers.NewCompareEvent(result.User, strconv.FormatFloat(confidence, 'f', -1, 32))
	uid := models2.Record(result.ID, result.User, con)
	fmt.Println(1)
	models2.InsertIntoCompareMsg(result.User, con, uid, result.CompareMsg)
	fmt.Println(2)
	models2.FileMv(result.User, result.ID, uid)
	fmt.Println(3)

	//publish <- controllers.NewCompareEvent(result.User, strconv.FormatFloat(confidence, 'f', -1, 32))
}

func runConsumer(maxInFlight int, nsqdAddr string) {
	count := 10
	for {
		consumer.Register("verify", "consume", maxInFlight, handleTest)
		err := consumer.Connect(nsqdAddr)
		if err == nil {
			break
		}
		time.Sleep(time.Second * 10)
		count -= 1
		if count == 0 {
			log.Println(err)
			os.Exit(1)
		}
	}

	consumer.Start(true)
}

func execVerifyByLocal(files []string) (float64, error) {

	// Create an *exec.Cmd
	cmd := exec.Command("/nvr/verify", files...)

	// Stdout buffer
	cmdOutput := &bytes.Buffer{}
	// Attach buffer to command
	cmd.Stdout = cmdOutput

	// Execute command
	err := cmd.Run() // will wait for command to return
	if err != nil {
		return float64(0), err
	}

	outs := cmdOutput.Bytes()
	if len(outs) > 0 {
		pattern := `^finish loading all nets\s+([0-9\.]+)\s+$`
		reg := regexp.MustCompile(pattern)
		result := reg.FindStringSubmatch(string(outs))
		if len(result) == 0 {
			return float64(0), fmt.Errorf("invaild outs of verify.exe")
		}

		confidence, err := strconv.ParseFloat(result[1], 32)
		if err != nil {
			return confidence, err
		}

		return confidence, nil
	}

	return float64(0), fmt.Errorf("failed to exec verify.exe")
}

func CurlByFacePlusPlus(files []string) (float64, error) {
	// Create an *exec.Cmd
	fmt.Println(files)
	cmdArgs := make([]string, 0)
	cmdArgs = append(cmdArgs, "https://api-cn.faceplusplus.com/facepp/v3/compare")
	cmdArgs = append(cmdArgs, "-F", "api_key="+*apiKey)
	cmdArgs = append(cmdArgs, "-F", "api_secret="+*apiSecret)
	cmdArgs = append(cmdArgs, "-F", "image_file1=@"+files[0])
	cmdArgs = append(cmdArgs, "-F", "image_file2=@"+files[1])

	cmd := exec.Command("curl", cmdArgs...)

	// Stdout buffer
	cmdOutput := &bytes.Buffer{}
	// Attach buffer to command
	cmd.Stdout = cmdOutput

	// Execute command
	err := cmd.Run() // will wait for command to return
	if err != nil {
		fmt.Println(err)
		return float64(0), err
	}

	outs := cmdOutput.Bytes()
	mp, err1 := response2Map(outs)
	if err1 != nil {
		return float64(0), err1
	}

	fmt.Println(mp["confidence"].(float64))

	return mp["confidence"].(float64), nil
}

func FacePlusPlus(files []string) (float64, error) {
	req := httplib.Post("https://api-cn.faceplusplus.com/facepp/v3/compare")
	req.Param("api_key", *apiKey)
	req.Param("api_secret", *apiSecret)
	req.PostFile("image_file1", files[0])
	req.PostFile("image_file2", files[1])

	outs, err := req.Bytes()
	if err != nil {
		log.Println("post error: ", err)
		return float64(0), err
	}
	mp, err1 := response2Map(outs)
	if err1 != nil {
		return float64(0), err1
	}

	if _, ok := mp["confidence"]; !ok {
		s, _ := strconv.ParseFloat("0", 64)
		return s, fmt.Errorf("confidence to low")
	}

	return mp["confidence"].(float64), nil
}

func LinkFace(files []string) (float64, error) {
	req := httplib.Post("https://cloudapi.linkface.cn/identity/historical_selfie_verification")
	req.Param("api_id", *apiKey)
	req.Param("api_secret", *apiSecret)
	req.PostFile("selfie_file", files[0])
	req.PostFile("historical_selfie_file", files[1])
	outs, err := req.Bytes()
	if err != nil {
		log.Println("post error: ", err)
		return float64(0), err
	}

	mp, err1 := response2Map(outs)
	if err1 != nil {
		return float64(0), err1
	}

	if _, ok := mp["confidence"]; !ok {
		s, _ := strconv.ParseFloat("0", 64)
		return s, fmt.Errorf("confidence to low")
	}

	fmt.Println(mp["confidence"].(float64))
	return mp["confidence"].(float64), nil
}

func locFacePlusPlus(files []string) (conf float64, err error) {
	feat1, err := postFile(extract, files[0])
	feat2, err := postFile(extract, files[1])

	compare, err := getRemoteData(compares, feat1, feat2)
	if err != nil {
		util.AddLog(err)
		return
	}
	conf = compare.Confidence //TODO different type
	return
}

//judge feature
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
func postFile(url, file string) (body Extract, err error) {
	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	// Add your image file
	f, err := os.Open(file)
	if err != nil {
		util.AddLog(err)
		return
	}
	defer f.Close()
	fw, err := w.CreateFormFile("image", file)
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
