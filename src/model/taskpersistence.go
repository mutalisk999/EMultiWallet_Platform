package model

import (
	"sync"
	"time"
	"github.com/go-xorm/xorm"
)

type tblTaskPersistence struct {
	Id            int    `xorm:"pk INTEGER autoincr"`
	Taskuuid      string `xorm:"VARCHAR(64) NOT NULL UNIQUE index"`
	Walletuuid    string `xorm:"VARCHAR(64)"`
	Trxuuid       string `xorm:"VARCHAR(64)"`
	Pushtype      int    `xorm:"INTEGER NOT NULL"`
	State         int    `xorm:"INTEGER NOT NULL"`
	Createtime    time.Time `xorm:"DATETIME"`
	Updatetime    time.Time `xorm:"DATETIME"`
}

type tblTaskPersistenceMgr struct {
	TableName string
	Mutex     *sync.Mutex
}

func (t *tblTaskPersistenceMgr) Init() {
	t.TableName = "tbl_task_persistence"
	t.Mutex = new(sync.Mutex)
}

func (t *tblTaskPersistenceMgr) InsertTask(dbSession *xorm.Session, TaskUuid string, WalletUuid string, TrxUuid string, TaskType int, State int) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var taskPersistence tblTaskPersistence
	taskPersistence.Taskuuid = TaskUuid
	taskPersistence.Walletuuid = WalletUuid
	taskPersistence.Trxuuid = TrxUuid
	taskPersistence.Pushtype = TaskType
	taskPersistence.State = State
	taskPersistence.Createtime = time.Now()

	_, err := dbSession.Insert(&taskPersistence)
	return err
}

func (t *tblTaskPersistenceMgr) UpdateTaskState(dbSession *xorm.Session, TaskId int, State int) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var taskPersistence tblTaskPersistence
	taskPersistence.Id = TaskId
	taskPersistence.State = State
	taskPersistence.Updatetime = time.Now()

	_, err := dbSession.Where("id=?", TaskId).Cols("state", "updatetime").Update(&taskPersistence)
	return err
}

func (t *tblTaskPersistenceMgr) GetTaskByTaskUuid(dbSession *xorm.Session, TaskUuid string) (bool, tblTaskPersistence, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	var taskPersistence tblTaskPersistence
	isFound, err := dbSession.Where("taskuuid=?", TaskUuid).Get(&taskPersistence)
	return isFound, taskPersistence, err
}

func (t *tblTaskPersistenceMgr) GetTasksByState(dbSession *xorm.Session, State int) ([]tblTaskPersistence, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	taskPersistences := make([]tblTaskPersistence, 0)

	err := dbSession.Where("state=?", State).Find(&taskPersistences)
	return taskPersistences, err
}


