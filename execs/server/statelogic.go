package main

import (
	. "github.com/Shnifer/magellan/commons"
	"github.com/Shnifer/magellan/v2"
	"os"
	"time"
)

func generateCommonData(common CommonData, stateData StateData, newState, prevState State) CommonData {
	defer LogFunc("generateCommonData " + newState.StateID + " " + newState.GalaxyID + " " + newState.ShipID)()

	//DROP inter clients params
	common.PilotData.HeatProduction = 0
	common.NaviData.SonarDir = common.PilotData.Ship.Ang

	sessionTime := time.Now().Sub(StartDateTime).Seconds()
	common.PilotData.SessionTime = sessionTime
	stateData.Galaxy.Update(sessionTime)

	switch newState.StateID {
	case STATE_cosmo:
		if prevState.StateID == STATE_login {
			common.PilotData.Ship.Pos =
				stateData.Galaxy.Points[DEFVAL.SolarStartLocationName].Pos
		}

		if prevState.StateID == STATE_warp {
			common.PilotData.Ship.Pos =
				v2.InDir(180 + common.PilotData.Ship.Ang).Mul(stateData.Galaxy.SpawnDistance)
		}
	case STATE_warp:
		common = toWarpCommonData(common, stateData, newState, prevState)
	}

	return common
}

func toWarpCommonData(common CommonData, stateData StateData, newState, prevState State) CommonData {
	fromSystem := prevState.GalaxyID
	var systemPoint GalaxyPoint
	var found bool
	for _, v := range stateData.Galaxy.Points {
		if v.ID == fromSystem {
			systemPoint = v
			found = true
			break
		}
	}
	if !found {
		Log(LVL_ERROR, "toWarpCommonData: can't find system", fromSystem, "on warp map!")
		return common
	}

	pos := systemPoint.Pos
	ang := common.PilotData.Ship.Ang
	spawnRange := systemPoint.WarpSpawnDistance

	ship := RBData{
		Pos:    pos.AddMul(v2.InDir(ang), spawnRange),
		Vel:    v2.InDir(ang).Mul(DEFVAL.StartWarpSpeed),
		AngVel: 0,
		Ang:    ang,
	}

	common.PilotData.Ship = ship
	common.NaviData.SonarDir = ship.Ang

	return common
}

func (rd *roomServer) IsValidState(roomName string, stateStr string) bool {
	rd.mu.RLock()
	defer rd.mu.RUnlock()

	var res bool
	state := State{}.Decode(stateStr)
	switch state.StateID {
	case STATE_login:
		res = state.GalaxyID == "" && state.ShipID == ""
	case STATE_cosmo:
		res = rd.isValidFlyShip(roomName, state.ShipID) && rd.isValidFlyGalaxy(state.GalaxyID)
	case STATE_warp:
		res = rd.isValidFlyShip(roomName, state.ShipID) && rd.isValidFlyGalaxy(state.GalaxyID)
	}

	if !res {
		server.AddCommand(roomName, CMD_STATECHANGEFAIL)
	}

	return res
}

//run internal mutex call
func (rd *roomServer) isValidFlyShip(roomName string, shipID string) bool {
	if roomName == "" || shipID == "" {
		return false
	}

	for room, state := range rd.curState {
		if room != roomName && state.ShipID == shipID {
			return false
		}
	}

	if _, err := os.Stat(DBPath + "BSP_" + shipID + ".json"); os.IsNotExist(err) {
		return false
	}

	return true
}

//run internal mutex call
func (rd *roomServer) isValidFlyGalaxy(galaxyID string) bool {
	if _, err := os.Stat(DBPath + "Galaxy_" + galaxyID + ".json"); os.IsNotExist(err) {
		return false
	}

	return true
}
