package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var config *viper.Viper
var logger *logrus.Logger
var server *http.Server
var instanceID string
var wg sync.WaitGroup

func init() {
	instanceID = uuid.New().String()

	config = viper.New()
	config.SetEnvPrefix("pinger")
	config.SetDefault("interval", time.Second*3)
	config.SetDefault("remote_url", "http://localhost:3001")
	config.SetDefault("request_body", fmt.Sprintf(`{"hello":"world","from":"pinger","id":"%s"}`, instanceID))
	config.SetDefault("request_Method", `get`)
	config.SetDefault("admin_addr", "0.0.0.0")
	config.SetDefault("admin_port", 3000)
	config.SetDefault("response_status_code", 200)
	config.AutomaticEnv()

	logger = logrus.New()
	logger.SetLevel(logrus.TraceLevel)
}

func main() {
	logger.Infof("starting pinger[%s] with:", instanceID)
	logger.Infof("\tinterval: %v", config.GetDuration("interval"))
	logger.Infof("\turl     : %s", config.GetString("remote_url"))
	logger.Infof("\tbody    : %s", config.GetString("request_body"))
	logger.Infof("\tmethod  : %s", config.GetString("request_method"))
	logger.Infof("\taddr    : %s", config.GetString("admin_addr"))
	logger.Infof("\tport    : %v", config.GetInt32("admin_port"))
	logger.Info("responding with:")
	logger.Infof("\tstatus  : %v", config.GetInt("response_status_code"))
	logger.Trace("\n")

	wg.Add(1)
	go startPingLoop(time.Tick(config.GetDuration("interval")))
	go startAdminServer()
	wg.Wait()
}

type PingerResponse struct {
	InstanceID string  `json:"instance_id"`
	NextHop    NextHop `json:"next_hop"`
}

type NextHop struct {
	Method string `json:"method"`
	URL    string `json:"url"`
	Body   string `json:"body"`
}

func createPingerResponse() []byte {
	response, err := json.Marshal(PingerResponse{
		InstanceID: instanceID,
		NextHop: NextHop{
			Method: strings.ToUpper(config.GetString("request_method")),
			URL:    config.GetString("remote_url"),
			Body:   config.GetString("request_body"),
		},
	})
	if err != nil {
		logger.Error(err)
		panic(err)
	}
	return response
}

func requestLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Infof("%s sent %s %s%s", r.RemoteAddr, r.Method, r.Host, r.URL)
		next.ServeHTTP(w, r)
	})
}

func startAdminServer() {
	bindAddr := fmt.Sprintf("%s:%v", config.GetString("admin_addr"), config.GetInt32("admin_port"))
	handler := mux.NewRouter()
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(config.GetInt("response_status_code"))
		w.Write(createPingerResponse())
	})
	server = &http.Server{
		Addr:    bindAddr,
		Handler: requestLoggerMiddleware(handler),
	}
	err := server.ListenAndServe()
	if err != nil {
		logger.Error(err)
		wg.Done()
	}
}

func startPingLoop(interval <-chan time.Time) {
	for {
		select {
		case <-interval:
			logger.Debugf("pinging '%s'...", config.GetString("remote_url"))
			requestBody := bytes.NewBuffer([]byte(config.GetString("request_body")))
			req, err := http.NewRequest(
				strings.ToUpper(config.GetString("request_method")),
				config.GetString("remote_url"),
				requestBody,
			)
			if err != nil {
				logger.Error(err)
				continue
			}
			client := http.Client{}
			res, err := client.Do(req)
			if err != nil {
				logger.Error(err)
			} else {
				switch res.StatusCode / 100 {
				case 2:
					logger.Info("2xx ok")
				case 3:
					logger.Info("3xx redirect")
				case 4:
					logger.Warn("4xx request error")
				case 5:
					logger.Warn("5xx remote error")
				default:
					logger.Debug(res)
				}
			}
		}
	}
}
