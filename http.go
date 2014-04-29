package beacons

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"labix.org/v2/mgo/bson"
	"net"
	"net/http"
	"time"
)

func init() {
	streams["http"] = NewHttpStream()
}

type HttpStream struct {
	config   map[string]interface{}
	service  *Service
	listener net.Listener
}

func NewHttpStream() Stream {
	return &HttpStream{}
}

func (hs *HttpStream) Init(config map[string]interface{}, service *Service) error {
	hs.config = config
	hs.service = service

	addr, ok := hs.config["addr"].(string)
	if !ok {
		addr = "127.0.0.1:5080"
	}
	tlsMap, _ := hs.config["tls"].(map[string]string)
	cert, _ := tlsMap["cert"]
	key, _ := tlsMap["key"]
	var err error
	hs.listener, err = net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	if cert != "" && key != "" {
		tlsConfig := &tls.Config{}
		tlsConfig.NextProtos = []string{"http/1.1"}
		tlsConfig.Certificates = make([]tls.Certificate, 1)
		tlsConfig.Certificates[0], err = tls.LoadX509KeyPair(cert, key)
		if err != nil {
			return err
		}
		hs.listener = tls.NewListener(hs.listener, tlsConfig)
	}
	http.HandleFunc("/", hs.handler)
	return nil
}

func (hs *HttpStream) Serve() error {
	return http.Serve(hs.listener, nil)
}

func (hs *HttpStream) Close() error {
	return hs.listener.Close()
}

func (hs *HttpStream) handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			writeJson(w, http.StatusBadRequest, err.Error())
			return
		}
		defer r.Body.Close()
		var e Entity
		if err := json.Unmarshal(data, &e); err != nil {
			writeJson(w, http.StatusBadRequest, err.Error())
			return
		}
		if e.Id == "" {
			e.Id = bson.NewObjectId()
		}
		if e.Time.IsZero() {
			e.Time = time.Now()
		}
		encoder := &httpEncoder{
			c: make(chan bson.ObjectId),
		}
		e.encoder = encoder
		go hs.service.Write(e)
		writeJson(w, http.StatusOK, encoder.Id())
		return
	}
	writeJson(w, http.StatusMethodNotAllowed, fmt.Errorf("%s", r.Method))
}

func (hs *HttpStream) Write(e Entity) error {
	return nil
}

type httpEncoder struct {
	c chan bson.ObjectId
}

func (encoder *httpEncoder) Encode(e interface{}) error {
	defer close(encoder.c)
	if resp, ok := e.(bson.ObjectId); ok {
		encoder.c <- resp
		return nil
	}
	return ErrTypeIssue
}

func (encoder *httpEncoder) Id() bson.ObjectId {
	return <-encoder.c
}

func writeJson(w http.ResponseWriter, status int, data interface{}) error {
	content, err := json.Marshal(data)
	if err != nil {
		return err
	}
	w.WriteHeader(status)
	_, err = w.Write(content)
	return err
}
