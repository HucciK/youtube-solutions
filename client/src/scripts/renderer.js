
require("../scripts/settings")
require("../scripts/checker")
require("../scripts/proxy")

let userPhotoEndpoint;
const { ipcRenderer } = require("electron");
const server = "http://youtubesolutions.ru"
const authEndpoint = `${server}/api/authKey`

class Error {
    constructor(message) {
        this.message;
        this.name;
    }
}

class AuthError extends Error {
    constructor(message) {
        super(message);
        this.name = "AuthError";
    }
}

async function loadingSuccess() {
    const response = await fetch(authEndpoint, {method: 'GET', headers:{"Authorization": localStorage.getItem("key")}});
    const data = await response.json();

    if (!data.valid) {
        throw new AuthError("ERR")
    }
}

let update = document.querySelector("#update")
ipcRenderer.on("auth-success", async (e, data, version) => {

    try {
        await loadingSuccess()
    } catch (err) {
        ipcRenderer.send("exit-signal")
        return
    }

    document.querySelector("#license_info_news").innerText = `VERSION: ${version}`

    let userPhotoEndpoint = `${server}/api/getUserPhoto?id=${data.key_info.owner}`
    let userInfoEndpoint = `${server}/api/getUserInfo?id=${data.key_info.owner}`

    let userInfo = await doRequest(userInfoEndpoint)

    document.querySelector("#owner_avatar").src = userPhotoEndpoint
    document.querySelector("#owner_name").innerText = userInfo.name.toUpperCase()
    document.querySelector("#license_info_type").innerText = `LICENSE TYPE: ${data.key_info.type.toUpperCase()}`
    document.querySelector("#license_info_expire").innerText = `LICENSE EXPIRATION: ${data.key_info.expire.toUpperCase()}`

    if (data.hasUpdates) {
        update.style.fill = "#9cdec1"
    }
})

update.addEventListener("click",  () => {
    ipcRenderer.send("check-updates")
    update.style.animation = "infinite update .7s linear"
})

ipcRenderer.on("updates-checked", () => {
    update.style.animation = "none"
    update.style.fill = "dimgrey"
})

let content = document.querySelectorAll(".content_item")

let sidebar = document.querySelectorAll(".sidebar_btn")
sidebar.forEach(btn => {
    btn.addEventListener("click", (e) => changeContent(e))
})

function changeContent(e) {
    switch (e.target.id) {
        case "checker":
            activeTab = findActiveTab()
            activeTab.classList.add("hidden")

            removeSidebarActive()
            e.target.classList.add("sidebar_active")

            document.querySelector("#checker_setup").classList.remove("hidden")
            break;
    
        case "uploader":
            break;

        case "seo":
            break;

        case "proxy":
            activeTab = findActiveTab()
            activeTab.classList.add("hidden")

            removeSidebarActive()
            e.target.classList.add("sidebar_active")

            document.querySelector("#proxy_setup").classList.remove("hidden")
            break;  
            
        case "settings":
            activeTab = findActiveTab()
            activeTab.classList.add("hidden")

            removeSidebarActive()
            e.target.classList.add("sidebar_active")

            document.querySelector("#settings_tab").classList.remove("hidden")
            break;
        case "exit":
            ipcRenderer.send("exit-signal")
    }
}

function findActiveTab() {
    content.forEach(tab => {
        if (!tab.classList.contains("hidden")) {
            active = tab
        }
    })
    return active
}

function removeSidebarActive() {
    sidebar.forEach(btn => {
        btn.classList.remove("sidebar_active")
    })
}

async function doRequest(url) {
    let response = await fetch(url);
    return await response.json();
}

document.querySelector("#change_theme").addEventListener("click", () => {
    document.querySelector(":root").style.setProperty("--theme-color", "red")
})