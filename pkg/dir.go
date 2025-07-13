package pkg

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"strconv"
	"strings"
	"time"
)

type EntryType string

const (
	File      EntryType = "file"
	Directory EntryType = "directory"
)

type DirectoryEntry struct {
	Name         string           `json:"name"`
	LastModified time.Time        `json:"last_modified"`
	Size         *int64           `json:"size,omitempty"`    // nil for directories
	Type         EntryType        `json:"type"`              // "file" or "directory"
	Entries      []DirectoryEntry `json:"entries,omitempty"` // subdirectories
}

func parseDirectoryListing(html string) ([]DirectoryEntry, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	var entries []DirectoryEntry

	doc.Find("pre").Each(func(i int, s *goquery.Selection) {
		if entries == nil {
			entries = make([]DirectoryEntry, 0)
		}
		lines := strings.Split(s.Text(), "\n")
		for _, line := range lines {
			// Skip blank lines or parent link
			if strings.TrimSpace(line) == "" || strings.Contains(line, "../") {
				continue
			}

			// Split line by fields from right to left
			fields := strings.Fields(line)
			if len(fields) < 3 {
				continue
			}

			// Extract date + time
			dateStr := fields[len(fields)-3]
			timeStr := fields[len(fields)-2]
			timestampStr := dateStr + " " + timeStr
			modified, err := time.Parse("02-Jan-2006 15:04", timestampStr)
			if err != nil {
				log.Printf("could not parse time %q: %v", timestampStr, err)
				continue
			}

			// Extract size
			sizeStr := fields[len(fields)-1]
			var size *int64
			var entryType EntryType

			if sizeStr == "-" {
				entryType = Directory
			} else {
				entryType = File
				if s, err := strconv.ParseInt(sizeStr, 10, 64); err == nil {
					size = &s
				}
			}

			// Extract name by trimming timestamp/size fields from line
			namePart := strings.TrimSpace(line[:strings.Index(line, dateStr)])
			name := strings.TrimSpace(namePart)

			entry := DirectoryEntry{
				Name:         strings.TrimSuffix(name, "/"),
				LastModified: modified,
				Size:         size,
				Type:         entryType,
			}

			entries = append(entries, entry)
		}
	})
	return entries, nil
}
