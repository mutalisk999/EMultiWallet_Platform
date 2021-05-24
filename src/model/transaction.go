package model

import (
	"errors"
	"sync"
	"time"
	"github.com/go-xorm/xorm"
)

type tblTransaction struct {
	Trxid         int    `xorm:"pk INTEGER autoincr"`
	Trxuuid       string `xorm:"VARCHAR(64) NOT NULL UNIQUE index"`
	Rawtrxid      string `xorm:"VARCHAR(128)"`
	Walletid      int    `xorm:"INT NOT NULL"`
	Coinid        int    `xorm:"INT NOT NULL"`
	Contractaddr  string `xorm:"VARCHAR(128)"`
	Acctid        int    `xorm:"INT NOT NULL"`
	Serverid      int    `xorm:"INT NOT NULL"`
	Fromaddr      string `xorm:"VARCHAR(128) NOT NULL"`
	Todetails     string `xorm:"VARCHAR(2048) NOT NULL"`
	Feecost       string `xorm:"VARCHAR(128)"`
	Trxtime       time.Time
	Needconfirm   int    `xorm:"INT NOT NULL"`
	Confirmed     int    `xorm:"INT NOT NULL"`
	Acctconfirmed string `xorm:"VARCHAR(1024) NOT NULL"`
	Signedtrxs    string `xorm:"VARCHAR(102400)"`
	Signedserverids string `xorm:"VARCHAR(1024)"`
	Fee           string `xorm:"VARCHAR(128)"`
	Gasprice      string `xorm:"VARCHAR(128)"`
	Gaslimit      string `xorm:"VARCHAR(128)"`
	State         int    `xorm:"INT NOT NULL"`
	Signature     string `xorm:"VARCHAR(256)"`
}

type tblTransactionMgr struct {
	TableName string
	Mutex     *sync.Mutex
}

func (t *tblTransactionMgr) Init() {
	t.TableName = "tbl_transaction"
	t.Mutex = new(sync.Mutex)
}

func (t *tblTransactionMgr) NewTransaction(dbSession *xorm.Session, trxUuid string, walletId int, coinId int, contractAddr string, acctId int,
	serverId int, fromAddr string, toDetails string, needConfirm int,
	fee string, gasPrice string, gasLimit string, signature string) (int, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	var transaction tblTransaction
	transaction.Trxuuid = trxUuid
	transaction.Walletid = walletId
	transaction.Coinid = coinId
	transaction.Contractaddr = contractAddr
	transaction.Acctid = acctId
	transaction.Serverid = serverId
 	transaction.Fromaddr = fromAddr
	transaction.Todetails = toDetails
	transaction.Trxtime = time.Now()
	transaction.Needconfirm = needConfirm
	transaction.Confirmed = 0
	transaction.Fee = fee
	transaction.Gasprice = gasPrice
	transaction.Gaslimit = gasLimit
	transaction.Signature = signature
	transaction.State = 0
	_, err := dbSession.Insert(&transaction)
	if err != nil {
		return 0, err
	}
	return transaction.Trxid, nil
}

func (t *tblTransactionMgr) GetTransactionById(dbSession *xorm.Session, trxId int) (bool, tblTransaction, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	var trx tblTransaction
	isFound, err := dbSession.Where("trxid=?", trxId).Get(&trx)
	if err != nil {
		return false, tblTransaction{}, err
	}
	if isFound {
		return true, trx, nil
	}
	return false, tblTransaction{}, errors.New("no find transaction")
}

func (t *tblTransactionMgr) GetTransactionByUuId(dbSession *xorm.Session, trxUuid string) (bool, tblTransaction, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	var trx tblTransaction
	isFound, err := dbSession.Where("trxuuid=?", trxUuid).Get(&trx)
	return isFound, trx, err
}

func (t *tblTransactionMgr) UpdateTransaction(dbSession *xorm.Session, transaction tblTransaction) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	_, err := dbSession.ID(transaction.Trxid).Update(transaction)
	if err != nil {
		return err
	}
	return nil
}

func (t *tblTransactionMgr) UpdateTransactionState(dbSession *xorm.Session, trxId int, state int) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	var transaction tblTransaction
	transaction.State = state
	_, err := dbSession.ID(trxId).Cols("state").Update(&transaction)
	if err != nil {
		return err
	}
	return nil
}

func (t *tblTransactionMgr) UpdateTransactionStateFeeCost(dbSession *xorm.Session, trxId int, state int, feeCost *string) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	var transaction tblTransaction
	transaction.State = state
	if feeCost != nil {
		transaction.Feecost = *feeCost
	}
	_, err := dbSession.ID(trxId).Update(&transaction)
	if err != nil {
		return err
	}
	return nil
}

func (t *tblTransactionMgr) GetTransactionsByState(dbSession *xorm.Session, state int) ([]tblTransaction, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	var trxs []tblTransaction
	err := dbSession.Where("state=?", state).Find(&trxs)
	return trxs, err
}

func (t *tblTransactionMgr) GetUnComfirmedTransactions(dbSession *xorm.Session) ([]tblTransaction, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	var trxs []tblTransaction
	err := dbSession.SQL("select * from tbl_transaction where rawtrxid is not null and rawtrxid != \"\" and state = 2").Find(&trxs)
	return trxs, err
}

func (t *tblTransactionMgr) GetTransactions(dbSession *xorm.Session, walletId []int, coinId []int, serverId []int, acctId []int, state []int,
	trxTime [2]string, offSet int, limit int) (int, []tblTransaction, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	dbSession1 := dbSession.Where("")
	if walletId != nil && len(walletId) != 0 {
		dbSession1 = dbSession1.In("walletid", walletId)
	}
	if coinId != nil && len(coinId) != 0 {
		dbSession1 = dbSession1.In("coinid", coinId)
	}
	if serverId != nil && len(serverId) != 0 {
		dbSession1 = dbSession1.In("serverid", serverId)
	}
	if acctId != nil && len(acctId) != 0 {
		dbSession1 = dbSession1.In("acctid", acctId)
	}
	if state != nil && len(state) != 0 {
		dbSession1 = dbSession1.In("state", state)
	}
	if trxTime[0] != "" {
		dbSession1 = dbSession1.And("trxtime > ?", trxTime[0])
	}
	if trxTime[1] != "" {
		dbSession1 = dbSession1.And("trxtime < ?", trxTime[1])
	}
	var trx tblTransaction
	total, err := dbSession1.Count(&trx)
	if err != nil {
		return 0, nil, err
	}

	dbSession2 := GetDBEngine().Where("")
	if walletId != nil && len(walletId) != 0 {
		dbSession2 = dbSession2.In("walletid", walletId)
	}
	if coinId != nil && len(coinId) != 0 {
		dbSession2 = dbSession2.In("coinid", coinId)
	}
	if serverId != nil && len(serverId) != 0 {
		dbSession2 = dbSession2.In("serverid", serverId)
	}
	if acctId != nil && len(acctId) != 0 {
		dbSession2 = dbSession2.In("acctid", acctId)
	}
	if state != nil && len(state) != 0 {
		dbSession2 = dbSession2.In("state", state)
	}
	if trxTime[0] != "" {
		dbSession2 = dbSession2.And("trxtime > ?", trxTime[0])
	}
	if trxTime[1] != "" {
		dbSession2 = dbSession2.And("trxtime < ?", trxTime[1])
	}
	trxs := make([]tblTransaction, 0)
	dbSession2.Limit(limit, offSet).Desc("trxtime").Find(&trxs)
	return int(total), trxs, nil
}
