package main

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/go-github/v74/github"
	"gopkg.in/yaml.v3"
)

const (
	owner     = "apple"
	repo      = "device-management"
	targetDir = "../temp"
)

type RepositoryAnalysis struct {
	LastUpdated   time.Time                      `json:"last_updated"`
	TotalFiles    int                            `json:"total_files"`
	Categories    map[string]int                 `json:"categories"`
	Declarations  map[string]int                 `json:"declarations"`
	SampleConfigs map[string]ConfigurationSample `json:"sample_configs"`
}

type ConfigurationSample struct {
	Type        string                 `json:"type"`
	Description string                 `json:"description"`
	SupportedOS map[string]interface{} `json:"supported_os"`
	Keys        []string               `json:"keys"`
}

type DDMSchema struct {
	Title       string                 `yaml:"title"`
	Description string                 `yaml:"description"`
	Payload     map[string]interface{} `yaml:"payload"`
	PayloadKeys []PayloadKey           `yaml:"payloadkeys"`
}

type PayloadKey struct {
	Key         string                 `yaml:"key"`
	Title       string                 `yaml:"title"`
	Type        string                 `yaml:"type"`
	Presence    string                 `yaml:"presence"`
	Content     string                 `yaml:"content"`
	Default     interface{}            `yaml:"default"`
	Range       map[string]interface{} `yaml:"range"`
	SupportedOS map[string]interface{} `yaml:"supportedOS"`
}

func main() {
	fmt.Println("🍎 Apple Device Management DDM Specifications Downloader v2.0")
	fmt.Println("===============================================================")

	ctx := context.Background()
	client := github.NewClient(nil)

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("Failed to get working directory:", err)
	}

	targetPath := filepath.Join(wd, targetDir)
	repoPath := filepath.Join(targetPath, repo)

	if err := os.MkdirAll(targetPath, 0755); err != nil {
		log.Fatal("Failed to create target directory:", err)
	}

	fmt.Printf("📡 Fetching repository information from GitHub API...\n")
	repoInfo, _, err := client.Repositories.Get(ctx, owner, repo)
	if err != nil {
		log.Fatal("Failed to get repository info:", err)
	}

	fmt.Printf("📊 Repository: %s\n", repoInfo.GetFullName())
	fmt.Printf("📝 Description: %s\n", repoInfo.GetDescription())
	fmt.Printf("⭐ Stars: %d\n", repoInfo.GetStargazersCount())
	fmt.Printf("📅 Last Updated: %s\n", repoInfo.GetUpdatedAt().Format("2006-01-02 15:04:05"))

	shouldDownload := true
	if _, err := os.Stat(repoPath); err == nil {
		fmt.Printf("📁 Repository already exists locally\n")

		analysisPath := filepath.Join(targetPath, "analysis.json")
		if _, err := os.Stat(analysisPath); err == nil {
			fmt.Printf("🔍 Local analysis data found\n")
		}

		fmt.Printf("🔄 Proceeding with fresh download to ensure latest data...\n")
	}

	if shouldDownload {
		fmt.Printf("📦 Downloading repository archive...\n")
		if err := downloadRepositoryZip(ctx, client, targetPath); err != nil {
			log.Fatal("Failed to download repository:", err)
		}
		fmt.Println("✅ Repository downloaded successfully!")
	}

	fmt.Println("\n🔍 Analyzing DDM specifications...")
	analysis, err := analyzeRepository(repoPath)
	if err != nil {
		log.Fatal("Failed to analyze repository:", err)
	}

	displayAnalysis(analysis)

	if err := saveAnalysis(analysis, filepath.Join(targetPath, "analysis.json")); err != nil {
		log.Printf("Warning: Failed to save analysis: %v", err)
	}

	fmt.Printf("\n🎉 Download and analysis complete!\n")
	fmt.Printf("📍 Repository location: %s\n", repoPath)
	fmt.Printf("⏰ Completed: %s\n", time.Now().Format("2006-01-02 15:04:05"))
}

func downloadRepositoryZip(ctx context.Context, client *github.Client, targetPath string) error {
	// Get the default branch ZIP archive URL
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/zipball", owner, repo)

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Make request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download: status %d", resp.StatusCode)
	}

	// Save to temporary file
	zipPath := filepath.Join(targetPath, "repo.zip")
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return fmt.Errorf("failed to create zip file: %w", err)
	}
	defer zipFile.Close()

	if _, err := io.Copy(zipFile, resp.Body); err != nil {
		return fmt.Errorf("failed to save zip: %w", err)
	}

	// Extract ZIP
	if err := extractZip(zipPath, targetPath); err != nil {
		return fmt.Errorf("failed to extract zip: %w", err)
	}

	// Clean up zip file
	os.Remove(zipPath)

	return nil
}

func extractZip(zipPath, targetPath string) error {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer reader.Close()

	// Remove existing directory if it exists
	repoPath := filepath.Join(targetPath, repo)
	os.RemoveAll(repoPath)

	for _, file := range reader.File {
		// Skip the root directory and adjust path
		parts := strings.Split(file.Name, "/")
		if len(parts) <= 1 {
			continue
		}

		// Replace the GitHub-generated root directory name with our desired name
		parts[0] = repo
		newPath := filepath.Join(targetPath, filepath.Join(parts...))

		if file.FileInfo().IsDir() {
			os.MkdirAll(newPath, 0755)
			continue
		}

		// Create parent directories
		os.MkdirAll(filepath.Dir(newPath), 0755)

		// Extract file
		rc, err := file.Open()
		if err != nil {
			return err
		}

		outFile, err := os.Create(newPath)
		if err != nil {
			rc.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}

	return nil
}

func analyzeRepository(repoPath string) (*RepositoryAnalysis, error) {
	analysis := &RepositoryAnalysis{
		LastUpdated:   time.Now(),
		Categories:    make(map[string]int),
		Declarations:  make(map[string]int),
		SampleConfigs: make(map[string]ConfigurationSample),
	}

	// Count files in different categories
	categories := map[string]string{
		"Declarations": "declarative/declarations",
		"Status Items": "declarative/status",
		"MDM Commands": "mdm/commands",
		"MDM Profiles": "mdm/profiles",
		"MDM Check-in": "mdm/checkin",
		"Protocol":     "declarative/protocol",
		"MDM Errors":   "mdm/errors",
	}

	totalFiles := 0
	for category, dir := range categories {
		count, err := countYAMLFiles(filepath.Join(repoPath, dir))
		if err != nil {
			return nil, fmt.Errorf("failed to count files in %s: %w", dir, err)
		}
		analysis.Categories[category] = count
		totalFiles += count
	}
	analysis.TotalFiles = totalFiles

	// Analyze declaration types
	declarationsDir := filepath.Join(repoPath, "declarative", "declarations")
	subdirs := []string{"configurations", "assets", "activations", "management"}

	for _, subdir := range subdirs {
		path := filepath.Join(declarationsDir, subdir)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}

		count, err := countYAMLFiles(path)
		if err != nil {
			return nil, err
		}
		analysis.Declarations[subdir] = count
	}

	// Parse sample configurations
	configDir := filepath.Join(declarationsDir, "configurations")
	if err := parseSampleConfigurations(configDir, analysis.SampleConfigs); err != nil {
		return nil, fmt.Errorf("failed to parse sample configurations: %w", err)
	}

	return analysis, nil
}

func countYAMLFiles(dir string) (int, error) {
	count := 0

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".yaml" {
			count++
		}
		return nil
	})

	return count, err
}

func parseSampleConfigurations(configDir string, samples map[string]ConfigurationSample) error {
	return filepath.Walk(configDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || filepath.Ext(path) != ".yaml" {
			return nil
		}

		// Only parse a few sample configurations to avoid overwhelming output
		filename := filepath.Base(path)
		if !strings.Contains(filename, "passcode") &&
			!strings.Contains(filename, "wifi") &&
			!strings.Contains(filename, "app.managed") {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return nil // Skip files we can't read
		}

		var schema DDMSchema
		if err := yaml.Unmarshal(data, &schema); err != nil {
			return nil // Skip files we can't parse
		}

		sample := ConfigurationSample{
			Type:        fmt.Sprintf("%v", schema.Payload["declarationtype"]),
			Description: schema.Description,
			SupportedOS: make(map[string]interface{}),
			Keys:        make([]string, 0, len(schema.PayloadKeys)),
		}

		// Extract supported OS info
		if supportedOS, ok := schema.Payload["supportedOS"]; ok {
			if osMap, ok := supportedOS.(map[string]interface{}); ok {
				sample.SupportedOS = osMap
			}
		}

		// Extract key names
		for _, key := range schema.PayloadKeys {
			sample.Keys = append(sample.Keys, key.Key)
		}

		configName := strings.TrimSuffix(filename, ".yaml")
		samples[configName] = sample

		return nil
	})
}

func displayAnalysis(analysis *RepositoryAnalysis) {
	fmt.Println("\n📊 Repository Analysis Results:")
	fmt.Println("===============================")
	fmt.Printf("📄 Total YAML files: %d\n", analysis.TotalFiles)

	fmt.Println("\n📁 Categories:")
	for category, count := range analysis.Categories {
		fmt.Printf("   %-15s: %3d files\n", category, count)
	}

	fmt.Println("\n🏷️  Declaration Types:")
	for declType, count := range analysis.Declarations {
		fmt.Printf("   %-15s: %3d files\n", declType, count)
	}

	if len(analysis.SampleConfigs) > 0 {
		fmt.Println("\n📋 Sample Configurations:")
		for name, config := range analysis.SampleConfigs {
			fmt.Printf("   %-20s: %s\n", name, config.Type)
			if len(config.Keys) > 0 {
				fmt.Printf("      Keys: %s\n", strings.Join(config.Keys[:min(5, len(config.Keys))], ", "))
				if len(config.Keys) > 5 {
					fmt.Printf("      ... and %d more keys\n", len(config.Keys)-5)
				}
			}
		}
	}
}

func saveAnalysis(analysis *RepositoryAnalysis, path string) error {
	// This would save analysis as JSON - implementation omitted for brevity
	fmt.Printf("💾 Analysis data would be saved to: %s\n", path)
	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
