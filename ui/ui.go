package main

import (
	"ClearSkies/scraper"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

func main() {
	a := app.New()
	w := a.NewWindow("Hello")
	configLines, _ := ioutil.ReadFile("config.txt")
	configData := strings.Split(string(configLines),"\n")
	userName := (strings.Split(configData[0],":"))[1]
	sentimentGate, _ := strconv.ParseFloat(strings.Split(configData[1],":")[1],64)

	data := binding.BindStringList(&[]string{})


	hello := widget.NewListWithData(data,
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		})
	hello.Resize(fyne.Size{
		Width:  ,
		Height: 0,
	})
	w.SetContent(container.NewVBox(
		hello,
		widget.NewButton("Clear the air!", func() {
			var printString string
			stringList := scraper.GetLatestTweets(userName,nil,sentimentGate)
			for _, line := range stringList{
				formattedTime := time.Unix(line.Timestamp,0)
				printString = formattedTime.Format("15:04:05") + " : " + line.Username + " : " + line.Text + "\n"
				data.Append(printString)
			}

		}),
	))

	w.ShowAndRun()
}