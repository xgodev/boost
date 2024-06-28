package boost

import (
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root          = "boost"
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

func init() {

	config.Add(bannerEnabled, true, "enable/disable boost banner")
	config.Add(phrase, "boost", "banner phrase")
	config.Add(fontName, "standard", "banner font. see https://github.com/common-nighthawk/go-figure")
	config.Add(color, "white", "banner color")
	config.Add(strict, true, "sets banner strict")
	config.Add(cfgEnabled, true, "enable/disable print boost configs")
	config.Add(maxLength, 25, "defines value max length")
}
