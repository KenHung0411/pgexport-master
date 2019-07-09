package task

import (
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gitlab.com/navyx/tools/pgexport/pkg/model"
	"gitlab.com/navyx/tools/pgexport/pkg/storage"
)

type DocumentSummaryTask struct {
	storage storage.Storage
	config  Config
}

func NewDocumentSummaryTask(storage storage.Storage, config Config) Task {
	return &DocumentSummaryTask{storage: storage, config: config}
}

func (task *DocumentSummaryTask) Execute() error {
	if !task.config.Enable {
		log.Printf("DocumentSummaryTask is not enabled, skip it.")
		return nil
	}
	db, err := model.NewPostgreSQL(task.config.Source)
	if err != nil {
		return errors.Wrapf(err, "connect to database")
	}
	defer db.Close()

	for _, table := range task.config.Tables {
		since := int64(0)
		var lastError error
		for {
			result, err := db.GetDocumentSummaryRecords(table.From, since, task.config.BatchSize)
			if err != nil {
				log.Warnf("get records from %s failed. %s", table.From, err)
				lastError = err
				break
			}
			if len(result.Records) == 0 {
				break
			}
			if err := task.storage.SaveDocumentSummaryRecords(table.To, result.Records); err != nil {
				log.Warnf("save record to %s failed. %s", table.To, err)
				lastError = err
				break
			}
			since = result.Next
		}
		if lastError == nil {
			progress := storage.ReplicateProgress{
				Progress: since,
				UpdateAt: time.Now().Unix(),
			}
			err = task.storage.UpdateReplicateProgress(DOCUMENT_SUMMARY_TASK_NAME, table.To, progress)
			if err != nil {
				log.Warnf("update replicated progress %s:%s failed. %s", DOCUMENT_SUMMARY_TASK_NAME, table.To, err)
			}
		}
	}

	return nil
}
