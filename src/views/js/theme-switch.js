console.log("LOADED THEME SWITCHER")
var isDarkTheme = true
var body = document.querySelector("body")
var themeIcn = document.querySelector("#theme-icon")

function switchThemes() {
    if(isDarkTheme){
        body.className = "light-theme"
        themeIcn.src = "/frontend/lightbulb-off.svg"
    } else {
        body.className = "dark-theme"
        themeIcn.src = "/frontend/lightbulb-on.svg"
    }
    isDarkTheme = !isDarkTheme
}