package core

import (
	"github.com/asdine/storm/v3"
	"github.com/op/go-logging"
	"golang.org/x/net/proxy"
	"time"
)

type Config struct {
	RepoPath string
	Logger   logging.Backend
	DB       *storm.DB
	Mnemonic string
	Created  time.Time
	Proxy    proxy.Dialer
}
