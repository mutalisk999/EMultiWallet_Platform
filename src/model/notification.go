package model

import (
	"sync"
	"time"
	"github.com/go-xorm/xorm"
)

type tblNotification struct {
	Notifyid     int       `xorm:"pk INTEGER autoincr"`
	Acctid       int       `xorm:"INT NOT NULL"`
	Wallettid    int       `xorm:"INT"`
	Trxid        int       `xorm:"INT"`
	Notifytype   int       `xorm:"INT NOT NULL"`
	Notification string    `xorm:"TEXT NOT NULL"`
	State        int       `xorm:"INT NOT NULL"`
	Reserved1    string    `xorm:"TEXT"`
	Reserved2    string    `xorm:"TEXT"`
	Createtime   time.Time `xorm:"created"`
	Updatetime   time.Time `xorm:"DATETIME"`
}

type tblNotificationMgr struct {
	TableName string
	Mutex     *sync.Mutex
}

func (t *tblNotificationMgr) Init() {
	t.TableName = "tbl_notification"
	t.Mutex = new(sync.Mutex)
}

func (t *tblNotificationMgr) ListNotifications(dbSession *xorm.Session, acctId int) ([]tblNotification, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	notifications := make([]tblNotification, 0)
	err := dbSession.Where("acctid = ?", acctId).And("state = ?", 0).
		Desc("createtime").Find(&notifications)
	if err != nil {
		return nil, err
	}
	return notifications, nil
}

func (t *tblNotificationMgr) GetNotificationCount(dbSession *xorm.Session, acctId int, state int) (int, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	var notification tblNotification
	total, err := dbSession.Where("acctid = ?", acctId).And("state = ?", state).Count(&notification)
	if err != nil {
		return 0, err
	}
	return int(total), nil
}

func (t *tblNotificationMgr) UpdateNotificationsState(dbSession *xorm.Session, notifyids []int, state int) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	var notification tblNotification
	notification.State = state
	_, err := dbSession.In("notifyid", notifyids).Cols("state").Update(&notification)
	if err != nil {
		return err
	}
	return nil
}

func (t *tblNotificationMgr) DeleteRegisterNotification(dbSession *xorm.Session, acctId int, notifyType int, reserved1 string) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	var notification tblNotification
	notification.Acctid = acctId
	notification.Notifytype = notifyType
	notification.Reserved1 = reserved1
	_, err := dbSession.Where("acctid=? and notifytype=? and Reserved1=?", acctId, notifyType, reserved1).Delete(&notification)
	if err != nil {
		return err
	}
	return nil
}

func (t *tblNotificationMgr) NewNotification(dbSession *xorm.Session, acctId *int, walletId *int, trxId *int,
	notifyType int, content string, state int, reserved1 string, reserved2 string) (int, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	var notification tblNotification
	if acctId != nil {
		notification.Acctid = *acctId
	}
	if walletId != nil {
		notification.Wallettid = *walletId
	}
	if trxId != nil {
		notification.Trxid = *trxId
	}
	notification.Notifytype = notifyType
	notification.Notification = content
	notification.State = state
	notification.Reserved1 = reserved1
	notification.Reserved2 = reserved2
	notification.Createtime = time.Now()
	notification.Updatetime = time.Now()
	_, err := dbSession.Insert(&notification)
	if err != nil {
		return 0, err
	}
	return notification.Notifyid, nil
}

func (t *tblNotificationMgr) DeleteNotification(dbSession *xorm.Session, notifyId *int, acctId *int, walletId *int, trxId *int,
	notifyType *int, state *int, reserved1 *string, reserved2 *string) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	var notification tblNotification
	if notifyId != nil {
		notification.Notifyid = *notifyId
	}
	if acctId != nil {
		notification.Acctid = *acctId
	}
	if walletId != nil {
		notification.Wallettid = *walletId
	}
	if trxId != nil {
		notification.Trxid = *trxId
	}
	if notifyType != nil {
		notification.Notifytype = *notifyType
	}
	if state != nil {
		notification.State = *state
	}
	if reserved1 != nil {
		notification.Reserved1 = *reserved1
	}
	if reserved2 != nil {
		notification.Reserved2 = *reserved2
	}
	_, err := dbSession.Delete(&notification)
	if err != nil {
		return err
	}
	return nil
}
