package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/myzhan/boomer"
)

var t *http.Transport
var c *http.Client
var fileIDs []string
var max = 100

func getContent() {
	// pre post many content
	// record all contentID
	id := fileIDs[rand.Intn(len(fileIDs))]
	fmt.Printf("id is :%s\n", id)
	start := time.Now()

	res, err := c.Get("http://10.0.0.28:30444/v1/file/" + id)
	if err != nil {
		boomer.RecordFailure("http", "getContent", time.Since(start).Nanoseconds()/int64(time.Millisecond), err.Error())
		return
	}
	length, err := io.Copy(ioutil.Discard, res.Body)
	elapsed := time.Since(start)

	if err != nil {
		boomer.RecordFailure("http", "getContent", elapsed.Nanoseconds()/int64(time.Millisecond), err.Error())
		return
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		boomer.RecordFailure("http", "getContent", elapsed.Nanoseconds()/int64(time.Millisecond), res.Status)
		return
	}
	boomer.RecordSuccess("http", "getContent", elapsed.Nanoseconds()/int64(time.Millisecond), length)
}

type avatarResponse struct {
	Status   int    `json:"status"`
	AvatarID string `json:"avatarID"`
}

func main() {
	t = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		MaxIdleConns:        10240,
		MaxIdleConnsPerHost: 10240,
		DisableKeepAlives:   true,
	}
	c = &http.Client{
		Timeout:   40 * time.Second,
		Transport: t,
	}
	bodyBuffer := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuffer)
	userIDWriter, _ := bodyWriter.CreateFormField("userID")
	_, _ = userIDWriter.Write([]byte("1"))
	fileWriter, _ := bodyWriter.CreateFormFile("file", "./tx.jpg")
	file, _ := os.Open("./tx.jpg")
	defer file.Close()

	_, _ = io.Copy(fileWriter, file)

	contentType := bodyWriter.FormDataContentType()
	_ = bodyWriter.Close()
	fileIDs = make([]string, max)
	for i := 0; i < max; i++ {
		req, _ := http.NewRequest("POST", "http://10.0.0.28:30444/v1/avatar/", bytes.NewReader(bodyBuffer.Bytes()))
		req.Header.Add("Content-Type", contentType)
		req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjI1OTI0ODY4OTksImlkIjoxLCJyb2xlIjoxfQ.i0YUDEcnJrdgLOMTMpsnp85mHEnJkueeqoWrQg-d_o0")
		resp, err := c.Do(req)
		if err != nil {
			panic(err)
		}
		if resp.StatusCode != http.StatusOK {
			panic(resp.StatusCode)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("read body err, %v\n", err)
			return
		}
		println("json:", string(body))

		var res avatarResponse
		if err = json.Unmarshal(body, &res); err != nil {
			fmt.Printf("Unmarshal err, %v\n", err)
			return
		}
		fileIDs[i] = res.AvatarID
	}

	id := fileIDs[rand.Intn(len(fileIDs))]
	fmt.Printf("id is :%s\n", id)

	_, err := c.Get("http://10.0.0.28:30444/v1/file/" + id)
	if err != nil {
		panic(err)
	}

	task1 := &boomer.Task{
		Name:   "getContent",
		Weight: 10,
		Fn:     getContent,
	}

	boomer.Run(task1)
}
