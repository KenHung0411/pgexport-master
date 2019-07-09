package main

import (
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gitlab.com/navyx/tools/pgexport/pkg/helper/formatter"
	"gitlab.com/navyx/tools/pgexport/pkg/storage"
	"gitlab.com/navyx/tools/pgexport/pkg/task"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v2"
)

const appName = "pgexport"

var appConfig = AppConfig{
	Prebooking:            task.Config{BatchSize: 100},
	Booking:               task.Config{BatchSize: 100},
	BookingConfirm:        task.Config{BatchSize: 100},
	BookingSummary:        task.Config{BatchSize: 100},
	BookingConfirmSummary: task.Config{BatchSize: 100},
	DocumentSummary:       task.Config{BatchSize: 100},
	Rate:                  task.Config{BatchSize: 100},
	RouteSchedule:         task.Config{BatchSize: 100},
	RateProvider:          task.Config{BatchSize: 100},
}

func init() {
	formatter := formatter.Formatter{}
	formatter.TimestampFormat = "2006-01-02 15:04:05"
	formatter.HideKeys = true
	formatter.ShowFullLevel = true
	formatter.NoColors = true
	formatter.NoFieldsColors = true

	log.SetFormatter(&formatter)
	log.SetOutput(os.Stdout)
	log.SetLevel(log.TraceLevel)
}

func main() {
	app := kingpin.New(appName, "PostgreSQL Data Exporter")
	app.Version("1.0.0")
	app.HelpFlag.Short('h')
	config := app.Flag("config", "Config File").Short('c').Default("config.yml").String()
	verbose := app.Flag("verbose", "Verbose mode.").Short('v').Bool()

	kingpin.MustParse(app.Parse(os.Args[1:]))

	if confContent, err := ioutil.ReadFile(*config); err != nil {
		log.Fatalf("Reading Config: %v", err)
	} else {
		confContent = []byte(os.ExpandEnv(string(confContent)))
		err = yaml.Unmarshal([]byte(confContent), &appConfig)
		if err != nil {
			log.Fatalf("Parse Config: %v", err)
		}
		if *verbose {
			log.Printf("Config: %+v", appConfig)
		}
	}

	if err := syncData(); err != nil {
		log.Fatalf("%v", err)
	}
}

func syncData() error {
	s, err := storage.NewDatabase(appConfig.Database)
	if err != nil {
		return err
	}
	defer s.Close()

	var t task.Task

	t = task.NewPrebookingTask(s, appConfig.Prebooking)
	if err := t.Execute(); err != nil {
		return errors.Wrapf(err, "execute prebooking task")
	}

	t = task.NewBookingTask(s, appConfig.Booking)
	if err := t.Execute(); err != nil {
		return errors.Wrapf(err, "execute booking task")
	}

	t = task.NewBookingConfirmTask(s, appConfig.BookingConfirm)
	if err := t.Execute(); err != nil {
		return errors.Wrapf(err, "execute booking confirm task")
	}

	t = task.NewBookingSummaryTask(s, appConfig.BookingSummary)
	if err := t.Execute(); err != nil {
		return errors.Wrapf(err, "execute booking summary task")
	}

	t = task.NewBookingConfirmSummaryTask(s, appConfig.BookingConfirmSummary)
	if err := t.Execute(); err != nil {
		return errors.Wrapf(err, "execute booking confirm summary task")
	}

	t = task.NewDocumentSummaryTask(s, appConfig.DocumentSummary)
	if err := t.Execute(); err != nil {
		return errors.Wrapf(err, "execute document summary task")
	}

	t = task.NewRateTask(s, appConfig.Rate)
	if err := t.Execute(); err != nil {
		return errors.Wrapf(err, "execute booking task")
	}

	t = task.NewRouteScheduleTask(s, appConfig.RouteSchedule)
	if err := t.Execute(); err != nil {
		return errors.Wrapf(err, "execute route schedule task")
	}

	t = task.NewRateProviderTask(s, appConfig.RateProvider)
	if err := t.Execute(); err != nil {
		return errors.Wrapf(err, "execute rate provider task")
	}
	return nil
}
