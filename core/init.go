package core

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/kelseyhightower/envconfig"
)

var (
	App *Application
)

type (
	Application struct {
		Name    string   `json:"name"`
		Port    string   `json:"port"`
		Version string   `json:"version"`
		Config  Config   `json:"app_config"`
		DB      *gorm.DB `json:"db"`
	}

	Config struct {
		Port    string `envconfig:"APPPORT"`
		DB_HOST string `envconfig:"DB_HOST"`
		DB_USER string `envconfig:"DB_USER"`
		DB_PASS string `envconfig:"DB_PASS"`
		DB_NAME string `envconfig:"DB_NAME"`
		DB_PORT string `envconfig:"DB_PORT"`
		DB_LOG  int    `envconfig:"DB_LOG"`
	}
)

func init() {
	var err error
	App = &Application{}

	if err = App.LoadConfigs(); err != nil {
		log.Printf("Load config error APP: %v", err)
	}

	if err = App.DatabaseInit(); err != nil {
		log.Printf("Load config error DB: %v", err)
	}
}

func (x *Application) LoadConfigs() error {
	err := envconfig.Process("myapp", &x.Config)
	x.Name = "gocommerce"
	x.Version = os.Getenv("APPVER")
	x.Port = x.Config.Port
	return err
}

func (x *Application) DatabaseInit() error {

	config := x.Config
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local", config.DB_USER, config.DB_PASS, config.DB_HOST, config.DB_PORT, config.DB_NAME)

	db, err := gorm.Open("mysql", dsn)
	db.LogMode(config.DB_LOG == 1)
	x.DB = db

	return err
}

func (x *Application) Close() (err error) {
	if err = x.DB.Close(); err != nil {
		return err
	}
	return nil
}
