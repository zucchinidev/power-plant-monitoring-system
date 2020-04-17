package www

import (
	"github.com/julienschmidt/httprouter"
	"github.com/zucchinidev/power-plant-monitoring-system/sensors/cmd/sensors/www/statusHandler"
	"github.com/zucchinidev/power-plant-monitoring-system/sensors/shared/ping"
	"net/http"
)

type Conf struct {
	Addr    string
	Version string
}

func Server(c Conf, pings []ping.Pinger) *http.Server {
	router := httprouter.New()
	router.GET("/status", statusHandler.Status(pings))
	return &http.Server{Addr: c.Addr, Handler: router}
}
