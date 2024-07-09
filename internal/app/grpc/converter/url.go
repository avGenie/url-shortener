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