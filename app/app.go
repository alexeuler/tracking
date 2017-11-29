package app

import (
	"github.com/up-finder/silk.web/app/setup"
	"github.com/up-finder/silk.web/app/utils"
	"os"
)

var Env *setup.Env

func init() {
	cmdLineArgs := utils.DecodeCmdLineArgs(os.Args)
	envName := "development"
	if val, ok := cmdLineArgs["e"]; ok {
		envName = val
	}
	Env = setup.Environment(envName)
	setup.Log(Env)
}
