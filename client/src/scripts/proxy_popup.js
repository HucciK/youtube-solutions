const remote = require("@electron/remote")
const wnd = remote.getCurrentWindow()

document.querySelector("#popup_close").addEventListener("click", () => {
    wnd.close()
})

document.querySelector("#save_input").addEventListener("click", () => {
    let newProxyList = {
        name: document.querySelector("#proxy_name").value,
        proxies: []
    }
    let proxyStore = JSON.parse(localStorage.getItem("proxy_lists"))
    if (proxyStore === null) {
        proxyStore = []
        localStorage.setItem("proxy_lists", JSON.stringify(proxyStore))
    }

    let name = document.querySelector("#proxy_name")
    let input = document.querySelector("#proxy_input")
    proxyStrings = input.value.split("\n")
    proxyStrings.forEach(proxyString => {
        ip = proxyString.split("@")[0].split(":")[0]
        port = proxyString.split("@")[0].split(":")[1]
        user = proxyString.split("@")[1].split(":")[0]
        password = proxyString.split("@")[1].split(":")[1]

        let proxyObject = {
            Ip: ip,
            Port: port,
            User: user,
            Password: password,
        }
        newProxyList.proxies.push(proxyObject)
    })
    proxyStore.push(newProxyList)
    localStorage.setItem("proxy_lists", JSON.stringify(proxyStore))

    input.value = ""
    input.placeholder = "Successfully saved"

    name.value = ""
    name.placeholder = "LIST NAME"
})