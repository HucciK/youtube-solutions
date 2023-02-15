module.exports = {
  packagerConfig: {
    asar: true,
    icon: __dirname + "/src/static/yt",
  },
  rebuildConfig: {},
  makers: [
    {
      name: '@electron-forge/maker-squirrel',
      config: {
        //loadingGif: "build/icon2.gif",
        iconUrl: __dirname + "/src/static/yt.ico",
        setupIcon: __dirname + "/src/static/yt.ico",
      },
    },
    {
      name: '@electron-forge/maker-zip',
      platforms: ['darwin'],
    },
    {
      name: '@electron-forge/maker-deb',
      config: {
        icon: __dirname + "/src/static/yt.ico"
      },
    },
    {
      name: '@electron-forge/maker-rpm',
      config: {},
    },
  ],
};
