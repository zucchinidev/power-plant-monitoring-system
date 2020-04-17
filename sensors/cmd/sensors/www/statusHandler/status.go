package statusHandler

import (
	"github.com/julienschmidt/httprouter"
	"github.com/zucchinidev/power-plant-monitoring-system/sensors/cmd/sensors/www/engine"
	"github.com/zucchinidev/power-plant-monitoring-system/sensors/shared/ping"
	"net/http"
	"strings"
)

type statusResp struct {
	Status string `json:"status"`
	Msg    string `json:"msg"`
}

func Status(pings []ping.Pinger) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		var errs []string
		for _, p := range pings {
			if err := p.Ping(); err != nil {
				errs = append(errs, err.Error())
			}
		}
		if len(errs) > 0 {
			engine.Respond(w, r, http.StatusInternalServerError, statusResp{Status: "DOWN", Msg: strings.Join(errs, " - ")})
		} else {
			engine.Respond(w, r, http.StatusOK, statusResp{Status: "UP"})
		}
	}
}
