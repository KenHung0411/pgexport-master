package task

import (
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gitlab.com/navyx/tools/pgexport/pkg/model"
	"gitlab.com/navyx/tools/pgexport/pkg/storage"
)

type RateTask struct {
	storage storage.Storage
	config  Config
}

func NewRateTask(storage storage.Storage, config Config) Task {
	return &RateTask{storage: storage, config: config}
}

func (task *RateTask) Execute() error {
	if !task.config.Enable {
		log.Printf("RateTask is not enabled, skip it.")
		return nil
	}
	db, err := model.NewPostgreSQL(task.config.Source)
	if err != nil {
		return errors.Wrapf(err, "connect to database")
	}
	defer db.Close()

	for _, table := range task.config.Tables {
		lastProgress, err := task.storage.GetReplicateProgress(RATE_TASK_NAME, table.From)
		if err != nil {
			return errors.Wrapf(err, "get last replicated progress")
		}
		for {
			since := lastProgress.Progress
			result, err := db.GetRates(table.From, since, task.config.BatchSize)
			if err != nil {
				return errors.Wrapf(err, "list rates")
			}
			if len(result.Records) == 0 {
				break
			}

			if err := task.storage.SaveRateRecords(table.To, result.Records); err != nil {
				return errors.Wrapf(err, "save rates")
			}

			lastProgress.Progress = result.Next
			lastProgress.UpdateAt = time.Now().Unix()
			err = task.storage.UpdateReplicateProgress(RATE_TASK_NAME, table.From, *lastProgress)
			if err != nil {
				return errors.Wrapf(err, "update replicated progress")
			}
		}
	}

	return nil
}
