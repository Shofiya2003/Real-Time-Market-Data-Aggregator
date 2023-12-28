package main

import (
	"context"
	"fmt"
	"real-time-aggregator/scraper"

	"time"

	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/log"
	"github.com/RichardKnop/machinery/v1/tasks"
	opentracing "github.com/opentracing/opentracing-go"
	// opentracing_log "github.com/opentracing/opentracing-go/log"
)

func main() {
	worker()
	send()
}

func startServer() (*machinery.Server, error) {
	cnf := &config.Config{
		DefaultQueue:    "machinery_tasks",
		ResultsExpireIn: 3600,
		Broker:          "redis://localhost:6379",
		ResultBackend:   "redis://localhost:6379",
		Redis: &config.RedisConfig{
			MaxIdle:                3,
			IdleTimeout:            240,
			ReadTimeout:            15,
			WriteTimeout:           15,
			ConnectTimeout:         15,
			NormalTasksPollPeriod:  1000,
			DelayedTasksPollPeriod: 500,
		},
	}

	// the machinery server stores the config info and the registered tasks
	server, err := machinery.NewServer(cnf)

	if err != nil {
		return nil, err
	}

	// Register tasks
	tasks := map[string]interface{}{
		"scrape_market_cap": scraper.Marketcap_Scraper{Url: "https://coinmarketcap.com/all/views/all/"}.Scrape,
	}

	return server, server.RegisterTasks(tasks)
}

// worker is used to consume the registered tasks, to start a worker we need a server with registered tasks
func worker() error {

	// get the server with registered tasks
	server, err := startServer()

	if err != nil {
		return err
	}
	worker := server.NewWorker("machinery_queue", 0)

	errHandler := func(err error) {
		log.ERROR.Println("I am error handler ", err)
	}

	postTaskHandler := func(signature *tasks.Signature) {
		log.INFO.Println("I am in post task handler: ", signature.Name)
	}

	preTaskHandler := func(signature *tasks.Signature) {
		log.INFO.Println("I am in pre task handler: ", signature.Name)
	}

	worker.SetErrorHandler(errHandler)
	worker.SetPostTaskHandler(postTaskHandler)
	worker.SetPreTaskHandler(preTaskHandler)

	return worker.Launch()
}

func send() error {

	server, err := startServer()
	if err != nil {
		return err
	}
	var (
		scrapeMarketCap tasks.Signature
	)

	var initTasks = func() {
		scrapeMarketCap = tasks.Signature{
			Name: "scrape_market_cap",
		}
	}

	span, ctx := opentracing.StartSpanFromContext(context.Background(), "send")
	defer span.Finish()

	initTasks()

	group, err := tasks.NewGroup(&scrapeMarketCap)
	if err != nil {
		return fmt.Errorf("Error creating group: %s", err.Error())
	}

	asyncResults, err := server.SendGroupWithContext(ctx, group, 10)
	if err != nil {
		return fmt.Errorf("Could not send group: %s", err.Error())
	}

	for _, asyncResult := range asyncResults {
		results, err := asyncResult.Get(time.Duration(time.Millisecond * 5))
		if err != nil {
			return fmt.Errorf("Getting task result failed with error: %s", err.Error())
		}

		log.INFO.Printf(
			"%v\n",
			tasks.HumanReadableResults(results),
		)

	}
	return nil
}
