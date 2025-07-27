package usecase

import (
	"errors"
	"fmt"
	"sync"

	"github.com/ribbinpo/mining-service/internal/application/domain"
	"github.com/ribbinpo/mining-service/internal/application/port"
)

type stoneUsecase struct {
	scraperRepository port.ScraperRepository
	stoneRepository   port.StoneRepository
}

func NewStoneUsecase(scraperRepository port.ScraperRepository, stoneRepository port.StoneRepository) port.StoneUsecase {
	return &stoneUsecase{scraperRepository: scraperRepository, stoneRepository: stoneRepository}
}

func (u *stoneUsecase) GetStoneByURL(url string) (*domain.StoneDomain, error) {
	stone, err := u.stoneRepository.GetStoneByURL(url)
	if err != nil {
		return nil, err
	}

	return stone, nil
}

func (u *stoneUsecase) RegisterDig(url string, refresh bool) error {
	// 1. found stone, is retry -> update stone
	// 2. found stone, is not retry -> throw error
	// 3. found stone, status is pending -> throw error
	// 4. not found stone -> create
	stone, err := u.GetStoneByURL(url)
	if err != nil {
		return err
	}
	// found stone
	if stone != nil {
		if stone.Status == domain.StoneStatusEnumPending {
			// throw error - stone is pending
			return errors.New("stone is pending")
		}
		if !refresh {
			// throw error - stone already exists
			return errors.New("stone already exists")
		} else {
			// update stone
			stone.Retry()
			err = u.stoneRepository.UpdateStone(stone)
			if err != nil {
				return err
			}
			return errors.New("stone already exists")
		}
	}

	// not found stone
	stone = domain.CreateStoneDomain(url)
	err = u.stoneRepository.CreateStone(stone)
	if err != nil {
		return err
	}

	return nil
}

func (u *stoneUsecase) ExecuteDigFromQueue() error {
	stones, err := u.stoneRepository.GetAllStones(&port.GetAllStonesOptions{
		Page:         1,
		PageSize:     1,
		StatusFilter: domain.StoneStatusEnumPending,
		SortBy:       "updated_at",
		SortOrder:    "asc",
	})
	if err != nil {
		return err
	}

	if len(stones) == 0 {
		return nil
	}

	stone := stones[0]

	stonePaths, err := u.ExecuteDig(stone.DomainURL)

	if err != nil {
		return err
	}

	for _, stonePath := range stonePaths {
		fmt.Println(stonePath)
	}

	stone.Paths = stonePaths
	stone.Status = domain.StoneStatusEnumCompleted

	u.stoneRepository.UpdateStone(stone)

	return nil
}

func (u *stoneUsecase) ExecuteDig(targetURL string) ([]string, error) {
	var (
		visited   = make(map[string]struct{})
		mutex     sync.Mutex
		queueURLs = make(chan string, 1000)
		wg        sync.WaitGroup // for workers only
	)

	// Use a separate wait group for producers (URLs being added)
	producerWg := sync.WaitGroup{}

	// Seed URL
	mutex.Lock()
	visited[targetURL] = struct{}{}
	mutex.Unlock()

	producerWg.Add(1)
	queueURLs <- targetURL

	// Start worker pool
	concurrency := 10
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for url := range queueURLs {
				fmt.Printf("[Worker %d] Scraping: %s\n", workerID, url)

				// Discover new URLs
				newURLs := u.scraperRepository.ScrapeURLs(url, true, false, true, true)

				for _, newURL := range newURLs {
					mutex.Lock()
					_, seen := visited[newURL]
					if !seen {
						visited[newURL] = struct{}{}
						producerWg.Add(1)
						queueURLs <- newURL
					}
					mutex.Unlock()
				}
				producerWg.Done()
			}
		}(i)
	}

	// Wait for all URL discovery to finish, then close queue
	go func() {
		producerWg.Wait()
		close(queueURLs) // All URLs added, safe to close
	}()

	// Wait for workers to finish
	wg.Wait()

	// Print all visited URLs
	visitedURLs := make([]string, 0, len(visited))
	for url := range visited {
		visitedURLs = append(visitedURLs, url)
	}

	return visitedURLs, nil
}
