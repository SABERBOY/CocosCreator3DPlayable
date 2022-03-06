System.register(["application.js"], function (_export, _context) {
  "use strict";

    var createApplication, canvas, $p, bcr;

    function loadJsListFile(url) {
        return new Promise(function (resolve, reject) {
            resolve();
        });
    }

    function fetchWasm(url) {
        return url;
        // return fetch(url).then(function (response) {
        //   return response.arrayBuffer();
        // });
        console.log("url:",url);
        url="cocos-js/assets/"+cc.path.basename(url)
        return base64DecToArr(window.res[url])
    }

    function findCanvas() {
        // Use canvas in outer context
        if (!canvas || canvas.tagName !== 'CANVAS') {
            console.error("unknown canvas id:", el);
        }

        var width = canvas.width;
        var height = canvas.height;
        var container = document.createElement('div');

        if (canvas && canvas.parentNode) {
            canvas.parentNode.insertBefore(container, canvas);
        }

        container.setAttribute('id', 'Cocos3dGameContainer');
        container.appendChild(canvas);
        var frame = container.parentNode === document.body ? document.documentElement : container.parentNode;
        addClass(canvas, 'gameCanvas');
        canvas.setAttribute('width', width || '480');
        canvas.setAttribute('height', height || '320');
        canvas.setAttribute('tabindex', '99');
        return {
            frame: frame,
            canvas: canvas,
            container: container
        };
    }

    function addClass(element, name) {
        var hasClass = (' ' + element.className + ' ').indexOf(' ' + name + ' ') > -1;

        if (!hasClass) {
            if (element.className) {
                element.className += ' ';
            }

            element.className += name;
        }
    }

    return {
        setters: [function (_applicationJs) {
            createApplication = _applicationJs.createApplication;
        }],
        execute: function () {
            canvas = document.getElementById('GameCanvas');
            $p = canvas.parentElement;
            bcr = $p.getBoundingClientRect();
            canvas.width = bcr.width;
            canvas.height = bcr.height;
            createApplication({
                loadJsListFile: loadJsListFile,
                fetchWasm: fetchWasm
            }).then(function (application) {
                return application.start({
                    findCanvas: findCanvas
                });
            })["catch"](function (err) {
                console.error(err);
            });
        }
    };
});