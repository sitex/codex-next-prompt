package tests_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

type pluginManifest struct {
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	Description string   `json:"description"`
	Author      author   `json:"author"`
	Homepage    string   `json:"homepage"`
	Repository  string   `json:"repository"`
	License     string   `json:"license"`
	Keywords    []string `json:"keywords"`
}

type author struct {
	Name string `json:"name"`
}

type hookManifest struct {
	Hooks map[string][]hookGroup `json:"hooks"`
}

type hookGroup struct {
	Matcher *string       `json:"matcher"`
	Hooks   []hookCommand `json:"hooks"`
}

type hookCommand struct {
	Type           string  `json:"type"`
	Command        string  `json:"command"`
	CommandWindows string  `json:"commandWindows"`
	Timeout        int     `json:"timeout"`
	StatusMessage  *string `json:"statusMessage"`
}

func Test_PluginManifest_declares_public_metadata_and_uses_default_hook_discovery(t *testing.T) {
	// Given
	manifestPath := filepath.Join("..", ".codex-plugin", "plugin.json")

	// When
	data := readFile(t, manifestPath)
	var manifest pluginManifest
	decodeJSON(t, manifestPath, data, &manifest)
	var fields map[string]json.RawMessage
	decodeJSON(t, manifestPath, data, &fields)

	// Then
	if manifest.Name != "codex-next-prompt" {
		t.Errorf("Then name = %q, want codex-next-prompt", manifest.Name)
	}
	if manifest.Version != "0.1.0" {
		t.Errorf("Then version = %q, want 0.1.0", manifest.Version)
	}
	if manifest.Description == "" || manifest.Author.Name == "" || manifest.Homepage == "" || manifest.Repository == "" || manifest.License != "MIT" || len(manifest.Keywords) == 0 {
		t.Error("Then public description, author, homepage, repository, MIT license, and keywords must be present")
	}
	for _, field := range []string{"hooks", "capabilities", "mcpServers", "apps", "auth"} {
		if _, exists := fields[field]; exists {
			t.Errorf("Then manifest must omit %q", field)
		}
	}
}

func Test_HookManifest_registers_fail_open_session_start_and_stop_commands(t *testing.T) {
	// Given
	manifestPath := filepath.Join("..", "hooks", "hooks.json")

	// When
	data := readFile(t, manifestPath)
	var manifest hookManifest
	decodeJSON(t, manifestPath, data, &manifest)

	// Then
	sessionStart := singleHookCommand(t, manifest, "SessionStart")
	if sessionStart.group.Matcher == nil || *sessionStart.group.Matcher != "startup|resume|clear|compact" {
		t.Errorf("Then SessionStart matcher = %v", sessionStart.group.Matcher)
	}
	assertCommand(t, sessionStart.command, "session-start")
	if sessionStart.command.StatusMessage == nil || *sessionStart.command.StatusMessage == "" {
		t.Error("Then SessionStart statusMessage must be present")
	}

	stop := singleHookCommand(t, manifest, "Stop")
	if stop.group.Matcher != nil {
		t.Errorf("Then Stop matcher = %q, want omitted", *stop.group.Matcher)
	}
	assertCommand(t, stop.command, "stop")
	if stop.command.StatusMessage != nil {
		t.Errorf("Then Stop statusMessage = %q, want omitted", *stop.command.StatusMessage)
	}
}

func Test_HookLaunchers_are_present_and_posix_launcher_is_executable(t *testing.T) {
	// Given
	launcherPaths := []string{
		filepath.Join("..", "hooks", "run"),
		filepath.Join("..", "hooks", "run.cmd"),
	}

	for _, launcherPath := range launcherPaths {
		t.Run(filepath.Base(launcherPath), func(t *testing.T) {
			// When
			info, err := os.Stat(launcherPath)

			// Then
			if err != nil {
				t.Fatalf("Then launcher %q must exist: %v", launcherPath, err)
			}
			if filepath.Base(launcherPath) == "run" && info.Mode().Perm()&0o111 == 0 {
				t.Errorf("Then POSIX launcher mode = %o, want executable", info.Mode().Perm())
			}
		})
	}
}

type registeredHook struct {
	group   hookGroup
	command hookCommand
}

func singleHookCommand(t *testing.T, manifest hookManifest, event string) registeredHook {
	t.Helper()

	groups := manifest.Hooks[event]
	if len(groups) != 1 {
		t.Fatalf("Then %s groups = %d, want 1", event, len(groups))
	}
	if len(groups[0].Hooks) != 1 {
		t.Fatalf("Then %s commands = %d, want 1", event, len(groups[0].Hooks))
	}
	return registeredHook{group: groups[0], command: groups[0].Hooks[0]}
}

func assertCommand(t *testing.T, command hookCommand, subcommand string) {
	t.Helper()

	wantPOSIX := "${PLUGIN_ROOT}/hooks/run " + subcommand
	if command.Type != "command" || command.Command != wantPOSIX {
		t.Errorf("Then %s command = %#v, want type command and %q", subcommand, command, wantPOSIX)
	}
	wantWindows := `cmd.exe /d /s /c ""%PLUGIN_ROOT%\hooks\run.cmd" ` + subcommand + `"`
	if command.CommandWindows != wantWindows {
		t.Errorf("Then %s commandWindows = %q, want %q", subcommand, command.CommandWindows, wantWindows)
	}
	if command.Timeout <= 0 || command.Timeout > 10 {
		t.Errorf("Then %s timeout = %d, want 1..10 seconds", subcommand, command.Timeout)
	}
}

func readFile(t *testing.T, path string) []byte {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Given required packaging file %q: read: %v", path, err)
	}
	return data
}

func decodeJSON[T any](t *testing.T, path string, data []byte, destination *T) {
	t.Helper()

	if err := json.Unmarshal(data, destination); err != nil {
		t.Fatalf("When decode %q: %v", path, err)
	}
}
