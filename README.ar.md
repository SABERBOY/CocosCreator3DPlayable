# CocosCreator3D PlayableAD

### [تم التعديل بناءً على هذا المشروع](https://github.com/Jecced/c3d2one)

### شكرا مرة أخرى على مشروع الرجل الكبير

## وصف README

-   [إنجليزي](README.en.md)
-   [الصينية المبسطة](README.zh-CN.md)
-   [الصينية التقليدية](README.zh-TW.md)
-   [هندي](README.hi.md)
-   [فرنسي](README.fr.md)
-   [عربى](README.ar.md)

### نسخة تاريخية

## [CocosCreator3D 3.2.0.0 تحديث](./version/Version%203.2.0.md)

## [CocosCreator3D 3.4.1](./version/Version%203.4.1.md)

## مبدأ الإنتاج

### لأن CocosCreator يستخدم الإصدار ثلاثي الأبعاد منذ ذلك الحين[نظام JS](https://github.com/systemjs/systemjs)لتحميل جميع موارد البرنامج النصي ، ولكن الطريقة التي يقوم بها SystemJS بتحميل البرامج النصية هي تحميل البرامج النصية من خلال طلبات الويب (ربما تكون الطريقة خاطئة) ، ولكن الشرط الأساسي لإنتاج الإعلانات القابلة للتشغيل هو أنه لا يُسمح بتحميل الموارد من خلال طلبات الشبكة ، وجميع الموارد ذات الصلة يجب أن يكون في نفس ملف Html ، لذلك تم تحميل الإصدار المبكر عن طريق تعديل SystemJS ، ولكن بعد العديد من الاختبارات (فترة اختبار طويلة وطويلة) ، وجدت بعض الطرق البارعة للغاية ، والتي يمكن القيام بها دون تعديل SystemJS.

#### [نظام](https://github.com/systemjs/systemjs/blob/main/docs/system-register.md)

بعد أن يقوم CocosCreator3D بتصدير مشروع الويب المحمول ، فإن العديد من البرامج النصية ذات الصلة بالنظام لها الخصائص التالية: المصفوفة الأولى من البرامج النصية المسجلة بواسطة SystemJS تحتوي على التبعيات ، والتي تمثل تلك الوحدات النمطية التي تعتمد عليها ، وسيتم تحميل الوحدات التابعة أولاً.

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

بعد ذلك ، يمكن أن تمر المعلمة الأولى لطريقة System.register في سلسلة ، والتي أسميها اسمًا مستعارًا ، وعادةً ما يكون اسم ملف js بالإضافة إلى المسار المقابل ، نظرًا لوجود العديد من ملفات js بنفس اسم الملف ، فستقوم هذه ببعض المعالجة المقابلة ، مثل ما يلي ، على سبيل المثال ، ملف js هو PPS.js ، والذي يمكن تعديله على النحو التالي:

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

تحتاج إلى الانتباه إلى الحالات التالية (مهمة جدًا) حيث يجب أن يكون المكان المحدد هو نفسه كما هو الحال في الجلب ، لأن كل تجميع سيغير اسم المتغير المقابل![](./pic/Snipaste_2022-04-01_13-48-28.jpg)

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

استبدل المعلمة الكاملة لـ "عنوان URL الجديد ()" بـ "الأصول / bullet.wasm-c98527b6.wasm" ، وهو أمر جيد جدًا. يجب ألا يتغير اسم ملف wasm في إصدار CocosCreator3D 3.4.1 أثناء الحزم (إذا كان تم تعديله ، يرجى استبداله باسم الملف الصحيح)
كالآتي:

```javascript
System.register([],(function(e,t){"use strict";return{execute:function(){
    e("default","assets/bullet.wasm-c98527b6.wasm")
}}}));

```

# ما يجب القيام به

[عملية تحميل wasm](./README.md#wasm加载文件处理)\|[معالجة الوارد لمسار WASM](./README.md#wasm路径传入处理)هذه العملية ضرورية ، لأنه يتم استخدام طريقة عملية جلب JavaScript ، ويمكن إجبار شفرة مصدر المحرك فقط على التعديل (إذا لم يتم استخدام wasm ، فيمكن تجاهله)
