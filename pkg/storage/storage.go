package storage

import (
	"gitlab.com/navyx/tools/pgexport/pkg/model"
)

// ReplicateProgress keeps the replicate progress
type ReplicateProgress struct {
	Progress int64
	UpdateAt int64
}

// Storage defines the abstract interface for accessing storage
type Storage interface {
	Close() error
	GetTableMaxID(tableName string) (int64, error)
	GetReplicateProgress(taskName string, tableName string) (*ReplicateProgress, error)
	UpdateReplicateProgress(taskName string, tableName string, progress ReplicateProgress) error
	SavePrebookings(tableName string, prebookings []*model.PrebookingRecord) error
	SaveBookingRecords(tableName string, bookings []*model.BookingRecord) error
	SaveBookingConfirmRecords(tableName string, confirms []*model.BookingConfirmRecord) error
	SaveRouteScheduleRecords(tableName string, schedules []*model.RouteScheduleRecord) error
	SaveBookingSummaryRecords(tableName string, records []*model.BookingSummaryRecord) error
	SaveBookingConfirmSummaryRecords(tableName string, records []*model.BookingConfirmSummaryRecord) error
	SaveDocumentSummaryRecords(tableName string, records []*model.BookingSummaryRecord) error
	SaveRateRecords(tableName string, records []*model.RateRecord) error
	SaveRateProviders(tableName string, records []*model.ProviderRecord) error
}
