package boost

import (
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root          = "boost"
	aName         = root + ".applicationName"
	banner        = root + ".banner"
	bannerEnabled = banner + ".enabled"
	phrase        = banner + ".phrase"
	fontName      = banner + ".fontName"
	color         = banner + ".color"
	strict        = banner + ".strict"
	prt           = root + ".print"
	cfg           = prt + ".config"
	maxLength     = cfg + ".maxLength"
	cfgEnabled    = cfg + ".enabled"
)

var applicationName = ""

func init() {
	config.Add(aName, "boost_application", "defines application name")
	config.Add(bannerEnabled, false, "enable/disable boost banner")
	config.Add(phrase, "boost", "banner phrase")
	config.Add(fontName, "standard", "banner font. see https://github.com/common-nighthawk/go-figure")
	config.Add(color, "white", "banner color")
	config.Add(strict, true, "sets banner strict")
	config.Add(cfgEnabled, false, "enable/disable print boost configs")
	config.Add(maxLength, 25, "defines value max length")
}

func ApplicationName() string {
	if applicationName == "" {
		applicationName = config.String(aName)
	}
	return applicationName
}
