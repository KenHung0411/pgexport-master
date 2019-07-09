package task

import (
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gitlab.com/navyx/tools/pgexport/pkg/model"
	"gitlab.com/navyx/tools/pgexport/pkg/storage"
)

type PrebookingTask struct {
	storage storage.Storage
	config  Config
}

func NewPrebookingTask(storage storage.Storage, config Config) Task {
	return &PrebookingTask{storage: storage, config: config}
}

func (task *PrebookingTask) Execute() error {
	if !task.config.Enable {
		log.Printf("PrebookingTask is not enabled, skip it.")
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
			result, err := db.ListPrebookings(table.From, since, task.config.BatchSize)
			if err != nil {
				log.Warnf("list prebookings from %s failed. %s", table.From, err)
				lastError = err
				break
			}
			if len(result.Prebookings) == 0 {
				break
			}
			if err := task.storage.SavePrebookings(table.To, result.Prebookings); err != nil {
				log.Warnf("save prebookings to %s failed. %s", table.To, err)
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
			err = task.storage.UpdateReplicateProgress(PREBOOKING_TASK_NAME, table.To, progress)
			if err != nil {
				log.Warnf("update replicated progress %s:%s failed. %s", PREBOOKING_TASK_NAME, table.To, err)
			}
		}
	}
	return nil
}
