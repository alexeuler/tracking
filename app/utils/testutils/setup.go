package testutils

import (
	"github.com/up-finder/silk.web/app"
	"github.com/up-finder/silk.web/app/db"
	"github.com/up-finder/silk.web/app/setup"
)

func Setup() {
	app.Env = setup.Environment("testing")
	setup.Log(app.Env)
	if db.Redis!=nil {
		db.Redis.Close()
	}
	db.Redis = db.NewRedis(app.Env)
	db.Redis.FlushDb()
	db.File = db.NewFileDB(app.Env)
}