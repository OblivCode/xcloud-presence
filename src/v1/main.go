package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/hugolgst/rich-go/client"
	"github.com/raff/godet"
)

var path string 
var id string = "ApplicationID" //rich presence application ID
var Title string
var URL string
var chromeOpen bool = false
var remote *godet.RemoteDebugger
var gerr error = nil
var tickerDuration time.Duration = 1 * time.Second
var condition int = 0
var status string

func main() { 
	datBytes, err := ioutil.ReadFile("browser.txt")
	if err != nil || len(datBytes) < 2 {
		defaultpath := "C:/Program Files/Google/Chrome/Application/chrome.exe"
		path = defaultpath
		writBytes := []byte(defaultpath)
		go ioutil.WriteFile("browser.txt", writBytes, 0644)
	} else { path = string(datBytes) }
	println(path)
	init := func() bool {
		
		lremote, err := godet.Connect("localhost:9222", true)
		remote = lremote
		gerr = err
		if err != nil {
			if !chromeOpen {
				openChrome()
				chromeOpen = true
			}
			return false
		}
		
		return true
	}
    for !init() && gerr != nil {
		time.Sleep(5 * time.Second)
	}

	if !IsXboxTab() {
		//args := ""
		openChrome()
	}

	err1 := client.Login(id)
	if err1 != nil {
		fmt.Println("RP client error: ", err1)
		return
	}
    Ticker()	
}
	

func openChrome() {
	println("Opening chrome...")
	cmd := exec.Command(path, "https://www.xbox.com/en-US/play", " --remote-debugging-port=9222")
	err := cmd.Start()
	if err != nil {
		println("CMD error: ",err)
	}

}

func Ticker() {
	ticker := time.NewTicker(tickerDuration)
	for range ticker.C {
		Update()
	}
}

func IsXboxTab() bool { //whether there is an xbox tab open
	tabs, _ := remote.TabList("")
	    for _, element := range tabs {
	        if(strings.Contains(strings.ToLower(element.URL), "xbox")) {
			    //fmt.Println(element.Description)
				URL = element.URL
			    Title = element.Title
				return true
		    }
	    }
		return false
}

func SetActivity(state string) { //change presence 
	status = state
	now := time.Now()
	err := client.SetActivity(client.Activity{
		State:      state,
		Details:    "Xbox Cloud Gaming",
		LargeImage: "LargeAssetID",
		LargeText:  "OblivCode",
		Timestamps: &client.Timestamps{
			Start: &now,
		},
	})
	
	if err != nil {
		println("RP Activity error: ", err)
	}
}

func Update() { //check status of the open xbox tab
//7 conditions
println(Title)
println(URL)
		if (!IsXboxTab()) {
			time.Sleep(5 * time.Second)
			os.Exit(1)
			return
		} else if Title == "Xbox Cloud Gaming" {
			if condition == 1 { return} else {condition = 1}
			SetActivity("Browsing the game library")
			
		} else if Title == "Xbox Cloud Gaming (Beta) on Xbox.com" {
			if condition == 2 { return } else {condition = 2}
			SetActivity("Signing in")
		}  else if strings.Contains(URL, "search") {
			if condition == 5 { return} else {condition = 5}
			SetActivity("Searching for a game")
		} else if strings.Contains(Title, "|") {
			idx := strings.Index(Title, "|")
			game := Title[:idx]
			if strings.Contains(game, "&amp;") {game = strings.Replace(game, "&amp;", "&", -1)}
			if strings.Contains(game, "&#39;") {game = strings.Replace(game, "&#39;", "'", -1)}
			if strings.Contains(game, "®") {game = strings.Replace(game, "®", "", -1)}
			println(game)
			if strings.Contains(URL, "games") || strings.Contains(URL, "gallery/all-games") {
				if condition == 3 { 
					if strings.Contains(status, game) {
						return
					}
				} else {condition = 3}
				SetActivity("Viewing " + game)
			}else if strings.Contains(URL, "gallery") {
				if condition == 7 {return} else {condition = 7}
				
				SetActivity("Browsing the " + game + " section")
			}else if strings.Contains(URL, "launch") {
				if condition == 4 { return} else {condition = 4}
				SetActivity("Playing " + game)
			} 
						
		} else {
			if condition == 6 { return} else {condition = 6}
			SetActivity("Loading")
		}
}
