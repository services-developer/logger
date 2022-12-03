package logger

import (
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

var (
	fatalTriggerPush = promauto.NewCounter(prometheus.CounterOpts{
		Name: "triggered_send_to_fatal_log_total",
		Help: "Total number of fatal level messages sent to the log",
	})
	errorTriggerPush = promauto.NewCounter(prometheus.CounterOpts{
		Name: "triggered_send_to_error_log_total",
		Help: "Total number of error level messages sent to the log",
	})
	warningTriggerPush = promauto.NewCounter(prometheus.CounterOpts{
		Name: "triggered_send_to_warning_log_total",
		Help: "Total number of warning level messages sent to the log",
	})
	infoTriggerPush = promauto.NewCounter(prometheus.CounterOpts{
		Name: "triggered_send_to_info_log_total",
		Help: "Total number of info level messages sent to the log",
	})
	debugTriggerPush = promauto.NewCounter(prometheus.CounterOpts{
		Name: "triggered_send_to_debug_log_total",
		Help: "Total number of debug level messages sent to the log",
	})
	traceTriggerPush = promauto.NewCounter(prometheus.CounterOpts{
		Name: "triggered_send_to_trace_log_total",
		Help: "Total number of trace level messages sent to the log",
	})
)

const (
	defaultLogLevel      = "2"
	defaultContainerName = "service-app-log"
)

var errorLevels = map[logrus.Level]string{
	logrus.PanicLevel: "panic",
	logrus.FatalLevel: "fatal",
	logrus.ErrorLevel: "error",
	logrus.WarnLevel:  "warning",
	logrus.InfoLevel:  "info",
	logrus.DebugLevel: "debug",
	logrus.TraceLevel: "trace",
}

func SendToPanicLog(message string) {
	fatalTriggerPush.Inc()
	pushLogger(message, logrus.PanicLevel)
	os.Exit(1)
}

func SendToFatalLog(message string) {
	fatalTriggerPush.Inc()
	pushLogger(message, logrus.FatalLevel)
	os.Exit(1)
}

func SendToErrorLog(message string) {
	errorTriggerPush.Inc()
	pushLogger(message, logrus.ErrorLevel)
}

func SendToWarningLog(message string) {
	warningTriggerPush.Inc()
	pushLogger(message, logrus.WarnLevel)
}

func SendToInfoLog(message string) {
	infoTriggerPush.Inc()
	pushLogger(message, logrus.InfoLevel)
}

func SendToDebugLog(message string) {
	debugTriggerPush.Inc()
	pushLogger(message, logrus.DebugLevel)
}

func SendToTraceLog(message string) {
	traceTriggerPush.Inc()
	pushLogger(message, logrus.TraceLevel)
}

// Запись сообщения в файл лога
func pushLogger(message string, currentLevel logrus.Level) {
	configLogLevel := os.Getenv("LOG_LEVEL")

	if len(configLogLevel) == 0 {
		configLogLevel = defaultLogLevel
	}

	levelValue, errLevel := strconv.Atoi(configLogLevel)
	var logLevel logrus.Level

	if errLevel != nil {
		log.Println(errLevel)
	} else {
		logLevel = logrus.Level(levelValue)
	}

	if currentLevel > logLevel {
		return
	}

	flag.Parse()
	logsFilePath := getLogFilePath()
	logFile, err := os.OpenFile(logsFilePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	logger := &logrus.Logger{
		Out:   logFile,
		Level: logrus.TraceLevel,
		Formatter: &easy.Formatter{
			TimestampFormat: "2006-01-02 15:04:05",
			LogFormat:       "[%time%] %msg%",
		},
	}

	levelMessage := errorLevels[currentLevel]
	logger.Printf("[%s] [%s] [%s] %s \n",
		getHostName(), getContainerName(), levelMessage, message)
}

// Получить имя контейнера
func getContainerName() string {
	containerName := os.Getenv("CONTAINER_NAME")

	if len(containerName) == 0 {
		containerName = defaultContainerName
	}

	return containerName
}

// Получить путь к файлу лога
func getLogFilePath() string {
	return fmt.Sprintf("./log/%s.log", getContainerName())
}

// Получить имя сервера
func getHostName() string {
	var hostName string
	hostNameFile, err := ioutil.ReadFile("/etc/hostname")
	if err != nil {
		serverName, _ := os.Hostname()
		hostName = serverName
	} else {
		hostName = strings.ReplaceAll(string(hostNameFile), "\n", "")
	}

	return hostName
}
