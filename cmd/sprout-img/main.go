package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/hum/sprout"
	"github.com/hum/sprout-img/internal"
)

var count int
var status string

type Data struct {
	categories []internal.PImageCategory
	servers    []internal.PServers
}

type Response struct {
	Count  int
	Status string
}

/*func handleWebhook(w http.ResponseWriter, r *http.Request) {

}*/

func getData(db *internal.Database) (*Data, error) {
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

	return &Data{
		categories: categories,
		servers:    servers,
	}, nil
}

func handleResponse(w http.ResponseWriter, r *http.Request) {
	status = "EMPTY"

	database, err := internal.CreateDb("db_config.json")
	if err != nil {
		panic(err)
	}
	defer database.Conn.Close()

	data, err := getData(database)
	if err != nil {
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
	count = 0
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
			status = "OK"
		}
	}

	response := Response{}
	w.Header().Set("Content-Type", "application/json")
	response.Status = status

	response.Count = count
	json, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshaling data %v\n", err)
	}
	w.Write(json)
}

func main() {
	http.HandleFunc("/", handleResponse)

	/*
	  TODO:
	  Serve as HTTPS -- generate private SSL cert and key -- to be able to handle webhooks

	  http.HandleFunc("/wh", handleWebhook)
	*/

	http.ListenAndServe(":3000", nil)
}
