package main

import (
	"fmt"
	"math/rand"
	"time"
)

var letters = "abcdefghijklmnopqrstuvwxyz"
var domainZones = []string{"com", "net", "ru", "org", "ua"}

func generateURLs(urls chan<- string, urlAmount int) {
	for i := 0; i < urlAmount; i++ {
		urls <- generateURL()
	}
	close(urls)
}

func generateURL() string {
	initRandomSleep(500, 1000)

	name := ""
	minNameLen := 3
	maxNameLen := 7
	nameLen := rand.Intn(maxNameLen+1-minNameLen) + minNameLen

	for len(name) != nameLen {
		name = name + string(letters[rand.Intn(len(letters))])
	}

	url := "http://" + name + "." + domainZones[rand.Intn(len(domainZones))]
	fmt.Println("Generated url is", url)
	return url
}

func initRandomSleep(minMS int, maxMS int) {
	r := rand.Intn(maxMS+1-minMS) + minMS
	time.Sleep(time.Duration(r) * time.Millisecond)
}
