var opts = document.querySelector("select#addtype")
var dirOpts = document.querySelector("#dir-opts")
var fileOpts = document.querySelector("#file-opts")

function getOpts() {
    if (opts.value == "dir") {
        dirOpts.style.visibility = "visible"
        fileOpts.style.visibility = "hidden"
    } else {
        dirOpts.style.visibility = "hidden"
        fileOpts.style.visibility = "visible"
    }
}