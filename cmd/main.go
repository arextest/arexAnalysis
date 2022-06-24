package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/arextest/arexAnalysis/arex"
	"github.com/gin-gonic/gin"

	_ "github.com/arextest/arexAnalysis/docs"

	"github.com/oklog/run"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           AREX Analysis
// @version         1.0
// @description     AREX Analysis to generate json-schema, to generate testcase auto.
// @termsOfService  http://swagger.io/terms/

// @contact.name   arexadmin@arex.com
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

func main() {
	serviceInit()
}

func serviceInit() {
	log.SetFormatter(&log.JSONFormatter{})

	var g run.Group
	{
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT, os.Interrupt)
		cancel := make(chan struct{})
		g.Add(
			//监控系统消息,优雅的退出
			func() error {
				select {
				case <-sigs:
					log.Warn("signal.Notify -> os.Exit(1)")
					os.Exit(1)
				case <-cancel:
					fmt.Println("Go routine 1 is closed")
					break
				}
				return nil
			},
			func(e error) {
				fmt.Println(e.Error())
				close(cancel)
			},
		)
	}
	{ // Prometheus monitor port 9090
		srv := &http.Server{Addr: ":9090", Handler: http.DefaultServeMux}
		http.Handle("/metrics", promhttp.Handler())
		g.Add(
			func() error {
				log.Info("Prometheus-Monitor-9090 start")
				err := srv.ListenAndServe()
				if err != http.ErrServerClosed {
					log.Errorf("Metrics Listen Failed. listen: %v", err)
					return err
				}
				return nil
			},
			func(e error) {
				fmt.Println("Prometheus-Monitor-9090 shutdown!")
				fmt.Println(e.Error())
				tempCtx, tcancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer tcancel()
				srv.Shutdown(tempCtx)
			},
		)
	}
	{ //  Start service of alert receiving server and port.
		engine := gin.Default()
		engine.Use(gin.Logger())
		engine.Use(gin.Recovery())
		// engine.Use(cors.Default())
		arex.InstallHandler(engine)
		engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

		srv := &http.Server{
			Addr:    ":8090",
			Handler: engine,
		}
		g.Add(
			func() error {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Error("Could not start listener")
					return err
				}
				return nil
			},
			func(err error) {
				tempCtx, tcancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer tcancel()
				srv.Shutdown(tempCtx)
				panic(err)
			},
		)
	}

	if err := g.Run(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
	log.Info("bigeyes terminated......................")
}
