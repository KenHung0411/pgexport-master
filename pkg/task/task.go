package task

import "gitlab.com/navyx/tools/pgexport/pkg/model"

// TaskConfig defines the configuration of a task
type Config struct {
	Enable    bool          `yaml:"enable"`
	Source    model.Config  `yaml:"source"`
	Tables    []TableConfig `yaml:"tables"`
	BatchSize int           `yaml:"batch_size"`
}

// TableConfig defines the source and destination table
type TableConfig struct {
	From string `yaml:"from"`
	To   string `yaml:"to"`
}

// Task defines the abstract interfaces
type Task interface {
	Execute() error
}

const (
	PREBOOKING_TASK_NAME              = "prebooking"
	BOOKING_TASK_NAME                 = "booking"
	BOOKING_CONFIRM_TASK_NAME         = "booking_confirm"
	BOOKING_SUMMARY_TASK_NAME         = "booking_summary"
	BOOKING_CONFIRM_SUMMARY_TASK_NAME = "booking_confirm_summary"
	DOCUMENT_SUMMARY_TASK_NAME        = "document_summary"
	ROUTE_SCHEDULE_TASK_NAME          = "route_schedule"
	RATE_TASK_NAME                    = "rate"
	RATE_PROVIDER_TASK_NAME           = "rate_provider"
)
