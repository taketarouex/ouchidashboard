package collector

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

func TestCollector(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	collectedLog := CollectLog{
		historyLog{0, time.Date(2020, 7, 31, 0, 0, 0, 0, time.Local)},
		historyLog{1, time.Date(2020, 7, 31, 0, 0, 0, 0, time.Local)},
		historyLog{2, time.Date(2020, 7, 31, 0, 0, 0, 0, time.Local)},
		historyLog{3, time.Date(2020, 7, 31, 0, 0, 0, 0, time.Local)},
	}
	fetcher := NewMockIFetcher(ctrl)
	fetcher.EXPECT().fetch().Return(collectedLog, nil)

	repository := NewMockIRepository(ctrl)
	repository.EXPECT().add(collectedLog).Return(nil)
	service := NewCollectorService(fetcher, repository)
	service.Collect()
}
