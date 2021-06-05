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

func (i *ITsql) SetupTest() {
	db, err := database.NewProvider().DB()
	must.NotFail(err)
	err = i.cleanUpDB(db)
	must.NotFail(err)
	err = i.setConnectionIntoMap(db)
	must.NotFail(err)
	err = i.initiateTransaction()
	must.NotFail(err)
}

func (i *ITsql) cleanUpDB(db *sql.DB) error {
	rootDir := i.rootDirectory()
	truncateStmt, err := ioutil.ReadFile(rootDir + "/" + i.config.DatabaseTruncateFile)
	if err != nil {
		return err
	}
	return database.NewProvider().CleanUpDB(db, string(truncateStmt))
}

func (i *ITsql) initiateTransaction() error {
	tx, err := i.startTransaction()
	if err != nil {
		return err
	}
	return i.setTransactionIntoMap(tx)
}

func (i *ITsql) startTransaction() (*sql.Tx, error) {
	db, err := i.getConnection()
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

func (i *ITsql) setTransactionIntoMap(tx *sql.Tx) error {
	i.initTransactionMapIfNeed()
	_, ok := i.transactions[i.T().Name()]
	if ok {
		return fmt.Errorf("Transaction map key conflicts (%s)", i.T().Name())
	}
	i.transactions[i.T().Name()] = tx
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

func (i *ITsql) TearDownTest(suiteName, testName string) {
	tx, err := i.getTransactionFromMap()
	must.NotFail(err)
	err = tx.Rollback()
	must.NotFail(err)
	i.deleteTransationFromMap()
	i.it.TearDownTest()
}

func (i *ITsql) getTransactionFromMap() (*sql.Tx, error) {
	tx, ok := i.transactions[i.T().Name()]
	if !ok {
		return nil, fmt.Errorf("Transaction of test (%s) is not set yet", i.T().Name())
	}

	return tx, nil
}

func (i *ITsql) deleteTransationFromMap() {
	delete(i.transactions, i.T().Name())
}
