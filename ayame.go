package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	logrus "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var AyameVersion = "19.08.0"

type AyameOptions struct {
	LogDir         string `yaml:"log_dir"`
	LogName        string `yaml:"log_name"`
	LogLevel       string `yaml:"log_level"`
	Addr           string `yaml:"addr"`
	Port           int    `yaml:"port"`
	OverWsPingPong bool   `yaml:"over_ws_ping_pong"`
	AuthWebhookURL string `yaml:"auth_webhook_url"`
	AllowOrigin    string `yaml:"allow_origin"`
}

var (
	// start options
	Options *AyameOptions
	logger  *logrus.Logger
)

// initialization
func init() {
	configFilePath := flag.String("c", "./config.yaml", "ayame configuration file path (yaml)")
	// yaml file set
	buf, err := ioutil.ReadFile(*configFilePath)
	if err != nil {
		log.Fatal("cannot open config file, err=", err)
	}
	// yaml parse
	err = yaml.Unmarshal(buf, &Options)
	if err != nil {
		log.Fatal("cannot parse config file, err=", err)
	}
}

func main() {
	flag.Parse()
	args := flag.Args()
	// argument processing
	if len(args) > 0 {
		if args[0] == "version" {
			fmt.Printf("WebRTC Signaling Server Ayame version %s\n", AyameVersion)
			return
		}
	}
	logger = setupLogger()
	url := fmt.Sprintf("%s:%d", Options.Addr, Options.Port)
	logger.Infof("WebRTC Signaling Server Ayame. version=%s", AyameVersion)
	logger.Infof("running on http://%s (Press Ctrl+C quit)", url)
	hub := newHub()
	go hub.run()

	// web file server for working sample service
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./sample/"+r.URL.Path[1:])
	})
	// /ws endpoint is same with /signaling for compatibility
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		signalingHandler(hub, w, r)
	})
	http.HandleFunc("/signaling", func(w http.ResponseWriter, r *http.Request) {
		signalingHandler(hub, w, r)
	})
	// timeout is 10 sec
	timeout := 10 * time.Second
	server := &http.Server{Addr: url, Handler: nil, ReadHeaderTimeout: timeout}
	err := server.ListenAndServe()
	if err != nil {
		logger.Fatal(err)
	}
}
