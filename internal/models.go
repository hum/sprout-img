package internal

import (
	"time"
)

type PImageCategory struct {
	tableName    struct{}   `pg:"p_image_category"`
	CategoryID   int        // category_id
	CategoryName string     // category_name
	AddTime      *time.Time // add_time
}

type PImages struct {
	ImageID    int        // image_id
	CategoryID int        // category_id
	Link       string     // link
	AddTime    *time.Time // add_time
}

type PServerImages struct {
	ImageID  int // image_id
	ServerID int // server_id
}

type PServers struct {
	ServerID   int    // server_id
	ServerName string // server_name
}
