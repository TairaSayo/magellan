package main

import (
	. "github.com/Shnifer/magellan/commons"
	"github.com/Shnifer/magellan/input"
	"github.com/Shnifer/magellan/v2"
)

const maneurDetailK = 5

func (s *cosmoScene) updateShipControl(dt float64) {
	s.procCruise()
	s.procControlTurn(dt)
	s.procControlForward(dt)
}

func (s *cosmoScene) procCruise() {
	if input.Get("cruiseonoff") {
		s.cruiseOn = !s.cruiseOn
		if s.cruiseOn {
			s.cruiseInput = input.GetF("forward")
		}
	}
	if input.Get("maneurdetailonoff") {
		s.maneurDetail = !s.maneurDetail
	}
}

func (s *cosmoScene) procControlTurn(dt float64) {
	turnInput := input.GetF("turn")
	if s.maneurDetail {
		turnInput /= maneurDetailK
	}
	massK := 1000 / Data.CalcCurMass()
	var min, max float64
	switch {
	case s.maneurLevel >= 0:
		max = s.maneurLevel + Data.SP.Shunter.Turn_acc*massK/100*dt
		min = s.maneurLevel - Data.SP.Shunter.Turn_slow*massK/100*dt
	case s.maneurLevel < 0:
		max = s.maneurLevel + Data.SP.Shunter.Turn_slow*massK/100*dt
		min = s.maneurLevel - Data.SP.Shunter.Turn_acc*massK/100*dt
	}
	s.maneurLevel = Clamp(turnInput, min, max)
	Data.PilotData.Ship.AngVel = s.maneurLevel * Data.SP.Shunter.Turn_max
}

func (s *cosmoScene) procControlForward(dt float64) {
	var thrustInput float64
	if s.cruiseOn {
		thrustInput = s.cruiseInput
	} else {
		thrustInput = input.GetF("forward")
	}

	massK := 1000 / Data.CalcCurMass()
	_ = massK
	var min, max float64
	switch {
	case s.thrustLevel >= 0:
		max = s.thrustLevel + Data.SP.March_engine.Thrust_acc/100*dt
		min = s.thrustLevel - Data.SP.March_engine.Thrust_slow/100*dt
	case s.thrustLevel < 0:
		max = s.thrustLevel + Data.SP.March_engine.Reverse_slow/100*dt
		min = s.thrustLevel - Data.SP.March_engine.Reverse_acc/100*dt
	}
	if Data.SP.March_engine.Reverse_max == 0 && min < 0 {
		min = 0
	}
	s.thrustLevel = Clamp(thrustInput, min, max)

	var accel float64
	switch {
	case s.thrustLevel >= 0:
		accel = s.thrustLevel * Data.SP.March_engine.Thrust_max / Data.CalcCurMass()
	case s.thrustLevel < 0:
		accel = s.thrustLevel * Data.SP.March_engine.Reverse_max / Data.CalcCurMass()
	}
	accelV := v2.InDir(Data.PilotData.Ship.Ang).Mul(accel)

	Data.PilotData.ThrustVector = accelV.Mul(DEFVAL.ThrustVectorK)
	//to general gravity calc
	//Data.PilotData.Ship.Vel.DoAddMul(accelV, dt)
}
