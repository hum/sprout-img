package internal

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/go-pg/pg/v10"
)

type Database struct {
	name   string
	server string
	Conn   *pg.DB
}

type DbConfig struct {
	User     string
	Password string
	Url      string
	Port     string
	DbName   string
}

func CreateDb(filepath string) (*Database, error) {
	if filepath == "" {
		return nil, errors.New("The filepath parameter to create db config is empty.")
	}

	dbConfig, err := parseDbConfig(filepath)
	if err != nil {
		return nil, err
	}

	db, err := createDbObj(dbConfig)
	if err != nil {
		return nil, err
	}

	return db, nil
}
func createDbObj(c *DbConfig) (*Database, error) {
	if c.User == "" || c.Password == "" || c.Url == "" || c.Port == "" || c.DbName == "" {
		return nil, errors.New("Config does not contain one or more of the necessary strings.")
	}

	opt, err := pg.ParseURL(c.String())
	if err != nil {
		return nil, err
	}

	db := &Database{
		name:   c.DbName,
		server: c.Url + c.Port,
		Conn:   pg.Connect(opt),
	}
	return db, nil
}

func (c *DbConfig) String() string {
	result := "postgres://"
	result += c.User + ":" + c.Password + "@"
	result += c.Url + ":" + c.Port + "/" + c.DbName + "?sslmode=disable"
	return result
}

func parseDbConfig(path string) (*DbConfig, error) {
	if path == "" {
		return nil, errors.New("Cannot read the config file. Empty path argument.")
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	db := &DbConfig{}
	err = json.Unmarshal(data, db)
	if err != nil {
		return nil, err
	}

	return db, nil
}
