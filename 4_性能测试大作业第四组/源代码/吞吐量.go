package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/myzhan/boomer"
)

var t *http.Transport
var c *http.Client

type Response struct {
	Status       int    `json:"status"`
	ContentID    string `json:"contentID"`
	ContentToken string `json:"contentToken"`
	FileID       string `json:"fileID"`
}

func contentThroughput() {
	start := time.Now()
	var response Response

	// Create Content
	file, _ := ioutil.ReadFile("tx.jpg")
	var buffer bytes.Buffer
	buffer.Write(file)
	buffer.Write(make([]byte, rand.Intn(31457280)+10485760)) // 10M ~ 40M

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("type", "1")
	part, err := writer.CreateFormFile("content", "tx.jpg")
	part.Write(buffer.Bytes())
	contentType := writer.FormDataContentType()
	writer.Close()

	req, err := http.NewRequest("POST", "http://10.0.0.28:30444/v1/content", body)
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjI1OTI0ODY4OTksImlkIjoxLCJyb2xlIjoxfQ.i0YUDEcnJrdgLOMTMpsnp85mHEnJkueeqoWrQg-d_o0")
	res, err := c.Do(req)

	{
		if err != nil {
			boomer.RecordFailure("http", "contentThroughput", time.Since(start).Nanoseconds()/int64(time.Millisecond), err.Error())
			return
		}
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			boomer.RecordFailure("http", "contentThroughput", time.Since(start).Nanoseconds()/int64(time.Millisecond), err.Error())
			return
		}
		defer res.Body.Close()
		elapsed := time.Since(start)
		if err != nil {
			boomer.RecordFailure("http", "contentThroughput", elapsed.Nanoseconds()/int64(time.Millisecond), err.Error())
			return
		}
		json.Unmarshal(data, &response)
		if res.StatusCode != 200 {
			boomer.RecordFailure("http", "contentThroughput", elapsed.Nanoseconds()/int64(time.Millisecond), res.Status)
			return
		}
		boomer.RecordSuccess("http", "contentThroughput", elapsed.Nanoseconds()/int64(time.Millisecond), int64(len(data)))
	}

	// Get Content
	res, err = c.Get("http://10.0.0.28:30444/v1/content/" + response.ContentID)

	{
		if err != nil {
			boomer.RecordFailure("http", "contentThroughput", time.Since(start).Nanoseconds()/int64(time.Millisecond), err.Error())
			return
		}
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			boomer.RecordFailure("http", "contentThroughput", time.Since(start).Nanoseconds()/int64(time.Millisecond), err.Error())
			return
		}
		defer res.Body.Close()
		elapsed := time.Since(start)
		if err != nil {
			boomer.RecordFailure("http", "contentThroughput", elapsed.Nanoseconds()/int64(time.Millisecond), err.Error())
			return
		}
		if res.StatusCode != 200 {
			boomer.RecordFailure("http", "contentThroughput", elapsed.Nanoseconds()/int64(time.Millisecond), res.Status)
			return
		}
		boomer.RecordSuccess("http", "contentThroughput", elapsed.Nanoseconds()/int64(time.Millisecond), int64(len(data)))
	}

	// Delete Content
	body = &bytes.Buffer{}
	writer = multipart.NewWriter(body)
	writer.WriteField("contentID", response.ContentID)
	writer.WriteField("contentToken", response.ContentToken)
	contentType = writer.FormDataContentType()
	writer.Close()
	req, err = http.NewRequest("DELETE", "http://10.0.0.28:30444/v1/content", bytes.NewReader(body.Bytes()))
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjI1OTI0ODY4OTksImlkIjoxLCJyb2xlIjoxfQ.i0YUDEcnJrdgLOMTMpsnp85mHEnJkueeqoWrQg-d_o0")
	res, err = c.Do(req)

	{
		if err != nil {
			boomer.RecordFailure("http", "contentThroughput", time.Since(start).Nanoseconds()/int64(time.Millisecond), err.Error())
			return
		}
		length, err := io.Copy(ioutil.Discard, res.Body)
		elapsed := time.Since(start)
		if err != nil {
			boomer.RecordFailure("http", "contentThroughput", elapsed.Nanoseconds()/int64(time.Millisecond), err.Error())
			return
		}
		defer res.Body.Close()

		if res.StatusCode != 200 {
			boomer.RecordFailure("http", "contentThroughput", elapsed.Nanoseconds()/int64(time.Millisecond), res.Status)
			return
		}
		boomer.RecordSuccess("http", "contentThroughput", elapsed.Nanoseconds()/int64(time.Millisecond), length)
	}
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

	task1 := &boomer.Task{
		Name:   "contentThroughput",
		Weight: 10,
		Fn:     contentThroughput,
	}

	boomer.Run(task1)
}
