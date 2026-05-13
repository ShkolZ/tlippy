# Tlippy

Tlippy is a fast tool for bulk-downloading Twitch clips by category or streamer.        
It's designed for users who want to download many clips at once instead of manually saving clips one by one.

---

# Getting Started

Download the latest release and choose the variant you want to use.

## Available Variants

### TUI
> Streamer downloads are not implemented yet (work in progress).

### CLI
> Streamer downloads are not implemented yet (work in progress).

### Basic Usage

```bash
./tlippy bulk -c <twitch_category> -t <time_range: 24h|7d|all> -a <clip_amount> -o <output_path>

# Download top 50 clips from League of Legends this week
./tlippy bulk -c "League of Legends" -t 7d -a 50 -o ./clips

# Download all-time top clips
./tlippy bulk -c Valorant -t all -a 100 -o ./valorant_clips
```

# Build Tutorial
Prerequisites:
- Go v1.26.2

Steps to build tlippy by yourself:
- Clone the repository
- Register your [Twitch Application](https://dev.twitch.tv/)
- Get CLIENT_ID and CLIENT_SECRET
- Rename .env.example into .env
- Put your CLIENT_ID and CLIENT_SECRET instead of placeholders
- Build tlippy with:
 ```bash
  go mod tidy
  go build ./cmd/tlippy-<variant>
 ```
