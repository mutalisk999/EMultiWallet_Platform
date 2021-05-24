package model

import (
	"sync"
	"time"
	"github.com/go-xorm/xorm"
)

type tblPendingTransaction struct {
	Id            int    `xorm:"pk INTEGER autoincr"`
	Trxuuid       string `xorm:"VARCHAR(64) NOT NULL UNIQUE index"`
	Coinid        int    `xorm:"INT NOT NULL"`
	Vintrxid      string `xorm:"VARCHAR(128)"`
	Vinvout       int    `xorm:"INT"`
	Fromaddress   string `xorm:"VARCHAR(128)"`
	Balance       string `xorm:"VARCHAR(64)"`
	Createtime    time.Time `xorm:"DATETIME"`
	Updatetime    time.Time `xorm:"DATETIME"`
}

type tblPendingTransactionMgr struct {
	TableName string
	Mutex     *sync.Mutex
}

func (t *tblPendingTransactionMgr) Init() {
	t.TableName = "tbl_pending_transaction"
	t.Mutex = new(sync.Mutex)
}

func (t *tblPendingTransactionMgr) NewPendingTransaction(dbSession *xorm.Session, trxUuid string, coinId int, vinTrxid string, vinVout int,
	fromAddress string, balance string) (int, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	var pendingTransaction tblPendingTransaction
	pendingTransaction.Trxuuid = trxUuid
	pendingTransaction.Coinid = coinId
	pendingTransaction.Vintrxid = vinTrxid
	pendingTransaction.Vinvout = vinVout
	pendingTransaction.Fromaddress = fromAddress
	pendingTransaction.Balance = balance
	pendingTransaction.Createtime = time.Now()
	pendingTransaction.Updatetime = time.Now()

	_, err := dbSession.Insert(&pendingTransaction)
	if err != nil {
		return 0, err
	}
	return pendingTransaction.Id, nil
}

func (t *tblPendingTransactionMgr) GetPendingTransactionByTrxUuid(dbSession *xorm.Session, trxUuid string) ([]tblPendingTransaction, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	pendingTransactions := make([]tblPendingTransaction, 0)
	err := dbSession.Where("trxuuid=?", trxUuid).Find(&pendingTransactions)
	return pendingTransactions, err
}

func (t *tblPendingTransactionMgr) GetPendingTransactionByCoinId(dbSession *xorm.Session, coinId int) ([]tblPendingTransaction, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	pendingTransactions := make([]tblPendingTransaction, 0)
	err := dbSession.Where("coinid=?", coinId).Find(&pendingTransactions)
	return pendingTransactions, err
}

func (t *tblPendingTransactionMgr) DeletePendingTransactionByTrxUuid(dbSession *xorm.Session, trxUuid string) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	pendingTransactions := make([]tblPendingTransaction, 0)
	err := dbSession.Where("trxuuid=?", trxUuid).Find(&pendingTransactions)
	if err != nil {
		return err
	}

	for _, pendingTransaction := range pendingTransactions {
		_, err = dbSession.Delete(pendingTransaction)
	}

	return err
}

func (t *tblPendingTransactionMgr) GetPendingTransactionByVinCount(dbSession *xorm.Session, coinId int, vinTrxId string, vinVout int) (int, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	var pendingTransactions tblPendingTransaction
	count, err := dbSession.Where("coinid=?", coinId).And("vintrxid=?", vinTrxId).And("vinvout=?", vinVout).Count(&pendingTransactions)

	if err != nil {
		return 0, err
	} else if count == 0 {
		return 0, nil
	} else {
		return int(count), nil
	}
	return 0, nil
}



