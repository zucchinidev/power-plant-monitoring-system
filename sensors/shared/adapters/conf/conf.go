package conf

import (
	"github.com/joeshaw/envdecode"
	"github.com/zucchinidev/power-plant-monitoring-system/sensors/shared/adapters/logger"
	"sync"
)

type Conf struct {
	SQLDBUser      string `env:"POWER_PLANT_MONITORING_SYSTEM_SQL_DB_USER,required"`
	SQLDBPass      string `env:"POWER_PLANT_MONITORING_SYSTEM_SQL_DB_PASS,required"`
	SQLDBName      string `env:"POWER_PLANT_MONITORING_SYSTEM_SQL_DB_NAME,required"`
	SQLDBHost      string `env:"POWER_PLANT_MONITORING_SYSTEM_SQL_DB_HOST,required"`
	Addr           string `env:"POWER_PLANT_MONITORING_SYSTEM_HTTP_SERVER_LISTEN_ADDRESS,required"`
	BrokerQUrl     string `env:"POWER_PLANT_MONITORING_SYSTEM_BROKER_URL,required"`
	BrokerExchange string `env:"POWER_PLANT_MONITORING_SYSTEM_BROKER_EXCHANGE,required"`
}

var config *Conf
var once sync.Once

func init() {
	config = &Conf{}
}

func C() *Conf {
	once.Do(func() {
		if err := envdecode.Decode(config); err != nil {
			logger.New().FatalError(err)
		}
	})
	return config
}
