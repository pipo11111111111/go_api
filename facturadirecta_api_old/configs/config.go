package configs

import "github.com/spf13/viper"

var cfg *config

type config struct {
	API APIConfing
	DB  DBConfig
	B   B2B
}

type APIConfing struct {
	Port string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Pass     string
	Database string
}
type B2B struct {
	Api_key    string
	Value string
	Account_id string
}

func init() {
	viper.SetDefault("api.port", "9000")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", "5432")
}

func Load() error {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()

	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	cfg = new(config)

	cfg.API = APIConfing{
		Port: viper.GetString("api.port"),
	}

	cfg.DB = DBConfig{
		Host:     viper.GetString("database.host"),
		Port:     viper.GetString("database.port"),
		User:     viper.GetString("database.user"),
		Pass:     viper.GetString("database.pass"),
		Database: viper.GetString("database.name"),
	}
	cfg.B = B2B{
		Api_key:    viper.GetString("b2b.api_key"),
		Value:    viper.GetString("b2b.value"),
		Account_id: viper.GetString("b2b.account_id"),
	}
	return nil
}

func GetDB() DBConfig {
	return cfg.DB
}

func GetServerPort() string {
	return cfg.API.Port
}

func GetB2B() B2B {
	return cfg.B
}
