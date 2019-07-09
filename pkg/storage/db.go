package storage

import (
	"database/sql"
	"fmt"

	"github.com/golang/protobuf/jsonpb"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gitlab.com/navyx/tools/pgexport/pkg/model"
)

// DBConfig holds the configuration used for instantiating a new Database.
type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	SSLMode  string `yaml:"sslmode"`
}

// Database holds the connection to database
type Database struct {
	config DBConfig
	db     *sql.DB
}

// NewDatabase returns the Database instance
func NewDatabase(config DBConfig) (*Database, error) {
	db, err := sql.Open("postgres", fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%d sslmode=%s",
		config.Username, config.Password, config.Database, config.Host, config.Port, config.SSLMode))

	if err != nil {
		return nil, errors.Wrapf(err, "open connection to database")
	}

	err = db.Ping()
	if err != nil {
		return nil, errors.Wrapf(err, "ping database")
	}
	return &Database{config, db}, nil
}

// Close disconnect the connection
func (db *Database) Close() error {
	if db.db == nil {
		return nil
	}

	return db.db.Close()
}

// 	GetTableMaxID return the maximum id of given table
func (db *Database) GetTableMaxID(tableName string) (int64, error) {
	log.Tracef("GetTableMaxID of %s", tableName)

	var maxID int64
	query := fmt.Sprintf(`SELECT id FROM %s ORDER BY id DESC LIMIT 1`, tableName)
	err := db.db.QueryRow(query).Scan(&maxID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, errors.Wrapf(err, "execute query")
	}

	log.Tracef("GetTableMaxID %s is %d", tableName, maxID)
	return maxID, nil
}

// SavePrebookings save prebookings to Database table
func (db *Database) SavePrebookings(tableName string, records []*model.PrebookingRecord) error {
	if len(records) == 0 {
		return nil
	}

	tx, err := db.db.Begin()
	if err != nil {
		return errors.Wrap(err, "start transaction")
	}
	defer tx.Rollback()

	createTempTable := `CREATE LOCAL TEMPORARY TABLE temp_prebookings ON COMMIT DROP AS SELECT * FROM prebookings LIMIT 0`
	if _, err := tx.Exec(createTempTable); err != nil {
		return errors.Wrap(err, "create temporary table")
	}

	insertTempTable := `INSERT INTO temp_prebookings (id, prebooking, submit_user_id, source_platform_id) VALUES ($1, $2, $3, $4)`
	stmt, err := tx.Prepare(insertTempTable)
	if err != nil {
		return errors.Wrap(err, "prepare statement")
	}
	defer stmt.Close()

	for _, record := range records {
		prebooking := record.Prebooking
		ma := jsonpb.Marshaler{}
		data, err := ma.MarshalToString(prebooking)
		if err != nil {
			return errors.Wrapf(err, "marshaling prebooking %s", prebooking.Id)
		}

		if _, err := stmt.Exec(prebooking.Id, data, record.SubmitUserID, record.SourcePlatformID); err != nil {
			return errors.Wrap(err, "insert temporary prebooking record")
		}
	}

	deleteDuplicated := fmt.Sprintf(`DELETE FROM %s p WHERE EXISTS (SELECT * FROM temp_prebookings t WHERE t.id = p.id)`, pq.QuoteIdentifier(tableName))
	if _, err := tx.Exec(deleteDuplicated); err != nil {
		return errors.Wrap(err, "delete duplicated record")
	}

	insertRecord := fmt.Sprintf(`INSERT INTO %s SELECT * FROM temp_prebookings`, pq.QuoteIdentifier(tableName))
	if _, err := tx.Exec(insertRecord); err != nil {
		return errors.Wrap(err, "insert prebooking record")
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "commit transaction")
	}

	return nil
}

// SaveBookingRecords save booking records to Database table
func (db *Database) SaveBookingRecords(tableName string, bookings []*model.BookingRecord) error {
	if len(bookings) == 0 {
		return nil
	}

	tx, err := db.db.Begin()
	if err != nil {
		return errors.Wrap(err, "start transaction")
	}
	defer tx.Rollback()

	query := fmt.Sprintf(`
INSERT INTO %s (
  id, region, application_number, version, is_draft, platform_id,
  created_by, updated_by, created_at, updated_at, data
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
`, tableName)

	stmt, err := tx.Prepare(query)
	if err != nil {
		return errors.Wrap(err, "prepare statement")
	}
	defer stmt.Close()

	for _, booking := range bookings {
		if _, err := stmt.Exec(booking.ID, booking.Region, booking.ApplicationNumber, booking.Version, booking.IsDraft,
			booking.PlatformID, booking.CreatedBy, booking.UpdatedBy, booking.CreatedAt, booking.UpdatedAt, booking.Data); err != nil {
			return errors.Wrap(err, "insert booking record")
		}
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "commit transaction")
	}

	return nil
}

// SaveBookingConfirmRecords save booking records to Database table
func (db *Database) SaveBookingConfirmRecords(tableName string, confirms []*model.BookingConfirmRecord) error {
	if len(confirms) == 0 {
		return nil
	}

	tx, err := db.db.Begin()
	if err != nil {
		return errors.Wrap(err, "start transaction")
	}
	defer tx.Rollback()

	query := fmt.Sprintf(`
INSERT INTO %s (id, region, application_number, version, created_at, data)
VALUES ($1, $2, $3, $4, $5, $6)
`, tableName)

	stmt, err := tx.Prepare(query)
	if err != nil {
		return errors.Wrap(err, "prepare statement")
	}
	defer stmt.Close()

	for _, confirm := range confirms {
		if _, err := stmt.Exec(confirm.ID, confirm.Region, confirm.ApplicationNumber, confirm.Version,
			confirm.CreatedAt, confirm.Data); err != nil {
			return errors.Wrap(err, "insert confirm record")
		}
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "commit transaction")
	}

	return nil
}

func (db *Database) GetReplicateProgress(taskName string, tableName string) (*ReplicateProgress, error) {
	if taskName == "" {
		return nil, errors.New("no task name")
	}
	if tableName == "" {
		return nil, errors.New("no table name")
	}
	log.Tracef("GetReplicateProgress for task %s on table %s", taskName, tableName)

	var progress, updateAt int64
	query := `SELECT progress, updated_at FROM replicate_progress WHERE task=$1 AND table_name=$2`
	err := db.db.QueryRow(query, taskName, tableName).Scan(&progress, &updateAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return &ReplicateProgress{}, nil
		}
		return nil, errors.Wrapf(err, "execute query")
	}

	log.Tracef("GetReplicateProgress for task %s on table %s: progress:%d updated_at:%d", taskName, tableName, progress, updateAt)
	return &ReplicateProgress{Progress: progress, UpdateAt: updateAt}, nil
}

func (db *Database) UpdateReplicateProgress(taskName string, tableName string, progress ReplicateProgress) error {
	if taskName == "" {
		return errors.New("no task name")
	}
	if tableName == "" {
		return errors.New("no table name")
	}
	log.Tracef("UpdateReplicateProgress for task %s on table %s:%+v", taskName, tableName, progress)

	query := `
INSERT INTO replicate_progress(task, table_name, progress, updated_at)
VALUES ($1, $2, $3, $4)
ON CONFLICT (task, table_name)
DO UPDATE SET progress=$3, updated_at=$4
`
	_, err := db.db.Exec(query, taskName, tableName, progress.Progress, progress.UpdateAt)
	if err != nil {
		return errors.Wrapf(err, "execute query")
	}

	return nil
}

func (db *Database) SaveRouteScheduleRecords(tableName string, schedules []*model.RouteScheduleRecord) error {
	if len(schedules) == 0 {
		return nil
	}

	tx, err := db.db.Begin()
	if err != nil {
		return errors.Wrap(err, "start transaction")
	}
	defer tx.Rollback()

	selectQuery := fmt.Sprintf(`
SELECT id FROM %s WHERE carrier=$1 AND vessel=$2 AND voyage=$3
`, tableName)

	selectStmt, err := tx.Prepare(selectQuery)
	if err != nil {
		return errors.Wrap(err, "prepare selection statement")
	}
	defer selectStmt.Close()

	insertQuery := fmt.Sprintf(`
INSERT INTO %s (carrier, vessel, voyage, lloyd_code, data) VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (upper(carrier), upper(vessel), upper(voyage))
DO UPDATE SET lloyd_code=EXCLUDED.lloyd_code, data=EXCLUDED.data
RETURNING id
`, tableName)

	insertStmt, err := tx.Prepare(insertQuery)
	if err != nil {
		return errors.Wrap(err, "prepare insertion statement")
	}
	defer insertStmt.Close()

	updateQuery := fmt.Sprintf(`
UPDATE %s SET lloyd_code=$2, data=$3 WHERE id=$1
`, tableName)
	updateStmt, err := tx.Prepare(updateQuery)
	if err != nil {
		return errors.Wrap(err, "prepare updating statement")
	}
	defer updateStmt.Close()

	for _, schedule := range schedules {
		key := fmt.Sprintf("%s/%s/%s", schedule.Carrier, schedule.Vessel, schedule.Voyage)
		var id int64

		if err = selectStmt.QueryRow(schedule.Carrier, schedule.Vessel, schedule.Voyage).Scan(&id); err != nil {
			if err != sql.ErrNoRows {
				return errors.Wrapf(err, "select schedule (%s)", key)
			}
			err = insertStmt.QueryRow(schedule.Carrier, schedule.Vessel, schedule.Voyage, schedule.LloydCode, schedule.Data).Scan(&id)
			if err != nil {
				return errors.Wrapf(err, "insert route schedule (%s)", key)
			}
			log.Tracef("route schedule %s inserted to row %d", key, id)
		} else {
			if _, err := updateStmt.Exec(id, schedule.LloydCode, schedule.Data); err != nil {
				return errors.Wrap(err, "update route schedule")
			}
			log.Tracef("route schedule %s updated to row %d", key, id)
		}
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "commit transaction")
	}

	return nil
}

func (db *Database) SaveBookingSummaryRecords(tableName string, records []*model.BookingSummaryRecord) error {
	if len(records) == 0 {
		return nil
	}

	tx, err := db.db.Begin()
	if err != nil {
		return errors.Wrap(err, "start transaction")
	}
	defer tx.Rollback()

	createTempTable := fmt.Sprintf(`CREATE LOCAL TEMPORARY TABLE temp_booking_summary ON COMMIT DROP AS SELECT * FROM %s LIMIT 0`, pq.QuoteIdentifier(tableName))
	if _, err := tx.Exec(createTempTable); err != nil {
		return errors.Wrap(err, "create temporary table")
	}

	insertTempTable := `
INSERT INTO temp_booking_summary (
    id, region, platform_id, application_number, carrier, status, is_deleted,
    created_by, updated_by, created_at, updated_at, data
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
`
	stmt, err := tx.Prepare(insertTempTable)
	if err != nil {
		return errors.Wrap(err, "prepare statement")
	}
	defer stmt.Close()

	for _, record := range records {
		region := sql.NullString{}
		carrier := sql.NullString{}
		updatedBy := sql.NullInt64{}
		updatedAt := sql.NullInt64{}
		if record.Region != "" {
			region = sql.NullString{String: record.Region, Valid: true}
		}
		if record.Carrier != "" {
			carrier = sql.NullString{String: record.Carrier, Valid: true}
		}
		if record.UpdatedBy > 0 {
			updatedBy = sql.NullInt64{Int64: record.UpdatedBy, Valid: true}
		}
		if record.UpdatedAt > 0 {
			updatedAt = sql.NullInt64{Int64: record.UpdatedAt, Valid: true}
		}

		if _, err := stmt.Exec(record.ID, region, record.PlatformID, record.AppNumber, carrier, record.Status, record.IsDeleted,
			record.CreatedBy, updatedBy, record.CreatedAt, updatedAt, record.Data); err != nil {
			return errors.Wrap(err, "insert temporary booking summary record")
		}
	}

	deleteDuplicated := fmt.Sprintf(`DELETE FROM %s p WHERE EXISTS (SELECT * FROM temp_booking_summary t WHERE t.id = p.id)`, pq.QuoteIdentifier(tableName))
	if _, err := tx.Exec(deleteDuplicated); err != nil {
		return errors.Wrap(err, "delete duplicated record")
	}

	insertRecord := fmt.Sprintf(`INSERT INTO %s SELECT * FROM temp_booking_summary`, pq.QuoteIdentifier(tableName))
	if _, err := tx.Exec(insertRecord); err != nil {
		return errors.Wrap(err, "insert booking summary record")
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "commit transaction")
	}

	return nil
}

func (db *Database) SaveBookingConfirmSummaryRecords(tableName string, records []*model.BookingConfirmSummaryRecord) error {
	if len(records) == 0 {
		return nil
	}

	tx, err := db.db.Begin()
	if err != nil {
		return errors.Wrap(err, "start transaction")
	}
	defer tx.Rollback()

	createTempTable := fmt.Sprintf(`CREATE LOCAL TEMPORARY TABLE temp_booking_confirm_summary ON COMMIT DROP AS SELECT * FROM %s LIMIT 0`, pq.QuoteIdentifier(tableName))
	if _, err := tx.Exec(createTempTable); err != nil {
		return errors.Wrap(err, "create temporary table")
	}

	insertTempTable := `
INSERT INTO temp_booking_confirm_summary (id, region, application_number, booking_number, status, ts, data)
VALUES ($1, $2, $3, $4, $5, $6, $7)
`
	stmt, err := tx.Prepare(insertTempTable)
	if err != nil {
		return errors.Wrap(err, "prepare statement")
	}
	defer stmt.Close()

	for _, record := range records {
		region := sql.NullString{}
		timestamp := sql.NullInt64{}
		if record.Region != "" {
			region = sql.NullString{String: record.Region, Valid: true}
		}
		if record.Timestamp > 0 {
			timestamp = sql.NullInt64{Int64: record.Timestamp, Valid: true}
		}

		if _, err := stmt.Exec(record.ID, region, record.AppNumber, record.BookingNumber,
			record.Status, timestamp, record.Data); err != nil {
			return errors.Wrap(err, "insert temporary booking confirm summary record")
		}
	}

	deleteDuplicated := fmt.Sprintf(`DELETE FROM %s p WHERE EXISTS (SELECT * FROM temp_booking_confirm_summary t WHERE t.id = p.id)`, pq.QuoteIdentifier(tableName))
	if _, err := tx.Exec(deleteDuplicated); err != nil {
		return errors.Wrap(err, "delete duplicated record")
	}

	insertRecord := fmt.Sprintf(`INSERT INTO %s SELECT * FROM temp_booking_confirm_summary`, pq.QuoteIdentifier(tableName))
	if _, err := tx.Exec(insertRecord); err != nil {
		return errors.Wrap(err, "insert booking confirm summary record")
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "commit transaction")
	}

	return nil
}

func (db *Database) SaveDocumentSummaryRecords(tableName string, records []*model.BookingSummaryRecord) error {
	if len(records) == 0 {
		return nil
	}

	tx, err := db.db.Begin()
	if err != nil {
		return errors.Wrap(err, "start transaction")
	}
	defer tx.Rollback()

	createTempTable := fmt.Sprintf(`CREATE LOCAL TEMPORARY TABLE temp_document_summary ON COMMIT DROP AS SELECT * FROM %s LIMIT 0`, pq.QuoteIdentifier(tableName))
	if _, err := tx.Exec(createTempTable); err != nil {
		return errors.Wrap(err, "create temporary table")
	}

	insertTempTable := `
INSERT INTO temp_document_summary (
    id, region, platform_id, application_number, carrier, status, is_deleted,
    created_by, updated_by, created_at, updated_at, data, confirm
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
`
	stmt, err := tx.Prepare(insertTempTable)
	if err != nil {
		return errors.Wrap(err, "prepare statement")
	}
	defer stmt.Close()

	for _, record := range records {
		region := sql.NullString{}
		carrier := sql.NullString{}
		updatedBy := sql.NullInt64{}
		updatedAt := sql.NullInt64{}
		if record.Region != "" {
			region = sql.NullString{String: record.Region, Valid: true}
		}
		if record.Carrier != "" {
			carrier = sql.NullString{String: record.Carrier, Valid: true}
		}
		if record.UpdatedBy > 0 {
			updatedBy = sql.NullInt64{Int64: record.UpdatedBy, Valid: true}
		}
		if record.UpdatedAt > 0 {
			updatedAt = sql.NullInt64{Int64: record.UpdatedAt, Valid: true}
		}

		if _, err := stmt.Exec(record.ID, region, record.PlatformID, record.AppNumber, carrier, record.Status, record.IsDeleted,
			record.CreatedBy, updatedBy, record.CreatedAt, updatedAt, record.Data, record.ConfirmData); err != nil {
			return errors.Wrap(err, "insert temporary confirm summary record")
		}
	}

	deleteDuplicated := fmt.Sprintf(`DELETE FROM %s p WHERE EXISTS (SELECT * FROM temp_document_summary t WHERE t.id = p.id)`, pq.QuoteIdentifier(tableName))
	if _, err := tx.Exec(deleteDuplicated); err != nil {
		return errors.Wrap(err, "delete duplicated record")
	}

	insertRecord := fmt.Sprintf(`INSERT INTO %s SELECT * FROM temp_document_summary`, pq.QuoteIdentifier(tableName))
	if _, err := tx.Exec(insertRecord); err != nil {
		return errors.Wrap(err, "insert confirm summary record")
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "commit transaction")
	}

	return nil
}

func (db *Database) SaveRateRecords(tableName string, records []*model.RateRecord) error {
	if len(records) == 0 {
		return nil
	}

	tx, err := db.db.Begin()
	if err != nil {
		return errors.Wrap(err, "start transaction")
	}
	defer tx.Rollback()

	query := fmt.Sprintf(`
INSERT INTO %s (
  origin, port_of_loading, port_of_discharge, destination, carrier_scac, container_type,
  provider_id, price, service_fee, service_code, contract_number, service_mode,
  effective_date, expiry_date, service_fee_effective_date, service_fee_expiry_date,
  promotion_discount, promotion_effective_date, promotion_expiry_date,
  commodity, demurrage, detention, outport, remarks, included, subject_to, version, created_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
  $11, $12, $13, $14, $15, $16, $17, $18, $19, $20,
  $21, $22, $23, $24, $25, $26, $27, $28
)
`, tableName)
	stmt, err := tx.Prepare(query)
	if err != nil {
		return errors.Wrap(err, "prepare statement")
	}
	defer stmt.Close()

	for _, record := range records {
		if _, err := stmt.Exec(record.PlaceOfReceopt, record.PortOfLoading, record.PortOfDischarge, record.PlaceOfDelivery,
			record.CarrierSCAC, record.ContainerType, record.ProviderID, record.Price, record.ServiceFee,
			record.ServiceCode, record.ContractNumber, record.ServiceMode,
			record.EffectiveDate, record.ExpiryDate, record.ServiceFeeEffectiveDate, record.ServiceFeeExpiryDate,
			record.PromotionDiscount, record.PromotionEffectiveDate, record.PromotionExpiryDate,
			record.Commodity, record.Demurrage, record.Detention, record.Outport, record.Remarks,
			record.Included, record.SubjectTo, record.Version, record.CreatedAt); err != nil {
			return errors.Wrapf(err, "insert rate record %d", record.ID)
		}
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "commit transaction")
	}

	return nil
}

// SaveRateProviders save rate providers to Database table
func (db *Database) SaveRateProviders(tableName string, records []*model.ProviderRecord) error {
	if len(records) == 0 {
		return nil
	}

	tx, err := db.db.Begin()
	if err != nil {
		return errors.Wrap(err, "start transaction")
	}
	defer tx.Rollback()

	query := fmt.Sprintf(`
INSERT INTO %s (id, code, name, version)
VALUES ($1, $2, $3, $4)
ON CONFLICT (id) DO UPDATE
SET code=EXCLUDED.code,
    name=EXCLUDED.name,
    version=EXCLUDED.version
`, tableName)
	stmt, err := tx.Prepare(query)
	if err != nil {
		return errors.Wrap(err, "prepare statement")
	}
	defer stmt.Close()

	for _, record := range records {
		if _, err := stmt.Exec(record.ID, record.Code, record.Name, record.Version); err != nil {
			return errors.Wrapf(err, "insert provider record %d", record.ID)
		}
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "commit transaction")
	}

	return nil
}
