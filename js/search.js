console.log("LOADED")

var searchTxt = document.querySelector('input.search-text')
var searchBtn = document.querySelector("a.search-btn")

function search() {
    console.log(searchTxt.value)
}

searchBtn.onclick = search