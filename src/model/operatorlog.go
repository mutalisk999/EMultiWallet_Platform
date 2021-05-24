package model

import (
	"sync"
	"time"
	"github.com/go-xorm/xorm"
)

type tblOperatorLog struct {
	Logid      int       `xorm:"pk INTEGER autoincr"`
	Acctid     int       `xorm:"INT NOT NULL"`
	Optype     int       `xorm:"INT NOT NULL"`
	Content    string    `xorm:"TEXT NOT NULL"`
	Createtime time.Time `xorm:"created"`
}

type tblOperationLogMgr struct {
	TableName string
	Mutex     *sync.Mutex
}

func (t *tblOperationLogMgr) Init() {
	t.TableName = "tbl_operator_log"
	t.Mutex = new(sync.Mutex)
}

func (t *tblOperationLogMgr) GetOperatorLogs(dbSession *xorm.Session, acctId []int, opType []int, opTime [2]string, offSet int, limit int) (int, []tblOperatorLog, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	dbSession1 := dbSession.Where("")
	if acctId != nil && len(acctId) != 0 {
		dbSession1 = dbSession1.In("acctid", acctId)
	}
	if opType != nil && len(opType) != 0 {
		dbSession1 = dbSession1.In("optype", opType)
	}
	if opTime[0] != "" {
		dbSession1 = dbSession1.And("createtime > ?", opTime[0])
	}
	if opTime[1] != "" {
		dbSession1 = dbSession1.And("createtime < ?", opTime[1])
	}
	var log tblOperatorLog
	total, err := dbSession1.Count(&log)
	if err != nil {
		return 0, nil, err
	}

	dbSession2 := dbSession.Where("")
	if acctId != nil && len(acctId) != 0 {
		dbSession2 = dbSession2.In("acctid", acctId)
	}
	if opType != nil && len(opType) != 0 {
		dbSession2 = dbSession2.In("optype", opType)
	}
	if opTime[0] != "" {
		dbSession2 = dbSession2.And("createtime > ?", opTime[0])
	}
	if opTime[1] != "" {
		dbSession2 = dbSession2.And("createtime < ?", opTime[1])
	}
	opLogs := make([]tblOperatorLog, 0)
	dbSession2.Limit(limit, offSet).Desc("createtime").Find(&opLogs)
	return int(total), opLogs, nil
}

func (t *tblOperationLogMgr) NewOperatorLog(dbSession *xorm.Session, acctId int, opType int, content string) (int, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	var log tblOperatorLog
	log.Acctid = acctId
	log.Optype = opType
	log.Content = content
	log.Createtime = time.Now()
	_, err := dbSession.Insert(&log)
	if err != nil {
		return 0, err
	}
	return log.Logid, nil
}
