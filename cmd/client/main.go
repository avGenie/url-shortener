package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/avGenie/url-shortener/cmd/client/client"
	"github.com/avGenie/url-shortener/cmd/client/random"
	"github.com/avGenie/url-shortener/internal/app/config"
	"github.com/avGenie/url-shortener/internal/app/entity"
	grpc_context "github.com/avGenie/url-shortener/internal/app/grpc/usecase/context"
	"github.com/avGenie/url-shortener/internal/app/logger"
	"github.com/avGenie/url-shortener/internal/app/models"
	pb "github.com/avGenie/url-shortener/proto"
	shortener "github.com/avGenie/url-shortener/proto"
)

const (
	maxCount     = 1000
	batchCount   = 3
	routineCount = 10
)

func main() {
	config, err := config.InitConfig()
	if err != nil {
		zap.L().Fatal("Failed to initialize config", zap.Error(err))
	}
	err = logger.Initialize(config)
	if err != nil {
		zap.L().Fatal("Failed to initialize logger", zap.Error(err))
	}

	config.NetAddr = fmt.Sprintf("http://%s", config.NetAddr)

	grpcTest(config)

	// stressTestHTTP(config)
}

func grpcTest(config config.Config) {
	conn, err := grpc.Dial(config.GRPCNetAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewShortenerClient(conn)

	// getOriginalGRPCURL(client, "be89c05e", "8c6c0dbc-22b8-4349-b33f-7204104bbd97")
	// getShortURL(client, "https://www.google.com", "8c6c0dbc-22b8-4349-b33f-7204104bbd97")
	// getAllURLs(client, "8c6c0dbc-22b8-4349-b33f-7204104bbd97")

	// urls := []*pb.BatchOriginalURLObject{
	// 	{
	// 		CorrelationID: "google1",
	// 		OriginalURL: "https://google1.com/",
	// 	},
	// 	{
	// 		CorrelationID: "google2",
	// 		OriginalURL: "https://google2.com/",
	// 	},
	// }
	// getBatchShortURL(client, "8c6c0dbc-22b8-4349-b33f-7204104bbd97", &pb.BatchRequest{Urls: urls})

	// deleteURLs := []*pb.DeleteObject{
	// 	{
	// 		ShortURL: "42b3e75f",
	// 	},
	// 	{
	// 		ShortURL: "77fca595",
	// 	},
	// 	{
	// 		ShortURL: "ac6bb669",
	// 	},
	// }
	// deleteShortURLs(client, "8c6c0dbc-22b8-4349-b33f-7204104bbd97", &pb.DeleteRequest{Urls: deleteURLs})

	getStatistic(client, "8c6c0dbc-22b8-4349-b33f-7204104bbd97")
}

func getStatistic(client shortener.ShortenerClient, userID string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ctx = grpc_context.SetUserIDContext(ctx, entity.UserID(userID))

	stat, err := client.GetStatistic(ctx, &emptypb.Empty{})
	if err != nil {
		zap.L().Error("getStatistic GetStatistic", zap.Error(err))

		return
	}

	fmt.Println(stat)
}

func deleteShortURLs(client shortener.ShortenerClient, userID string, request *pb.DeleteRequest) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ctx = grpc_context.SetUserIDContext(ctx, entity.UserID(userID))

	_, err := client.DeleteURLs(ctx, request)
	if err != nil {
		zap.L().Error("deleteShortURLs DeleteURLs", zap.Error(err))

		return
	}

	fmt.Println("no errors while deleting")
}

func getBatchShortURL(client shortener.ShortenerClient, userID string, req *pb.BatchRequest) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ctx = grpc_context.SetUserIDContext(ctx, entity.UserID(userID))

	urls, err := client.GetBatchShortURL(ctx, req)
	if err != nil {
		zap.L().Error("getAllUserURL GetAllUserURL", zap.Error(err))

		return
	}

	fmt.Println(urls)
}

func getAllURLs(client shortener.ShortenerClient, userID string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ctx = grpc_context.SetUserIDContext(ctx, entity.UserID(userID))

	urls, err := client.GetAllUserURL(ctx, &emptypb.Empty{})
	if err != nil {
		zap.L().Error("getAllUserURL GetAllUserURL", zap.Error(err))

		return
	}

	fmt.Println(urls)
}

func getShortURL(client shortener.ShortenerClient, url, userID string) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	ctx = grpc_context.SetUserIDContext(ctx, entity.UserID(userID))

	original, err := client.GetShortURL(ctx, &shortener.OriginalURL{Url: url})
	if err != nil {
		zap.L().Error("getShortURL GetShortURL", zap.Error(err))

		return
	}

	fmt.Println(original)
}

func getOriginalGRPCURL(client shortener.ShortenerClient, url, userID string) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	ctx = grpc_context.SetUserIDContext(ctx, entity.UserID(userID))

	original, err := client.GetOriginalURL(ctx, &shortener.ShortURL{Url: url})
	if err != nil {
		zap.L().Error("getOriginalGRPCURL GetOriginalURL", zap.Error(err))

		return
	}

	fmt.Println(original)
}

func stressTestHTTP(config config.Config) {
	now := time.Now()
	wg := &sync.WaitGroup{}
	wg.Add(routineCount)
	for i := 0; i < routineCount; i++ {
		go testPostRequest(config, wg)
	}

	wg.Add(routineCount)
	for i := 0; i < routineCount; i++ {
		go testPostShortenRequest(config, wg)
	}

	wg.Add(routineCount)
	for i := 0; i < routineCount; i++ {
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
	c := client.NewClient(config.NetAddr)
	var cookie *http.Cookie

	for i := 0; i < maxCount; i++ {
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
	c := client.NewClient(netAddr)
	var cookie *http.Cookie

	for i := 0; i < maxCount; i++ {
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
	c := client.NewClient(netAddr)
	var cookie *http.Cookie

	for i := 0; i < maxCount; i++ {
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
	batch := make(models.ReqBatch, 0, batchCount)
	for i := 0; i < batchCount; i++ {
		randomString := random.GenerateRandomString()
		url := random.GenerateURL(randomString)
		batch = append(batch, models.BatchObjectRequest{
			ID:  randomString,
			URL: url,
		})
	}

	return json.Marshal(batch)
}
