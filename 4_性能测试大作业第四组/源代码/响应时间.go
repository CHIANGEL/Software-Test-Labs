package main

import (
	"crypto/tls"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/myzhan/boomer"
)

var t *http.Transport
var c *http.Client
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var max = 100

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func goodSearch() {
	rand.Seed(time.Now().UnixNano())
	start := time.Now()
	param := url.Values{}
	param.Add("limit", strconv.Itoa(rand.Intn(30)+30))
	param.Add("offset", strconv.Itoa(rand.Intn(10)))
	res, err := c.Get("http://10.0.0.28:30444/v1/sellInfo?" + param.Encode())
	if err != nil {
		boomer.RecordFailure("http", "goodSearch", time.Since(start).Nanoseconds()/int64(time.Millisecond), err.Error())
		return
	}
	length, err := io.Copy(ioutil.Discard, res.Body)
	elapsed := time.Since(start)

	if err != nil {
		boomer.RecordFailure("http", "goodSearch", elapsed.Nanoseconds()/int64(time.Millisecond), err.Error())
		return
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		boomer.RecordFailure("http", "goodSearch", elapsed.Nanoseconds()/int64(time.Millisecond), res.Status)
		return
	}
	boomer.RecordSuccess("http", "goodSearch", elapsed.Nanoseconds()/int64(time.Millisecond), length)

	// click one sellinfo (50% probability)
	if rand.Intn(2) == 0 {
		start := time.Now()
		res, err := c.Get("http://10.0.0.28:30444/v1/sellInfo/50")
		if err != nil {
			boomer.RecordFailure("http", "getSellInfo", time.Since(start).Nanoseconds()/int64(time.Millisecond), err.Error())
			return
		}
		length, err := io.Copy(ioutil.Discard, res.Body)
		elapsed := time.Since(start)

		if err != nil {
			boomer.RecordFailure("http", "getSellInfo", elapsed.Nanoseconds()/int64(time.Millisecond), err.Error())
			return
		}
		defer res.Body.Close()

		if res.StatusCode != 200 {
			boomer.RecordFailure("http", "getSellInfo", elapsed.Nanoseconds()/int64(time.Millisecond), res.Status)
			return
		}
		boomer.RecordSuccess("http", "getSellInfo", elapsed.Nanoseconds()/int64(time.Millisecond), length)
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

	// add SellInfo
	// for i := 0; i < max; i++ {
	// 	bodyBuffer := &bytes.Buffer{}
	// 	bodyWriter := multipart.NewWriter(bodyBuffer)
	// 	userIDWriter, _ := bodyWriter.CreateFormField("userID")
	// 	_, _ = userIDWriter.Write([]byte("1"))
	// 	validTimeWriter, _ := bodyWriter.CreateFormField("validTime")
	// 	_, _ = validTimeWriter.Write([]byte("1609434061"))
	// 	goodNameWriter, _ := bodyWriter.CreateFormField("goodName")
	// 	_, _ = goodNameWriter.Write([]byte(randSeq(10)))

	// 	contentType := bodyWriter.FormDataContentType()
	// 	_ = bodyWriter.Close()

	// 	req, _ := http.NewRequest("POST", "http://10.0.0.28:30444/v1/sellInfo", bytes.NewReader(bodyBuffer.Bytes()))
	// 	req.Header.Add("Content-Type", contentType)
	// 	req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjI1OTI0ODY4OTksImlkIjoxLCJyb2xlIjoxfQ.i0YUDEcnJrdgLOMTMpsnp85mHEnJkueeqoWrQg-d_o0")
	// 	resp, err := c.Do(req)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	if resp.StatusCode != http.StatusOK {
	// 		panic(resp.StatusCode)
	// 	}
	// 	body, err := ioutil.ReadAll(resp.Body)
	// 	if err != nil {
	// 		fmt.Printf("read body err, %v\n", err)
	// 		return
	// 	}
	// 	println("json:", string(body))
	// }

	// start test
	task1 := &boomer.Task{
		Name:   "goodSearch",
		Weight: 10,
		Fn:     goodSearch,
	}

	boomer.Run(task1)
}
