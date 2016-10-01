package lol

import "errors"

var (
	// ErrNoSuchRegion is returned if region not found.
	ErrNoSuchRegion = errors.New("no such region")
)

// Region represents a league of legends service area.
//
type Region int32

const (
	Global Region = 0

	PBE  Region = 2
	NA   Region = 10
	EUW  Region = 11
	EUNE Region = 12
	KR   Region = 13
	BR   Region = 14
	TR   Region = 15
	RU   Region = 16
	LAS  Region = 17
	LAN  Region = 18
	OCE  Region = 19
	JP   Region = 20
)

var regions = [...]Region{
	PBE,
	NA,
	EUW,
	EUNE,
	KR,
	BR,
	TR,
	RU,
	LAS,
	LAN,
	OCE,
	JP,
}

var regionByName = map[string]Region{
	"global": Global,
	"pbe":    PBE,
	"na":     NA,
	"euw":    EUW,
	"eune":   EUNE,
	"kr":     KR,
	"br":     BR,
	"tr":     TR,
	"ru":     RU,
	"las":    LAS,
	"lan":    LAN,
	"oce":    OCE,
	"jp":     JP,
}

var regionByPlatformID = map[string]Region{
	"pbe":  PBE,
	"na1":  NA,
	"euw1": EUW,
	"eun1": EUNE,
	"kr":   KR,
	"br1":  BR,
	"tr1":  TR,
	"ru":   RU,
	"la2":  LAS,
	"la1":  LAN,
	"oc1":  OCE,
	"jp1":  JP,
}

var nameByRegion = map[Region]string{
	Global: "global",
	PBE:    "pbe",
	NA:     "na",
	EUW:    "euw",
	EUNE:   "eune",
	KR:     "kr",
	BR:     "br",
	TR:     "tr",
	RU:     "ru",
	LAS:    "las",
	LAN:    "lan",
	OCE:    "oce",
	JP:     "jp",
}

var platformIDByRegion = map[Region]string{
	Global: "",
	PBE:    "pbe",
	NA:     "na1",
	EUW:    "euw1",
	EUNE:   "eun1",
	KR:     "kr",
	BR:     "br1",
	TR:     "tr1",
	RU:     "ru",
	LAS:    "la2",
	LAN:    "la1",
	OCE:    "oc1",
	JP:     "jp1",
}

var hostByRegion = map[Region]string{
	Global: "global.api.pvp.net",
	PBE:    "pbe.api.pvp.net",
	NA:     "na.api.pvp.net",
	EUW:    "euw.api.pvp.net",
	EUNE:   "eune.api.pvp.net",
	KR:     "kr.api.pvp.net",
	BR:     "br.api.pvp.net",
	TR:     "tr.api.pvp.net",
	RU:     "ru.api.pvp.net",
	LAS:    "las.api.pvp.net",
	LAN:    "lan.api.pvp.net",
	OCE:    "oce.api.pvp.net",
	JP:     "jp.api.pvp.net",
}

// RegionByName gets a region by name.
func RegionByName(name string) (Region, error) {
	r, ok := regionByName[name]
	if !ok {
		return 0, ErrNoSuchRegion
	}
	return r, nil
}

// RegionByPlatformID gets a region by platform id.
func RegionByPlatformID(id string) (Region, error) {
	r, ok := regionByPlatformID[id]
	if !ok {
		return 0, ErrNoSuchRegion
	}
	return r, nil
}

// Name returns name of the region.
func (r Region) Name() string {
	return nameByRegion[r]
}

// IsGlobal returns true if this region is global.
func (r Region) IsGlobal() bool {
	return r == Global
}

// String implements fmt.Stringer
func (r Region) String() string {
	return r.Name()
}

// PlatformID returns the ID for observer api.
func (r Region) PlatformID() string {
	return platformIDByRegion[r]
}

// Host returns hostname for api call.
func (r Region) Host() string {
	return hostByRegion[r]
}

func (r Region) baseURL() string {
	return "https://" + r.Host()
}

// Regions returns all region except 'Global'
func Regions() []Region {
	return regions[:]
}
