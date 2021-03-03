package sproutimg

import (
	"github.com/hum/sprout"
	"github.com/hum/sprout-img/internal"
	"log"
	"net/http"
	"os"
	"strconv"
)

var databaseHandler *internal.DatabaseHandler

type imageData struct {
	categories []internal.PImageCategory
	servers    []internal.PServers
}

func getData(db *internal.Database) (*imageData, error) {
	var categories []internal.PImageCategory
	var servers []internal.PServers

	err := db.Conn.Model(&categories).Select()
	if err != nil {
		return nil, err
	}

	err = db.Conn.Model(&servers).Select()
	if err != nil {
		return nil, err
	}

	return &imageData{
		categories: categories,
		servers:    servers,
	}, nil
}

func HandleImageCollection(w http.ResponseWriter, r *http.Request) {
  if databaseHandler == nil {
    databaseHandler = &internal.DatabaseHandler{}
  }
	// Only supports Reddit for now
	status := DEFAULT

	/*
	  TODO:
	  Create a database pool. Opening up a new connection is expensive
	*/
  dbUri := os.Getenv("POSTGRE_URI")
  database, err := databaseHandler.Connect(dbUri)
  if err != nil {
    /*
      TODO:
      handle errors properly
    */
    panic(err)
  }
	defer database.Conn.Close()

	data, err := getData(database)
	if err != nil {
		/*
		   TODO:
		   handle errors properly.
		*/
		panic(err)
	}

	s := sprout.New()
	reddit := s.Reddit()
	reddit.UseAPI = true

	reddit.Conf = &sprout.Config{
		Username:     os.Getenv("REDDIT_USERNAME"),
		Password:     os.Getenv("REDDIT_PASSWORD"),
		UserAgent:    os.Getenv("REDDIT_USER_AGENT"),
		ClientID:     os.Getenv("REDDIT_CLIENT_ID"),
		ClientSecret: os.Getenv("REDDIT_CLIENT_SECRET"),
	}

	/*
	  TODO:
	  Handle concurrently
	*/
	count := 0
	for _, c := range data.categories {
		result, err := reddit.Get(c.CategoryName, 100)
		if err != nil {
			log.Println(err)
			continue
		}

		for _, post := range result.Posts {
			image := &internal.PImages{
				CategoryID: c.CategoryID,
				Link:       post.Link,
			}

			_, err := database.Conn.Model(image).Insert()
			if err != nil {
				continue
			}
			count++

			for _, s := range data.servers {
				serverImage := &internal.PServerImages{
					ServerID: s.ServerID,
					ImageID:  image.ImageID,
				}

				_, err := database.Conn.Model(serverImage).Insert()
				if err != nil {
					continue
				}
			}
			status = OK
		}
	}

	response := Response{
		Status: status,
		Data:   map[string]string{"Count": strconv.Itoa(count)},
		/*
		   TODO:
		   Possibly add more data to return?
		   E.g. Return what source collected what amount of images
		*/
	}

	err = setHeader(w, map[string]string{"Content-Type": "application/json"})
	if err != nil {
		panic(err)
	}

	obj, err := marshalResponse(response)
	if err != nil {
		panic(err)
	}
	w.Write(obj)
}
