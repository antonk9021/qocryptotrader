package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/antonk9021/qocryptotrader/common"
	"github.com/antonk9021/qocryptotrader/config"
	"github.com/antonk9021/qocryptotrader/core"
	"github.com/antonk9021/qocryptotrader/database"
	dbPSQL "github.com/antonk9021/qocryptotrader/database/drivers/postgres"
	dbsqlite3 "github.com/antonk9021/qocryptotrader/database/drivers/sqlite3"
	"github.com/antonk9021/qocryptotrader/database/repository"
	"github.com/thrasher-corp/goose"
)

var (
	dbConn         *database.Instance
	configFile     string
	defaultDataDir string
	migrationDir   string
	command        string
	args           string
)

func openDBConnection(cfg *database.Config) (err error) {
	if cfg.Driver == database.DBPostgreSQL {
		dbConn, err = dbPSQL.Connect(cfg)
		if err != nil {
			return fmt.Errorf("database failed to connect: %v, some features that utilise a database will be unavailable", err)
		}
		return nil
	} else if cfg.Driver == database.DBSQLite || cfg.Driver == database.DBSQLite3 {
		dbConn, err = dbsqlite3.Connect(cfg.Database)
		if err != nil {
			return fmt.Errorf("database failed to connect: %v, some features that utilise a database will be unavailable", err)
		}
		return nil
	}
	return errors.New("no connection established")
}

func main() {
	fmt.Println("GoCryptoTrader database migration tool")
	fmt.Println(core.Copyright)
	fmt.Println()

	flag.StringVar(&command, "command", "", "command to run status|up|up-by-one|up-to|down|create")
	flag.StringVar(&args, "args", "", "arguments to pass to goose")
	flag.StringVar(&configFile, "config", config.DefaultFilePath(), "config file to load")
	flag.StringVar(&defaultDataDir, "datadir", common.GetDefaultDataDir(runtime.GOOS), "default data directory for GoCryptoTrader files")
	flag.StringVar(&migrationDir, "migrationdir", database.MigrationDir, "override migration folder")

	flag.Parse()

	var conf config.Config
	err := conf.LoadConfig(configFile, true)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if !conf.Database.Enabled {
		fmt.Println("Database support is disabled")
		os.Exit(1)
	}

	err = openDBConnection(&conf.Database)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	drv := repository.GetSQLDialect()
	if drv == database.DBSQLite || drv == database.DBSQLite3 {
		fmt.Printf("Database file: %s\n", conf.Database.Database)
	} else {
		fmt.Printf("Connected to: %s\n", conf.Database.Host)
	}

	if command == "" {
		_ = goose.Run("status", dbConn.SQL, drv, migrationDir, "")
		fmt.Println()
		flag.Usage()
		return
	}

	if err = goose.Run(command, dbConn.SQL, drv, migrationDir, args); err != nil {
		fmt.Println(err)
	}
}
