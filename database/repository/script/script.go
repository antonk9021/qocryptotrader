package script

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
	"github.com/antonk9021/qocryptotrader/database"
	modelPSQL "github.com/antonk9021/qocryptotrader/database/models/postgres"
	modelSQLite "github.com/antonk9021/qocryptotrader/database/models/sqlite3"
	"github.com/antonk9021/qocryptotrader/database/repository"
	"github.com/antonk9021/qocryptotrader/log"
	"github.com/thrasher-corp/sqlboiler/boil"
	"github.com/volatiletech/null"
)

// Event inserts a new script event into database with execution details (script name time status hash of script)
func Event(id, name, path string, data null.Bytes, executionType, status string, time time.Time) {
	if database.DB.SQL == nil {
		return
	}

	ctx := context.TODO()
	ctx = boil.SkipTimestamps(ctx)
	tx, err := database.DB.SQL.BeginTx(ctx, nil)
	if err != nil {
		log.Errorf(log.DatabaseMgr, "Event transaction begin failed: %v", err)
		return
	}

	if repository.GetSQLDialect() == database.DBSQLite3 {
		query := modelSQLite.ScriptWhere.ScriptID.EQ(id)
		f, errQry := modelSQLite.Scripts(query).Exists(ctx, tx)
		if errQry != nil {
			log.Errorf(log.DatabaseMgr, "Query failed: %v", errQry)
			err = tx.Rollback()
			if err != nil {
				log.Errorf(log.DatabaseMgr, "Event Transaction rollback failed: %v", err)
			}
			return
		}
		var tempEvent = modelSQLite.Script{}
		if !f {
			newUUID, errUUID := uuid.NewV4()
			if errUUID != nil {
				log.Errorf(log.DatabaseMgr, "Failed to generate UUID: %v", errUUID)
				_ = tx.Rollback()
				return
			}

			tempEvent.ID = newUUID.String()
			tempEvent.ScriptID = id
			tempEvent.ScriptName = name
			tempEvent.ScriptPath = path
			tempEvent.ScriptData = data
			err = tempEvent.Insert(ctx, tx, boil.Infer())
			if err != nil {
				log.Errorf(log.DatabaseMgr, "Event insert failed: %v", err)
				err = tx.Rollback()
				if err != nil {
					log.Errorf(log.DatabaseMgr, "Event Transaction rollback failed: %v", err)
				}
				return
			}
		} else {
			tempEvent.ID = id
		}

		tempScriptExecution := &modelSQLite.ScriptExecution{
			ScriptID:        id,
			ExecutionTime:   time.UTC().String(),
			ExecutionStatus: status,
			ExecutionType:   executionType,
		}
		err = tempEvent.AddScriptExecutions(ctx, tx, true, tempScriptExecution)
		if err != nil {
			log.Errorf(log.DatabaseMgr, "Event insert failed: %v", err)
			err = tx.Rollback()
			if err != nil {
				log.Errorf(log.DatabaseMgr, "Event Transaction rollback failed: %v", err)
			}
			return
		}
	} else {
		var tempEvent = modelPSQL.Script{
			ScriptID:   id,
			ScriptName: name,
			ScriptPath: path,
			ScriptData: data,
		}
		err = tempEvent.Upsert(ctx, tx, true, []string{"script_id"}, boil.Whitelist("last_executed_at"), boil.Infer())
		if err != nil {
			log.Errorf(log.DatabaseMgr, "Event insert failed: %v", err)
			err = tx.Rollback()
			if err != nil {
				log.Errorf(log.DatabaseMgr, "Event Transaction rollback failed: %v", err)
			}
			return
		}

		tempScriptExecution := &modelPSQL.ScriptExecution{
			ExecutionTime:   time.UTC(),
			ExecutionStatus: status,
			ExecutionType:   executionType,
		}

		err = tempEvent.AddScriptExecutions(ctx, tx, true, tempScriptExecution)
		if err != nil {
			log.Errorf(log.DatabaseMgr, "Event insert failed: %v", err)
			err = tx.Rollback()
			if err != nil {
				log.Errorf(log.DatabaseMgr, "Event Transaction rollback failed: %v", err)
			}
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Errorf(log.DatabaseMgr, "Event Transaction commit failed: %v", err)
	}
}
