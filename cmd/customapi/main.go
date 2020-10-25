package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aidenwallis/customapi2/pkg/cache/redis"
	"github.com/aidenwallis/customapi2/pkg/config"
	"github.com/aidenwallis/customapi2/pkg/data"
	"github.com/aidenwallis/customapi2/pkg/server"
	"github.com/aidenwallis/customapi2/pkg/twitchapi"
	"github.com/aidenwallis/customapi2/pkg/version"
)

func main() {
	cfg := config.New()
	cacheImpl := redis.New(cfg.Redis)
	twitchAPI := twitchapi.New(cfg.Twitch)
	dataImpl := data.New(&data.Config{
		Cache:     cacheImpl,
		TwitchAPI: twitchAPI,
	})
	srv := server.New(cfg.ServerAddr, &version.Config{
		Cache:     cacheImpl,
		Data:      dataImpl,
		TwitchAPI: twitchAPI,
	})

	quitChan := make(chan bool, 1)

	go twitchAPI.Start()

	go func() {
		log.Println("Server started on " + cfg.ServerAddr)
		if err := srv.Start(); err != nil {
			select {
			case <-quitChan:
				// intentional close
				return

			default:
				log.Println("Failed to start server: " + err.Error())
			}
		}
	}()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	quitChan <- true

	if err := srv.Stop(ctx); err != nil {
		log.Println("Failed to stop server: " + err.Error())
	}

	if err := cacheImpl.Close(ctx); err != nil {
		log.Println("Failed to close cache client: " + err.Error())
	}

	log.Println("goodbye!")
}
