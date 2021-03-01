console.log("LOADED")

var searchTxt = document.querySelector('input.search-text')
var searchBtn = document.querySelector("a.search-btn")

function search() {
    searchBtn.href = '/search?' + searchTxt.value
}

searchBtn.onclick = search