const shell = require('electron').shell

const FoundChannelType = "found channel"
const SavingValidType  = "saving valid"
const ErrorNotifyType  = "error notify"
const CheckStatusType  = "check status"

let checkPathField = document.querySelector("#check_path")
let savePathField = document.querySelector("#save_path")

let errors = 0;

window.addEventListener("paths-storage", setActualPath)

function setActualPath() {
    checkPath = localStorage.getItem("toCheckDir")
    savePath = localStorage.getItem("saveDir")

    if (checkPath == null || undefined) {
        checkPathField.innerText = "Check path: Path not specified"
    } else {
        checkPathField.innerText = "Check path: " + checkPath
    }

    if (savePath == null || undefined) {
        savePathField.innerText = "Save path: Path not specified"
    } else {
        savePathField.innerText = "Save path: " + savePath
    }
}

let start_btn = document.querySelector("#checker_start_button")
start_btn.addEventListener("click", () => {
    let err = "";

    let socket = new WebSocket("ws://localhost:8080/checker", ["token"])
    let results = document.querySelector("#checker_result_info_wrap")
    let workLog = document.querySelector("#checker_status")
    let checkStatus = document.querySelector("#checker_result")

    results.innerHTML = ``

    let selectedProxyList = getSelectedProxyList()
    b = {
        paths: {
            check_path: localStorage.getItem("toCheckDir"),
            save_path: localStorage.getItem("saveDir")
        },
        proxy: (selectedProxyList === undefined) ? null : selectedProxyList.proxies
    }

    if (b.paths.check_path === null || undefined) {
        checkPath = document.querySelector("#check_path")
        checkPath.style.color = "#ff7f7f"
        setTimeout(() => {
            checkPath.style.color = "var(--default-text-color)"
        }, 1000)

        return
    }

    if (b.paths.save_path === null || undefined) {
        savePath = document.querySelector("#save_path")
        savePath.style.color = "#ff7f7f"
        setTimeout(() => {
            savePath.style.color = "var(--default-text-color)"
        }, 1000)

        return;
    }

    socket.onopen = () => {
        workLog.innerText = "Working"
        socket.send(JSON.stringify(b))
    }

    socket.onmessage = (e) => {
        upd = (JSON.parse(e.data))
        if (upd.type === FoundChannelType){
            let monetize = "-"
            let brand = "No"
            let verified = "-"

            if (upd.data.geo === "") {
                upd.data.geo = "-"
            }

            if (upd.data.monetize === true) {
                monetize = "$"
            }

            if (upd.data.brand === true) {
                brand = "Yes"
            }

            if (upd.data.verified === true) {
                verified = `
                    <svg id="verified">
                        <use xlink:href="#verified_icon"></use>
                    </svg>`
            }

            results.innerHTML += `
                
                <div class="checker_result_info">
                  <div class="checker_data_section">
                    <div class="sections_content monetization">${monetize}</div>
                  </div>  
                  
                  <div class="checker_data_section">
                    <div class="section_content">${verified}</div>
                  </div> 
                  
                  <div class="checker_data_section">
                    <div class="section_content subs_count">${upd.data.subs_count}</div>
                  </div>

                  <div class="checker_data_section">
                    <div class="section_content views_count">${upd.data.views_count}</div>
                  </div>

                  <div class="checker_data_section">
                    <div class="section_content">${upd.data.reg_date}</div>
                  </div>

                  <div class="checker_data_section">
                    <div class="section_content">${brand}</div>
                  </div> 

                  <div class="checker_data_section">
                    <div class="section_content geo">${upd.data.geo}</div>
                  </div>

                  <div class="checker_data_section">
                    <div class="section_content"><a class="channel-url disabled-url" href="https://www.youtube.com/channel/${upd.data.id}">Loading..</a></div>
                  </div>
                </div>
            `
        }

        if (upd.type == CheckStatusType) {
            checkStatus.innerText = `${upd.data.valid}/${upd.data.errors}/${upd.data.checked}`
        }

        if (upd.type == SavingValidType) {
            workLog.innerText = "Saving valid..."
        }

        if (upd.type == ErrorNotifyType) {
            checkStatus.innerText = `${upd.data.valid}/${upd.data.errors}/${upd.data.checked}`
            console.log(upd)
            if (upd.error.Op === "Get") {
                workLog.innerText = "YouTube rate limits error"
            }

            if (upd.error === "error while saving valid") {
                workLog.innerText = "Error while saving valid"
            }
        }
    }

    socket.onerror = (e) => {
        err = e.type
    }

    socket.onclose = (e) => {
        if (err === "error") {
            workLog.innerText = "Something went wrong"
            return;
        }

        start_btn.disabled = false
        if (workLog.innerText === "Error while saving valid") {
            activateUrls()
            socket.close(1000, "finish")
            return;
        }

        if (errors !== 0) {
            workLog.innerText = "Saved with some errors"
            activateUrls()
            socket.close(1000, "finish")
            return;

        }

        workLog.innerText = "Successfully saved!"
        activateUrls()
        socket.close(1000, "finish")
    }

})

function activateUrls() {
    urls = document.querySelectorAll(".channel-url")
    urls.forEach((url) => {
        url.classList.remove("disabled-url")
        url.innerText = "Link"
        url.addEventListener("click", (e) => {
            e.preventDefault()
            shell.openExternal(url.href)
        })
    })
}

function getSelectedProxyList() {
    let selected;
    document.querySelectorAll(".proxy_list_content").forEach((list, i) => {
        if (list.classList.contains("selected")) {
           let storedLists = JSON.parse(localStorage.getItem("proxy_lists"))
           selected = storedLists[i]
        }
    })
    return selected
}

let sortButtons = document.querySelectorAll(".section_title_sort")
sortButtons.forEach(button => {
    button.addEventListener("click", (e) => {
        if (button.classList.contains("asc")) {
            switch (button.id) {
                case "views_sort":
                    button.classList.remove("asc")
                    button.classList.add("desc")
                    sortDesc(".views_count")
                    break
                case "subs_sort":
                    button.classList.remove("asc")
                    button.classList.add("desc")
                    sortDesc(".subs_count")
                    break
                case "geo_sort":
                    button.classList.remove("asc")
                    button.classList.add("desc")
                    sortDesc(".geo")
                    break
            }

            return
        }

        if (button.classList.contains("desc")) {
            switch (button.id) {
                case "views_sort":
                    button.classList.remove("desc")
                    button.classList.add("asc")
                    sortAsc(".views_count")
                    break
                case "subs_sort":
                    button.classList.remove("desc")
                    button.classList.add("asc")
                    sortAsc(".subs_count")
                    break
                case "geo_sort":
                    button.classList.remove("desc")
                    button.classList.add("asc")
                    sortAsc(".geo")
                    break
            }

        }

        switch (button.id){
            case "views_sort":
                button.classList.add("asc")
                sortAsc(".views_count")
                break
            case "subs_sort":
                button.classList.add("asc")
                sortAsc(".subs_count")
                break
            case "geo_sort":
                button.classList.remove("desc")
                button.classList.add("asc")
                sortAsc(".geo")
                break
        }

    })
})


function sortAsc(sortingField) {
    let parsed = [];

    if (sortingField === ".geo") {
        let resultRows = document.querySelectorAll(".checker_result_info");
        resultRows.forEach((row, i) => {
            let rowInfo = {
                html: row,
                count: row.querySelector(sortingField).innerText
            }

            parsed.push(rowInfo)
        })

        parsed.sort(function (a, b) {
            if (a.count.toLowerCase() < b.count.toLowerCase()) {
                return -1;
            }
            if (a.count.toLowerCase() > b.count.toLowerCase()) {
                return 1;
            }
            return 0;
        })

        let resultsField = document.querySelector("#checker_result_info_wrap")
        resultsField.innerHTML = ``
        parsed.forEach(data => {
            resultsField.append(data.html)
        })

        return;
    }

    let resultRows = document.querySelectorAll(".checker_result_info");
    resultRows.forEach((row, i) => {
        let rowInfo = {
            html: row,
            count: parseInt(row.querySelector(sortingField).innerText)
        }

        parsed.push(rowInfo)
    })

    if (parsed.length === 0) {
        return
    }


    parsed.sort(function (a, b) {
        if (a.count > b.count) {
            return -1
        }

        if (a.count < b.count) {
            return 1;
        }

        return 0;
    })

    let resultsField = document.querySelector("#checker_result_info_wrap")
    resultsField.innerHTML = ``
    parsed.forEach(data => {
        resultsField.append(data.html)
    })
}

function sortDesc(sortingField) {


    let parsed = [];

    if (sortingField === ".geo") {
        let resultRows = document.querySelectorAll(".checker_result_info");
        resultRows.forEach((row, i) => {
            let rowInfo = {
                html: row,
                count: row.querySelector(sortingField).innerText
            }

            parsed.push(rowInfo)
        })

        parsed.sort(function (a, b) {
            if (a.count.toLowerCase() < b.count.toLowerCase()) {
                return 1;
            }
            if (a.count.toLowerCase() > b.count.toLowerCase()) {
                return -1;
            }
            return 0;
        })

        let resultsField = document.querySelector("#checker_result_info_wrap")
        resultsField.innerHTML = ``
        parsed.forEach(data => {
            resultsField.append(data.html)
        })

        return;
    }

    let resultRows = document.querySelectorAll(".checker_result_info");
    resultRows.forEach((row, i) => {
        let rowInfo = {
            html: row,
            count: parseInt(row.querySelector(sortingField).innerText)
        }

        parsed.push(rowInfo)
    })

    if (parsed.length === 0) {
        return
    }

    parsed.sort(function (a, b) {
        if (a.count > b.count) {
            return 1
        }

        if (a.count < b.count) {
            return -1;
        }

        return 0;
    })

    let resultsField = document.querySelector("#checker_result_info_wrap")
    resultsField.innerHTML = ``
    parsed.forEach(data => {
        resultsField.append(data.html)
    })
}

setActualPath()