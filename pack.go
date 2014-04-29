package beacons

import (
	"encoding/gob"
	"fmt"
	"io"
	"labix.org/v2/mgo/bson"
	"net"
	"time"
)

var (
	ErrTypeIssue = fmt.Errorf("Type issue")
)

func init() {
	gob.Register(&Entity{})
}

type Entity struct {
	Id   bson.ObjectId
	Tag  map[string]bool
	Time time.Time
	Data string

	encoder Encoder
}

func (r *Entity) Response() error {
	if r.encoder == nil {
		return nil
	}
	return r.encoder.Encode(r.Id)
}

func IsFatal(err error) (fatal bool, e error) {
	if opErr, ok := err.(*net.OpError); ok { // is OpError
		fatal = opErr.Temporary() == false
		e = opErr
	} else { // isn't OpError
		fatal = true
		if err != io.EOF {
			e = err
		}
	}
	return
}
