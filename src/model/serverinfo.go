package model

import (
	"time"
	"sync"
	"github.com/kataras/iris/core/errors"
	"github.com/go-xorm/xorm"
)

type tblServerInfo struct {
	Serverid   		int       `xorm:"INT NOT NULL"`
	Servername 		string    `xorm:"VARCHAR(64) NOT NULL UNIQUE"`
	Islocalserver  	bool      `xorm:"BOOL NOT NULL"`
	Serverpubkey    string    `xorm:"VARCHAR(256)"`
	Serverstartindex   int       `xorm:"INT NOT NULL"`
	Serverstatus   	int       `xorm:"INT NOT NULL"`
	Createtime 		time.Time `xorm:"DATETIME"`
	Updatetime   	time.Time `xorm:"DATETIME"`
}

type tblServerInfoMgr struct {
	TableName string
	Mutex     *sync.Mutex
}

func (t *tblServerInfoMgr) Init() {
	t.TableName = "tbl_server_info"
	t.Mutex = new(sync.Mutex)
}

func (t *tblServerInfoMgr) InsertServerInfo(dbSession *xorm.Session, ServerId int, ServerName string, IsLocalServer bool, ServerPubKey string,
	ServerStartIndex int, ServerStatus int) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	var serverInfo tblServerInfo
	count, err := dbSession.Where("serverid=?", ServerId).Count(&serverInfo)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("serverid already exist!")
	}
	serverInfo.Serverid = ServerId
	serverInfo.Servername = ServerName
	serverInfo.Islocalserver = IsLocalServer
	serverInfo.Serverpubkey = ServerPubKey
	serverInfo.Serverstartindex = ServerStartIndex
	serverInfo.Serverstatus = ServerStatus
	serverInfo.Createtime = time.Now()
	serverInfo.Updatetime = time.Now()

	_, err = dbSession.Insert(&serverInfo)
	return err
}

func (t *tblServerInfoMgr) GetServerInfoById(dbSession *xorm.Session, ServerId int) (bool, tblServerInfo, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	var serverInfo tblServerInfo
	isFound, err := dbSession.Where("serverid=?", ServerId).Get(&serverInfo)
	return isFound, serverInfo, err
}

func (t *tblServerInfoMgr) GetServerInfoCountById(dbSession *xorm.Session, ServerId int) (bool, tblServerInfo, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	var serverInfo tblServerInfo
	isFound, err := dbSession.Where("serverid=?", ServerId).Get(&serverInfo)
	return isFound, serverInfo, err
}

func (t *tblServerInfoMgr) GetLocalServerInfo(dbSession *xorm.Session) (bool, tblServerInfo, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	var serverInfo tblServerInfo
	isFound, err := dbSession.Where("islocalserver=?", true).Get(&serverInfo)
	return isFound, serverInfo, err
}

func (t *tblServerInfoMgr) GetAllServerInfo(dbSession *xorm.Session) ([]tblServerInfo, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	serverInfos := make([]tblServerInfo, 0)
	err := dbSession.Where("").Find(&serverInfos)
	return serverInfos, err
}

