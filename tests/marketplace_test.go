package tests_test

import (
	"encoding/json"
	"path/filepath"
	"testing"
)

type marketplaceCatalog struct {
	Name      string               `json:"name"`
	Interface marketplaceInterface `json:"interface"`
	Plugins   []marketplacePlugin  `json:"plugins"`
}

type marketplaceInterface struct {
	DisplayName string `json:"displayName"`
}

type marketplacePlugin struct {
	Name     string            `json:"name"`
	Source   marketplaceSource `json:"source"`
	Policy   marketplacePolicy `json:"policy"`
	Category string            `json:"category"`
}

type marketplaceSource struct {
	Source string `json:"source"`
	Path   string `json:"path"`
}

type marketplacePolicy struct {
	Installation   string `json:"installation"`
	Authentication string `json:"authentication"`
}

func Test_MarketplaceCatalog_exposes_packaged_local_plugin_without_auth_or_network_promises(t *testing.T) {
	// Given
	catalogPath := filepath.Join("..", ".agents", "plugins", "marketplace.json")

	// When
	data := readFile(t, catalogPath)
	var catalog marketplaceCatalog
	decodeJSON(t, catalogPath, data, &catalog)
	var fields map[string]json.RawMessage
	decodeJSON(t, catalogPath, data, &fields)
	var rawCatalog struct {
		Plugins []map[string]json.RawMessage `json:"plugins"`
	}
	decodeJSON(t, catalogPath, data, &rawCatalog)

	// Then
	if catalog.Name != "codex-next-prompt" || catalog.Interface.DisplayName == "" {
		t.Errorf("Then catalog name = %q and displayName = %q", catalog.Name, catalog.Interface.DisplayName)
	}
	if len(catalog.Plugins) != 1 {
		t.Fatalf("Then plugins = %d, want 1", len(catalog.Plugins))
	}
	plugin := catalog.Plugins[0]
	if plugin.Name != "codex-next-prompt" || plugin.Source.Source != "local" || plugin.Source.Path != "./" {
		t.Errorf("Then plugin name = %q and packaged source = %#v", plugin.Name, plugin.Source)
	}
	if plugin.Policy.Installation != "AVAILABLE" || plugin.Policy.Authentication != "ON_INSTALL" || plugin.Category != "Productivity" {
		t.Errorf("Then marketplace availability metadata = %#v", plugin)
	}
	for _, field := range []string{"auth", "network", "permissions", "telemetry"} {
		if _, exists := fields[field]; exists {
			t.Errorf("Then catalog must omit %q", field)
		}
		if _, exists := rawCatalog.Plugins[0][field]; exists {
			t.Errorf("Then plugin entry must omit %q", field)
		}
	}
}
