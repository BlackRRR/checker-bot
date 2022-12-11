package config

import (
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"path/filepath"
)

type Config struct {
	Pgx    *pgxpool.Config
	DB     DB     `yaml:"db"`
	Bots   []Bot  `yaml:"bots"`
	Server Server `yaml:"server"`
}

type DB struct {
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	Host         string `yaml:"host"`
	Port         string `yaml:"port"`
	DBName       string `yaml:"db_name"`
	PoolMaxConns int    `yaml:"pool_max_conns"`
}

type Bot struct {
	AMSBotType string      `yaml:"ams_bot_type"`
	Components []Component `yaml:"components"`
}

type Component struct {
	Token   string `yaml:"token"`
	BotLang string `yaml:"bot_lang"`
}

type Server struct {
	IP              string `yaml:"ip"`
	AdminRout       string `yaml:"admin"`
	IncomeInfoRoute string `yaml:"income_info"`
}

func InitConfig() (*Config, string, error) {
	filename, _ := filepath.Abs("config/config.yaml")
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	var config Config

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}

	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?pool_max_conns=%d",
		config.DB.User,
		config.DB.Password,
		config.DB.Host,
		config.DB.Port,
		config.DB.DBName,
		config.DB.PoolMaxConns)

	connForMigrations := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DB.Host,
		config.DB.Port,
		config.DB.User,
		config.DB.Password,
		config.DB.DBName)

	fmt.Println(connString)
	pgxConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, "", errors.Wrap(err, "`Init config` failed to parse config")
	}

	config.Pgx = pgxConfig

	return &config, connForMigrations, nil
}
