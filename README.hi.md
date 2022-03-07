# CocosCreator3D खेलने योग्यAD

### [इस परियोजना के आधार पर संशोधित](https://github.com/Jecced/c3d2one)

### बड़े आदमी की परियोजना के लिए फिर से धन्यवाद

### --------------------------

### ऐतिहासिक संस्करण

## [CocosCreator3D 3.2.0](./version/Version%203.2.0.md)

## [CocosCreator3D 3.4.1](./version/Version%203.4.1.md)

## उत्पादन सिद्धांत

### क्योंकि CocosCreator तब से 3D संस्करण का उपयोग कर रहा है[सिस्टमजेएस](https://github.com/systemjs/systemjs)सभी स्क्रिप्ट संसाधनों को लोड करने के लिए, लेकिन जिस तरह से सिस्टमजेएस स्क्रिप्ट लोड करता है वह वेब अनुरोधों के माध्यम से स्क्रिप्ट लोड करना है (शायद विधि गलत है), लेकिन प्लेएबल विज्ञापन उत्पादन के लिए शर्त यह है कि इसे नेटवर्क अनुरोधों के माध्यम से संसाधनों को लोड करने की अनुमति नहीं है, सभी संबंधित संसाधन यह एक ही एचटीएमएल फाइल में होना चाहिए, इसलिए सिस्टमजेएस को संशोधित करके प्रारंभिक संस्करण लोड किया गया था, लेकिन कई परीक्षणों (परीक्षण की लंबी, लंबी अवधि) के बाद, मुझे कुछ बहुत ही सरल विधियां मिलीं, जिन्हें सिस्टमजेएस को संशोधित किए बिना किया जा सकता है।

#### [सिस्टम.रजिस्टर](https://github.com/systemjs/systemjs/blob/main/docs/system-register.md)

CocosCreator3D द्वारा एक मोबाइल वेब प्रोजेक्ट निर्यात करने के बाद, कई सिस्टम-संबंधित स्क्रिप्ट में निम्नलिखित विशेषताएं होती हैं। SystemJS द्वारा पंजीकृत स्क्रिप्ट की पहली सरणी में निर्भरताएँ होती हैं, जो उन मॉड्यूल का प्रतिनिधित्व करती हैं जो उन पर निर्भर करते हैं, और आश्रित मॉड्यूल पहले लोड किए जाएंगे।

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

然后System.register方法的第一个参数可以传入string，这个我称之为别名，通常是js文件的名称加上对应路径，因为有很多相同文件名的js文件，这些会做一些对应处理，像下面这样，比如该js文件为PPS.js,可以修改如下：

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

उपरोक्त का अर्थ है कि जब PPS.js लोड किया जाता है, तो PPY.js फ़ाइल पहले लोड की जाएगी, और फिर PPS.js लोड की जाएगी।
नीचे PPYjs फ़ाइल है,

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

फिर PPY.js और PPS.js को Html फ़ाइल में इंजेक्ट करें, जैसा कि निम्न है

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

(यदि आप सामान्य रूप से उदाहरण चलाना चाहते हैं, तो आधार SystemJS मॉड्यूल को इंजेक्ट करना है, अन्यथा यह सामान्य रूप से नहीं चलेगा) नीचे दिए गए कोड को सफलतापूर्वक निष्पादित किया जाएगा

以上是简单的示例，该工程就是通过该方法做了大量对应处理，然后生成单个Html文件的

## WASM लोड फ़ाइल प्रसंस्करण

चूंकि भौतिकी इंजन मॉड्यूल की आवश्यकता हो सकती है, निश्चित रूप से, WASM मॉड्यूल भी चुना जा सकता है। इस मामले में, आपको इंजन स्रोत कोड को संशोधित करने की आवश्यकता है (खुले "प्रोजेक्ट निर्देशिका/बिल्ड/वेब-मोबाइल/कोकोस-जेएस/तत्काल- \*\*.js" इस फ़ाइल के समान)
प्रमुख फ़ील्ड खोजें![](./pic/Snipaste_2022-03-06_19-31-09.jpg)

```javascript
WebAssembly.instantiate()
```

जैसा कि चित्र में दिखाया गया है, आप पाएंगे कि स्क्रिप्ट में फ़ेच विधि का उपयोग किया जाता है, जिसे PlayableAD द्वारा अनुमति नहीं है, इसलिए आप केवल संबंधित फ़ेच विधि को बदल सकते हैं और इसे निम्नानुसार बदल सकते हैं (पथ का पता आता है[WASM पथ आने वाली प्रसंस्करण](./README.md#wasm路径传入处理))

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

इस स्क्रिप्ट को आने वाले wasm संसाधन का पथ मिलता है, सभी टेक्स्ट संसाधन और बाइनरी संसाधन window.res में मौजूद होते हैं

## WASM पथ आने वाली प्रसंस्करण

इस फ़ाइल के समान "प्रोजेक्ट डायरेक्टरी/बिल्ड/वेब-मोबाइल/कोकोस-जेएस/बुलेट.वासम-\*\*\*\*.जेएस" खोलें

```javascript
System.register([],(function(e,t){"use strict";return{execute:function(){
  e("default",new URL("assets/bullet.wasm-c98527b6.wasm",t.meta.url).href)}
  }}));

```

'नए URL ()' के पूरे पैरामीटर को 'assets/bullet.wasm-c98527b6.wasm' से बदलें, जो बहुत अच्छा है। CocosCreator3D 3.4.1 संस्करण में wasm फ़ाइल का नाम पैकेजिंग के दौरान नहीं बदलना चाहिए (यदि यह संशोधित किया गया है, कृपया सही फ़ाइल नाम से बदलें)
निम्नलिखित नुसार:

```javascript
System.register([],(function(e,t){"use strict";return{execute:function(){
    e("default","assets/bullet.wasm-c98527b6.wasm")
}}}));

```

# क्या किया जाना चाहिए

[वासम लोडिंग प्रोसेसिंग](./README.md#wasm加载文件处理)\|[WASM पथ आने वाली प्रसंस्करण](./README.md#wasm路径传入处理)
该操作是必须的，因为用到了JavaScript fetch操作方法，只能强制修改引擎源码（如果没有使用wasm可以忽略）
