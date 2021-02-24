package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/hum/sprout"
	"github.com/hum/sprout-img/internal"
)

var count int

type Data struct {
	categories []internal.PImageCategory
	servers    []internal.PServers
}

type Response struct {
	Count  int
	Status string
}

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

func getConfig() *sprout.Config {
	return &sprout.Config{
		Username:     os.Getenv("REDDIT_USERNAME"),
		Password:     os.Getenv("REDDIT_PASSWORD"),
		UserAgent:    os.Getenv("REDDIT_USER_AGENT"),
		ClientID:     os.Getenv("REDDIT_CLIENT_ID"),
		ClientSecret: os.Getenv("REDDIT_CLIENT_SECRET"),
	}
}

func handleResponse(w http.ResponseWriter, r *http.Request) {
	database, err := internal.CreateDb("db_config.json")
	if err != nil {
		panic(err)
	}
	defer database.Conn.Close()

	data, err := getData(database)
	if err != nil {
		panic(err)
	}

	subreddits := make([]string, len(data.categories))
	for _, c := range data.categories {
		subreddits = append(subreddits, c.CategoryName)
	}

	sprout := sprout.New()
	reddit := sprout.Reddit()
	reddit.UseAPI = true
	reddit.Conf = getConfig()

	sub, err := reddit.Get(subreddits, 100)
	if err != nil {
		log.Println(err)
		continue
	}

	count = 0
	for _, category := range data.categories {
		for _, post := range sub[category.CategoryName].Posts {
			image := &internal.PImages{
				CategoryID: category.CategoryID,
				Link:       post.Link,
			}

			_, err := database.Conn.Model(image).Insert()
			if err != nil {
				log.Printf("Error inserting data into DB: %v\n", err)
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
					log.Printf("Error inserting data into DB: %v\n", err)
					continue
				}
			}
		}
	}

	result := Response{}
	w.Header().Set("Content-Type", "application/json")
	if count != 0 {
		result.Status = "400 OK"
	} else {
		result.Status = "200 ERROR"
	}

	result.Count = count
	json, err := json.Marshal(result)
	if err != nil {
		log.Printf("Error marshaling data %v\n", err)
	}

	w.Write(json)
}

func main() {
	http.HandleFunc("/", handleResponse)
	http.ListenAndServe(":3000", nil)
}
