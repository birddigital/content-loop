package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

// Config holds all configuration for ContentLoop
type Config struct {
	Env          string            `yaml:"env" envconfig:"ENV" default:"production"`
	Author       string            `yaml:"author" envconfig:"AUTHOR" default:"Adrien Bird"`
	BlogURL      string            `yaml:"blog_url" envconfig:"BLOG_URL" default:"https://adrienbird.net/blog"`
	BlogPath     string            `yaml:"blog_path" envconfig:"BLOG_PATH" default:"~/sources/adrienbird/src/app/blog/posts/"`
	Social       SocialConfig      `yaml:"social"`
	OpenAI       OpenAIConfig      `yaml:"openai"`
	Anthropic    AnthropicConfig   `yaml:"anthropic"`
	Scheduling   SchedulingConfig  `yaml:"scheduling"`
	Analytics    AnalyticsConfig   `yaml:"analytics"`
}

// SocialConfig holds social media platform configurations
type SocialConfig struct {
	LinkedIn   LinkedInConfig   `yaml:"linkedin"`
	Twitter    TwitterConfig    `yaml:"twitter"`
	YouTube    YouTubeConfig    `yaml:"youtube"`
	Threads    ThreadsConfig    `yaml:"threads"`
}

// LinkedInConfig holds LinkedIn-specific settings
type LinkedInConfig struct {
	Enabled      bool   `yaml:"enabled" envconfig:"LINKEDIN_ENABLED" default:"true"`
	AccessToken string `yaml:"access_token" envconfig:"LINKEDIN_ACCESS_TOKEN"`
	OrganisationID string `yaml:"organisation_id" envconfig:"LINKEDIN_ORGANISATION_ID"`
}

// TwitterConfig holds Twitter/X-specific settings
type TwitterConfig struct {
	Enabled      bool   `yaml:"enabled" envconfig:"TWITTER_ENABLED" default:"true"`
	AccessToken  string `yaml:"access_token" envconfig:"TWITTER_ACCESS_TOKEN"`
	AccessSecret string `yaml:"access_secret" envconfig:"TWITTER_ACCESS_SECRET"`
	APIKey       string `yaml:"api_key" envconfig:"TWITTER_API_KEY"`
	APISecret    string `yaml:"api_secret" envconfig:"TWITTER_API_SECRET"`
	BearerToken  string `yaml:"bearer_token" envconfig:"TWITTER_BEARER_TOKEN"`
}

// YouTubeConfig holds YouTube-specific settings
type YouTubeConfig struct {
	Enabled         bool   `yaml:"enabled" envconfig:"YOUTUBE_ENABLED" default:"true"`
	AccessToken     string `yaml:"access_token" envconfig:"YOUTUBE_ACCESS_TOKEN"`
	RefreshToken    string `yaml:"refresh_token" envconfig:"YOUTUBE_REFRESH_TOKEN"`
	ClientID        string `yaml:"client_id" envconfig:"YOUTUBE_CLIENT_ID"`
	ClientSecret    string `yaml:"client_secret" envconfig:"YOUTUBE_CLIENT_SECRET"`
	ChannelID       string `yaml:"channel_id" envconfig:"YOUTUBE_CHANNEL_ID"`
	DefaultPlaylist string `yaml:"default_playlist" envconfig:"YOUTUBE_DEFAULT_PLAYLIST"`
}

// ThreadsConfig holds Threads-specific settings
type ThreadsConfig struct {
	Enabled     bool   `yaml:"enabled" envconfig:"THREADS_ENABLED" default:"true"`
	AccessToken string `yaml:"access_token" envconfig:"THREADS_ACCESS_TOKEN"`
	UserID      string `yaml:"user_id" envconfig:"THREADS_USER_ID"`
}

// OpenAIConfig holds OpenAI API settings
type OpenAIConfig struct {
	Enabled   bool   `yaml:"enabled" envconfig:"OPENAI_ENABLED" default:"true"`
	APIKey    string `yaml:"api_key" envconfig:"OPENAI_API_KEY"`
	Model     string `yaml:"model" envconfig:"OPENAI_MODEL" default:"gpt-4o"`
	MaxTokens int    `yaml:"max_tokens" envconfig:"OPENAI_MAX_TOKENS" default:"4000"`
}

// AnthropicConfig holds Anthropic API settings
type AnthropicConfig struct {
	Enabled   bool   `yaml:"enabled" envconfig:"ANTHROPIC_ENABLED" default:"true"`
	APIKey    string `yaml:"api_key" envconfig:"ANTHROPIC_API_KEY"`
	Model     string `yaml:"model" envconfig:"ANTHROPIC_MODEL" default:"claude-sonnet-4-20250514"`
	MaxTokens int    `yaml:"max_tokens" envconfig:"ANTHROPIC_MAX_TOKENS" default:"8000"`
}

// SchedulingConfig holds scheduling settings
type SchedulingConfig struct {
	Timezone     string    `yaml:"timezone" envconfig:"TZ" default:"America/New_York"`
	BlogCadence  int       `yaml:"blog_cadence" envconfig:"BLOG_CADENCE" default:"3"` // Blog post every N hours
	SocialCadence int      `yaml:"social_cadence" envconfig:"SOCIAL_CADENCE" default:"1"` // Social post every N hours
	QuietHours   QuietHours `yaml:"quiet_hours"`
}

// QuietHours defines when not to post
type QuietHours struct {
	Enabled bool   `yaml:"enabled" envconfig:"QUIET_HOURS_ENABLED" default:"false"`
	Start   string `yaml:"start" envconfig:"QUIET_HOURS_START" default:"22:00"`
	End     string `yaml:"end" envconfig:"QUIET_HOURS_END" default:"08:00"`
	Timezone string `yaml:"timezone" envconfig:"QUIET_HOURS_TIMEZONE" default:"America/New_York"`
}

// AnalyticsConfig holds analytics integration settings
type AnalyticsConfig struct {
	GoogleAnalytics GoogleAnalyticsConfig `yaml:"google_analytics"`
	Plausible       PlausibleConfig       `yaml:"plausible"`
	CustomAnalytics CustomAnalyticsConfig `yaml:"custom"`
}

// GoogleAnalyticsConfig holds GA4 settings
type GoogleAnalyticsConfig struct {
	Enabled    bool   `yaml:"enabled" envconfig:"GA_ENABLED" default:"true"`
	MeasurementID string `yaml:"measurement_id" envconfig:"GA_MEASUREMENT_ID"`
	APISecret  string `yaml:"api_secret" envconfig:"GA_API_SECRET"`
}

// PlausibleConfig holds Plausible analytics settings
type PlausibleConfig struct {
	Enabled bool   `yaml:"enabled" envconfig:"PLAUSIBLE_ENABLED" default:"false"`
	Domain  string `yaml:"domain" envconfig:"PLAUSIBLE_DOMAIN"`
	APIKey  string `yaml:"api_key" envconfig:"PLAUSIBLE_API_KEY"`
}

// CustomAnalyticsConfig holds custom analytics endpoint settings
type CustomAnalyticsConfig struct {
	Enabled  bool   `yaml:"enabled" envconfig:"CUSTOM_ANALYTICS_ENABLED" default:"false"`
	Endpoint string `yaml:"endpoint" envconfig:"CUSTOM_ANALYTICS_ENDPOINT"`
	APIKey   string `yaml:"api_key" envconfig:"CUSTOM_ANALYTICS_API_KEY"`
}

// Category represents a content category with its topics and cadence
type Category struct {
	Name     string   `yaml:"name"`
	Topics   []string `yaml:"topics"`
	Tags     []string `yaml:"tags"`
	Cadence  int      `yaml:"cadence"`  // Hours between posts in this category
	Priority int      `yaml:"priority"` // Lower = higher priority
}

// Platform represents a social media platform
type Platform struct {
	Name        string      `yaml:"name"`
	Enabled     bool        `yaml:"enabled"`
	Formats     []string    `yaml:"formats"`
	Limits      PostLimits  `yaml:"limits"`
	OptimalTimes []string   `yaml:"optimal_times"`
}

// PostLimits defines posting limits for a platform
type PostLimits struct {
	MaxChars    int `yaml:"max_chars"`
	MaxHashtags int `yaml:"max_hashtags"`
	MaxEmojis   int `yaml:"max_emojis"`
}

// Load loads configuration from file and environment variables
func Load() (*Config, error) {
	cfg := &Config{}

	// Try to load from YAML file first
	configPaths := []string{
		"config.yaml",
		"~/.content-loop/config.yaml",
		"/etc/content-loop/config.yaml",
	}

	for _, path := range configPaths {
		expandedPath := expandPath(path)
		if data, err := os.ReadFile(expandedPath); err == nil {
			if err := yaml.Unmarshal(data, cfg); err != nil {
				return nil, fmt.Errorf("parse config file %s: %w", path, err)
			}
			break
		}
	}

	// Override with environment variables
	if err := envconfig.Process("", cfg); err != nil {
		return nil, fmt.Errorf("process env config: %w", err)
	}

	return cfg, nil
}

func expandPath(path string) string {
	if len(path) > 0 && path[0] == '~' {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[1:])
	}
	return path
}
