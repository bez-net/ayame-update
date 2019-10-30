/*
	main function of ayame signaling server package
*/
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/ambelovsky/gosf"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var AyameVersion = "19.04.12"

type AyameOptions struct {
	LogDir         string `yaml:"log_dir"`
	LogName        string `yaml:"log_name"`
	LogLevel       string `yaml:"log_level"`
	Addr           string `yaml:"addr"`
	PortPlain      int    `yaml:"port_plain"`
	PortSecure     int    `yaml:"port_secure"`
	OverWsPingPong bool   `yaml:"over_ws_ping_pong"`
	AuthWebhookURL string `yaml:"auth_webhook_url"`
	AllowOrigin    string `yaml:"allow_origin"`
	MaxSessions    int    `yaml:"max_sessions"`
	FileDir        string
}

var (
	// start options
	Options *AyameOptions
	logger  *logrus.Logger
)

// initialization from config
func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	configFilePath := flag.String("c", "./config.yaml", "ayame configuration file path (yaml)")
	// yaml file read
	buf, err := ioutil.ReadFile(*configFilePath)
	if err != nil {
		log.Fatal("cannot open config file, err=", err)
	}
	// yaml data parse
	err = yaml.Unmarshal(buf, &Options)
	if err != nil {
		log.Fatal("cannot parse config file, err=", err)
	}
}

func main() {
	defer sendDebugStackToSlack()

	flag.Parse()
	args := flag.Args()
	// argument processing
	if len(args) > 0 {
		if args[0] == "version" {
			log.Printf("WebRTC Signaling Server Ayame version=%s", AyameVersion)
			return
		}
	}

	// NOTICE: I will not use logrus for readability
	// logger = setupLogger()

	// CAUTION: don't use localhost in url
	urlPlain := fmt.Sprintf(":%d", Options.PortPlain)
	urlSecure := fmt.Sprintf(":%d", Options.PortSecure)
	log.Printf("WebRTC Signaling Server, Ayame v%s", AyameVersion)
	log.Printf("running on http://<server>%s and https://<server>%s (Press Ctrl+C quit)", urlPlain, urlSecure)

	hub := newHub(fmt.Sprintf("Ayame v%s (%d)", AyameVersion, Options.MaxSessions))
	go hub.run()

	setupServerAPI(hub)

	// start servers for protocols supported
	go runPlainServer(urlPlain)
	go runSecureServer(urlSecure)
	go runSocketioServer(hub) // support ws and wss at the same time

	runSelfChecker(hub)
}

// Setting API endpoints for signalling
func setupServerAPI(hub *Hub) {
	// web file server for working a sample page
	// http.HandleFunc("/asset", func(w http.ResponseWriter, r *http.Request) {
	// 	http.ServeFile(w, r, r.URL.Path[1:])
	// })
	//http.Handle("/", http.FileServer(http.Dir("asset")))	// the simplest use
	fs := http.FileServer(http.Dir("asset"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Belows are API endpoints
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		// log.Printf("/hello")
		fmt.Fprintf(w, "{\"msg\":\"hello\"}")
	})
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		// log.Printf("/upload")
		uploadHandler(hub, w, r)
	})
	http.HandleFunc("/fetch", func(w http.ResponseWriter, r *http.Request) {
		// log.Printf("/upload")
		fetchHandler(hub, w, r)
	})
	// /ws endpoint is same with /signaling for compatibility
	http.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
		// log.Printf("/admin")
		adminHandler(hub, w, r)
	})
	http.HandleFunc("/event/", func(w http.ResponseWriter, r *http.Request) {
		// log.Printf("/event")
		eventHandler(hub, w, r)
	})
	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		// log.Printf("/chat")
		chatHandler(hub, w, r)
	})
	http.HandleFunc("/signal", func(w http.ResponseWriter, r *http.Request) {
		// log.Printf("/signal")
		signalHandler(hub, w, r)
	})
}

// Plain server supporint http and ws
func runPlainServer(url string) {
	timeout := 10 * time.Second
	server := &http.Server{Addr: url, Handler: nil, ReadHeaderTimeout: timeout}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

// Secure server supporting https and wss
func runSecureServer(url string) {
	timeout := 10 * time.Second
	server := &http.Server{Addr: url, Handler: nil, ReadHeaderTimeout: timeout}
	err := server.ListenAndServeTLS("certs/cert.pem", "certs/key.pem")
	if err != nil {
		log.Fatal(err)
	}
}

// Socket.io(gosf) server for plain and secure connections
func runSocketioServer(hub *Hub) {
	gosf.Listen("message", handleSignalMessage)
	gosf.Startup(map[string]interface{}{"port": 9999})
	log.Printf("socket.io closed")
}

func runSelfChecker(hub *Hub) {
	// chatting daemon function
	go procChatMessages()

	check := CheckConfig{level: 10.0, period: 10}
	disk := DiskUsage("/")
	summary := StringDiskUsage(disk)
	log.Println(summary)

	// checking the status
	for {
		event := &EventInfo{content: "event test"}
		hub.event <- event

		log.Printf("The service is now alive with %d min interval checking", check.period)
		time.Sleep(10 * time.Minute)
		disk := DiskUsage("/")
		CheckDiskWarning(disk, check.level)
		time.Sleep(time.Duration(check.period) * time.Minute)
	}
}
