package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"os"
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
	Title       string  `xml:"title"`
	Link        string  `xml:"link"`
	Description string  `xml:"description"`
	Items       []Item  `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
	Description string `xml:"description"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("⚠️ .env 파일을 로드하지 못했어요:", err)
	}

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

	posts, _, err := client.Subreddit.NewPosts(context.Background(), "rstats", &reddit.ListOptions{
		Limit: 100,
	})
	if err != nil {
		panic(err)
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
			Title:       "r/rstats - 최근 24시간 글",
			Link:        "https://www.reddit.com/r/rstats/",
			Description: "Reddit r/rstats 서브레딧에서 최근 하루 동안 작성된 글들",
			Items:       items,
		},
	}

	file, err := os.Create("docs/rss.xml")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	enc := xml.NewEncoder(file)
	enc.Indent("", "  ")
	if err := enc.Encode(rss); err != nil {
		panic(err)
	}

	fmt.Println("✅ rss.xml 생성 완료!")
}
