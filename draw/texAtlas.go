package draw

import (
	"bytes"
	"encoding/json"
	"github.com/Shnifer/magellan/graph"
	. "github.com/Shnifer/magellan/log"
	"github.com/Shnifer/magellan/static"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"strconv"
)

type TexAtlasRec struct {
	FileName string
	Sx, Sy   int
	Count    int
	Smooth   bool
}

func (t TexAtlasRec) recKey() string {
	if t.Smooth {
		return "1~" + t.FileName
	} else {
		return "0~" + t.FileName
	}
}

type TexAtlas map[string]TexAtlasRec

const atlasFN = "atlas.json"
const defaultTexFN = "MAGIC_default.png"

var atlas TexAtlas

func InitTexAtlas() {
	saveAtlasExample("example_" + atlasFN)
	data, err := static.Load("textures", atlasFN)
	if err != nil {
		panic("Can't find tex atlas file " + atlasFN)
	}
	atlas = make(TexAtlas)
	err = json.Unmarshal(data, &atlas)
	if err != nil {
		panic(err)
	}
}

func atlasLoader(filename string) (io.Reader, error) {
	return static.Read("textures", filename)
}

//return tex and error
func getAtlasTexWithError(name string) (graph.Tex, error) {
	rec, ok := atlas[name]
	if !ok {
		return graph.Tex{}, errors.New("Not found atlas")
	}

	tex, err := graph.GetTex(rec.FileName, rec.Smooth, rec.Sx, rec.Sy, rec.Count, atlasLoader)
	if err != nil {
		return graph.Tex{}, err
	}
	return tex, nil
}

//return tex. if can't find return default tex
func getAtlasTex(name string) graph.Tex {
	tex, err := getAtlasTexWithError(name)
	if err != nil {
		Log(LVL_ERROR, "can't GetAtlasTex, name:", name)
		tex, err = graph.GetTex(defaultTexFN, false, 0, 0, 0, atlasLoader)
		if err != nil {
			panic(err)
		}
	}
	return tex
}

//return tex. if can't find return default tex
// panic on error
func GetAtlasTex(name string) graph.Tex {
	cacheKey := ""
	if rec, ok := atlas[name]; ok {
		cacheKey = "n~" + rec.recKey()
		if tex, exist := graph.CheckTexCache(cacheKey); exist {
			return tex
		}
	}
	tex := getAtlasTex(name)

	if cacheKey != "" {
		graph.StoreTexCache(cacheKey, tex)
	}

	return tex
}

func GetAtlasRoundTex(name string) graph.Tex {
	cacheKey := ""
	if rec, ok := atlas[name]; ok {
		cacheKey = "r~" + rec.recKey()
		if tex, exist := graph.CheckTexCache(cacheKey); exist {
			return tex
		}
	}
	tex := getAtlasTex(name)
	tex = graph.RoundTex(tex)

	if cacheKey != "" {
		graph.StoreTexCache(cacheKey, tex)
	}

	return tex
}

func GetSlidingAtlasTex(name string) graph.Tex {
	cacheKey := ""
	if rec, ok := atlas[name]; ok {
		cacheKey = "s~" + rec.recKey()
		if tex, exist := graph.CheckTexCache(cacheKey); exist {
			return tex
		}
	}
	tex := getAtlasTex(name)
	tex = graph.SlidingTex(tex)

	if cacheKey != "" {
		graph.StoreTexCache(cacheKey, tex)
	}

	return tex
}

func NewAtlasSprite(atlasName string, params graph.CamParams) *graph.Sprite {
	return graph.NewSprite(GetAtlasTex(atlasName), params)
}

func NewAtlasRoundSprite(atlasName string, params graph.CamParams) *graph.Sprite {
	return graph.NewSprite(GetAtlasRoundTex(atlasName), params)
}

func NewAtlasSpriteHUD(atlasName string) *graph.Sprite {
	return graph.NewSpriteHUD(GetAtlasTex(atlasName))
}

func NewAtlasFrame9HUD(atlasName string, w, h int, layer int) *graph.Frame9HUD {
	var sprites [9]*graph.Sprite
	for i := 0; i < 9; i++ {
		tex, err := getAtlasTexWithError(atlasName + strconv.Itoa(i))
		if err != nil {
			continue
		}
		sprites[i] = graph.NewSpriteHUD(tex)

	}
	return graph.NewFrame9(sprites, float64(w), float64(h), layer)
}

func saveAtlasExample(fn string) {
	exAtlas := make(map[string]TexAtlasRec)
	exAtlas["name"] = TexAtlasRec{
		FileName: "filename.png",
		Sx:       0,
		Sy:       0,
		Count:    1,
	}
	buf, err := json.Marshal(exAtlas)
	if err != nil {
		panic(err)
	}
	identbuf := bytes.Buffer{}
	json.Indent(&identbuf, buf, "", "  ")
	err = ioutil.WriteFile(fn, identbuf.Bytes(), 0)
	if err != nil {
		panic("can't write texture atlas example " + err.Error())
	}
}
