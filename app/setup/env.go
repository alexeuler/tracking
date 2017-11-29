package setup

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
)

const (
	APP_ID = "github.com/up-finder/silk.web"
)

//Struct for config.json - user-defined parameters
type Config struct {
	DB struct {
		File struct {
			NumberOfFiles   int
			NumberOfRetries int
			RetryDelay      int
			Path            string
			ExportPath      string
		}
		Redis struct {
			Host     string
			Password string
			DB       int64
		}
	}
	Log struct {
		Path  string
		Level string
	}
	System struct {
		PidPath string
	}
	Server struct {
		Port int
	}
}

//Application environment
type Env struct {
	Config
	Root        string
	Production  bool
	Development bool
	Staging     bool
	Testing     bool
}

// Environment initializer, reads config from config folder specified by the name argument
// Name = production | staging | development | testing
// The root dir for the dev and test environment is based on gopath and APP_ID,
// o/w it is the dir, the app was launched from
func Environment(name string) (env *Env) {
	env = new(Env)
	env.Production = (name == "production")
	env.Staging = (name == "staging")
	env.Testing = (name == "testing")
	if !(env.Production || env.Testing || env.Staging) {
		name = "development"
		env.Development = true
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Panicf("Env initializer: Could not get current dir: %v", err)
	}
	var gopath string
	if env.Development || env.Testing {
		gopath = os.Getenv("GOPATH")
		if gopath == "" {
			log.Panicf("Env initializer: Could not get gopath environment variable")
		}
	}
	switch name {
	case "development", "testing":
		env.Root = path.Join(gopath, "src", APP_ID)
	case "staging", "production":
		env.Root = wd
	}
	configPath := path.Join(env.Root, fmt.Sprintf("config/%s.json", name))
	mergeConfig(configPath, env)
	return env
}

// Loads json file by path and merges it into environment
var mergeConfig = func(path string, env *Env) {
	file, err := os.Open(path)
	if err != nil {
		log.Panic("Could not open %s file: %v", path, err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&env)
	if err != nil {
		log.Panic("Could not parse %s file: %v", path, err)
	}
}
