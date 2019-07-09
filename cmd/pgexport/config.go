package main

import (
	"gitlab.com/navyx/tools/pgexport/pkg/storage"
	"gitlab.com/navyx/tools/pgexport/pkg/task"
)

// AppConfig hold the application settings
type AppConfig struct {
	Database              storage.DBConfig `yaml:"database"`
	Prebooking            task.Config      `yaml:"prebooking`
	Booking               task.Config      `yaml:"booking""`
	BookingConfirm        task.Config      `yaml:"booking_confirm""`
	BookingSummary        task.Config      `yaml:"booking_summary"`
	BookingConfirmSummary task.Config      `yaml:"booking_confirm_summary"`
	DocumentSummary       task.Config      `yaml:"document_summary"`
	Rate                  task.Config      `yaml:"rate"`
	RouteSchedule         task.Config      `yaml:"route_schedule"`
	RateProvider          task.Config      `yaml:"rate_provider"`
}
