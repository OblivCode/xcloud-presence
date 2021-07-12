# xcloud-presence<br />
A Discord Rich Presence for the Xbox Cloud Gaming website.<br />
Using Raff's [godet](https://github.com/raff/godet) Remote client for Chrome DevTools, hugolgst's [rich-go](https://github.com/hugolgst/rich-go/) discord rich presence implementation and mitchellh's [go-ps](https://github.com/mitchellh/go-ps) library.<br />

#How to use:<br />
You only have to run the executable and it will automatically open the xbox cloud website.<br />
Add the path to your chrome browser into browser.txt if your chrome executable is not in it's default path.<br />
#Notices:<br />
-You will need to restart the app if you close chrome.<br />
-If the rich presence is not appearing then close chrome and run the app. The app opens chrome with extra parameters for the app to work.<br />
#Dependencies:<br />
"https://github.com/raff/godet",<br />
"https://github.com/hugolgst/rich-go/client"<br />
"github.com/mitchellh/go-ps"
