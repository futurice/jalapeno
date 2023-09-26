"use strict";(self.webpackChunkjalapeno=self.webpackChunkjalapeno||[]).push([[207],{3905:(t,e,a)=>{a.d(e,{Zo:()=>k,kt:()=>s});var n=a(7294);function l(t,e,a){return e in t?Object.defineProperty(t,e,{value:a,enumerable:!0,configurable:!0,writable:!0}):t[e]=a,t}function r(t,e){var a=Object.keys(t);if(Object.getOwnPropertySymbols){var n=Object.getOwnPropertySymbols(t);e&&(n=n.filter((function(e){return Object.getOwnPropertyDescriptor(t,e).enumerable}))),a.push.apply(a,n)}return a}function i(t){for(var e=1;e<arguments.length;e++){var a=null!=arguments[e]?arguments[e]:{};e%2?r(Object(a),!0).forEach((function(e){l(t,e,a[e])})):Object.getOwnPropertyDescriptors?Object.defineProperties(t,Object.getOwnPropertyDescriptors(a)):r(Object(a)).forEach((function(e){Object.defineProperty(t,e,Object.getOwnPropertyDescriptor(a,e))}))}return t}function p(t,e){if(null==t)return{};var a,n,l=function(t,e){if(null==t)return{};var a,n,l={},r=Object.keys(t);for(n=0;n<r.length;n++)a=r[n],e.indexOf(a)>=0||(l[a]=t[a]);return l}(t,e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(t);for(n=0;n<r.length;n++)a=r[n],e.indexOf(a)>=0||Object.prototype.propertyIsEnumerable.call(t,a)&&(l[a]=t[a])}return l}var d=n.createContext({}),o=function(t){var e=n.useContext(d),a=e;return t&&(a="function"==typeof t?t(e):i(i({},e),t)),a},k=function(t){var e=o(t.components);return n.createElement(d.Provider,{value:e},t.children)},m="mdxType",u={inlineCode:"code",wrapper:function(t){var e=t.children;return n.createElement(n.Fragment,{},e)}},N=n.forwardRef((function(t,e){var a=t.components,l=t.mdxType,r=t.originalType,d=t.parentName,k=p(t,["components","mdxType","originalType","parentName"]),m=o(a),N=l,s=m["".concat(d,".").concat(N)]||m[N]||u[N]||r;return a?n.createElement(s,i(i({ref:e},k),{},{components:a})):n.createElement(s,i({ref:e},k))}));function s(t,e){var a=arguments,l=e&&e.mdxType;if("string"==typeof t||l){var r=a.length,i=new Array(r);i[0]=N;var p={};for(var d in e)hasOwnProperty.call(e,d)&&(p[d]=e[d]);p.originalType=t,p[m]="string"==typeof t?t:l,i[1]=p;for(var o=2;o<r;o++)i[o]=a[o];return n.createElement.apply(null,i)}return n.createElement.apply(null,a)}N.displayName="MDXCreateElement"},7480:(t,e,a)=>{a.r(e),a.d(e,{assets:()=>d,contentTitle:()=>i,default:()=>u,frontMatter:()=>r,metadata:()=>p,toc:()=>o});var n=a(7462),l=(a(7294),a(3905));const r={sidebar_position:4,slug:"/api/",title:"API Reference"},i=void 0,p={unversionedId:"api",id:"api",title:"API Reference",description:"check",source:"@site/docs/api.md",sourceDirName:".",slug:"/api/",permalink:"/jalapeno/api/",draft:!1,editUrl:"https://github.com/futurice/jalapeno/tree/main/docs/docs/api.md",tags:[],version:"current",sidebarPosition:4,frontMatter:{sidebar_position:4,slug:"/api/",title:"API Reference"},sidebar:"tutorialSidebar",previous:{title:"Usage",permalink:"/jalapeno/usage/"},next:{title:"Contributing",permalink:"/jalapeno/contributing/"}},d={},o=[{value:"check",id:"check",level:2},{value:"Usage: <code>jalapeno check RECIPE_NAME</code>",id:"usage-jalapeno-check-recipe_name",level:4},{value:"Options",id:"check-options",level:4},{value:"create",id:"create",level:2},{value:"Usage: <code>jalapeno create RECIPE_NAME</code>",id:"usage-jalapeno-create-recipe_name",level:4},{value:"Examples",id:"examples",level:4},{value:"Options",id:"create-options",level:4},{value:"eject",id:"eject",level:2},{value:"Usage: <code>jalapeno eject</code>",id:"usage-jalapeno-eject",level:4},{value:"Options",id:"eject-options",level:4},{value:"execute",id:"execute",level:2},{value:"Usage: <code>jalapeno execute RECIPE_PATH</code>",id:"usage-jalapeno-execute-recipe_path",level:4},{value:"Options",id:"execute-options",level:4},{value:"pull",id:"pull",level:2},{value:"Usage: <code>jalapeno pull URL</code>",id:"usage-jalapeno-pull-url",level:4},{value:"Options",id:"pull-options",level:4},{value:"push",id:"push",level:2},{value:"Usage: <code>jalapeno push RECIPE_PATH TARGET_URL</code>",id:"usage-jalapeno-push-recipe_path-target_url",level:4},{value:"Options",id:"push-options",level:4},{value:"test",id:"test",level:2},{value:"Usage: <code>jalapeno test RECIPE_PATH</code>",id:"usage-jalapeno-test-recipe_path",level:4},{value:"Options",id:"test-options",level:4},{value:"upgrade",id:"upgrade",level:2},{value:"Usage: <code>jalapeno upgrade RECIPE_PATH</code>",id:"usage-jalapeno-upgrade-recipe_path",level:4},{value:"Options",id:"upgrade-options",level:4},{value:"validate",id:"validate",level:2},{value:"Usage: <code>jalapeno validate RECIPE_PATH</code>",id:"usage-jalapeno-validate-recipe_path",level:4},{value:"why",id:"why",level:2},{value:"Usage: <code>jalapeno why FILEPATH</code>",id:"usage-jalapeno-why-filepath",level:4},{value:"Options",id:"why-options",level:4}],k={toc:o},m="wrapper";function u(t){let{components:e,...a}=t;return(0,l.kt)(m,(0,n.Z)({},k,a,{components:e,mdxType:"MDXLayout"}),(0,l.kt)("h1",{id:"cli"},"CLI"),(0,l.kt)("h2",{id:"check"},"check"),(0,l.kt)("p",null,"TODO"),(0,l.kt)("h4",{id:"usage-jalapeno-check-recipe_name"},"Usage: ",(0,l.kt)("inlineCode",{parentName:"h4"},"jalapeno check RECIPE_NAME")),(0,l.kt)("h4",{id:"check-options"},"Options"),(0,l.kt)("table",null,(0,l.kt)("thead",{parentName:"table"},(0,l.kt)("tr",{parentName:"thead"},(0,l.kt)("th",{parentName:"tr",align:null},"Name"),(0,l.kt)("th",{parentName:"tr",align:null},"Type"),(0,l.kt)("th",{parentName:"tr",align:null},"Default"),(0,l.kt)("th",{parentName:"tr",align:null},"Description"))),(0,l.kt)("tbody",{parentName:"table"},(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"--ca-file")),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"string")),(0,l.kt)("td",{parentName:"tr",align:null}),(0,l.kt)("td",{parentName:"tr",align:null},"Server certificate authority file for the remote registry")),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"--dir"),", ",(0,l.kt)("inlineCode",{parentName:"td"},"-d")),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"string")),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},".")),(0,l.kt)("td",{parentName:"tr",align:null},"Sets working directory")),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"--insecure")),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"bool")),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"false")),(0,l.kt)("td",{parentName:"tr",align:null},"Allow connections to SSL registry without certs")),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"--password"),", ",(0,l.kt)("inlineCode",{parentName:"td"},"-p")),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"string")),(0,l.kt)("td",{parentName:"tr",align:null}),(0,l.kt)("td",{parentName:"tr",align:null},"Registry password or identity token")),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"--plain-http")),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"bool")),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"false")),(0,l.kt)("td",{parentName:"tr",align:null},"Allow insecure connections to registry without SSL check")),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"--registry-config")),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"stringArray")),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"[]")),(0,l.kt)("td",{parentName:"tr",align:null},"Path of the authentication file")),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"--username"),", ",(0,l.kt)("inlineCode",{parentName:"td"},"-u")),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"string")),(0,l.kt)("td",{parentName:"tr",align:null}),(0,l.kt)("td",{parentName:"tr",align:null},"Registry username")))),(0,l.kt)("hr",null),(0,l.kt)("h2",{id:"create"},"create"),(0,l.kt)("p",null,"TODO"),(0,l.kt)("pre",null,(0,l.kt)("code",{parentName:"pre"},"my-recipe\n\u251c\u2500\u2500 recipe.yml\n\u251c\u2500\u2500 templates\n\u2502   \u2514\u2500\u2500 README.md\n\u2514\u2500\u2500 tests\n    \u2514\u2500\u2500 defaults\n        \u251c\u2500\u2500 test.yml\n        \u2514\u2500\u2500 files\n            \u2514\u2500\u2500 README.md\n")),(0,l.kt)("h4",{id:"usage-jalapeno-create-recipe_name"},"Usage: ",(0,l.kt)("inlineCode",{parentName:"h4"},"jalapeno create RECIPE_NAME")),(0,l.kt)("h4",{id:"examples"},"Examples"),(0,l.kt)("pre",null,(0,l.kt)("code",{parentName:"pre"},"jalapeno create my-recipe\n")),(0,l.kt)("h4",{id:"create-options"},"Options"),(0,l.kt)("table",null,(0,l.kt)("thead",{parentName:"table"},(0,l.kt)("tr",{parentName:"thead"},(0,l.kt)("th",{parentName:"tr",align:null},"Name"),(0,l.kt)("th",{parentName:"tr",align:null},"Type"),(0,l.kt)("th",{parentName:"tr",align:null},"Default"),(0,l.kt)("th",{parentName:"tr",align:null},"Description"))),(0,l.kt)("tbody",{parentName:"table"},(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"--dir"),", ",(0,l.kt)("inlineCode",{parentName:"td"},"-d")),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"string")),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},".")),(0,l.kt)("td",{parentName:"tr",align:null},"Sets working directory")))),(0,l.kt)("hr",null),(0,l.kt)("h2",{id:"eject"},"eject"),(0,l.kt)("p",null,"Remove all the files and directories that are for Jalapeno internal use, and leave only the rendered project files."),(0,l.kt)("h4",{id:"usage-jalapeno-eject"},"Usage: ",(0,l.kt)("inlineCode",{parentName:"h4"},"jalapeno eject")),(0,l.kt)("h4",{id:"eject-options"},"Options"),(0,l.kt)("table",null,(0,l.kt)("thead",{parentName:"table"},(0,l.kt)("tr",{parentName:"thead"},(0,l.kt)("th",{parentName:"tr",align:null},"Name"),(0,l.kt)("th",{parentName:"tr",align:null},"Type"),(0,l.kt)("th",{parentName:"tr",align:null},"Default"),(0,l.kt)("th",{parentName:"tr",align:null},"Description"))),(0,l.kt)("tbody",{parentName:"table"},(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"--dir"),", ",(0,l.kt)("inlineCode",{parentName:"td"},"-d")),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"string")),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},".")),(0,l.kt)("td",{parentName:"tr",align:null},"Sets working directory")))),(0,l.kt)("hr",null),(0,l.kt)("h2",{id:"execute"},"execute"),(0,l.kt)("p",null,"TODO"),(0,l.kt)("h4",{id:"usage-jalapeno-execute-recipe_path"},"Usage: ",(0,l.kt)("inlineCode",{parentName:"h4"},"jalapeno execute RECIPE_PATH")),(0,l.kt)("h4",{id:"execute-options"},"Options"),(0,l.kt)("table",null,(0,l.kt)("thead",{parentName:"table"},(0,l.kt)("tr",{parentName:"thead"},(0,l.kt)("th",{parentName:"tr",align:null},"Name"),(0,l.kt)("th",{parentName:"tr",align:null},"Type"),(0,l.kt)("th",{parentName:"tr",align:null},"Default"),(0,l.kt)("th",{parentName:"tr",align:null},"Description"))),(0,l.kt)("tbody",{parentName:"table"},(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"--ca-file")),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"string")),(0,l.kt)("td",{parentName:"tr",align:null}),(0,l.kt)("td",{parentName:"tr",align:null},"Server certificate authority file for the remote registry")),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"--dir"),", ",(0,l.kt)("inlineCode",{parentName:"td"},"-d")),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"string")),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},".")),(0,l.kt)("td",{parentName:"tr",align:null},"Sets working directory")),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"--insecure")),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"bool")),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"false")),(0,l.kt)("td",{parentName:"tr",align:null},"Allow connections to SSL registry without certs")),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"--password"),", ",(0,l.kt)("inlineCode",{parentName:"td"},"-p")),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"string")),(0,l.kt)("td",{parentName:"tr",align:null}),(0,l.kt)("td",{parentName:"tr",align:null},"Registry password or identity token")),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"--plain-http")),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"bool")),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"false")),(0,l.kt)("td",{parentName:"tr",align:null},"Allow insecure connections to registry without SSL check")),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"--registry-config")),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"stringArray")),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"[]")),(0,l.kt)("td",{parentName:"tr",align:null},"Path of the authentication file")),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"--reuse-sauce-values"),", ",(0,l.kt)("inlineCode",{parentName:"td"},"-r")),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"bool")),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"false")),(0,l.kt)("td",{parentName:"tr",align:null},"By default each sauce has their own set of values even if the variable names are same in both recipes. Setting this to ",(0,l.kt)("inlineCode",{parentName:"td"},"true")," will reuse previous sauce values if the variable name match")),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"--set"),", ",(0,l.kt)("inlineCode",{parentName:"td"},"-s")),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"stringArray")),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"[]")),(0,l.kt)("td",{parentName:"tr",align:null},"Predefine values to be used in the templates. Example: ",(0,l.kt)("inlineCode",{parentName:"td"},'--set "MY_VAR=foo"'))),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"--username"),", ",(0,l.kt)("inlineCode",{parentName:"td"},"-u")),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"string")),(0,l.kt)("td",{parentName:"tr",align:null}),(0,l.kt)("td",{parentName:"tr",align:null},"Registry username")))),(0,l.kt)("hr",null),(0,l.kt)("h2",{id:"pull"},"pull"),(0,l.kt)("p",null,"TODO"),(0,l.kt)("h4",{id:"usage-jalapeno-pull-url"},"Usage: ",(0,l.kt)("inlineCode",{parentName:"h4"},"jalapeno pull URL")),(0,l.kt)("h4",{id:"pull-options"},"Options"),(0,l.kt)("table",null,(0,l.kt)("thead",{parentName:"table"},(0,l.kt)("tr",{parentName:"thead"},(0,l.kt)("th",{parentName:"tr",align:null},"Name"),(0,l.kt)("th",{parentName:"tr",align:null},"Type"),(0,l.kt)("th",{parentName:"tr",align:null},"Default"),(0,l.kt)("th",{parentName:"tr",align:null},"Description"))),(0,l.kt)("tbody",{parentName:"table"},(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"--ca-file")),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"string")),(0,l.kt)("td",{parentName:"tr",align:null}),(0,l.kt)("td",{parentName:"tr",align:null},"Server certificate authority file for the remote registry")),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"--dir"),", ",(0,l.kt)("inlineCode",{parentName:"td"},"-d")),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"string")),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},".")),(0,l.kt)("td",{parentName:"tr",align:null},"Sets working directory")),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"--insecure")),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"bool")),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"false")),(0,l.kt)("td",{parentName:"tr",align:null},"Allow connections to SSL registry without certs")),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"--password"),", ",(0,l.kt)("inlineCode",{parentName:"td"},"-p")),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"string")),(0,l.kt)("td",{parentName:"tr",align:null}),(0,l.kt)("td",{parentName:"tr",align:null},"Registry password or identity token")),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"--plain-http")),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"bool")),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"false")),(0,l.kt)("td",{parentName:"tr",align:null},"Allow insecure connections to registry without SSL check")),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"--registry-config")),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"stringArray")),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"[]")),(0,l.kt)("td",{parentName:"tr",align:null},"Path of the authentication file")),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"--username"),", ",(0,l.kt)("inlineCode",{parentName:"td"},"-u")),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"string")),(0,l.kt)("td",{parentName:"tr",align:null}),(0,l.kt)("td",{parentName:"tr",align:null},"Registry username")))),(0,l.kt)("hr",null),(0,l.kt)("h2",{id:"push"},"push"),(0,l.kt)("p",null,"TODO"),(0,l.kt)("h4",{id:"usage-jalapeno-push-recipe_path-target_url"},"Usage: ",(0,l.kt)("inlineCode",{parentName:"h4"},"jalapeno push RECIPE_PATH TARGET_URL")),(0,l.kt)("h4",{id:"push-options"},"Options"),(0,l.kt)("table",null,(0,l.kt)("thead",{parentName:"table"},(0,l.kt)("tr",{parentName:"thead"},(0,l.kt)("th",{parentName:"tr",align:null},"Name"),(0,l.kt)("th",{parentName:"tr",align:null},"Type"),(0,l.kt)("th",{parentName:"tr",align:null},"Default"),(0,l.kt)("th",{parentName:"tr",align:null},"Description"))),(0,l.kt)("tbody",{parentName:"table"},(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"--ca-file")),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"string")),(0,l.kt)("td",{parentName:"tr",align:null}),(0,l.kt)("td",{parentName:"tr",align:null},"Server certificate authority file for the remote registry")),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"--insecure")),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"bool")),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"false")),(0,l.kt)("td",{parentName:"tr",align:null},"Allow connections to SSL registry without certs")),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"--password"),", ",(0,l.kt)("inlineCode",{parentName:"td"},"-p")),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"string")),(0,l.kt)("td",{parentName:"tr",align:null}),(0,l.kt)("td",{parentName:"tr",align:null},"Registry password or identity token")),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"--plain-http")),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"bool")),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"false")),(0,l.kt)("td",{parentName:"tr",align:null},"Allow insecure connections to registry without SSL check")),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"--registry-config")),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"stringArray")),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"[]")),(0,l.kt)("td",{parentName:"tr",align:null},"Path of the authentication file")),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"--username"),", ",(0,l.kt)("inlineCode",{parentName:"td"},"-u")),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"string")),(0,l.kt)("td",{parentName:"tr",align:null}),(0,l.kt)("td",{parentName:"tr",align:null},"Registry username")))),(0,l.kt)("hr",null),(0,l.kt)("h2",{id:"test"},"test"),(0,l.kt)("p",null,"TODO"),(0,l.kt)("h4",{id:"usage-jalapeno-test-recipe_path"},"Usage: ",(0,l.kt)("inlineCode",{parentName:"h4"},"jalapeno test RECIPE_PATH")),(0,l.kt)("h4",{id:"test-options"},"Options"),(0,l.kt)("table",null,(0,l.kt)("thead",{parentName:"table"},(0,l.kt)("tr",{parentName:"thead"},(0,l.kt)("th",{parentName:"tr",align:null},"Name"),(0,l.kt)("th",{parentName:"tr",align:null},"Type"),(0,l.kt)("th",{parentName:"tr",align:null},"Default"),(0,l.kt)("th",{parentName:"tr",align:null},"Description"))),(0,l.kt)("tbody",{parentName:"table"},(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"--create"),", ",(0,l.kt)("inlineCode",{parentName:"td"},"-c")),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"bool")),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"false")),(0,l.kt)("td",{parentName:"tr",align:null},"Create a new test case")),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"--update-snapshots"),", ",(0,l.kt)("inlineCode",{parentName:"td"},"-u")),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"bool")),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"false")),(0,l.kt)("td",{parentName:"tr",align:null},"Update test file snapshots")))),(0,l.kt)("hr",null),(0,l.kt)("h2",{id:"upgrade"},"upgrade"),(0,l.kt)("p",null,"Upgrade a recipe in a project with a newer version."),(0,l.kt)("h4",{id:"usage-jalapeno-upgrade-recipe_path"},"Usage: ",(0,l.kt)("inlineCode",{parentName:"h4"},"jalapeno upgrade RECIPE_PATH")),(0,l.kt)("h4",{id:"upgrade-options"},"Options"),(0,l.kt)("table",null,(0,l.kt)("thead",{parentName:"table"},(0,l.kt)("tr",{parentName:"thead"},(0,l.kt)("th",{parentName:"tr",align:null},"Name"),(0,l.kt)("th",{parentName:"tr",align:null},"Type"),(0,l.kt)("th",{parentName:"tr",align:null},"Default"),(0,l.kt)("th",{parentName:"tr",align:null},"Description"))),(0,l.kt)("tbody",{parentName:"table"},(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"--dir"),", ",(0,l.kt)("inlineCode",{parentName:"td"},"-d")),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"string")),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},".")),(0,l.kt)("td",{parentName:"tr",align:null},"Sets working directory")),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"--reuse-sauce-values"),", ",(0,l.kt)("inlineCode",{parentName:"td"},"-r")),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"bool")),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"false")),(0,l.kt)("td",{parentName:"tr",align:null},"By default each sauce has their own set of values even if the variable names are same in both recipes. Setting this to ",(0,l.kt)("inlineCode",{parentName:"td"},"true")," will reuse previous sauce values if the variable name match")),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"--set"),", ",(0,l.kt)("inlineCode",{parentName:"td"},"-s")),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"stringArray")),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"[]")),(0,l.kt)("td",{parentName:"tr",align:null},"Predefine values to be used in the templates. Example: ",(0,l.kt)("inlineCode",{parentName:"td"},'--set "MY_VAR=foo"'))))),(0,l.kt)("hr",null),(0,l.kt)("h2",{id:"validate"},"validate"),(0,l.kt)("p",null,"TODO"),(0,l.kt)("h4",{id:"usage-jalapeno-validate-recipe_path"},"Usage: ",(0,l.kt)("inlineCode",{parentName:"h4"},"jalapeno validate RECIPE_PATH")),(0,l.kt)("hr",null),(0,l.kt)("h2",{id:"why"},"why"),(0,l.kt)("p",null,"Explains where a file comes from in the project, e.g. is the file create by a recipe or user"),(0,l.kt)("h4",{id:"usage-jalapeno-why-filepath"},"Usage: ",(0,l.kt)("inlineCode",{parentName:"h4"},"jalapeno why FILEPATH")),(0,l.kt)("h4",{id:"why-options"},"Options"),(0,l.kt)("table",null,(0,l.kt)("thead",{parentName:"table"},(0,l.kt)("tr",{parentName:"thead"},(0,l.kt)("th",{parentName:"tr",align:null},"Name"),(0,l.kt)("th",{parentName:"tr",align:null},"Type"),(0,l.kt)("th",{parentName:"tr",align:null},"Default"),(0,l.kt)("th",{parentName:"tr",align:null},"Description"))),(0,l.kt)("tbody",{parentName:"table"},(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"--dir"),", ",(0,l.kt)("inlineCode",{parentName:"td"},"-d")),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"string")),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},".")),(0,l.kt)("td",{parentName:"tr",align:null},"Sets working directory")))),(0,l.kt)("hr",null))}u.isMDXComponent=!0}}]);