package setup

import (
	"os"
	"path"
	"testing"
)

func TestEnv(t *testing.T) {
	smergeConfig := mergeConfig
	defer func() { mergeConfig = smergeConfig }()

	prodpath, _ := os.Getwd()
	devpath := path.Join(os.Getenv("GOPATH"), "src/github.com/up-finder/silk.web")

	testCases := []struct {
		name           string
		expectedRoot   string
		expectedConfig string
	}{
		{"development", devpath, path.Join(devpath, "config/development.json")},
		{"", devpath, path.Join(devpath, "config/development.json")},
		{"staging", prodpath, path.Join(prodpath, "config/staging.json")},
		{"production", prodpath, path.Join(prodpath, "config/production.json")},
		{"testing", devpath, path.Join(devpath, "config/testing.json")},
	}
	for _, tc := range testCases {
		mergeConfig = func(path string, env *Env) {
			if path != tc.expectedConfig {
				t.Fatalf("Init Env: Expected %s config path, got %s", tc.expectedConfig, path)
			}
		}
		got := Environment(tc.name)
		if tc.expectedRoot != got.Root {
			t.Errorf("On input: %v, expected env.root to be: %v, got: %v", tc.name, tc.expectedRoot, got.Root)
		}
	}
}
