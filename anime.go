package main

import (
	"encoding/json"
	"log"
)

type AnimeImageReponse struct {
	URL string `json:"url"`
}

func GetAnimeImage() string {
	data := ClientRequest(ANIME_REQUEST_URL, nil, nil)

	animeImageReponse := &AnimeImageReponse{}
	err := json.Unmarshal(data, animeImageReponse)
	if err != nil {
		log.Print("anime image response error: ", err)
	}
	return animeImageReponse.URL
}
