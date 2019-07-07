package src

import (
	"github.com/hr3lxphr6j/bililive-go/src/api"
	"github.com/juju/errors"
	"sync"
)

type IRecorderMngr interface {
	AddRecorder(*Recorder) error
	RemoveRecorder(id api.LiveId) error
	GetRecorder(id api.LiveId) *Recorder
	GetAllRecorder() []*Recorder
}

var iRecorderMngrInstance IRecorderMngr
var once sync.Once

func GetIRecorderMngr() IRecorderMngr {
	once.Do(func() {
		iRecorderMngrInstance = &recorderMngrImpl{
			recordersMap: make(map[api.LiveId]*Recorder),
		}
	})
	return iRecorderMngrInstance
}

type recorderMngrImpl struct {
	recordersMap map[api.LiveId]*Recorder
	rwlock       sync.RWMutex
	count        int
}

func (self *recorderMngrImpl) AddRecorder(rcd *Recorder) error {
	self.rwlock.Lock()
	defer self.rwlock.Unlock()
	if _, exist := self.recordersMap[rcd.Live.GetLiveId()]; exist {
		return errors.New("recorder already exist")
	} else {
		rcd.IndexID = self.count
		self.count++
		self.recordersMap[rcd.Live.GetLiveId()] = rcd
	}
	return nil
}

func (self *recorderMngrImpl) RemoveRecorder(id api.LiveId) (err error) {
	self.rwlock.Lock()
	defer self.rwlock.Unlock()
	if recorder, exist := self.recordersMap[id]; exist {
		recorder.Stop()
		delete(self.recordersMap, id)
	} else {
		err = errors.New("recorder not exist")
	}
	return
}

func (self *recorderMngrImpl) GetRecorder(id api.LiveId) *Recorder {
	self.rwlock.RLock()
	defer self.rwlock.RUnlock()
	if recorder, exist := self.recordersMap[id]; exist {
		return recorder
	}
	return nil
}

func (self *recorderMngrImpl) GetAllRecorder() []*Recorder {
	self.rwlock.RLock()
	defer self.rwlock.RUnlock()
	list := make([]*Recorder, 0)
	for _, v := range self.recordersMap {
		list = append(list, v)
	}
	return list
}
