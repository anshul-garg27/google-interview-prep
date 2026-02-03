package middleware

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"init.bulbul.tv/bulbul-backend/saas-gateway/config"
	"init.bulbul.tv/bulbul-backend/saas-gateway/header"
	"init.bulbul.tv/bulbul-backend/saas-gateway/util"
	"io"
	"io/ioutil"
	"strings"
	"time"
)

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func RequestLogger(config config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		gc := util.GatewayContextFromGinContext(c, config)

		buf, _ := ioutil.ReadAll(c.Request.Body)
		rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
		rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf)) //We have to create a new Buffer, because rdr1 will be read.

		body := readBody(rdr1)
		c.Request.Body = rdr2

		w := &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = w
		//Process request
		c.Next()

		if gc.Config.LOG_LEVEL < 3 {
			// log with response body
			gc.Logger.Error().Msg(fmt.Sprintf("%s - %s - [%s] \"%s %s %s %s %d %s \"%s\"\n%s\n%s\n%s\n",
				c.Request.Header.Get(header.RequestID),
				c.ClientIP(),
				time.Now().Format(time.RFC1123),
				c.Request.Method,
				path,
				c.Request.Header.Get(header.ApolloOpName),
				c.Request.Proto,
				c.Writer.Status(),
				"Response time: ", time.Now().Sub(start),
				c.Request.UserAgent(),
				c.Request.Header,
				c.Errors.ByType(gin.ErrorTypePrivate).String(),
			))
			gc.Logger.Printf("~~~~ QUERY (%s) ~~~\n%s\n", c.Request.Header.Get(header.RequestID), strings.ReplaceAll(query, "\\n", "\n"))
			gc.Logger.Printf("~~~~ QUERY BODY (%s)~~~\n%s", c.Request.Header.Get(header.RequestID), body)
			gc.Logger.Printf("~~~~ RESPONSE (%s) ~~~\n%s", c.Request.Header.Get(header.RequestID), w.body.String())
		} else {
			// no response body
			gc.Logger.Error().Msg(fmt.Sprintf("%s - %s - [%s] \"%s %s %s %s %d %s %s \"%s\"\n%s\n%s\n",
				c.Request.Header.Get(header.RequestID),
				c.ClientIP(),
				time.Now().Format(time.RFC1123),
				c.Request.Method,
				path,
				c.Request.Header.Get(header.ApolloOpName),
				c.Request.Proto,
				c.Writer.Status(),
				"Response time: ", time.Now().Sub(start),
				c.Request.UserAgent(),
				c.Request.Header,
				c.Errors.ByType(gin.ErrorTypePrivate).String(),
			))
			gc.Logger.Printf("~~~~ QUERY (%s) ~~~\n%s\n", c.Request.Header.Get(header.RequestID), strings.ReplaceAll(query, "\\n", "\n"))
			gc.Logger.Printf("~~~~ QUERY BODY (%s) ~~~\n%s\n", c.Request.Header.Get(header.RequestID), body)
		}

	}
}

func readBody(reader io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)

	s := buf.String()
	return s
}
