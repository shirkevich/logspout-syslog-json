package json

import (
	"encoding/json"
	"errors"
	"log"
	"net"
	"os"

	"github.com/gliderlabs/logspout/router"
)

var hostname string

func init() {
	hostname, _ = os.Hostname()
	router.AdapterFactories.Register(NewJSONAdapter, "json")
}

// LogstashAdapter is an adapter that streams UDP JSON to Logstash.
type JSONAdapter struct {
	conn  net.Conn
	route *router.Route
}

// NewLogstashAdapter creates a LogstashAdapter with UDP as the default transport.
func NewJSONAdapter(route *router.Route) (router.LogAdapter, error) {
	transport, found := router.AdapterTransports.Lookup(route.AdapterTransport("udp"))
	if !found {
		return nil, errors.New("unable to find adapter: " + route.Adapter)
	}

	conn, err := transport.Dial(route.Address, route.Options)
	if err != nil {
		return nil, err
	}

	return &JSONAdapter{
		route: route,
		conn:  conn,
	}, nil
}

//func NewJSONMessage
func NewJSONMessage(m *router.Message) ([]byte, error) {
	msg := JSONMessage{
		Message: m.Data,
		Time:    uint(m.Time.Unix()),
		Source:  hostname,
		Docker: DockerInfo{
			Name:     m.Container.Name,
			ID:       m.Container.ID,
			Image:    m.Container.Config.Image,
			Hostname: m.Container.Config.Hostname,
			Labels:   m.Container.Config.Labels,
			Env:      EnvToMap(m.Container.Config.Env),
		},
	}
	js, err := json.Marshal(msg)
	if err != nil {
		log.Println("json:", err)
		return []byte{}, err
	}
	return js, nil
}

// Stream implements the router.LogAdapter interface.
func (a *JSONAdapter) Stream(logstream chan *router.Message) {
	for m := range logstream {
		msg, err := NewJSONMessage(m)
		if err != nil {
			log.Println("json:", err)
			continue
		}
		log.Println("output:", string(msg))
		_, err = a.conn.Write(msg)
		if err != nil {
			log.Println("json:", err)
			continue
		}
	}
}

func EnvToMap(envValues map) (map) {
	result := make(map[string]string, 0, len(envValues))
	for _, env := range envValues {
		s := strings.Split(env, "=")
		result[s[0]] = s[1]
	}
	return result
}

type DockerInfo struct {
	Name     string      `json:"name"`
	ID       string      `json:"id"`
	Image    string      `json:"image"`
	Hostname string      `json:"hostname"`
	Labels   interface{} `json:"labels"`
	Env      interface{} `json:"env"`
}

// LogstashMessage is a simple JSON input to Logstash.
type JSONMessage struct {
	Message string     `json:"message"`
	Time    uint       `json:"time"`
	Source  string     `json:"source"`
	Docker  DockerInfo `json:"docker"`
}
