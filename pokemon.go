package main

import (
	"strconv"

	"github.com/tidwall/gjson"
)

type PokemonInfo struct {
	Name     string `json:"name"`
	ImageUrl string `json:"image_url"`
}

func GetPokemonImage(id int) *PokemonInfo {
	pokeID := strconv.Itoa(id)
	data := ClientRequest(POKEMON_REQUEST_URL+pokeID, nil, nil)

	byteData := string(data)

	name := gjson.Get(byteData, "name")
	image := gjson.Get(byteData, "sprites.other.official-artwork.front_default")

	pokemonInfo := &PokemonInfo{
		Name:     name.Str,
		ImageUrl: image.Str,
	}

	return pokemonInfo
}
