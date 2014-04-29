package beacons

import (
	"fmt"
	"testing"
	"time"
)

type testStream struct {
	c chan bool
}

func (s *testStream) Start() error {
	s.c = make(chan bool)
	return nil
}

func (s *testStream) Serve() error {
	<-s.c
	return nil
}

func (s *testStream) Close() error {
	close(s.c)
	return nil
}

func (s *testStream) Name() string {
	return "test"
}

type testErrStream struct {
	c chan bool
}

func (s *testErrStream) Start() error {
	s.c = make(chan bool)
	return fmt.Errorf("Testing Error")
}

func (s *testErrStream) Serve() error {
	<-s.c
	return nil
}

func (s *testErrStream) Close() error {
	close(s.c)
	return nil
}

func (s *testErrStream) Name() string {
	return "test"
}

func TestService(t *testing.T) {
	service := New()
	service.AddStream(new(testStream))
	go func() {
		if err := service.Serve(); err != nil {
			t.Error(err)
		}
	}()
	timer := time.NewTimer(time.Second)
	<-timer.C
	service.Close()
}

func TestErrService(t *testing.T) {
	service := New()
	service.AddStream(new(testErrStream))
	if err := service.Serve(); err == nil {
		t.Error(fmt.Errorf("An error should be raised."))
		return
	} else {
		t.Log(err)
	}
	timer := time.NewTimer(time.Second)
	<-timer.C
	service.Close()
}
