package converter

import (
	"fmt"

	"github.com/avGenie/url-shortener/internal/app/models"
	storage "github.com/avGenie/url-shortener/internal/app/storage/api/model"
)

func ConvertStorageBatchToOutBatch(batch storage.Batch, uriPrefix string) models.ResBatch {
	outBatch := make(models.ResBatch, 0, len(batch))
	for _, obj := range batch {
		if len(obj.ShortURL) == 0 {
			continue
		}

		outObj := models.BatchObjectResponse{
			ID:  obj.ID,
			URL: fmt.Sprintf("%s/%s", uriPrefix, obj.ShortURL),
		}

		outBatch = append(outBatch, outObj)
	}

	return outBatch
}
