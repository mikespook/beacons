package beacons

import (
	"fmt"
	"github.com/mikespook/golib/iptpool"
	"github.com/mikespook/golib/log"
	"labix.org/v2/mgo/bson"
	"strings"
	"sync"
)

var (
	streams           = make(map[string]Stream)
	ErrStreamNotFound = fmt.Errorf("Stream not found")
)

type Stream interface {
	Init(map[string]interface{}, *Service) error
	Serve() error
	Write(Entity) error
	Close() error
}

type Encoder interface {
	Encode(interface{}) error
}

type Service struct {
	iptPool     *iptpool.IptPool
	config      Config
	streams     map[string]Stream
	processChan chan bson.ObjectId
	data        map[bson.ObjectId]Entity

	sync.RWMutex
	wg sync.WaitGroup

	ErrorHandler func(error)
}

func New(config Config) (*Service, error) {
	service := &Service{
		iptPool:     iptpool.NewIptPool(newLuaIpt),
		processChan: make(chan bson.ObjectId),
		streams:     make(map[string]Stream),
		data:        make(map[bson.ObjectId]Entity),
	}
	service.config = config
	service.iptPool.OnCreate = func(ipt iptpool.ScriptIpt) error {
		ipt.Init(config.Script)
		ipt.Bind("Pass", service.Pass)
		return nil
	}
	for name, config := range service.config.Stream {
		if err := service.addStream(name, streams[name], config); err != nil {
			return nil, err
		}
	}
	return service, nil
}

func (s *Service) addStream(name string, stream Stream, config map[string]interface{}) error {
	s.streams[name] = stream
	return s.streams[name].Init(config, s)
}

func (s *Service) Serve() error {
	for name, stream := range s.streams {
		log.Messagef("The stream %s is starting.", name)
		s.wg.Add(1)
		go func(name string, stream Stream) {
			if err := stream.Serve(); err != nil {
				if !strings.Contains(err.Error(), "use of closed network connection") {
					s.err(err)
				}
			}
			s.wg.Done()
			log.Messagef("The stream %s is closed.", name)
		}(name, stream)
	}
	for id := range s.processChan {
		go func(id bson.ObjectId) {
			s.RLock()
			defer s.RUnlock()
			ipt := s.iptPool.Get()
			defer s.iptPool.Put(ipt)
			r, ok := s.data[id]
			if ok {
				if err := ipt.Exec("", r); err != nil {
					s.processChan <- id
					s.err(err)
				} else {
					s.Done(id)
				}
			}
		}(id)
	}
	return nil
}

func (s *Service) Pass(name string, e Entity) error {
	stream, ok := s.streams[name]
	if ok {
		return stream.Write(e)
	}
	return ErrStreamNotFound
}

func (s *Service) Write(e Entity) {
	s.Lock()
	defer s.Unlock()
	s.data[e.Id] = e
	s.processChan <- e.Id
	if err := e.Response(); err != nil {
		s.err(err)
	}
}

func (s *Service) Done(id bson.ObjectId) {
	s.Lock()
	defer s.Unlock()
	delete(s.data, id)
}

func (s *Service) Close() {
	close(s.processChan)
	for _, v := range s.streams {
		if err := v.Close(); err != nil {
			s.err(err)
		}
	}
	s.iptPool.Free()
	s.wg.Wait()
}

func (s *Service) err(err error) {
	if s.ErrorHandler != nil {
		s.ErrorHandler(err)
	}
}
