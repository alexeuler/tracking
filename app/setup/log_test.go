package setup

import (
	logger "github.com/Sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestLog(t *testing.T) {
	cases := []struct {
		env           Env
		inputLevel    string
		expectedLevel logger.Level
		out           string
	}{
		{
			env:           Env{Development: true},
			inputLevel:    "Info",
			expectedLevel: logger.InfoLevel,
			out:           "stdout",
		},
		{
			env:           Env{Testing: true},
			inputLevel:    "Null",
			expectedLevel: logger.ErrorLevel,
			out:           "null",
		},
		{
			env:           Env{Staging: true},
			inputLevel:    "Debug",
			expectedLevel: logger.DebugLevel,
			out:           "file",
		},
		{
			env:           Env{Production: true},
			inputLevel:    "Error",
			expectedLevel: logger.ErrorLevel,
			out:           "file",
		},
	}

	for _, c := range cases {
		gopath := os.Getenv("GOPATH")
		c.env.Root = path.Join(gopath, "src", APP_ID)
		c.env.Log.Level = c.inputLevel
		if c.env.Production || c.env.Staging {
			c.env.Log.Path = "tmp/test.log"
		}
		Log(&c.env)
		if logger.GetLevel() != c.expectedLevel {
			t.Errorf("Init Log: expected %s level, got %s level", c.env.Log.Level, c.expectedLevel)
		}
		switch c.out {
		case "null":
			if logger.StandardLogger().Out != ioutil.Discard {
				t.Errorf("Init Log: Env: %v: expected Discard out, got %v out", c.env, logger.StandardLogger().Out)
			}
		case "file":
			if _, ok := logger.StandardLogger().Out.(*os.File); !ok || logger.StandardLogger().Out == os.Stdout {
				t.Errorf("Init Log: Env: %v: expected File out, got %v out", c.env, logger.StandardLogger().Out)
			}
		case "stdout":
			if logger.StandardLogger().Out != os.Stdout {
				t.Errorf("Init Log: Env: %v: expected Stdout out, got %v out", c.env, logger.StandardLogger().Out)
			}
		}
	}

}
