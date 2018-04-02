package main

import (
	. "github.com/Shnifer/magellan/commons"
	"sync"
)

type pilotData struct {
	mu sync.RWMutex

	//state data
	state  State
	bsp    CBSP
	galaxy CGalaxy

	//common data
	ship CShipPos
}

var Data pilotData

func (pd pilotData) getStateData(data []byte) chan struct{} {
	done := make(chan struct{})

	go func() {
		//anyway done, even with error
		defer close(done)

		//get state data
		md, err := CMapData{}.Decode(data)
		if err != nil {
			panic("Weird state data:")
		}

		pd.mu.Lock()
		if bspDat, ok := md[PARTSTATE_BSP]; ok {
			pd.bsp = CBSP{}.Decode([]byte(bspDat))
		} else {
			pd.bsp = CBSP{}
		}

		if galaxyDat, ok := md[PARTSTATE_Galaxy]; ok {
			pd.galaxy = CGalaxy{}.Decode([]byte(galaxyDat))
		} else {
			pd.galaxy = CGalaxy{}
		}

		initSceneState()

		pd.mu.Unlock()

	}()

	return done
}

func (pd pilotData) commonSend() []byte {
	pd.mu.RLock()
	defer pd.mu.RUnlock()

	dat := pd.ship.Encode()
	md := make(CMapData)
	md[PARTCOMMON_ShipPos] = dat
	res, err := md.Encode()
	if err != nil {
		panic("CommonSend " + err.Error())
	}
	return []byte(res)
}

func (pd pilotData) commonRecv(buf []byte) {
	pd.mu.Lock()
	defer pd.mu.Unlock()

	md, err := CMapData{}.Decode(buf)
	if err != nil {
		panic("pilotData.commonRecv Can't decode mapData " + err.Error())
	}

	if part, ok := md[PARTCOMMON_ShipPos]; ok {
		pd.ship = CShipPos{}.Decode([]byte(part))
	}
}
