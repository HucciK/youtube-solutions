const { ipcRenderer} = require("electron")

let pathSelect = document.querySelectorAll(".settings_path_select_button")
pathSelect.forEach(btn => {
    btn.addEventListener("click", (e) => {
        switch (e.target.id) {
            case "select_toCheck_path":
                ipcRenderer.send("select-check-dir")
                ipcRenderer.on("selected-check-dir", (e, result) => {
                    localStorage.setItem("toCheckDir", result)
                    window.dispatchEvent(new Event('paths-storage'))
                })
                break;
        
            case "select_save_path":
                ipcRenderer.send("select-save-dir")
                ipcRenderer.on("selected-save-dir", (e, result) => {
                    localStorage.setItem("saveDir", result)
                    window.dispatchEvent(new Event('paths-storage'))
                })
                break;
        }
    })
});