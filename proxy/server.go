package proxy

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"sync"
)

const eventsKey = "events"

type server struct {
	redis    *redis.Client
	pushLock *sync.Mutex
	popLock  *sync.Mutex
	maxItems int64
}

func (s *server) startServer(port int) error {
	http.HandleFunc("/push", func(w http.ResponseWriter, r *http.Request) {
		s.offerWebHookEvent(r)
		w.WriteHeader(204)
	})

	http.HandleFunc("/pop", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		event, err := s.pollWebHookEvent(r)

		if len(event) == 0 {
			w.WriteHeader(204)
		} else if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(map[string]string {
				"error": err.Error(),
			})
		} else {
			hash := sha1.Sum([]byte(event))
			w.WriteHeader(200)
			w.Header().Set("X-DRONE-WEBHOOK-SHA1", fmt.Sprintf("%x", hash))
			w.Write([]byte(event))
		}
	})

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	logrus.Infof("Listening for traffic on port %d.", port)
	return err
}

func (s *server) pollWebHookEvent(r *http.Request) (string, error) {
	reportError := func(msg string) {
		logrus.Errorf("Error dumping web hook events. %s", msg)
	}

	s.popLock.Lock()
	defer s.popLock.Unlock()

	event, err := s.redis.RPop(eventsKey).Result()
	if err != nil {
		if err == redis.Nil {
			logrus.Debug("No event.")
			return "", nil
		}

		reportError(err.Error())
		return "", err
	}

	if len(event) > 0 {
		logrus.Info("Polled one event.")
	} else {
		logrus.Debug("No event available.")
	}

	return event, nil
}

func(s *server) offerWebHookEvent(r *http.Request) {
	reportError := func(msg string) {
		logrus.Errorf("Error processing %s traffic from %s. %s", r.Method, r.Host, msg)
	}

	logrus.Infof("Receiving %s traffic from %s.\n", r.Method, r.Host)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		reportError(err.Error())
		return
	}

	eventBytes, err := json.Marshal(&Event{
		Method: 	r.Method,
		Headers: 	r.Header,
		Body:		string(body),
	})
	if err != nil {
		reportError(err.Error())
		return
	}

	s.pushLock.Lock()
	defer s.pushLock.Unlock()

	if _, err = s.redis.LPush(eventsKey, string(eventBytes)).Result(); err != nil {
		reportError(err.Error())
		return
	}

	logrus.Infof("Saved events from %s", r.Host)

	if itemsNum, err := s.redis.LLen(eventsKey).Result(); err != nil {
		reportError(err.Error())
		return
	} else if s.maxItems > 0 && itemsNum >= s.maxItems {
		s.redis.RPop(eventsKey)
	}
}

type Event struct {
	Method		string					`json:"0"`
	Headers		map[string][]string		`json:"1"`
	Body		string					`json:"2"`
}