package config

import (
	"strings"
	"time"

	gostreamv1 "github.com/lordvidex/gostream/pkg/api/gostream/v1"
)

// TODO: add config validations
// Config ...
type Config struct {
	Server Server `mapstructure:"server"`
	Client Client `mapstructure:"client"`
}

// Server ...
type Server struct {
	LogFile       string `mapstructure:"log_file"`
	DSN           string `mapstructure:"dsn"`
	RedisURL      string `mapstructure:"redis_url"`
	GRPCPort      int64  `mapstructure:"grpc_port"`
	HTTPPort      int64  `mapstructure:"http_port"`
	RunMigrations bool   `mapstructure:"run_migrations"`
	DryRun        bool   `mapstructure:"dry_run"`
}

// Client ...
type Client struct {
	Name              string
	LogFile           string `mapstructure:"log_file"`
	Servers           []string
	Entities          SpaceSeparatedEntities
	Connections       int
	ConnectionTimeout time.Duration `mapstructure:"connection_timeout"`
	DryRun            bool          `mapstructure:"dry_run"`
}

// SpaceSeparatedEntities ...
type SpaceSeparatedEntities []string

// UnmarshalText ...
func (s SpaceSeparatedEntities) Values() []gostreamv1.Entity {
	arr := make([]gostreamv1.Entity, 0, len(s))
	for _, str := range s {
		switch strings.ToLower(str) {
		case "all":
			arr = append(arr, gostreamv1.Entity_ENTITY_UNSPECIFIED)
			continue
		case "pet":
			arr = append(arr, gostreamv1.Entity_ENTITY_PET)
			continue
		case "user":
			arr = append(arr, gostreamv1.Entity_ENTITY_USER)
			continue
		}

		if v, ok := gostreamv1.Entity_value[str]; ok {
			arr = append(arr, gostreamv1.Entity(v))
		}
	}

	return arr
}
