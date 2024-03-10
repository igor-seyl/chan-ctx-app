package main

import (
	"fmt"
	"math/rand"
	"time"
)

// можно объединить все под единый блок var, и в его скобках указать все нужные переменные
var letters = "abcdefghijklmnopqrstuvwxyz"
var domainZones = []string{"com", "net", "ru", "org", "ua"}

// молодец, что указываешь в параметрах явно, что канал на запись
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
	// кстати, а в чем смысл такой формулы?
	nameLen := rand.Intn(maxNameLen+1-minNameLen) + minNameLen

	for len(name) != nameLen {
		name = name + string(letters[rand.Intn(len(letters))])
	}

	// вместо сложения строк лучше использовать fmt.Sprintf()
	url := "http://" + name + "." + domainZones[rand.Intn(len(domainZones))]
	fmt.Println("Generated url is", url)
	return url
}

func initRandomSleep(minMS int, maxMS int) {
	r := rand.Intn(maxMS+1-minMS) + minMS
	time.Sleep(time.Duration(r) * time.Millisecond)
}
