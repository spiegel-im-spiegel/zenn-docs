package main

import (
	"fmt"
	"os"
	"time"

	"github.com/mmcdole/gofeed"
)

func main() {
	// ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	// go func() {
	// 	defer cancel()
	// 	_, err := gofeed.NewParser().ParseURLWithContext("https://zenn.dev/spiegel/feed", ctx)
	// 	if err != nil {
	// 		return
	// 	}
	// }()

	feed, err := gofeed.NewParser().ParseURL("https://zenn.dev/spiegel/feed")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	if feed == nil {
		fmt.Fprintln(os.Stderr, "No data")
		return
	}
	fmt.Println(feed.Title)
	fmt.Println(feed.FeedType, feed.FeedVersion)
	for _, item := range feed.Items {
		if item == nil {
			break
		}
		fmt.Println(item.Title)
		fmt.Println("\t->", item.Link)
		fmt.Println("\t->", item.PublishedParsed.Format(time.RFC3339))
	}
}
