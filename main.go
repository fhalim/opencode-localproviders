package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	baseURL := flag.String("base-url", "", "OpenAI-compatible base URL (e.g., http://localhost:1234/v1 or http://localhost:1234/)")
	baseURLShort := flag.String("u", "", "")
	provider := flag.String("provider", "", "Provider key name (e.g., lmstudio)")
	providerShort := flag.String("p", "", "")
	name := flag.String("name", "", "Display name (defaults to titlecase provider + ' (local)')")
	nameShort := flag.String("n", "", "")
	configPath := flag.String("config", "", "Path to opencode.json (defaults to ~/.config/opencode/opencode.json)")
	configShort := flag.String("c", "", "")
	dryRun := flag.Bool("dry-run", false, "Print generated block without modifying config")

	flag.Parse()

	// Handle short flags
	if *baseURLShort != "" {
		baseURL = baseURLShort
	}
	if *providerShort != "" {
		provider = providerShort
	}
	if *nameShort != "" {
		name = nameShort
	}
	if *configShort != "" {
		configPath = configShort
	}

	// Validate required flags
	if *baseURL == "" || *provider == "" {
		fmt.Fprintf(os.Stderr, "Error: --base-url and --provider are required\n")
		os.Exit(1)
	}

	// Normalize base URL
	normalized := normalizeURL(*baseURL)

	// Fetch models
	models, err := fetchModels(normalized)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching models: %v\n", err)
		os.Exit(1)
	}

	// Generate provider block
	displayName := *name
	if displayName == "" {
		displayName = titleCase(*provider) + " (local)"
	}

	providerBlock := generateProviderBlock(displayName, normalized, models)

	// Dry run: just print
	if *dryRun {
		jsonBytes, _ := json.MarshalIndent(providerBlock, "", "  ")
		fmt.Printf("%s\n", jsonBytes)
		return
	}

	// Resolve config path
	if *configPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
			os.Exit(1)
		}
		*configPath = filepath.Join(home, ".config/opencode/opencode.json")
	}

	// Update config
	if err := updateConfig(*configPath, *provider, providerBlock); err != nil {
		fmt.Fprintf(os.Stderr, "Error updating config: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Updated %s with %d models for provider '%s'\n", *configPath, len(models), *provider)
}

func normalizeURL(baseURL string) string {
	// Remove trailing slashes
	baseURL = strings.TrimRight(baseURL, "/")

	// If doesn't end with /v1, append it
	if !strings.HasSuffix(baseURL, "/v1") {
		baseURL += "/v1"
	}

	return baseURL
}

func fetchModels(baseURL string) ([]string, error) {
	resp, err := http.Get(baseURL + "/models")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var apiResp struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	models := make([]string, len(apiResp.Data))
	for i, m := range apiResp.Data {
		models[i] = m.ID
	}

	return models, nil
}

func generateProviderBlock(displayName, baseURL string, models []string) map[string]interface{} {
	modelMap := make(map[string]map[string]string)
	for _, m := range models {
		modelMap[m] = map[string]string{"name": m}
	}

	return map[string]interface{}{
		"npm":  "@ai-sdk/openai-compatible",
		"name": displayName,
		"options": map[string]string{
			"baseURL": baseURL,
		},
		"models": modelMap,
	}
}

func updateConfig(configPath string, providerName string, providerBlock map[string]interface{}) error {
	// Read config
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	var config map[string]json.RawMessage
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	// Ensure provider key exists
	if config["provider"] == nil {
		config["provider"] = json.RawMessage([]byte("{}"))
	}

	var providers map[string]json.RawMessage
	if err := json.Unmarshal(config["provider"], &providers); err != nil {
		return fmt.Errorf("failed to parse provider block: %w", err)
	}

	// Marshal new provider entry
	blockJSON, err := json.Marshal(providerBlock)
	if err != nil {
		return fmt.Errorf("failed to marshal provider block: %w", err)
	}

	providers[providerName] = blockJSON

	// Remarshal providers
	providersJSON, err := json.Marshal(providers)
	if err != nil {
		return fmt.Errorf("failed to remarshal providers: %w", err)
	}

	config["provider"] = providersJSON

	// Write back with indentation
	outData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, outData, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

func titleCase(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(string(s[0])) + strings.ToLower(s[1:])
}
