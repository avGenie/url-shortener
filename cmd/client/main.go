package main

import (
	"encoding/json"
	"fmt"
	"github.com/avGenie/url-shortener/cmd/client/client"
	"github.com/avGenie/url-shortener/cmd/client/random"
	"github.com/avGenie/url-shortener/internal/app/config"
	"github.com/avGenie/url-shortener/internal/app/logger"
	"github.com/avGenie/url-shortener/internal/app/models"
	"go.uber.org/zap"
	"net/http"
	"sync"
	"time"
)

const (
	MaxCount     = 1000
	BatchCount   = 3
	RoutineCount = 10
)

func main() {
	config := config.InitConfig()
	err := logger.Initialize(config)
	if err != nil {
		panic(err.Error())
	}

	config.NetAddr = fmt.Sprintf("http://%s", config.NetAddr)

	now := time.Now()
	wg := &sync.WaitGroup{}
	wg.Add(RoutineCount)
	for i := 0; i < RoutineCount; i++ {
		go testPostRequest(config, wg)
	}

	wg.Add(RoutineCount)
	for i := 0; i < RoutineCount; i++ {
		go testPostShortenRequest(config, wg)
	}

	wg.Add(RoutineCount)
	for i := 0; i < RoutineCount; i++ {
		go testPostShortenBatchRequest(config, wg)
	}

	wg.Wait()
	fmt.Printf("spended time: %s\n", time.Since(now))
}

func testPostRequest(config config.Config, wg *sync.WaitGroup) {
	defer wg.Done()
	postRequest(config)
}

func testPostShortenRequest(config config.Config, wg *sync.WaitGroup) {
	defer wg.Done()
	postShortenRequest(config)
}

func testPostShortenBatchRequest(config config.Config, wg *sync.WaitGroup) {
	defer wg.Done()
	postShortenBatchRequest(config)
}

func postRequest(config config.Config) {
	c := client.New(config.NetAddr)
	var cookie *http.Cookie

	for i := 0; i < MaxCount; i++ {
		time.Sleep(random.SleepDuration(10, 50) * time.Millisecond)
		res, err := c.SendPostRequest([]byte(random.GenerateRandomURL()), cookie)
		if err != nil {
			zap.L().Error("PostShortenRequest SendPostRequest", zap.Error(err))
			continue
		}
		defer res.Body.Close()

		if cookie == nil {
			cookies := res.Cookies()
			cookie = cookies[0]
		}
	}
}

func postShortenRequest(config config.Config) {
	netAddr := fmt.Sprintf("%s/api/shorten", config.NetAddr)
	c := client.New(netAddr)
	var cookie *http.Cookie

	for i := 0; i < MaxCount; i++ {
		time.Sleep(random.SleepDuration(10, 50) * time.Millisecond)
		request := models.Request{
			URL: random.GenerateRandomURL(),
		}

		data, err := json.Marshal(request)
		if err != nil {
			zap.L().Error("PostShortenRequest marshall", zap.Error(err))
			continue
		}

		res, err := c.SendPostRequest(data, cookie)
		if err != nil {
			zap.L().Error("PostShortenRequest SendPostRequest", zap.Error(err))
			continue
		}
		defer res.Body.Close()

		if cookie == nil {
			cookies := res.Cookies()
			cookie = cookies[0]
		}
	}
}

func postShortenBatchRequest(config config.Config) {
	netAddr := fmt.Sprintf("%s/api/shorten/batch", config.NetAddr)
	c := client.New(netAddr)
	var cookie *http.Cookie

	for i := 0; i < MaxCount; i++ {
		time.Sleep(random.SleepDuration(10, 50) * time.Millisecond)
		data, err := createBatchData()
		if err != nil {
			zap.L().Error("PostShortenRequest marshall", zap.Error(err))
			continue
		}

		res, err := c.SendPostRequest(data, cookie)
		if err != nil {
			zap.L().Error("PostShortenRequest SendPostRequest", zap.Error(err))
			continue
		}
		defer res.Body.Close()

		if cookie == nil {
			cookies := res.Cookies()
			cookie = cookies[0]
		}
	}
}

func createBatchData() ([]byte, error) {
	batch := make(models.ReqBatch, 0, BatchCount)
	for i := 0; i < BatchCount; i++ {
		randomString := random.GenerateRandomString()
		url := random.GenerateURL(randomString)
		batch = append(batch, models.BatchObjectRequest{
			ID:  randomString,
			URL: url,
		})
	}

	return json.Marshal(batch)
}
