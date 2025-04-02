package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	"github.com/vartanbeno/go-reddit/v2/reddit"
)

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Items       []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
	Description string `xml:"description"`
}

func generateFeedFor(subreddit string, client *reddit.Client) error {
	posts, _, err := client.Subreddit.NewPosts(context.Background(), subreddit, &reddit.ListOptions{
		Limit: 100,
	})
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	oneDayAgo := now.Add(-24 * time.Hour)

	var items []Item
	for _, post := range posts {
		postTime := post.Created.Time
		if postTime.After(oneDayAgo) {
			items = append(items, Item{
				Title:       post.Title,
				Link:        "https://reddit.com" + post.Permalink,
				PubDate:     postTime.Format(time.RFC1123Z),
				Description: post.Body,
			})
		}
	}

	rss := RSS{
		Version: "2.0",
		Channel: Channel{
			Title:       fmt.Sprintf("r/%s - 최근 24시간 글", subreddit),
			Link:        fmt.Sprintf("https://www.reddit.com/r/%s/", subreddit),
			Description: fmt.Sprintf("Reddit r/%s 서브레딧에서 최근 하루 동안 작성된 글들", subreddit),
			Items:       items,
		},
	}

	outputPath := filepath.Join("docs", subreddit+".xml")
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := xml.NewEncoder(file)
	enc.Indent("", "  ")
	return enc.Encode(rss)
}

func main() {
	_ = godotenv.Load()

	credentials := reddit.Credentials{
		ID:       os.Getenv("REDDIT_CLIENT_ID"),
		Secret:   os.Getenv("REDDIT_CLIENT_SECRET"),
		Username: os.Getenv("REDDIT_USERNAME"),
		Password: os.Getenv("REDDIT_PASSWORD"),
	}
	userAgent := os.Getenv("REDDIT_USER_AGENT")

	client, err := reddit.NewClient(credentials, reddit.WithUserAgent(userAgent))
	if err != nil {
		panic(err)
	}

	subreddits := []string{"rstats", "rprogramming"}

	for _, sub := range subreddits {
		fmt.Printf("🔄 %s 피드 생성 중...\n", sub)
		if err := generateFeedFor(sub, client); err != nil {
			fmt.Printf("⚠️ %s 처리 중 에러: %v\n", sub, err)
		}
	}
}
