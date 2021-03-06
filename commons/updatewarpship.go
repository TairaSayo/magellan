package commons

import (
	. "github.com/Shnifer/magellan/log"
	"github.com/Shnifer/magellan/v2"
	"math"
)

func UpdateWarpAndShip(data TData, sumT float64, dt float64, powN float64) {
	if data.Galaxy == nil {
		Log(LVL_ERROR, "UpdateWarpAndShip called with nil Galaxy")
		return
	}
	if data.PilotData == nil {
		Log(LVL_ERROR, "UpdateWarpAndShip called with nil PilotData")
		return
	}

	sessionTime := data.PilotData.SessionTime
	distortion := data.PilotData.Distortion
	dir := data.PilotData.Dir

	//fast return, in fact we have to go out from warp
	if distortion == 0 {
		sessionTime += dt
		//final update for all and every object, slow but once
		data.Galaxy.Update(sessionTime)
		data.Galaxy.fixedTimeRest = 0
		data.PilotData.SessionTime = sessionTime
		return
	}

	galaxy := data.Galaxy
	ship := data.PilotData.Ship
	//warp update COPYPASTE warpPredictor
	gravK := math.Pow(distortion, powN)
	vel := VelDistWarpK * distortion

	var grav v2.V2
	ship.Vel = v2.InDir(dir).Mul(vel)
	sumT += galaxy.fixedTimeRest
	for sumT >= dt {
		sessionTime += dt
		sumT -= dt

		grav = SumWarpGravityAcc(ship.Pos, galaxy).Mul(gravK)
		ship.Vel.DoAddMul(grav, dt)
		ship.Vel = ship.Vel.Normed().Mul(vel)
		ship.Pos.DoAddMul(ship.Vel, dt)
		ship.Vel = ship.Vel.Rotate(data.PilotData.Ship.AngVel * dt)
	}
	ship.Ang = ship.Vel.Dir()
	//final update for all and every object, slow but once
	data.Galaxy.Update(sessionTime)

	data.PilotData.Ship = ship
	data.PilotData.Dir = ship.Vel.Dir()
	data.Galaxy.fixedTimeRest = sumT
	data.PilotData.SessionTime = sessionTime
}
