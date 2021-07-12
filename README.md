# xcloud-presence<br />
A Discord Rich Presence for the Xbox Cloud Gaming website.<br />
Using Raff's [godet](https://github.com/raff/godet) Remote client for Chrome DevTools, hugolgst's [rich-go](https://github.com/hugolgst/rich-go/) discord rich presence implementation and mitchellh's [go-ps](https://github.com/mitchellh/go-ps) library.<br />

#How to use:<br />
You only have to run the executable and it will automatically open the xbox cloud website.<br />
Add the path to your chrome browser into browser.txt if your chrome executable is not in it's default path.<br />
#Notices:<br />
-You will need to use the custom shortcuts for the app to work in the background.<br />
-If the rich presence is not appearing then close chrome and run a custom shortcut. The shortcuts open chrome with extra parameters so the app can work.<br />
#Dependencies:<br />
"https://github.com/raff/godet",<br />
"https://github.com/hugolgst/rich-go/client"<br />
"github.com/mitchellh/go-ps"
