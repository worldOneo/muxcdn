package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"

	"github.com/worldOneo/muxcdn/loadbalancer/host"
)

const configFilePath = "cdnlb.conf.json"

// LoadBalancer Holding everything
type LoadBalancer struct {
	config *Config
	hosts  []*host.Server
}

// Config for the proxy
type Config struct {
	Addr       string       `json:"addr"`
	Hosts      []HostConfig `json:"hosts"`
	TLS        bool         `json:"TLS"`
	CertFile   string       `json:"certfile"`
	KeyFile    string       `json:"keyfile"`
	HealthTime int          `json:"healthCheckSecs"`
}

// HostConfig a config for a specific host
type HostConfig struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

var (
	serverIndex = 0
	lb          *LoadBalancer
)

func main() {
	fmt.Println("Loading Config!")
	cnf, confErr := loadConfig()
	if cnf != nil && confErr != nil {
		fmt.Println("Config file created! Stoping server!")
		return
	}

	if cnf == nil && confErr != nil {
		fmt.Println("Couldn't create/load config! Stoping server!")
		return
	}

	hosts := loadHosts(cnf)
	lb = &LoadBalancer{
		config: cnf,
		hosts:  hosts,
	}

	runHealthChecks()
	r := mux.NewRouter()

	r.HandleFunc("/{f:.*}", serve)
	if cnf.TLS {
		fmt.Println("Starting Server with SSL!")
		go func() {
			err := http.ListenAndServeTLS(cnf.Addr, cnf.CertFile, cnf.KeyFile, r)
			fmt.Printf("An Error occured in the SSL Server: %v\n", err)
		}()
	} else {
		fmt.Println("Starting Server!")
		go func() {
			err := http.ListenAndServe(cnf.Addr, r)
			fmt.Printf("An error occured in the Server: %v\n", err)
		}()
	}

	for true {
		fmt.Scanln()
		cnf, confErr = loadConfig()
		if confErr != nil {
			fmt.Println("Couldn't reload config!")
		} else if cnf != nil {
			lb.config = cnf
			fmt.Println("Reloaded Config!")
		}
		lb.hosts = loadHosts(lb.config)
	}
}

func logOut(s string, r *http.Request, server *host.Server) {
	var f = "init"
	if server != nil {
		f = "[%s] [%s]-> local -> [%s](%s) %s \n"
		fmt.Printf(f, time.Now().Format(time.UnixDate), r.RemoteAddr, server.Name, server.URL, s)
	} else {
		f = "[%s] [%s]-> local -X %s \n"
		fmt.Printf(f, time.Now().Format(time.UnixDate), r.RemoteAddr, s)
	}
}

func loadConfig() (*Config, error) {
	var config Config
	configFile, err := os.Open(configFilePath)
	defer configFile.Close()
	if err != nil {
		config := Config{
			Addr: "0.0.0.0:8080",
			Hosts: []HostConfig{
				HostConfig{
					Name: "US-1",
					URL:  "http://123.123.123.123:8080/",
				},
			},
			TLS:        false,
			CertFile:   "cert.crt",
			KeyFile:    "key.pem",
			HealthTime: 2,
		}

		data, jsonErr := json.MarshalIndent(config, "", "    ")
		if jsonErr != nil {
			return nil, jsonErr
		}
		writingErr := ioutil.WriteFile(configFilePath, data, 0644)

		if writingErr != nil {
			return nil, writingErr
		}
		return &config, err
	}
	jsonParser := json.NewDecoder(configFile)
	parseErr := jsonParser.Decode(&config)

	if parseErr != nil {
		return nil, parseErr
	}
	return &config, nil
}

func serve(w http.ResponseWriter, r *http.Request) {
	server, err := getAvailableServer(lb.hosts)
	if err != nil {
		logOut(fmt.Sprintf("Failed [%v]", err), r, server)
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("503 No Service Available"))
		w.Header().Set("Content-Type", "text/plain")
		return
	}
	logOut("Connected", r, server)
	server.Handle(w, r)
}

func runHealthChecks() {
	go func() {
		for true {
			for _, e := range lb.hosts {
				e.Status.Online, _ = e.IsRunning()
			}
			time.Sleep(time.Duration(lb.config.HealthTime) * time.Second)
		}
	}()
}

func getAvailableServer(servers []*host.Server) (*host.Server, error) {
	for i := 0; i < len(servers); i++ {
		server := getNextServer(servers)
		if server.Status.Online {
			return server, nil
		}
	}
	return nil, fmt.Errorf("No Available server")
}

func getNextServer(servers []*host.Server) *host.Server {
	serverIndex = (serverIndex + 1) % len(servers)
	return servers[serverIndex]
}

func loadHosts(cnf *Config) []*host.Server {
	hosts := make([]*host.Server, 0)

	for _, h := range cnf.Hosts {
		server, err := host.New(h.Name, h.URL)
		if err != nil {
			fmt.Printf("Error occured while creating host [%v]\n", err)
			continue
		}
		fmt.Printf("Loaded host: [%s](%s)\n", server.Name, server.URL)
		hosts = append(hosts, server)
	}
	return hosts
}
