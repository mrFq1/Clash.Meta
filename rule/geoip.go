package rules

import (
	"fmt"
	"strings"

	C "github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/log"
	"github.com/Dreamacro/clash/rule/geodata"
	"github.com/Dreamacro/clash/rule/geodata/router"
	_ "github.com/Dreamacro/clash/rule/geodata/standard"
)

type GEOIP struct {
	country      string
	adapter      string
	noResolveIP  bool
	network      C.NetWork
	geoIPMatcher *router.GeoIPMatcher
}

func (g *GEOIP) RuleType() C.RuleType {
	return C.GEOIP
}

func (g *GEOIP) Match(metadata *C.Metadata) bool {
	ip := metadata.DstIP
	if ip == nil {
		return false
	}
	return g.geoIPMatcher.Match(ip)
}

func (g *GEOIP) Adapter() string {
	return g.adapter
}

func (g *GEOIP) Payload() string {
	return g.country
}

func (g *GEOIP) ShouldResolveIP() bool {
	return !g.noResolveIP
}

func (g *GEOIP) NetWork() C.NetWork {
	return g.network
}

func NewGEOIP(country string, adapter string, noResolveIP bool, network C.NetWork) (*GEOIP, error) {
	geoLoaderName := "standard"
	//geoLoaderName := "memconservative"
	geoLoader, err := geodata.GetGeoDataLoader(geoLoaderName)
	if err != nil {
		return nil, fmt.Errorf("[GeoIP] %s", err.Error())
	}

	records, err := geoLoader.LoadGeoIP(strings.ReplaceAll(country, "!", ""))
	if err != nil {
		return nil, fmt.Errorf("[GeoIP] %s", err.Error())
	}

	geoIP := &router.GeoIP{
		CountryCode:  country,
		Cidr:         records,
		ReverseMatch: strings.Contains(country, "!"),
	}

	geoIPMatcher, err := router.NewGeoIPMatcher(geoIP)

	if err != nil {
		return nil, fmt.Errorf("[GeoIP] %s", err.Error())
	}

	log.Infoln("Start initial GeoIP rule %s => %s, records: %d", country, adapter, len(records))

	geoip := &GEOIP{
		country:      country,
		adapter:      adapter,
		noResolveIP:  noResolveIP,
		network:      network,
		geoIPMatcher: geoIPMatcher,
	}

	return geoip, nil
}