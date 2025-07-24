package toolmanager

import (
  "encoding/json"
  "errors"
	"io/fs"
  "fmt"
  "net/http"
  "path/filepath"
  "plugin"
  "strings"

  "github.com/johnjallday/dolphin-tool-calling-agent/pkg/tools"
)

var (
  // PluginDir is where your .so plugin files live. Change as needed.
  PluginDir = "./plugins"
)

// ToolPacks scans PluginDir (recursively), opens every “.so” it finds,
// looks up the PluginPackage() symbol, and returns the slice of ToolPackage.
func ToolPacks() ([]tools.ToolPackage, error) {
  var packs []tools.ToolPackage

  walkErr := filepath.WalkDir(PluginDir, func(path string, d fs.DirEntry, err error) error {
    if err != nil {
      // e.g. permission error
      return err
    }
    if d.IsDir() {
      return nil
    }
    if filepath.Ext(path) != ".so" {
      return nil
    }

    p, err := plugin.Open(path)
    if err != nil {
      fmt.Printf("⚠️  failed to open plugin %q: %v\n", path, err)
      return nil // keep walking
    }

    sym, err := p.Lookup("PluginPackage")
    if err != nil {
      fmt.Printf("⚠️  %q does not export PluginPackage(): %v\n", path, err)
      return nil
    }

    constructor, ok := sym.(func() tools.ToolPackage)
    if !ok {
      fmt.Printf("⚠️  %q: PluginPackage has wrong signature\n", path)
      return nil
    }

    pkg := constructor()
    packs = append(packs, pkg)
    return nil
  })

  if walkErr != nil {
    return nil, fmt.Errorf("walking plugin dir %q: %w", PluginDir, walkErr)
  }
  if len(packs) == 0 {
    return nil, errors.New("no plugins loaded")
  }
  return packs, nil
}

// CheckVersion fetches the latest GitHub release for each discovered
// ToolPackage and compares it to the local Version field.
// Prints out which are up-to-date and which have updates available.
func CheckVersion() error {
  packs, err := ToolPacks()
  if err != nil {
    return err
  }

  type ghRelease struct {
    TagName string `json:"tag_name"`
  }

  for _, pkg := range packs {
    owner, repo, err := parseGitHubRepo(pkg.Link)
    if err != nil {
      fmt.Printf("%s: invalid GitHub link %q: %v\n", pkg.Name, pkg.Link, err)
      continue
    }

    url := fmt.Sprintf(
      "https://api.github.com/repos/%s/%s/releases/latest",
      owner, repo,
    )

    resp, err := http.Get(url)
    if err != nil {
      fmt.Printf("%s: HTTP error: %v\n", pkg.Name, err)
      continue
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
      fmt.Printf("%s: unexpected HTTP status %s\n", pkg.Name, resp.Status)
      continue
    }

    var rel ghRelease
    if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
      fmt.Printf("%s: JSON decode error: %v\n", pkg.Name, err)
      continue
    }

    if rel.TagName == pkg.Version {
      fmt.Printf("%s: up-to-date (%s)\n", pkg.Name, pkg.Version)
    } else {
      fmt.Printf("%s: update available: local v%s → remote %s\n",
        pkg.Name, pkg.Version, rel.TagName)
    }
  }

  return nil
}

// parseGitHubRepo extracts owner and repo from a GitHub HTTPS URL.
// e.g. https://github.com/foo/bar or https://github.com/foo/bar/ → ("foo","bar",nil)
func parseGitHubRepo(raw string) (owner, repo string, err error) {
  const prefix = "https://github.com/"
  if !strings.HasPrefix(raw, prefix) {
    return "", "", fmt.Errorf("must start with %q", prefix)
  }
  suffix := strings.TrimPrefix(raw, prefix)
  suffix = strings.TrimSuffix(suffix, "/")
  parts := strings.Split(suffix, "/")
  if len(parts) < 2 {
    return "", "", fmt.Errorf("expect owner/repo, got %q", suffix)
  }
  return parts[0], parts[1], nil
}
