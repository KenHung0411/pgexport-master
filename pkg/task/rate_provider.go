package task

import (
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gitlab.com/navyx/tools/pgexport/pkg/model"
	"gitlab.com/navyx/tools/pgexport/pkg/storage"
)

type RateProviderTask struct {
	storage storage.Storage
	config  Config
}

func NewRateProviderTask(storage storage.Storage, config Config) Task {
	return &RateProviderTask{storage: storage, config: config}
}

func (task *RateProviderTask) Execute() error {
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
			result, err := db.ListRateProviders(table.From, since, task.config.BatchSize)
			if err != nil {
				log.Warnf("list rate providers from %s failed. %s", table.From, err)
				lastError = err
				break
			}
			if len(result.Records) == 0 {
				break
			}
			if err := task.storage.SaveRateProviders(table.To, result.Records); err != nil {
				log.Warnf("save rate providers to %s failed. %s", table.To, err)
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
			err = task.storage.UpdateReplicateProgress(RATE_PROVIDER_TASK_NAME, table.To, progress)
			if err != nil {
				log.Warnf("update replicated progress %s:%s failed. %s", RATE_PROVIDER_TASK_NAME, table.To, err)
			}
		}
	}
	return nil
}
