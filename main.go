package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type SiteContent struct {
	Title   string
	Date    time.Time
	Content string
}

func main() {
	// для оптимизации ресурсов, я бы посоветовал завести отдельную переменную под кол-во ссылок,
	// и использовать ее для указания размера канала и как параметр для generateURLs
	urls := make(chan string)
	go generateURLs(urls, 6)
	// В целом, в рамках такого проекта не критично, но вообще это плохая пратика передавать
	// контекст background вот так в функцию. Он создается где-то в мейне как отдельная переменная,
	// и дальше используется уже нужным образом.
	results := ParallelDownload(context.Background(), urls, 3)
	printResults(&results)
}

func ParallelDownload(ctx context.Context, urls <-chan string, numWorkers int) map[string]SiteContent {
	result := map[string]SiteContent{}
	// вг и мутексы всегда объявляются как указатели на эти структуры, иначе может быть очень плохо в случае их "копирования" при передаче
	wg := sync.WaitGroup{}
	mutex := sync.Mutex{}

	for i := 0; i < numWorkers; i++ {
		// это плохая практика, лучше сразу добавить горутины по количеству воркеров
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Выглядит так, что цикл фор тут не нужен вообще. Бесконечное ожидание будет само по себе.
			for {
				select {
				// Этот кейс у тебя никогда не выполнится, потому что ты здесь нигде контекст не отменяешь
				// Вместо этого ты его прокидываешь ниже в функцию, и уже там отменяешь. Но дочерние контексты не оказывают
				// влияние на родительский. Только сигналы завершения родительского окажут влияние на дочерние.
				case <-ctx.Done():
					fmt.Println("Ctx is done")
					return
				case url, isOpen := <-urls:
					// молодец, что предусмотрел это
					if !isOpen {
						return
					}
					siteContent, hasError := DownloadSiteContent(ctx, url)
					// имхо, что-то можно писать и в случае ошибки
					if !hasError {
						// Молодец, что предусмотрел это, мапы и правда не считаются потокобезопасными (даже для чтения).
						// Но в целом, если ты уверен, что у тебя всегда будут только уникальные урлы, то это может быть несколько необязательно.
						// В любом случае, можно использовать еще sync.Map{} и не париться, она как раз потокобезопасная.
						mutex.Lock()
						result[url] = siteContent
						mutex.Unlock()
						fmt.Println("Downloaded", url, siteContent)
					}
				}
			}
		}()
	}

	wg.Wait()
	return result
}

// Возвращение hasError не понятно, если честно. Лучше возврашать err (его можно кастомно объявить), либо как-то иначе переиграть этот флаг,
// чтобы отменять родительский контекст
func DownloadSiteContent(ctx context.Context, url string) (siteContent SiteContent, hasError bool) {
	fmt.Println("Start Downloading", url)

	// выглядит так, что в канал пишет только горутина одна, но ты делаешь его "безмерным".
	// В целом, по архитектуре вижу, что у тебя эта функция запускается в цикле, и у тебя каждый раз будет
	// создаваться новый канал для одной записи, и в него будет писать запущенная новая горутинка.
	// Это не шибко оптимально, выглядит так, что сюда надо передать общий канал и пусть в него эти горутины и пишут.
	scChan := make(chan SiteContent)
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*3500)
	defer cancel()

	go func() {
		initRandomSleep(1500, 5000)
		scChan <- SiteContent{
			Title:   "Заголовок сайта " + url,
			Date:    time.Now(),
			Content: "Содержание сайта " + url,
		}
	}()

	select {
	// Вот здесь как раз правильно, этот кейс может произойти, но он не окажет влияние на родительский контекст
	case <-ctx.Done():
		fmt.Println("Time`s up for", url)
		return SiteContent{}, true
	case siteContent := <-scChan:
		return siteContent, false
	}
}

func printResults(result *map[string]SiteContent) {
	fmt.Println("--- Download results ---")
	count := 1
	// Мапа сама по себе передается по ссылке, память тут ты не экономишь :)
	for url, siteContent := range *result {
		fmt.Println(count, url, siteContent)
		count++
	}
}
