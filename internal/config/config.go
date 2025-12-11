package config

type Redis struct {
	Host     string `mapstructure:"host"`
	Port     uint16 `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Pass     string `mapstructure:"pass"`
	Stream   string `mapstructure:"stream"`
	Group    string `mapstructure:"group"`
	Consumer string `mapstructure:"consumer"`
	DB       int    `mapstructure:"db"`
	MaxLen   int64  `mapstructure:"max_len"`
}

type Postgres struct {
	Host            string `mapstructure:"host"`
	Port            uint16 `mapstructure:"port"`
	DB              string `mapstructure:"db"`
	User            string `mapstructure:"user"`
	Pass            string `mapstructure:"pass"`
	SSLMode         string `mapstructure:"sslmode"`
	MaxConns        int32  `mapstructure:"max_conns"`
	MinConns        int32  `mapstructure:"min_conns"`
	MaxConnLifetime string `mapstructure:"max_conn_lifetime"`
	MaxConnIdleTime string `mapstructure:"max_conn_idle_time"`
	HealthTimeout   string `mapstructure:"health_timeout"`
}

type App struct {
	Redis    Redis    `mapstructure:"redis"`
	Postgres Postgres `mapstructure:"postgres"`
}
