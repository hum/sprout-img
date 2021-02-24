package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/hum/sprout-img/internal"
	"github.com/turnage/graw/reddit"
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

	bot, err := reddit.NewBotFromAgentFile("image.agent", 0)
	if err != nil {
		panic(err)
	}

	count = 0
	for _, category := range data.categories {
		log.Printf("Fetching image links from: %s\n", category.CategoryName)

		harvest, err := bot.Listing("/r/"+category.CategoryName, "")
		if err != nil {
			log.Println(err)
			continue
		}

		for _, post := range harvest.Posts {
			if strings.Contains(post.URL, "/comments/") {
				continue
			}

			image := &internal.PImages{
				CategoryID: category.CategoryID,
				Link:       post.URL,
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
