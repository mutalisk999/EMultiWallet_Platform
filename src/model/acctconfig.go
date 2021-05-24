package model

import (
	"fmt"
	"github.com/kataras/iris/core/errors"
	"sync"
	"time"
	"github.com/go-xorm/xorm"
)

type tblAcctConfig struct {
	Acctid     int       `xorm:"pk INTEGER autoincr"`
	Cellphone  string    `xorm:"VARCHAR(64) NOT NULL"`
	Realname   string    `xorm:"VARCHAR(64) NOT NULL"`
	Idcard     string    `xorm:"VARCHAR(64) NOT NULL"`
	Pubkey     string    `xorm:"VARCHAR(512) NOT NULL UNIQUE"`
	Role       int       `xorm:"INT NOT NULL"`
	State      int       `xorm:"INT NOT NULL"`
	Createtime time.Time `xorm:"DATETIME"`
	Updatetime time.Time `xorm:"DATETIME"`
}

type tblAcctConfigMgr struct {
	TableName string
	Mutex     *sync.Mutex
}

func (t *tblAcctConfigMgr) Init() {
	t.TableName = "tbl_acct_config"
	t.Mutex = new(sync.Mutex)
}

func (t *tblAcctConfigMgr) VerifyUnique(dbSession *xorm.Session, CellPhone string, RealName string, IdCard string, Pubkey string) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var acct tblAcctConfig
	state := []int{0,1,2}
	result, err := dbSession.In("state",state).Where("cellphone=? or realname=? or idcard=? or pubkey=?", CellPhone,RealName,IdCard,Pubkey).Get(&acct)
	if result {
		return errors.New("key already exist!")
	}

	return err
}

func (t *tblAcctConfigMgr) InsertAcct(dbSession *xorm.Session, CellPhone string, RealName string, IdCard string, Pubkey string) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var acct tblAcctConfig
	count, err := dbSession.Count(&acct)
	acct.Cellphone = CellPhone
	acct.Realname = RealName
	acct.Idcard = IdCard
	acct.Pubkey = Pubkey
	acct.State = 0
	acct.Createtime = time.Now()
	acct.Updatetime = time.Now()
	if err != nil {
		fmt.Println("account count ",err.Error())
	}
	if count > 0 {
		acct.Role = 1
	} else {
		acct.Role = 0
		acct.State = 1
	}

	_, err = dbSession.Insert(&acct)
	return err
}

func (t *tblAcctConfigMgr) GetAdminId(dbSession *xorm.Session) int {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var acct tblAcctConfig
	result, err := dbSession.Where("role=0").Get(&acct)
	if err != nil {
		return -1
	}
	if !result {
		return -1
	}
	return acct.Acctid
}

func (t *tblAcctConfigMgr) UpdateAcct(dbSession *xorm.Session, acctid int, CellPhone string, RealName string, IdCard string, Pubkey string, State int) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var acct tblAcctConfig
	result, err := dbSession.Where("acctid=?", acctid).Get(&acct)
	if err != nil {
		return err
	}
	if result {
		acct.Cellphone = CellPhone
		acct.Idcard = IdCard
		acct.Pubkey = Pubkey
		acct.State = State
		acct.Updatetime = time.Now()
		_, err = dbSession.Where("acctid=?", acctid).Update(&acct)
		return err
	} else {
		return errors.New("not find account")
	}
}

func (t *tblAcctConfigMgr) UpdateAcctState(dbSession *xorm.Session, acctid int, State int) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var acct tblAcctConfig
	result, err := dbSession.Where("acctid=?", acctid).Get(&acct)
	if err != nil {
		return err
	}
	if result {
		acct.State = State
		acct.Updatetime = time.Now()
		_, err = dbSession.Where("acctid=?", acctid).Cols("state").Update(&acct)
		return err
	} else {
		return errors.New("not find account")
	}
}

func (t *tblAcctConfigMgr) ActiveAcct(dbSession *xorm.Session, acctid int) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var acct tblAcctConfig
	result, err := dbSession.Where("acctid=?", acctid).Get(&acct)
	if err != nil {
		return err
	}
	if result {
		acct.State = 1
		acct.Updatetime = time.Now()
		_, err = dbSession.Where("acctid=?", acctid).Update(&acct)
		return err
	} else {
		return errors.New("not find account")
	}
}

func (t *tblAcctConfigMgr) FreezeAcct(dbSession *xorm.Session, acctid int) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var acct tblAcctConfig
	result, err := dbSession.Where("acctid=?", acctid).Get(&acct)
	if err != nil {
		return err
	}
	if result {
		acct.State = 2
		acct.Updatetime = time.Now()
		_, err = dbSession.Where("acctid=?", acctid).Update(&acct)
		return err
	} else {
		return errors.New("not find account")
	}
}

func (t *tblAcctConfigMgr) DeleteAcct(dbSession *xorm.Session, acctid int) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var acct tblAcctConfig
	acct.Acctid = acctid
	_, err := dbSession.Delete(acct)
	return err
}

func (t *tblAcctConfigMgr) FindAccountByPubkey(dbSession *xorm.Session, pubkey string) (tblAcctConfig, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var acct tblAcctConfig
	result, err := dbSession.Where("pubkey=?", pubkey).Get(&acct)
	if result {
		return acct, err
	} else {
		return acct, errors.New("Not Found")
	}
}

func (t *tblAcctConfigMgr) ListNormalAccount(dbSession *xorm.Session, state []int, limit, offset int) ([]tblAcctConfig, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var acct []tblAcctConfig
	err := dbSession.In("state", state).And("role=?", 1).Limit(limit, offset).Find(&acct)

	return acct, err
}

func (t *tblAcctConfigMgr) GetNormalAccountCount(dbSession *xorm.Session, state []int) (int64, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var acct tblAcctConfig
	count, err := dbSession.In("state", state).And("role=?", 1).Count(&acct)

	return count, err
}

func (t *tblAcctConfigMgr) GetAccountById(dbSession *xorm.Session, acctId int) (tblAcctConfig, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var acct tblAcctConfig
	exist, err := dbSession.Where("acctid=?", acctId).Get(&acct)
	if err != nil {
		return acct, err
	}
	if !exist {
		return acct, errors.New("Account Not Found")
	}

	return acct, err
}

func (t *tblAcctConfigMgr) GetAccountCount(dbSession *xorm.Session) (int64, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var acct tblAcctConfig
	count, err := dbSession.Count(&acct)
	return count,err
}

func (t *tblAcctConfigMgr) GetAccountIdByPubkey(dbSession *xorm.Session, pubkey string) (int, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var acct tblAcctConfig
	exist, err := dbSession.Where("pubkey=?", pubkey).Get(&acct)
	if err != nil {
		return -1, err
	}
	if !exist {
		return -1, errors.New("Account Not Found")
	}

	return acct.Acctid, err
}

func (t *tblAcctConfigMgr) GetAccountsByIds(dbSession *xorm.Session, acctIds []int) ([]tblAcctConfig, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var accts []tblAcctConfig
	if len(acctIds)>0{
		err := dbSession.In("acctid", acctIds).Find(&accts)
		if err != nil {
			return nil, err
		}
	}else if len(acctIds)==0{
		err := dbSession.Find(&accts)
		if err != nil {
			return nil, err
		}
	}

	return accts, nil
}
