package config

type Config struct {
	Client Client `mapstructure:"client"`
	Server Server `mapstructure:"server"`
}

type Server struct {
	LogFile       string `mapstructure:"log_file"`
	DSN           string `mapstructure:"dsn"`
	RedisURL      string `mapstructure:"redis_url"`
	GRPCPort      int64  `mapstructure:"grpc_port"`
	HTTPPort      int64  `mapstructure:"http_port"`
	RunMigrations bool   `mapstructure:"run_migrations"`
	DryRun        bool   `mapstructure:"dry_run"`
}

type Client struct {
	Servers     []string
	Connections int
	Entities    []string
	DryRun      bool `mapstructure:"dry_run"`
}
