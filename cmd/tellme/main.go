package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/yozel/tellme"
	"github.com/yozel/tellme/pkg/configure"
)

var (
	conf = &Configuration{}

	cmd  string
	args []string
)

type Configuration struct {
	TelegramToken  string `configure:"name:telegram-token;env:TELEGRAM_TOKEN;default:;usage:telegram bot token;required:true"`
	TelegramChatID int64  `configure:"name:telegram-chat-id;env:TELEGRAM_CHAT_ID;default:;usage:telegram chat id;required:true"`
}

// func main() {
// 	conf := &Configuration{}
// 	err := Parse(conf)
// 	if err != nil {
// 		panic(err)
// 	}
// })

func prepare() {
	err := configure.Parse(conf)
	if err != nil {
		panic(err)
	}

	// flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "missing required command\n")
		os.Exit(2)
	}
	cmd = flag.Arg(0)
	args = flag.Args()[1:]
}

func main() {
	prepare()

	factory := tellme.TeeCmdFactory{
		TelegramToken:  conf.TelegramToken,
		TelegramChatID: conf.TelegramChatID,
	}

	cmd := factory.NewTeeCmd(cmd, args...)

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		s := <-sigc
		cmd.Process.Signal(s)
	}()

	cmd.Run()
	err := cmd.SendNotification()
	if err != nil {
		log.Fatal(err)
	}
	if cmd.ProcessState == nil {
		log.Fatal("process state is nil")
	}
	os.Exit(cmd.ProcessState.ExitCode())
}
