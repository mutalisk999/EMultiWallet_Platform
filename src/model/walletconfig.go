package model

import (
	"github.com/kataras/iris/core/errors"
	"sync"
	"time"
	"github.com/go-xorm/xorm"
)

type tblWalletConfig struct {
	Walletid     int       `xorm:"pk INTEGER autoincr"`
	Walletuuid   string    `xorm:"VARCHAR(64) NOT NULL UNIQUE index"`
	Coinid       int       `xorm:"INTEGER"`
	Walletname   string    `xorm:"VARCHAR(64) NOT NULL UNIQUE index"`
	Serverkeys         string    `xorm:"VARCHAR(64)"`
	Createserver int       `xorm:"INTEGER NOT NULL"`
	Keycount     int       `xorm:"INTEGER"`
	Needkeysigcount int    `xorm:"INTEGER"`
	Address      string    `xorm:"VARCHAR(64) NOT NULL UNIQUE index"`
	Destaddress  string    `xorm:"TEXT"`
	Needsigcount int       `xorm:"INTEGER NOT NULL"`
	Fee          string    `xorm:"VARCHAR(64)"`
	Gasprice     string    `xorm:"VARCHAR(64)"`
	Gaslimit     string    `xorm:"VARCHAR(64)"`
	State        int       `xorm:"INTEGER NOT NULL"`
	Createtime   time.Time `xorm:"DATETIME"`
	Updatetime   time.Time `xorm:"DATETIME"`
}
type tblWalletConfigMgr struct {
	TableName string
	Mutex     *sync.Mutex
}
type WalletRelationinfo struct {
	tblWalletConfig       `xorm:"extends"`
	tblAcctWalletRelation `xorm:"extends"`
}

func (WalletRelationinfo) TableName() string {
	return "tbl_wallet_config"
}

func (t *tblWalletConfigMgr) Init() {
	t.TableName = "tbl_wallet_config"
	t.Mutex = new(sync.Mutex)
}

func (t *tblWalletConfigMgr) InsertWallet(dbSession *xorm.Session, AssetId int, Walletuuid string, Walletname string, Serverkeys string,
	Createserver int, Keycount int, Needkeysigcount int, Address string, Destaddress string, Needsigcount int,
	Fee string, GasPrice string, GasLimit string, State int) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var wallet tblWalletConfig
	wallet.Address = Address
	wallet.Walletuuid = Walletuuid
	wallet.Coinid = AssetId
	wallet.Createtime = time.Now()
	wallet.Destaddress = Destaddress
	wallet.Serverkeys = Serverkeys
	wallet.Createserver = Createserver
	wallet.Keycount = Keycount
	wallet.Needkeysigcount = Needkeysigcount
	wallet.Needsigcount = Needsigcount
	wallet.Fee = Fee
	wallet.Gasprice = GasPrice
	wallet.Gaslimit = GasLimit
	wallet.State = State
	wallet.Updatetime = time.Now()
	wallet.Walletname = Walletname
	_, err := dbSession.Insert(&wallet)
	return err
}

//func (t *tblWalletConfigMgr) UpdateWallet(walletid int, Walletname string, Destaddress string, Needsigcount int, Fee string, GasPrice string, GasLimit string, State int) error {
//	t.Mutex.Lock()
//	defer t.Mutex.Unlock()
//	var wallet tblWalletConfig
//	result, err := GetDBEngine().Where("walletid=?", walletid).Get(&wallet)
//	if err != nil {
//		return err
//	}
//	if result {
//		wallet.Destaddress = Destaddress
//		wallet.Needsigcount = Needsigcount
//		wallet.Fee = Fee
//		wallet.Gasprice = GasPrice
//		wallet.Gaslimit = GasLimit
//		wallet.State = State
//		wallet.Updatetime = time.Now()
//		wallet.Walletname = Walletname
//		_, err := GetDBEngine().Where("walletid=?", walletid).Cols("needsigcount","state","updatetime","walletname","destaddress","fee","gasprice","gaslimit").Update(&wallet)
//		return err
//	} else {
//		return errors.New("not find wallet")
//	}
//
//}

func (t *tblWalletConfigMgr) UpdateWallet(dbSession *xorm.Session, walletid int, Needsigcount int) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var wallet tblWalletConfig
	result, err := dbSession.Where("walletid=?", walletid).Get(&wallet)
	if err != nil {
		return err
	}
	if result {
		wallet.Needsigcount = Needsigcount
		wallet.Updatetime = time.Now()
		_, err := dbSession.Where("walletid=?", walletid).Cols("needsigcount").Update(&wallet)
		return err
	} else {
		return errors.New("not find wallet")
	}
}

func (t *tblWalletConfigMgr) ListWallets(dbSession *xorm.Session, coinids []int, state []int, acctids []int, offset int, limit int) ([]tblWalletConfig, int64, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	result := make([]tblWalletConfig, 0)
	walls := make([]tblAcctWalletRelation, 0)
	wallids := make(map[int]int, 0)
	wids:=make([]int, 0)
	if len(acctids)!=0{
		err :=dbSession.In("acctid",acctids).Find(&walls)
		if err != nil {
			return result, 0, err
		}
		for _,rela :=range walls{
			wallids[rela.Walletid]=0
		}
	}
	for k,_:=range wallids{
		wids=append(wids, k)
	}
	if len(acctids)!=0&&len(wids)==0{
		return result,0,nil
	}
	dbSe := dbSession.Where("")
	if len(wids)!=0{
		dbSe = dbSe.In("walletid", wids)
	}
	if len(coinids)!=0{
		dbSe = dbSe.In("coinid", coinids)
	}
	if len(state) != 0 {
		dbSe = dbSe.In("state", state)
	}
	err := dbSe.Limit(limit,offset).Find(&result)
	if err != nil {
		return result, 0, err
	}
	dbSe = dbSession.Where("")
	if len(wids)!=0{
		dbSe = dbSe.In("walletid", wids)
	}
	if len(coinids)!=0{
		dbSe = dbSe.In("coinid", coinids)
	}
	if len(state) != 0 {
		dbSe = dbSe.In("state", state)
	}
	var EmptyWa tblWalletConfig
	total,err:= dbSe.Count(EmptyWa)
	if err != nil {
		return result, 0, err
	}
	return result, total, nil
}

func (t *tblWalletConfigMgr) GetWalletById(dbSession *xorm.Session, id int) (tblWalletConfig, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var wallet tblWalletConfig
	result, err := dbSession.Where("walletid=?", id).Get(&wallet)
	if err != nil {
		return wallet, err
	}
	if result {
		return wallet, nil
	}
	return wallet, errors.New("no find wallet")
}

func (t *tblWalletConfigMgr) GetWalletByUUId(dbSession *xorm.Session, uuid string) (tblWalletConfig, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var wallet tblWalletConfig
	result, err := dbSession.Where("walletuuid=?", uuid).Get(&wallet)
	if err != nil {
		return wallet, err
	}
	if result {
		return wallet, nil
	}
	return wallet, errors.New("no find wallet")
}

func (t *tblWalletConfigMgr) GetWalletsByIds(dbSession *xorm.Session, ids []int) ([]tblWalletConfig, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var wallets []tblWalletConfig
	err := dbSession.In("walletid", ids).Find(&wallets)
	if err != nil {
		return nil, err
	}
	return wallets, nil
}

func (t *tblWalletConfigMgr) GetWalletsByUUIds(dbSession *xorm.Session, uuids []int) ([]tblWalletConfig, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var wallets []tblWalletConfig
	err := dbSession.In("walletuuid", uuids).Find(&wallets)
	if err != nil {
		return nil, err
	}
	return wallets, nil
}

func (t *tblWalletConfigMgr) GetWalletByName(dbSession *xorm.Session, name string) (tblWalletConfig, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var wallet tblWalletConfig
	result, err := dbSession.Where("walletname=?", name).Get(&wallet)
	if err != nil {
		return wallet, err
	}
	if result {
		return wallet, nil
	}
	return wallet, errors.New("no find wallet")
}

func (t *tblWalletConfigMgr) ChangeWalletState(dbSession *xorm.Session, id int, sta int) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var wallet tblWalletConfig
	result, err := dbSession.Where("walletid=?", id).Get(&wallet)
	if err != nil {
		return err
	}
	if result {
		wallet.Walletid = id
		wallet.Updatetime = time.Now()
		wallet.State = sta
		_, err := dbSession.Where("walletid=?", id).Cols("state","updatetime").Update(&wallet)
		return err
	}
	return errors.New("no find wallet")
}

func (t *tblWalletConfigMgr) ChangeWalletStateByUuid(dbSession *xorm.Session, walletUuid string, sta int) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var wallet tblWalletConfig
	result, err := dbSession.Where("walletuuid=?", walletUuid).Get(&wallet)
	if err != nil {
		return err
	}
	if result {
		wallet.Walletuuid = walletUuid
		wallet.Updatetime = time.Now()
		wallet.State = sta
		_, err := dbSession.Where("walletuuid=?", walletUuid).Cols("state","updatetime").Update(&wallet)
		return err
	}
	return errors.New("no find wallet")
}

func (t *tblWalletConfigMgr) ActiviteWallet(dbSession *xorm.Session, id int) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var wallet tblWalletConfig
	result, err := dbSession.Where("walletid=?", id).Get(&wallet)
	if err != nil {
		return err
	}
	if result {
		wallet.State = 1
		wallet.Updatetime = time.Now()
		_, err := dbSession.Where("walletid=?", id).Update(&wallet)
		return err
	}
	return errors.New("no find wallet")
}

func (t *tblWalletConfigMgr) FreezeWallet(dbSession *xorm.Session, id int) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var wallet tblWalletConfig
	result, err := dbSession.Where("walletid=?", id).Get(&wallet)
	if err != nil {
		return err
	}
	if result {
		wallet.State = 0
		wallet.Updatetime = time.Now()
		_, err := dbSession.Where("walletid=?", id).Cols("updatetime").Update(&wallet)
		return err
	}
	return errors.New("no find wallet")
}

func (t *tblWalletConfigMgr) DeleteWallet(dbSession *xorm.Session, id int) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var wallet tblWalletConfig
	wallet.Walletid = id
	_, err := dbSession.Delete(wallet)
	return err
}
