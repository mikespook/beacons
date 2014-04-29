package beacons

import (
	"github.com/aarzilli/golua/lua"
	"github.com/mikespook/golib/iptpool"
	"github.com/mikespook/golib/log"
	"github.com/stevedonovan/luar"
	"path"
)

const (
	luaModule = "beacons"
)

type LuaIpt struct {
	module string
	state  *lua.State
	path   string
}

func newLuaIpt() iptpool.ScriptIpt {
	return &LuaIpt{}
}

func (luaipt *LuaIpt) Exec(name string, params interface{}) error {
	f := path.Join(luaipt.path, "beacons.lua")
	luaipt.Bind("Request", params)
	return luaipt.state.DoFile(f)
}

func (luaipt *LuaIpt) Init(path string) error {
	luaipt.state = luar.Init()
	luaipt.Bind("Debugf", log.Debugf)
	luaipt.Bind("Debug", log.Debug)
	luaipt.Bind("Messagef", log.Messagef)
	luaipt.Bind("Message", log.Message)
	luaipt.Bind("Warningf", log.Warningf)
	luaipt.Bind("Warning", log.Warning)
	luaipt.Bind("Errorf", log.Errorf)
	luaipt.Bind("Error", log.Error)
	luaipt.path = path
	return nil
}

func (luaipt *LuaIpt) Final() error {
	luaipt.state.Close()
	return nil
}

func (luaipt *LuaIpt) Bind(name string, item interface{}) error {
	luar.Register(luaipt.state, luaModule, luar.Map{
		name: item,
	})
	return nil
}
