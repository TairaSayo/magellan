package commons

import (
	"encoding/json"
	"github.com/Shnifer/magellan/network/storage"
)

const (
	BUILDING_BLACKBOX = iota
	BUILDING_MINE
	BUILDING_BEACON
)

type Building struct {
	FullKey string

	Type int
	//where is it
	GalaxyID string
	//for mines
	PlanetID string
	//beckon and boxes are auto placed on far reach of system

	Message string
	//for mine
	OwnerID string
}

func (b Building) Encode() []byte {
	buf, err := json.Marshal(b)
	if err != nil {
		Log(LVL_ERROR, "can't marshal stateData", err)
		return nil
	}
	return buf
}

func (Building) Decode(buf []byte) (b Building, err error) {
	err = json.Unmarshal(buf, &b)
	if err != nil {
		return Building{}, err
	}
	return b, nil
}

func EventToCommand(e storage.Event) string {
	buf, err := json.Marshal(e)
	if err != nil {
		Log(LVL_ERROR, "can't marshal event", err)
		return ""
	}
	return CMD_BUILDING + string(buf)
}