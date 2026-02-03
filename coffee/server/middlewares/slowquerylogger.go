package middlewares

import (
	"coffee/helpers"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

func SlowQueryLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		customWriter := w.(*helpers.ResponseWriterWithContent)
		requestStartTime := time.Now()
		requestTimeThreshold := 5 * time.Second
		defer LogSlowQueries(r, customWriter.InputBodyCopy, requestStartTime, requestTimeThreshold)
		next.ServeHTTP(customWriter, r)
	})
}

var (
	once       sync.Once
	loggerInit *logrus.Logger
	file       *os.File
)

func GetSlowQueryLogger() *logrus.Logger {
	once.Do(func() {
		loggerInit = initSlowQueryLogger()
	})
	return loggerInit
}
func initSlowQueryLogger() *logrus.Logger {
	l := logrus.New()

	// Set the output to a log file
	currentDir, err := os.Getwd()
	if err != nil {
		l.Error(err)
	}
	currentTime := time.Now().UTC()
	dateFormat := "2006-01-02" // This format represents yyyy-mm-dd
	formattedDate := currentTime.Format(dateFormat)
	fileName := formattedDate + "-slow.log"
	logFilePath := filepath.Join(currentDir, "logs", fileName)

	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		l.Out = file
	} else {
		l.Info("Failed to log to file, using default stderr")
	}

	l.SetFormatter(&logrus.JSONFormatter{})

	return l
}
func LogSlowQueries(r *http.Request, requestBody []byte, requestStartTime time.Time, requestTimeThreshold time.Duration) {
	elapsedTime := time.Since(requestStartTime).Seconds()
	elapsedTimeDuration := time.Duration(elapsedTime * float64(time.Second))

	if elapsedTimeDuration > requestTimeThreshold {
		logger := GetSlowQueryLogger()
		headers := make(map[string][]string)
		for name, values := range r.Header {
			headers[name] = values
		}
		logger.WithFields(logrus.Fields{
			"Type":        "SLOW REQUEST",
			"Message":     "Took " + strconv.FormatFloat(float64(elapsedTime), 'f', -1, 64) + " seconds",
			"RequestType": r.Method,
			"RequestUrl":  r.URL.String(),
			"RequestBody": string(requestBody),
			"Headers":     headers,
		}).Warn()

		if file != nil {
			file.Close()
		}
	}
}
