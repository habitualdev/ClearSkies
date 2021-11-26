package main

import (
	"ClearSkies/ui"
	"errors"
	"fmt"
	"os"
)

func main(){


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

	ui.StartUi()
}