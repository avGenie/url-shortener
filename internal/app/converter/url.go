package converter

import (
	"github.com/avGenie/url-shortener/internal/app/entity"
	"github.com/avGenie/url-shortener/internal/app/models"
)

func ConvertBatchObjectReqToURL(obj models.BatchObjectRequest) (*entity.URL, error) {
	url, err := entity.NewURL(obj.URL)
	if err != nil {
		return nil, err
	}

	return url, nil
}

func ConvertBatchReqToURL(batch models.ReqBatch) (models.ReqURLBatch, error) {
	urls := make(models.ReqURLBatch, 0, len(batch))

	for _, obj := range batch {
		url, err := ConvertBatchObjectReqToURL(obj)
		if err != nil {
			return nil, err
		}

		outObj := models.ReqURLBatchObject{
			Obj: obj,
			URL: *url,
		}

		urls = append(urls, outObj)
	}

	return urls, nil
}
