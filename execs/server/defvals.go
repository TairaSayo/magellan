package main

import (
	"bytes"
	"encoding/json"
	. "github.com/Shnifer/magellan/commons"
	"io/ioutil"
	"log"
)

const DefValPath = "./"
const roleName = "server"

type tDefVals struct {
	Port        string
	NeededRoles []string
	NodeName    string

	//in ms
	RoomUpdatePeriod      int
	SubscribeUpdatePeriod int
	LastSeenTimeout       int

	DoProf bool

	StartWarpSpeed         float64
	SolarStartLocationName string
}

var DEFVAL tDefVals

func setDefDef() {
	DEFVAL = tDefVals{
		Port: ":8000",
		NeededRoles: []string{
			ROLE_Pilot,
			ROLE_Navi,
			ROLE_Engi},
		SolarStartLocationName: "magellan",
		NodeName:               "storage01",
		SubscribeUpdatePeriod:  250,
	}
}

func init() {
	setDefDef()

	exfn := DefValPath + "example_ini_" + roleName + ".json"
	exbuf, err := json.Marshal(DEFVAL)
	identbuf := bytes.Buffer{}
	json.Indent(&identbuf, exbuf, "", "    ")
	if err := ioutil.WriteFile(exfn, identbuf.Bytes(), 0); err != nil {
		log.Println("can't even write ", exfn)
	}

	fn := DefValPath + "ini_" + roleName + ".json"

	buf, err := ioutil.ReadFile(fn)
	if err != nil {
		log.Println("cant read ", fn, "using default")
		return
	}
	json.Unmarshal(buf, &DEFVAL)
}
