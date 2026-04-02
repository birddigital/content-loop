package generator

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/birddigital/content-loop/internal/config"
)

// Config holds generator configuration
type Config struct {
	Categories      []config.Category
	Platforms       []config.Platform
	Author          string
	BlogURL         string
	DefaultYear     int
	MinPostInterval time.Duration
}

// Generator creates content for blog and social media
type Generator struct {
	cfg         *Config
	categoryIdx int
	topicIdx    map[string]int
}

// New creates a new Generator
func New(cfg *Config) *Generator {
	topicIdx := make(map[string]int)
	for _, cat := range cfg.Categories {
		topicIdx[cat.Name] = 0
	}

	return &Generator{
		cfg:      cfg,
		topicIdx: topicIdx,
	}
}

// ContentPlan defines what content to generate in a cycle
type ContentPlan struct {
	GenerateBlog    bool
	GenerateSocial  bool
	GenerateEngagement bool
	Category        string
	Topic           string
}

// PlanContent determines what content should be generated this tick
func (g *Generator) PlanContent(tick time.Time) *ContentPlan {
	hour := tick.Hour()
	minute := tick.Minute()

	plan := &ContentPlan{
		GenerateSocial:     true,  // Always generate social content
		GenerateEngagement: true,  // Always generate engagement
	}

	// Blog posts every 3 hours at :00
	if minute == 0 && hour%3 == 0 {
		plan.GenerateBlog = true
		plan.Category = g.selectCategory()
		plan.Topic = g.selectTopic(plan.Category)
	}

	return plan
}

// selectCategory chooses the next category based on priority and rotation
func (g *Generator) selectCategory() string {
	// Find categories with highest priority (lowest number)
	minPriority := 999
	var candidates []string

	for _, cat := range g.cfg.Categories {
		if cat.Priority < minPriority {
			minPriority = cat.Priority
			candidates = []string{cat.Name}
		} else if cat.Priority == minPriority {
			candidates = append(candidates, cat.Name)
		}
	}

	// Rotate through candidates
	if len(candidates) > 0 {
		selected := candidates[g.categoryIdx%len(candidates)]
		g.categoryIdx++
		return selected
	}

	return g.cfg.Categories[0].Name
}

// selectTopic chooses the next topic from a category
func (g *Generator) selectTopic(categoryName string) string {
	for _, cat := range g.cfg.Categories {
		if cat.Name == categoryName && len(cat.Topics) > 0 {
			topic := cat.Topics[g.topicIdx[categoryName]%len(cat.Topics)]
			g.topicIdx[categoryName]++
			return topic
		}
	}
	return "general"
}

// BlogPost represents a generated blog post
type BlogPost struct {
	Title       string
	Description string
	Slug        string
	Date        time.Time
	Category    string
	Tags        []string
	Content     string
}

// GenerateBlogPost generates a complete blog post
func (g *Generator) GenerateBlogPost(ctx context.Context, topic string) (*BlogPost, error) {
	// Find category for this topic
	var category config.Category
	for _, cat := range g.cfg.Categories {
		for _, t := range cat.Topics {
			if t == topic {
				category = cat
				break
			}
		}
	}

	now := time.Now()
	slug := fmt.Sprintf("%s-%s", topic, now.Format("20060102-1504"))

	post := &BlogPost{
		Title:       g.generateTitle(topic),
		Description: g.generateDescription(topic),
		Slug:        slug,
		Date:        now,
		Category:    category.Name,
		Tags:        append([]string{topic}, category.Tags...),
		Content:     g.generateBlogContent(topic, category),
	}

	return post, nil
}

// generateTitle creates a compelling blog title
func (g *Generator) generateTitle(topic string) string {
	titles := map[string]string{
		"deep-research-swarm":          "Building Research Corpora with Agent Swarms",
		"cloud-cost-optimization":      "Cloud Cost Optimization: The Complete Guide",
		"youtube-algorithm-2026":       "YouTube Algorithm 2026: What Actually Works",
		"viral-loop-implementation":    "Building Viral Loops: A Technical Guide",
		"mvp-development-timeline":     "MVP Development: From Idea to Launch in 4 Weeks",
		"go-concurrency-patterns":      "Go Concurrency Patterns: Beyond the Basics",
		"technical-seo-automation":     "Technical SEO Automation with AI Agent Swarms",
		"product-market-fit-validation": "Validating Product-Market Fit: A Data-Driven Approach",
		"b2b-saas-growth":              "B2B SaaS Growth: The Complete Playbook",
	}

	if title, ok := titles[topic]; ok {
		return title
	}

	// Default title generation
	return strings.Title(strings.ReplaceAll(topic, "-", " "))
}

// generateDescription creates SEO-friendly description
func (g *Generator) generateDescription(topic string) string {
	return fmt.Sprintf("How %s works in practice. Real examples, code samples, and implementation strategies for modern development.", topic)
}

// generateBlogContent generates the full blog post markdown
func (g *Generator) generateBlogContent(topic string, category config.Category) string {
	now := time.Now()

	return fmt.Sprintf(`---
title: "%s"
description: "%s"
date: %s
slug: %s
tags: %s
category: "%s"
author: "%s"
---

# %s

%s

## The Challenge

Every %s project starts with the same problem: too much complexity, not enough clarity.

## The Solution

%s

## Implementation

%s

## Results

| Metric | Before | After |
|--------|--------|-------|
| Performance | baseline | 2-5x improvement |
| Time to market | weeks | days |
| Code quality | technical debt | maintainable |

## What's Next

%s

---

**Need help with %s?** We specialize in building production-grade solutions. [Get in touch](/contact) to learn more.

*Tags: %s*
`,
		g.generateTitle(topic),
		g.generateDescription(topic),
		now.Format("2006-01-02"),
		topic,
		fmt.Sprintf(`["%s"]`, strings.Join(append([]string{topic}, category.Tags...), `", "`)),
		category.Name,
		g.cfg.Author,
		g.generateTitle(topic),
		g.generateIntroduction(topic),
		strings.Title(strings.ReplaceAll(topic, "-", " ")),
		g.generateSolution(topic),
		g.generateCodeExample(topic),
		g.generateNextSteps(topic),
		topic,
		strings.Join(append([]string{topic}, category.Tags...), ", "),
	)
}

// generateIntroduction creates the post introduction
func (g *Generator) generateIntroduction(topic string) string {
	intros := map[string]string{
		"deep-research-swarm":      "Research is the foundation of good engineering decisions. But comprehensive research takes time—time most projects don't have. Agent swarms can perform parallel research across multiple facets simultaneously, building permanent knowledge corpora in hours instead of weeks.",
		"cloud-cost-optimization":  "Cloud bills spiral out of control silently. Over-provisioned instances, forgotten resources, inefficient architectures—these hidden costs add up to thousands per month. Automated cost analysis identifies and fixes these issues before they drain your budget.",
		"youtube-algorithm-2026":   "The YouTube algorithm of 2026 rewards consistency, engagement, and watch time—but the specifics are always shifting. Technical channels face unique challenges: complex topics, longer production cycles, and a narrower audience. The right strategy cuts through the noise.",
		"viral-loop-implementation": "Viral loops are the holy grail of growth. Each user brings in more users automatically, creating compounding growth that scales without ad spend. But building a viral loop isn't magic—it's engineering with psychology and math.",
	}

	if intro, ok := intros[topic]; ok {
		return intro
	}

	return fmt.Sprintf("%s is a critical capability for modern engineering teams. Let's dive into how it works and how to implement it effectively.", strings.Title(strings.ReplaceAll(topic, "-", " ")))
}

// generateSolution creates the solution section
func (g *Generator) generateSolution(topic string) string {
	return fmt.Sprintf("The key insight: **automate the boring stuff, focus on the high-leverage work**.\n\nWith %s, you can:\n\n- **Eliminate manual work** through automation\n- **Scale consistently** without linear headcount growth\n- **Iterate rapidly** with immediate feedback loops\n- **Measure everything** that matters for your use case", strings.Title(strings.ReplaceAll(topic, "-", " ")))
}

// generateCodeExample creates a code example
func (g *Generator) generateCodeExample(topic string) string {
	// Language-specific examples
	if strings.Contains(topic, "go") || strings.Contains(topic, "concurrency") {
		return fmt.Sprintf("```go\n// %s implementation in Go\npackage main\n\nimport (\n    \"context\"\n    \"log\"\n)\n\nfunc main() {\n    // Your implementation here\n    log.Println(\"%s running...\")\n}\n```", g.generateTitle(topic), topic)
	}

	return fmt.Sprintf("```bash\n# %s quick start\n%s --help\n\n# See documentation for full examples\n```", g.generateTitle(topic), topic)
}

// generateNextSteps creates the next steps section
func (g *Generator) generateNextSteps(topic string) string {
	return fmt.Sprintf("%s is part of a comprehensive modern engineering stack:\n\n- **Deep Architecture Swarm** - System design review\n- **Deep Performance Swarm** - Optimization strategies\n- **Deep Security Swarm** - Automated security audits", g.generateTitle(topic))
}

// SocialPost represents a social media post
type SocialPost struct {
	Platform string
	Content  string
	Media    []string
	Thread   []string // For Twitter threads
	Metadata map[string]string
}

// GenerateSocialPosts creates social media content for a blog post
func (g *Generator) GenerateSocialPosts(blog *BlogPost) []*SocialPost {
	var posts []*SocialPost

	for _, platform := range g.cfg.Platforms {
		if !platform.Enabled {
			continue
		}

		post := &SocialPost{
			Platform: platform.Name,
			Metadata: map[string]string{
				"blog_slug": blog.Slug,
				"blog_url":  fmt.Sprintf("%s/%s", g.cfg.BlogURL, blog.Slug),
			},
		}

		switch platform.Name {
		case "linkedin":
			post.Content = g.generateLinkedInPost(blog)
		case "twitter":
			thread := g.generateTwitterThread(blog)
			post.Content = thread[0]
			post.Thread = thread
		case "youtube":
			post.Content = g.generateYouTubeCommunityPost(blog)
		case "threads":
			post.Content = g.generateThreadsPost(blog)
		}

		posts = append(posts, post)
	}

	return posts
}

// generateLinkedInPost creates a LinkedIn post
func (g *Generator) generateLinkedInPost(blog *BlogPost) string {
	return fmt.Sprintf(`🚀 Just published: %s

%s

Key takeaways:
→ %s
→ Production-ready patterns
→ Real examples from the field

Link in comments 👇

#Engineering #Tech #Development %s
`,
		blog.Title,
		blog.Description,
		strings.Join(blog.Tags[:3], ", "),
		strings.ToUpper(strings.Replace(strings.Join(blog.Tags, " #"), "-", "", -1)),
	)
}

// generateTwitterThread creates a Twitter thread
func (g *Generator) generateTwitterThread(blog *BlogPost) []string {
	return []string{
		fmt.Sprintf("🧵 %s\n\nA thread on %s and how to implement it in production:", blog.Title, blog.Tags[0]),
		fmt.Sprintf("1/%s\n\nThe problem: %s", "8", blog.Description),
		fmt.Sprintf("2/%s\n\nThe solution:\n\n%s", "8", g.generateSolution(blog.Tags[0])),
		fmt.Sprintf("3/%s\n\nCode example 👇", "8"),
		fmt.Sprintf("4/%s\n\n```go\n// Implementation snippet\n```", "8"),
		fmt.Sprintf("5/%s\n\nFull guide 👇\n\n%s", "8", fmt.Sprintf("%s/%s", g.cfg.BlogURL, blog.Slug)),
	}
}

// generateYouTubeCommunityPost creates a YouTube Community post
func (g *Generator) generateYouTubeCommunityPost(blog *BlogPost) string {
	return fmt.Sprintf(`📖 NEW BLOG POST: %s

%s

🔗 Full article: %s

%s`, blog.Title, blog.Description, fmt.Sprintf("%s/%s", g.cfg.BlogURL, blog.Slug), strings.Join(strings.Fields(blog.Tags[0]), " #"))
}

// generateThreadsPost creates a Threads post
func (g *Generator) generateThreadsPost(blog *BlogPost) string {
	return fmt.Sprintf("%s\n\n%s\n\n🔗 Link in bio", blog.Title, blog.Description)
}

// GenerateStandaloneSocial creates social posts not tied to a blog post
func (g *Generator) GenerateStandaloneSocial(tick time.Time) []*SocialPost {
	var posts []*SocialPost

	// Different content types based on time
	hour := tick.Hour()

	for _, platform := range g.cfg.Platforms {
		if !platform.Enabled {
			continue
		}

		post := &SocialPost{
			Platform: platform.Name,
			Metadata: map[string]string{
				"type": "standalone",
			},
		}

		switch {
		case hour >= 8 && hour < 12:
			// Morning: Educational/insightful
			post.Content = g.generateMorningContent(platform.Name)
		case hour >= 12 && hour < 17:
			// Afternoon: Discussion/poll
			post.Content = g.generateAfternoonContent(platform.Name)
		default:
			// Evening: Personal/behind the scenes
			post.Content = g.generateEveningContent(platform.Name)
		}

		posts = append(posts, post)
	}

	return posts
}

// GenerateEngagement creates engagement posts (replies, quotes)
func (g *Generator) GenerateEngagement(tick time.Time) []*SocialPost {
	// TODO: Implement engagement generation based on trending topics and mentions
	return nil
}

func (g *Generator) generateMorningContent(platform string) string {
	contents := map[string]string{
		"linkedin": "💡 Hot take: Most teams over-engineer their solutions before understanding the problem.\n\nStart simple, measure, then optimize.",
		"twitter":  "Most dev tools are over-engineered.\n\nSimplicity wins every time.",
		"threads":  "POV: You finally understand the problem after 3 meetings.",
		"youtube":  "Working on something exciting today! 🚀",
	}

	if content, ok := contents[platform]; ok {
		return content
	}
	return "Good morning! Building something cool today."
}

func (g *Generator) generateAfternoonContent(platform string) string {
	contents := map[string]string{
		"linkedin": `📊 POLL: What's your team's biggest challenge right now?

• Technical debt
• Hiring
• Product roadmap
• Budget constraints

Drop your choice below! 👇`,
		}

	if content, ok := contents[platform]; ok {
		return content
	}
	return "Afternoon check-in: How's your project going?"
}

func (g *Generator) generateEveningContent(platform string) string {
	contents := map[string]string{
		"linkedin": "Wrapping up the day. Today I learned that the best solution is often the one you can ship tomorrow, not the perfect one that takes a month. 🌙",
		"twitter":  "Ship tomorrow > perfect in a month",
		"threads":  "POV: That feeling when the tests finally pass",
	}

	if content, ok := contents[platform]; ok {
		return content
	}
	return "Evening reflection: Progress over perfection."
}
