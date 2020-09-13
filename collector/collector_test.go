package collector

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/tktkc72/ouchi/enum"
)

func TestCollector_Collect(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("success", func(t *testing.T) {
		collectLogs := []CollectLog{
			{0, time.Date(2020, 7, 31, 0, 0, 0, 0, time.Local), enum.Temperature, "test"},
			{1, time.Date(2020, 7, 31, 0, 0, 0, 0, time.Local), enum.Humidity, "test"},
			{2, time.Date(2020, 7, 31, 0, 0, 0, 0, time.Local), enum.Illumination, "test"},
			{3, time.Date(2020, 7, 31, 0, 0, 0, 0, time.Local), enum.Motion, "test"},
		}
		fetcher := NewMockIFetcher(ctrl)
		fetcher.EXPECT().fetch("testID").Return(collectLogs, nil)

		repository := NewMockIRepository(ctrl)
		repository.EXPECT().SourceID().Return("testID", nil)
		repository.EXPECT().Add(collectLogs).Return(nil)
		service := NewCollectorService(fetcher, repository)
		if err := service.Collect(); err != nil {
			t.Error("fail to collect")
		}
	})
	t.Run("error get sourceID", func(t *testing.T) {
		fetcher := NewMockIFetcher(ctrl)
		repository := NewMockIRepository(ctrl)
		repository.EXPECT().SourceID().Return("", errors.Errorf("fail to get sourceID"))
		service := NewCollectorService(fetcher, repository)
		if err := service.Collect(); err == nil || err.Error() != "fail to get sourceID" {
			t.Errorf("expect fail to get sourceID but err: %v", err)
		}
	})
	t.Run("error fetch", func(t *testing.T) {
		fetcher := NewMockIFetcher(ctrl)
		fetcher.EXPECT().fetch("testID").Return(nil, errors.Errorf("fail to fetch"))
		repository := NewMockIRepository(ctrl)
		repository.EXPECT().SourceID().Return("testID", nil)
		service := NewCollectorService(fetcher, repository)
		if err := service.Collect(); err == nil || err.Error() != "fail to fetch" {
			t.Errorf("expect fail to fetch but err: %v", err)
		}
	})
	t.Run("error add", func(t *testing.T) {
		collectLogs := []CollectLog{
			{0, time.Date(2020, 7, 31, 0, 0, 0, 0, time.Local), enum.Temperature, "test"},
			{1, time.Date(2020, 7, 31, 0, 0, 0, 0, time.Local), enum.Humidity, "test"},
			{2, time.Date(2020, 7, 31, 0, 0, 0, 0, time.Local), enum.Illumination, "test"},
			{3, time.Date(2020, 7, 31, 0, 0, 0, 0, time.Local), enum.Motion, "test"},
		}
		fetcher := NewMockIFetcher(ctrl)
		fetcher.EXPECT().fetch("testID").Return(collectLogs, nil)
		repository := NewMockIRepository(ctrl)
		repository.EXPECT().SourceID().Return("testID", nil)
		repository.EXPECT().Add(collectLogs).Return(errors.Errorf("fail to add"))
		service := NewCollectorService(fetcher, repository)
		if err := service.Collect(); err == nil || err.Error() != "fail to add" {
			t.Errorf("expect fail to add but err: %v", err)
		}
	})
}
