package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gitlab.com/navyx/tools/pgexport/pkg/helper/formatter"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v2"
)

const appName = "migrator"

var appConfig = AppConfig{}

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
	app := kingpin.New(appName, "Database Migration Toolkit")
	app.Version("1.0.0")
	app.HelpFlag.Short('h')
	config := app.Flag("config", "Config File").Short('c').Default("database.yml").String()
	verbose := app.Flag("verbose", "Verbose mode.").Short('v').Bool()

	cmdMigrate := app.Command("migrate", "Run Database Migration")
	migratePath := cmdMigrate.Flag("path", "Path to the Migrations Folder").Short('p').Default("migrations").String()
	cmdRollback := app.Command("rollback", "Run Database Rollback")
	rollbackPath := cmdRollback.Flag("path", "Path to the Migrations Folder").Short('p').Default("migrations").String()
	rollbackStep := cmdRollback.Flag("step", "Rollback N steps").Default("1").Int()

	command := kingpin.MustParse(app.Parse(os.Args[1:]))

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

	var err error
	switch command {
	case cmdMigrate.FullCommand():
		err = dbMigrate(*migratePath)

	case cmdRollback.FullCommand():
		err = dbRollback(*rollbackPath, *rollbackStep)
	}

	if err != nil {
		log.Fatalf("%v", err)
	}
}

func dbMigrate(migratePath string) error {
	conn, err := getConneciton()
	if err != nil {
		return errors.Wrapf(err, "get database connection")
	}
	defer conn.Close()
	migrator, err := pop.NewFileMigrator(migratePath, conn)
	if err != nil {
		return errors.Wrapf(err, "create database migrator")
	}
	return migrator.Up()
}

func dbRollback(rollbackPath string, step int) error {
	conn, err := getConneciton()
	if err != nil {
		return errors.Wrapf(err, "get database connection")
	}
	defer conn.Close()
	migrator, err := pop.NewFileMigrator(rollbackPath, conn)
	if err != nil {
		return errors.Wrapf(err, "create database migrator")
	}
	return migrator.Down(step)
}

func getConneciton() (*pop.Connection, error) {
	details := &pop.ConnectionDetails{
		Dialect:    "postgres",
		Database:   appConfig.Database.Database,
		Host:       appConfig.Database.Host,
		Port:       fmt.Sprintf("%d", appConfig.Database.Port),
		User:       appConfig.Database.Username,
		Password:   appConfig.Database.Password,
		RawOptions: fmt.Sprintf("sslmode=%s", appConfig.Database.SSLMode),
	}
	return pop.NewConnection(details)
}
