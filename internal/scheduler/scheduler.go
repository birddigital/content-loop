package scheduler

import (
	"context"
	"log"
	"time"
)

// Config holds scheduler configuration
type Config struct {
	HourlyTick func(ctx context.Context, tick time.Time) error
}

// Scheduler runs content generation on a schedule
type Scheduler struct {
	cfg       *Config
	ticker    *time.Ticker
	running   bool
}

// New creates a new Scheduler
func New(cfg *Config) *Scheduler {
	return &Scheduler{
		cfg: cfg,
	}
}

// Run starts the scheduler loop
func (s *Scheduler) Run(ctx context.Context) error {
	log.Println("📅 Scheduler starting...")

	// Start hourly ticker (runs at the top of every hour)
	s.ticker = time.NewTicker(time.Hour)
	defer s.ticker.Stop()

	s.running = true

	// Run immediately on start
	if err := s.cfg.HourlyTick(ctx, time.Now()); err != nil {
		log.Printf("⚠️  Initial tick failed: %v", err)
	}

	// Then run every hour
	for {
		select {
		case <-ctx.Done():
			log.Println("📅 Scheduler stopped")
			return ctx.Err()

		case tick := <-s.ticker.C:
			log.Printf("⏰ Tick at %s", tick.Format(time.RFC3339))
			if err := s.cfg.HourlyTick(ctx, tick); err != nil {
				log.Printf("❌ Tick failed: %v", err)
			}
		}
	}
}

// Stop gracefully stops the scheduler
func (s *Scheduler) Stop() {
	s.running = false
	if s.ticker != nil {
		s.ticker.Stop()
	}
}

// Status returns the current scheduler status
func (s *Scheduler) Status() string {
	if s.running {
		return "running"
	}
	return "stopped"
}
