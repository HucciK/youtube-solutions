const authEndpoint = "http://youtubesolutions.ru/api/authKey"

const { ipcRenderer, shell } = require("electron")
const preloader = document.querySelector("#preloader")
tryAuth()

document.querySelector("#auth-button").addEventListener("click", ()  => {
    let key = document.querySelector("#auth-key-input").value

    localStorage.setItem("key", key)
    tryAuth()
})

async function tryAuth() {
    let key = localStorage.getItem("key")
    if (key == null) {
        return
    }

    //запустить прелоадер

    preloader.style.display = "flex"
    setTimeout(() => {
        preloader.style.opacity = "1"
    }, 300)

   const response = await fetch(authEndpoint, {method: 'GET', headers:{"Authorization": key}});
   const data = await response.json();

   if (data.valid) {
       ipcRenderer.send("auth-success", data)
   } else {
       setTimeout(() => {
           preloader.style.opacity = "0"
       }, 300)
       document.querySelector("#auth-key-input").value = ""
       document.querySelector("#auth-key-input").placeholder = "Invalid key or logging IP"
       preloader.style.display = "none"
   }
}

let links = document.querySelectorAll(".links-content")
links.forEach(link => {
    link.addEventListener("click", (e) => {
        e.preventDefault()
        shell.openExternal(link.href)
    })
})