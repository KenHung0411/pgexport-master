package main

import "gitlab.com/navyx/tools/pgexport/pkg/storage"

// AppConfig hold the application settings
type AppConfig struct {
	Database storage.DBConfig `yaml:"database"`
}
