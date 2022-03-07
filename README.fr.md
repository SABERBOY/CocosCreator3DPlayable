# CocosCreator3D PlayableAD

### [Modifié sur la base de ce projet](https://github.com/Jecced/c3d2one)

### Merci encore pour le projet du grand gars

## Description LISEZMOI

-   [Anglais](README.en.md)
-   [Chinois simplifié](README.zh-CN.md)
-   [chinois traditionnel](README.zh-TW.md)
-   [hindi](README.hi.md)
-   [Française](README.fr.md)
-   [arabe](README.ar.md)

### version historique

## [CocosCreator3D 3.2.0](./version/Version%203.2.0.md)

## [CocosCreator3D 3.4.1](./version/Version%203.4.1.md)

## Principe de fabrication

### 因为CocosCreator 从3d版本后就开始通过用[SystemJS](https://github.com/systemjs/systemjs)pour charger toutes les ressources de script, mais la façon dont SystemJS charge les scripts consiste à charger des scripts via des requêtes Web (peut-être que la méthode est erronée), mais la condition préalable à la production d'annonces jouables est qu'il n'est pas autorisé à charger des ressources via des requêtes réseau, toutes les ressources associées Il doit être dans le même fichier Html, donc la première version a été chargée en modifiant SystemJS, mais après de nombreux tests (une longue, longue période de test), j'ai trouvé des méthodes très ingénieuses, qui peuvent être faites sans modifier SystemJS.

#### [System.register](https://github.com/systemjs/systemjs/blob/main/docs/system-register.md)

Une fois que CocosCreator3D a exporté un projet Web mobile, de nombreux scripts liés au système présentent les caractéristiques suivantes : le premier tableau de scripts enregistrés par SystemJS contient des dépendances, qui représentent les modules qui en dépendent, et les modules dépendants seront chargés en premier.

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

Ensuite, le premier paramètre de la méthode System.register peut passer une chaîne, que j'appelle un alias, généralement le nom du fichier js plus le chemin correspondant, car il existe de nombreux fichiers js avec le même nom de fichier, ceux-ci feront certains traitement correspondant, Comme le suivant, par exemple, le fichier js est PPS.js, qui peut être modifié comme suit :

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

Ce qui précède signifie que lorsque PPS.js est chargé, le fichier PPY.js sera chargé en premier, puis PPS.js sera chargé.
Ci-dessous le fichier PPYjs,

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

Injectez ensuite PPY.js et PPS.js dans le fichier Html, comme suit

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

(Si vous voulez exécuter l'exemple normalement, le principe est d'injecter le module SystemJS, sinon il ne fonctionnera pas normalement) Le code en bas sera exécuté avec succès

Ce qui précède est un exemple simple. Le projet effectue de nombreux traitements correspondants via cette méthode, puis génère un seul fichier Html.

## Traitement du fichier de chargement WASM

Etant donné que le module du moteur physique peut être requis, bien sûr, le module WASM peut également être sélectionné. Dans ce cas, vous devez modifier le code source du moteur (ouvrez "répertoire du projet/build/web-mobile/cocos-js/instantiated- \*\*.js" similaire à ce fichier)
Rechercher des champs clés![](./pic/Snipaste_2022-03-06_19-31-09.jpg)

```javascript
WebAssembly.instantiate()
```

Comme le montre la figure, vous constaterez que la méthode de récupération est utilisée dans le script, ce qui n'est pas autorisé par PlayableAD, vous ne pouvez donc modifier que la méthode de récupération correspondante et la remplacer comme suit (l'adresse du chemin provient de[Traitement entrant du chemin WASM](./README.md#wasm路径传入处理))

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

Ce script obtient le chemin de la ressource wasm entrante, toutes les ressources de texte et les ressources binaires existent dans window.res

## Traitement entrant du chemin WASM

Ouvrez "répertoire du projet/build/web-mobile/cocos-js/bullet.wasm-\*\*\*\*.js" similaire à ce fichier

```javascript
System.register([],(function(e,t){"use strict";return{execute:function(){
  e("default",new URL("assets/bullet.wasm-c98527b6.wasm",t.meta.url).href)}
  }}));

```

Remplacez tout le paramètre de 'new URL()' par 'assets/bullet.wasm-c98527b6.wasm', ce qui est très bien. Le nom du fichier wasm dans la version 3.4.1 de CocosCreator3D ne devrait pas changer lors de l'empaquetage (s'il est modifié, veuillez le remplacer par le nom de fichier correct)
comme suit:

```javascript
System.register([],(function(e,t){"use strict";return{execute:function(){
    e("default","assets/bullet.wasm-c98527b6.wasm")
}}}));

```

# ce qu'il faut faire

[traitement de chargement wasm](./README.md#wasm加载文件处理)\|[Traitement entrant du chemin WASM](./README.md#wasm路径传入处理)Cette opération est nécessaire, car la méthode d'opération de récupération JavaScript est utilisée et seul le code source du moteur peut être forcé à être modifié (si wasm n'est pas utilisé, il peut être ignoré)
