package main

import (
	"bytes"
	"encoding/json"
	. "github.com/Shnifer/magellan/commons"
	"github.com/Shnifer/magellan/network/storage"
	"github.com/Shnifer/magellan/static"
	"io/ioutil"
)

func (rd *roomServer) loadStateData(state State) (sd StateData, subscribe chan storage.Event) {
	defer LogFunc("loadStateData")()

	if state.ShipID != "" {
		sd.BSP = loadShipState(state.ShipID)
	}

	if state.GalaxyID != "" {
		sd.Galaxy = loadGalaxyState(state.GalaxyID)
		sd.Buildings, subscribe = loadBuildingsAndSubscribe(rd.storage, state.GalaxyID)
	}

	return sd, subscribe
}

//TODO: look in DB
func loadShipState(shipID string) *BSP {
	var res BSP

	buf, err := static.Load("DB", "bsp_"+shipID+".json")

	if err != nil {
		Log(LVL_ERROR, "Can't open file for ShipID ", shipID)
		return nil
	}
	err = json.Unmarshal(buf, &res)
	if err != nil {
		Log(LVL_ERROR, "can't unmarshal file for ship", shipID)
		return nil
	}
	return &res
}

//TODO: look in DB
func loadGalaxyState(GalaxyID string) *Galaxy {
	var res Galaxy
	buf, err := static.Load("DB", "galaxy_"+GalaxyID+".json")
	if err != nil {
		Log(LVL_ERROR, "Can't open file for galaxyID ", GalaxyID)
		return nil
	}

	err = json.Unmarshal(buf, &res)
	if err != nil {
		Log(LVL_ERROR, "can't unmarshal file for galaxy", GalaxyID)
		return nil
	}

	//First restore ID's
	for id, v := range res.Points {
		if v.ID == "" {
			v.ID = id
			res.Points[id] = v
		}
	}
	//Second - recalc lvls!
	res.RecalcLvls()

	return &res
}

func loadBuildingsAndSubscribe(storage *storage.Storage, GalaxyID string) (builds map[string]Building, subscribe chan storage.Event) {
	diskData, subscribe := storage.Subscribe(GalaxyID)
	builds = make(map[string]Building, len(diskData))
	for fullKey, data := range diskData {
		b, err := Building{}.Decode([]byte(data))
		if err != nil {
			Log(LVL_ERROR, "Wrong diskData", string(data))
			continue
		}
		b.FullKey = fullKey
		builds[fullKey] = b
	}
	return builds, subscribe
}

//save examples of DB data
func init() {
	//do not save it now
	//saveDataExamples(DBPath)
}

func saveDataExamples(path string) {
	bsp, _ := json.Marshal(BSP{})

	bufBsp := bytes.Buffer{}
	json.Indent(&bufBsp, bsp, "", "    ")
	ioutil.WriteFile(path+"example_bsp.json", bufBsp.Bytes(), 0)

	galaxy := Galaxy{Points: make(map[string]*GalaxyPoint)}
	galaxy.Points["samplePoint"] = &GalaxyPoint{
		ParentID:          "parentID",
		Orbit:             100,
		Period:            80,
		Mass:              10,
		WarpInDistance:    100,
		WarpSpawnDistance: 80,
		ScanData:          "Eurika!",
		Emissions:         []Emission{{Type: "Heat", MainRange: 100, FarValue: 200, MainValue: 100, FarRange: 200}},
	}
	galaxyStr, _ := json.Marshal(galaxy)
	//bufGalaxy := bytes.Buffer{}
	//json.Indent(&bufGalaxy, galaxyStr, "", "    ")
	ioutil.WriteFile(path+"example_galaxy.json", galaxyStr, 0)
}
