package main

import (
	"ClearSkies/scraper"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func main(){
	var userName string
	var sentimentGate float64

	if _, err := os.Stat("config.txt"); errors.Is(err, os.ErrNotExist) {
		fmt.Println("config.txt does not exist. Creating file...")
		fmt.Println("Enter twitter handle: ")
		var twitterHandle string
		fmt.Scanln(&twitterHandle)
		if file, fileErr :=os.OpenFile("config.txt", os.O_CREATE|os.O_RDWR, 0644); fileErr != nil{
			fmt.Println("Could not create file. Exiting...")
			os.Exit(1)
		} else{
			file.Write([]byte("username:"+twitterHandle+"\nsentiment:-.5"))
			file.Close()
			fmt.Println("Empty config.txt created. Modify as needed.")
			fmt.Println("NOTE: Sentiment is a score from 1 to -1 decribing the positivity of the tweet.")
			fmt.Println("NOTE: A sentiment of -.5 is a good baseline to screen for obscene comments")
		}
	}
	configLines, err := ioutil.ReadFile("config.txt")
	configData := strings.Split(string(configLines),"\n")
	userName = (strings.Split(configData[0],":"))[1]
	sentimentGate, _ = strconv.ParseFloat(strings.Split(configData[1],":")[1],64)
	fileLines, err := ioutil.ReadFile("ban.list")
	if err != nil {
		println(scraper.GetLatestTweets(userName, nil, sentimentGate))
		}else{
			banList := strings.Split(string(fileLines),"\n")
			println(scraper.GetLatestTweets(userName, banList, sentimentGate))
	}
}