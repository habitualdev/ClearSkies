package ui

import (
	"ClearSkies/scraper"
	"crypto/sha256"
	"fmt"
	g "github.com/AllenDang/giu"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var name string

var banList []string

var tweetCollector []string

var wnd = g.NewMasterWindow("Clear Skies", 800, 600, 0)

func testHash(hash string, hashedHistory []string) bool {
	for _, testHash := range hashedHistory {
		if hash == testHash {
			return true
		}
	}
	return false
}

func getTweets() {
	configLines, _ := ioutil.ReadFile("config.txt")
	configData := strings.Split(string(configLines), "\n")
	userName := (strings.Split(configData[0], ":"))[1]
	sentimentGate, _ := strconv.ParseFloat(strings.Split(configData[1], ":")[1], 64)
	var hashedHistory []string
	var printString string

	for {
		fileLines, err := ioutil.ReadFile("ban.list")
		if err != nil {
			banList = nil
		} else {
			banList = strings.Split(string(fileLines), "\n")
		}
		for n, existingLine := range tweetCollector {
			for _, badguy := range banList {
				if len(existingLine) > 4 {
					if matched, _ := regexp.MatchString(": "+badguy+" :", existingLine); matched {
						tweetCollector[n] = "BANNED USER - " + badguy
					}
				}
			}
		}
		stringList := scraper.GetLatestTweets(userName, banList, sentimentGate)
		for _, line := range stringList {
			formattedTime := time.Unix(line.Timestamp, 0)
			printString = formattedTime.Format("15:04:05") + " : " + line.Username + " : " + line.Text + "\n"
			tempHashSum := fmt.Sprintf("%x", sha256.Sum256([]byte(printString)))
			if !testHash(tempHashSum, hashedHistory) {
				hashedHistory = append(hashedHistory, tempHashSum)
				tweetCollector = append(tweetCollector, printString)
			}
		}
		time.Sleep(5 * time.Second)
	}
}

func writeToBan() {
	banList, _ := os.OpenFile("ban.list", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	defer banList.Close()
	banList.Write([]byte(name + "\n"))

}

func buildTweetRows() []*g.TableRowWidget {
	if len(tweetCollector) == 0 {
		return nil
	} else {
		rows := make([]*g.TableRowWidget, len(tweetCollector))
		for i := range rows {
			rows[i] = g.TableRow(
				g.Label(tweetCollector[i]).Wrapped(true),
			)
		}
		return rows
	}
}

func loop() {
	w1 := g.SingleWindow()

	w1Layout := g.Layout{
		g.Label("User to block"),
		g.InputText(&name).OnChange(writeToBan),
		g.Label("Clear Skies"),
		g.Table().FastMode(false).Rows(buildTweetRows()...).Columns(),
	}
	w1.Layout(w1Layout)
}

func StartUi() {
	go getTweets()
	wnd.Run(loop)
}
