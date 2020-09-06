package collector

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

func TestCollector_Collect(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	collectLogs := []collectLog{
		{0, time.Date(2020, 7, 31, 0, 0, 0, 0, time.Local), temperature, "test"},
		{1, time.Date(2020, 7, 31, 0, 0, 0, 0, time.Local), humidity, "test"},
		{2, time.Date(2020, 7, 31, 0, 0, 0, 0, time.Local), illumination, "test"},
		{3, time.Date(2020, 7, 31, 0, 0, 0, 0, time.Local), motion, "test"},
	}
	fetcher := NewMockIFetcher(ctrl)
	fetcher.EXPECT().fetch("testID").Return(collectLogs, nil)

	repository := NewMockIRepository(ctrl)
	repository.EXPECT().sourceID().Return("testID", nil)
	repository.EXPECT().add(collectLogs).Return(nil)
	service := NewCollectorService(fetcher, repository)
	service.Collect()
}
