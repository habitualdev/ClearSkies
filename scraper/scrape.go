package scraper

import (
	"context"
	"github.com/grassmudhorses/vader-go/lexicon"
	"github.com/grassmudhorses/vader-go/sentitext"
	twitterscraper "github.com/n0madic/twitter-scraper"
	"sort"
	"strings"
)

type TweetEntry struct{
	Text string
	Username string
	Sentiment float64
	Timestamp int64
}

type listTemplate []TweetEntry

func (l listTemplate) Len() int {
	return len(l)
}

func (l listTemplate) Less(i, j int) bool {
	return l[i].Timestamp > l[j].Timestamp
}

func (l listTemplate) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
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
			parsedText := sentitext.Parse(tweet.Text, lexicon.DefaultLexicon)
			sentiment := sentitext.PolarityScore(parsedText)
			if sentiment.Compound > sentimentGate {
				tempTweet := TweetEntry{
					Text:      strings.ReplaceAll(tweet.Text,"\n"," "),
					Username:  tweet.Username,
					Sentiment: sentiment.Compound,
					Timestamp: tweet.Timestamp,
				}
				tweetList = append(tweetList, tempTweet)
			}
		}
	}
	sort.Sort(listTemplate(tweetList))

	return tweetList
}

