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

var startTime time.Time

var sentimentGate float64

var banList []string

var configList []string

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
	sentimentGate, _ = strconv.ParseFloat(strings.Split(configData[1], ":")[1], 64)
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

func buildBanRows() []*g.TableRowWidget {
	fileLines, err := ioutil.ReadFile("ban.list")
	if err != nil {
		banList = nil
	} else {
		banList = strings.Split(string(fileLines), "\n")
	}
	rows := make([]*g.TableRowWidget, len(banList))
	for i := range rows {
		rows[i] = g.TableRow(
			g.Label(banList[i]).Wrapped(true),
		)
	}
	return rows
}

func buildConfigRows() []*g.TableRowWidget {
	fileLines, err := ioutil.ReadFile("config.txt")
	if err != nil {
		configList = nil
	} else {
		configList = strings.Split(string(fileLines), "\n")
	}
	rows := make([]*g.TableRowWidget, len(configList) + 2)

	rows[0] = g.TableRow(g.Label(configList[0]).Wrapped(true))
	rows[1] = g.TableRow(g.Label(configList[1]).Wrapped(true))
	rows[2] = g.TableRow(g.Label("Uptime: " + strconv.Itoa(int(time.Since(startTime).Seconds())) + " seconds"))
	rows[3] = g.TableRow(g.Label("Made by Habitual"))

	return rows
}


func loop() {

	w1Layout := g.Layout{
		g.Label("Clear Skies"),
		g.Table().FastMode(false).Rows(buildTweetRows()...).Columns(),
	}

	w2Layout := g.Layout{
		g.Label("User to block"),
		g.InputText(&name),
		g.Table().Rows(buildBanRows()...),
		g.Custom(func() {
			if g.IsKeyPressed(g.KeyEnter) {
				banList, _ := os.OpenFile("ban.list", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
				defer banList.Close()
				banList.Write([]byte(name + "\n"))
				name = ""
			}
		}),
	}

	w3Layout := g.Layout{
		g.Label("Running Configurations"),
		g.Table().Rows(buildConfigRows()...),

	}


	w1 := g.SingleWindow()
	tabLayout := g.TabBar().TabItems(
		g.TabItem("Twitter Feed").Layout(w1Layout),
		g.TabItem("Ban List").Layout(w2Layout),
		g.TabItem("Configuration").Layout(w3Layout),
	)

	w1.Layout(tabLayout)
}

func StartUi() {
	startTime = time.Now()
	go getTweets()
	wnd.Run(loop)
}
