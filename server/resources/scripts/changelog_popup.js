const remote = require("@electron/remote");
const wnd = remote.getCurrentWindow();

document.querySelector("#popup_close").addEventListener("click", () => {
    wnd.close()
})