package beacons

import (
	"encoding/gob"
	"labix.org/v2/mgo/bson"
	"net"
)

func init() {
	streams["gob-up"] = NewGobUpStream()
	streams["gob-down"] = NewGobDownStream()
}

type GobUpStream struct {
	config   map[string]interface{}
	service  *Service
	listener net.Listener
}

func NewGobUpStream() Stream {
	return &GobUpStream{}
}

func (gs *GobUpStream) Init(config map[string]interface{}, service *Service) error {
	gs.config = config
	gs.service = service
	n, ok := config["net"].(string)
	if !ok {
		n = "tcp"
	}
	a, ok := config["addr"].(string)
	if ok {
		a = "127.0.0.1:5001"
	}
	var err error
	if gs.listener, err = net.Listen(n, a); err != nil {
		return err
	}
	return nil
}

func (gs *GobUpStream) Serve() error {
	for {
		conn, err := gs.listener.Accept()
		if err != nil {
			return err
		}
		go gs.newConn(conn)
	}
}

func (gs *GobUpStream) Close() error {
	return gs.listener.Close()
}

func (gs *GobUpStream) Write(e Entity) error {
	return e.Response()
}

func (gs *GobUpStream) newConn(conn net.Conn) {
	decoder := gob.NewDecoder(conn)
	encoder := gob.NewEncoder(conn)
	var e Entity
	for {
		if err := decoder.Decode(&e); err != nil {
			if fatal, _ := IsFatal(err); fatal {
				return
			} else {
				continue
			}
		}
		e.encoder = encoder
		go gs.service.Write(e)
	}
}

type GobDownStream struct {
	config  map[string]interface{}
	service *Service
	conn    net.Conn
	encoder *gob.Encoder
	decoder *gob.Decoder
}

func NewGobDownStream() Stream {
	return &GobDownStream{}
}

func (gs *GobDownStream) Init(config map[string]interface{}, service *Service) error {
	gs.config = config
	gs.service = service
	n, ok := config["net"].(string)
	if !ok {
		n = "tcp"
	}
	a, ok := config["addr"].(string)
	if ok {
		a = "127.0.0.1:5001"
	}
	var err error
	if gs.conn, err = net.Dial(n, a); err != nil {
		return err
	}
	gs.decoder = gob.NewDecoder(gs.conn)
	gs.encoder = gob.NewEncoder(gs.conn)
	return nil
}

func (gs *GobDownStream) Serve() error {
	var id bson.ObjectId
	for {
		if err := gs.decoder.Decode(&id); err != nil {
			if fatal, _ := IsFatal(err); fatal {
				return err
			} else {
				continue
			}
		}
		go gs.service.Done(id)
	}
}

func (gs *GobDownStream) Close() error {
	return gs.conn.Close()
}

func (gs *GobDownStream) Write(e Entity) error {
	return gs.encoder.Encode(e)
}
