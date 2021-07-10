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
var path, Title, URL, status, response string 
var port = "9222"
var id string = "750690488809422933"
var chromeOpen bool = false
var remote *godet.RemoteDebugger
var tickerDuration time.Duration = 1 * time.Second
var condition int = 0

var updateChannel chan int = make(chan int, 1)
/* 
1 - break ticker
2 - start ticker
*/
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
		
		lremote, err := godet.Connect("localhost:" + port, false)
		remote = lremote
		if err != nil {
			if !chromeOpen {
				openChrome()
				chromeOpen = true
			}
			return false
		}
		return true
	}
	if !init() {
		for remote == nil {
			if remote != nil {
				break
			}
			init()
			time.Sleep(1 * time.Second)
		}
	}
    
	if !IsXboxTab() {
		//args := ""
		openChrome()
	}

	err1 := client.Login(id)
	if err1 != nil {
		fmt.Println("RP client error: ", err1.Error())
		return
	}
    Ticker()	
}
	

func openChrome() {
	println("Opening chrome...")
	cmd := exec.Command(path, "https://www.xbox.com/en-US/play", " --remote-debugging-port=9222")
	err := cmd.Start()
	if err != nil {
		println("CMD error: ",err.Error())
		return
	}
}

func Ticker() {
	println("Ticker")
	ticker := time.NewTicker(tickerDuration)
	for range ticker.C {
		go IntegrityCheck()
		if len(updateChannel) > 0 {
			value := <-updateChannel
			switch value {
			case 1:
				//ticker.Stop()
				break
			}
		}
		Update()
	}
	
}



func IsXboxTab() bool {
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

func SetActivity(state string) {
	status = state
	now := time.Now()
	err := client.SetActivity(client.Activity{
		State:      state,
		Details:    "Xbox Cloud Gaming",
		LargeImage: "xboxlogo",
		LargeText:  "OblivCode",
		Timestamps: &client.Timestamps{
			Start: &now,
		},
	})
	
	if err != nil {
		println("RP Activity error: ", err.Error())
	}
	
}

func IntegrityCheck() {
	_, err := remote.Version(); if err != nil {
		updateChannel <- 1
		println("Fatal error: Remote not found")
		RestartPrompt()
	}
}

func RestartPrompt() {
	println("Would you like to restart the app? (y/n)(default: n)")
	fmt.Scanln(&response)
	response = strings.ToLower(response)
	if strings.HasPrefix(response, "y") {
		println("Restarting")
		procAttr := new(os.ProcAttr)
		procAttr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}
		os.StartProcess(os.Args[0], []string{}, procAttr)
	}
	os.Exit(0)
}
func Update() {
/*7 conditions 
println(Title)
println(URL)*/
		if (!IsXboxTab()) {
			println("no xbox tab")
			go client.Logout()
			
			for !IsXboxTab() {
				time.Sleep(time.Duration(3))
			}
			client.Login(id)
		} 
		if !strings.HasPrefix(URL, "https://www.xbox.com/") {return}
		if Title == "Xbox Cloud Gaming" {
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
			game = strings.Replace(game, "&amp;", "&", -1)
			game = strings.Replace(game, "&#39;", "'", -1)
			game = strings.Replace(game, "®", "", -1)
			//println(game)
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