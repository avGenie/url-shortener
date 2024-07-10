package converter

import (
	"github.com/avGenie/url-shortener/internal/app/models"

	pb "github.com/avGenie/url-shortener/proto"
)

// AllURLsBatchToAllURLsResponse Converts model AllUrlsBatch struct to proto AllUrlsResponse
func AllURLsBatchToAllURLsResponse(batch models.AllUrlsBatch) *pb.AllUrlsResponse {
	urlsResponse := make([]*pb.UrlsResponse, 0, len(batch))

	for _, urls := range batch {
		response := &pb.UrlsResponse{
			ShortURL: urls.ShortURL,
			OriginalURL: urls.OriginalURL,
		}

		urlsResponse = append(urlsResponse, response)
	}

	return &pb.AllUrlsResponse{
		Urls: urlsResponse,
	}
}

// BatchRequestToReqBatch Converts proto BatchRequest to model ReqBatch struct
func BatchRequestToReqBatch(request *pb.BatchRequest) models.ReqBatch {
	outBatch := make(models.ReqBatch, 0, len(request.Urls))

	for _, val := range request.GetUrls() {
		batch := models.BatchObjectRequest{
			ID: val.GetCorrelationID(),
			URL: val.GetOriginalURL(),
		}

		outBatch = append(outBatch, batch)
	}

	return outBatch
}

// ResBatchToBatchResponse Converts model ResBatch struct to proto BatchResponse
func ResBatchToBatchResponse(resBatch models.ResBatch) *pb.BatchResponse {
	outBatch := make([]*pb.BatchShortURLObject, 0, len(resBatch))

	for _, val := range resBatch {
		batch := &pb.BatchShortURLObject{
			CorrelationID: val.ID,
			ShortURL: val.URL,
		}

		outBatch = append(outBatch, batch)
	}

	return &pb.BatchResponse{
		Urls: outBatch,
	}
}

// DeleteRequestToReqDeletedURLBatch Converts proto DeleteRequest to model ReqDeletedURLBatch struct
func DeleteRequestToReqDeletedURLBatch(request *pb.DeleteRequest) models.ReqDeletedURLBatch {
	output := make(models.ReqDeletedURLBatch, 0, len(request.Urls))

	for _, val := range request.GetUrls() {
		output = append(output, val.GetShortURL())
	}

	return output
}
