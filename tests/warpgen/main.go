package main

import (
	"io/ioutil"
	"encoding/json"
	"image"
	"github.com/Shnifer/magellan/commons"
	"github.com/Shnifer/magellan/v2"
	"math/rand"
	"time"
	"strconv"
	"math"
	"golang.org/x/image/colornames"
)

const scale = 100

func init(){
	rand.Seed(time.Now().UnixNano())
	usedId = make(map[string]struct{},0)
}

func main(){
	buf,err:=ioutil.ReadFile("starpos.json")
	if err!=nil{
		panic(err)
	}
	var pts []image.Point
	err=json.Unmarshal(buf, &pts)
	if err!=nil{
		panic(err)
	}
	var gal commons.Galaxy
	gal.Points = make(map[string]*commons.GalaxyPoint)
	var flag bool
	for _,pt:=range pts{
		p:=commons.GalaxyPoint{
			Pos: pos(pt).Mul(scale),
			Type: commons.GPT_WARP,
			Size: 1,
			Mass: okr(1+rand.Float64()),
			WarpSpawnDistance: 15,
			WarpInDistance: 10,
			Color: colornames.White,
		}
		id:=genID()
		if !flag{
			flag=true
			id = "solar"
		}
		gal.Points[id] = &p
	}
	res,err:=json.Marshal(gal)
	ioutil.WriteFile("galaxy_warp.json", res, 0)
}

func pos(pt image.Point) v2.V2{
	v:= v2.V2{X: float64(pt.X), Y: float64(pt.Y)}.Add(v2.RandomInCircle(1))
	v.X=okr(v.X)
	v.Y=okr(v.Y)
	return v
}

var usedId map[string]struct{}
func genID() string{
	for {
		res:=randLetter()+randLetter()+strconv.Itoa(rand.Intn(10))
		if _,exist:=usedId[res]; !exist{
			usedId[res] = struct{}{}
			return res
		}
	}
}

func randLetter() string{
	n:=byte(rand.Intn(26))
	s:=[]byte("A")[0]
	return string([]byte{s+n})
}

func okr(x float64) float64 {
	const sgn = 100
	return math.Floor(x*sgn) / sgn
}