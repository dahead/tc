package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type TagData struct {
	Word     string   `json:"word"`
	Count    int      `json:"count"`
	Files    []string `json:"files"`
	FontSize int      `json:"fontSize"`
	Color    string   `json:"color"`
	X        float64  `json:"x"`
	Y        float64  `json:"y"`
}

type TemplateData struct {
	Tags      []*TagData `json:"tags"`
	TotalTags int        `json:"totalTags"`
	Layout    string     `json:"layout"`
	MaxCount  int        `json:"maxCount"`
	MinCount  int        `json:"minCount"`
}

var colors = []string{
	"#FF6B6B", "#4ECDC4", "#45B7D1", "#96CEB4",
	"#FFEAA7", "#DDA0DD", "#98D8C8", "#F7DC6F",
	"#FF8A65", "#81C784", "#64B5F6", "#FFB74D",
	"#F06292", "#90A4AE", "#A1887F", "#FFF176",
}

var wordRegex = regexp.MustCompile(`[a-zA-Z]{2,}`)

func main() {
	// Parse command line flags
	sortName := flag.Bool("name", false, "Sort tags by name (default: by count descending)")
	amount := flag.Int("amount", 100, "Maximum number of tags to display")
	output := flag.String("output", "tagcloud.html", "Output HTML file path")
	templateDir := flag.String("template", "./templates", "Template directory path")
	flag.Parse()

	// Check if stdin has data
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		fmt.Println("No input received from STDIN.")
		fmt.Println("Usage: cat list.txt | ./app [-name] [-amount 100] [-output tagcloud.html] [-template ./templates]")
		fmt.Println("   or: echo 'file1.mp4\\nfile2.txt' | ./app")
		return
	}

	// Read from STDIN
	scanner := bufio.NewScanner(os.Stdin)
	var filenames []string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			filenames = append(filenames, line)
		}
	}

	if len(filenames) == 0 {
		fmt.Println("No input received. Usage: cat list.txt | ./app [-name] [-amount 100] [-output tagcloud.html] [-template ./templates]")
		fmt.Println(("Create a file list: find \"$(pwd)\" -name \"*.mp4\" -type f ! -path '*/.*' > ~/list.txt"))
		return
	}

	// Process filenames to create tag data
	tagMap := make(map[string]*TagData)

	for _, filename := range filenames {
		// Extract just the filename without path
		baseName := filepath.Base(filename)
		words := extractWords(baseName)
		for _, word := range words {
			if _, exists := tagMap[word]; !exists {
				tagMap[word] = &TagData{
					Word:  word,
					Count: 0,
					Files: []string{},
				}
			}
			tagMap[word].Count++
			tagMap[word].Files = append(tagMap[word].Files, filename)
		}
	}

	// Convert to slice
	var tags []*TagData
	for _, tag := range tagMap {
		tags = append(tags, tag)
	}

	if len(tags) == 0 {
		fmt.Println("No valid words found in filenames")
		return
	}

	// Sort based on flags
	if *sortName {
		sort.Slice(tags, func(i, j int) bool {
			return tags[i].Word < tags[j].Word
		})
	} else {
		// Default: sort by count descending
		sort.Slice(tags, func(i, j int) bool {
			return tags[i].Count > tags[j].Count
		})
	}

	// Limit to specified amount
	if len(tags) > *amount {
		tags = tags[:*amount]
	}

	// Calculate font sizes and colors
	maxCount := tags[0].Count // Since we sort by count desc by default
	minCount := tags[len(tags)-1].Count

	for i, tag := range tags {
		// Font size between 12 and 48px
		if maxCount == minCount {
			tag.FontSize = 24
		} else {
			tag.FontSize = 12 + int(float64(tag.Count-minCount)/float64(maxCount-minCount)*36)
		}
		tag.Color = colors[i%len(colors)]
	}

	// Position tags in spiral layout
	positionTagsSpiral(tags)

	// Prepare template data
	templateData := &TemplateData{
		Tags:      tags,
		TotalTags: len(tags),
		Layout:    "Spiral",
		MaxCount:  maxCount,
		MinCount:  minCount,
	}

	// Generate HTML
	err := generateHTML(templateData, *templateDir, *output)
	if err != nil {
		fmt.Printf("Error generating HTML: %v\n", err)
		return
	}

	fmt.Printf("Tag cloud generated: %s\n", *output)
	fmt.Printf("Total tags: %d\n", len(tags))

	// Debug: Print top 5 tags
	fmt.Println("Top tags:")
	for i, tag := range tags {
		if i >= 5 {
			break
		}
		fmt.Printf("  %s (%d files)\n", tag.Word, tag.Count)
	}
}

func extractWords(filename string) []string {
	matches := wordRegex.FindAllString(filename, -1)

	var result []string
	for _, word := range matches {
		word = strings.ToLower(word)
		if len(word) > 2 { // Filter out very short words
			result = append(result, word)
		}
	}

	return result
}

func positionTagsSpiral(tags []*TagData) {
	centerX, centerY := 50.0, 50.0 // Center of the container (in percentage)
	radius := 5.0
	angle := 0.0
	angleStep := 0.5
	radiusStep := 1.5

	for i, tag := range tags {
		if i == 0 {
			// Place the first (most frequent) tag in the center
			tag.X = centerX - 5 // Offset a bit for better centering
			tag.Y = centerY - 2
		} else {
			// Calculate spiral position
			x := centerX + radius*math.Cos(angle)
			y := centerY + radius*math.Sin(angle)

			// Ensure tags stay within bounds
			if x < 5 {
				x = 5
			} else if x > 85 {
				x = 85
			}
			if y < 5 {
				y = 5
			} else if y > 85 {
				y = 85
			}

			tag.X = x
			tag.Y = y

			// Update spiral parameters with more spacing
			angle += angleStep
			radius += radiusStep / 5.0 // Increase spacing (was /10.0)
		}
	}
}

func generateHTML(data *TemplateData, templateDir, outputFile string) error {
	// Read HTML template
	htmlPath := filepath.Join(templateDir, "index.html")
	htmlContent, err := os.ReadFile(htmlPath)
	if err != nil {
		return fmt.Errorf("failed to read HTML template: %v", err)
	}

	// Create template with custom functions
	tmpl := template.New("tagcloud").Funcs(template.FuncMap{
		"toJSON": func(v interface{}) string {
			jsonBytes, _ := json.Marshal(v)
			return string(jsonBytes)
		},
		"templateDir": func() string {
			return templateDir
		},
	})

	tmpl, err = tmpl.Parse(string(htmlContent))
	if err != nil {
		return fmt.Errorf("failed to parse HTML template: %v", err)
	}

	// Create output file
	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer file.Close()

	// Execute template
	err = tmpl.Execute(file, data)
	if err != nil {
		return fmt.Errorf("failed to execute template: %v", err)
	}

	return nil
}
