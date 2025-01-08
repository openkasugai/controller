package main

import (
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"
)

func init() {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	config.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.In(time.Local).Format("2006-01-02T15:04:05.000Z07:00"))
	}
	l, err := config.Build()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(l)
}

const configFilePath = "/config/config.yaml"

func main() {
	logger := zap.L()
	defer logger.Sync()

	waitForConfigFile := make(chan struct{})
	go func() {
		defer close(waitForConfigFile)
		for {
			files, err := filepath.Glob(configFilePath)
			if err != nil {
				logger.Fatal("config file list faild")
			}
			if 0 < len(files) {
				return
			}
			time.Sleep(3 * time.Second)
		}
	}()

	<-waitForConfigFile

	done := make(chan struct{})

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		s := <-sig
		defer close(done)
		logger.Fatal("Notify stop signal", zap.Any("Signal", s))
	}()

	reloadSig := make(chan os.Signal, 1)
	signal.Notify(reloadSig, syscall.SIGHUP)
	go func() {
		for {
			time.Sleep(3 * time.Second)

			buf, err := os.ReadFile(configFilePath)
			if err != nil {
				logger.Error("config read failed", zap.Error(err))
				continue
			}
			config := &struct {
				Command []string `yaml:"command"`
			}{}
			err = yaml.Unmarshal(buf, config)
			if err != nil {
				logger.Error("config read failed", zap.Error(err))
				continue
			}
			name := config.Command[0]
			args := []string{}
			if 1 < len(config.Command) {
				args = config.Command[1:]
			}
			cmd := exec.Command(name, args...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			err = cmd.Start()
			if err != nil {
				logger.Fatal("prosess start failed", zap.Any("cmd", err))
			}

			_, closed := <-reloadSig
			err = cmd.Process.Kill()
			if err != nil {
				logger.Error("prosess kill failed", zap.Any("prosess", cmd.Process))
			}
			if closed {
				break
			}
		}
	}()

	<-done
	close(reloadSig)
}
