window.System = window.System || {
  register: function () {},
  import: function () { return Promise.resolve({ assetManager: { downloader: { register: function () {} } } }); }
};
