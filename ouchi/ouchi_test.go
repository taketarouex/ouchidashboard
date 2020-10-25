package ouchi

import (
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	"github.com/tktkc72/ouchidashboard/enum"
)

func TestOuchi_GetLogs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedLogs := []Log{
		{30, time.Date(2020, 7, 31, 0, 0, 0, 0, time.Local), time.Date(2020, 7, 31, 0, 0, 0, 0, time.Local)},
		{31, time.Date(2020, 7, 31, 1, 0, 0, 0, time.Local), time.Date(2020, 7, 31, 1, 0, 0, 0, time.Local)},
		{32, time.Date(2020, 7, 31, 2, 0, 0, 0, time.Local), time.Date(2020, 7, 31, 2, 0, 0, 0, time.Local)},
		{33, time.Date(2020, 7, 31, 3, 0, 0, 0, time.Local), time.Date(2020, 7, 31, 3, 0, 0, 0, time.Local)},
	}
	start := time.Date(2020, 7, 31, 0, 0, 0, 0, time.Local)
	end := time.Date(2020, 7, 31, 10, 0, 0, 0, time.Local)
	repository := NewMockIRepository(ctrl)
	service := NewOuchi(repository)

	t.Run("success no option", func(t *testing.T) {
		repository.EXPECT().Fetch("test", enum.Temperature, start, end, 0, enum.Asc).Return(expectedLogs, nil)
		logs, err := service.GetLogs("test", enum.Temperature, start, end)
		if err != nil {
			t.Errorf("failed to get temperature log, due to: %v", err)
		}
		if !reflect.DeepEqual(expectedLogs, logs) {
			t.Errorf("expect: %v, but got: %v", expectedLogs, logs)
		}
	})
	t.Run("success set limit", func(t *testing.T) {
		repository.EXPECT().Fetch("test", enum.Temperature, start, end, 3, enum.Asc).Return(expectedLogs[0:2], nil)
		logs, err := service.GetLogs("test", enum.Temperature, start, end, Limit(3))
		if err != nil {
			t.Errorf("failed to get temperature log, due to: %v", err)
		}
		if !reflect.DeepEqual(expectedLogs[0:2], logs) {
			t.Errorf("expect: %v, but got: %v", expectedLogs[0:2], logs)
		}
	})
	t.Run("success set order", func(t *testing.T) {
		reversedLogs := []Log{
			{33, time.Date(2020, 7, 31, 3, 0, 0, 0, time.Local), time.Date(2020, 7, 31, 3, 0, 0, 0, time.Local)},
			{32, time.Date(2020, 7, 31, 2, 0, 0, 0, time.Local), time.Date(2020, 7, 31, 2, 0, 0, 0, time.Local)},
			{31, time.Date(2020, 7, 31, 1, 0, 0, 0, time.Local), time.Date(2020, 7, 31, 1, 0, 0, 0, time.Local)},
			{30, time.Date(2020, 7, 31, 0, 0, 0, 0, time.Local), time.Date(2020, 7, 31, 0, 0, 0, 0, time.Local)},
		}
		repository.EXPECT().Fetch("test", enum.Temperature, start, end, 0, enum.Desc).Return(reversedLogs, nil)
		logs, err := service.GetLogs("test", enum.Temperature, start, end, Order(enum.Desc))
		if err != nil {
			t.Errorf("failed to get temperature log, due to: %v", err)
		}
		if !reflect.DeepEqual(reversedLogs, logs) {
			t.Errorf("expect: %v, but got: %v", reversedLogs, logs)
		}
	})
	t.Run("fail to Fetch", func(t *testing.T) {
		repository.EXPECT().Fetch("test", enum.Temperature, start, end, 0, enum.Asc).Return(nil, errors.New("failed to fetch"))
		if _, err := service.GetLogs("test", enum.Temperature, start, end); err == nil {
			t.Error("expect error but nil")
		}
	})
}

func TestOuchi_GetRoomNames(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repository := NewMockIRepository(ctrl)
	service := NewOuchi(repository)

	t.Run("success to fetch", func(t *testing.T) {
		expected := []string{"hoge", "fuga", "bar"}
		repository.EXPECT().FetchRoomNames().Return(expected, nil)
		roomNames, err := service.GetRoomNames()
		if err != nil {
			t.Errorf("failed to get room names, due to: %v", err)
		}
		if !cmp.Equal(expected, roomNames) {
			t.Errorf("expect: %v, but got: %v", expected, roomNames)
		}
	})
	t.Run("fail to fetch", func(t *testing.T) {
		repository.EXPECT().FetchRoomNames().Return([]string{}, errors.New("failed to fetch"))
		if _, err := service.GetRoomNames(); err == nil {
			t.Error("expect error but nil")
		}
	})
}
