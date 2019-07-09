package task

import (
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gitlab.com/navyx/tools/pgexport/pkg/model"
	"gitlab.com/navyx/tools/pgexport/pkg/storage"
)

type RouteScheduleTask struct {
	storage storage.Storage
	config  Config
}

func NewRouteScheduleTask(storage storage.Storage, config Config) Task {
	return &RouteScheduleTask{storage: storage, config: config}
}

func (task *RouteScheduleTask) Execute() error {
	if !task.config.Enable {
		log.Printf("RouteScheduleTask is not enabled, skip it.")
		return nil
	}
	db, err := model.NewPostgreSQL(task.config.Source)
	if err != nil {
		return errors.Wrapf(err, "connect to database")
	}
	defer db.Close()

	for _, table := range task.config.Tables {
		lastProgress, err := task.storage.GetReplicateProgress(ROUTE_SCHEDULE_TASK_NAME, table.From)
		if err != nil {
			return errors.Wrapf(err, "get last replicated progress")
		}
		for {
			since := lastProgress.Progress
			result, err := db.GetRouteSchedules(table.From, since, task.config.BatchSize)
			if err != nil {
				return errors.Wrapf(err, "list route schedules")
			}
			if len(result.Records) == 0 {
				break
			}

			if err := task.storage.SaveRouteScheduleRecords(table.To, result.Records); err != nil {
				return errors.Wrapf(err, "save route schedules")
			}

			lastProgress.Progress = result.Next
			lastProgress.UpdateAt = time.Now().Unix()
			err = task.storage.UpdateReplicateProgress(ROUTE_SCHEDULE_TASK_NAME, table.From, *lastProgress)
			if err != nil {
				return errors.Wrapf(err, "update replicated progress")
			}
		}
	}

	return nil
}
