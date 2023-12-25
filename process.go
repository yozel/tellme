package tellme

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"syscall"

	"github.com/yozel/tellme/internal/markdown"
)

type TeeCmdFactory struct {
	TelegramToken  string
	TelegramChatID int64
}

func (f *TeeCmdFactory) NewTeeCmd(name string, args ...string) *Cmd {
	t := exec.Command(name, args...)
	// var stdoutBuffer, stderrBuffer, stdallBuffer *bytes.Buffer
	var (
		stdoutBuffer = bytes.NewBuffer(nil)
		stderrBuffer = bytes.NewBuffer(nil)
		stdallBuffer = bytes.NewBuffer(nil)
	)
	t.Stdout = stdoutBuffer
	t.Stderr = stderrBuffer
	t.Stdout = io.MultiWriter(stdoutBuffer, stdallBuffer, os.Stdout)
	t.Stderr = io.MultiWriter(stderrBuffer, stdallBuffer, os.Stderr)
	return &Cmd{
		Cmd:            t,
		telegramToken:  f.TelegramToken,
		telegramChatID: f.TelegramChatID,
		Stdout:         stdoutBuffer,
		Stderr:         stderrBuffer,
		Stdall:         stdallBuffer,
	}
}

type Cmd struct {
	*exec.Cmd
	telegramToken  string
	telegramChatID int64
	Stdout         io.Reader
	Stderr         io.Reader
	Stdall         io.Reader
	Err            error
}

func (c *Cmd) RenderResult() (string, error) {
	if c.Cmd.ProcessState == nil {
		return "", fmt.Errorf("process state is nil")
	}
	status := c.Cmd.ProcessState.Sys().(syscall.WaitStatus)
	c.Cmd.ProcessState.Success()

	stdallByte, err := io.ReadAll(c.Stdall)
	if err != nil {
		return "", err
	}

	doc := &markdown.Document{}
	doc = doc.Normal("Command:\n").Code("sh", c.Cmd.String()).Normal("\n")
	doc = doc.Normal("Output:\n").Code("text", string(stdallByte)).Normal("\n")
	doc = doc.Normal("Exit code: ").InlineCode("%d", c.Cmd.ProcessState.ExitCode()).Normal("\n")
	doc = doc.Normal("Received signal: ").InlineCode("%d", status.Signal()).Normal("\n")
	return doc.String(), nil
}

func (c *Cmd) Run() {
	c.Err = c.Cmd.Run()
	log.Printf("Command: %s", c.Cmd.String())
	log.Printf("Exit code: %s", c.Err)
}

func (c *Cmd) SendNotification() error {
	if c.Err == nil {
		return nil
	}

	r, err := c.RenderResult()
	if err != nil {
		return err
	}
	err = SendNotification(c.telegramToken, c.telegramChatID, r)
	if err != nil {
		return err
	}

	return c.Err
}
