package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

// QuoteData :
type QuoteData struct {
	ID         string   `json:"id"`
	Quote      string   `json:"quote"`
	Length     string   `json:"length"`
	Author     string   `json:"author"`
	Tags       []string `json:"tags"`
	Category   string   `json:"category"`
	Date       string   `json:"date"`
	Permalink  string   `json:"parmalink"`
	Title      string   `json:"title"`
	Background string   `json:"Background"`
}

// QuoteResponse :
type QuoteResponse struct {
	Success  APISuccess   `json:"success"`
	Contents QuoteContent `json:"contents"`
}

// QuoteContent :
type QuoteContent struct {
	Quotes    []QuoteData `json:"quotes"`
	Copyright string      `json:"copyright"`
}

// APISuccess :
type APISuccess struct {
	Total string `json:"total"`
}

// ==================================================================================================
func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome! Please hit the `/qod` API to get the quote of the day."))
}

func quoteOfTheDayHandler(client *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currentTime := time.Now()
		date := currentTime.Format("2006-01-02")

		val, err := client.Get(date).Result()
		if err == redis.Nil {
			log.Println("Cache miss for date :", date)
			quoteResp, err := getQuoteFromAPI()
			if err != nil {
				w.Write([]byte("Sorry! We could not get the Quote of the Day. Please try again."))
				return
			}
			quote := quoteResp.Contents.Quotes[0].Quote
			client.Set(date, quote, 24*time.Hour)
			w.Write([]byte(quote))
		} else {
			log.Println("Cache Hit for date :", date)
			w.Write([]byte(val))
		}
	}
}

// ==================================================================================================
// ==================================================================================================
func main() {
	// Create Redis Client
	client := redis.NewClient(&redis.Options{
		Addr:     getEnv("REDIS_URL", "localhost:6379"),
		Password: getEnv("REDIS_PASSWORD", ""),
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		log.Fatal(err)
	}

	// Create Server and Route Handlers
	r := mux.NewRouter()

	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/qod", quoteOfTheDayHandler(client))

	srv := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start Server
	go func() {
		log.Println("Starting Server")
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	// Graceful Shutdown
	waitForShutdown(srv)
}

func waitForShutdown(srv *http.Server) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	<-interruptChan

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	srv.Shutdown(ctx)

	log.Println("Shutting down")
	os.Exit(0)
}

func getQuoteFromAPI() (*QuoteResponse, error) {
	APIURL := "http://quotes.rest/qod.json"
	resp, err := http.Get(APIURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	log.Println("Quote API Returned: ", resp.StatusCode, http.StatusText(resp.StatusCode))

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		quoteResp := &QuoteResponse{}
		json.NewDecoder(resp.Body).Decode(quoteResp)
		return quoteResp, nil
	}

	return nil, errors.New("Could not get quote from API")

}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// codes from : https://www.callicoder.com/docker-compose-multi-container-orchestration-golang/
