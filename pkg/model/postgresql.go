package model

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	_ "github.com/lib/pq" // import postgres driver
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gitlab.com/navyx/tools/pgexport/pkg/proto/booking_proto"
	pb "gitlab.com/navyx/tools/pgexport/pkg/proto/prebooking_proto"
)

// PostgreSQL holds the connection pool to the database
type PostgreSQL struct {
	db *sql.DB
}

// NewPostgreSQL returns the PostgreSQL instance
func NewPostgreSQL(config Config) (DB, error) {
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		url.PathEscape(config.Username),
		url.PathEscape(config.Password),
		url.PathEscape(config.Host),
		config.Port,
		url.PathEscape(config.Database),
		url.QueryEscape(config.SSLMode),
	)

	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, errors.Wrap(err, "open connection to database")
	}

	if err := db.Ping(); err != nil {
		return nil, errors.Wrap(err, "ping database")
	}
	return &PostgreSQL{db}, nil
}

// Close disconnect the connection
func (postgres *PostgreSQL) Close() error {
	if postgres != nil {
		if postgres.db != nil {
			return postgres.db.Close()
		}
	}
	return nil
}

// ListPrebookings returns the collection of pre-bookings
func (postgres *PostgreSQL) ListPrebookings(tableName string, since int64, limit int) (*PrebookingResult, error) {
	log.Tracef("ListPreBookings since: %d limit: %d", since, limit)

	query := fmt.Sprintf(`
SELECT b.id,
       b.uuid,
       b.source_platform_id,
       b.status,
       b.shipment_name,
       b.submit_user_id,
       b.submit_user_name,
       b.volume,
       b.volume_unit,
       b.weight,
       b.weight_unit,
       b.created_at,
       b.updated_at,
       b.notification_groups_bytes,
       b.notes,
       b.history,
       q.content,
       b.containers,
       b.updated_by_rbd_user_id,
       b.updated_by_rbd_platform_id,
       b.updated_by_rbd_at
FROM %s b
LEFT JOIN quotes q ON (q.id = b.quote_id)
WHERE b.id > $1
ORDER BY b.id
LIMIT $2
`, tableName)

	var id, sourcePlatformID int64
	var uuid string
	var shipmentName, submitUserName sql.NullString
	var status int32
	var submitUserID sql.NullInt64
	var volume, volumeUnit, weight, weightUnit sql.NullInt64
	var notificationGroups, notes, history, content, containers []byte
	var createdAt, updatedAt time.Time
	var rbdUpdateUserID, rbdUpdatePlatformID, rbdUpdatedAt sql.NullInt64

	rows, err := postgres.db.Query(query, since, limit)
	if err != nil {
		return nil, errors.Wrapf(err, "execute query")
	}
	defer rows.Close()

	records := make([]*PrebookingRecord, 0)
	for rows.Next() {
		err := rows.Scan(&id, &uuid, &sourcePlatformID, &status, &shipmentName, &submitUserID, &submitUserName,
			&volume, &volumeUnit, &weight, &weightUnit, &createdAt, &updatedAt,
			&notificationGroups, &notes, &history, &content,
			&containers, &rbdUpdateUserID, &rbdUpdatePlatformID, &rbdUpdatedAt)
		if err != nil {
			return nil, errors.Wrapf(err, "scan columns")
		}
		prebooking := &pb.Prebooking{
			Id:               uuid,
			Status:           pb.PrebookingStatus(status),
			Quote:            &pb.Quote{},
			UpdatePlatformId: uint32(sourcePlatformID),
		}
		if shipmentName.Valid {
			prebooking.ShipmentName = truncateNonASCII(shipmentName.String)
		}
		if submitUserID.Valid {
			prebooking.UpdateUserId = uint32(submitUserID.Int64)
		}
		if submitUserName.Valid {
			prebooking.UpdateUserName = truncateNonASCII(submitUserName.String)
		}
		if volume.Valid {
			prebooking.Volume = uint32(volume.Int64)
		}
		if volumeUnit.Valid {
			prebooking.VolumeUnit = pb.VolumeUnitType(int32(volumeUnit.Int64))
		}
		if weight.Valid {
			prebooking.Weight = uint32(weight.Int64)
		}
		if weightUnit.Valid {
			prebooking.WeightUnit = pb.WeightUnitType(int32(weightUnit.Int64))
		}
		if pbTime, err := ptypes.TimestampProto(createdAt); err != nil {
			log.Warnf("convert created_at timestamp err: %s\n", err)
			return nil, errors.Wrap(err, "convert created_at timestamp")
		} else {
			prebooking.CreatedAt = pbTime
		}
		if pbTime, err := ptypes.TimestampProto(updatedAt); err != nil {
			log.Warnf("convert updated_at timestamp err: %s\n", err)
			return nil, errors.Wrap(err, "convert updated_at timestamp")
		} else {
			prebooking.UpdatedAt = pbTime
		}

		var notificationGroupList pb.NotificationGroupList
		if err := proto.Unmarshal(notificationGroups, &notificationGroupList); err != nil {
			log.Warnf("unmarshal notification group list err: %s\n", err)
			return nil, errors.Wrap(err, "unmarshal notification group list")
		} else {
			prebooking.NotificationGroups = append(prebooking.NotificationGroups, notificationGroupList.NotificationGroups...)
		}

		var pbNotes pb.Dictionary
		if err := proto.Unmarshal(notes, &pbNotes); err != nil {
			log.Warnf("unmarshal notes err: %s\n", err)
			return nil, errors.Wrap(err, "unmarshal notes")
		} else {
			prebooking.Notes = pbNotes.Pairs
		}

		var pbHistoryList pb.UpdateHistoryList
		if err := proto.Unmarshal(history, &pbHistoryList); err != nil {
			log.Warnf("unmarshal update history list err: %s\n", err)
			return nil, errors.Wrap(err, "unmarshal update history list")
		} else {
			prebooking.History = append(prebooking.History, pbHistoryList.List...)
		}

		if err := proto.Unmarshal(content, prebooking.Quote); err != nil {
			log.Warnf("unmarshal quote err: %s\n", err)
			return nil, errors.Wrap(err, "unmarshal quote")
		}

		var containerList PrebookingContainers
		if err := json.Unmarshal(containers, &containerList); err != nil {
			log.Warnf("unmarshal containers err: %s\n", err)
			return nil, errors.Wrap(err, "unmarshal containers")
		} else {
			pbContainerList := make([]*pb.Container, 0)
			for _, c := range containerList.Containers {
				pbContainerList = append(pbContainerList, &pb.Container{
					Id:           c.ID,
					Status:       pb.ContainerStatus(c.Status),
					Type:         pb.ContainerType(c.Type),
					BookingNo:    c.BookingNumber,
					Document:     c.Document,
					UpdatedAt:    c.UpdatedAt,
					RejectReason: c.RejectReason,
				})
			}
			prebooking.Containers = pbContainerList
		}

		prebooking.UpdatedByRbdUserId = rbdUpdateUserID.Int64
		prebooking.UpdatedByRbdPlatformId = rbdUpdatePlatformID.Int64
		prebooking.UpdatedByRbdAt = rbdUpdatedAt.Int64

		userID := submitUserID.Int64
		platformID := sourcePlatformID
		if len(prebooking.History) > 0 {
			userID = int64(prebooking.History[0].UpdateUserId)
			platformID = int64(prebooking.History[0].UpdatePlatformId)
		}

		records = append(records, &PrebookingRecord{
			SubmitUserID:     userID,
			SourcePlatformID: platformID,
			Prebooking:       prebooking,
		})
	}
	err = rows.Err()
	if err != nil {
		return nil, errors.Wrapf(err, "row iterator")
	}

	maxID := since
	if len(records) > 0 {
		maxID = id
	}
	log.Tracef("ListPreBookings got %d records, nextID=%d", len(records), maxID)
	return &PrebookingResult{Since: since, Limit: limit, Next: maxID, Prebookings: records}, nil
}

// GetBookingRecords returns the collection from booking/si/vgm
func (postgres *PostgreSQL) GetBookingRecords(tableName string, since int64, limit int) (*BookingResult, error) {
	log.Tracef("GetBookingRecords %s since: %d limit: %d", tableName, since, limit)

	query := fmt.Sprintf(`
SELECT id, region, application_number, version, is_draft, platform_id, created_by, updated_by, created_at, updated_at, data
FROM %s
WHERE id > $1 ORDER BY id LIMIT $2
`, tableName)

	var id, platformID, version int64
	var appNumber string
	var region sql.NullString
	var isDraft bool
	var createdBy, updatedBy, createdAt, updatedAt int64
	var data []byte

	rows, err := postgres.db.Query(query, since, limit)
	if err != nil {
		return nil, errors.Wrapf(err, "execute query")
	}
	defer rows.Close()

	records := make([]*BookingRecord, 0)
	for rows.Next() {
		err := rows.Scan(&id, &region, &appNumber, &version, &isDraft, &platformID,
			&createdBy, &updatedBy, &createdAt, &updatedAt, &data)
		if err != nil {
			return nil, errors.Wrapf(err, "scan columns")
		}

		record := &BookingRecord{
			ID:                id,
			PlatformID:        platformID,
			Region:            region.String,
			ApplicationNumber: appNumber,
			Version:           version,
			IsDraft:           isDraft,
			CreatedBy:         createdBy,
			UpdatedBy:         updatedBy,
			CreatedAt:         createdAt,
			UpdatedAt:         updatedAt,
			Data:              data,
		}
		records = append(records, record)
	}

	err = rows.Err()
	if err != nil {
		return nil, errors.Wrapf(err, "row iterator")
	}

	maxID := since
	if len(records) > 0 {
		maxID = id
	}
	log.Tracef("GetBookingRecords %s got %d records, nextID=%d", tableName, len(records), maxID)
	return &BookingResult{Since: since, Limit: limit, Next: maxID, Records: records}, nil
}

// GetBookingConfirmRecords returns the collection from booking/si/vgm
func (postgres *PostgreSQL) GetBookingConfirmRecords(tableName string, since int64, limit int) (*BookingConfirmResult, error) {
	log.Tracef("GetBookingConfirmRecords %s since: %d limit: %d", tableName, since, limit)

	query := fmt.Sprintf(`
SELECT id, region, application_number, version, created_at, data
FROM %s
WHERE id > $1 ORDER BY id LIMIT $2
`, tableName)

	var id, version int64
	var region sql.NullString
	var appNumber string
	var createdAt int64
	var data []byte

	rows, err := postgres.db.Query(query, since, limit)
	if err != nil {
		return nil, errors.Wrapf(err, "execute query")
	}
	defer rows.Close()

	records := make([]*BookingConfirmRecord, 0)
	for rows.Next() {
		err := rows.Scan(&id, &region, &appNumber, &version, &createdAt, &data)
		if err != nil {
			return nil, errors.Wrapf(err, "scan columns")
		}

		record := &BookingConfirmRecord{
			ID:                id,
			Region:            region.String,
			ApplicationNumber: appNumber,
			Version:           version,
			CreatedAt:         createdAt,
			Data:              data,
		}
		records = append(records, record)
	}

	err = rows.Err()
	if err != nil {
		return nil, errors.Wrapf(err, "row iterator")
	}

	maxID := since
	if len(records) > 0 {
		maxID = id
	}
	log.Tracef("GetBookingConfirmRecords %s got %d records, nextID=%d", tableName, len(records), maxID)
	return &BookingConfirmResult{Since: since, Limit: limit, Next: maxID, Records: records}, nil
}

// GetRouteSchedules returns the collection from booking/si/vgm
func (postgres *PostgreSQL) GetRouteSchedules(tableName string, since int64, limit int) (*RouteScheduleResult, error) {
	log.Tracef("GetRouteSchedules %s since: %d limit: %d", tableName, since, limit)

	query := fmt.Sprintf(`
SELECT id, data
FROM %s
WHERE id > $1 ORDER BY id LIMIT $2
`, tableName)

	var id int64
	var data json.RawMessage

	rows, err := postgres.db.Query(query, since, limit)
	if err != nil {
		return nil, errors.Wrapf(err, "execute query")
	}
	defer rows.Close()

	records := make([]*RouteScheduleRecord, 0)
	for rows.Next() {
		err := rows.Scan(&id, &data)
		if err != nil {
			return nil, errors.Wrapf(err, "scan columns")
		}

		schedule := &booking_proto.CarrierSchedule{}
		if err = json.Unmarshal(data, &schedule); err != nil {
			return nil, errors.Wrapf(err, "parse route schedule for row %d", id)
		}

		record := &RouteScheduleRecord{
			ID:        id,
			Carrier:   strings.ToUpper(schedule.GetVessel().GetCarrier()),
			Vessel:    strings.ToUpper(schedule.GetVessel().GetName()),
			Voyage:    strings.ToUpper(schedule.GetVoyage()),
			LloydCode: schedule.GetVessel().GetLloydCode(),
			Data:      data,
		}
		records = append(records, record)
	}

	err = rows.Err()
	if err != nil {
		return nil, errors.Wrapf(err, "row iterator")
	}

	maxID := since
	if len(records) > 0 {
		maxID = id
	}

	log.Tracef("GetRouteSchedules %s got %d records, nextID=%d", tableName, len(records), maxID)
	return &RouteScheduleResult{Since: since, Limit: limit, Next: maxID, Records: records}, nil
}

// GetBookingSummaryRecords returns the booking summary records
func (postgres *PostgreSQL) GetBookingSummaryRecords(tableName string, since int64, limit int) (*BookingSummaryResult, error) {
	log.Tracef("GetBookingSummaryRecords %s since: %d limit: %d", tableName, since, limit)

	query := fmt.Sprintf(`
SELECT
    id, region, platform_id, application_number, carrier, status, is_deleted,
    created_by, updated_by, created_at, updated_at, data
FROM %s
WHERE id > $1 ORDER BY id LIMIT $2
`, tableName)

	var id, platformID, status int64
	var region sql.NullString
	var appNumber string
	var carrier sql.NullString
	var isDeleted sql.NullBool
	var createdBy, createdAt int64
	var updatedBy, updatedAt sql.NullInt64
	var data []byte

	rows, err := postgres.db.Query(query, since, limit)
	if err != nil {
		return nil, errors.Wrapf(err, "execute query")
	}
	defer rows.Close()

	records := make([]*BookingSummaryRecord, 0)
	for rows.Next() {
		err := rows.Scan(&id, &region, &platformID, &appNumber, &carrier, &status, &isDeleted,
			&createdBy, &updatedBy, &createdAt, &updatedAt, &data)
		if err != nil {
			return nil, errors.Wrapf(err, "scan columns")
		}

		record := &BookingSummaryRecord{
			ID:         id,
			Region:     region.String,
			PlatformID: platformID,
			AppNumber:  appNumber,
			Carrier:    carrier.String,
			Status:     status,
			IsDeleted:  isDeleted.Bool,
			CreatedBy:  createdBy,
			UpdatedBy:  updatedBy.Int64,
			CreatedAt:  createdAt,
			UpdatedAt:  updatedAt.Int64,
			Data:       data,
		}
		records = append(records, record)
	}

	err = rows.Err()
	if err != nil {
		return nil, errors.Wrapf(err, "row iterator")
	}

	maxID := since
	if len(records) > 0 {
		maxID = id
	}
	log.Tracef("GetBookingSummaryRecords %s got %d records, nextID=%d", tableName, len(records), maxID)
	return &BookingSummaryResult{Since: since, Limit: limit, Next: maxID, Records: records}, nil
}

// GetBookingConfirmSummaryRecords returns the booking confirm summary records
func (postgres *PostgreSQL) GetBookingConfirmSummaryRecords(tableName string, since int64, limit int) (*BookingConfirmSummaryResult, error) {
	log.Tracef("GetBookingConfirmSummaryRecords %s since: %d limit: %d", tableName, since, limit)

	query := fmt.Sprintf(`
SELECT id, region, application_number, booking_number, status, ts, data
FROM %s
WHERE id > $1 ORDER BY id LIMIT $2
`, tableName)

	var id, status int64
	var region sql.NullString
	var appNumber, bookingNumber string
	var timestamp sql.NullInt64
	var data []byte

	rows, err := postgres.db.Query(query, since, limit)
	if err != nil {
		return nil, errors.Wrapf(err, "execute query")
	}
	defer rows.Close()

	records := make([]*BookingConfirmSummaryRecord, 0)
	for rows.Next() {
		err := rows.Scan(&id, &region, &appNumber, &bookingNumber, &status, &timestamp, &data)
		if err != nil {
			return nil, errors.Wrapf(err, "scan columns")
		}

		record := &BookingConfirmSummaryRecord{
			ID:            id,
			Region:        region.String,
			AppNumber:     appNumber,
			BookingNumber: bookingNumber,
			Status:        status,
			Timestamp:     timestamp.Int64,
			Data:          data,
		}
		records = append(records, record)
	}

	err = rows.Err()
	if err != nil {
		return nil, errors.Wrapf(err, "row iterator")
	}

	maxID := since
	if len(records) > 0 {
		maxID = id
	}
	log.Tracef("GetBookingSummaryRecords %s got %d records, nextID=%d", tableName, len(records), maxID)
	return &BookingConfirmSummaryResult{Since: since, Limit: limit, Next: maxID, Records: records}, nil
}

// GetDocumentSummaryRecords returns si/vgm summary records
func (postgres *PostgreSQL) GetDocumentSummaryRecords(tableName string, since int64, limit int) (*BookingSummaryResult, error) {
	log.Tracef("GetDocumentSummaryRecords %s since: %d limit: %d", tableName, since, limit)

	query := fmt.Sprintf(`
SELECT
    id, region, platform_id, application_number, carrier, status, is_deleted,
    created_by, updated_by, created_at, updated_at, data, confirm
FROM %s
WHERE id > $1 ORDER BY id LIMIT $2
`, tableName)

	var id, platformID, status int64
	var region sql.NullString
	var appNumber string
	var carrier sql.NullString
	var isDeleted sql.NullBool
	var createdBy, createdAt int64
	var updatedBy, updatedAt sql.NullInt64
	var data []byte
	var confirmData *[]byte

	rows, err := postgres.db.Query(query, since, limit)
	if err != nil {
		return nil, errors.Wrapf(err, "execute query")
	}
	defer rows.Close()

	records := make([]*BookingSummaryRecord, 0)
	for rows.Next() {
		err := rows.Scan(&id, &region, &platformID, &appNumber, &carrier, &status, &isDeleted,
			&createdBy, &updatedBy, &createdAt, &updatedAt, &data, &confirmData)
		if err != nil {
			return nil, errors.Wrapf(err, "scan columns")
		}

		record := &BookingSummaryRecord{
			ID:          id,
			Region:      region.String,
			PlatformID:  platformID,
			AppNumber:   appNumber,
			Carrier:     carrier.String,
			Status:      status,
			IsDeleted:   isDeleted.Bool,
			CreatedBy:   createdBy,
			UpdatedBy:   updatedBy.Int64,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt.Int64,
			Data:        data,
			ConfirmData: confirmData,
		}
		records = append(records, record)
	}

	err = rows.Err()
	if err != nil {
		return nil, errors.Wrapf(err, "row iterator")
	}

	maxID := since
	if len(records) > 0 {
		maxID = id
	}
	log.Tracef("GetConfirmSummaryRecords %s got %d records, nextID=%d", tableName, len(records), maxID)
	return &BookingSummaryResult{Since: since, Limit: limit, Next: maxID, Records: records}, nil
}

// GetRates returns the collection from rates
func (postgres *PostgreSQL) GetRates(tableName string, since int64, limit int) (*RateResult, error) {
	log.Tracef("GetRates %s since: %d limit: %d", tableName, since, limit)

	query := fmt.Sprintf(`
SELECT
  r.id,r.origin,r.port_of_loading,r.port_of_discharge,r.destination,r.carrier_scac,r.container_type,
  r.provider_id,r.price,r.service_fee,r.service_code,r.contract_number,r.service_mode,
  r.effective_date,r.expiry_date,r.service_fee_effective_date,r.service_fee_expiry_date,
  r.promotion_discount,r.promotion_effective_date,r.promotion_expiry_date,
  r.commodity,r.demurrage,r.detention,r.outport,r.remarks,r.included,r.subject_to,r.version,r.created_at
FROM %s r
WHERE id > $1 ORDER BY id LIMIT $2
`, tableName)

	var id int64
	rows, err := postgres.db.Query(query, since, limit)
	if err != nil {
		return nil, errors.Wrapf(err, "execute query")
	}
	defer rows.Close()

	records := make([]*RateRecord, 0)
	for rows.Next() {
		r := &RateRecord{}
		err := rows.Scan(&id, &r.PlaceOfReceopt, &r.PortOfLoading, &r.PortOfDischarge, &r.PlaceOfDelivery,
			&r.CarrierSCAC, &r.ContainerType,
			&r.ProviderID, &r.Price, &r.ServiceFee, &r.ServiceCode, &r.ContractNumber, &r.ServiceMode,
			&r.EffectiveDate, &r.ExpiryDate, &r.ServiceFeeEffectiveDate, &r.ServiceFeeExpiryDate,
			&r.PromotionDiscount, &r.PromotionEffectiveDate, &r.PromotionExpiryDate,
			&r.Commodity, &r.Demurrage, &r.Detention, &r.Outport, &r.Remarks, &r.Included, &r.SubjectTo,
			&r.Version, &r.CreatedAt)
		if err != nil {
			return nil, errors.Wrapf(err, "scan columns")
		}
		r.ID = id
		records = append(records, r)
	}

	err = rows.Err()
	if err != nil {
		return nil, errors.Wrapf(err, "row iterator")
	}

	maxID := since
	if len(records) > 0 {
		maxID = id
	}

	log.Tracef("GetRates %s got %d records, nextID=%d", tableName, len(records), maxID)
	return &RateResult{Since: since, Limit: limit, Next: maxID, Records: records}, nil
}

// GetRates returns the collection from rates
func (postgres *PostgreSQL) ListRateProviders(tableName string, since int64, limit int) (*RateProviderResult, error) {
	log.Tracef("ListRateProviders %s since: %d limit: %d", tableName, since, limit)

	query := fmt.Sprintf(`
SELECT p.id, p.code, p.name, p.version
FROM %s p
WHERE id > $1 ORDER BY id LIMIT $2
`, tableName)

	var id int64
	rows, err := postgres.db.Query(query, since, limit)
	if err != nil {
		return nil, errors.Wrapf(err, "execute query")
	}
	defer rows.Close()

	records := make([]*ProviderRecord, 0)
	for rows.Next() {
		r := &ProviderRecord{}
		if err := rows.Scan(&id, &r.Code, &r.Name, &r.Version); err != nil {
			return nil, errors.Wrapf(err, "scan columns")
		}
		r.ID = id
		records = append(records, r)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrapf(err, "row iterator")
	}

	maxID := since
	if len(records) > 0 {
		maxID = id
	}

	log.Tracef("ListRateProviders %s got %d records, nextID=%d", tableName, len(records), maxID)
	return &RateProviderResult{Since: since, Limit: limit, Next: maxID, Records: records}, nil
}

func truncateNonASCII(s string) string {
	var buffer strings.Builder
	for _, r := range s {
		if r >= utf8.RuneSelf {
			break
		}
		buffer.WriteRune(r)
	}
	return buffer.String()
}
