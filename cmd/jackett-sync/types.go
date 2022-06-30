package main

// TrackerDefinition describes a Jackett tracker.
type TrackerDefinition struct {
	ID          string `yaml:"id"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Language    string `yaml:"language"`
	Type        string `yaml:"type"`
}
