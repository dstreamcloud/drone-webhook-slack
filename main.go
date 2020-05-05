package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/drone/drone-go/plugin/webhook"
	"github.com/dstreamcloud/drone-webhook-slack/plugin"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Secret   string `envconfig:"DRONE_WEBHOOK_SECRET"`
	Addr     string `envconfig:"DRONE_PLUGIN_ADDR"`
	SlackURL string `envconfig:"DRONE_PLUGIN_SLACK_URL"`
}

func main() {
	cfg := &Config{}
	envconfig.MustProcess("", cfg)
	p := plugin.New(cfg.SlackURL)
	var handler http.Handler
	if os.Getenv("IS_DEVELOPMENT") == "1" {
		handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			req := &webhook.Request{}
			if err := json.NewDecoder(r.Body).Decode(req); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			err := p.Deliver(r.Context(), req)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write(nil)
		})
	} else {
		handler = webhook.Handler(p, cfg.Secret, logrus.StandardLogger())
	}
	http.Handle("/", handler)
	logrus.Fatal(http.ListenAndServe(cfg.Addr, nil))
}
