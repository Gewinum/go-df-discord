package main

import (
	"fmt"
	"github.com/Gewinum/go-df-discord/client"
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/pelletier/go-toml/v2"
	"github.com/sirupsen/logrus"
	"os"
)

var apiInstance *client.Api

type BindCommand struct {
}

func (c BindCommand) Run(source cmd.Source, output *cmd.Output) {
	plr, ok := source.(*player.Player)
	if !ok {
		output.Printf("You must run this command as a player")
		return
	}
	codeInfo, err := apiInstance.IssueCode(plr.XUID())
	if err != nil {
		output.Printf(err.Error())
		return
	}
	output.Printf("Your code is %s", codeInfo.Code)
}

func main() {
	api, err := client.NewApi("http://127.0.0.1:8080", "aaaa-bbb-cc")
	if err != nil {
		panic(err)
	}
	apiInstance = api

	cmd.Register(cmd.New("bind", "Bind your discord account to minecraft", make([]string, 0), BindCommand{}))

	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{ForceColors: true}
	logger.Level = logrus.DebugLevel

	chat.Global.Subscribe(chat.StdoutSubscriber{})

	conf, err := readConfig(logger)
	if err != nil {
		logger.Fatalln(err)
	}

	srv := conf.New()
	srv.CloseOnProgramEnd()

	srv.Listen()
	for srv.Accept(nil) {
	}
}

// readConfig reads the configuration from the config.toml file, or creates the
// file if it does not yet exist.
func readConfig(log server.Logger) (server.Config, error) {
	c := server.DefaultConfig()
	var zero server.Config
	if _, err := os.Stat("config.toml"); os.IsNotExist(err) {
		data, err := toml.Marshal(c)
		if err != nil {
			return zero, fmt.Errorf("encode default config: %v", err)
		}
		if err := os.WriteFile("config.toml", data, 0644); err != nil {
			return zero, fmt.Errorf("create default config: %v", err)
		}
		return c.Config(log)
	}
	data, err := os.ReadFile("config.toml")
	if err != nil {
		return zero, fmt.Errorf("read config: %v", err)
	}
	if err := toml.Unmarshal(data, &c); err != nil {
		return zero, fmt.Errorf("decode config: %v", err)
	}
	return c.Config(log)
}
