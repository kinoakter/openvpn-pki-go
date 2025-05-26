package app

import (
	"context"
	"github.com/kinoakter/openvpn-pki-go/internal/api/handler"
	"github.com/kinoakter/openvpn-pki-go/internal/config"
	"github.com/kinoakter/openvpn-pki-go/internal/db/psql"
	"github.com/kinoakter/openvpn-pki-go/internal/db/psql/repository"
	"github.com/kinoakter/openvpn-pki-go/internal/domain/service"
	"github.com/kinoakter/openvpn-pki-go/log"
	"github.com/labstack/echo/v4"
	agentHttp "gitlab.com/vpn-tube/vpt-agent/pkg/http"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	ctx       context.Context
	stop      context.CancelFunc
	config    *config.Config
	apiServer *http.Server
}

func NewApp() *App {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	return &App{
		ctx:    ctx,
		stop:   stop,
		config: config.MustLoad(),
	}
}

func (a *App) Run() error {
	if err := a.initApp(); err != nil {
		return err
	}

	return nil
}

func (a *App) initApp() error {
	dbConn := psql.MustConnect(a.ctx, a.config.DatabaseURL)
	psql.Migrate(a.ctx, dbConn)

	// Repositories
	caRepository := repository.NewCaRepository(dbConn)
	serverCertRepository := repository.NewServerCertificateRepository(dbConn)
	clientCertRepository := repository.NewClientCertificateRepository(dbConn)

	// Services
	caService := service.NewCAService(a.ctx, caRepository)
	certificateService := service.NewServerCertificateService(a.ctx, serverCertRepository, caRepository)
	clientCertService := service.NewClientCertificateService(a.ctx, clientCertRepository, caRepository, serverCertRepository)
	ovpnService := service.NewOVPNService(a.ctx, caRepository, serverCertRepository, clientCertRepository)

	router := echo.New()
	openVpnRouter := router.Group("/api/v1/ovpn")
	handler.RegisterHealthyCheck(router)
	handler.NewCAHandler(caService).RegisterRoutes(openVpnRouter)
	handler.NewServerCertHandler(certificateService).RegisterRoutes(openVpnRouter)
	handler.NewClientCertHandler(clientCertService).RegisterRoutes(openVpnRouter)
	handler.NewOVPNHandler(ovpnService).RegisterRoutes(openVpnRouter)

	a.apiServer = agentHttp.StartHttpServer(a.config.VptAgent.Http, router)

	return nil
}

func (a *App) Stop() {
	a.stop()
}

func (a *App) WaitForExit() {
	<-a.ctx.Done()

	// Restore default behavior on the interrupt signal and notify a client of shutdown.
	a.stop()

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := a.apiServer.Shutdown(ctx); err != nil {
		log.Fatalf("API server forced to shutdown: %v", err)
	}

	select {
	case <-ctx.Done():
		log.Fatalf("vpt srv forced to shutdown")
	default:
		log.Infof("shut down gracefully")
	}
}
