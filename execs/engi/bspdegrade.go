package main

import (
	. "github.com/Shnifer/magellan/commons"
	"github.com/Shnifer/magellan/ranma"
)

func e(emi string) float64 {
	return 1 + Data.EngiData.Emissions[emi]*(DEFVAL.EmissionDegradePercent/100)
}

func getBoostPow(sysN int) float64 {
	res := 1.0
	for _, b := range Data.EngiData.Boosts {
		if b.LeftTime > 0 && b.SysN == sysN {
			res += b.Power / 100
		}
	}
	return res
}

func CalculateBSPDegrade(ranma *ranma.Ranma) (res BSPDegrade) {
	var b uint16
	var p float64
	k := func(x uint16) float64 {
		var n uint16
		for i := 0; i < 16; i++ {
			n += x & 1
			x = x >> 1
		}
		return float64(n)
	}
	f := func(v uint16, mask uint16) float64 {
		x := v & mask
		return k(x) / k(mask) * DEFVAL.RanmaMaxDegradePercent
	}
	d := func(v uint16, mask uint16) float64 {
		x := f(v, mask)
		if x > 1 {
			return 0
		} else {
			return 1 - x
		}
	}
	u := func(v uint16, mask uint16) float64 {
		return 1 + f(v, mask)
	}

	b = ranma.GetOut(SYS_MARCH)
	p = getBoostPow(SYS_MARCH)
	res.March_engine.Thrust_max = p * d(b, 36411) / e(EMI_ACCEL)
	res.March_engine.Thrust_acc = p * d(b, 18863)
	res.March_engine.Thrust_slow = p * d(b, 9590)
	res.March_engine.Reverse_max = p * d(b, 29125) / e(EMI_REVERSE)
	res.March_engine.Reverse_acc = p * d(b, 48721)
	res.March_engine.Reverse_slow = p * d(b, 56008)
	res.March_engine.Heat_prod = 1 / p * u(b, 5822) * e(EMI_ENGINE_HEAT)

	b = ranma.GetOut(SYS_WARP)
	p = getBoostPow(SYS_WARP)
	res.Warp_engine.Distort_max = p * d(b, 36411) / e(EMI_DIST_DOWN) * e(EMI_DIST_UP)
	res.Warp_engine.Warp_enter_consumption = 1 / p * u(b, 18863)
	res.Warp_engine.Distort_acc = p * d(b, 9590)
	res.Warp_engine.Distort_slow = p * d(b, 29125)
	res.Warp_engine.Consumption = 1 / p * u(b, 48721) * e(EMI_FUEL)
	res.Warp_engine.Turn_speed = p * d(b, 56008) / e(EMI_WARP_TURN)
	res.Warp_engine.Turn_consumption = 1 / p * u(b, 5822)

	b = ranma.GetOut(SYS_SHUNTER)
	p = getBoostPow(SYS_SHUNTER)
	res.Shunter.Turn_max = p * d(b, 36411) / e(EMI_TURN)
	res.Shunter.Turn_acc = p * d(b, 18863)
	res.Shunter.Turn_slow = p * d(b, 9590)
	res.Shunter.Strafe_max = p * d(b, 29125) / e(EMI_STRAFE)
	res.Shunter.Strafe_acc = p * d(b, 48721)
	res.Shunter.Strafe_slow = p * d(b, 56008)
	res.Shunter.Heat_prod = 1 / p * u(b, 5822) * e(EMI_ENGINE_HEAT)

	b = ranma.GetOut(SYS_RADAR)
	p = getBoostPow(SYS_RADAR)
	res.Radar.Range_Max = p * d(b, 36411) / e(EMI_RADAR_COSMOS) / e(EMI_RADAR_WARP)
	res.Radar.Angle_Min = 1 / p * u(b, 18863)
	res.Radar.Angle_Max = p * d(b, 9590)
	checkAngles(Data.BSP.Radar.Angle_Min, Data.BSP.Radar.Angle_Max, &res.Radar.Angle_Min, &res.Radar.Angle_Max)
	ak := e(EMI_RADAR_ANG_UP) / e(EMI_RADAR_ANG_DOWN)
	res.Radar.Angle_Min *= ak
	res.Radar.Angle_Max *= ak

	res.Radar.Angle_Change = p * d(b, 29125)
	res.Radar.Range_Change = p * d(b, 56008)
	res.Radar.Rotate_Speed = p * d(b, 5822)

	b = ranma.GetOut(SYS_SCANNER)
	p = getBoostPow(SYS_SCANNER)
	res.Scanner.ScanRange = p * d(b, 36410) / e(EMI_SCAN_RADIUS)
	res.Scanner.ScanSpeed = p * d(b, 18862) / e(EMI_SCAN_SPEED)
	res.Scanner.DropRange = p * d(b, 9590) / e(EMI_DROP_RADIUS)
	res.Scanner.DropSpeed = p * d(b, 4830) / e(EMI_DROP_SPEED)

	b = ranma.GetOut(SYS_FUEL)
	p = getBoostPow(SYS_FUEL)
	res.Fuel_tank.Fuel_Protection = p * d(b, 8095)
	res.Fuel_tank.Radiation_def = k(b & 58355) //K - Count of bits

	b = ranma.GetOut(SYS_LSS)
	p = getBoostPow(SYS_LSS)
	res.Lss.Thermal_def = p * d(b, 36411)
	res.Lss.Co2_level = k(b & 18862) //K - Count of bits
	res.Lss.Air_prepare_speed = p * d(b, 9590)
	res.Lss.Lightness = p * d(b, 4831)

	b = ranma.GetOut(SYS_SHIELD)
	p = getBoostPow(SYS_SHIELD)
	res.Shields.Radiation_def = p * d(b, 36411) / e(EMI_DEF_RADI)
	res.Shields.Disinfect_level = p * d(b, 18863)
	res.Shields.Mechanical_def = p * d(b, 9590) / e(EMI_DEF_MECH)
	res.Shields.Heat_reflection = p * d(b, 48721) / e(EMI_DEF_HEAT)
	res.Shields.Heat_capacity = p * d(b, 56008)
	res.Shields.Heat_sink = p * d(b, 5822)

	return res
}

func checkAngles(bMin, bMax float64, dMin, dMax *float64) {
	min := bMin * (*dMin)
	max := bMax * (*dMax)
	if min <= max {
		return
	}
	v := Clamp((min+max)/2, bMin, bMax)
	*dMin = v / bMin
	*dMax = v / bMax
}
