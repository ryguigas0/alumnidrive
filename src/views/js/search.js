console.log("LOADED SEARCH BAR")

var searchTxt = document.querySelector('input.search-text')
var searchBtn = document.querySelector("a.search-btn")

function search() {
    searchBtn.href = '/search?id=' + searchTxt.value
}

searchBtn.onclick = search