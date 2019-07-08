package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	logrus "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var AyameVersion = "19.02.1"

type AyameOptions struct {
	LogDir         string `yaml:"log_dir"`
	LogName        string `yaml:"log_name"`
	LogLevel       string `yaml:"log_level"`
	Addr           string `yaml:"addr"`
	Port           int    `yaml:"port"`
	OverWsPingPong bool   `yaml:"over_ws_ping_pong"`
	AuthWebhookUrl string `yaml:"auth_webhook_url"`
}

var (
	// start options
	Options *AyameOptions
	logger  *logrus.Logger
)

// Initialize
func init() {
	configFilePath := flag.String("c", "./config.yaml", "ayame の設定ファイルへのパス(yaml)")
	// yaml ファイルを読み込み
	buf, err := ioutil.ReadFile(*configFilePath)
	if err != nil {
		// 読み込めない場合 Fatal で終了
		log.Fatal("cannot open config file, err=", err)
	}
	// yaml をパース
	err = yaml.Unmarshal(buf, &Options)
	if err != nil {
		// パースに失敗した場合 Fatal で終了
		log.Fatal("cannot parse config file, err=", err)
	}
}

func main() {
	flag.Parse()
	args := flag.Args()
	// 引数の処理
	if len(args) > 0 {
		if args[0] == "version" {
			fmt.Printf("WebRTC Signaling Server Ayame version %s", AyameVersion)
			return
		}
	}
	logger = setupLogger()
	url := fmt.Sprintf("%s:%d", Options.Addr, Options.Port)
	logger.Infof("WebRTC Signaling Server Ayame. version=%s", AyameVersion)
	logger.Infof("running on http://%s (Press Ctrl+C quit)", url)
	hub := newHub()
	go hub.run()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./sample/"+r.URL.Path[1:])
	})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		wsHandler(hub, w, r)
	})
	logger.Fatal(http.ListenAndServe(url, nil))
}
