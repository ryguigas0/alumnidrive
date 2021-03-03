var isDarkTheme = true
var body = document.querySelector("body")
var themeIcn = document.querySelector("#theme-icon")

function switchThemes() {
    if(isDarkTheme){
        body.className = "light-theme"
        themeIcn.src = "/icons/lightbulb-off.svg"
    } else {
        body.className = "dark-theme"
        themeIcn.src = "/icons/lightbulb-on.svg"
    }
    isDarkTheme = !isDarkTheme
}