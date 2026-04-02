package platform

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/birddigital/content-loop/internal/config"
	"github.com/birddigital/content-loop/internal/generator"
)

// Manager handles all social media platform integrations
type Manager struct {
	cfg      *config.Config
	clients  map[string]Publisher
}

// Publisher defines the interface for publishing to social platforms
type Publisher interface {
	Publish(ctx context.Context, post *generator.SocialPost) error
	Validate(post *generator.SocialPost) error
}

// NewManager creates a new platform Manager
func NewManager(cfg *config.Config) *Manager {
	m := &Manager{
		cfg:     cfg,
		clients: make(map[string]Publisher),
	}

	// Initialize enabled platforms
	if cfg.Social.LinkedIn.Enabled {
		m.clients["linkedin"] = NewLinkedInPublisher(&cfg.Social.LinkedIn)
	}
	if cfg.Social.Twitter.Enabled {
		m.clients["twitter"] = NewTwitterPublisher(&cfg.Social.Twitter)
	}
	if cfg.Social.YouTube.Enabled {
		m.clients["youtube"] = NewYouTubePublisher(&cfg.Social.YouTube)
	}
	if cfg.Social.Threads.Enabled {
		m.clients["threads"] = NewThreadsPublisher(&cfg.Social.Threads)
	}

	return m
}

// Publish posts content to the specified platform
func (m *Manager) Publish(ctx context.Context, post *generator.SocialPost) error {
	client, ok := m.clients[post.Platform]
	if !ok {
		return fmt.Errorf("no client for platform: %s", post.Platform)
	}

	// Validate the post
	if err := client.Validate(post); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Rate limiting: don't spam
	if err := m.rateLimit(post.Platform); err != nil {
		return err
	}

	// Publish
	if err := client.Publish(ctx, post); err != nil {
		return fmt.Errorf("publish failed: %w", err)
	}

	log.Printf("✅ Published to %s: %s", post.Platform, truncate(post.Content, 50))
	return nil
}

// rateLimit ensures we don't exceed platform limits
func (m *Manager) rateLimit(platform string) error {
	// TODO: Implement proper rate limiting with Redis/DB
	// For now, just sleep a bit
	time.Sleep(500 * time.Millisecond)
	return nil
}

// LinkedInPublisher publishes to LinkedIn
type LinkedInPublisher struct {
	cfg *config.LinkedInConfig
}

func NewLinkedInPublisher(cfg *config.LinkedInConfig) *LinkedInPublisher {
	return &LinkedInPublisher{cfg: cfg}
}

func (p *LinkedInPublisher) Validate(post *generator.SocialPost) error {
	if len(post.Content) > 3000 {
		return fmt.Errorf("post too long: %d chars (max 3000)", len(post.Content))
	}
	return nil
}

func (p *LinkedInPublisher) Publish(ctx context.Context, post *generator.SocialPost) error {
	// For now, just log (dry run mode)
	log.Printf("[LinkedIn - DRY RUN] Would publish: %s", truncate(post.Content, 100))

	// TODO: Implement actual LinkedIn API call
	// POST https://api.linkedin.com/v2/ugcPosts
	// with UGC post container
	return nil
}

// TwitterPublisher publishes to Twitter/X
type TwitterPublisher struct {
	cfg *config.TwitterConfig
}

func NewTwitterPublisher(cfg *config.TwitterConfig) *TwitterPublisher {
	return &TwitterPublisher{cfg: cfg}
}

func (p *TwitterPublisher) Validate(post *generator.SocialPost) error {
	if len(post.Content) > 280 {
		return fmt.Errorf("tweet too long: %d chars (max 280)", len(post.Content))
	}
	return nil
}

func (p *TwitterPublisher) Publish(ctx context.Context, post *generator.SocialPost) error {
	// For now, just log (dry run mode)
	log.Printf("[Twitter - DRY RUN] Would publish: %s", truncate(post.Content, 100))

	// TODO: Implement actual Twitter API v2 call
	// POST https://api.twitter.com/2/tweets
	return nil
}

// YouTubePublisher publishes to YouTube Community tab
type YouTubePublisher struct {
	cfg *config.YouTubeConfig
}

func NewYouTubePublisher(cfg *config.YouTubeConfig) *YouTubePublisher {
	return &YouTubePublisher{cfg: cfg}
}

func (p *YouTubePublisher) Validate(post *generator.SocialPost) error {
	if len(post.Content) > 5000 {
		return fmt.Errorf("community post too long: %d chars (max 5000)", len(post.Content))
	}
	return nil
}

func (p *YouTubePublisher) Publish(ctx context.Context, post *generator.SocialPost) error {
	// For now, just log (dry run mode)
	log.Printf("[YouTube - DRY RUN] Would publish to Community: %s", truncate(post.Content, 100))

	// TODO: Implement actual YouTube API call
	// POST https://www.googleapis.com/youtube/v3/comments
	return nil
}

// ThreadsPublisher publishes to Threads
type ThreadsPublisher struct {
	cfg *config.ThreadsConfig
}

func NewThreadsPublisher(cfg *config.ThreadsConfig) *ThreadsPublisher {
	return &ThreadsPublisher{cfg: cfg}
}

func (p *ThreadsPublisher) Validate(post *generator.SocialPost) error {
	if len(post.Content) > 500 {
		return fmt.Errorf("threads post too long: %d chars (max 500)", len(post.Content))
	}
	return nil
}

func (p *ThreadsPublisher) Publish(ctx context.Context, post *generator.SocialPost) error {
	// For now, just log (dry run mode)
	log.Printf("[Threads - DRY RUN] Would publish: %s", truncate(post.Content, 100))

	// TODO: Implement actual Threads API call
	// POST https://graph.threads.net/v1.0/me/threads
	return nil
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
