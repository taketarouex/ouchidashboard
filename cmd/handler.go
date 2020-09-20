package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/labstack/echo"
	"github.com/tenntenn/natureremo"
	"github.com/tktkc72/ouchidashboard/collector"
	"github.com/tktkc72/ouchidashboard/enum"
	"github.com/tktkc72/ouchidashboard/ouchi"
	"github.com/tktkc72/ouchidashboard/repository"
)

func getLogsHandler(c echo.Context) error {
	projectID := os.Getenv("GCP_PROJECT")
	rootPath := os.Getenv("FIRESTORE_ROOT_PATH")

	ctx := context.Background()
	firestoreClient, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Printf("failed to create firestore client due to: %v", err)
		return echo.ErrInternalServerError
	}
	defer firestoreClient.Close()

	roomName := c.Param("roomName")
	repository, err := repository.NewRepository(firestoreClient, rootPath, roomName, &collector.NowTime{})
	if err != nil {
		log.Printf("create repository: %v", err)
		if ouchi.IsNoRoom(err) {
			return echo.ErrBadRequest
		}
		return echo.ErrInternalServerError
	}

	service := ouchi.NewOuchi(repository)

	logType, err := enum.ParseLogType(c.Param("logType"))
	if err != nil {
		log.Printf("failed to parse logtype: %s", c.Param("logType"))
		return echo.ErrBadRequest
	}

	start, end, err := parseStartEnd(c)
	if err != nil {
		return err
	}

	options, err := parseOptions(c)
	if err != nil {
		return err
	}

	logs, err := service.GetLogs(logType, start, end, options...)
	if err != nil {
		return echo.ErrInternalServerError
	}
	return c.JSON(http.StatusOK, logs)
}

func parseOptions(c echo.Context) ([]ouchi.GetOption, error) {
	options := []ouchi.GetOption{}
	limitParam := c.QueryParam("limit")
	if limitParam != "" {
		limit, err := strconv.Atoi(limitParam)
		if err != nil {
			log.Printf("err: %v, query parameter limit is expected to be int but got: %s", err, limitParam)
			return []ouchi.GetOption{}, echo.ErrBadRequest
		}
		options = append(options, ouchi.Limit(limit))
	}
	orderParam := c.QueryParam("order")
	if orderParam != "" {
		order, err := enum.ParseOrder(orderParam)
		if err != nil {
			log.Printf("err: %v, query parameter order is expected to be Order but got: %s", err, orderParam)
			return []ouchi.GetOption{}, echo.ErrBadRequest
		}
		options = append(options, ouchi.Order(order))
	}
	return options, nil
}

func parseStartEnd(c echo.Context) (time.Time, time.Time, error) {
	startParam := c.QueryParam("start")
	if startParam == "" {
		log.Print("query parameter start is needed")
		return time.Time{}, time.Time{}, echo.ErrBadRequest
	}
	start, err := time.Parse(time.RFC3339, startParam)
	if err != nil {
		log.Printf("start format is RFC3339 but got: %s", startParam)
		return time.Time{}, time.Time{}, echo.ErrBadRequest
	}
	endParam := c.QueryParam("end")
	if endParam == "" {
		log.Print("query parameter end is needed")
		return time.Time{}, time.Time{}, echo.ErrBadRequest
	}
	end, err := time.Parse(time.RFC3339, endParam)
	if err != nil {
		log.Printf("end format is RFC3339 but got: %s", endParam)
		return time.Time{}, time.Time{}, echo.ErrBadRequest
	}
	return start, end, nil
}

func collectorHandler(c echo.Context) error {
	accessToken := os.Getenv("NATURE_REMO_ACCESS_TOKEN")
	projectID := os.Getenv("GCP_PROJECT")
	rootPath := os.Getenv("FIRESTORE_ROOT_PATH")

	m := new(collector.Message)
	if err := c.Bind(m); err != nil {
		return err
	}
	errorChannel := make(chan error, len(m.RoomNames))
	for _, roomName := range m.RoomNames {
		go collect(accessToken, roomName, projectID, rootPath, errorChannel)
	}
	for range m.RoomNames {
		err := <-errorChannel
		if err != nil {
			log.Printf("collect: %v", err)
			if ouchi.IsNoRoom(err) {
				return echo.ErrBadRequest
			}
			return echo.ErrInternalServerError
		}
	}
	return nil
}

func collect(accessToken, roomName, projectID, rootPath string, c chan error) {
	natureremoClient := natureremo.NewClient(accessToken)
	fetcher := collector.NewFetcher(natureremoClient)

	ctx := context.Background()
	firestoreClient, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		c <- err
		return
	}
	defer firestoreClient.Close()
	repository, err := repository.NewRepository(firestoreClient, rootPath, roomName, &collector.NowTime{})
	if err != nil {
		c <- err
		return
	}

	service := collector.NewCollectorService(fetcher, repository)
	err = service.Collect()
	if err != nil {
		c <- err
		return
	}
	c <- nil
}
