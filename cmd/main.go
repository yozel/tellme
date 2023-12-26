package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/yozel/tellme"
)

var (
	telegramToken  string
	telegramChatID int64

	cmd  string
	args []string
)

func prepare() {
	flag.StringVar(&telegramToken, "telegram-token", "", "telegram bot token")
	flag.Int64Var(&telegramChatID, "telegram-chat-id", 0, "telegram chat id")
	required := []string{"telegram-token", "telegram-chat-id"}
	flag.Parse()

	seen := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) { seen[f.Name] = true })
	msgs := []string{}
	for _, req := range required {
		if !seen[req] {
			msgs = append(msgs, fmt.Sprintf("missing required -%s argument/flag\n", req))
		}
	}
	if len(msgs) > 0 {
		for _, msg := range msgs {
			fmt.Fprintf(os.Stderr, msg)
		}
		os.Exit(2)
	}

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
		TelegramToken:  telegramToken,
		TelegramChatID: telegramChatID,
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
