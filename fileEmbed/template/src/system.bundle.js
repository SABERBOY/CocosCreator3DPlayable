! function (t) {
    const e = "undefined" != typeof self,
        n = "undefined" != typeof document,
        r = e ? self : global;
    let o;
    if (n) {
        const t = document.querySelector("base[href]");
        t && (o = t.href)
    }
    if (!o && "undefined" != typeof location) {
        o = location.href.split("#")[0].split("?")[0];
        const t = o.lastIndexOf("/"); - 1 !== t && (o = o.slice(0, t + 1))
    }
    const i = /\\/g;

    function s(t, e) {
        if (-1 !== t.indexOf("\\") && (t = t.replace(i, "/")), "/" === t[0] && "/" === t[1]) return e.slice(0, e.indexOf(":") + 1) + t;
        if ("." === t[0] && ("/" === t[1] || "." === t[1] && ("/" === t[2] || 2 === t.length && (t += "/")) || 1 === t.length && (t += "/")) || "/" === t[0]) {
            const n = e.slice(0, e.indexOf(":") + 1);
            let r;
            if ("/" === e[n.length + 1] ? "file:" !== n ? (r = e.slice(n.length + 2), r = r.slice(r.indexOf("/") + 1)) : r = e.slice(8) : r = e.slice(n.length + ("/" === e[n.length])), "/" === t[0]) return e.slice(0, e.length - r.length - 1) + t;
            const o = r.slice(0, r.lastIndexOf("/") + 1) + t,
                i = [];
            let s = -1;
            for (let t = 0; t < o.length; t++) - 1 !== s ? "/" === o[t] && (i.push(o.slice(s, t + 1)), s = -1) : "." === o[t] ? "." !== o[t + 1] || "/" !== o[t + 2] && t + 2 !== o.length ? "/" === o[t + 1] || t + 1 === o.length ? t += 1 : s = t : (i.pop(), t += 2) : s = t;
            return -1 !== s && i.push(o.slice(s)), e.slice(0, e.length - r.length) + i.join("")
        }
    }

    function c(t, e) {
        return s(t, e) || (-1 !== t.indexOf(":") ? t : s("./" + t, e))
    }

    function l(t, e) {
        for (let n in e) t[n] = e[n];
        return t
    }

    function u(t, e, n, r, o) {
        for (let i in t) {
            const c = s(i, n) || i,
                l = t[i];
            if ("string" != typeof l) continue;
            const u = p(r, s(l, n) || l, o);
            u ? e[c] = u : h(i, l, "bare specifier did not resolve")
        }
    }

    function f(t, e) {
        if (e[t]) return t;
        let n = t.length;
        do {
            const r = t.slice(0, n + 1);
            if (r in e) return r
        } while (-1 !== (n = t.lastIndexOf("/", n - 1)))
    }

    function a(t, e) {
        const n = f(t, e);
        if (n) {
            const r = e[n];
            if (null === r) return;
            if (!(t.length > n.length && "/" !== r[r.length - 1])) return r + t.slice(n.length);
            h(n, r, "should have a trailing '/'")
        }
    }

    function h(t, e, n) {
        console.warn("Package target " + n + ", resolving target '" + e + "' for " + t)
    }

    function p(t, e, n) {
        let r = n && f(n, t.scopes);
        for (; r;) {
            const n = a(e, t.scopes[r]);
            if (n) return n;
            r = f(r.slice(0, r.lastIndexOf("/")), t.scopes)
        }
        return a(e, t.imports) || -1 !== e.indexOf(":") && e
    }
    const d = "undefined" != typeof Symbol,
        g = d && Symbol.toStringTag,
        m = d ? Symbol() : "@";

    function y() {
        this[m] = {}
    }
    const v = y.prototype;
    let E;
    v.prepareImport = function () {}, v.import = function (t, e) {
        const n = this;
        return Promise.resolve(n.prepareImport()).then((function () {
            return n.resolve(t, e)
        })).then((function (t) {
            const e = function t(e, n, r) {
                let o = e[m][n];
                if (o) return o;
                const i = [],
                    s = Object.create(null);
                g && Object.defineProperty(s, g, {
                    value: "Module"
                });
                let c = Promise.resolve().then((function () {
                    return e.instantiate(n, r)
                })).then((function (t) {
                    if (!t) throw Error("Module " + n + " did not instantiate");
                    const r = t[1]((function (t, e) {
                        o.h = !0;
                        let n = !1;
                        if ("object" != typeof t) t in s && s[t] === e || (s[t] = e, n = !0);
                        else {
                            for (let e in t) {
                                let r = t[e];
                                e in s && s[e] === r || (s[e] = r, n = !0)
                            }
                            t.__esModule && (s.__esModule = t.__esModule)
                        }
                        if (n)
                            for (let t = 0; t < i.length; t++) i[t](s);
                        return e
                    }), 2 === t[1].length ? {
                        import: function (t) {
                            return e.import(t, n)
                        },
                        meta: e.createContext(n)
                    } : void 0);
                    return o.e = r.execute || function () {}, [t[0], r.setters || []]
                }));
                const l = c.then((function (r) {
                    return Promise.all(r[0].map((function (o, i) {
                        const s = r[1][i];
                        return Promise.resolve(e.resolve(o, n)).then((function (r) {
                            const o = t(e, r, n);
                            return Promise.resolve(o.I).then((function () {
                                return s && (o.i.push(s), !o.h && o.I || s(o.n)), o
                            }))
                        }))
                    }))).then((function (t) {
                        o.d = t
                    }))
                }));
                return l.catch((function (t) {
                    o.e = null, o.er = t
                })), o = e[m][n] = {
                    id: n,
                    i: i,
                    n: s,
                    I: c,
                    L: l,
                    h: !1,
                    d: void 0,
                    e: void 0,
                    er: void 0,
                    E: void 0,
                    C: void 0
                }
            }(n, t);
            return e.C || function (t, e) {
                return e.C = function t(e, n, r) {
                    if (!r[n.id]) return r[n.id] = !0, Promise.resolve(n.L).then((function () {
                        return Promise.all(n.d.map((function (n) {
                            return t(e, n, r)
                        })))
                    }))
                }(t, e, {}).then((function () {
                    return function t(e, n, r) {
                        if (r[n.id]) return;
                        if (r[n.id] = !0, !n.e) {
                            if (n.er) throw n.er;
                            return n.E ? n.E : void 0
                        }
                        let o;
                        return n.d.forEach((function (n) {
                            {
                                const i = t(e, n, r);
                                i && (o = o || []).push(i)
                            }
                        })), o ? Promise.all(o).then(i) : i();

                        function i() {
                            try {
                                let t = n.e.call(O);
                                if (t) return t = t.then((function () {
                                    n.C = n.n, n.E = null
                                })), n.E = n.E || t;
                                n.C = n.n
                            } catch (t) {
                                throw n.er = t, t
                            } finally {
                                n.L = n.I = void 0, n.e = null
                            }
                        }
                    }(t, e, {})
                })).then((function () {
                    return e.n
                }))
            }(n, e)
        }))
    }, v.createContext = function (t) {
        return {
            url: t
        }
    }, v.register = function (t, e) {
        E = [t, e]
    }, v.getRegister = function () {
        const t = E;
        return E = void 0, t
    };
    const O = Object.freeze(Object.create(null));
    r.System = new y;
    const b = v.register;
    let S, w;

    function x() {
        Array.prototype.forEach.call(document.querySelectorAll("script[type=systemjs-module]"), (function (t) {
            t.src && System.import("import:" === t.src.slice(0, 7) ? t.src.slice(7) : c(t.src, o))
        }))
    }
    v.register = function (t, e) {
        b.call(this, t, e)
    }, v.createScript = function (t) {
        const e = document.createElement("script");
        return e.charset = "utf-8", e.async = !0, e.crossOrigin = "anonymous", e.src = t, e
    }, n && window.addEventListener("error", (function (t) {
        S = t.filename, w = t.error
    })), v.instantiate = function (t, e) {
        const n = this;
        return new Promise((function (r, o) {
            const i = v.createScript(t);
            i.addEventListener("error", (function () {
                o(Error("Error loading " + t + (e ? " from " + e : "")))
            })), i.addEventListener("load", (function () {
                document.head.removeChild(i), S === t ? o(w) : r(n.getRegister())
            })), document.head.appendChild(i)
        }))
    }, n && (window.addEventListener("DOMContentLoaded", x), x());
    let I = o;
    const P = {
        imports: {},
        scopes: {}
    };
    let R = {
        imports: {},
        scopes: {}
    };
    v.resolve = function (t, e) {
        return p(R, s(t, e = e || I) || t, e) || function (t, e) {
            throw Error("Unable to resolve specifier '" + t + (e ? "' from " + e : "'"))
        }(t, e)
    };
    const C = {
        setBaseUrl: function (t) {
            I = t
        },
        setImportMap: function (t, e) {
            R = function (t, e, n) {
                const r = {
                    imports: l({}, n.imports),
                    scopes: l({}, n.scopes)
                };
                if (t.imports && u(t.imports, r.imports, e, n, null), t.scopes)
                    for (let o in t.scopes) {
                        const i = c(o, e);
                        u(t.scopes[o], r.scopes[i] || (r.scopes[i] = {}), e, n, i)
                    }
                return r
            }(t, e || I, P)
        },
        hookInstantiationOverSchema: function (t, e) {
            const n = v.instantiate;
            v.instantiate = function (r, o) {
                const i = r.substr(0, t.length) === t ? r.substr(t.length) : null;
                return null === i ? n.call(this, r, o) : e.call(this, i, o)
            }
        }
    };
    v.patches = C,
        function (t) {
            const e = t.System;
            s(e);
            const n = e.constructor.prototype,
                r = e.constructor,
                o = function () {
                    r.call(this), s(this)
                };
            let i;

            function s(t) {
                t.registerRegistry = Object.create(null)
            }
            o.prototype = n, e.constructor = o;
            const c = n.register;
            n.register = function (t, e, n) {
                if ("string" != typeof t) return c.apply(this, arguments);
                const r = [e, n];
                t = new URL(t, window.location.href).href;
                return this.registerRegistry[t] = r, i || (i = r, Promise.resolve().then((function () {
                    i = null
                }))), c.apply(this, arguments)
            };
            const l = n.resolve;
            n.resolve = function (t, e) {
                try {
                    return l.call(this, t, e)
                } catch (e) {
                    if (t in this.registerRegistry) return t;
                    throw e
                }
            };
            const u = n.instantiate;
            n.instantiate = function (t, e) {
                const n = this.registerRegistry[t];
                return n ? (this.registerRegistry[t] = null, n) : u.call(this, t, e)
            };
            const f = n.getRegister;
            n.getRegister = function () {
                const t = f.call(this),
                    e = i || t;
                return i = null, e
            }
        }("undefined" != typeof self ? self : global), t.REGISTRY = m, t.patches = C, t.systemJSPrototype = v
}({});