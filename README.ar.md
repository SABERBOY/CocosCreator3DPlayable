# CocosCreator3D PlayableAD

### [تم التعديل بناءً على هذا المشروع](https://github.com/Jecced/c3d2one)

### 再次感谢该大佬的项目

## وصف README

-   [إنجليزي](README.en.md)
-   [الصينية المبسطة](README.zh-CN.md)
-   [الصينية التقليدية](README.zh-TW.md)
-   [هندي](README.hi.md)
-   [الفرنسية](README.fr.md)
-   [عربى](README.ar.md)

### نسخة تاريخية

## [CocosCreator3D 3.2.0.0 تحديث](./version/Version%203.2.0.md)

## [CocosCreator3D 3.4.1](./version/Version%203.4.1.md)

## مبدأ الإنتاج

### لأن CocosCreator يستخدم الإصدار ثلاثي الأبعاد منذ ذلك الحين[نظام JS](https://github.com/systemjs/systemjs)来加载所有的脚本资源，但是SystemJS加载脚本的方式是通过web请求方式（可能方法说得不对）来加载脚本，但是Playable广告制作的前提条件就是不允许通过网络请求来加载资源，所有有关的资源必须要在同一个Html文件里面，所以早期版本是通过魔改SystemJS来加载，但是经过多方面测试（很久很久的测试），发现一些很巧妙的方法，可以不魔改SystemJS就可行

#### [نظام](https://github.com/systemjs/systemjs/blob/main/docs/system-register.md)

بعد أن يقوم CocosCreator3D بتصدير مشروع ويب محمول ، فإن العديد من البرامج النصية ذات الصلة بالنظام لها الخصائص التالية: المصفوفة الأولى من البرامج النصية المسجلة بواسطة SystemJS تحتوي على التبعيات ، والتي تمثل تلك الوحدات النمطية التي تعتمد عليها ، وسيتم تحميل الوحدات التابعة أولاً.

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

ما ورد أعلاه يعني أنه عند تحميل PPS.js ، سيتم تحميل ملف PPY.js أولاً ، ثم سيتم تحميل PPS.js.
يوجد أدناه ملف PPYjs ،

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

ثم قم بحقن PPY.js و PPS.js في ملف Html ، على غرار ما يلي

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

(إذا كنت ترغب في تشغيل المثال بشكل طبيعي ، فإن الفرضية هي حقن وحدة SystemJS ، وإلا فلن تعمل بشكل طبيعي) سيتم تنفيذ الكود الموجود في الأسفل بنجاح

ما ورد أعلاه هو مثال بسيط ، حيث يقوم المشروع بالكثير من المعالجة المقابلة من خلال هذه الطريقة ، ثم يقوم بإنشاء ملف Html واحد.

## WASM معالجة ملف تحميل

نظرًا لأن وحدة محرك الفيزياء قد تكون مطلوبة ، بالطبع ، يمكن أيضًا تحديد وحدة WASM. في هذه الحالة ، تحتاج إلى تعديل شفرة مصدر المحرك (افتح "دليل المشروع / البناء / الويب المحمول / cocos-js / يتم إنشاء مثيل- \*\*. js "مشابه لهذا الملف)
ابحث عن الحقول الرئيسية![](./pic/Snipaste_2022-03-06_19-31-09.jpg)

```javascript
WebAssembly.instantiate()
```

كما هو موضح في الشكل ، ستجد أن طريقة الجلب مستخدمة في البرنامج النصي ، وهو ما لا يسمح به PlayableAD ، لذلك يمكنك فقط تغيير طريقة الجلب المقابلة واستبدالها على النحو التالي (يأتي عنوان المسار من[معالجة الوارد لمسار WASM](./README.md#wasm路径传入处理))

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

يحصل هذا البرنامج النصي على مسار مورد wasm الوارد ، وتوجد جميع موارد النص والموارد الثنائية في window.res

## معالجة الوارد لمسار WASM

افتح "project directory / build / web-mobile / cocos-js / bullet.wasm - \*\*\*\*. js" مشابه لهذا الملف

```javascript
System.register([],(function(e,t){"use strict";return{execute:function(){
  e("default",new URL("assets/bullet.wasm-c98527b6.wasm",t.meta.url).href)}
  }}));

```

استبدل المعلمة الكاملة لـ "عنوان URL الجديد ()" بـ "الأصول / bullet.wasm-c98527b6.wasm" ، وهو أمر جيد جدًا. يجب ألا يتغير اسم ملف wasm في إصدار CocosCreator3D 3.4.1 أثناء التعبئة (إذا كان تم تعديله ، يرجى استبدال اسم الملف الصحيح)
كالآتي:

```javascript
System.register([],(function(e,t){"use strict";return{execute:function(){
    e("default","assets/bullet.wasm-c98527b6.wasm")
}}}));

```

# ما يجب القيام به

[عملية تحميل wasm](./README.md#wasm加载文件处理)\|[معالجة الوارد لمسار WASM](./README.md#wasm路径传入处理)هذه العملية ضرورية ، لأنه يتم استخدام طريقة عملية جلب JavaScript ، ويمكن فقط إجبار رمز مصدر المحرك على التعديل (إذا لم يتم استخدام wasm ، فيمكن تجاهله)
