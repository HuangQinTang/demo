package config

import "time"

type tomlConfig struct {
	Server Server `toml:"server"`
	App    App    `toml:"app"`
	Verify Verify `toml:"verify"`
	Redis  Redis  `toml:"redis"`
	Nsq    Nsq    `toml:"nsq"`
}

type Server struct {
	Address    string `toml:"address"`
	ListenPort string `toml:"listenPort"`
}

type App struct {
	ClientDebug bool `toml:"clientDebug"`
	ServerDebug bool `toml:"serverDebug"`
}

type Verify struct {
	IllegalStr []string `toml:"illegalStr"`
}

type Redis struct {
	Address     string        `toml:"address"`
	MaxIdle     int           `toml:"maxIdle"`
	MaxActive   int           `toml:"maxActive"`
	IdleTimeout time.Duration `toml:"idleTimeout"`
}

type Nsq struct {
	TestTopic   string `toml:"testTopic"`
	TestChannel string `toml:"testChannel"`
	MesTopic    string `toml:"mesTopic"`
	MesChannel  string `toml:"mesChannel"`
	LookupAddr  string `toml:"lookupAddr"`
	Nsqd1Tcp    string `toml:"nsqd1Tcp"`
	Nsqd1Http   string `toml:"nsqd1Http"`
	MaxInFlight int    `toml:"maxInFlight"`
}
