# Filename Tag Cloud Generator

A Go application that generates interactive HTML tag clouds from filenames, with support for video playback.

## Features

- Extract words from filenames to create visual tag clouds
- Interactive spiral layout with hover effects
- Click tags to see all files containing that word
- Video files open in new browser tabs
- Non-video files copy to clipboard
- Responsive design with animations

## Build

```bash
go build -o tc main.go
```

## Requirements

- Go 1.16+
- Modern web browser
- Template files in ./templates/ directory

## Usage

```bash
# Basic usage
find /path/to/files -name "*.mp4" | ./app

# Sort by name instead of frequency
find /path/to/files -type f | ./app -name

# Limit to top 50 tags
ls *.* | ./app -amount 50

# Custom output file
cat filelist.txt | ./app -output my-tagcloud.html
```

## Options

-name: Sort tags alphabetically (default: by frequency)
-amount N: Show top N tags (default: 100)
-output FILE: Output HTML file (default: tagcloud.html)
-template DIR: Template directory (default: ./templates)

## File Structure
