package toolmanager

import (
  "encoding/json"
  "fmt"
  "io"
  "net/http"
  "os"
  "path/filepath"
  "strings"
)

type ghRelease struct {
  TagName string `json:"tag_name"`
  Assets  []struct {
    Name               string `json:"name"`
    BrowserDownloadURL string `json:"browser_download_url"`
  } `json:"assets"`
}

// EnsurePluginDir makes sure PluginDir exists and returns its absolute path.
func EnsurePluginDir() (string, error) {
  dir := PluginDir
  if !filepath.IsAbs(dir) {
    wd, err := os.Getwd()
    if err != nil {
      return "", err
    }
    dir = filepath.Join(wd, dir)
  }
  if err := os.MkdirAll(dir, 0o755); err != nil {
    return "", err
  }
  return dir, nil
}

// DownloadReleaseSO downloads all .so assets for the given repo@tag (or "latest")
// into your PluginDir.  If tag=="latest" we hit /releases/latest, otherwise
// /releases/tags/{tag}.  Returns a slice of local filenames.
func DownloadReleaseSO(repoURL, tag string) ([]string, error) {
  owner, repo, err := parseGitHubRepo(repoURL)
  if err != nil {
    return nil, err
  }

  apiPath := "latest"
  if tag != "" && tag != "latest" {
    apiPath = "tags/" + tag
  }
  apiURL := fmt.Sprintf(
    "https://api.github.com/repos/%s/%s/releases/%s",
    owner, repo, apiPath,
  )

  resp, err := http.Get(apiURL)
  if err != nil {
    return nil, err
  }
  defer resp.Body.Close()
  if resp.StatusCode != http.StatusOK {
    return nil, fmt.Errorf("GitHub API %q returned %d", apiURL, resp.StatusCode)
  }

  var rel ghRelease
  if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
    return nil, err
  }

  pluginsDir, err := EnsurePluginDir()
  if err != nil {
    return nil, err
  }

  var downloaded []string
  for _, a := range rel.Assets {
    if !strings.HasSuffix(a.Name, ".so") {
      continue
    }
    dst := filepath.Join(pluginsDir, a.Name)
    // skip if already present
    if _, err := os.Stat(dst); err == nil {
      downloaded = append(downloaded, dst)
      continue
    }
    out, err := os.Create(dst)
    if err != nil {
      return nil, err
    }
    defer out.Close()

    dl, err := http.Get(a.BrowserDownloadURL)
    if err != nil {
      return nil, err
    }
    defer dl.Body.Close()
    if dl.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("download %q failed: %d", a.Name, dl.StatusCode)
    }
    if _, err := io.Copy(out, dl.Body); err != nil {
      return nil, err
    }
    downloaded = append(downloaded, dst)
  }
  if len(downloaded) == 0 {
    return nil, fmt.Errorf("no .so assets found in %s@%s", repoURL, rel.TagName)
  }
  return downloaded, nil
}
