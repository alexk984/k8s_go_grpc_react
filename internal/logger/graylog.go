package logger

import (
	"os"

	"github.com/Graylog2/go-gelf/gelf"
	"github.com/sirupsen/logrus"
)

// SetupGraylogLogger настраивает логирование в Graylog
func SetupGraylogLogger(serviceName string) *logrus.Logger {
	log := logrus.New()

	// Получаем адрес Graylog из переменной окружения
	graylogAddr := os.Getenv("GRAYLOG_ADDR")
	if graylogAddr == "" {
		graylogAddr = "localhost:12201" // значение по умолчанию
	}

	// Создаем GELF writer
	gelfWriter, err := gelf.NewWriter(graylogAddr)
	if err != nil {
		log.WithError(err).Warn("Failed to create GELF writer, using stdout")
		return log
	}

	// Создаем hook для отправки логов в Graylog
	hook := &GraylogHook{
		Writer: gelfWriter,
		Extra: map[string]interface{}{
			"service": serviceName,
			"version": "1.0.0",
		},
	}

	log.AddHook(hook)
	log.SetFormatter(&logrus.JSONFormatter{})

	return log
}

// GraylogHook реализует logrus.Hook для отправки логов в Graylog
type GraylogHook struct {
	Writer *gelf.Writer
	Extra  map[string]interface{}
}

func (hook *GraylogHook) Fire(entry *logrus.Entry) error {
	// Создаем GELF сообщение
	gelfMsg := gelf.Message{
		Version:  "1.1",
		Host:     getHostname(),
		Short:    entry.Message,
		TimeUnix: float64(entry.Time.Unix()),
		Level:    gelfLevel(entry.Level),
		Extra:    make(map[string]interface{}),
	}

	// Добавляем дополнительные поля
	for k, v := range hook.Extra {
		gelfMsg.Extra["_"+k] = v
	}

	// Добавляем поля из entry
	for k, v := range entry.Data {
		gelfMsg.Extra["_"+k] = v
	}

	return hook.Writer.WriteMessage(&gelfMsg)
}

func (hook *GraylogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostname
}

func gelfLevel(level logrus.Level) int32 {
	switch level {
	case logrus.PanicLevel:
		return 0
	case logrus.FatalLevel:
		return 1
	case logrus.ErrorLevel:
		return 3
	case logrus.WarnLevel:
		return 4
	case logrus.InfoLevel:
		return 6
	case logrus.DebugLevel:
		return 7
	case logrus.TraceLevel:
		return 7
	default:
		return 6
	}
}
