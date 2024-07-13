package logger

import (
	"io"
	"log"
	"os"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	SugaredLogger *zap.SugaredLogger
	accessLogger  *zap.SugaredLogger
	panicLogger   *zap.SugaredLogger
	Logger        *zap.Logger
)

// 四个级别的日志句柄
var (
	FatalWriter  io.Writer
	AccessWriter io.Writer
	PanicWriter  io.Writer
	InfoWriter   io.Writer
	WarnWriter   io.Writer
)

var rotateMap = map[string]int{
	"Day":    1,
	"Hour":   1,
	"Minute": 1,
}

func InitLogger(levelRotate string, consoleLog, duration bool) {
	// 校验第一级文件夹上的时间
	// windows有些电脑砂锅面使用time.LoadLocation会失败，因为缺少一个文件，可以使用下面的方法替代
	timeLocal, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		timeLocal = time.FixedZone("CST", 8*3600)
	}
	// 日志按照格式分割，默认为hour级别
	if len(levelRotate) > 0 {
		if _, ok := rotateMap[levelRotate]; !ok {
			log.Fatal("LevelRotate error,levelRotate Must in Day|Hour|Minute")
		}
	}

	// 设置基本日志格式
	encoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		MessageKey:  "msg",
		LevelKey:    "level",
		EncodeLevel: zapcore.CapitalLevelEncoder,
		TimeKey:     "timestamp",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.In(timeLocal).Format("2006-01-02 15:04:05"))
		},
		CallerKey:    "caller",
		EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	})

	// access日志格式
	accessEncoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:     "timestamp",
		MessageKey:  "access",
		EncodeLevel: zapcore.CapitalLevelEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.In(timeLocal).Format("2006-01-02 15:04:05"))
		},
		EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	})

	// 实现三个判断日志等级的interface
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.WarnLevel
	})
	warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.WarnLevel
	})

	fatalLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.FatalLevel
	})

	accessLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return true

	})

	// 获取 info、warn，fatal,access日志文件的io.Writer 抽象
	InfoWriter = getWriter("logging/info", levelRotate)
	WarnWriter = getWriter("logging/error", levelRotate)
	FatalWriter = getWriter("logging/fatal", levelRotate)
	AccessWriter = getWriter("logging/access", "")
	PanicWriter = getWriter("logging/panic", "")

	// 创建具体的Logger
	// duration -- bool  是否需要持久化到磁盘中
	// var core = zapcore.NewTee()
	// var accessCore = zapcore.NewTee()
	// core := zapcore.NewTee()
	var core, accessCore, panicCore zapcore.Core
	if duration && consoleLog {
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(InfoWriter)), infoLevel),
			zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(WarnWriter)), warnLevel),
			zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(FatalWriter)), fatalLevel),
		)
		accessCore = zapcore.NewCore(accessEncoder, zapcore.NewMultiWriteSyncer( zapcore.AddSync(AccessWriter)), accessLevel)
		panicCore = zapcore.NewCore(accessEncoder, zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(PanicWriter)), accessLevel)
	} else if duration && !consoleLog {
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(zapcore.AddSync(InfoWriter)), infoLevel),
			zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(zapcore.AddSync(WarnWriter)), warnLevel),
			zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(zapcore.AddSync(FatalWriter)), fatalLevel),
		)
		accessCore = zapcore.NewCore(accessEncoder, zapcore.NewMultiWriteSyncer(zapcore.AddSync(AccessWriter)), accessLevel)
		panicCore = zapcore.NewCore(accessEncoder, zapcore.NewMultiWriteSyncer(zapcore.AddSync(PanicWriter)), accessLevel)
	} else {
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), infoLevel),
			zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), warnLevel),
			zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), fatalLevel),
		)

	}

	// 需要传入 zap.AddCaller() 才会显示打日志点的文件名和行数
	Logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	SugaredLogger = Logger.Sugar()

	AccessLogger := zap.New(accessCore, zap.AddCaller(), zap.AddCallerSkip(1))
	accessLogger = AccessLogger.Sugar()

	PanicLogger := zap.New(panicCore, zap.AddCaller(), zap.AddCallerSkip(1))
	panicLogger = PanicLogger.Sugar()
}

// 日志保存10天
func getWriter(filename string, rotate string) io.Writer {
	var fileName string
	var rotateTime time.Duration

	switch {
	case rotate == "Day":
		fileName = filename + "-%Y%m%d.log"
		rotateTime = time.Hour * 24
	case rotate == "Hour":
		fileName = filename + "-%Y%m%d%H.log"
		rotateTime = time.Hour
	case rotate == "Minute":
		fileName = filename + "-%Y%m%d%H%M.log"
		rotateTime = time.Minute
	default:
		fileName = filename + "-%Y%m%d%H.log"
		rotateTime = time.Hour
	}

	if strings.Contains(fileName, "fatal") {
		fileName = filename + ".log"
	}

	hook, err := rotatelogs.New(
		fileName,
		rotatelogs.WithMaxAge(time.Hour*10*24), // 日志保留10天
		rotatelogs.WithRotationTime(rotateTime),
	)

	if err != nil {
		Panic(err)
	}
	return hook
}

func Access(template string, args ...interface{}) {
	accessLogger.Infof(template, args...)
}

func PanicLog(template string, args ...interface{}) {
	panicLogger.Infof(template, args...)
}

func Debugf(template string, args ...interface{}) {
	SugaredLogger.Debugf(template, args...)
}

func Info(args ...interface{}) {
	SugaredLogger.Info(args...)
}

func Infof(template string, args ...interface{}) {
	SugaredLogger.Infof(template, args...)
}

func Warn(args ...interface{}) {
	SugaredLogger.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	SugaredLogger.Warnf(template, args...)
}

func Error(args ...interface{}) {
	SugaredLogger.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	SugaredLogger.Errorf(template, args...)
}

func Panic(args ...interface{}) {
	SugaredLogger.Error(args...)
	SugaredLogger.Panic(args...)
}

func Panicf(template string, args ...interface{}) {
	SugaredLogger.Errorf(template, args...)
	SugaredLogger.Panicf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	SugaredLogger.Errorf(template, args...)
	SugaredLogger.Fatalf(template, args...)
}

func Fatal(args ...interface{}) {
	SugaredLogger.Fatal(args...)
	SugaredLogger.Fatal(args...)
}
