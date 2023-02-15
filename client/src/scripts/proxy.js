const { ipcRenderer } = require("electron")

const proxyTab = document.querySelector("#proxy_content")

let add = document.querySelector("#add_proxy")
add.addEventListener("click", addProxyList)

let remove = document.querySelector("#delete_proxy")
remove.addEventListener("click", deleteProxyList)

let test = document.querySelector("#proxy_speed_test")
test.addEventListener("click", speedTest)

window.addEventListener("storage", (e) => {
    if (e.key !== "proxy_lists") {
        return
    }
    updateLists()
})

function addProxyList() {
    ipcRenderer.send("add-proxy-list")
}

function deleteProxyList() {
    let proxies = JSON.parse(localStorage.getItem("proxy_lists"))
    let selectedProxyList = document.querySelector("#checker_proxy_list")
    let lists = document.querySelectorAll(".proxy_list_content")
    lists.forEach((list, i) => {
        if (list.classList.contains("selected")) {
            proxies.splice(i, 1)
            localStorage.setItem("proxy_lists", JSON.stringify(proxies))

            e = new Event("storage")
            e.key = "proxy_lists"
            window.dispatchEvent(e)
        }
    })
    disableProxyList()
    selectedProxyList.innerText = `Proxy not selected`
}

function disableProxyList() {
    let lists = document.querySelectorAll(".proxy_list_content")
    lists.forEach(list => {
        if (list.classList.contains("selected")) {
            list.classList.remove("selected")
        }
    })
}

function speedTest() {
    let selectedProxyList = getSelectedProxyList();

    if (selectedProxyList.list === undefined || selectedProxyList.render === undefined) {
        test.style.fill = "#ff4545"
        setTimeout(() => {
            test.style.fill = "#c2c2c2"
        }, 1500)
        return
    }

    body = {
        proxy_list: selectedProxyList.list.proxies
    }

    let socket = new WebSocket("ws://localhost:8080/proxy_checker", ["token"]);
    socket.onopen = () => {
        socket.send(JSON.stringify(body))
    }

    let row = 0;
    let proxyPing = selectedProxyList.render.querySelectorAll(".proxy_data > .proxy_ping")
    socket.onmessage = (e) => {
        upd = (JSON.parse(e.data))
        if (upd.data === null) {
            proxyPing[row].innerText = "NaN"
            row++
            return
        }

        proxyPing[row].innerText = upd.data
        row++
    }

    socket.onclose = () => {
        row = 0

        console.log("CLOSE")

        setTimeout(() => {
            proxyPing.forEach(row => {
                row.innerText = ``
            })
        }, 20000)
    }
}

function updateLists() {
    proxyTab.innerHTML = ``

    let proxyLists = JSON.parse(localStorage.getItem("proxy_lists"))
    let selectedProxyList = document.querySelector("#checker_proxy_list")
    if (proxyLists === null) {
        return
    }

    proxyLists.forEach(list => {

       proxyTab.insertAdjacentHTML("beforeend",
           `
            <div class="proxy_list_content">
            <div class="proxy_title">${list.name}</div>
             <div class="proxy_list">
                 
             </div>
            </div>`
       )

       let currentContentList = document.querySelectorAll(".proxy_list_content")[document.querySelectorAll(".proxy_list_content").length-1]
        currentContentList.addEventListener("click", (e) => {
           if (e.target.classList.contains("selected")) {
               disableProxyList()
               selectedProxyList.innerText = `Proxy not selected`
               return
           }
           disableProxyList()
            currentContentList.classList.toggle("selected")
           selectedProxyList.innerText = `Proxy ${currentContentList.querySelector(".proxy_title").innerText}`
       })

       let currentProxyList = document.querySelectorAll(".proxy_list")[document.querySelectorAll(".proxy_list").length-1]
       list.proxies.forEach(proxyString => {
           currentProxyList.insertAdjacentHTML("beforeend",
               `
                <div class="proxy_data">
                      <div class="proxy_ping"></div>
                      <div class="proxy_string">${proxyString.Ip}:${proxyString.Port}@${proxyString.User}:${proxyString.Password}</div>
                </div>
               `
           )
       })
    })
}

function getSelectedProxyList(){
    let render;
    let selected;
    document.querySelectorAll(".proxy_list_content").forEach((list, i) => {
        if (list.classList.contains("selected")) {
            let storedLists = JSON.parse(localStorage.getItem("proxy_lists"))
            render = list
            selected = storedLists[i]
        }
    })
    return {
        render: render,
        list: selected,
    }
}

updateLists()