package model

import (
	"database/sql"

	"github.com/lib/pq"
	pb "gitlab.com/navyx/tools/pgexport/pkg/proto/prebooking_proto"
)

// Config holds the configuration used for instantiating a new Database.
type Config struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	SSLMode  string `yaml:"sslmode"`
}

type PrebookingInterface interface {
	ListPrebookings(tableName string, since int64, limit int) (*PrebookingResult, error)
}

type BookingInterface interface {
	GetBookingRecords(tableName string, since int64, limit int) (*BookingResult, error)
	GetBookingConfirmRecords(tableName string, since int64, limit int) (*BookingConfirmResult, error)
}

type RouteScheduleInterface interface {
	GetRouteSchedules(tableName string, since int64, limit int) (*RouteScheduleResult, error)
}

type RateInterface interface {
	GetRates(tableName string, since int64, limit int) (*RateResult, error)
}

type BookingSummaryInterface interface {
	GetBookingSummaryRecords(tableName string, since int64, limit int) (*BookingSummaryResult, error)
	GetBookingConfirmSummaryRecords(tableName string, since int64, limit int) (*BookingConfirmSummaryResult, error)
	GetDocumentSummaryRecords(tableName string, since int64, limit int) (*BookingSummaryResult, error)
}

type RateProviderInterface interface {
	ListRateProviders(tableName string, since int64, limit int) (*RateProviderResult, error)
}

// DB defines the abstract interface for accessing database
type DB interface {
	Close() error
	PrebookingInterface
	BookingInterface
	RouteScheduleInterface
	BookingSummaryInterface
	RateInterface
	RateProviderInterface
}

type PrebookingRecord struct {
	SubmitUserID     int64
	SourcePlatformID int64
	Prebooking       *pb.Prebooking
}

// PrebookingResult model
type PrebookingResult struct {
	Since       int64
	Limit       int
	Next        int64
	Prebookings []*PrebookingRecord
}

// BookingRecord data model
type BookingRecord struct {
	ID                int64
	PlatformID        int64
	Region            string
	ApplicationNumber string
	Version           int64
	IsDraft           bool
	CreatedBy         int64
	UpdatedBy         int64
	CreatedAt         int64
	UpdatedAt         int64
	Data              []byte
}

// BookingResult model
type BookingResult struct {
	Since   int64
	Limit   int
	Next    int64
	Records []*BookingRecord
}

// BookingConfirmRecord data model
type BookingConfirmRecord struct {
	ID                int64
	Region            string
	ApplicationNumber string
	Version           int64
	CreatedAt         int64
	Data              []byte
}

// BookingResult model
type BookingConfirmResult struct {
	Since   int64
	Limit   int
	Next    int64
	Records []*BookingConfirmRecord
}

// RouteScheduleRecord data model
type RouteScheduleRecord struct {
	ID        int64
	Carrier   string
	Vessel    string
	Voyage    string
	LloydCode string
	Data      []byte
}

// RouteScheduleResult model
type RouteScheduleResult struct {
	Since   int64
	Limit   int
	Next    int64
	Records []*RouteScheduleRecord
}

// BookingSummaryRecord data model
type BookingSummaryRecord struct {
	ID          int64
	Region      string
	PlatformID  int64
	AppNumber   string
	Carrier     string
	Status      int64
	IsDeleted   bool
	CreatedBy   int64
	UpdatedBy   int64
	CreatedAt   int64
	UpdatedAt   int64
	Data        []byte
	ConfirmData *[]byte
}

//BookingSummaryResult model
type BookingSummaryResult struct {
	Since   int64
	Limit   int
	Next    int64
	Records []*BookingSummaryRecord
}

// BookingConfirmSummaryRecord data model
type BookingConfirmSummaryRecord struct {
	ID            int64
	Region        string
	AppNumber     string
	BookingNumber string
	Status        int64
	Timestamp     int64
	Data          []byte
}

//BookingConfirmSummaryResult model
type BookingConfirmSummaryResult struct {
	Since   int64
	Limit   int
	Next    int64
	Records []*BookingConfirmSummaryRecord
}

// RateRecord data model
type RateRecord struct {
	ID                      int64
	PlaceOfReceopt          string
	PlaceOfDelivery         string
	PortOfLoading           string
	PortOfDischarge         string
	CarrierSCAC             string
	ContainerType           int32
	Price                   float64
	EffectiveDate           pq.NullTime
	ExpiryDate              pq.NullTime
	ProviderID              int64
	ContractNumber          sql.NullString
	ServiceCode             sql.NullString
	ServiceMode             sql.NullString
	Commodity               sql.NullString
	Demurrage               sql.NullString
	Detention               sql.NullString
	Outport                 sql.NullString
	Remarks                 sql.NullString
	Included                sql.NullString
	SubjectTo               sql.NullString
	ServiceFee              float64
	ServiceFeeEffectiveDate pq.NullTime
	ServiceFeeExpiryDate    pq.NullTime
	PromotionDiscount       sql.NullFloat64
	PromotionEffectiveDate  pq.NullTime
	PromotionExpiryDate     pq.NullTime
	Version                 int64
	CreatedAt               pq.NullTime
}

// RateResult model
type RateResult struct {
	Since   int64
	Limit   int
	Next    int64
	Records []*RateRecord
}

type PrebookingContainer struct {
	ID            string `json:"id,omitempty"`
	Status        int32  `json:"status,omitempty"`
	Type          int32  `json:"type,omitempty"`
	BookingNumber string `json:"booking_number,omitempty"`
	Document      string `json:"document,omitempty"`
	UpdatedAt     int64  `json:"updated_at,omitempty"`
	RejectReason  string `json:"reject_reason,omitempty"`
}

// PrebookingContainers holds several 'Container' together.
type PrebookingContainers struct {
	Containers []PrebookingContainer `json:"containers"`
}

type ProviderRecord struct {
	ID      int64
	Code    string
	Name    string
	Version int64
}

// RateProviderResult model
type RateProviderResult struct {
	Since   int64
	Limit   int
	Next    int64
	Records []*ProviderRecord
}
