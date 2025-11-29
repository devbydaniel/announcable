(function(){"use strict";try{if(typeof document<"u"){var n=document.createElement("style");n.appendChild(document.createTextNode(":root{--announcable-background: #ffffff;--announcable-foreground: #0a0a0b;--announcable-card: #ffffff;--announcable-card-foreground: #0a0a0b;--announcable-border: #e5e5e5;--announcable-radius: .5rem}.announcable-widget{box-sizing:border-box}.announcable-widget *,.announcable-widget *:before,.announcable-widget *:after{box-sizing:inherit}")),document.head.appendChild(n)}}catch(e){console.error("vite-plugin-css-injected-by-js",e)}})();
(function(w){typeof define=="function"&&define.amd?define(w):w()})(function(){"use strict";/**
 * @license
 * Copyright 2019 Google LLC
 * SPDX-License-Identifier: BSD-3-Clause
 */var pt,ft,gt;const w=globalThis,be=w.ShadowRoot&&(w.ShadyCSS===void 0||w.ShadyCSS.nativeShadow)&&"adoptedStyleSheets"in Document.prototype&&"replace"in CSSStyleSheet.prototype,me=Symbol(),He=new WeakMap;let We=class{constructor(e,t,r){if(this._$cssResult$=!0,r!==me)throw Error("CSSResult is not constructable. Use `unsafeCSS` or `css` instead.");this.cssText=e,this.t=t}get styleSheet(){let e=this.o;const t=this.t;if(be&&e===void 0){const r=t!==void 0&&t.length===1;r&&(e=He.get(t)),e===void 0&&((this.o=e=new CSSStyleSheet).replaceSync(this.cssText),r&&He.set(t,e))}return e}toString(){return this.cssText}};const bt=i=>new We(typeof i=="string"?i:i+"",void 0,me),g=(i,...e)=>{const t=i.length===1?i[0]:e.reduce((r,s,o)=>r+(n=>{if(n._$cssResult$===!0)return n.cssText;if(typeof n=="number")return n;throw Error("Value passed to 'css' function must be a 'css' function result: "+n+". Use 'unsafeCSS' to pass non-literal values, but take care to ensure page security.")})(s)+i[o+1],i[0]);return new We(t,i,me)},mt=(i,e)=>{if(be)i.adoptedStyleSheets=e.map(t=>t instanceof CSSStyleSheet?t:t.styleSheet);else for(const t of e){const r=document.createElement("style"),s=w.litNonce;s!==void 0&&r.setAttribute("nonce",s),r.textContent=t.cssText,i.appendChild(r)}},Be=be?i=>i:i=>i instanceof CSSStyleSheet?(e=>{let t="";for(const r of e.cssRules)t+=r.cssText;return bt(t)})(i):i;/**
 * @license
 * Copyright 2017 Google LLC
 * SPDX-License-Identifier: BSD-3-Clause
 */const{is:$t,defineProperty:_t,getOwnPropertyDescriptor:vt,getOwnPropertyNames:yt,getOwnPropertySymbols:wt,getPrototypeOf:At}=Object,_=globalThis,Ve=_.trustedTypes,xt=Ve?Ve.emptyScript:"",$e=_.reactiveElementPolyfillSupport,z=(i,e)=>i,ee={toAttribute(i,e){switch(e){case Boolean:i=i?xt:null;break;case Object:case Array:i=i==null?i:JSON.stringify(i)}return i},fromAttribute(i,e){let t=i;switch(e){case Boolean:t=i!==null;break;case Number:t=i===null?null:Number(i);break;case Object:case Array:try{t=JSON.parse(i)}catch{t=null}}return t}},te=(i,e)=>!$t(i,e),Ge={attribute:!0,type:String,converter:ee,reflect:!1,useDefault:!1,hasChanged:te};Symbol.metadata??(Symbol.metadata=Symbol("metadata")),_.litPropertyMetadata??(_.litPropertyMetadata=new WeakMap);let k=class extends HTMLElement{static addInitializer(e){this._$Ei(),(this.l??(this.l=[])).push(e)}static get observedAttributes(){return this.finalize(),this._$Eh&&[...this._$Eh.keys()]}static createProperty(e,t=Ge){if(t.state&&(t.attribute=!1),this._$Ei(),this.prototype.hasOwnProperty(e)&&((t=Object.create(t)).wrapped=!0),this.elementProperties.set(e,t),!t.noAccessor){const r=Symbol(),s=this.getPropertyDescriptor(e,r,t);s!==void 0&&_t(this.prototype,e,s)}}static getPropertyDescriptor(e,t,r){const{get:s,set:o}=vt(this.prototype,e)??{get(){return this[t]},set(n){this[t]=n}};return{get:s,set(n){const c=s==null?void 0:s.call(this);o==null||o.call(this,n),this.requestUpdate(e,c,r)},configurable:!0,enumerable:!0}}static getPropertyOptions(e){return this.elementProperties.get(e)??Ge}static _$Ei(){if(this.hasOwnProperty(z("elementProperties")))return;const e=At(this);e.finalize(),e.l!==void 0&&(this.l=[...e.l]),this.elementProperties=new Map(e.elementProperties)}static finalize(){if(this.hasOwnProperty(z("finalized")))return;if(this.finalized=!0,this._$Ei(),this.hasOwnProperty(z("properties"))){const t=this.properties,r=[...yt(t),...wt(t)];for(const s of r)this.createProperty(s,t[s])}const e=this[Symbol.metadata];if(e!==null){const t=litPropertyMetadata.get(e);if(t!==void 0)for(const[r,s]of t)this.elementProperties.set(r,s)}this._$Eh=new Map;for(const[t,r]of this.elementProperties){const s=this._$Eu(t,r);s!==void 0&&this._$Eh.set(s,t)}this.elementStyles=this.finalizeStyles(this.styles)}static finalizeStyles(e){const t=[];if(Array.isArray(e)){const r=new Set(e.flat(1/0).reverse());for(const s of r)t.unshift(Be(s))}else e!==void 0&&t.push(Be(e));return t}static _$Eu(e,t){const r=t.attribute;return r===!1?void 0:typeof r=="string"?r:typeof e=="string"?e.toLowerCase():void 0}constructor(){super(),this._$Ep=void 0,this.isUpdatePending=!1,this.hasUpdated=!1,this._$Em=null,this._$Ev()}_$Ev(){var e;this._$ES=new Promise(t=>this.enableUpdating=t),this._$AL=new Map,this._$E_(),this.requestUpdate(),(e=this.constructor.l)==null||e.forEach(t=>t(this))}addController(e){var t;(this._$EO??(this._$EO=new Set)).add(e),this.renderRoot!==void 0&&this.isConnected&&((t=e.hostConnected)==null||t.call(e))}removeController(e){var t;(t=this._$EO)==null||t.delete(e)}_$E_(){const e=new Map,t=this.constructor.elementProperties;for(const r of t.keys())this.hasOwnProperty(r)&&(e.set(r,this[r]),delete this[r]);e.size>0&&(this._$Ep=e)}createRenderRoot(){const e=this.shadowRoot??this.attachShadow(this.constructor.shadowRootOptions);return mt(e,this.constructor.elementStyles),e}connectedCallback(){var e;this.renderRoot??(this.renderRoot=this.createRenderRoot()),this.enableUpdating(!0),(e=this._$EO)==null||e.forEach(t=>{var r;return(r=t.hostConnected)==null?void 0:r.call(t)})}enableUpdating(e){}disconnectedCallback(){var e;(e=this._$EO)==null||e.forEach(t=>{var r;return(r=t.hostDisconnected)==null?void 0:r.call(t)})}attributeChangedCallback(e,t,r){this._$AK(e,r)}_$ET(e,t){var o;const r=this.constructor.elementProperties.get(e),s=this.constructor._$Eu(e,r);if(s!==void 0&&r.reflect===!0){const n=(((o=r.converter)==null?void 0:o.toAttribute)!==void 0?r.converter:ee).toAttribute(t,r.type);this._$Em=e,n==null?this.removeAttribute(s):this.setAttribute(s,n),this._$Em=null}}_$AK(e,t){var o,n;const r=this.constructor,s=r._$Eh.get(e);if(s!==void 0&&this._$Em!==s){const c=r.getPropertyOptions(s),a=typeof c.converter=="function"?{fromAttribute:c.converter}:((o=c.converter)==null?void 0:o.fromAttribute)!==void 0?c.converter:ee;this._$Em=s;const f=a.fromAttribute(t,c.type);this[s]=f??((n=this._$Ej)==null?void 0:n.get(s))??f,this._$Em=null}}requestUpdate(e,t,r){var s;if(e!==void 0){const o=this.constructor,n=this[e];if(r??(r=o.getPropertyOptions(e)),!((r.hasChanged??te)(n,t)||r.useDefault&&r.reflect&&n===((s=this._$Ej)==null?void 0:s.get(e))&&!this.hasAttribute(o._$Eu(e,r))))return;this.C(e,t,r)}this.isUpdatePending===!1&&(this._$ES=this._$EP())}C(e,t,{useDefault:r,reflect:s,wrapped:o},n){r&&!(this._$Ej??(this._$Ej=new Map)).has(e)&&(this._$Ej.set(e,n??t??this[e]),o!==!0||n!==void 0)||(this._$AL.has(e)||(this.hasUpdated||r||(t=void 0),this._$AL.set(e,t)),s===!0&&this._$Em!==e&&(this._$Eq??(this._$Eq=new Set)).add(e))}async _$EP(){this.isUpdatePending=!0;try{await this._$ES}catch(t){Promise.reject(t)}const e=this.scheduleUpdate();return e!=null&&await e,!this.isUpdatePending}scheduleUpdate(){return this.performUpdate()}performUpdate(){var r;if(!this.isUpdatePending)return;if(!this.hasUpdated){if(this.renderRoot??(this.renderRoot=this.createRenderRoot()),this._$Ep){for(const[o,n]of this._$Ep)this[o]=n;this._$Ep=void 0}const s=this.constructor.elementProperties;if(s.size>0)for(const[o,n]of s){const{wrapped:c}=n,a=this[o];c!==!0||this._$AL.has(o)||a===void 0||this.C(o,void 0,n,a)}}let e=!1;const t=this._$AL;try{e=this.shouldUpdate(t),e?(this.willUpdate(t),(r=this._$EO)==null||r.forEach(s=>{var o;return(o=s.hostUpdate)==null?void 0:o.call(s)}),this.update(t)):this._$EM()}catch(s){throw e=!1,this._$EM(),s}e&&this._$AE(t)}willUpdate(e){}_$AE(e){var t;(t=this._$EO)==null||t.forEach(r=>{var s;return(s=r.hostUpdated)==null?void 0:s.call(r)}),this.hasUpdated||(this.hasUpdated=!0,this.firstUpdated(e)),this.updated(e)}_$EM(){this._$AL=new Map,this.isUpdatePending=!1}get updateComplete(){return this.getUpdateComplete()}getUpdateComplete(){return this._$ES}shouldUpdate(e){return!0}update(e){this._$Eq&&(this._$Eq=this._$Eq.forEach(t=>this._$ET(t,this[t]))),this._$EM()}updated(e){}firstUpdated(e){}};k.elementStyles=[],k.shadowRootOptions={mode:"open"},k[z("elementProperties")]=new Map,k[z("finalized")]=new Map,$e==null||$e({ReactiveElement:k}),(_.reactiveElementVersions??(_.reactiveElementVersions=[])).push("2.1.1");/**
 * @license
 * Copyright 2017 Google LLC
 * SPDX-License-Identifier: BSD-3-Clause
 */const R=globalThis,se=R.trustedTypes,Fe=se?se.createPolicy("lit-html",{createHTML:i=>i}):void 0,Ye="$lit$",v=`lit$${Math.random().toFixed(9).slice(2)}$`,Ke="?"+v,Ct=`<${Ke}>`,A=document,U=()=>A.createComment(""),L=i=>i===null||typeof i!="object"&&typeof i!="function",_e=Array.isArray,Ot=i=>_e(i)||typeof(i==null?void 0:i[Symbol.iterator])=="function",ve=`[ 	
\f\r]`,q=/<(?:(!--|\/[^a-zA-Z])|(\/?[a-zA-Z][^>\s]*)|(\/?$))/g,Xe=/-->/g,Je=/>/g,x=RegExp(`>|${ve}(?:([^\\s"'>=/]+)(${ve}*=${ve}*(?:[^ 	
\f\r"'\`<>=]|("|')|))|$)`,"g"),Ze=/'/g,Qe=/"/g,et=/^(?:script|style|textarea|title)$/i,tt=i=>(e,...t)=>({_$litType$:i,strings:e,values:t}),l=tt(1),ie=tt(2),E=Symbol.for("lit-noChange"),b=Symbol.for("lit-nothing"),st=new WeakMap,C=A.createTreeWalker(A,129);function it(i,e){if(!_e(i)||!i.hasOwnProperty("raw"))throw Error("invalid template strings array");return Fe!==void 0?Fe.createHTML(e):e}const kt=(i,e)=>{const t=i.length-1,r=[];let s,o=e===2?"<svg>":e===3?"<math>":"",n=q;for(let c=0;c<t;c++){const a=i[c];let f,m,p=-1,$=0;for(;$<a.length&&(n.lastIndex=$,m=n.exec(a),m!==null);)$=n.lastIndex,n===q?m[1]==="!--"?n=Xe:m[1]!==void 0?n=Je:m[2]!==void 0?(et.test(m[2])&&(s=RegExp("</"+m[2],"g")),n=x):m[3]!==void 0&&(n=x):n===x?m[0]===">"?(n=s??q,p=-1):m[1]===void 0?p=-2:(p=n.lastIndex-m[2].length,f=m[1],n=m[3]===void 0?x:m[3]==='"'?Qe:Ze):n===Qe||n===Ze?n=x:n===Xe||n===Je?n=q:(n=x,s=void 0);const y=n===x&&i[c+1].startsWith("/>")?" ":"";o+=n===q?a+Ct:p>=0?(r.push(f),a.slice(0,p)+Ye+a.slice(p)+v+y):a+v+(p===-2?c:y)}return[it(i,o+(i[t]||"<?>")+(e===2?"</svg>":e===3?"</math>":"")),r]};class H{constructor({strings:e,_$litType$:t},r){let s;this.parts=[];let o=0,n=0;const c=e.length-1,a=this.parts,[f,m]=kt(e,t);if(this.el=H.createElement(f,r),C.currentNode=this.el.content,t===2||t===3){const p=this.el.content.firstChild;p.replaceWith(...p.childNodes)}for(;(s=C.nextNode())!==null&&a.length<c;){if(s.nodeType===1){if(s.hasAttributes())for(const p of s.getAttributeNames())if(p.endsWith(Ye)){const $=m[n++],y=s.getAttribute(p).split(v),ge=/([.?@])?(.*)/.exec($);a.push({type:1,index:o,name:ge[2],strings:y,ctor:ge[1]==="."?St:ge[1]==="?"?Pt:ge[1]==="@"?Tt:re}),s.removeAttribute(p)}else p.startsWith(v)&&(a.push({type:6,index:o}),s.removeAttribute(p));if(et.test(s.tagName)){const p=s.textContent.split(v),$=p.length-1;if($>0){s.textContent=se?se.emptyScript:"";for(let y=0;y<$;y++)s.append(p[y],U()),C.nextNode(),a.push({type:2,index:++o});s.append(p[$],U())}}}else if(s.nodeType===8)if(s.data===Ke)a.push({type:2,index:o});else{let p=-1;for(;(p=s.data.indexOf(v,p+1))!==-1;)a.push({type:7,index:o}),p+=v.length-1}o++}}static createElement(e,t){const r=A.createElement("template");return r.innerHTML=e,r}}function S(i,e,t=i,r){var n,c;if(e===E)return e;let s=r!==void 0?(n=t._$Co)==null?void 0:n[r]:t._$Cl;const o=L(e)?void 0:e._$litDirective$;return(s==null?void 0:s.constructor)!==o&&((c=s==null?void 0:s._$AO)==null||c.call(s,!1),o===void 0?s=void 0:(s=new o(i),s._$AT(i,t,r)),r!==void 0?(t._$Co??(t._$Co=[]))[r]=s:t._$Cl=s),s!==void 0&&(e=S(i,s._$AS(i,e.values),s,r)),e}class Et{constructor(e,t){this._$AV=[],this._$AN=void 0,this._$AD=e,this._$AM=t}get parentNode(){return this._$AM.parentNode}get _$AU(){return this._$AM._$AU}u(e){const{el:{content:t},parts:r}=this._$AD,s=((e==null?void 0:e.creationScope)??A).importNode(t,!0);C.currentNode=s;let o=C.nextNode(),n=0,c=0,a=r[0];for(;a!==void 0;){if(n===a.index){let f;a.type===2?f=new W(o,o.nextSibling,this,e):a.type===1?f=new a.ctor(o,a.name,a.strings,this,e):a.type===6&&(f=new It(o,this,e)),this._$AV.push(f),a=r[++c]}n!==(a==null?void 0:a.index)&&(o=C.nextNode(),n++)}return C.currentNode=A,s}p(e){let t=0;for(const r of this._$AV)r!==void 0&&(r.strings!==void 0?(r._$AI(e,r,t),t+=r.strings.length-2):r._$AI(e[t])),t++}}class W{get _$AU(){var e;return((e=this._$AM)==null?void 0:e._$AU)??this._$Cv}constructor(e,t,r,s){this.type=2,this._$AH=b,this._$AN=void 0,this._$AA=e,this._$AB=t,this._$AM=r,this.options=s,this._$Cv=(s==null?void 0:s.isConnected)??!0}get parentNode(){let e=this._$AA.parentNode;const t=this._$AM;return t!==void 0&&(e==null?void 0:e.nodeType)===11&&(e=t.parentNode),e}get startNode(){return this._$AA}get endNode(){return this._$AB}_$AI(e,t=this){e=S(this,e,t),L(e)?e===b||e==null||e===""?(this._$AH!==b&&this._$AR(),this._$AH=b):e!==this._$AH&&e!==E&&this._(e):e._$litType$!==void 0?this.$(e):e.nodeType!==void 0?this.T(e):Ot(e)?this.k(e):this._(e)}O(e){return this._$AA.parentNode.insertBefore(e,this._$AB)}T(e){this._$AH!==e&&(this._$AR(),this._$AH=this.O(e))}_(e){this._$AH!==b&&L(this._$AH)?this._$AA.nextSibling.data=e:this.T(A.createTextNode(e)),this._$AH=e}$(e){var o;const{values:t,_$litType$:r}=e,s=typeof r=="number"?this._$AC(e):(r.el===void 0&&(r.el=H.createElement(it(r.h,r.h[0]),this.options)),r);if(((o=this._$AH)==null?void 0:o._$AD)===s)this._$AH.p(t);else{const n=new Et(s,this),c=n.u(this.options);n.p(t),this.T(c),this._$AH=n}}_$AC(e){let t=st.get(e.strings);return t===void 0&&st.set(e.strings,t=new H(e)),t}k(e){_e(this._$AH)||(this._$AH=[],this._$AR());const t=this._$AH;let r,s=0;for(const o of e)s===t.length?t.push(r=new W(this.O(U()),this.O(U()),this,this.options)):r=t[s],r._$AI(o),s++;s<t.length&&(this._$AR(r&&r._$AB.nextSibling,s),t.length=s)}_$AR(e=this._$AA.nextSibling,t){var r;for((r=this._$AP)==null?void 0:r.call(this,!1,!0,t);e!==this._$AB;){const s=e.nextSibling;e.remove(),e=s}}setConnected(e){var t;this._$AM===void 0&&(this._$Cv=e,(t=this._$AP)==null||t.call(this,e))}}class re{get tagName(){return this.element.tagName}get _$AU(){return this._$AM._$AU}constructor(e,t,r,s,o){this.type=1,this._$AH=b,this._$AN=void 0,this.element=e,this.name=t,this._$AM=s,this.options=o,r.length>2||r[0]!==""||r[1]!==""?(this._$AH=Array(r.length-1).fill(new String),this.strings=r):this._$AH=b}_$AI(e,t=this,r,s){const o=this.strings;let n=!1;if(o===void 0)e=S(this,e,t,0),n=!L(e)||e!==this._$AH&&e!==E,n&&(this._$AH=e);else{const c=e;let a,f;for(e=o[0],a=0;a<o.length-1;a++)f=S(this,c[r+a],t,a),f===E&&(f=this._$AH[a]),n||(n=!L(f)||f!==this._$AH[a]),f===b?e=b:e!==b&&(e+=(f??"")+o[a+1]),this._$AH[a]=f}n&&!s&&this.j(e)}j(e){e===b?this.element.removeAttribute(this.name):this.element.setAttribute(this.name,e??"")}}class St extends re{constructor(){super(...arguments),this.type=3}j(e){this.element[this.name]=e===b?void 0:e}}class Pt extends re{constructor(){super(...arguments),this.type=4}j(e){this.element.toggleAttribute(this.name,!!e&&e!==b)}}class Tt extends re{constructor(e,t,r,s,o){super(e,t,r,s,o),this.type=5}_$AI(e,t=this){if((e=S(this,e,t,0)??b)===E)return;const r=this._$AH,s=e===b&&r!==b||e.capture!==r.capture||e.once!==r.once||e.passive!==r.passive,o=e!==b&&(r===b||s);s&&this.element.removeEventListener(this.name,this,r),o&&this.element.addEventListener(this.name,this,e),this._$AH=e}handleEvent(e){var t;typeof this._$AH=="function"?this._$AH.call(((t=this.options)==null?void 0:t.host)??this.element,e):this._$AH.handleEvent(e)}}class It{constructor(e,t,r){this.element=e,this.type=6,this._$AN=void 0,this._$AM=t,this.options=r}get _$AU(){return this._$AM._$AU}_$AI(e){S(this,e)}}const ye=R.litHtmlPolyfillSupport;ye==null||ye(H,W),(R.litHtmlVersions??(R.litHtmlVersions=[])).push("3.3.1");const rt=(i,e,t)=>{const r=(t==null?void 0:t.renderBefore)??e;let s=r._$litPart$;if(s===void 0){const o=(t==null?void 0:t.renderBefore)??null;r._$litPart$=s=new W(e.insertBefore(U(),o),o,void 0,t??{})}return s._$AI(i),s};/**
 * @license
 * Copyright 2017 Google LLC
 * SPDX-License-Identifier: BSD-3-Clause
 */const O=globalThis;let d=class extends k{constructor(){super(...arguments),this.renderOptions={host:this},this._$Do=void 0}createRenderRoot(){var t;const e=super.createRenderRoot();return(t=this.renderOptions).renderBefore??(t.renderBefore=e.firstChild),e}update(e){const t=this.render();this.hasUpdated||(this.renderOptions.isConnected=this.isConnected),super.update(e),this._$Do=rt(t,this.renderRoot,this.renderOptions)}connectedCallback(){var e;super.connectedCallback(),(e=this._$Do)==null||e.setConnected(!0)}disconnectedCallback(){var e;super.disconnectedCallback(),(e=this._$Do)==null||e.setConnected(!1)}render(){return E}};d._$litElement$=!0,d.finalized=!0,(pt=O.litElementHydrateSupport)==null||pt.call(O,{LitElement:d});const we=O.litElementPolyfillSupport;we==null||we({LitElement:d}),(O.litElementVersions??(O.litElementVersions=[])).push("4.2.1");/**
 * @license
 * Copyright 2017 Google LLC
 * SPDX-License-Identifier: BSD-3-Clause
 */const u=i=>(e,t)=>{t!==void 0?t.addInitializer(()=>{customElements.define(i,e)}):customElements.define(i,e)};/**
 * @license
 * Copyright 2017 Google LLC
 * SPDX-License-Identifier: BSD-3-Clause
 */const Nt={attribute:!0,type:String,converter:ee,reflect:!1,hasChanged:te},jt=(i=Nt,e,t)=>{const{kind:r,metadata:s}=t;let o=globalThis.litPropertyMetadata.get(s);if(o===void 0&&globalThis.litPropertyMetadata.set(s,o=new Map),r==="setter"&&((i=Object.create(i)).wrapped=!0),o.set(t.name,i),r==="accessor"){const{name:n}=t;return{set(c){const a=e.get.call(this);e.set.call(this,c),this.requestUpdate(n,a,i)},init(c){return c!==void 0&&this.C(n,void 0,i,c),c}}}if(r==="setter"){const{name:n}=t;return function(c){const a=this[n];e.call(this,c),this.requestUpdate(n,a,i)}}throw Error("Unsupported decorator location: "+r)};function h(i){return(e,t)=>typeof t=="object"?jt(i,e,t):((r,s,o)=>{const n=s.hasOwnProperty(o);return s.constructor.createProperty(o,r),n?Object.getOwnPropertyDescriptor(s,o):void 0})(i,e,t)}/**
 * @license
 * Copyright 2017 Google LLC
 * SPDX-License-Identifier: BSD-3-Clause
 */const oe={PENDING:1,COMPLETE:2},Mt=Symbol();let ne=class{get taskComplete(){return this.t||(this.i===1?this.t=new Promise((e,t)=>{this.o=e,this.h=t}):this.i===3?this.t=Promise.reject(this.l):this.t=Promise.resolve(this.u)),this.t}constructor(e,t,r){var o;this.p=0,this.i=0,(this._=e).addController(this);const s=typeof t=="object"?t:{task:t,args:r};this.v=s.task,this.j=s.args,this.m=s.argsEqual??Dt,this.k=s.onComplete,this.A=s.onError,this.autoRun=s.autoRun??!0,"initialValue"in s&&(this.u=s.initialValue,this.i=2,this.O=(o=this.T)==null?void 0:o.call(this))}hostUpdate(){this.autoRun===!0&&this.S()}hostUpdated(){this.autoRun==="afterUpdate"&&this.S()}T(){if(this.j===void 0)return;const e=this.j();if(!Array.isArray(e))throw Error("The args function must return an array");return e}async S(){const e=this.T(),t=this.O;this.O=e,e===t||e===void 0||t!==void 0&&this.m(t,e)||await this.run(e)}async run(e){var n,c,a,f,m;let t,r;e??(e=this.T()),this.O=e,this.i===1?(n=this.q)==null||n.abort():(this.t=void 0,this.o=void 0,this.h=void 0),this.i=1,this.autoRun==="afterUpdate"?queueMicrotask(()=>this._.requestUpdate()):this._.requestUpdate();const s=++this.p;this.q=new AbortController;let o=!1;try{t=await this.v(e,{signal:this.q.signal})}catch(p){o=!0,r=p}if(this.p===s){if(t===Mt)this.i=0;else{if(o===!1){try{(c=this.k)==null||c.call(this,t)}catch{}this.i=2,(a=this.o)==null||a.call(this,t)}else{try{(f=this.A)==null||f.call(this,r)}catch{}this.i=3,(m=this.h)==null||m.call(this,r)}this.u=t,this.l=r}this._.requestUpdate()}}abort(e){var t;this.i===1&&((t=this.q)==null||t.abort(e))}get value(){return this.u}get error(){return this.l}get status(){return this.i}render(e){var t,r,s,o;switch(this.i){case 0:return(t=e.initial)==null?void 0:t.call(e);case 1:return(r=e.pending)==null?void 0:r.call(e);case 2:return(s=e.complete)==null?void 0:s.call(e,this.value);case 3:return(o=e.error)==null?void 0:o.call(e,this.error);default:throw Error("Unexpected status: "+this.i)}}};const Dt=(i,e)=>i===e||i.length===e.length&&i.every((t,r)=>!te(t,e[r])),ot="announcable_last_opened",nt="announcable_client_id";function at(){try{const i="__storage_test__";return localStorage.setItem(i,i),localStorage.removeItem(i),!0}catch{return!1}}function lt(i){if(!at())return console.warn("[Announcable] localStorage not available"),null;try{return localStorage.getItem(i)}catch(e){return console.error("[Announcable] Error reading from localStorage:",e),null}}function ct(i,e){if(!at())return console.warn("[Announcable] localStorage not available"),!1;try{return localStorage.setItem(i,e),!0}catch(t){return t instanceof Error&&t.name==="QuotaExceededError"?console.error("[Announcable] localStorage quota exceeded"):console.error("[Announcable] Error writing to localStorage:",t),!1}}function zt(){const i=lt(ot);return i&&!isNaN(parseInt(i))?i:null}function Rt(i){return!i||isNaN(parseInt(i))?(console.error("[Announcable] Invalid timestamp format:",i),!1):ct(ot,i)}function Ut(){return Date.now().toString()}function Lt(){const i=lt(nt);if(i)return i;const e="xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx".replace(/[xy]/g,t=>{const r=Math.random()*16|0;return(t==="x"?r:r&3|8).toString(16)});return ct(nt,e),e}class qt{constructor(e,t){this.isOpen=!1,this.lastOpened=null,this.anchors=null,this.boundToggleWidget=this.toggleWidget.bind(this),this.observerSetupAttempted=!1,this.mutationObserver=null,this.host=e,this.querySelector=t,e.addController(this),this.lastOpened=zt()}hostConnected(){this.setupAnchorListeners(),this.querySelector&&typeof MutationObserver<"u"&&this.setupMutationObserver()}hostDisconnected(){this.cleanup()}setupAnchorListeners(){if(this.querySelector)try{const e=document.querySelectorAll(this.querySelector);if(!e||e.length===0){this.observerSetupAttempted||(console.warn(`[Announcable] No elements found for selector: ${this.querySelector}`),console.warn("[Announcable] Widget will use floating button instead"));return}this.anchors&&this.removeAnchorListeners(),this.anchors=e,this.anchors.forEach(t=>{t.removeEventListener("click",this.boundToggleWidget),t.addEventListener("click",this.boundToggleWidget),t.style.cursor||(t.style.cursor="pointer")}),console.debug(`[Announcable] Attached listeners to ${this.anchors.length} anchor(s)`)}catch(e){console.error("[Announcable] Error setting up anchor listeners:",e)}}setupMutationObserver(){if(this.querySelector)try{this.mutationObserver=new MutationObserver(()=>{this.setupAnchorListeners()}),this.mutationObserver.observe(document.body,{childList:!0,subtree:!0}),this.observerSetupAttempted=!0}catch(e){console.error("[Announcable] Error setting up MutationObserver:",e)}}removeAnchorListeners(){this.anchors&&this.anchors.length>0&&this.anchors.forEach(e=>{try{e.removeEventListener("click",this.boundToggleWidget)}catch(t){console.error("[Announcable] Error removing listener:",t)}})}cleanup(){if(this.removeAnchorListeners(),this.mutationObserver){try{this.mutationObserver.disconnect()}catch(e){console.error("[Announcable] Error disconnecting MutationObserver:",e)}this.mutationObserver=null}this.anchors=null}toggleWidget(e){e&&(e.preventDefault(),e.stopPropagation()),this.setIsOpen(!this.isOpen)}setIsOpen(e){if(this.isOpen=e,e){const t=Ut();this.lastOpened=t,Rt(t)||console.warn("[Announcable] Failed to save last opened timestamp")}this.host.requestUpdate()}}class Ht{constructor(e,t){this.anchors=null,this.mutationObserver=null,this.retryTimeout=null,this.RETRY_DELAY=500,this.host=e,this.querySelector=t,e.addController(this)}hostConnected(){this.queryAnchors(),this.querySelector&&typeof MutationObserver<"u"?this.setupMutationObserver():this.querySelector&&!this.anchors&&this.scheduleRetry()}hostDisconnected(){this.cleanup()}queryAnchors(){if(this.querySelector)try{document.querySelector(this.querySelector);const e=document.querySelectorAll(this.querySelector);if(!e||e.length===0){console.debug(`[Announcable] No anchor elements found for selector: ${this.querySelector}`),this.anchors=null;return}(!this.anchors||this.anchors.length!==e.length||Array.from(this.anchors).some((r,s)=>r!==e[s]))&&(this.anchors=e,console.debug(`[Announcable] Found ${e.length} anchor element(s)`),this.host.requestUpdate())}catch(e){e instanceof DOMException&&e.name==="SyntaxError"?console.error(`[Announcable] Invalid CSS selector: ${this.querySelector}`,e):console.error("[Announcable] Error querying anchor elements:",e),this.anchors=null}}setupMutationObserver(){if(this.querySelector)try{this.mutationObserver=new MutationObserver(()=>{this.queryAnchors()}),this.mutationObserver.observe(document.body,{childList:!0,subtree:!0})}catch(e){console.error("[Announcable] Error setting up MutationObserver:",e)}}scheduleRetry(){this.retryTimeout&&clearTimeout(this.retryTimeout),this.retryTimeout=window.setTimeout(()=>{console.debug("[Announcable] Retrying anchor query..."),this.queryAnchors()},this.RETRY_DELAY)}cleanup(){if(this.mutationObserver){try{this.mutationObserver.disconnect()}catch(e){console.error("[Announcable] Error disconnecting MutationObserver:",e)}this.mutationObserver=null}this.retryTimeout&&(clearTimeout(this.retryTimeout),this.retryTimeout=null),this.anchors=null}}const P="https://release-notes.danielbenner.de";class Wt{constructor(e,t){this.host=e,this.orgId=t,e.addController(this),this.task=new ne(e,async([r])=>{const s=`${P}/api/release-notes/${r}?for=widget`,o=await fetch(s,{method:"GET",headers:{"Content-Type":"application/json"}});if(!o.ok)throw new Error("Failed to fetch release notes");const{data:n}=await o.json();return n||[]},()=>[this.orgId])}hostConnected(){}hostDisconnected(){}}class Bt{constructor(e,t){this.host=e,this.orgId=t,e.addController(this),this.task=new ne(e,async([r])=>{const s=`${P}/api/widget-config/${r}`,o=await fetch(s,{method:"GET",headers:{"Content-Type":"application/json"}});if(!o.ok)throw new Error("Failed to fetch widget config");const{data:n}=await o.json();return{...n,border_radius:parseInt(n.border_radius),border_width:parseInt(n.border_width)}},()=>[this.orgId])}hostConnected(){}hostDisconnected(){}}class Vt{constructor(e,t){this.host=e,this.orgId=t,e.addController(this),this.task=new ne(e,async([r])=>{const s=`${P}/api/release-notes/${r}/status?for=widget`,o=await fetch(s);if(!o.ok)throw new Error("Failed to fetch release note status");return(await o.json()).data||[]},()=>[this.orgId])}hostConnected(){}hostDisconnected(){}}var Gt=Object.defineProperty,Ft=Object.getOwnPropertyDescriptor,Ae=(i,e,t,r)=>{for(var s=r>1?void 0:r?Ft(e,t):e,o=i.length-1,n;o>=0;o--)(n=i[o])&&(s=(r?n(e,t,s):n(s))||s);return r&&s&&Gt(e,t,s),s};let B=class extends d{constructor(){super(...arguments),this.variant="default",this.size="default"}render(){return l`
      <button class="variant-${this.variant} size-${this.size}">
        <slot></slot>
      </button>
    `}};B.styles=g`
    :host {
      display: inline-flex;
      align-items: center;
      justify-content: center;
    }

    button {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      white-space: nowrap;
      border-radius: 0.375rem;
      font-size: 0.875rem;
      font-weight: 500;
      transition: colors 0.2s;
      cursor: pointer;
      border: 1px solid transparent;
      background: none;
    }

    button:focus-visible {
      outline: none;
      ring: 2px solid var(--announcable-border);
      ring-offset: 2px;
    }

    button:disabled {
      pointer-events: none;
      opacity: 0.5;
    }

    /* Variants */
    .variant-default {
      background-color: #0a0a0b;
      color: #ffffff;
      border-color: var(--announcable-border);
    }

    .variant-default:hover {
      background-color: rgba(10, 10, 11, 0.9);
    }

    .variant-ghost:hover {
      background-color: rgba(0, 0, 0, 0.05);
    }

    .variant-link {
      text-decoration: underline;
      text-underline-offset: 4px;
    }

    .variant-link:hover {
      text-decoration: none;
    }

    /* Sizes */
    .size-default {
      height: 2.5rem;
      padding: 0.5rem 1rem;
    }

    .size-sm {
      height: 2.25rem;
      padding: 0.5rem 0.75rem;
    }

    .size-lg {
      height: 2.75rem;
      padding: 0.5rem 2rem;
    }

    .size-icon {
      height: 2.5rem;
      width: 2.5rem;
      padding: 0;
    }
  `,Ae([h({type:String})],B.prototype,"variant",2),Ae([h({type:String})],B.prototype,"size",2),B=Ae([u("ui-button")],B);var Yt=Object.defineProperty,Kt=Object.getOwnPropertyDescriptor,xe=(i,e,t,r)=>{for(var s=r>1?void 0:r?Kt(e,t):e,o=i.length-1,n;o>=0;o--)(n=i[o])&&(s=(r?n(e,t,s):n(s))||s);return r&&s&&Yt(e,t,s),s};let Ce=class extends d{render(){return l`<div class="indicator"></div>`}};Ce.styles=g`
    :host {
      display: inline-block;
    }

    .indicator {
      background-color: #ef4444;
      border-radius: 50%;
      width: 0.375rem;
      height: 0.375rem;
      transform: translate(0.25rem, -0.25rem);
    }
  `,Ce=xe([u("ui-indicator")],Ce);let Oe=class extends d{constructor(){super(...arguments),this.indicatorElement=null}connectedCallback(){super.connectedCallback(),this.attachIndicator()}disconnectedCallback(){super.disconnectedCallback(),this.removeIndicator()}attachIndicator(){if(!this.anchorElement)return;const i=document.createElement("div");i.style.position="absolute",i.style.top="0",i.style.right="0",i.style.backgroundColor="red",i.style.borderRadius="50%",i.style.width="8px",i.style.height="8px",i.style.transform="translate(50%, -50%)",i.style.zIndex="9999",this.indicatorElement=i,this.anchorElement.appendChild(i)}removeIndicator(){if(this.indicatorElement&&this.anchorElement){try{this.anchorElement.removeChild(this.indicatorElement)}catch{}this.indicatorElement=null}}render(){return l``}};xe([h({type:Object})],Oe.prototype,"anchorElement",2),Oe=xe([u("ui-anchor-indicator")],Oe);var Xt=Object.getOwnPropertyDescriptor,T=(i,e,t,r)=>{for(var s=r>1?void 0:r?Xt(e,t):e,o=i.length-1,n;o>=0;o--)(n=i[o])&&(s=n(s)||s);return s};let ke=class extends d{render(){return l`<slot></slot>`}};ke.styles=g`
    :host {
      display: block;
      border-radius: 0.5rem;
      border: 1px solid var(--announcable-border);
      background-color: var(--announcable-card);
      color: var(--announcable-card-foreground);
      box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
    }
  `,ke=T([u("ui-card")],ke);let Ee=class extends d{render(){return l`<slot></slot>`}};Ee.styles=g`
    :host {
      display: flex;
      flex-direction: column;
      gap: 0.375rem;
      padding: 1.5rem;
    }
  `,Ee=T([u("ui-card-header")],Ee);let Se=class extends d{render(){return l`<h3><slot></slot></h3>`}};Se.styles=g`
    :host {
      display: block;
    }

    h3 {
      font-size: 1.5rem;
      font-weight: 600;
      line-height: 1;
      letter-spacing: -0.025em;
      margin: 0;
    }
  `,Se=T([u("ui-card-title")],Se);let Pe=class extends d{render(){return l`<p><slot></slot></p>`}};Pe.styles=g`
    :host {
      display: block;
    }

    p {
      font-size: 0.875rem;
      color: rgba(10, 10, 11, 0.6);
      margin: 0;
    }
  `,Pe=T([u("ui-card-description")],Pe);let Te=class extends d{render(){return l`<slot></slot>`}};Te.styles=g`
    :host {
      display: block;
      padding: 1.5rem;
      padding-top: 0;
    }
  `,Te=T([u("ui-card-content")],Te);let Ie=class extends d{render(){return l`<slot></slot>`}};Ie.styles=g`
    :host {
      display: flex;
      align-items: center;
      padding: 1.5rem;
      padding-top: 0;
    }
  `,Ie=T([u("ui-card-footer")],Ie);var Jt=Object.getOwnPropertyDescriptor,Zt=(i,e,t,r)=>{for(var s=r>1?void 0:r?Jt(e,t):e,o=i.length-1,n;o>=0;o--)(n=i[o])&&(s=n(s)||s);return s};let ht=class extends d{render(){return l`
      <ui-card>
        <ui-card-header>
          <ui-card-title>Error</ui-card-title>
        </ui-card-header>
        <ui-card-content>
          <p>There was an error loading the release notes.</p>
        </ui-card-content>
      </ui-card>
    `}};ht=Zt([u("ui-error-panel")],ht);/**
 * @license
 * Copyright 2020 Google LLC
 * SPDX-License-Identifier: BSD-3-Clause
 */const Qt=i=>i.strings===void 0;/**
 * @license
 * Copyright 2017 Google LLC
 * SPDX-License-Identifier: BSD-3-Clause
 */const es={CHILD:2},ts=i=>(...e)=>({_$litDirective$:i,values:e});class ss{constructor(e){}get _$AU(){return this._$AM._$AU}_$AT(e,t,r){this._$Ct=e,this._$AM=t,this._$Ci=r}_$AS(e,t){return this.update(e,t)}update(e,t){return this.render(...t)}}/**
 * @license
 * Copyright 2017 Google LLC
 * SPDX-License-Identifier: BSD-3-Clause
 */const V=(i,e)=>{var r;const t=i._$AN;if(t===void 0)return!1;for(const s of t)(r=s._$AO)==null||r.call(s,e,!1),V(s,e);return!0},ae=i=>{let e,t;do{if((e=i._$AM)===void 0)break;t=e._$AN,t.delete(i),i=e}while((t==null?void 0:t.size)===0)},dt=i=>{for(let e;e=i._$AM;i=e){let t=e._$AN;if(t===void 0)e._$AN=t=new Set;else if(t.has(i))break;t.add(i),os(e)}};function is(i){this._$AN!==void 0?(ae(this),this._$AM=i,dt(this)):this._$AM=i}function rs(i,e=!1,t=0){const r=this._$AH,s=this._$AN;if(s!==void 0&&s.size!==0)if(e)if(Array.isArray(r))for(let o=t;o<r.length;o++)V(r[o],!1),ae(r[o]);else r!=null&&(V(r,!1),ae(r));else V(this,i)}const os=i=>{i.type==es.CHILD&&(i._$AP??(i._$AP=rs),i._$AQ??(i._$AQ=is))};class ns extends ss{constructor(){super(...arguments),this._$AN=void 0}_$AT(e,t,r){super._$AT(e,t,r),dt(this),this.isConnected=e._$AU}_$AO(e,t=!0){var r,s;e!==this.isConnected&&(this.isConnected=e,e?(r=this.reconnected)==null||r.call(this):(s=this.disconnected)==null||s.call(this)),t&&(V(this,e),ae(this))}setValue(e){if(Qt(this._$Ct))this._$Ct._$AI(e,this);else{const t=[...this._$Ct._$AH];t[this._$Ci]=e,this._$Ct._$AI(t,this,0)}}disconnected(){}reconnected(){}}/**
 * @license
 * Copyright 2020 Google LLC
 * SPDX-License-Identifier: BSD-3-Clause
 */const as=()=>new ls;class ls{}const Ne=new WeakMap,cs=ts(class extends ns{render(i){return b}update(i,[e]){var r;const t=e!==this.G;return t&&this.G!==void 0&&this.rt(void 0),(t||this.lt!==this.ct)&&(this.G=e,this.ht=(r=i.options)==null?void 0:r.host,this.rt(this.ct=i.element)),b}rt(i){if(this.isConnected||(i=void 0),typeof this.G=="function"){const e=this.ht??globalThis;let t=Ne.get(e);t===void 0&&(t=new WeakMap,Ne.set(e,t)),t.get(this.G)!==void 0&&this.G.call(this.ht,void 0),t.set(this.G,i),i!==void 0&&this.G.call(this.ht,i)}else this.G.value=i}get lt(){var i,e;return typeof this.G=="function"?(i=Ne.get(this.ht??globalThis))==null?void 0:i.get(this.G):(e=this.G)==null?void 0:e.value}disconnected(){this.lt===this.ct&&this.rt(void 0)}reconnected(){this.rt(this.ct)}});function je(){return Lt()}class hs{constructor(e,t){this.hasTrackedView=!1,this.observer=null,this.element=null,this.timeoutId=null,this.retryCount=0,this.MAX_RETRIES=3,this.SETUP_DELAY=100,this.RETRY_DELAY=1e3,this.trackCtaClick=()=>{this.sendMetric("cta_click")},this.host=e,this.releaseNoteId=t.releaseNoteId,this.orgId=t.orgId,e.addController(this)}hostConnected(){}hostDisconnected(){this.cleanup()}cleanup(){if(this.observer){try{this.observer.disconnect()}catch(e){console.error("[Announcable] Error disconnecting observer:",e)}this.observer=null}this.timeoutId&&(clearTimeout(this.timeoutId),this.timeoutId=null),this.element=null}setElement(e){if(!e){console.warn("[Announcable] setElement called with null element");return}this.element=e,!this.hasTrackedView&&(this.timeoutId&&clearTimeout(this.timeoutId),this.timeoutId=window.setTimeout(()=>{this.setupIntersectionObserver()},this.SETUP_DELAY))}setupIntersectionObserver(){if(!(!this.element||this.hasTrackedView))try{const e=document.querySelector("[data-scroll-area-viewport]");if(typeof IntersectionObserver>"u"){console.error("[Announcable] IntersectionObserver not supported"),this.hasTrackedView=!0,this.sendMetric("view");return}this.observer=new IntersectionObserver(t=>{t.forEach(r=>{r.isIntersecting&&r.intersectionRatio>=.5&&!this.hasTrackedView&&(this.hasTrackedView=!0,this.sendMetric("view"),this.observer&&(this.observer.disconnect(),this.observer=null))})},{threshold:.5,root:e,rootMargin:"0px"}),this.observer.observe(this.element)}catch(e){if(console.error("[Announcable] Error setting up IntersectionObserver:",e),this.retryCount<this.MAX_RETRIES){this.retryCount++;const t=this.RETRY_DELAY*Math.pow(2,this.retryCount-1);console.log(`[Announcable] Retrying observer setup in ${t}ms (attempt ${this.retryCount}/${this.MAX_RETRIES})`),this.timeoutId=window.setTimeout(()=>{this.setupIntersectionObserver()},t)}else console.error("[Announcable] Max retries reached for IntersectionObserver setup"),this.hasTrackedView=!0,this.sendMetric("view")}}async sendMetric(e){try{const t=je(),r=await fetch(`${P}/api/release-notes/${this.orgId}/metrics`,{method:"POST",headers:{"Content-Type":"application/json"},body:JSON.stringify({release_note_id:this.releaseNoteId,metric_type:e,client_id:t})});if(!r.ok)throw new Error(`HTTP ${r.status}: ${r.statusText}`);console.debug(`[Announcable] Tracked ${e} for release note ${this.releaseNoteId}`)}catch(t){console.error(`[Announcable] Failed to send ${e} metric:`,t)}}}class ds{constructor(e,t){this.isPending=!1,this.error=null,this.host=e,this.releaseNoteId=t.releaseNoteId,this.orgId=t.orgId,this.clientId=t.clientId,e.addController(this),this.task=new ne(e,async([r,s,o])=>{if(!o)return{is_liked:!1};const n=await fetch(`${P}/api/release-notes/${s}/${r}/like?clientId=${o}`,{method:"GET"});if(!n.ok)throw new Error("Failed to get like state");return await n.json()},()=>[this.releaseNoteId,this.orgId,this.clientId])}hostConnected(){}hostDisconnected(){}async toggleLike(){if(this.clientId){this.isPending=!0,this.error=null,this.host.requestUpdate();try{if(!(await fetch(`${P}/api/release-notes/${this.orgId}/${this.releaseNoteId}/like`,{method:"POST",headers:{"Content-Type":"application/json"},body:JSON.stringify({release_note_id:this.releaseNoteId,client_id:this.clientId})})).ok)throw new Error("Failed to toggle like");this.task.run()}catch(e){this.error=e instanceof Error?e:new Error("Unknown error"),console.error("Failed to toggle like:",e)}finally{this.isPending=!1,this.host.requestUpdate()}}}get isLiked(){return this.task.status===oe.COMPLETE&&this.task.value?this.task.value.is_liked:!1}}var us=Object.getOwnPropertyDescriptor,ps=(i,e,t,r)=>{for(var s=r>1?void 0:r?us(e,t):e,o=i.length-1,n;o>=0;o--)(n=i[o])&&(s=n(s)||s);return s};let Me=class extends d{render(){return l`
      <div class="skeleton">
        <slot></slot>
      </div>
    `}};Me.styles=g`
    :host {
      display: block;
    }

    .skeleton {
      animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
      border-radius: 0.375rem;
      background-color: rgba(10, 10, 11, 0.1);
    }

    @keyframes pulse {
      0%, 100% {
        opacity: 1;
      }
      50% {
        opacity: 0.5;
      }
    }
  `,Me=ps([u("ui-skeleton")],Me);var fs=Object.defineProperty,gs=Object.getOwnPropertyDescriptor,le=(i,e,t,r)=>{for(var s=r>1?void 0:r?gs(e,t):e,o=i.length-1,n;o>=0;o--)(n=i[o])&&(s=(r?n(e,t,s):n(s))||s);return r&&s&&fs(e,t,s),s};let I=class extends d{constructor(){super(...arguments),this.class="",this.size=24,this.filled=!1}render(){return ie`
      <svg
        class=${this.class}
        width=${this.size}
        height=${this.size}
        viewBox="0 0 24 24"
        fill=${this.filled?"currentColor":"none"}
        stroke="currentColor"
        stroke-width="2"
        stroke-linecap="round"
        stroke-linejoin="round"
      >
        <path d="M7 10v12"></path>
        <path d="M15 5.88 14 10h5.83a2 2 0 0 1 1.92 2.56l-2.33 8A2 2 0 0 1 17.5 22H4a2 2 0 0 1-2-2v-8a2 2 0 0 1 2-2h2.76a2 2 0 0 0 1.79-1.11L12 2a3.13 3.13 0 0 1 3 3.88Z"></path>
      </svg>
    `}};I.styles=g`
    :host {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      vertical-align: middle;
    }

    svg {
      display: block;
    }
  `,le([h({type:String})],I.prototype,"class",2),le([h({type:Number})],I.prototype,"size",2),le([h({type:Boolean})],I.prototype,"filled",2),I=le([u("icon-thumbs-up")],I);var bs=Object.defineProperty,ms=Object.getOwnPropertyDescriptor,De=(i,e,t,r)=>{for(var s=r>1?void 0:r?ms(e,t):e,o=i.length-1,n;o>=0;o--)(n=i[o])&&(s=(r?n(e,t,s):n(s))||s);return r&&s&&bs(e,t,s),s};let G=class extends d{constructor(){super(...arguments),this.class="",this.size=24}render(){return ie`
      <svg 
        class=${this.class}
        width=${this.size} 
        height=${this.size} 
        viewBox="0 0 24 24" 
        fill="none" 
        stroke="currentColor" 
        stroke-width="2" 
        stroke-linecap="round" 
        stroke-linejoin="round"
      >
        <path d="M15 3h6v6"></path>
        <path d="M10 14 21 3"></path>
        <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"></path>
      </svg>
    `}};G.styles=g`
    :host {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      vertical-align: middle;
    }

    svg {
      display: block;
    }
  `,De([h({type:String})],G.prototype,"class",2),De([h({type:Number})],G.prototype,"size",2),G=De([u("icon-external-link")],G);var $s=Object.defineProperty,_s=Object.getOwnPropertyDescriptor,N=(i,e,t,r)=>{for(var s=r>1?void 0:r?_s(e,t):e,o=i.length-1,n;o>=0;o--)(n=i[o])&&(s=(r?n(e,t,s):n(s))||s);return r&&s&&$s(e,t,s),s};let ze=class extends d{render(){return l`
      <div class="list">
        <slot></slot>
      </div>
    `}};ze.styles=g`
    :host {
      display: block;
    }

    .list {
      display: flex;
      flex-direction: column;
      gap: 1.5rem;
    }
  `,ze=N([u("release-notes-list")],ze);let F=class extends d{constructor(){super(...arguments),this.cardRef=as()}connectedCallback(){super.connectedCallback(),this.metricsController=new hs(this,{releaseNoteId:this.releaseNote.id,orgId:this.config.org_id});const i=je();this.likesController=new ds(this,{releaseNoteId:this.releaseNote.id,orgId:this.config.org_id,clientId:i})}updated(){this.cardRef.value&&this.metricsController.setElement(this.cardRef.value)}handleImageError(i){const e=i.target;console.error(`Image failed to load for ${this.releaseNote.title}`,this.releaseNote.imageSrc,i),e.style.display="none"}render(){const i=this.releaseNote.cta_label_override||this.config.cta_text,e=this.config.release_page_baseurl,t=this.releaseNote.cta_href_override||`${e}#${this.releaseNote.id}`,r=je();return l`
      <ui-card
        ${cs(this.cardRef)}
        style="
          border-radius: ${this.config.release_note_border_radius}px;
          border-color: ${this.config.release_note_border_color};
          border-width: ${this.config.release_note_border_width}px;
          color: ${this.config.release_note_font_color};
          background-color: ${this.config.release_note_bg_color};
        "
      >
        <ui-card-header style="padding-bottom: 1rem;">
          <ui-card-title>${this.releaseNote.title}</ui-card-title>
          <ui-card-description style="color: ${this.config.release_note_font_color}">
            ${this.releaseNote.date||""}
          </ui-card-description>
        </ui-card-header>
        <ui-card-content>
          <div class="content">
            ${this.releaseNote.media_link?l`
              <div class="media-container">
                <iframe
                  src="${this.releaseNote.media_link}"
                  allow="fullscreen"
                  allowfullscreen
                  loading="lazy"
                  referrerpolicy="no-referrer"
                  sandbox="allow-scripts allow-presentation allow-same-origin"
                  title="${this.releaseNote.title}"
                ></iframe>
              </div>
            `:this.releaseNote.imageSrc?l`
              <div class="media-container">
                <img
                  src="${this.releaseNote.imageSrc}"
                  alt="${this.releaseNote.title}"
                  @error=${this.handleImageError}
                />
              </div>
            `:""}

            ${this.releaseNote.text?l`
              <div class="text">${this.releaseNote.text}</div>
            `:""}

            ${this.config.enable_likes||!this.releaseNote.hide_cta?l`
              <div class="actions">
                ${this.config.enable_likes?l`
                  <div class="action-wrapper">
                    <button
                      class="like-button"
                      @click=${()=>this.likesController.toggleLike()}
                      ?disabled=${!r}
                    >
                      <span>
                        ${this.likesController.isLiked?this.config.unlike_button_text:this.config.like_button_text}
                      </span>
                      <icon-thumbs-up
                        class="icon-small"
                        .size=${12}
                        ?filled=${this.likesController.isLiked}
                      ></icon-thumbs-up>
                    </button>
                  </div>
                `:""}

                ${this.releaseNote.hide_cta?"":l`
                  <div class="action-wrapper">
                    <a
                      class="cta-link"
                      href="${t}"
                      target="_blank"
                      @click=${this.metricsController.trackCtaClick}
                    >
                      <span class="cta-content">
                        ${i}
                        <icon-external-link class="icon-small" .size=${12}></icon-external-link>
                      </span>
                    </a>
                  </div>
                `}
              </div>
            `:""}
          </div>
        </ui-card-content>
      </ui-card>
    `}};F.styles=g`
    :host {
      display: block;
    }

    .content {
      width: 100%;
      display: flex;
      flex-direction: column;
      gap: 1rem;
    }

    .media-container {
      position: relative;
      width: 100%;
      aspect-ratio: 16 / 9;
    }

    .media-container iframe {
      position: absolute;
      top: 0;
      left: 0;
      width: 100%;
      height: 100%;
    }

    .media-container img {
      width: 100%;
      height: auto;
    }

    .text {
      white-space: pre-wrap;
    }

    .actions {
      width: 100%;
      display: flex;
      justify-content: space-around;
      margin-top: 0.5rem;
    }

    .action-wrapper {
      width: 100%;
      display: flex;
      justify-content: center;
      margin: 0 auto;
    }

    .like-button {
      display: inline-flex;
      align-items: center;
      gap: 0.25rem;
      background: none;
      border: none;
      cursor: pointer;
      padding: 0;
      font: inherit;
      color: inherit;
    }

    .like-button:disabled {
      opacity: 0.5;
      cursor: not-allowed;
    }

    .like-button span {
      font-size: 0.875rem;
    }

    .cta-link {
      text-decoration: none;
      color: inherit;
    }

    .cta-content {
      display: flex;
      align-items: center;
      gap: 0.25rem;
    }

    .icon-small {
      width: 0.75rem;
      height: 0.75rem;
    }
  `,N([h({type:Object})],F.prototype,"config",2),N([h({type:Object})],F.prototype,"releaseNote",2),F=N([u("release-note-entry")],F);let ce=class extends d{render(){const i=this.config.widget_bg_color,e=this.config.release_note_border_radius;return l`
      <ui-card
        style="
          border-radius: ${this.config.release_note_border_radius}px;
          border-color: ${this.config.release_note_border_color};
          border-width: ${this.config.release_note_border_width}px;
          color: ${this.config.release_note_font_color};
          background-color: ${this.config.release_note_bg_color};
        "
      >
        <ui-card-header style="padding-bottom: 1rem;">
          <div class="header-skeletons">
            <ui-skeleton
              class="sk-title"
              style="
                background-color: ${i};
                border-radius: ${e}px;
              "
            ></ui-skeleton>
            <ui-skeleton
              class="sk-date"
              style="
                background-color: ${i};
                border-radius: ${e}px;
              "
            ></ui-skeleton>
          </div>
        </ui-card-header>
        <ui-card-content>
          <div class="content">
            <ui-skeleton
              class="sk-image"
              style="
                background-color: ${i};
                border-radius: ${e}px;
              "
            ></ui-skeleton>
            <div class="text-skeletons">
              <ui-skeleton
                class="sk-text"
                style="
                  background-color: ${i};
                  border-radius: ${e}px;
                "
              ></ui-skeleton>
              <ui-skeleton
                class="sk-text"
                style="
                  background-color: ${i};
                  border-radius: ${e}px;
                "
              ></ui-skeleton>
              <ui-skeleton
                class="sk-text"
                style="
                  background-color: ${i};
                  border-radius: ${e}px;
                "
              ></ui-skeleton>
              <ui-skeleton
                class="sk-text-short"
                style="
                  background-color: ${i};
                  border-radius: ${e}px;
                "
              ></ui-skeleton>
            </div>
            <div class="center">
              <ui-skeleton
                class="sk-cta"
                style="
                  background-color: ${i};
                  border-radius: ${e}px;
                "
              ></ui-skeleton>
            </div>
          </div>
        </ui-card-content>
      </ui-card>
    `}};ce.styles=g`
    :host {
      display: block;
    }

    .header-skeletons {
      display: flex;
      flex-direction: column;
      gap: 0.5rem;
    }

    .content {
      width: 100%;
      display: flex;
      flex-direction: column;
      gap: 1rem;
    }

    .text-skeletons {
      display: flex;
      flex-direction: column;
      gap: 0.5rem;
    }

    .center {
      width: 100%;
      display: flex;
      justify-content: center;
    }

    .sk-title {
      height: 1.75rem;
      width: 75%;
    }

    .sk-date {
      height: 1rem;
      width: 25%;
    }

    .sk-image {
      height: 12rem;
      width: 100%;
    }

    .sk-text {
      height: 1rem;
      width: 100%;
    }

    .sk-text-short {
      height: 1rem;
      width: 75%;
    }

    .sk-cta {
      height: 1rem;
      width: 6rem;
    }
  `,N([h({type:Object})],ce.prototype,"config",2),ce=N([u("release-note-skeleton")],ce);var vs=Object.getOwnPropertyDescriptor,ys=(i,e,t,r)=>{for(var s=r>1?void 0:r?vs(e,t):e,o=i.length-1,n;o>=0;o--)(n=i[o])&&(s=n(s)||s);return s};let Re=class extends d{render(){return l`
      <div
        class="scroll-container"
        data-scroll-area-viewport
      >
        <slot></slot>
      </div>
    `}};Re.styles=g`
    :host {
      display: block;
      height: var(--scroll-area-height, auto);
    }

    .scroll-container {
      position: relative;
      overflow-y: auto;
      overflow-x: hidden;
      height: 100%;
      max-height: var(--scroll-area-max-height, none);
    }

    /* Custom scrollbar styling */
    .scroll-container::-webkit-scrollbar {
      width: 10px;
    }

    .scroll-container::-webkit-scrollbar-track {
      background: transparent;
    }

    .scroll-container::-webkit-scrollbar-thumb {
      background: var(--announcable-border);
      border-radius: 5px;
    }

    .scroll-container::-webkit-scrollbar-thumb:hover {
      background: rgba(229, 229, 229, 0.8);
    }
  `,Re=ys([u("ui-scroll-area")],Re);var ws=Object.defineProperty,As=Object.getOwnPropertyDescriptor,Ue=(i,e,t,r)=>{for(var s=r>1?void 0:r?As(e,t):e,o=i.length-1,n;o>=0;o--)(n=i[o])&&(s=(r?n(e,t,s):n(s))||s);return r&&s&&ws(e,t,s),s};let Y=class extends d{constructor(){super(...arguments),this.class="",this.size=24}render(){return ie`
      <svg 
        class=${this.class}
        width=${this.size} 
        height=${this.size} 
        viewBox="0 0 24 24" 
        fill="none" 
        stroke="currentColor" 
        stroke-width="2" 
        stroke-linecap="round" 
        stroke-linejoin="round"
      >
        <path d="M18 6 6 18"></path>
        <path d="m6 6 12 12"></path>
      </svg>
    `}};Y.styles=g`
    :host {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      vertical-align: middle;
    }

    svg {
      display: block;
    }
  `,Ue([h({type:String})],Y.prototype,"class",2),Ue([h({type:Number})],Y.prototype,"size",2),Y=Ue([u("icon-x")],Y);var xs=Object.defineProperty,Cs=Object.getOwnPropertyDescriptor,Le=(i,e,t,r)=>{for(var s=r>1?void 0:r?Cs(e,t):e,o=i.length-1,n;o>=0;o--)(n=i[o])&&(s=(r?n(e,t,s):n(s))||s);return r&&s&&xs(e,t,s),s};let K=class extends d{handleClose(){this.dispatchEvent(new CustomEvent("close",{bubbles:!0,composed:!0}))}render(){var t;const{title:i,description:e}=this.config;return l`
      <div
        class="popover"
        style="
          --widget-border-radius: ${this.config.widget_border_radius}px;
          --widget-border-width: ${this.config.widget_border_width}px;
          --widget-border-color: ${this.config.widget_border_color};
          --widget-bg-color: ${this.config.widget_bg_color};
          --widget-font-color: ${this.config.widget_font_color};
          --widget-font-family: ${((t=this.init.font_family)==null?void 0:t.join(","))||"inherit"};
        "
      >
        <div class="actions">
          <ui-button size="icon" variant="ghost">
            <icon-external-link class="icon-medium" .size=${16}></icon-external-link>
          </ui-button>
          <ui-button size="icon" variant="ghost" @click=${this.handleClose}>
            <icon-x class="icon-medium" .size=${16}></icon-x>
          </ui-button>
        </div>
        <ui-card-header>
          <ui-card-title class="title">
            ${i}
          </ui-card-title>
          ${e?l`
            <ui-card-description style="color: ${this.config.widget_font_color}">
              ${e}
            </ui-card-description>
          `:""}
        </ui-card-header>
        <ui-card-content>
          <ui-scroll-area
            class="scroll-content"
            style="--scroll-area-height: 32rem;"
          >
            <slot></slot>
          </ui-scroll-area>
        </ui-card-content>
      </div>
    `}};K.styles=g`
    :host {
      display: block;
    }

    a {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      text-decoration: none;
      color: inherit;
    }

    .popover {
      width: 32rem;
      position: fixed;
      bottom: 5rem;
      right: 1rem;
      z-index: 9999;
      border-radius: var(--widget-border-radius, 0.5rem);
      border-width: var(--widget-border-width, 1px);
      border-style: solid;
      border-color: var(--widget-border-color, #e5e5e5);
      background-color: var(--widget-bg-color, #ffffff);
      color: var(--widget-font-color, #0a0a0b);
      font-family: var(--widget-font-family, inherit);
    }

    .actions {
      position: absolute;
      top: 0;
      right: 0;
      padding: 0.5rem;
      display: flex;
      align-items: center;
      gap: 0.25rem;
    }

    .icon-medium {
      width: 1rem;
      height: 1rem;
    }

    .title {
      font-size: 1.125rem;
      display: inline-flex;
      align-items: center;
    }

    .scroll-content {
      height: 32rem;
    }
  `,Le([h({type:Object})],K.prototype,"config",2),Le([h({type:Object})],K.prototype,"init",2),K=Le([u("widget-popover")],K);var Os=Object.defineProperty,ks=Object.getOwnPropertyDescriptor,he=(i,e,t,r)=>{for(var s=r>1?void 0:r?ks(e,t):e,o=i.length-1,n;o>=0;o--)(n=i[o])&&(s=(r?n(e,t,s):n(s))||s);return r&&s&&Os(e,t,s),s};let j=class extends d{constructor(){super(...arguments),this.isOpen=!1,this.title="",this.description=""}handleBackdropClick(i){i.target===i.currentTarget&&this.dispatchEvent(new CustomEvent("close",{bubbles:!0,composed:!0}))}handleKeyDown(i){i.key==="Escape"&&this.dispatchEvent(new CustomEvent("close",{bubbles:!0,composed:!0}))}connectedCallback(){super.connectedCallback(),document.addEventListener("keydown",this.handleKeyDown.bind(this))}disconnectedCallback(){super.disconnectedCallback(),document.removeEventListener("keydown",this.handleKeyDown.bind(this))}render(){return this.isOpen?l`
      <div class="backdrop" @click=${this.handleBackdropClick}>
        <dialog
          open
          role="dialog"
          aria-modal="true"
          aria-labelledby="dialog-title"
          tabindex="-1"
        >
          <div class="content">
            <div class="header">
              <div class="actions">
                <slot name="actions"></slot>
              </div>
              <h2 id="dialog-title">${this.title}</h2>
              ${this.description?l` <p class="description">${this.description}</p> `:""}
            </div>
            <slot></slot>
          </div>
        </dialog>
      </div>
    `:l``}};j.styles=g`
    :host {
      display: contents;
    }

    .backdrop {
      position: fixed;
      inset: 0;
      background-color: rgba(0, 0, 0, 0.7);
      backdrop-filter: blur(4px);
      display: flex;
      align-items: center;
      justify-content: center;
      z-index: 9999;
    }

    dialog {
      border: none;
      background: transparent;
      padding: 1rem;
      max-width: 40rem;
      width: 100%;
      box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1);
    }

    .content {
      display: grid;
      gap: 1rem;
      background-color: var(--dialog-bg-color, #ffffff);
      color: var(--dialog-font-color, #0a0a0b);
      border-radius: var(--dialog-border-radius, 0.5rem);
      border-width: var(--dialog-border-width, 1px);
      border-style: solid;
      border-color: var(--dialog-border-color, #e5e5e5);
      padding: 1.5rem;
      font-family: var(--dialog-font-family, inherit);
    }

    .header {
      position: relative;
    }

    .actions {
      position: absolute;
      top: 0;
      right: 0;
      transform: translate(12px, -12px);
      display: flex;
      align-items: center;
      gap: 0.25rem;
    }

    h2 {
      font-size: 1.125rem;
      font-weight: 600;
      letter-spacing: -0.025em;
      margin: 0;
    }

    .description {
      margin-top: 0.375rem;
      font-size: 0.875rem;
      color: rgba(10, 10, 11, 0.6);
    }
  `,he([h({type:Boolean})],j.prototype,"isOpen",2),he([h({type:String})],j.prototype,"title",2),he([h({type:String})],j.prototype,"description",2),j=he([u("ui-dialog")],j);var Es=Object.defineProperty,Ss=Object.getOwnPropertyDescriptor,de=(i,e,t,r)=>{for(var s=r>1?void 0:r?Ss(e,t):e,o=i.length-1,n;o>=0;o--)(n=i[o])&&(s=(r?n(e,t,s):n(s))||s);return r&&s&&Es(e,t,s),s};let M=class extends d{constructor(){super(...arguments),this.isOpen=!1}handleClose(){this.dispatchEvent(new CustomEvent("close",{bubbles:!0,composed:!0}))}render(){var i;return l`
      <ui-dialog
        .isOpen=${this.isOpen}
        .title=${this.config.title}
        .description=${this.config.description}
        @close=${this.handleClose}
        style="
          --dialog-border-radius: ${this.config.widget_border_radius}px;
          --dialog-border-color: ${this.config.widget_border_color};
          --dialog-border-width: ${this.config.widget_border_width}px;
          --dialog-bg-color: ${this.config.widget_bg_color};
          --dialog-font-color: ${this.config.widget_font_color};
          --dialog-font-family: ${((i=this.init.font_family)==null?void 0:i.join(","))||"inherit"};
        "
      >
        <div slot="actions">
          <ui-button size="icon" variant="ghost">
            <a href="${this.config.release_page_baseurl}" target="_blank">
              <icon-external-link class="icon-medium" .size=${16}></icon-external-link>
            </a>
          </ui-button>
          <ui-button size="icon" variant="ghost" @click=${this.handleClose}>
            <icon-x class="icon-medium" .size=${16}></icon-x>
          </ui-button>
        </div>
        <ui-scroll-area
          class="scroll-content"
          style="--scroll-area-height: 32rem;"
        >
          <slot></slot>
        </ui-scroll-area>
      </ui-dialog>
    `}};M.styles=g`
    :host {
      display: block;
    }

    a {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      text-decoration: none;
      color: inherit;
    }

    .icon-medium {
      width: 1rem;
      height: 1rem;
    }

    .scroll-content {
      height: 32rem;
    }
  `,de([h({type:Object})],M.prototype,"config",2),de([h({type:Object})],M.prototype,"init",2),de([h({type:Boolean})],M.prototype,"isOpen",2),M=de([u("widget-modal")],M);var Ps=Object.defineProperty,Ts=Object.getOwnPropertyDescriptor,ue=(i,e,t,r)=>{for(var s=r>1?void 0:r?Ts(e,t):e,o=i.length-1,n;o>=0;o--)(n=i[o])&&(s=(r?n(e,t,s):n(s))||s);return r&&s&&Ps(e,t,s),s};let D=class extends d{constructor(){super(...arguments),this.isOpen=!1}handleClose(){this.dispatchEvent(new CustomEvent("close",{bubbles:!0,composed:!0}))}render(){var t;const{title:i,description:e}=this.config;return l`
      <div
        class="sidebar ${this.isOpen?"open":"closed"}"
        style="
          --widget-bg-color: ${this.config.widget_bg_color};
          --widget-font-color: ${this.config.widget_font_color};
          --widget-font-family: ${((t=this.init.font_family)==null?void 0:t.join(","))||"inherit"};
        "
      >
        <div class="actions">
          <ui-button size="icon" variant="ghost">
            <icon-external-link class="icon-medium" .size=${16}></icon-external-link>
          </ui-button>
          <ui-button size="icon" variant="ghost" @click=${this.handleClose}>
            <icon-x class="icon-medium" .size=${16}></icon-x>
          </ui-button>
        </div>
        <ui-card-header>
          <ui-card-title class="title">
            ${i}
          </ui-card-title>
          ${e?l`
            <ui-card-description style="color: ${this.config.widget_font_color}">
              ${e}
            </ui-card-description>
          `:""}
        </ui-card-header>
        <ui-card-content class="content">
          <slot></slot>
        </ui-card-content>
      </div>
    `}};D.styles=g`
    :host {
      display: block;
    }

    a {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      text-decoration: none;
      color: inherit;
    }

    .sidebar {
      width: 32rem;
      position: fixed;
      top: 0;
      right: 0;
      height: 100vh;
      transition: transform 0.3s ease-in-out;
      z-index: 9999;
      border-radius: 0;
      background-color: var(--widget-bg-color, #ffffff);
      color: var(--widget-font-color, #0a0a0b);
      font-family: var(--widget-font-family, inherit);
    }

    .sidebar.open {
      transform: translateX(0);
    }

    .sidebar.closed {
      transform: translateX(100%);
    }

    .actions {
      position: fixed;
      top: 0;
      right: 0;
      padding: 0.5rem;
      display: flex;
      align-items: center;
      gap: 0.25rem;
    }

    .icon-medium {
      width: 1rem;
      height: 1rem;
    }

    .title {
      font-size: 1.125rem;
      display: inline-flex;
      align-items: center;
    }

    .content {
      height: calc(100vh - 120px);
      overflow-y: auto;
    }
  `,ue([h({type:Object})],D.prototype,"config",2),ue([h({type:Object})],D.prototype,"init",2),ue([h({type:Boolean})],D.prototype,"isOpen",2),D=ue([u("widget-sidebar")],D);var Is=Object.defineProperty,Ns=Object.getOwnPropertyDescriptor,pe=(i,e,t,r)=>{for(var s=r>1?void 0:r?Ns(e,t):e,o=i.length-1,n;o>=0;o--)(n=i[o])&&(s=(r?n(e,t,s):n(s))||s);return r&&s&&Is(e,t,s),s};let X=class extends d{constructor(){super(...arguments),this.isOpen=!1}handleClose(){this.dispatchEvent(new CustomEvent("close",{bubbles:!0,composed:!0}))}render(){var e;switch(((e=this.config)==null?void 0:e.widget_type)||"popover"){case"sidebar":return l`
          <widget-sidebar
            .config=${this.config}
            .init=${this.init}
            .isOpen=${this.isOpen}
            @close=${this.handleClose}
          >
            <slot></slot>
          </widget-sidebar>
        `;case"modal":return l`
          <widget-modal
            .config=${this.config}
            .init=${this.init}
            .isOpen=${this.isOpen}
            @close=${this.handleClose}
          >
            <slot></slot>
          </widget-modal>
        `;default:return l`
          <widget-popover
            .config=${this.config}
            .init=${this.init}
            @close=${this.handleClose}
          >
            <slot></slot>
          </widget-popover>
        `}}};pe([h({type:Object})],X.prototype,"config",2),pe([h({type:Object})],X.prototype,"init",2),pe([h({type:Boolean})],X.prototype,"isOpen",2),X=pe([u("widget-container")],X);var js=Object.defineProperty,Ms=Object.getOwnPropertyDescriptor,qe=(i,e,t,r)=>{for(var s=r>1?void 0:r?Ms(e,t):e,o=i.length-1,n;o>=0;o--)(n=i[o])&&(s=(r?n(e,t,s):n(s))||s);return r&&s&&js(e,t,s),s};let J=class extends d{constructor(){super(...arguments),this.class="",this.size=24}render(){return ie`
      <svg 
        class=${this.class}
        width=${this.size} 
        height=${this.size} 
        viewBox="0 0 24 24" 
        fill="none" 
        stroke="currentColor" 
        stroke-width="2" 
        stroke-linecap="round" 
        stroke-linejoin="round"
      >
        <rect x="3" y="8" width="18" height="4" rx="1"></rect>
        <path d="M12 8v13"></path>
        <path d="M19 12v7a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2v-7"></path>
        <path d="M7.5 8a2.5 2.5 0 0 1 0-5A4.8 8 0 0 1 12 8a4.8 8 0 0 1 4.5-5 2.5 2.5 0 0 1 0 5"></path>
      </svg>
    `}};J.styles=g`
    :host {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      vertical-align: middle;
    }

    svg {
      display: block;
    }
  `,qe([h({type:String})],J.prototype,"class",2),qe([h({type:Number})],J.prototype,"size",2),J=qe([u("icon-gift")],J);var Ds=Object.defineProperty,zs=Object.getOwnPropertyDescriptor,Z=(i,e,t,r)=>{for(var s=r>1?void 0:r?zs(e,t):e,o=i.length-1,n;o>=0;o--)(n=i[o])&&(s=(r?n(e,t,s):n(s))||s);return r&&s&&Ds(e,t,s),s};let fe=class extends d{connectedCallback(){super.connectedCallback(),this.toggleController=new qt(this,this.init.anchor_query_selector),this.anchorsController=new Ht(this,this.init.anchor_query_selector),this.statusTask=new Vt(this,this.init.org_id)}updated(){const i=this.hasUnseenReleaseNotes();this.anchorsController.anchors&&this.anchorsController.anchors.forEach(e=>{this.updateIndicatorDataset(e,i)}),this.shouldInstantOpen()&&!this.toggleController.isOpen&&this.anchorsController.anchors&&(this.toggleController.setIsOpen(!0),this.anchorsController.anchors.forEach(e=>{this.updateInstantOpenDataset(e,!0)}))}hasUnseenReleaseNotes(){if(this.statusTask.task.status!==oe.COMPLETE)return!1;const i=this.statusTask.task.value,e=this.toggleController.lastOpened;if(!i)return!1;if(!e)return!0;const t=parseInt(e);return i.some(r=>r.last_update_on?new Date(r.last_update_on).getTime()>t:!1)}shouldInstantOpen(){if(this.statusTask.task.status!==oe.COMPLETE)return!1;const i=this.statusTask.task.value,e=this.toggleController.lastOpened;if(!i)return!1;if(!e)return!0;const t=parseInt(e);return i.some(r=>!r.last_update_on||r.attention_mechanism!=="instant_open"?!1:new Date(r.last_update_on).getTime()>t)}updateIndicatorDataset(i,e){if(!i)return;const t=e?"true":"false";i.dataset.newReleaseNotes=t}updateInstantOpenDataset(i,e){if(!i)return;const t=e?"true":"false";i.dataset.instantOpen=t}render(){const i=this.hasUnseenReleaseNotes(),e=this.anchorsController.anchors&&!this.init.hide_indicator&&i;return l`
      ${e?Array.from(this.anchorsController.anchors).map(t=>l`
              <ui-anchor-indicator .anchorElement=${t}></ui-anchor-indicator>
            `):""}

      ${this.init.anchor_query_selector?"":l`
            <ui-button
              class="floating-button"
              size="icon"
              style="border-radius: 0.5rem"
              @click=${()=>this.toggleController.setIsOpen(!this.toggleController.isOpen)}
            >
              ${i?l`<ui-indicator class="indicator-wrapper"></ui-indicator>`:""}
              <icon-gift class="icon-medium" .size=${16}></icon-gift>
            </ui-button>
          `}

      ${this.toggleController.isOpen?l`
            <widget-content
              .init=${this.init}
              .isOpen=${this.toggleController.isOpen}
              @close=${()=>this.toggleController.setIsOpen(!1)}
            ></widget-content>
          `:""}
    `}};fe.styles=g`
    :host {
      display: block;
    }

    .floating-button {
      position: fixed;
      z-index: 50;
      bottom: 1rem;
      right: 1rem;
    }

    .indicator-wrapper {
      position: absolute;
      top: 0;
      right: 0;
    }

    .icon-medium {
      width: 1rem;
      height: 1rem;
    }
  `,Z([h({type:Object})],fe.prototype,"init",2),fe=Z([u("announcable-app")],fe);let Q=class extends d{constructor(){super(...arguments),this.isOpen=!1}connectedCallback(){super.connectedCallback(),this.notesTask=new Wt(this,this.init.org_id),this.configTask=new Bt(this,this.init.org_id)}handleClose(){this.dispatchEvent(new CustomEvent("close",{bubbles:!0,composed:!0}))}render(){const i=this.configTask.task.status===oe.PENDING,e=this.notesTask.task.error,t=this.configTask.task.error;if(e&&console.error(e),t&&console.error(t),!(!i&&!e&&!t))return l``;const s=this.configTask.task.value;return s?l`
      <div class="container">
        <widget-container
          .config=${s}
          .init=${this.init}
          .isOpen=${this.isOpen}
          @close=${this.handleClose}
        >
          ${e||t?l`<ui-error-panel></ui-error-panel>`:this.notesTask.task.render({pending:()=>l`
                  <release-notes-list>
                    <release-note-skeleton .config=${s}></release-note-skeleton>
                  </release-notes-list>
                `,complete:o=>l`
                  <release-notes-list>
                    ${o.map(n=>l`
                        <release-note-entry
                          .config=${s}
                          .releaseNote=${n}
                        ></release-note-entry>
                      `)}
                  </release-notes-list>
                `,error:o=>(console.error("Error loading release notes:",o),l`<ui-error-panel></ui-error-panel>`)})}
        </widget-container>
      </div>
    `:l``}};Q.styles=g`
    :host {
      display: block;
    }

    .container {
      position: relative;
    }
  `,Z([h({type:Object})],Q.prototype,"init",2),Z([h({type:Boolean})],Q.prototype,"isOpen",2),Q=Z([u("widget-content")],Q);function Rs(i){return!i||i.length===0?"":`
    .announcable-widget {
      font-family: ${i.map(r=>r.includes(" ")&&!r.startsWith('"')&&!r.startsWith("'")?`"${r}"`:r).join(", ")} !important;
    }
    
    .announcable-widget * {
      font-family: inherit !important;
    }
  `}function ut(i){if(!i.org_id){console.error("[Announcable] org_id is required to initialize widget");return}try{const e=document.createElement("div");if(e.id="announcable-widget-root",e.className="announcable-widget",document.body.appendChild(e),i.font_family&&i.font_family.length>0){const t=document.createElement("style");t.textContent=Rs(i.font_family),document.head.appendChild(t),console.debug(`[Announcable] Applied custom fonts: ${i.font_family.join(", ")}`)}rt(l`<announcable-app .init=${i}></announcable-app>`,e),console.debug("[Announcable] Widget initialized successfully")}catch(e){console.error("[Announcable] Failed to initialize widget:",e)}}window.AnnouncableWidget={init:i=>{ut(i)}},window.ReleaseBeaconWidget={init:i=>{console.warn("[Announcable] ReleaseBeaconWidget is deprecated, use AnnouncableWidget instead"),ut(i)}},window.announcable_init&&((ft=window.AnnouncableWidget)!=null&&ft.init)?(console.log("[Announcable] Auto-initializing with announcable_init config"),window.AnnouncableWidget.init(window.announcable_init)):window.release_beacon_widget_init&&((gt=window.ReleaseBeaconWidget)!=null&&gt.init)?(console.log("[Announcable] Auto-initializing with legacy release_beacon_widget_init config"),window.ReleaseBeaconWidget.init(window.release_beacon_widget_init)):console.warn("[Announcable] No auto-init config found. Call AnnouncableWidget.init() manually.")});
