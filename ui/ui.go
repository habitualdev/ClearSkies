package ui

import (
	"ClearSkies/scraper"
	"crypto/sha256"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
	"time"
)




var banList []string

var tweetCollector []string

func testHash(hash string, hashedHistory []string) bool{
	for _, testHash := range hashedHistory {
		if hash == testHash{
			return true
		}
	}
	return false
}
func StartUi() {
	a := app.New()
	w := a.NewWindow("Clear Skies")
	configLines, _ := ioutil.ReadFile("config.txt")
	configData := strings.Split(string(configLines),"\n")
	userName := (strings.Split(configData[0],":"))[1]
	sentimentGate, _ := strconv.ParseFloat(strings.Split(configData[1],":")[1],64)
	data := binding.BindStringList(&tweetCollector)
	tweetList := widget.NewListWithData(data,
		func() fyne.CanvasObject {
			widge := widget.NewLabel("")
			widge.Wrapping = fyne.TextTruncate
			tempsize := widge.Size()
			tempsize.Height = tempsize.Height*2
			widge.Resize(tempsize)
			return widge
		},func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		})

	startButton := widget.NewButton("Clear the air!", func() {
		var hashedHistory []string
		var printString string
			go func() {
				for {
					fileLines, err := ioutil.ReadFile("ban.list")
					if err != nil {
						banList = nil
					}else{
						banList = strings.Split(string(fileLines),"\n")
					}
					for n, existingLine := range tweetCollector {
						for _, badguy := range banList {
							if len(existingLine) > 4 {
								if matched, _ := regexp.MatchString(": " + badguy + " :", existingLine); matched {
									println(badguy, ":", existingLine)
									data.SetValue(n, "BANNED USER - " + badguy)
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
						data.Append(printString)
					}
				}
				time.Sleep(5 * time.Second)
			}
		}()
	})
	banEntry := widget.NewEntry()
	banEntry.OnSubmitted = func(string){
		existingBans, _ := ioutil.ReadFile("ban.list")
		ioutil.WriteFile("ban.list",[]byte(string(existingBans) + "\n" + banEntry.Text),0644)
		banEntry.SetText("")
	}
	w.SetContent(container.NewBorder(
		banEntry,
		startButton,
		nil,
		nil,
		tweetList,
	))
	w.Resize(fyne.Size{800,600})
	w.ShowAndRun()
}