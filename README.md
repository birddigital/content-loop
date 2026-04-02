# ContentLoop

24/7 automated content generation and social media distribution for adrienbird.net

## What It Does

ContentLoop runs continuously, generating and publishing:

| Frequency | Content Type | Platforms |
|-----------|--------------|-----------|
| Every 3 hours | Blog post | adrienbird.net |
| Every hour | Social posts | LinkedIn, Twitter/X, YouTube, Threads |
| Continuous | Engagement | Replies, quotes, shares |

## Content Categories

- **Agent Swarms** - AI automation frameworks
- **Golang Development** - Performance, concurrency, microservices
- **Cloud & DevOps** - Cost optimization, Kubernetes, security
- **Security & Compliance** - Audits, GDPR, HIPAA, patent search
- **Digital Marketing** - SEO automation, analytics
- **Growth Hacking** - Viral loops, SaaS growth, PLG
- **YouTube Strategy** - Algorithm, content planning
- **Product Development** - MVPs, validation, prototyping
- **Case Studies** - Real project examples

## Installation

```bash
cd ~/sources/standalone-projects/content-loop
chmod +x install.sh
./install.sh
```

## Configuration

Edit `config.yaml` with your API credentials:

```yaml
social:
  linkedin:
    access_token: "YOUR_TOKEN"
  twitter:
    bearer_token: "YOUR_TOKEN"
  youtube:
    access_token: "YOUR_TOKEN"
```

## Usage

### Start manually
```bash
./content-loop
```

### Check status
```bash
# macOS
launchctl list | grep content-loop

# Linux
systemctl status content-loop
```

### View logs
```bash
# macOS
tail -f ~/.local/var/log/content-loop.out.log

# Linux
journalctl -u content-loop -f
```

### Stop
```bash
# macOS
launchctl unload ~/Library/LaunchAgents/com.birddigital.content-loop.plist

# Linux
systemctl stop content-loop
```

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      ContentLoop                            │
├─────────────────────────────────────────────────────────────┤
│  ┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────┐ │
│  │Scheduler │───▶│Generator │───▶│ Platform │───▶│Social│ │
│  │(hourly)  │    │          │    │ Manager  │    │ APIs │ │
│  └──────────┘    └──────────┘    └──────────┘    └──────┘ │
│       │               │                  │                 │
│       ▼               ▼                  ▼                 │
│  ┌──────────┐    ┌──────────┐    ┌──────────┐             │
│  │  Config  │    │ Blog Post│    │ Analytics│             │
│  └──────────┘    │ Deploy   │    │ Feedback │             │
│                  └──────────┘    └──────────┘             │
└─────────────────────────────────────────────────────────────┘
```

## Content Flow

1. **Hourly Tick** - Scheduler triggers at :00
2. **Plan Content** - Select category and topic based on rotation
3. **Generate Blog** - Create markdown post with frontmatter
4. **Deploy Blog** - SCP to bird-m1, add to Git
5. **Generate Social** - Create platform-specific posts
6. **Distribute** - Publish to LinkedIn, Twitter, YouTube, Threads
7. **Track** - Log engagement for analytics feedback

## Development

```bash
# Run once for testing
go run . --once

# Build
go build -o content-loop .

# Test content generation (dry run)
./content-loop --dry-run
```

## License

MIT
