package main

import (
	"database/sql"
	"errors"
	"flag"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type baseConfig struct {
	HttpAddress     string `yaml:"http_address"`
	StorageURL      string `yaml:"storage_url"`
	TemplatesPath   string `yaml:"templates_path"`
	GreetingAPIHost string `yaml:"greetingAPI_host"`
}

type AppConfigFile struct {
	Base baseConfig `yaml:"base"`
}

type RuntimeState struct {
	Config      AppConfigFile
	dbType      string
	db          *sql.DB
	authcookies map[string]cookieInfo
	cookiemutex sync.Mutex
	jwtTokenmap map[string]tokenInfo
	tokenmutex  sync.Mutex
}

type tokenInfo struct {
	Token     string
	ExpiresAt time.Time
}

type cookieInfo struct {
	Username  string
	ExpiresAt time.Time
}

type Response struct {
	OutResponse string
}

var (
	configFilename = flag.String("config", "config.yml", "The filename of the configuration")
)

const (
	cookieExpirationHours = 6
	jwtTokenExpirationMin = 30
	cookieName            = "auth_cookie"

	loginPath  = "/login-signup-page"
	checklogin = "/checklogin"
	register   = "/register"
	logoutpath ="/logout"
	indexPath    = "/"
	greetingPath = "/greeting"
	cssPath      = "/css/"
	imagesPath   = "/images/"
	jsPath       = "/js/"
)

//parses the config file and initialize required structs
func loadConfig(configFilename string) (RuntimeState, error) {

	var state RuntimeState

	if _, err := os.Stat(configFilename); os.IsNotExist(err) {
		err = errors.New("mising config file failure")
		return state, err
	}

	//ioutil.ReadFile returns a byte slice (i.e)(source)
	source, err := ioutil.ReadFile(configFilename)
	if err != nil {
		err = errors.New("cannot read config file")
		return state, err
	}

	//Unmarshall(source []byte,out interface{})decodes the source byte slice/value and puts them in out.
	err = yaml.Unmarshal(source, &state.Config)

	if err != nil {
		err = errors.New("Cannot parse config file")
		log.Printf("Source=%s", source)
		return state, err
	}
	err = initDB(&state)
	if err != nil {
		return state, err
	}
	state.authcookies = make(map[string]cookieInfo)
	state.jwtTokenmap = make(map[string]tokenInfo)
	return state, err
}

func main() {
	flag.Parse()

	state, err := loadConfig(*configFilename)
	if err != nil {
		panic(err)
	}

	http.Handle(loginPath, http.HandlerFunc(state.loginPageHandler))

	http.Handle(greetingPath, http.HandlerFunc(state.GreetingHandler))

	http.Handle(indexPath, http.HandlerFunc(state.IndexHandler))

	http.Handle(checklogin, http.HandlerFunc(state.Checklogin))

	http.Handle(register, http.HandlerFunc(state.RegisterUser))
	http.Handle(logoutpath,http.HandlerFunc(state.logoutHanler))

	fs := http.FileServer(http.Dir(state.Config.Base.TemplatesPath))
	http.Handle(cssPath, fs)
	http.Handle(imagesPath, fs)
	http.Handle(jsPath, fs)
	log.Fatal(http.ListenAndServe(state.Config.Base.HttpAddress, nil))
}
