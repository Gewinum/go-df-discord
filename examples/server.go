package main

import (
	"github.com/Gewinum/go-df-discord/server"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	sgn := make(chan os.Signal, 1)
	signal.Notify(sgn, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	srv := server.NewServer("aaaa-bbb-cc", "token", &server.Opts{
		Logger: logger,
	})
	srv.Bot().RegisterCommands("leave-empty-if-global")
	go func() {
		err := srv.ServeWeb(":8080")
		if err != nil {
			return
		}
	}()
	<-sgn
	logger.Info("Shutting the server down")
}
