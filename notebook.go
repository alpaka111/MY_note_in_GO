package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

type Note struct {
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Tags      []string  `json:"tags"`
	Timestamp time.Time `json:"timestamp"`
}

type Notebook struct {
	Notes    []Note
	filename string
	password string
}

func NewNote(title, content string, tags []string) *Note {
	return &Note{
		Title:     title,
		Content:   content,
		Tags:      tags,
		Timestamp: time.Now(),
	}
}

func NewNotebook(filename, password string) *Notebook {
	return &Notebook{
		Notes:    []Note{},
		filename: filename,
		password: password,
	}
}

func (n *Notebook) AddNote(note *Note) {
	n.Notes = append(n.Notes, *note)
}

func (n *Notebook) DeleteNote(index int) {
	if index >= 0 && index < len(n.Notes) {
		n.Notes = append(n.Notes[:index], n.Notes[index+1:]...)
	}
}

func (n *Notebook) Search(query string) []Note {
	var results []Note
	query = strings.ToLower(query)

	for _, note := range n.Notes {
		if strings.Contains(strings.ToLower(note.Title), query) ||
			strings.Contains(strings.ToLower(note.Content), query) {
			results = append(results, note)
			continue
		}
		
		// Search in tags
		for _, tag := range note.Tags {
			if strings.Contains(strings.ToLower(tag), query) {
				results = append(results, note)
				break
			}
		}
	}

	return results
}

func (n *Notebook) CountWords() int {
	total := 0
	for _, note := range n.Notes {
		total += len(strings.Fields(note.Content))
	}
	return total
}

func (n *Notebook) CountTags() int {
	tagMap := make(map[string]bool)
	for _, note := range n.Notes {
		for _, tag := range note.Tags {
			tagMap[tag] = true
		}
	}
	return len(tagMap)
}

func (n *Notebook) GetTagCloud() map[string]int {
	tagCount := make(map[string]int)
	for _, note := range n.Notes {
		for _, tag := range note.Tags {
			tagCount[tag]++
		}
	}
	return tagCount
}

func (n *Notebook) GetRecentNotes(count int) []Note {
	sorted := make([]Note, len(n.Notes))
	copy(sorted, n.Notes)
	
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Timestamp.After(sorted[j].Timestamp)
	})
	
	if len(sorted) > count {
		return sorted[:count]
	}
	return sorted
}

func (n *Notebook) SortNotes(mode sortMode) {
	switch mode {
	case sortByDate:
		sort.Slice(n.Notes, func(i, j int) bool {
			return n.Notes[i].Timestamp.After(n.Notes[j].Timestamp)
		})
	case sortByTitle:
		sort.Slice(n.Notes, func(i, j int) bool {
			return n.Notes[i].Title < n.Notes[j].Title
		})
	case sortByTags:
		sort.Slice(n.Notes, func(i, j int) bool {
			if len(n.Notes[i].Tags) == 0 {
				return false
			}
			if len(n.Notes[j].Tags) == 0 {
				return true
			}
			return n.Notes[i].Tags[0] < n.Notes[j].Tags[0]
		})
	}
}

func (n *Notebook) GetSortedNotes(mode sortMode) []Note {
	sorted := make([]Note, len(n.Notes))
	copy(sorted, n.Notes)
	
	switch mode {
	case sortByDate:
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].Timestamp.After(sorted[j].Timestamp)
		})
	case sortByTitle:
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].Title < sorted[j].Title
		})
	case sortByTags:
		sort.Slice(sorted, func(i, j int) bool {
			if len(sorted[i].Tags) == 0 {
				return false
			}
			if len(sorted[j].Tags) == 0 {
				return true
			}
			return sorted[i].Tags[0] < sorted[j].Tags[0]
		})
	}
	
	return sorted
}

func (n *Notebook) Save() error {
	// Serialize notes to JSON
	data, err := json.Marshal(n.Notes)
	if err != nil {
		return err
	}

	// Encrypt
	encrypted := encrypt(data, n.password)

	// Create file header
	passwordHash := hashPassword(n.password)
	header := fmt.Sprintf("ALPAKA\nVERSION:1.0\nHASH:%s\n---ENCRYPTED---\n", passwordHash)

	// Write to file
	file, err := os.Create(n.filename)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.WriteString(header); err != nil {
		return err
	}

	if _, err := file.Write(encrypted); err != nil {
		return err
	}

	return nil
}

func LoadNotebook(filename, password string) (*Notebook, error) {
	// Read file
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Parse header
	lines := strings.Split(string(data), "\n")
	if len(lines) < 4 || lines[0] != "ALPAKA" {
		return nil, fmt.Errorf("invalid file format")
	}

	// Extract password hash
	hashLine := lines[2]
	if !strings.HasPrefix(hashLine, "HASH:") {
		return nil, fmt.Errorf("missing password hash")
	}
	storedHash := hashLine[5:]

	// Verify password
	if hashPassword(password) != storedHash {
		return nil, fmt.Errorf("invalid password")
	}

	// Find encrypted data start
	encryptedStart := strings.Index(string(data), "---ENCRYPTED---\n") + 16
	if encryptedStart == 15 {
		return nil, fmt.Errorf("missing encrypted data")
	}

	encrypted := data[encryptedStart:]

	// Decrypt
	decrypted := decrypt(encrypted, password)

	// Deserialize JSON
	var notes []Note
	if err := json.Unmarshal(decrypted, &notes); err != nil {
		return nil, fmt.Errorf("decryption error: %v", err)
	}

	return &Notebook{
		Notes:    notes,
		filename: filename,
		password: password,
	}, nil
}

// Simple XOR encryption (dla demonstracji - w produkcji użyj AES)
func encrypt(data []byte, password string) []byte {
	key := deriveKey(password, len(data)+256)
	result := make([]byte, len(data))

	for i := 0; i < len(data); i++ {
		result[i] = data[i] ^ key[i%len(key)]
	}

	return result
}

func decrypt(data []byte, password string) []byte {
	// XOR is symmetric
	return encrypt(data, password)
}

func deriveKey(password string, length int) []byte {
	key := make([]byte, length)

	for i := 0; i < length; i++ {
		key[i] = password[i%len(password)] ^ byte(i&0xFF)
	}

	return key
}

func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", hash)
}
