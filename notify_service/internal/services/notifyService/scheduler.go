package notifyService

import (
	"context"
	"log"
	"time"
)

type Scheduler struct {
	service *NotifyService
	ticker  *time.Ticker
	done    chan bool
}

func NewScheduler(service *NotifyService, interval time.Duration) *Scheduler {
	return &Scheduler{
		service: service,
		ticker:  time.NewTicker(interval),
		done:    make(chan bool),
	}
}

func (s *Scheduler) Start(ctx context.Context) {
	log.Printf("Starting scheduler")

	// Выполняем сразу при старте
	s.service.CheckNewArticlesForAllSubscriptions(ctx)

	go func() {
		for {
			select {
			case <-s.ticker.C:
				s.service.CheckNewArticlesForAllSubscriptions(ctx)
			case <-s.done:
				return
			case <-ctx.Done():
				s.Stop()
				return
			}
		}
	}()
}

func (s *Scheduler) Stop() {
	log.Println("Stopping scheduler...")
	s.ticker.Stop()
	s.done <- true
	s.service.Close()
}
