(function () {
  function cleanUrl(url) {
    url = String(url || "");
    var q = url.indexOf("?");
    if (q >= 0) url = url.slice(0, q);
    if (url.indexOf("./") === 0) url = url.slice(2);
    return url;
  }

  function readResource(url) {
    return window.res && window.res[cleanUrl(url)];
  }

  function base64ToArray(base64) {
    var raw = atob(base64 || "");
    var out = new Uint8Array(raw.length);
    for (var i = 0; i < raw.length; i++) out[i] = raw.charCodeAt(i);
    return out;
  }

  function downloadJson(url, options, onComplete) {
    onComplete(null, readResource(url));
  }

  function downloadText(url, options, onComplete) {
    var value = readResource(url);
    onComplete(null, typeof value === "string" ? value : JSON.stringify(value));
  }

  function downloadArrayBuffer(url, options, onComplete) {
    onComplete(null, base64ToArray(readResource(url)).buffer);
  }

  function downloadImage(url, options, onComplete) {
    var image = new Image();
    image.onload = function () { onComplete(null, image); };
    image.onerror = function () { onComplete(new Error("load image failed: " + url)); };
    image.src = "data:image/png;base64," + readResource(url);
  }

  function downloadAudio(url, options, onComplete) {
    var audio = document.createElement("audio");
    var blob = new Blob([base64ToArray(readResource(url))], { type: "audio/mpeg" });
    audio.src = URL.createObjectURL(blob);
    audio.oncanplaythrough = function () { onComplete(null, audio); };
    audio.onerror = function () { onComplete(new Error("load audio failed: " + url)); };
  }

  function downloadFont(url, options, onComplete) {
    var name = cleanUrl(url).split("/").pop().replace(/\.(ttf|tmp)$/i, "") + "_LABEL";
    var source = "url(data:application/x-font-woff;charset=utf-8;base64," + readResource(url) + ") format(\"woff\")";
    var font = new FontFace(name, source);
    document.fonts.add(font);
    font.load().then(function () { onComplete(null, name); }, function () { onComplete(null, name); });
  }

  function downloadBundle(nameOrUrl, options, onComplete) {
    var url = /^(?:\w+:\/\/|\.+\/).+/.test(nameOrUrl) ? nameOrUrl : "assets/" + nameOrUrl;
    downloadJson(url + "/config.json", options, function (err, config) {
      if (config) config.base = url + "/";
      if (err) return onComplete(err);
      downloadScript(url + "/index.js", options, function (scriptErr) {
        onComplete(scriptErr, config);
      });
    });
  }

  function downloadScript(url, options, onComplete) {
    var moduleId = cleanUrl(url).indexOf("resources") >= 0 ? "chunks:///_virtual/resources" : "chunks:///_virtual/main";
    System.import(moduleId).then(function () { onComplete(null); }, onComplete);
  }

  function mainFileName(fileName) {
    var idx = fileName.lastIndexOf(".");
    return idx === -1 ? fileName : fileName.substring(0, idx);
  }

  function downloadCCON(url, options, onComplete) {
    downloadJson(url, options, function (err, json) {
      if (err) return onComplete(err);
      var preface = cc.internal.parseCCONJson(json);
      Promise.all(preface.chunks.map(function (chunk) {
        return new Promise(function (resolve, reject) {
          downloadArrayBuffer(mainFileName(url) + chunk, {}, function (chunkErr, buffer) {
            if (chunkErr) reject(chunkErr);
            else resolve(new Uint8Array(buffer));
          });
        });
      })).then(function (chunks) {
        onComplete(null, new cc.internal.CCON(preface.document, chunks));
      }, onComplete);
    });
  }

  function downloadCCONB(url, options, onComplete) {
    downloadArrayBuffer(url, options, function (err, buffer) {
      if (err) return onComplete(err);
      try {
        onComplete(null, cc.internal.decodeCCONBinary(new Uint8Array(buffer)));
      } catch (decodeErr) {
        onComplete(decodeErr);
      }
    });
  }

  window.regLoading = function () {
    return System.import("cc").then(function (engine) {
      var downloaders = {
        ".png": downloadImage, ".jpg": downloadImage, ".jpeg": downloadImage, ".bmp": downloadImage,
        ".gif": downloadImage, ".webp": downloadImage, ".pvr": downloadArrayBuffer, ".pkm": downloadArrayBuffer,
        ".astc": downloadArrayBuffer, ".txt": downloadText, ".xml": downloadText, ".vsh": downloadText,
        ".fsh": downloadText, ".atlas": downloadText, ".tmx": downloadText, ".tsx": downloadText,
        ".json": downloadJson, ".ExportJson": downloadJson, ".plist": downloadText, ".fnt": downloadText,
        ".ttf": downloadFont, ".mp3": downloadAudio, ".ogg": downloadAudio, ".wav": downloadAudio,
        ".m4a": downloadAudio, ".ccon": downloadCCON, ".cconb": downloadCCONB, ".binary": downloadArrayBuffer,
        ".bin": downloadArrayBuffer, ".dbbin": downloadArrayBuffer, ".skel": downloadArrayBuffer,
        ".wasm": downloadArrayBuffer, ".js": downloadScript, "bundle": downloadBundle, "default": downloadText
      };
      Object.keys(downloaders).forEach(function (key) {
        engine.assetManager.downloader.register(key, downloaders[key]);
      });
      return engine;
    });
  };
})();
