package fetcher

import (
	"io"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"

	"regexp"

	"github.com/PuerkitoBio/fetchbot"
	"github.com/PuerkitoBio/goquery"
	"github.com/Sirupsen/logrus"
	"github.com/fatih/set"
)

const (
	basePath string = "http://ratedata.gaincapital.com"
)

var (
	years  = set.New()
	months = set.New()
	zips   = set.New()
	pairs  = set.New()

	log           *logrus.Logger
	yearQueue     *fetchbot.Queue
	monthQueue    *fetchbot.Queue
	pairQueue     *fetchbot.Queue
	downloadQueue *fetchbot.Queue

	// sigChan receives os signals.
	sigChan = make(chan os.Signal, 1)

	// complete is used to report processing is done.
	complete = make(chan error)

	// shutdown provides system wide notification.
	shutdown = make(chan struct{})
)

func handleYears(ctx *fetchbot.Context, res *http.Response, err error) {
	log.Infoln("Getting available years...", ctx.Cmd.URL().String())
	if err != nil {
		log.Errorf("error: %s\n", err)
		return
	}
	log.Infof("[%d] %s %s\n", res.StatusCode, ctx.Cmd.Method(), ctx.Cmd.URL())
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
			}

		}
	})

}

func handleMonths(ctx *fetchbot.Context, res *http.Response, err error) {
	log.Infoln("Fetching months...", ctx.Cmd.URL().String())

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
	log.Infoln("Fetching pairs...", ctx.Cmd.URL().String())
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
				pairName := link[:7]

				if !pairs.Has(pairName) {
					return
				}

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
	log.Infoln("Downloading...", ctx.Cmd.URL().String())
	file := regexp.MustCompile(`\w\w\w_.*`)
	directory := regexp.MustCompile(`\/\w\w\w\w\/.*?\/`)

	filePath := "./" + directory.FindString(ctx.Cmd.URL().String()) + file.FindString(ctx.Cmd.URL().String())

	if _, pathErr := os.Stat(filePath); pathErr == nil {
		log.Warn("File already exists.")
		return
	}

	log.Infoln("Saving to path:", filePath)

	out, fileErr := os.Create(filePath)
	if fileErr != nil {
		log.Error(fileErr)
	}
	defer out.Close()
	io.Copy(out, res.Body)
	log.Infoln("Done.")

}

// Entry is the entry point for the CLI app
func Entry(from int, to int, pairsFlag []string) {
	log = logrus.New()
	log.Infoln("Running fetch from", from, "to", to, "for pairs", pairsFlag)
	for _, p := range pairsFlag {
		pairs.Add(p)
	}

	log.Infoln("Starting fetcher app")
	signal.Notify(sigChan, os.Interrupt)

	log.Infoln("Launching scraper")
	go processor(from, to)

ControlLoop:
	for {
		select {
		case <-sigChan:
			log.Error("OS interrupt received")
			yearQueue.Cancel()
			monthQueue.Cancel()
			pairQueue.Cancel()
			downloadQueue.Cancel()
			close(shutdown)
			sigChan = nil

		case err := <-complete:
			log.Infof("Scraper Completed: Error[%s]", err)
			break ControlLoop
		}
	}
	log.Println("Fetcher complete")
}

func scrape(from int, to int) {
	for i := from; i < to+1; i++ {
		years.Add(strconv.Itoa(i))
	}
	log.Infoln("Starting scraper...")

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

	complete <- nil
}
