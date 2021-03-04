console.log("LOADED SEARCH BAR")

var searchTxt = document.querySelector('input.search-text')
var searchBtn = document.querySelector("a.search-btn")

function search() {
    if (/[0-9]/.test(searchTxt.value)){
        searchBtn.href = '/search?id=' + searchTxt.value
    } else {
        searchBtn.href = '/search?name=' + searchTxt.value
    }
    
}

searchBtn.onclick = search