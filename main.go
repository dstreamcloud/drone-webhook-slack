package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/drone/drone-go/plugin/webhook"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Secret string `envconfig:"DRONE_YAML_SECRET"`
	Addr   string `envconfig:"DRONE_PLUGIN_ADDR"`
}

type plugin struct {
	url string
}

func (p *plugin) Deliver(ctx context.Context, req *webhook.Request) error {
	if req.Event == webhook.EventBuild && req.Action == "completed" {
		buf := bytes.NewBuffer(nil)
		json.NewEncoder(buf).Encode(map[string]interface{}{
			"blocks": []map[string]interface{}{{
				"type": "section",
				"text": map[string]interface{}{
					"type": "mrkdwn",
					"text": fmt.Sprintf("%s/%s [Build %d](%s) %s", req.Repo.Namespace, req.Repo.Name, req.Build.ID, req.Build.Link, req.Build.Status),
				},
			}},
		})
		res, err := http.Post(p.url, "application/json", buf)
		if err != nil {
			return err
		}
		defer res.Body.Close()
		if res.StatusCode/100 != 2 {
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				return err
			}
			return errors.New(string(body))
		}
	}
	return nil
}

func main() {
	cfg := &Config{}
	envconfig.MustProcess("", cfg)
	p := &plugin{}
	handler := webhook.Handler(p, cfg.Secret, logrus.StandardLogger())
	http.Handle("/", handler)
	logrus.Fatal(http.ListenAndServe(cfg.Addr, nil))
}
