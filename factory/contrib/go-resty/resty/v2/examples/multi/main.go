package main

import (
	"context"
	logplugin "github.com/xgodev/boost/factory/contrib/go-resty/resty/v2/plugins/local/wrapper/log"
	"os"

	r "github.com/go-resty/resty/v2"
	"github.com/xgodev/boost/config"
	"github.com/xgodev/boost/factory/contrib/go-resty/resty/v2"
	ilog "github.com/xgodev/boost/factory/local/wrapper/log"
	"github.com/xgodev/boost/wrapper/log"
)

const (

	// config google client
	googleConfigPath          = "app.resty.googleConfigPath"
	googlePluginsConfigPath   = googleConfigPath + ".plugins"
	googleLogPluginConfigPath = googlePluginsConfigPath + ".log"

	// config bing client
	bingConfigPath          = "app.resty.bingConfigPath"
	bingPluginsConfigPath   = bingConfigPath + ".plugins"
	bingLogPluginConfigPath = bingPluginsConfigPath + ".log"
)

func init() {

	os.Setenv("BOOST_FACTORY_LOGRUS_CONSOLE_LEVEL", "INFO")

	os.Setenv("APP_RESTY_GOOGLE_HOST", "http://www.google.com")
	os.Setenv("APP_RESTY_ACOM_HOST", "https://www.bing.com.br")
	os.Setenv("APP_RESTY_ACOM_PLUGINS_LOG_LEVEL", "INFO")

	resty.ConfigAdd(bingConfigPath)
	logplugin.ConfigAdd(bingLogPluginConfigPath)

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

	var bingLogPlugin *logplugin.Log
	bingLogPlugin, err = logplugin.NewLogWithConfigPath(bingLogPluginConfigPath)
	if err != nil {
		log.Fatal(err)
	}

	var clientAcom *r.Client
	if clientAcom, err = resty.NewClientWithConfigPath(ctx, bingConfigPath, bingLogPlugin.Register); err != nil {
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
