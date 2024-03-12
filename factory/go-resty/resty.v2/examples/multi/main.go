package main

import (
	"context"
	"os"

	r "github.com/go-resty/resty/v2"
	"github.com/xgodev/boost/config"
	"github.com/xgodev/boost/factory/go-resty/resty.v2"
	logplugin "github.com/xgodev/boost/factory/go-resty/resty.v2/plugins/contrib/americanas-go/log.v1"
	ilog "github.com/xgodev/boost/factory/xgodev/boost.v1/log"
	"github.com/xgodev/boost/log"
)

const (

	// config google client
	googleConfigPath          = "app.resty.googleConfigPath"
	googlePluginsConfigPath   = googleConfigPath + ".plugins"
	googleLogPluginConfigPath = googlePluginsConfigPath + ".log"

	// config americanas client
	acomConfigPath          = "app.resty.acomConfigPath"
	acomPluginsConfigPath   = acomConfigPath + ".plugins"
	acomLogPluginConfigPath = acomPluginsConfigPath + ".log"
)

func init() {

	os.Setenv("IGNITE_LOGRUS_CONSOLE_LEVEL", "INFO")

	os.Setenv("APP_RESTY_GOOGLE_HOST", "http://www.googleConfigPath.com")
	os.Setenv("APP_RESTY_ACOM_HOST", "https://www.americanas.com.br")
	os.Setenv("APP_RESTY_ACOM_PLUGINS_LOG_LEVEL", "INFO")

	resty.ConfigAdd(acomConfigPath)
	logplugin.ConfigAdd(acomLogPluginConfigPath)

	resty.ConfigAdd(googleConfigPath)
	logplugin.ConfigAdd(googleLogPluginConfigPath)
}

func main() {

	config.Load()
	ilog.New()

	ctx := context.Background()
	logger := log.FromContext(ctx)

	var err error

	// ACOM CALL

	var acomLogPlugin *logplugin.Log
	acomLogPlugin, err = logplugin.NewLogWithConfigPath(acomLogPluginConfigPath)
	if err != nil {
		log.Fatal(err)
	}

	var clientAcom *r.Client
	if clientAcom, err = resty.NewClientWithConfigPath(ctx, acomConfigPath, acomLogPlugin.Register); err != nil {
		log.Fatal(err)
	}

	var responseAcom *r.Response
	if responseAcom, err = clientAcom.R().Get("/"); err != nil {
		log.Fatal(err)
	}

	if responseAcom != nil {
		logger.Infof(responseAcom.String())
	}

	// GOOGLE CALL

	var googleLogPlugin *logplugin.Log
	if googleLogPlugin, err = logplugin.NewLogWithConfigPath(googleLogPluginConfigPath); err != nil {
		log.Fatal(err)
	}

	var clientGoogle *r.Client
	if clientGoogle, err = resty.NewClientWithConfigPath(ctx, googleConfigPath, googleLogPlugin.Register); err != nil {
		log.Fatal(err)
	}

	var responseGoogle *r.Response
	if responseGoogle, err = clientGoogle.R().Get("/"); err != nil {
		log.Fatal(err)
	}

	if responseGoogle != nil {
		logger.Infof(responseGoogle.String())
	}

}
