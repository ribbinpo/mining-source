package main

import (
	// "sync"

	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/ribbinpo/mining-service/internal/application/usecase"
	"github.com/ribbinpo/mining-service/internal/config"
	"github.com/ribbinpo/mining-service/internal/frameworks/database"
	repository "github.com/ribbinpo/mining-service/internal/frameworks/database/repository"
	"github.com/ribbinpo/mining-service/internal/frameworks/webscraping"
)

func main() {
	config := config.NewConfig("./.env")

	databaseInstance := database.NewDatabase(config.Db.Dsn)
	db := databaseInstance.Connect()
	scraperRepository := webscraping.NewScraperRepository()
	stoneRepository := repository.NewStoneRepository(db)
	stoneUsecase := usecase.NewStoneUsecase(scraperRepository, stoneRepository)

	s, err := gocron.NewScheduler()
	if err != nil {
		panic(fmt.Errorf("failed to initialize scheduler: %w", err))
	}

	messageChannel := make(chan string)

	j, err := s.NewJob(
		gocron.DurationJob(
			60*time.Second,
		),
		gocron.NewTask(
			func() {
				start := time.Now()

				stoneUsecase.ExecuteDigFromQueue()

				elapsed := time.Since(start)
				messageChannel <- fmt.Sprintf("Time taken: %s\n", elapsed)
			},
		),
	)
	if err != nil {
		panic(fmt.Errorf("failed to create job: %w", err))
	}

	fmt.Printf("Job: %v\n", j.ID())

	go s.Start()

	go func() {
		for message := range messageChannel {
			fmt.Println(message)
		}
	}()

	// Wait for interrupt (Ctrl+C) to gracefully shut down
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig

	fmt.Println("Shutting down...")
	s.Shutdown()
}
