"use strict";(self.webpackChunkjalapeno=self.webpackChunkjalapeno||[]).push([[827],{3905:(e,t,a)=>{a.d(t,{Zo:()=>s,kt:()=>g});var r=a(7294);function n(e,t,a){return t in e?Object.defineProperty(e,t,{value:a,enumerable:!0,configurable:!0,writable:!0}):e[t]=a,e}function i(e,t){var a=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),a.push.apply(a,r)}return a}function l(e){for(var t=1;t<arguments.length;t++){var a=null!=arguments[t]?arguments[t]:{};t%2?i(Object(a),!0).forEach((function(t){n(e,t,a[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(a)):i(Object(a)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(a,t))}))}return e}function o(e,t){if(null==e)return{};var a,r,n=function(e,t){if(null==e)return{};var a,r,n={},i=Object.keys(e);for(r=0;r<i.length;r++)a=i[r],t.indexOf(a)>=0||(n[a]=e[a]);return n}(e,t);if(Object.getOwnPropertySymbols){var i=Object.getOwnPropertySymbols(e);for(r=0;r<i.length;r++)a=i[r],t.indexOf(a)>=0||Object.prototype.propertyIsEnumerable.call(e,a)&&(n[a]=e[a])}return n}var p=r.createContext({}),c=function(e){var t=r.useContext(p),a=t;return e&&(a="function"==typeof e?e(t):l(l({},t),e)),a},s=function(e){var t=c(e.components);return r.createElement(p.Provider,{value:t},e.children)},u="mdxType",d={inlineCode:"code",wrapper:function(e){var t=e.children;return r.createElement(r.Fragment,{},t)}},m=r.forwardRef((function(e,t){var a=e.components,n=e.mdxType,i=e.originalType,p=e.parentName,s=o(e,["components","mdxType","originalType","parentName"]),u=c(a),m=n,g=u["".concat(p,".").concat(m)]||u[m]||d[m]||i;return a?r.createElement(g,l(l({ref:t},s),{},{components:a})):r.createElement(g,l({ref:t},s))}));function g(e,t){var a=arguments,n=t&&t.mdxType;if("string"==typeof e||n){var i=a.length,l=new Array(i);l[0]=m;var o={};for(var p in t)hasOwnProperty.call(t,p)&&(o[p]=t[p]);o.originalType=e,o[u]="string"==typeof e?e:n,l[1]=o;for(var c=2;c<i;c++)l[c]=a[c];return r.createElement.apply(null,l)}return r.createElement.apply(null,a)}m.displayName="MDXCreateElement"},2175:(e,t,a)=>{a.r(t),a.d(t,{assets:()=>p,contentTitle:()=>l,default:()=>d,frontMatter:()=>i,metadata:()=>o,toc:()=>c});var r=a(7462),n=(a(7294),a(3905));const i={sidebar_position:3,slug:"/usage/",title:"Usage"},l="Usage [WIP]",o={unversionedId:"usage",id:"usage",title:"Usage",description:"Getting started",source:"@site/docs/usage.md",sourceDirName:".",slug:"/usage/",permalink:"/jalapeno/usage/",draft:!1,editUrl:"https://github.com/futurice/jalapeno/tree/main/docs/site/docs/usage.md",tags:[],version:"current",sidebarPosition:3,frontMatter:{sidebar_position:3,slug:"/usage/",title:"Usage"},sidebar:"tutorialSidebar",previous:{title:"Installation",permalink:"/jalapeno/installation/"},next:{title:"API Reference",permalink:"/jalapeno/api/"}},p={},c=[{value:"Getting started",id:"getting-started",level:2},{value:"Variables",id:"variables",level:2},{value:"Validation",id:"validation",level:3},{value:"Conditional variables",id:"conditional-variables",level:3},{value:"Templating",id:"templating",level:2},{value:"Pushing recipe to Container Registry",id:"pushing-recipe-to-container-registry",level:2},{value:"Executing a recipe from local path",id:"executing-a-recipe-from-local-path",level:3},{value:"Executing a recipe from Container registry",id:"executing-a-recipe-from-container-registry",level:3},{value:"Checking updates for a recipe",id:"checking-updates-for-a-recipe",level:3}],s={toc:c},u="wrapper";function d(e){let{components:t,...a}=e;return(0,n.kt)(u,(0,r.Z)({},s,a,{components:t,mdxType:"MDXLayout"}),(0,n.kt)("h1",{id:"usage-wip"},"Usage ","[WIP]"),(0,n.kt)("h2",{id:"getting-started"},"Getting started"),(0,n.kt)("p",null,(0,n.kt)("inlineCode",{parentName:"p"},"jalapeno create my-recipe")),(0,n.kt)("ul",null,(0,n.kt)("li",{parentName:"ul"},"Add variable"),(0,n.kt)("li",{parentName:"ul"},"Use that variable in templates")),(0,n.kt)("p",null,(0,n.kt)("inlineCode",{parentName:"p"},"jalapeno execute my-recipe")),(0,n.kt)("h2",{id:"variables"},"Variables"),(0,n.kt)("h3",{id:"validation"},"Validation"),(0,n.kt)("h3",{id:"conditional-variables"},"Conditional variables"),(0,n.kt)("h2",{id:"templating"},"Templating"),(0,n.kt)("p",null,"Templates are done by using ",(0,n.kt)("a",{parentName:"p",href:"https://pkg.go.dev/text/template"},"Go templates"),"."),(0,n.kt)("p",null,"The following context is available on the templates:"),(0,n.kt)("ul",null,(0,n.kt)("li",{parentName:"ul"},(0,n.kt)("inlineCode",{parentName:"li"},"Recipe"),": Metadata object of the recipe",(0,n.kt)("ul",{parentName:"li"},(0,n.kt)("li",{parentName:"ul"},(0,n.kt)("inlineCode",{parentName:"li"},"Recipe.Name"),": The name of the recipe"),(0,n.kt)("li",{parentName:"ul"},(0,n.kt)("inlineCode",{parentName:"li"},"Recipe.Version"),": The current version of the recipe"))),(0,n.kt)("li",{parentName:"ul"},(0,n.kt)("inlineCode",{parentName:"li"},"ID"),": UUID which is generated after the first execution of the recipe. It will keep its value over upgrades. Can be used to generate unique pseudo-random values which stays the same over the upgrades, for example ",(0,n.kt)("inlineCode",{parentName:"li"},"my-resource-{{ sha1sum .ID | trunc 5 }}")),(0,n.kt)("li",{parentName:"ul"},(0,n.kt)("inlineCode",{parentName:"li"},"Variables"),": Object which contains the values of the variables defined for the recipe. Example: ",(0,n.kt)("inlineCode",{parentName:"li"},"{{ .Variables.FOO }}"))),(0,n.kt)("h2",{id:"pushing-recipe-to-container-registry"},"Pushing recipe to Container Registry"),(0,n.kt)("h3",{id:"executing-a-recipe-from-local-path"},"Executing a recipe from local path"),(0,n.kt)("h3",{id:"executing-a-recipe-from-container-registry"},"Executing a recipe from Container registry"),(0,n.kt)("h3",{id:"checking-updates-for-a-recipe"},"Checking updates for a recipe"),(0,n.kt)("p",null,(0,n.kt)("inlineCode",{parentName:"p"},"jalapeno check")))}d.isMDXComponent=!0}}]);