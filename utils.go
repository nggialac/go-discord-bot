package main

import (
	"io/ioutil"
	"log"
	"math/rand"
	"time"
)

func RandomInt(max, min int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func GetImageFromFilePath(filePath string) ([]byte, error) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	return file, err
}
