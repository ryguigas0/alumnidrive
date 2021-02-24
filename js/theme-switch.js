var isDarkTheme = true
var themeBtn = document.querySelector("button#toggle-theme")
var body = document.querySelector("body")
var themeIcn = document.querySelector("#theme-icon")

function switchThemes() {
    if(isDarkTheme){
        body.className = "light-theme"
        themeIcn.src = "../icons/lightbulb-off.svg"
        switchIcons("light")
    } else {
        body.className = "dark-theme"
        themeIcn.src = "../icons/lightbulb-on.svg"
        switchIcons("dark")
    }
    isDarkTheme = !isDarkTheme
}

function switchIcons(theme) {
    
}