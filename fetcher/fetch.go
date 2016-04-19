package fetcher

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"regexp"

	"github.com/PuerkitoBio/fetchbot"
	"github.com/PuerkitoBio/goquery"
	"github.com/Sirupsen/logrus"
	"github.com/fatih/set"
)

const (
	basePath string = "http://ratedata.gaincapital.com"
)

var years = set.New()
var months = set.New()
var zips = set.New()
var log *logrus.Logger
var yearQueue *fetchbot.Queue
var monthQueue *fetchbot.Queue
var pairQueue *fetchbot.Queue
var downloadQueue *fetchbot.Queue

func handleYears(ctx *fetchbot.Context, res *http.Response, err error) {
	log.Info("Getting available years...", ctx.Cmd.URL().String())
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}
	fmt.Printf("[%d] %s %s\n", res.StatusCode, ctx.Cmd.Method(), ctx.Cmd.URL())
	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		log.Error(err)
		return
	}
	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		link, exists := s.Attr("href")
		if exists {
			link = strings.TrimPrefix(link, ".\\")
			if years.Has(link) {
				uri := ctx.Cmd.URL().String() + "/" + link
				monthQueue.SendStringGet(uri)
				log.Info("Sending to monthQueue:", uri)
			}

		}
	})

}

func handleMonths(ctx *fetchbot.Context, res *http.Response, err error) {
	log.Info("Fetching months...", ctx.Cmd.URL().String())

	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		log.Error(err)
		return
	}

	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		link, exists := s.Attr("href")
		if exists {
			link = strings.TrimPrefix(link, ".\\")
			link = strings.Replace(link, " ", "%20", -1)
			if strings.HasPrefix(link, "0") || strings.HasPrefix(link, "1") {
				uri := ctx.Cmd.URL().String() + "/" + link
				pairQueue.SendStringGet(uri)
			}
		}
	})
}

func handlePairs(ctx *fetchbot.Context, res *http.Response, err error) {
	log.Info("Fetching pairs...", ctx.Cmd.URL().String())
	r := regexp.MustCompile(`\w.*zip`)
	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		log.Error(err)
		return
	}
	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		link, exists := s.Attr("href")
		if exists {
			link = strings.TrimPrefix(link, ".\\")
			if r.MatchString(link) {
				uri := ctx.Cmd.URL().String() + "/" + link
				file := regexp.MustCompile(`\w\w\w_.*`)
				directory := regexp.MustCompile(`\/\w\w\w\w\/.*?\/`)
				folderErr := os.MkdirAll("./"+directory.FindString(uri), os.ModePerm)

				if folderErr != nil {
					log.Error(err)
					return
				}

				filePath := "./" + directory.FindString(uri) + file.FindString(uri)

				if _, pathErr := os.Stat(filePath); pathErr == nil {
					log.Warn("File already exists.")
					return
				}

				downloadQueue.SendStringGet(uri)

			}
		}
	})
}

func handleDownload(ctx *fetchbot.Context, res *http.Response, err error) {
	log.Info("Downloading...", ctx.Cmd.URL().String())
	file := regexp.MustCompile(`\w\w\w_.*`)
	directory := regexp.MustCompile(`\/\w\w\w\w\/.*?\/`)

	filePath := "./" + directory.FindString(ctx.Cmd.URL().String()) + file.FindString(ctx.Cmd.URL().String())

	if _, pathErr := os.Stat(filePath); pathErr == nil {
		log.Warn("File already exists.")
		return
	}

	log.Info("Saving to path:", filePath)

	out, fileErr := os.Create(filePath)
	if fileErr != nil {
		log.Error(fileErr)
	}
	defer out.Close()
	io.Copy(out, res.Body)
	log.Info("Done.")

}

// Do begins fetching from URL url
func Do(from int, to int) {
	log = logrus.New()

	for i := from; i < to+1; i++ {
		years.Add(strconv.Itoa(i))
	}
	log.Info("Starting scraper...")

	fetchYearBot := fetchbot.New(fetchbot.HandlerFunc(handleYears))
	fetchMonthBot := fetchbot.New(fetchbot.HandlerFunc(handleMonths))
	fetchPairBot := fetchbot.New(fetchbot.HandlerFunc(handlePairs))
	downloadBot := fetchbot.New(fetchbot.HandlerFunc(handleDownload))

	yearQueue = fetchYearBot.Start()
	monthQueue = fetchMonthBot.Start()
	pairQueue = fetchPairBot.Start()
	downloadQueue = downloadBot.Start()

	yearQueue.SendStringGet(basePath)

	yearQueue.Close()
	monthQueue.Close()
	pairQueue.Close()
	downloadQueue.Close()

	log.Info("Done.")

}

func spawnWorker(tasks chan string, wg *sync.WaitGroup) {
	for cmd := range tasks {
		log.Info("Fetching:", cmd)
	}
	wg.Done()

}
