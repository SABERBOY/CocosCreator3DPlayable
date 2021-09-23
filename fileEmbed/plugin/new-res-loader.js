// 新的资源载入方式脚本

/** 官网范例,反正看不懂
 * - https://developer.mozilla.org/zh-CN/docs/Web/API/WindowBase64/Base64_encoding_and_decoding#Solution_1_%E2%80%93_JavaScript's_UTF-16_%3E_base64
 */
function b64ToUint6(nChr) {
    return nChr > 64 && nChr < 91
        ? nChr - 65 : nChr > 96 && nChr < 123
            ? nChr - 71 : nChr > 47 && nChr < 58
                ? nChr + 4 : nChr === 43
                    ? 62 : nChr === 47
                        ? 63 : 0
}

/** 官网范例+1,看不懂+1,作用是将base64编码的字符串转为ArrayBuffer */
function base64DecToArr(sBase64, nBlockSize) {
    var sB64Enc = sBase64.replace(/[^A-Za-z0-9\+\/]/g, ""), nInLen = sB64Enc.length
    var nOutLen = nBlockSize ? Math.ceil((nInLen * 3 + 1 >>> 2) / nBlockSize) * nBlockSize : nInLen * 3 + 1 >>> 2
    var aBytes = new Uint8Array(nOutLen)
    for (var nMod3, nMod4, nUint24 = 0, nOutIdx = 0, nInIdx = 0; nInIdx < nInLen; nInIdx++) {
        nMod4 = nInIdx & 3
        nUint24 |= b64ToUint6(sB64Enc.charCodeAt(nInIdx)) << 18 - 6 * nMod4
        if (nMod4 === 3 || nInLen - nInIdx === 1) {
            for (nMod3 = 0; nMod3 < 3 && nOutIdx < nOutLen; nMod3++ , nOutIdx++) {
                aBytes[nOutIdx] = nUint24 >>> (16 >>> nMod3 & 24) & 255;
            }
            nUint24 = 0
        }
    }
    return aBytes
}

function uint6ToB64(nUint6) {

    return nUint6 < 26 ?
        nUint6 + 65
        : nUint6 < 52 ?
            nUint6 + 71
            : nUint6 < 62 ?
                nUint6 - 4
                : nUint6 === 62 ?
                    43
                    : nUint6 === 63 ?
                        47
                        :
                        65;

}

function base64EncArr(arraybuffer) {

    let aBytes = new Uint8Array(arraybuffer.buffer)

    var eqLen = (3 - (aBytes.length % 3)) % 3, sB64Enc = "";

    for (var nMod3, nLen = aBytes.length, nUint24 = 0, nIdx = 0; nIdx < nLen; nIdx++) {
        nMod3 = nIdx % 3;
        /* Uncomment the following line in order to split the output in lines 76-character long: */
        /*
        if (nIdx > 0 && (nIdx * 4 / 3) % 76 === 0) { sB64Enc += "\r\n"; }
        */
        nUint24 |= aBytes[nIdx] << (16 >>> nMod3 & 24);
        if (nMod3 === 2 || aBytes.length - nIdx === 1) {
            sB64Enc += String.fromCharCode(uint6ToB64(nUint24 >>> 18 & 63), uint6ToB64(nUint24 >>> 12 & 63), uint6ToB64(nUint24 >>> 6 & 63), uint6ToB64(nUint24 & 63));
            nUint24 = 0;
        }
    }

    return eqLen === 0 ?
        sB64Enc
        :
        sB64Enc.substring(0, sB64Enc.length - eqLen) + (eqLen === 1 ? "=" : "==");

}

function regLoading() {
    return new Promise(function (resolve, reject) {
        let ccc = System["import"]('cc')
        ccc.then(function (engine) {
            const downloadImage =
                (url/*: string*/, options/*: IDownloadParseOptions*/, onComplete/*: CompleteCallback*/) => {
                    console.log("downloadImage:", url)
                    const img = new Image()
                    img.src = "data:image/png;base64," + (window.res[url])   // 注意需要给base64编码添加前缀
                    img.onload = function () {
                        onComplete(null, img)
                    }
                };

            const downloadBlob =
                (/* url: string, options: IDownloadParseOptions, onComplete: CompleteCallback */) => {
                    // console.log("downloadBundle:", nameOrUrl)
                    // options.xhrResponseType = 'blob';
                    // downloadFile(url, options, options.onFileProgress, onComplete);
                };

            const downloadJson =
                (url/*: string*/, options/*: IDownloadParseOptions*/, onComplete/*: CompleteCallback<Record<string, any>>*/) => {
                    console.log("downloadJson:", url)
                    if (url.startsWith("./")) {
                        url = url.substring(2)
                    }
                    onComplete(null, window.res[url])
                };

            const downloadArrayBuffer =
                (url/*: string*/, options/*: IDownloadParseOptions*/, onComplete/*: CompleteCallback*/) => {
                    console.log("downloadArrayBuffer:", url)
                    const arr = base64DecToArr(JSON.stringify(window.res[url]))
                    onComplete(null, arr)
                };

            /*            var __audioSupport = cc.sys.__audioSupport;
                        console.log("__audioSupport:",cc.sys)
                        var formatSupport = __audioSupport.format;
                        var context = __audioSupport.context;*/

            function base64toBlob(base64, type) {
                var bstr = atob(base64, type),
                    n = bstr.length,
                    u8arr = new Uint8Array(n);
                while (n--) {
                    u8arr[n] = bstr.charCodeAt(n)
                }
                return new Blob([u8arr], {
                    type: type,
                })
            }

            function base64toArray(base64) {
                var bstr = atob(base64),
                    n = bstr.length,
                    u8arr = new Uint8Array(n);
                while (n--) {
                    u8arr[n] = bstr.charCodeAt(n)
                }

                return u8arr;
            }

            const loadDomAudio = (url, onComplete) => {
                var dom = document.createElement('audio');
                dom.muted = true;
                dom.muted = false;
                var data = window.res[url.split("?")[0]];
                data = base64toBlob(data, "audio/mpeg");
                if (window.URL) {
                    dom.src = window.URL.createObjectURL(data);
                } else {
                    dom.src = data;
                }
                dom["url"] = dom.src
                var clearEvent = function () {
                    clearTimeout(timer);
                    dom.removeEventListener("canplaythrough", success, false);
                    dom.removeEventListener("error", failure, false);
                    /*if (__audioSupport.USE_LOADER_EVENT)
                        dom.removeEventListener(__audioSupport.USE_LOADER_EVENT, success, false);*/
                };

                var timer = setTimeout(function () {
                    if (dom.readyState === 0)
                        failure();
                    else
                        success();
                }, 8000);

                var success = function () {
                    clearEvent();
                    onComplete && onComplete(null, dom);
                };

                var failure = function () {
                    clearEvent();
                    var message = 'load audio failure - ' + url;
                    cc.log(message);
                    onComplete && onComplete(new Error(message));
                };

                dom.addEventListener("canplaythrough", success, false);
                dom.addEventListener("error", failure, false);
                /*if (__audioSupport.USE_LOADER_EVENT)
                    dom.addEventListener(__audioSupport.USE_LOADER_EVENT, success, false);*/
                return dom;
            }

            const loadWebAudio = (url, onComplete) => {
                let context = new (window.AudioContext || window.webkitAudioContext);
                if (!context) onComplete(new Error('Audio Downloader: no web audio context.'), null);
                var data = window.res[url];
                data = base64toArray(data);
                if (data) {
                    context["decodeAudioData"](data.buffer, function (buffer) {
                        onComplete(null, buffer);
                    }, function () {
                        onComplete('decode error - ' + url, null);
                    });
                } else {
                    onComplete('request error - ' + url, null);
                }
            }

            const downloadAudio = (url, options, onComplete) => {
                // loadWebAudio(url, onComplete);
                loadDomAudio(url, onComplete);
                /*if (formatSupport.length === 0) {
                    return new Error('Audio Downloader: audio not supported on this browser!');
                }

                // If WebAudio is not supported, load using DOM mode
                if (!__audioSupport.WEB_AUDIO) {
                    loadDomAudio(url, onComplete);
                }
                else {
                    loadWebAudio(url, onComplete);
                }*/
            }

            const downloadText =
                (url/*: string*/, options/*: IDownloadParseOptions*/, onComplete/*: CompleteCallback*/) => {
                    console.log("downloadText:", url)
                    onComplete(null, JSON.stringify(window.res[url]))
                };
            const REGEX = /^(?:\w+:\/\/|\.+\/).+/;

            const getScriptById = (id) => {
                return document.getElementById(id).innerText;
            }

            const downloadBundle =
                (nameOrUrl/*: string*/, options/*: IBundleOptions*/, onComplete/*: CompleteCallback*/) => {
                    console.log("downloadBundle:", nameOrUrl)
                    // onComplete(null, null)
                    let url = nameOrUrl
                    if (!REGEX.test(url)) {
                        url = `assets/${url}`;
                    }
                    const config = `${url}/config.json`
                    let count = 0;
                    let error = null;
                    let out = null;
                    downloadJson(config, options, (err, response) => {
                        error = err
                        out = response
                        if (out) {
                            out.base = `${url}/`;
                        }
                        if (++count === 2) {
                            onComplete(error, out);
                        }
                    })
                    const jspath = `${url}/index.js`;
                    downloadScript(jspath, options, (err) => {
                        error = err;
                        if (++count === 2) {
                            onComplete(err, out);
                        }
                    });
                };
            const downloaded = {};
            const downloadScript = (
                url/*: string*/,
                options/*: IBundleOptions*/,
                onComplete/*: CompleteCallback,*/
            ) =>/*: HTMLScriptElement | null*/ {
                console.log("downloadScript:", url)
                if (url.search("main") !== -1) {
                    System.import("chunks:///_virtual/main").then(() => {
                        console.log("SUCCESS!!")
                        onComplete(null);
                    });
                } else {
                    System.import("chunks:///_virtual/resources").then(() => {
                        console.log("SUCCESS!!")
                        onComplete(null);
                    });
                }
                return;
                // no need to load script again
                /* if (downloaded[url]) {
                     if (onComplete) {
                         onComplete(null);
                     }
                     return null;
                 }
                 const ps = {}
                 ps[url.toString()] = "./" + url.toString()
                 // ps["PS"] = "./"+url.toString()
                 var importMapJson = {"imports": ps}
                 console.log("href:", window.location.href)
                 const importMapAbsUrl = new URL(url, window.location.href).href;
                 console.log("importMap:", importMapJson)
                 console.log("importMapAbsUrl:", importMapAbsUrl)
                 // System.constructor.prototype.patches.setImportMap(importMapJson, importMapAbsUrl);
                 // System.import(ps[url.toString()])
                 console.log("INDEX:", getScriptById("assets/main/index.js"))
                 return eval(getScriptById("assets/main/index.js"));*/
            }

            const getFontFamily = (fontHandle) => {
                var ttfIndex = fontHandle.lastIndexOf(".ttf");
                if (ttfIndex === -1) {
                    ttfIndex = fontHandle.lastIndexOf(".tmp");
                }
                if (ttfIndex === -1) return fontHandle;

                var slashPos = fontHandle.lastIndexOf("/");
                var fontFamilyName;
                if (slashPos === -1) {
                    fontFamilyName = fontHandle.substring(0, ttfIndex) + "_LABEL";
                } else {
                    fontFamilyName = fontHandle.substring(slashPos + 1, ttfIndex) + "_LABEL";
                }
                return fontFamilyName;
            }

            const downloadFont = (url, options, onComplete) => {
                console.log("downloadFont:", url)
                let fontFamilyName = getFontFamily(url);
                let data = "url(data:application/x-font-woff;charset=utf-8;base64,PASTE-BASE64-HERE) format(\"woff\")";
                data = data.replace("PASTE-BASE64-HERE", window.res[url]);

                let fontFace = new FontFace(fontFamilyName, data);
                document.fonts.add(fontFace);

                fontFace.load();
                fontFace.loaded.then(function () {
                    onComplete(null, fontFamilyName);
                }, function () {
                    cc.warnID(4933, fontFamilyName);
                    onComplete(null, fontFamilyName);
                });
            }
            const downloadCCON = (url/*: string*/, options/*: IDownloadParseOptions*/, onComplete/*: CompleteCallback<CCON>*/) => {
                downloadJson(url, options, (err, json) => {
                    if (err) {
                        onComplete(err);
                        return;
                    }
                    const cconPreface = parseCCONJson(json);
                    const chunkPromises = Promise.all(cconPreface.chunks.map((chunk) => new Promise < Uint8Array > ((resolve, reject) => {
                        downloadArrayBuffer(`${/*path.*/mainFileName(url)}${chunk}`, {}, (errChunk, chunkBuffer) => {
                            if (err) {
                                reject(err);
                            } else {
                                resolve(new Uint8Array(chunkBuffer));
                            }
                        });
                    })));
                    chunkPromises.then((chunks) => {
                        const ccon = new CCON(cconPreface.document, chunks);
                        onComplete(null, ccon);
                    }).catch((err) => {
                        onComplete(err);
                    });
                });
            };

            const downloadCCONB = (url/*: string*/, options/*: IDownloadParseOptions*/, onComplete/*: CompleteCallback<CCON>*/) => {
                downloadArrayBuffer(url, options, (err, arrayBuffer/*: ArrayBuffer*/) => {
                    if (err) {
                        onComplete(err);
                        return;
                    }
                    try {
                        const ccon = decodeCCONBinary(new Uint8Array(arrayBuffer));
                        onComplete(null, ccon);
                    } catch (err) {
                        onComplete(err);
                    }
                });
            };

            const parseCCONJson = (json/*: unknown*/) => {
                const cconPreface = json /*as CCONPreface*/;

                return {
                    chunks: cconPreface.chunks,
                    document: cconPreface.document,
                };
            }

            const mainFileName = (fileName/*: string*/) => {
                if (fileName) {
                    const idx = fileName.lastIndexOf('.');
                    if (idx !== -1) {
                        return fileName.substring(0, idx);
                    }
                }
                return fileName;
            }

            class CCON {
                constructor(document/*: unknown*/, chunks/*: Uint8Array[]*/) {
                    this._document = document;
                    this._chunks = chunks;
                }

                get document() {
                    return this._document;
                }

                get chunks() {
                    return this._chunks;
                }

                /*private*/
                _document/*: unknown*/;
                /*private*/
                _chunks/*: Uint8Array[]*/;
            }

            const VERSION = 1;

            const MAGIC = 0x4E4F4343;

            const CHUNK_ALIGN_AS = 8;
            const decodeCCONBinary = (bytes/*: Uint8Array*/) => {
                if (bytes.length < 16) {
                    throw new Error("13102");
                }

                const dataView = new DataView(
                    bytes.buffer,
                    bytes.byteOffset,
                    bytes.byteLength,
                );

                const magic = dataView.getUint32(0, true);
                if (magic !== MAGIC) {
                    throw new Error("13100");
                }

                const version = dataView.getUint32(4, true);
                if (version !== VERSION) {
                    throw new Error("13101"/*, version*/);
                }

                const dataByteLength = dataView.getUint32(8, true);
                if (dataByteLength !== dataView.byteLength) {
                    throw new Error("13102");
                }

                let chunksStart = 12;

                const jsonDataLength = dataView.getUint32(chunksStart, true);
                chunksStart += 4;
                const jsonData = new Uint8Array(dataView.buffer, chunksStart + dataView.byteOffset, jsonDataLength);
                chunksStart += jsonDataLength;
                const jsonString = decodeJson(jsonData);
                let json/*: unknown*/;
                try {
                    json = JSON.parse(jsonString);
                } catch (err) {
                    throw new Error(err);
                }

                const chunks/*: Uint8Array[]*/ = [];
                while (chunksStart < dataView.byteLength) {
                    if (chunksStart % CHUNK_ALIGN_AS !== 0) {
                        const padding = CHUNK_ALIGN_AS - chunksStart % CHUNK_ALIGN_AS;
                        chunksStart += padding;
                    }
                    const chunkDataLength = dataView.getUint32(chunksStart, true);
                    chunksStart += 4;
                    chunks.push(new Uint8Array(dataView.buffer, chunksStart + dataView.byteOffset, chunkDataLength));
                    chunksStart += chunkDataLength;
                }

                if (chunksStart !== dataView.byteLength) {
                    throw new Error("13102");
                }

                return new CCON(json, chunks);
            }

            const decodeJson = (data/*: Uint8Array*/) => {
                if (typeof TextDecoder !== 'undefined') {
                    return new TextDecoder().decode(data);
                } else if ('Buffer' in globalThis) {
                    const {Buffer} = (globalThis/* as unknown as { Buffer: BufferConstructor }*/);
                    // eslint-disable-next-line no-buffer-constructor
                    return Buffer.from(data.buffer, data.byteOffset, data.byteLength).toString();
                } else {
                    throw new Error("13104");
                }
            }

            let downloaders = {
                // Images
                '.png': downloadImage,
                '.jpg': downloadImage,
                '.bmp': downloadImage,
                '.jpeg': downloadImage,
                '.gif': downloadImage,
                '.ico': downloadImage,
                '.tiff': downloadImage,
                '.webp': downloadImage,
                '.image': downloadImage,
                '.pvr': downloadArrayBuffer,
                '.pkm': downloadArrayBuffer,
                '.astc': downloadArrayBuffer,

                // Txt
                '.txt': downloadText,
                '.xml': downloadText,
                '.vsh': downloadText,
                '.fsh': downloadText,
                '.atlas': downloadText,
                '.tmx': downloadText,
                '.tsx': downloadText,
                '.json': downloadJson,
                '.ExportJson': downloadJson,
                '.plist': downloadText,
                '.fnt': downloadText,
                '.zip': downloadText,// 按需添加

                // Font
                '.ttf': downloadFont,
                // MP3
                '.mp3': downloadAudio,
                '.ogg': downloadAudio,
                '.wav': downloadAudio,
                '.m4a': downloadAudio,

                '.ccon': downloadCCON,
                '.cconb': downloadCCONB,

                // Binary
                '.binary': downloadArrayBuffer,
                '.bin': downloadArrayBuffer,
                '.dbbin': downloadArrayBuffer,
                '.skel': downloadArrayBuffer,

                '.js': downloadScript,

                bundle: downloadBundle,

                default: downloadText,
            };
            for (const key in downloaders) {
                if (Object.prototype.hasOwnProperty.call(downloaders, key)) {
                    const element = downloaders[key];
                    engine.assetManager.downloader.register(key, element);
                }
            }
            resolve(engine)
        })
    })
}