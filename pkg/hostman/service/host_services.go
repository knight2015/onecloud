package service

import (
	"context"
	"os"

	"yunion.io/x/log"
	"yunion.io/x/onecloud/pkg/appsrv"
	"yunion.io/x/onecloud/pkg/cloudcommon"
	"yunion.io/x/onecloud/pkg/cloudcommon/service"
	"yunion.io/x/onecloud/pkg/hostman/guestman"
	"yunion.io/x/onecloud/pkg/hostman/hostinfo"
	"yunion.io/x/onecloud/pkg/hostman/options"
)

type SHostService struct {
	service.SServiceBase
}

func (host *SHostService) StartService() {
	cloudcommon.ParseOptions(&options.HostOptions, &options.HostOptions.Options, os.Args, "host.conf")
	// Hostinfo.Init()
	// Firewall.Init()

	var c = make(chan struct{})
	cloudcommon.InitAuth(&options.HostOptions.Options, func() {
		log.Infof("Auth complete!!")
		// TODO
		hostinfo.Instance().StartRegister()
		guestman.Init(options.HostOptions.ServersPath)
		close(c)
	})
	app := cloudcommon.InitApp(&options.HostOptions.Options)
	host.TrapSignals(func() { host.quitSignalHandler(app) })
	host.InitHandlers(app)
	<-c // wait host info registered
	cloudcommon.ServeForever(app, &options.HostOptions)
}

func (host *SHostService) quitSignalHandler(app *appsrv.Application) {
	// TODO
	/* cloud/yunion/server/clouds/common/handler/__init__.py -> stop() */
	err := app.ShowDown(context.Background())
	if err != nil {
		log.Errorln(err.Error())
	}
	guestman.GetWorkManager().Stop()
}

func (host *SHostService) initHandlers(app *appsrv.Application) {
	guestman.AddGuestTaskHandler("", app)
}
