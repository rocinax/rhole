package rhole

import (
	"net/http"
	"net/http/httputil"
	"path"
	"strings"
	"text/template"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Rhole Rocinax Rhole
type Rhole struct {
	Dump    bool
	Default Dummy
	Dummies []Dummy
}

// Dummy :
type Dummy struct {
	Name     string
	Path     string
	Method   string
	Template string
	Header   map[string]string
}

// ServeHTTP Custrum HTTP Handle Function
func (r Rhole) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	for i := 0; i < len(r.Dummies); i++ {

		if strings.HasPrefix(strings.ToUpper(req.RequestURI), strings.ToUpper(r.Dummies[i].Path)) {
			// output dump log
			logrus.WithFields(r.formatDumpLog(req)).Info(r.Dummies[i].Name)

			// response dummy text
			ResponseDummy(rw, req, r.Dummies[i])
			return
		}
	}

	// output dump log
	logrus.WithFields(r.formatDumpLog(req)).Info(r.Default.Name)
	ResponseDummy(rw, req, r.Default)
}

// ResponseDummy :
func ResponseDummy(rw http.ResponseWriter, req *http.Request, dummy Dummy) {
	// set advanced header set
	for key, value := range dummy.Header {
		rw.Header().Set(key, value)
	}

	// response data
	rw.WriteHeader(http.StatusOK)
	dummyTemplate := path.Join(viper.GetString("TemplateDir"), dummy.Template)
	templateHTML := template.Must(template.ParseFiles(dummyTemplate))
	templateHTML.Execute(rw, map[string]string{
		"ErrorDescription": "Contents Not Found",
		"Status":           "404",
		"Error":            "Not Found",
	})
	return
}

func (r Rhole) formatDumpLog(req *http.Request) map[string]interface{} {
	var dump string
	var remoteAddr string
	var remotePort string

	if r.Dump {
		dumpByte, _ := httputil.DumpRequest(req, true)
		dump = string(dumpByte)
	}

	if req.Header.Get("X-Read-IP") != "" {
		remoteAddr = req.Header.Get("X-Read-IP")
		remotePort = ""
	} else if req.Header.Get("X-Forwarded-For") != "" {
		remoteAddr = strings.Split(req.Header.Get("X-Forwarded-For"), ",")[0]
		remotePort = ""
	} else {
		remoteAddr = strings.Split(req.RemoteAddr, ":")[0]
		remotePort = strings.Split(req.RemoteAddr, ":")[1]
	}

	return logrus.Fields{
		"type":           "dump",
		"app":            "rhole",
		"request_id":     uuid.New().String(),
		"remote_address": remoteAddr,
		"remote_port":    remotePort,
		"host":           req.Host,
		"method":         req.Method,
		"request_uri":    req.RequestURI,
		"protocol":       req.Proto,
		"referer":        req.Referer(),
		"user_agent":     req.UserAgent(),
		"header":         req.Header,
		"body":           req.PostForm,
		"timestamp":      time.Now().Format(time.RFC3339),
		"_dump":          dump,
	}
}
