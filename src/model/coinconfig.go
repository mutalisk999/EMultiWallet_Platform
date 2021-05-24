package model

import (
	"github.com/kataras/iris/core/errors"
	"sync"
	"time"
	"github.com/go-xorm/xorm"
)

type tblCoinConfig struct {
	Coinid     int       `xorm:"pk INTEGER autoincr"`
	Coinsymbol string    `xorm:"VARCHAR(16) NOT NULL UNIQUE"`
	Ip         string    `xorm:"VARCHAR(64) NOT NULL"`
	Rpcport    int       `xorm:"INT NOT NULL"`
	Rpcuser    string    `xorm:"VARCHAR(64)"`
	Rpcpass    string    `xorm:"VARCHAR(64)"`
	State      int       `xorm:"INT NOT NULL"`
	Createtime time.Time `xorm:"DATETIME"`
	Updatetime time.Time `xorm:"DATETIME"`
}

type tblCoinConfigMgr struct {
	TableName string
	Mutex     *sync.Mutex
}

func (t *tblCoinConfigMgr) Init() {
	t.TableName = "tbl_coin_config"
	t.Mutex = new(sync.Mutex)
}

func (t *tblCoinConfigMgr) ListCoins(dbSession *xorm.Session) ([]tblCoinConfig, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var coins []tblCoinConfig
	err := dbSession.Find(&coins)
	return coins, err
}

func (t *tblCoinConfigMgr) InsertCoin(dbSession *xorm.Session, coinSymbol string, ip string, port int, rpcUserName, rpcPassword string, state int) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var coin tblCoinConfig
	coin.Coinsymbol = coinSymbol
	coin.Ip = ip
	coin.Rpcport = port
	coin.Rpcuser = rpcUserName
	coin.Rpcpass = rpcPassword
	coin.State = state
	coin.Createtime = time.Now()
	coin.Updatetime = time.Now()
	_, err := dbSession.Insert(&coin)
	return err
}

func (t *tblCoinConfigMgr) InsertCoinWithCoinId(dbSession *xorm.Session, coinId int, coinSymbol string, ip string, port int, rpcUserName, rpcPassword string, state int) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var coin tblCoinConfig
	coin.Coinid = coinId
	coin.Coinsymbol = coinSymbol
	coin.Ip = ip
	coin.Rpcport = port
	coin.Rpcuser = rpcUserName
	coin.Rpcpass = rpcPassword
	coin.State = state
	coin.Createtime = time.Now()
	coin.Updatetime = time.Now()
	_, err := dbSession.Insert(&coin)
	return err
}

func (t *tblCoinConfigMgr) UpdateCoin(dbSession *xorm.Session, coinId int, coinSymbol string, ip string, port int, rpcUserName, rpcPassword string) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var coin tblCoinConfig
	exist, err := dbSession.Where("coinid=?", coinId).Get(&coin)
	if !exist {
		return errors.New("coinId not found!")
	}
	if err != nil {
		return err
	}

	// not allow to modify coin symbol
	if coin.Coinsymbol != coinSymbol {
		return errors.New("not allow to modify coinSymbol")
	}

	coin.Coinsymbol = coinSymbol
	coin.Ip = ip
	coin.Rpcport = port
	coin.Rpcuser = rpcUserName
	coin.Rpcpass = rpcPassword
	coin.Updatetime = time.Now()
	_, err = dbSession.Where("coinid=?", coinId).Update(&coin)
	return err
}

func (t *tblCoinConfigMgr) UpdateCoinState(dbSession *xorm.Session, coinId int, state int) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var coin tblCoinConfig
	exist, err := dbSession.Where("coinid=?", coinId).Get(&coin)
	if !exist {
		return errors.New("coinId not found!")
	}
	if err != nil {
		return err
	}
	coin.State = state
	coin.Updatetime = time.Now()
	_, err = dbSession.Where("coinid=?", coinId).Update(&coin)
	return err
}

func (t *tblCoinConfigMgr) GetCoin(dbSession *xorm.Session, coinId int) (tblCoinConfig, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var coin tblCoinConfig
	result, err := dbSession.Where("coinid=?", coinId).Get(&coin)
	if !result {
		return coin, errors.New("key not found")
	}
	return coin, err
}

func (t *tblCoinConfigMgr) GetCoin2(dbSession *xorm.Session, coinId int) (bool, tblCoinConfig, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var coin tblCoinConfig
	isFound, err := dbSession.Where("coinid=?", coinId).Get(&coin)
	return isFound, coin, err
}

func (t *tblCoinConfigMgr) GetCoins(dbSession *xorm.Session, coinIds []int) ([]tblCoinConfig, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var coins []tblCoinConfig
	err := dbSession.In("coinid", coinIds).Find(&coins)
	if err != nil {
		return nil, err
	}
	return coins, err
}

func (t *tblCoinConfigMgr) GetCoinBySymbol(dbSession *xorm.Session, symbol string) (tblCoinConfig, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var coin tblCoinConfig
	result, err := dbSession.Where("coinsymbol=?", symbol).Get(&coin)
	if !result {
		return coin, errors.New("symbol not found")
	}
	return coin, err
}
