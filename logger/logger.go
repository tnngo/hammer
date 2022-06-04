package logger

import (
	"io"
	"os"
	"time"

	"github.com/tnngo/lad"
	"github.com/tnngo/lad/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger interface {
	Mode() zapcore.Core
}

type Level int8

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
)

func LadLevel(l Level) zapcore.Level {
	switch l {
	case DebugLevel:
		return zapcore.DebugLevel
	case InfoLevel:
		return zapcore.InfoLevel
	case WarnLevel:
		return zapcore.WarnLevel
	case ErrorLevel:
		return zapcore.ErrorLevel
	default:
		return zapcore.DebugLevel
	}
}

type Console struct {
}

func (c *Console) Mode() zapcore.Core {
	write := zapcore.AddSync(io.MultiWriter(os.Stdout))
	config := lad.NewProductionEncoderConfig()
	config.EncodeTime = timeFormat
	// 控制台输出颜色
	config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	// 定义日志核心
	return zapcore.NewCore(
		// 控制台
		zapcore.NewConsoleEncoder(config),
		write,
		LadLevel(DebugLevel),
	)
}

func (c *Console) Build() {
	lad.ReplaceGlobals(lad.New(c.Mode(), lad.AddCaller()))
}

type File struct {
	Level Level `json:"level"`
	// Filename 日志文件名称。
	Filename string `json:"filename"`
	// MaxSize 日志最大尺寸，默认为100MB。
	MaxSize int `json:"max_size" yaml:"max_size"`
	// MaxBackups 最大备份数量。
	MaxBackups int `json:"max_backups" yaml:"max_backups"`
	// MaxAge 最大保存时间。
	MaxAge int `json:"max_age" yaml:"max_age"`
	// Compress 是否压缩打包。
	Compress bool `json:"compress"`
}

func (f *File) Mode() zapcore.Core {
	hook := &lumberjack.Logger{
		Filename:   f.Filename,
		MaxSize:    f.MaxSize,
		MaxBackups: f.MaxBackups,
		MaxAge:     f.MaxAge,
		Compress:   f.Compress,
	}
	write := zapcore.AddSync(io.MultiWriter(hook))
	config := lad.NewProductionEncoderConfig()
	config.EncodeTime = timeFormat
	return zapcore.NewCore(
		zapcore.NewConsoleEncoder(config),
		write,
		LadLevel(f.Level),
	)
}

func (f *File) Build() {
	lad.ReplaceGlobals(lad.New(f.Mode(), lad.AddCaller()))
}

// 日志时间格式。
func timeFormat(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}
