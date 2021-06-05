package goit

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"

	"github.com/tungnt/goit/database"
	"github.com/tungnt/goit/must"
)

type ITsql struct {
	it
	transactions map[string]*sql.Tx
}

func (i *ITsql) BeforeTest(suiteName, testName string) {
	db, err := database.NewProvider().DB()
	must.NotFail(err)
	err = i.cleanUpDB(db)
	must.NotFail(err)
	err = i.setConnectionIntoMap(testName, db)
	must.NotFail(err)
	err = i.initiateTransaction(testName)
	must.NotFail(err)
}

func (i *ITsql) AfterTest(suiteName, testName string) {
	tx, err := i.getTransactionFromMap(testName)
	must.NotFail(err)
	err = tx.Rollback()
	must.NotFail(err)
	i.deleteTransationFromMap(testName)
	i.it.AfterTest(suiteName, testName)
}

func (i *ITsql) cleanUpDB(db *sql.DB) error {
	rootDir := i.rootDirectory()
	truncateStmt, err := ioutil.ReadFile(rootDir + "/" + i.config.DatabaseTruncateFile)
	if err != nil {
		return err
	}
	return database.NewProvider().CleanUpDB(db, string(truncateStmt))
}

func (i *ITsql) initiateTransaction(testName string) error {
	tx, err := i.startTransaction(testName)
	if err != nil {
		return err
	}
	return i.setTransactionIntoMap(testName, tx)
}

func (i *ITsql) startTransaction(testName string) (*sql.Tx, error) {
	db, err := i.getConnection(testName)
	if err != nil {
		return nil, err
	}
	options := &sql.TxOptions{Isolation: sql.LevelSerializable}
	tx, err := db.BeginTx(context.Background(), options)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (i *ITsql) setTransactionIntoMap(testName string, tx *sql.Tx) error {
	i.initTransactionMapIfNeed()
	_, ok := i.transactions[testName]
	if ok {
		return fmt.Errorf("Transaction map key conflicts (%s)", testName)
	}
	i.transactions[testName] = tx
	return nil
}

func (i *ITsql) initTransactionMapIfNeed() {
	mutex.Lock()
	defer mutex.Unlock()
	if i.transactions != nil {
		return
	}

	i.transactions = make(map[string]*sql.Tx)
}

func (i *ITsql) getTransactionFromMap(testName string) (*sql.Tx, error) {
	tx, ok := i.transactions[testName]
	if !ok {
		return nil, fmt.Errorf(
			`Transaction of test "%s" is not set yet. Is "%s" your test function name?`,
			testName,
			testName,
		)
	}

	return tx, nil
}

func (i *ITsql) deleteTransationFromMap(testName string) {
	delete(i.transactions, testName)
}
