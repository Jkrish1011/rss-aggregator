package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/Jkrish1011/rss-aggregator/internal/database"
	"github.com/google/uuid"
)

func startScraping(db *database.Queries, concurrency int, timeBetweenRequest time.Duration) {
	log.Printf("Scraping on %v goroutines every %v duration", concurrency, timeBetweenRequest)
	// To generate a request every passed in timeBetweenRequest we use the in-built ticker function
	ticker := time.NewTicker(timeBetweenRequest)
	// Empty for loop for immediate execution
	// ticker.C is the channel where the ticks(notification for every timeBetweenRequest) will be notified
	for ; ; <-ticker.C {

		feeds, err := db.GetNextFeedsToFetch(
			context.Background(),
			int32(concurrency),
		)

		if err != nil {
			log.Println("Error fetching feeds: %v", err)
			continue
		}

		// At a single time, many different goroutines are spawn up and given tasks to execute
		wg := &sync.WaitGroup{}

		for _, feed := range feeds {
			// Adds to the counter of wait group
			wg.Add(1)

			go scrapeFeed(db, wg, feed)
		}
		// Waits for every item in wait group to finish
		wg.Wait()
	}
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	//Will mark goroutine as done at the end
	defer wg.Done()

	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Error marking feed as fetched\n", err)
		return
	}

	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Println("Error in fetching url", err)
	}

	for _, item := range rssFeed.Channel.Item {
		description := sql.NullString{}
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}

		pubAt, err := time.Parse(time.RFC1123Z, item.PubDate)

		if err != nil {
			log.Printf("Couldn't parse date: [%v]. Err: %v\n", item.PubDate, err)
		}
		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Description: description,
			PublishedAt: pubAt,
			Url:         item.Link,
			FeedID:      feed.ID,
		})

		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			log.Printf("Couldn't create a new post with title : [%v]", item.Title)
		}
	}
	log.Printf("Feed %v Collected. %v posts Found", feed.Name, len(rssFeed.Channel.Item))
}
