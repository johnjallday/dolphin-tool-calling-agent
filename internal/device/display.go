package device

import (
	"os/exec"
	"strings"
)

type DisplayDevice struct {
	Name       string
	Resolution string
	IsMain     bool
}

func GetCurrentDisplay() *DisplayDevice {
	out, err := exec.Command("system_profiler", "SPDisplaysDataType").Output()
	if err != nil {
		return nil
	}
	lines := strings.Split(string(out), "\n")
	var inDisplaysSection bool
	var display *DisplayDevice

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "Displays:" {
			inDisplaysSection = true
			continue
		}
		if inDisplaysSection {
			// Look for name line like 'SAMSUNG:'
			if strings.HasSuffix(line, ":") && !strings.Contains(line, "Display") {
				// Found a display name
				name := strings.TrimSuffix(line, ":")
				display = &DisplayDevice{Name: name}
				// Now look ahead for resolution and main
				for j := i + 1; j < len(lines); j++ {
					subline := strings.TrimSpace(lines[j])
					if strings.HasSuffix(subline, ":") && !strings.Contains(subline, "Display") {
						break // next display
					}
					if strings.HasPrefix(subline, "Resolution:") {
						display.Resolution = strings.TrimSpace(subline[len("Resolution:"):])
					}
					if strings.HasPrefix(subline, "Main Display:") && strings.Contains(subline, "Yes") {
						display.IsMain = true
					}
				}
				break // Only return the first display, or you can build a slice if you want all
			}
		}
	}
	return display
}
