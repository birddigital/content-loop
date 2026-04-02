package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"time"

	"github.com/birddigital/content-loop/internal/config"
	"github.com/birddigital/content-loop/internal/generator"
	"github.com/birddigital/content-loop/internal/platform"
	"github.com/birddigital/content-loop/internal/scheduler"
)

const version = "1.0.0"

// Content categories covering all expertise areas
var contentCategories = []config.Category{
	{
		Name: "Agent Swarms",
		Topics: []string{
			"deep-research-swarm", "deep-gap-swarm", "deep-explorer-swarm",
			"deep-security-swarm", "deep-compliance-swarm", "deep-threat-swarm",
			"deep-dependency-swarm", "deep-supply-chain-swarm",
		},
		Tags:       []string{"agent-swarm", "automation", "ai"},
		Cadence:    6, // Every 6 hours
		Priority:   1,
	},
	{
		Name: "Golang Development",
		Topics: []string{
			"go-performance-profiling", "go-concurrency-patterns",
			"go-microservices-architecture", "go-memory-optimization",
			"go-grpc-best-practices", "go-testing-patterns",
		},
		Tags:     []string{"golang", "performance", "engineering"},
		Cadence:  4,
		Priority: 2,
	},
	{
		Name: "Cloud & DevOps",
		Topics: []string{
			"cloud-cost-optimization", "kubernetes-security",
			"infrastructure-as-code", "observability-strategy",
			"disaster-recovery-planning", "multi-cloud-architecture",
		},
		Tags:     []string{"devops", "cloud", "aws", "kubernetes"},
		Cadence:  5,
		Priority: 2,
	},
	{
		Name: "Security & Compliance",
		Topics: []string{
			"automated-security-audits", "gdpr-compliance-automation",
			"threat-modeling-automation", "supply-chain-security",
			"patent-prior-art-search", "hipaa-compliance",
		},
		Tags:     []string{"security", "compliance", "automation"},
		Cadence:  8,
		Priority: 3,
	},
	{
		Name: "Digital Marketing",
		Topics: []string{
			"technical-seo-automation", "ai-marketing-stacks",
			"marketing-analytics-automation", "content-seo-strategy",
			"programmatic-seo", "marketing-attribution",
		},
		Tags:     []string{"marketing", "seo", "automation"},
		Cadence:  4,
		Priority: 2,
	},
	{
		Name: "Growth Hacking",
		Topics: []string{
			"viral-loop-implementation", "b2b-saas-growth",
			"product-led-growth", "referral-automation",
			"onboarding-optimization", "activation-rate-improvement",
		},
		Tags:     []string{"growth", "saas", "product"},
		Cadence:  6,
		Priority: 2,
	},
	{
		Name: "YouTube Strategy",
		Topics: []string{
			"youtube-algorithm-2026", "technical-channel-growth",
			"video-seo-optimization", "youtube-analytics-automation",
			"thumbnail-automation", "content-calendar-strategy",
		},
		Tags:     []string{"youtube", "video", "growth"},
		Cadence:  12,
		Priority: 3,
	},
	{
		Name: "Product Development",
		Topics: []string{
			"mvp-development-timeline", "product-market-fit-validation",
			"agile-development-automation", "feature-prioritization",
			"user-research-automation", "prototyping-workflows",
		},
		Tags:     []string{"product", "mvp", "development"},
		Cadence:  5,
		Priority: 2,
	},
	{
		Name: "Case Studies",
		Topics: []string{
			"monolith-to-microservices", "offline-first-mobile",
			"real-time-analytics-platform", "fintech-security",
			"healthcare-compliance", "e-commerce-scalability",
		},
		Tags:     []string{"case-study", "architecture"},
		Cadence:  24, // Daily
		Priority: 1,
	},
}

// Social platforms with their content formats
var socialPlatforms = []config.Platform{
	{
		Name:    "linkedin",
		Enabled: true,
		Formats: []string{"long-form", "carousel", "poll"},
		Limits: config.PostLimits{
			MaxChars:    3000,
			MaxHashtags: 10,
			MaxEmojis:   5,
		},
		OptimalTimes: []string{"08:00", "12:00", "17:00"},
	},
	{
		Name:    "twitter",
		Enabled: true,
		Formats: []string{"thread", "short", "image"},
		Limits: config.PostLimits{
			MaxChars:    280,
			MaxHashtags: 3,
			MaxEmojis:   2,
		},
		OptimalTimes: []string{"09:00", "15:00", "20:00"},
	},
	{
		Name:    "youtube",
		Enabled: true,
		Formats: []string{"video-script", "shorts-script", "community"},
		Limits: config.PostLimits{
			MaxChars:    5000,
			MaxHashtags: 15,
			MaxEmojis:   10,
		},
		OptimalTimes: []string{"10:00", "14:00", "18:00"},
	},
	{
		Name:    "threads",
		Enabled: true,
		Formats: []string{"short", "carousel", "video"},
		Limits: config.PostLimits{
			MaxChars:    500,
			MaxHashtags: 5,
			MaxEmojis:   3,
		},
		OptimalTimes: []string{"07:00", "19:00", "22:00"},
	},
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.Printf("🚀 ContentLoop v%s starting...", version)

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize generator with all categories
	gen := generator.New(&generator.Config{
		Categories:      contentCategories,
		Platforms:       socialPlatforms,
		Author:          "Adrien Bird",
		BlogURL:         "https://adrienbird.net/blog",
		DefaultYear:     2026,
		MinPostInterval: 1 * time.Hour,
	})

	// Initialize platform managers
	platforms := platform.NewManager(cfg)

	// Initialize scheduler
	sched := scheduler.New(&scheduler.Config{
		HourlyTick: func(ctx context.Context, tick time.Time) error {
			return hourlyLoop(ctx, gen, platforms, cfg, tick)
		},
	})

	// Run the content loop
	if err := sched.Run(context.Background()); err != nil {
		log.Fatalf("Scheduler failed: %v", err)
	}
}

// hourlyLoop runs the complete content generation and distribution loop
func hourlyLoop(ctx context.Context, gen *generator.Generator, platforms *platform.Manager, cfg *config.Config, tick time.Time) error {
	log.Printf("⏰ Hourly tick: %s", tick.Format(time.RFC3339))

	// 1. Determine what content to generate this hour
	contentPlan := gen.PlanContent(tick)

	// 2. Generate blog post if scheduled
	if contentPlan.GenerateBlog {
		blogPost, err := gen.GenerateBlogPost(ctx, contentPlan.Topic)
		if err != nil {
			log.Printf("❌ Failed to generate blog: %v", err)
			return err
		}

		// 3. Deploy blog post
		if err := deployBlogPost(blogPost); err != nil {
			log.Printf("❌ Failed to deploy blog: %v", err)
			return err
		}

		log.Printf("✅ Blog post published: %s", blogPost.Slug)

		// 4. Generate social media posts for this blog
		socialPosts := gen.GenerateSocialPosts(blogPost)

		// 5. Distribute to platforms
		for _, post := range socialPosts {
			if err := platforms.Publish(ctx, post); err != nil {
				log.Printf("❌ Failed to publish to %s: %v", post.Platform, err)
				continue
			}
			log.Printf("✅ Published to %s: %s", post.Platform, truncate(post.Content, 50))
		}
	}

	// 6. Generate standalone social content (non-blog)
	if contentPlan.GenerateSocial {
		standalonePosts := gen.GenerateStandaloneSocial(tick)
		for _, post := range standalonePosts {
			if err := platforms.Publish(ctx, post); err != nil {
				log.Printf("❌ Failed to publish standalone: %v", err)
				continue
			}
			log.Printf("✅ Standalone social published to %s", post.Platform)
		}
	}

	// 7. Generate and post engagement content (replies, quotes)
	if contentPlan.GenerateEngagement {
		engagementPosts := gen.GenerateEngagement(tick)
		for _, post := range engagementPosts {
			if err := platforms.Publish(ctx, post); err != nil {
				log.Printf("❌ Failed to publish engagement: %v", err)
				continue
			}
		}
	}

	log.Printf("📊 Hourly loop complete. Generated: blog=%v, social=%v, engagement=%v",
		contentPlan.GenerateBlog,
		contentPlan.GenerateSocial,
		contentPlan.GenerateEngagement)

	return nil
}

// deployBlogPost copies the blog post to the adrienbird.net site
func deployBlogPost(post *generator.BlogPost) error {
	// Write to /tmp first
	tmpPath := fmt.Sprintf("/tmp/blog-%s.md", post.Slug)
	if err := os.WriteFile(tmpPath, []byte(post.Content), 0644); err != nil {
		return fmt.Errorf("write temp file: %w", err)
	}

	// Copy to bird-m1 via SSH
	cmd := fmt.Sprintf("scp %s bird@m1.local:~/sources/adrienbird/src/app/blog/posts/", tmpPath)
	if err := runCommand(cmd); err != nil {
		return fmt.Errorf("scp to bird-m1: %w", err)
	}

	// Clean up temp file
	os.Remove(tmpPath)

	return nil
}

func runCommand(cmd string) error {
	log.Printf("🔧 Running: %s", cmd)
	// Execute the command using os/exec
	return exec.Command("sh", "-c", cmd).Run()
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
