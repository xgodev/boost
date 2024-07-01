package main

import (
	"context"
	"github.com/xgodev/boost"
	log2 "github.com/xgodev/boost/factory/contrib/go-resty/resty/v2/plugins/local/wrapper/log"
	c "github.com/xgodev/boost/factory/contrib/spf13/cobra/v1"
	"os"

	r "github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
	"github.com/xgodev/boost/factory/contrib/go-resty/resty/v2"
	ilog "github.com/xgodev/boost/factory/local/wrapper/log"
	"github.com/xgodev/boost/wrapper/log"
)

const (

	// config google client
	google = "app.resty.google"

	// config bing client
	bing          = "app.resty.bing"
	bingPlugins   = bing + ".plugins"
	bingLogPlugin = bingPlugins + ".log"
)

func init() {

	os.Setenv("BOOST_FACTORY_LOGRUS_CONSOLE_LEVEL", "INFO")

	os.Setenv("APP_RESTY_GOOGLE_HOST", "http://www.google.com")
	os.Setenv("APP_RESTY_SITE_HOST", "https://www.bing.com")
	os.Setenv("APP_RESTY_SITE_PLUGINS_LOG_LEVEL", "INFO")

	resty.ConfigAdd(google)
	resty.ConfigAdd(bing)

	log2.ConfigAdd(bingLogPlugin)
}

func main() {

	boost.Start()

	ctx := context.Background()
	logger := log.FromContext(ctx)

	cmds := []*cobra.Command{
		{
			Use:  "google",
			Long: "google call",
			RunE: func(cmd *cobra.Command, args []string) error {
				return call(ctx, google)
			},
		},
		{
			Use:  "bing",
			Long: "bing call",
			RunE: func(cmd *cobra.Command, args []string) error {

				bingLogP, err := log2.NewLogWithConfigPath(bingLogPlugin)
				if err != nil {
					return err
				}

				return call(ctx, bing, bingLogP.Register)
			},
		},
	}

	rootCMD := &cobra.Command{
		Version: "1.0.0",
	}

	if err := c.Run(rootCMD, cmds...); err != nil {
		logger.Errorf(err.Error())
	}

	// go run main.go -> show options
	// go run main.go bing -> call bing
	// go run main.go google -> call google
}

func call(ctx context.Context, path string, plugins ...resty.Plugin) error {

	logger := log.FromContext(ctx)

	var err error

	var client *r.Client
	if client, err = resty.NewClientWithConfigPath(ctx, path, plugins...); err != nil {
		return err
	}

	var response *r.Response
	if response, err = client.R().Get("/"); err != nil {
		return err
	}

	if response != nil {
		logger.Infof(response.String())
	}

	return nil
}
