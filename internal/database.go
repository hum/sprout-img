package internal

import (
  "fmt"
  "github.com/go-pg/pg/v10"
)

type Database struct {
  Conn *pg.DB
}

type DatabaseHandler struct {
  db *Database
}

func createDatabase(uri string) (*Database, error) {
  opt, err := pg.ParseURL(uri)
  if err != nil {
    return nil, err
  }

  return &Database{
    Conn: pg.Connect(opt),
  }, nil
}

func checkConnection(d *Database) bool {
  _, err := d.Conn.Exec("SELECT 1");
  if err != nil {
    return false
  }
  return true
}

func (dh *DatabaseHandler) Connect(uri string) (*Database, error) {
  if uri == "" {
    return nil, fmt.Errorf("URI string for Connect() is empty.")
  }

  if dh.db == nil {
    db, err := createDatabase(uri)
    if err != nil {
      return nil, err
    }
    dh.db = db
    return db, nil
  }

  if ok := checkConnection(dh.db); ok {
    return dh.db, nil
  } else {
    dh.db = nil
    return dh.Connect(uri)
  } 
}
