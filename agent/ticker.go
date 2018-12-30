package agent

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"github.com/imulab/drone-webhook-proxy/proxy"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type ticker struct {
	interval 		int64
	proxyUrl		string
	droneHookUrl	string
}

func (t *ticker) start() {
	ticker := time.NewTicker(time.Duration(t.interval) * time.Second)
	go func() {
		for range ticker.C {
			event := t.getOneEvent()
			if event != nil {
				t.postOneEvent(event)
			}
		}
	}()
}

func (t *ticker) getOneEvent() *proxy.Event {
	logrus.Debug("Polling remote proxy for one event.")

	resp, err := http.Get(t.proxyUrl)
	if err != nil {
		logrus.Error("Remote server returned error.", err)
		return nil
	}

	if resp.StatusCode == 200 {
		logrus.Info("Remote server return a new event.")
		event := &proxy.Event{}
		if err := json.NewDecoder(resp.Body).Decode(event); err != nil {
			logrus.Errorf("Error decoding received event.", err)
		}
		return event
	} else {
		logrus.Debugf("Remote server returned status %d.", resp.StatusCode)
		return nil
	}
}

func (t *ticker) postOneEvent(event *proxy.Event) {
	req, err := http.NewRequest(event.Method, t.droneHookUrl, bytes.NewBuffer([]byte(event.Body)))
	if err != nil {
		logrus.Errorf("Error posting received event.", err)
		return
	}

	for k, vl := range event.Headers {
		for _, v := range vl {
			req.Header.Add(k, v)
		}
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		logrus.Errorf("Error posting received event.", err)
		return
	}

	logrus.Infof("Posting new event received status %d", resp.StatusCode)
}