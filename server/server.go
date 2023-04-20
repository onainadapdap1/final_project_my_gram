package server

import (
	"fmt"

	g "github.com/incubus8/go/pkg/gin"
	"github.com/onainadapdap1/dev/kode/my_gram/driver"
	"github.com/onainadapdap1/dev/kode/my_gram/router"
	"github.com/rs/zerolog/log"
	"github.com/subosito/gotenv"
)

func init() {
	gotenv.Load()
}


func StartServer() {
	addr := driver.Config.ServiceHost + ":" + driver.Config.ServicePort
	fmt.Println(driver.Config.ServicePort,"<<<<<")
	conf := g.Config {
		ListenAddr: addr,
		Handler: router.Router(),
		OnStarting: func () {
			log.Info().Msg("Your service is up and running at "+ addr)
		},
	}
	g.Run(conf)
}