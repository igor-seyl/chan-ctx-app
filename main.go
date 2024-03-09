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
	urls := make(chan string)
	go generateURLs(urls, 6)
	results := ParallelDownload(context.Background(), urls, 3)
	printResults(&results)
}

func ParallelDownload(ctx context.Context, urls <-chan string, numWorkers int) map[string]SiteContent {
	result := map[string]SiteContent{}
	wg := sync.WaitGroup{}
	mutex := sync.Mutex{}

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					fmt.Println("Ctx is done")
					return
				case url, isOpen := <-urls:
					if !isOpen {
						return
					}
					siteContent, hasError := DownloadSiteContent(ctx, url)
					if !hasError {
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

func DownloadSiteContent(ctx context.Context, url string) (siteContent SiteContent, hasError bool) {
	fmt.Println("Start Downloading", url)

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
	for url, siteContent := range *result {
		fmt.Println(count, url, siteContent)
		count++
	}
}
