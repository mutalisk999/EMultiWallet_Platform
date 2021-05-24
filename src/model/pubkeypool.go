package model

import (
	"github.com/kataras/iris/core/errors"
	"sync"
	"time"
	"github.com/go-xorm/xorm"
)

type ServerKeyInfo struct {
	ServerId   int
	StartIndex int
	KeyIndex   int
	PubKey     string
}

type tblPubkeyPool struct {
	Serverid   int       `xorm:"INT NOT NULL"`
	Keyindex   int       `xorm:"INT NOT NULL"`
	Pubkey     string    `xorm:"VARCHAR(256) NOT NULL UNIQUE"`
	Isused     bool      `xorm:"BOOL NOT NULL"`
	Createtime time.Time `xorm:"DATETIME"`
	Usedtime   time.Time `xorm:"DATETIME"`
}

type tblPubKeyPoolMgr struct {
	TableName string
	Mutex     *sync.Mutex
}

func (t *tblPubKeyPoolMgr) Init() {
	t.TableName = "tbl_pubkey_pool"
	t.Mutex = new(sync.Mutex)
}

func (t *tblPubKeyPoolMgr) InsertPubkey(dbSession *xorm.Session, ServerId int, KeyIndex int, Pubkey string) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	var pubkey_record tblPubkeyPool
	count, err := dbSession.Where("serverid=?", ServerId).And("keyindex=?", KeyIndex).Count(&pubkey_record)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("key already exist!")
	}
	pubkey_record.Serverid = ServerId
	pubkey_record.Keyindex = KeyIndex
	pubkey_record.Pubkey = Pubkey
	pubkey_record.Isused = false
	pubkey_record.Createtime = time.Now()
	pubkey_record.Usedtime = time.Now()

	_, err = dbSession.Insert(&pubkey_record)
	return err
}

func (t *tblPubKeyPoolMgr) UpdatePubkey(dbSession *xorm.Session, ServerId int, KeyIndex int, Pubkey string, isuse bool) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	var pubkey_record tblPubkeyPool
	exist, err := dbSession.Where("serverid=?", ServerId).And("keyindex=?", KeyIndex).Get(&pubkey_record)
	if err != nil {
		return err
	}
	if !exist {
		return errors.New("key not found!")
	}
	pubkey_record.Serverid = ServerId
	pubkey_record.Keyindex = KeyIndex
	pubkey_record.Pubkey = Pubkey
	pubkey_record.Isused = isuse
	pubkey_record.Usedtime = time.Now()

	_, err = dbSession.Where("serverid=?", ServerId).And("keyindex=?", KeyIndex).Cols("isused").Update(&pubkey_record)
	return err
}

func (t *tblPubKeyPoolMgr) UsePubkey(dbSession *xorm.Session, ServerId int, KeyIndex int) (string, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	var pubkey_record tblPubkeyPool
	exist, err := dbSession.Where("serverid=?", ServerId).And("keyindex=?", KeyIndex).Get(&pubkey_record)
	if err != nil {
		return "", err
	}
	if !exist {
		return "", errors.New("key not found!")
	}
	if pubkey_record.Isused {
		return "", errors.New("key has been used before")
	}
	pubkey_record.Serverid = ServerId
	pubkey_record.Keyindex = KeyIndex
	pubkey_record.Isused = true
	pubkey_record.Usedtime = time.Now()

	_, err = dbSession.Where("serverid=?", ServerId).And("keyindex=?", KeyIndex).Cols("isused").Update(&pubkey_record)
	return pubkey_record.Pubkey, err
}

func (t *tblPubKeyPoolMgr) GetAnUnusedKeyIndex(dbSession *xorm.Session, ServerId int, StartIndex int) (int, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	var pubkey_record tblPubkeyPool
	exist, err := dbSession.Where("serverid=?", ServerId).And("keyindex>=?", StartIndex).
		And("isused=?", false).Get(&pubkey_record)
	if err != nil {
		return -1, errors.New("query key error!")
	}
	if !exist {
		return -1, errors.New("key not found!")
	}
	pubkey_record.Isused = true
	pubkey_record.Usedtime = time.Now()

	_, err = dbSession.Where("serverid=?", ServerId).And("keyindex=?", pubkey_record.Keyindex).Update(&pubkey_record)
	return pubkey_record.Keyindex, err
}

func (t *tblPubKeyPoolMgr) GetUnusedServerKeys(dbSession *xorm.Session, ServerKeys []ServerKeyInfo) ([]ServerKeyInfo, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	ServerKeysRet := make([]ServerKeyInfo, 0)
	for _, serverKey := range ServerKeys {
		var pubkey_record tblPubkeyPool
		exist, err := dbSession.Where("serverid=?", serverKey.ServerId).And("keyindex>=?", serverKey.StartIndex).
			And("isused=?", false).Get(&pubkey_record)
		if err != nil {
			return nil, errors.New("query key error!")
		}
		if !exist {
			return nil, errors.New("key not found!")
		}
		serverKey.PubKey = pubkey_record.Pubkey
		serverKey.KeyIndex = pubkey_record.Keyindex
		ServerKeysRet = append(ServerKeysRet, serverKey)
	}

	for _, serverKey := range ServerKeysRet {
		var pubkey_record tblPubkeyPool
		pubkey_record.Serverid = serverKey.ServerId
		pubkey_record.Keyindex = serverKey.KeyIndex
		pubkey_record.Pubkey = serverKey.PubKey
		pubkey_record.Isused = true
		pubkey_record.Usedtime = time.Now()
		_, err := dbSession.Where("serverid=?", serverKey.ServerId).And("keyindex=?", serverKey.KeyIndex).
			Cols("isused").Update(&pubkey_record)
		if err != nil {
			return nil, err
		}
	}

	return ServerKeysRet, nil
}

func (t *tblPubKeyPoolMgr) RollBackUsedServerKeys(dbSession *xorm.Session, ServerKeys []ServerKeyInfo) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	for _, serverKey := range ServerKeys {
		var pubkey_record tblPubkeyPool
		pubkey_record.Serverid = serverKey.ServerId
		pubkey_record.Keyindex = serverKey.KeyIndex
		pubkey_record.Pubkey = serverKey.PubKey
		pubkey_record.Isused = false
		_, err := dbSession.Where("serverid=?", serverKey.ServerId).And("keyindex=?", serverKey.KeyIndex).
			Cols("isused", "usedtime").Update(&pubkey_record)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *tblPubKeyPoolMgr) QueryPubKeyByKeyIndex(dbSession *xorm.Session, ServerId int, KeyIndex int) (string, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var pubkey_record tblPubkeyPool
	exist, err := dbSession.Where("serverid=?", ServerId).And("keyindex=?", KeyIndex).And("isused=?", true).Get(&pubkey_record)
	if err != nil {
		return "", err
	}
	if !exist {
		return "", errors.New("key not found!")
	}
	return pubkey_record.Pubkey, nil
}

func (t *tblPubKeyPoolMgr) LoadPubKeysByServerId(dbSession *xorm.Session, ServerId int) ([]tblPubkeyPool, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	serverKeyInfos := make([]tblPubkeyPool, 0)
	err := dbSession.Where("serverid=?", ServerId).Find(&serverKeyInfos)
	if err != nil {
		return nil, err
	}
	return serverKeyInfos, nil
}

func (t *tblPubKeyPoolMgr) UpdateServerIdByServerId(dbSession *xorm.Session, ServerId int, UpdatedServerId int) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	var pubkey_record tblPubkeyPool
	pubkey_record.Serverid = UpdatedServerId
	_, err := dbSession.Where("serverid=?", ServerId).Cols("serverid").Update(&pubkey_record)
	if err != nil {
		return err
	}
	return nil
}

func (t *tblPubKeyPoolMgr) GetPubkeyByIdIndex(dbSession *xorm.Session,keyindex,serverid int) (*tblPubkeyPool,error){
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	var oneKey tblPubkeyPool
	exist,err:=dbSession.Where("serverid=? and keyindex=?",serverid,keyindex).Get(&oneKey)
	if err!=nil{
		return nil,err
	}
	if !exist{
		return nil,errors.New("key not found!")
	}
	return &oneKey,nil
}
