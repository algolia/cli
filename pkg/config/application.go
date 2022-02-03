package config

type Application struct {
	Name string
	ID   string `mapstructure:"application_id"`

	AdminAPIKey      string `mapstructure:"admin_api_key"`
	UsageAPIKey      string `mapstructure:"usage_api_key"`
	MonitoringAPIKey string `mapstructure:"monitoring_api_key"`
}
