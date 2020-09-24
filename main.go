package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"muxcdn/cache"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

const configFilePath = "muxcdn.conf.json"

// Config The default cdn config
type Config struct {
	Addr      string   `json:"addr"`
	Workdir   string   `json:"workdir"`
	Whitelist []string `json:"whitelist"`
	TLS       bool     `json:"TLS"`
	CertFile  string   `json:"certfile"`
	KeyFile   string   `json:"keyfile"`
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
		fmt.Println("Couldn't create config! Stoping server!")
		return
	}

	chDirErr := os.Chdir(cnf.Workdir)

	if chDirErr != nil {
		fmt.Println("Couldn't changedir into working dir! Stoping server!")
		return
	}

	dir, getDirErr := os.Getwd()
	dir += string(os.PathSeparator)

	if getDirErr != nil {
		fmt.Println("Couldn't get working directory! Stoping server!")
		return
	}

	c := cache.NewCache(int64(time.Minute))
	server = &Server{
		config:     cnf,
		workingDir: dir,
		cache:      c,
	}

	r := mux.NewRouter()

	r.HandleFunc("/{file:.+}", func(w http.ResponseWriter, r *http.Request) {
		f := mux.Vars(r)["file"]
		logOut(fmt.Sprintf("REQ [%s]", f), r, true)

		deny := true
		indirName := dir + f
		for _, elem := range cnf.Whitelist {
			matched, _ := filepath.Match(elem, indirName)
			if matched {
				deny = false
			}
		}

		if deny {
			logOut(fmt.Sprintf("DEN [%v] 404", f), r, false)
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 page not found"))
			return
		}

		data, err := readContent(f, c)
		if err == nil {
			endings := strings.Split(f, ".")
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
			logOut(fmt.Sprintf("RES [%s] 200", f), r, false)
			w.Write(data)
			return
		}

		logOut(fmt.Sprintf("ERR [%v] 404", err), r, false)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 page not found"))
	})

	go http.ListenAndServe(cnf.Addr, r)
	fmt.Printf("Server is running on %s!\n", cnf.Addr)
	fmt.Scanln()
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
			Addr:      "0.0.0.0:8080",
			Workdir:   ".",
			Whitelist: []string{"/var/www/html/*"},
			TLS:       false,
			KeyFile:   "key.pem",
			CertFile:  "key.crt",
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
	jsonParser.Decode(&config)
	return &config, nil
}

func deliverContent(w http.ResponseWriter, r *http.Request) {
	f := mux.Vars(r)["file"]
	logOut(fmt.Sprintf("REQ [%s]", f), r, true)

	deny := true
	indirName := server.workingDir + f
	for _, elem := range server.config.Whitelist {
		matched, _ := filepath.Match(elem, indirName)
		if matched {
			deny = false
		}
	}

	if deny {
		logOut(fmt.Sprintf("DEN [%v] 404", f), r, false)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 page not found"))
		return
	}

	data, err := readContent(f, server.cache)
	if err == nil {
		endings := strings.Split(f, ".")
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
		logOut(fmt.Sprintf("RES [%s] 200", f), r, false)
		w.Write(data)
		return
	}

	logOut(fmt.Sprintf("ERR [%v] 404", err), r, false)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 page not found"))
}
