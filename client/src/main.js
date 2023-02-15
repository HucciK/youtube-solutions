const { app, BrowserWindow, autoUpdater, ipcMain, dialog } = require('electron');
const remoteMain = require('@electron/remote/main')
const exec = require('child_process').execFile
const path = require('path');

const server = 'http://youtubesolutions.ru'
const url = `${server}/devUpdate`
autoUpdater.setFeedURL({ url })

let local;
remoteMain.initialize()

if (require('electron-squirrel-startup')) {
  app.quit();
}

function startLocalhost() {
  localhost = exec(path.join(__dirname, '/cmd.exe'), function(err, data) {
  });
  local = localhost.pid
  console.log("LOCAL", local)
}

const createWindow = () => {
  const mainWindow = new BrowserWindow({
    width: 1000,
    height: 600,
    icon: __dirname + "/static/yt.ico",
    title: "YouTube Solutions",
    //frame: false,
    //transparent:true,
    //resizable: false,
    webPreferences: {
      contextIsolation: false,
      nodeIntegration: true,
    },
  });
  process.env.MAIN_WINDOW_ID = mainWindow.id;
  //mainWindow.removeMenu()
  mainWindow.webContents.openDevTools()

  mainWindow.loadFile(path.join(__dirname, '/templates/auth.html'));
  remoteMain.enable(mainWindow.webContents);
};

const createChildWindow = (w, h, title, path) => {
  const childWindow = new BrowserWindow({
    width: w,
    height: h,
    icon: __dirname + "/static/yt.ico",
    title: title,
    parent: getMainWindow(),
    //frame: false,
    //transparent:true,
    //resizable: false,
    webPreferences: {
      contextIsolation: false,
      nodeIntegration: true,
    },
  })
  //childWindow.removeMenu()

  childWindow.loadFile(path)
  remoteMain.enable(childWindow.webContents);
}

const getMainWindow = () => {
  const ID = process.env.MAIN_WINDOW_ID * 1;
  return BrowserWindow.fromId(ID)
}

app.on('ready', createWindow);

app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') {
    app.quit();
  }
});

app.on('activate', () => {
  if (BrowserWindow.getAllWindows().length === 0) {
    createWindow();
  }
});

ipcMain.on("select-check-dir", async (e, arg) => {
  let mainWindow = getMainWindow()
  let result = await dialog.showOpenDialog(mainWindow, {
    properties: ['openDirectory']
  })

  mainWindow.webContents.send("selected-check-dir", result.filePaths[0])
})

ipcMain.on("select-save-dir", async (e, arg) => {
  let mainWindow = getMainWindow()
  let result = await dialog.showOpenDialog(mainWindow, {
    properties: ['openDirectory']
  })

  mainWindow.webContents.send("selected-save-dir", result.filePaths[0])
})

ipcMain.on("auth-success", async (e, data) => {
  let mainWindow = getMainWindow()

  let version = app.getVersion()
  data.hasUpdates = version !== data.latest_version;

  console.log(data)

  await mainWindow.loadFile(path.join(__dirname, '/templates/index.html'));
  mainWindow.webContents.send("auth-success", data, version)
  startLocalhost()
})

ipcMain.on("add-proxy-list", () => {
  let w = 0
  BrowserWindow.getAllWindows().forEach((window) => {
    if (window.title === "Proxy") {
      w++
    }
  })

  if (w >= 1) {
    return
  }

  let file = path.join(__dirname, '/templates/proxy.html')

  createChildWindow(440, 560, "Proxy", file)
})

ipcMain.on("check-updates", async () => {
  let mainWindow = getMainWindow()

  autoUpdater.checkForUpdates();

  autoUpdater.addListener("error", (error) => {
    mainWindow.webContents.send("updates-checked")
  })

  autoUpdater.on("update-not-available", () => {
    mainWindow.webContents.send("updates-checked")
  })

  autoUpdater.on("update-available", () => {
    console.log("вывод окна с инфой об апдейте")
  })

  autoUpdater.on('update-downloaded', async (event, releaseNotes, releaseName) => {
    const dialogOpts = {
      type: 'info',
      buttons: ['Restart', 'Later'],

      title: 'Application Update',
      message: process.platform === 'win32' ? releaseNotes : releaseName,
      detail:
          'A new version has been downloaded. Restart the application to apply the updates.',
    }

    await dialog.showMessageBox(dialogOpts).then((returnValue) => {
      if (returnValue.response === 0) autoUpdater.quitAndInstall()
    })

    mainWindow.webContents.send("updates-checked")

  })
})

ipcMain.on("exit-signal", () => {
  try {
    process.kill(local)
  } catch (e) {
    app.quit()
    return
  }
  app.quit()
})
