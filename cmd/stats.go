package cmd

//Stats contains the statistics about the installed mods
type Stats struct {
	Installed int `json:"installed"`
	Updated   int `json:"updated"`
	Removed   int `json:"removed"`

	Warnings []string `json:"warnings,omitempty"`
}

//AddWarning to the statistics
func (stats *Stats) AddWarning(warning string) {
	stats.Warnings = append(stats.Warnings, warning)
}
