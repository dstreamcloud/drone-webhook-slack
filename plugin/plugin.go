package plugin

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/drone/drone-go/plugin/webhook"
)

type Plugin struct {
	url string
}

func New(url string) *Plugin {
	return &Plugin{url: url}
}

func (p *Plugin) Deliver(ctx context.Context, req *webhook.Request) error {
	if req.Event == webhook.EventBuild && req.Action == "completed" {
		buf := bytes.NewBuffer(nil)
		json.NewEncoder(buf).Encode(map[string]interface{}{
			"blocks": []map[string]interface{}{{
				"type": "section",
				"text": map[string]interface{}{
					"type": "mrkdwn",
					"text": fmt.Sprintf("%s/%s <%s|Build %d> %s", req.Repo.Namespace, req.Repo.Name, req.Build.Link, req.Build.ID, req.Build.Status),
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
