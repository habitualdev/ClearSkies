package scraper

import (
	"context"
	twitterscraper "github.com/n0madic/twitter-scraper"
	"github.com/grassmudhorses/vader-go/lexicon"
	"github.com/grassmudhorses/vader-go/sentitext"
)

type TweetEntry struct{
	Text string
	Username string
	Sentiment float64
	Timestamp int64
}


func GetLatestTweets(user string, badList []string, sentimentGate float64) []TweetEntry{
	var tweetList []TweetEntry
	var printCheck bool
	scraper := twitterscraper.New()
	scraper.SetSearchMode(twitterscraper.SearchLatest)
	for tweet := range scraper.SearchTweets(context.Background(), "@" + user, 50) {
		printCheck = true
		if tweet.Username == user{
			continue
		}
		for _, badGuy := range badList{
			if tweet.Username == badGuy {
				printCheck = false
			}
		}
		if printCheck {
			parsedtext := sentitext.Parse(tweet.Text, lexicon.DefaultLexicon)
			sentiment := sentitext.PolarityScore(parsedtext)
			if sentiment.Compound > sentimentGate {
				tempTweet := TweetEntry{
					Text:      tweet.Text,
					Username:  tweet.Username,
					Sentiment: sentiment.Compound,
					Timestamp: tweet.Timestamp,
				}
				tweetList = append(tweetList, tempTweet)
			}
		}
	}
	return tweetList
}

