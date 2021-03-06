package main

import (
	"bytes"
	"encoding/json"
	. "github.com/Shnifer/magellan/commons"
	. "github.com/Shnifer/magellan/log"
	"io/ioutil"
)

const DefValPath = "./"
const roleName = "server"

type tDefVals struct {
	Port        string
	NeededRoles []string
	NodeName    string

	//in ms
	RoomUpdatePeriod       int
	SubscribeUpdatePeriod  int
	LastSeenTimeout        int
	OtherShipsUpdatePeriod int
	DiskRefreshPeriod      int

	DoProf    bool
	DebugPort string
	DebugMode bool

	StartWarpSpeed         float64
	SolarStartLocationName string
	MinDistortion          float64

	//in ms
	StoragePath     string
	LocalLogPath    string
	RestorePath     string
	RestoreTimeoutS int
	LogTimeoutMs    int
	LogRetryMinMs   int
	LogRetryMaxMs   int
	LogIP           string
	LogHostName     string

	//storage exchanger
	GameExchPort     string
	GameExchAddrs    []string
	GameExchPeriodMs int
	LogExchPort      string
	LogExchAddrs     []string
	LogExchPeriodMs  int

	WormHoleExitDistK float64
	WormHoleExitVel   float64

	NormPressure float64

	ShipsRequestHyServerAddr        string
	ShipDeadRequestHyServerAddr     string
	ShipReturnedRequestHyServerAddr string
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
		OtherShipsUpdatePeriod: 1000,
		LogTimeoutMs:           1000,
		LogRetryMinMs:          10,
		LogRetryMaxMs:          60000,
		DiskRefreshPeriod:      30,
	}
}

func init() {
	setDefDef()

	exfn := DefValPath + "example_ini_" + roleName + ".json"
	exbuf, err := json.Marshal(DEFVAL)
	identbuf := bytes.Buffer{}
	json.Indent(&identbuf, exbuf, "", "    ")
	if err := ioutil.WriteFile(exfn, identbuf.Bytes(), 0); err != nil {
		Log(LVL_WARN, "can't even write ", exfn)
	}

	fn := DefValPath + "ini_" + roleName + ".json"

	buf, err := ioutil.ReadFile(fn)
	if err != nil {
		Log(LVL_WARN, "cant read ", fn, "using default")
		return
	}
	json.Unmarshal(buf, &DEFVAL)
}
