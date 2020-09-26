package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/worldOneo/muxcdn/host/cache"
)

const configFilePath = "muxcdn.conf.json"

// Config The default cdn config
type Config struct {
	DefaultFile string   `json:"defaultfile"`
	Addr        string   `json:"addr"`
	Workdir     string   `json:"workdir"`
	Whitelist   []string `json:"whitelist"`
	Blacklist   []string `json:"blacklist"`
	TLS         bool     `json:"TLS"`
	CertFile    string   `json:"certfile"`
	KeyFile     string   `json:"keyfile"`
	RecacheTime int64    `json:"recacheNano"`
}

// Server the server which contains every information needed
type Server struct {
	config     *Config
	workingDir string
	cache      *cache.Cache
}

var server *Server

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

	dir := cnf.Workdir
	dir += string(os.PathSeparator)

	c := cache.NewCache(cnf.RecacheTime)

	server = &Server{
		config:     cnf,
		workingDir: dir,
		cache:      c,
	}

	r := mux.NewRouter()

	r.HandleFunc("/{file:.+}", deliverContent)
	r.HandleFunc("/", deliverDefault)

	if server.config.TLS {
		fmt.Println("Starting Server with SSL!")
		go func() {
			err := http.ListenAndServeTLS(server.config.Addr, server.config.CertFile, server.config.KeyFile, r)
			fmt.Printf("An Error occured in the SSL Server: %v\n", err)
		}()
	} else {
		fmt.Println("Starting Server!")
		go func() {
			err := http.ListenAndServe(server.config.Addr, r)
			fmt.Printf("An error occured in the Server: %v\n", err)
		}()
	}

	fmt.Printf("Server is running on %s!\n", server.config.Addr)
	for true {
		fmt.Scanln()
		cnf, confErr = loadConfig()
		if confErr != nil {
			fmt.Println("Couldn't reload config!")
		} else if cnf != nil {
			server.config = cnf
			fmt.Println("Reloaded Config!")
		}
	}
}

func readContent(file string, c *cache.Cache) ([]byte, error) {
	return c.Get(file, func() ([]byte, error) {
		return ioutil.ReadFile(file)
	})
}

func logOut(s string, r *http.Request, in bool) {
	var f = "init"
	if in {
		f = "[%s] [%s]->local %s \n"
	} else {
		f = "[%s] local->[%s] %s \n"
	}
	fmt.Printf(f, time.Now().Format(time.UnixDate), r.RemoteAddr, s)
}

func loadConfig() (*Config, error) {
	var config Config
	configFile, err := os.Open(configFilePath)
	defer configFile.Close()
	if err != nil {
		config := Config{
			DefaultFile: "index.html",
			Addr:        "0.0.0.0:8080",
			Workdir:     "/var/www/html",
			Whitelist:   []string{"/var/www/html/*"},
			Blacklist:   []string{},
			TLS:         false,
			KeyFile:     "key.pem",
			CertFile:    "key.crt",
			RecacheTime: int64(time.Minute),
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

func deliverContent(w http.ResponseWriter, r *http.Request) {
	f := mux.Vars(r)["file"]
	logOut(fmt.Sprintf("REQ [%s]", f), r, true)

	indirName := server.workingDir + f
	allow := anyMatch(indirName, server.config.Whitelist)

	if allow {
		allow = !anyMatch(indirName, server.config.Blacklist)
	}

	if !allow {
		logOut(fmt.Sprintf("DEN [%v] 404", f), r, false)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 page not found"))
		return
	}

	err := sendFile(w, r, indirName)

	if err != nil {
		logOut(fmt.Sprintf("ERR [%v] 404", err), r, false)
		send404(w)
	}
}

func deliverDefault(w http.ResponseWriter, r *http.Request) {
	logOut(fmt.Sprintf("REQ [%s]", "/"), r, true)

	if server.config.DefaultFile == "" {
		send404(w)
		return
	}

	err := sendFile(w, r, server.config.DefaultFile)

	if err != nil {
		logOut(fmt.Sprintf("ERR [%v] 404", err), r, false)
		send404(w)
	}
}

func sendFile(w http.ResponseWriter, r *http.Request, file string) error {
	data, err := readContent(file, server.cache)
	if err != nil {
		return err
	}
	endings := strings.Split(file, ".")
	if len(endings) > 0 {
		ending := endings[len(endings)-1]
		contenttype, exists := TypeMap[ending]
		if exists {
			w.Header().Set("Content-Type", contenttype)
		} else {
			w.Header().Set("Content-Type", "text/html")
		}
	} else {
		w.Header().Set("Content-Type", "text/html")
	}
	w.WriteHeader(http.StatusOK)
	logOut(fmt.Sprintf("RES [%s] 200", file), r, false)
	w.Write(data)
	return nil
}

func send404(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 page not found"))
}

func anyMatch(path string, list []string) bool {
	for _, elem := range list {
		matched, _ := filepath.Match(elem, path)
		if matched {
			return true
		}
	}
	return false
}
