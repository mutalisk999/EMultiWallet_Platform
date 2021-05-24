package model

import (
	"sync"
	"github.com/kataras/iris/core/errors"
	"github.com/go-xorm/xorm"
)

type tblSessionPersistence struct {
	Sessionid      string    `xorm:"VARCHAR(64) NOT NULL"`
	Sessionvalue   string    `xorm:"VARCHAR(1024) NOT NULL"`
}

type tblSessionPersistenceMgr struct {
	TableName string
	Mutex     *sync.Mutex
}

func (t *tblSessionPersistenceMgr) Init() {
	t.TableName = "tbl_session_persistence"
	t.Mutex = new(sync.Mutex)
}

func (t *tblSessionPersistenceMgr) GetSessionValue(dbSession *xorm.Session, SessionId string) (tblSessionPersistence, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var sessionPersistence tblSessionPersistence
	result, err := dbSession.Where("sessionid=?", SessionId).Get(&sessionPersistence)
	if result {
		return sessionPersistence, err
	} else {
		return sessionPersistence, errors.New("Not Found")
	}
}

func (t *tblSessionPersistenceMgr) LoadAllSessionValue(dbSession *xorm.Session) ([]tblSessionPersistence, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	sessionPersistences := make([]tblSessionPersistence, 0)
	err := dbSession.Find(&sessionPersistences)
	if err != nil {
		return nil, err
	}
	return sessionPersistences, nil
}

func (t *tblSessionPersistenceMgr) InsertSessionValue(dbSession *xorm.Session, SessionId string, SessionValue string) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var sessionPersistence tblSessionPersistence
	sessionPersistence.Sessionid = SessionId
	sessionPersistence.Sessionvalue = SessionValue

	_, err := dbSession.Insert(&sessionPersistence)

	return err
}

func (t *tblSessionPersistenceMgr) UpdateSessionValue(dbSession *xorm.Session, SessionId string, SessionValue string) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var sessionPersistence tblSessionPersistence
	result, err := dbSession.Where("sessionid=?", SessionId).Get(&sessionPersistence)
	if err != nil {
		return err
	}
	if result {
		sessionPersistence.Sessionvalue = SessionValue
		_, err = dbSession.Where("sessionid=?", SessionId).Update(&sessionPersistence)
		return err
	} else {
		return errors.New("not find session")
	}
}

func (t *tblSessionPersistenceMgr) DeleteSessionValue(dbSession *xorm.Session, SessionId string) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var sessionPersistence tblSessionPersistence
	sessionPersistence.Sessionid = SessionId
	_, err := dbSession.Delete(sessionPersistence)
	return err
}


