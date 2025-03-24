/* Copyright 2024 NTT Corporation , FUJITSU LIMITED */

package logging

import (
	// "fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/ini.v1"
	"log"
	"os"
	"time"
)

const (
	configName = "./config.ini"
	logFolder  = "./log/"
)

type ConfigList struct {
	LogFile string
}

var Config ConfigList

/* Initial settings for log file output */
func init() {
	cfg, err := ini.Load(configName)
	timeNow := time.Now()
	const layout = "20060102_150412_"

	if err != nil {
		log.Printf("Failed to read file: %v", err)
		os.Exit(1)
	}

	Config = ConfigList{
		LogFile: string(timeNow.Format(layout)) + cfg.Section("LogSettings").Key("log_file").String(),
	}
}

// Change the date and time display of the log
func JSTTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(time.RFC3339Nano))
}

// Log output settings
func SettingZapLogger() *zap.Logger {
	zc := zap.NewDevelopmentConfig()
	zc.OutputPaths = []string{"stdout", logFolder + Config.LogFile}
	zc.ErrorOutputPaths = zc.OutputPaths
	zc.DisableStacktrace = true
	zc.DisableCaller = false
	zc.EncoderConfig.EncodeTime = JSTTimeEncoder
	z, _ := zc.Build()
	return z
}
