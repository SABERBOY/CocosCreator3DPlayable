# CocosCreator3D PlayableAD

### [Modified based on this project](https://github.com/Jecced/c3d2one)

### 再次感谢该大佬的项目

### --------------------------

### historic version

## [CocosCreator3D 3.2.0](./version/Version%203.2.0.md)

## [CocosCreator3D 3.4.1](./version/Version%203.4.1.md)

## Production principle

### Because CocosCreator has been using the 3d version since[SystemJS](https://github.com/systemjs/systemjs)to load all script resources, but the way SystemJS loads scripts is to load scripts through web requests (maybe the method is wrong), but the prerequisite for Playable ad production is that it is not allowed to load resources through network requests, all related resources It must be in the same Html file, so the early version was loaded by modifying SystemJS, but after many tests (a long, long period of testing), I found some very ingenious methods, which can be done without modifying SystemJS.

#### [System.register](https://github.com/systemjs/systemjs/blob/main/docs/system-register.md)

After CocosCreator3D exports the mobile web project, many system-related scripts have the following characteristics. The first array of scripts registered by SystemJS contains dependencies, which represent those modules that depend on them, and the dependent modules will be loaded first.

```javascript
System.register(['dependency'], function (_export, _context) {
  var dep;
  return {
    setters: [function (_dep) {
      dep = _dep;
    }],
    execute: function () {
      _export({
        name: 'value'
      });
    }
  };
});
```

Then the first parameter of the System.register method can pass in a string, which I call an alias, usually the name of the js file plus the corresponding path, because there are many js files with the same file name, these will do some corresponding processing, Like the following, for example, the js file is PPS.js, which can be modified as follows:

```javascript
System.register('./PPS.js'，['./PPY.js'], function (_export, _context) {
  var dep;
  return {
    setters: [function (_dep) {
      dep = _dep;
    }],
    execute: function () {
      _export({
        name: 'value'
      });
    }
  };
});
```

The above means that when PPS.js is loaded, the PPY.js file will be loaded first, and then PPS.js will be loaded.
Below is the PPYjs file,

```javascript
System.register('./PPY.js'，[], function (_export, _context) {
  var dep;
  return {
    setters: [function (_dep) {
      dep = _dep;
    }],
    execute: function () {
      _export({
        name: 'value'
      });
    }
  };
});
```

Then inject PPY.js and PPS.js into the Html file, similar to the following

```html

<script>
  System.register('./PPS.js'，['./PPY.js'], function (_export, _context) {
  var dep;
  return {
    setters: [function (_dep) {
      dep = _dep;
    }],
    execute: function () {
      _export({
        name: 'value'
      });
    }
  };
});
</script>
<script>
   System.register('./PPY.js'，[], function (_export, _context) {
  var dep;
  return {
    setters: [function (_dep) {
      dep = _dep;
    }],
    execute: function () {
      _export({
        name: 'value'
      });
    }
  };
});
</script>
<script>
   System.import('./PPS.js').catch(function (err) {
        console.error(err);
    })
</script>
```

(If you want to run the example normally, the premise is to inject the SystemJS module, otherwise it will not run normally) The code at the bottom will be executed successfully

The above is a simple example. The project does a lot of corresponding processing through this method, and then generates a single Html file.

## WASM load file processing

Since the physics engine module may be required, of course, the WASM module may also be selected. In this case, you need to modify the engine source code (open "project directory/build/web-mobile/cocos-js/instantiated-\*\*.js" similar to this file)
Find key fields![](./pic/Snipaste_2022-03-06_19-31-09.jpg)

```javascript
WebAssembly.instantiate()
```

As shown in the figure, you will find that the fetch method is used in the script, which is not allowed by PlayableAD, so you can only change the corresponding fetch method and replace it as follows (the path address comes from[WASM path incoming processing](./README.md#wasm路径传入处理))

```javascript
let url = "cocos-js/" + (e)
let wa = base64DecToArr(window.res[url])
WebAssembly.instantiate(wa, o).then((function (e) {
    var n = e.instance.exports;
    Object.assign(qJ, n),
        t()
}
), n) 
```

This script gets the path of the incoming wasm resource, all text resources and binary resources exist in window.res

## WASM path incoming processing

Open "project directory/build/web-mobile/cocos-js/bullet.wasm-\*\*\*\*.js" similar to this file

```javascript
System.register([],(function(e,t){"use strict";return{execute:function(){
  e("default",new URL("assets/bullet.wasm-c98527b6.wasm",t.meta.url).href)}
  }}));

```

Replace the entire parameter of 'new URL()' with 'assets/bullet.wasm-c98527b6.wasm', which is very good. The name of the wasm file in the CocosCreator3D 3.4.1 version should not change during packaging (if it is modified, Please replace with the correct file name)
as follows:

```javascript
System.register([],(function(e,t){"use strict";return{execute:function(){
    e("default","assets/bullet.wasm-c98527b6.wasm")
}}}));

```

# what must be done

[wasm loading processing](./README.md#wasm加载文件处理)\|[WASM path incoming processing](./README.md#wasm路径传入处理)This operation is necessary, because the JavaScript fetch operation method is used, and only the engine source code can be forced to be modified (if wasm is not used, it can be ignored)
