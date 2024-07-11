"use strict";(self.webpackChunkjalapeno=self.webpackChunkjalapeno||[]).push([[802],{177:(e,a,n)=>{n.r(a),n.d(a,{assets:()=>A,contentTitle:()=>V,default:()=>q,frontMatter:()=>S,metadata:()=>T,toc:()=>D});var l=n(4848),t=n(8453),r=n(6540),i=n(4164),s=n(3104),o=n(6347),c=n(205),u=n(7485),d=n(1682),h=n(679);function p(e){return r.Children.toArray(e).filter((e=>"\n"!==e)).map((e=>{if(!e||(0,r.isValidElement)(e)&&function(e){const{props:a}=e;return!!a&&"object"==typeof a&&"value"in a}(e))return e;throw new Error(`Docusaurus error: Bad <Tabs> child <${"string"==typeof e.type?e.type:e.type.name}>: all children of the <Tabs> component should be <TabItem>, and every <TabItem> should have a unique "value" prop.`)}))?.filter(Boolean)??[]}function m(e){const{values:a,children:n}=e;return(0,r.useMemo)((()=>{const e=a??function(e){return p(e).map((e=>{let{props:{value:a,label:n,attributes:l,default:t}}=e;return{value:a,label:n,attributes:l,default:t}}))}(n);return function(e){const a=(0,d.X)(e,((e,a)=>e.value===a.value));if(a.length>0)throw new Error(`Docusaurus error: Duplicate values "${a.map((e=>e.value)).join(", ")}" found in <Tabs>. Every value needs to be unique.`)}(e),e}),[a,n])}function b(e){let{value:a,tabValues:n}=e;return n.some((e=>e.value===a))}function g(e){let{queryString:a=!1,groupId:n}=e;const l=(0,o.W6)(),t=function(e){let{queryString:a=!1,groupId:n}=e;if("string"==typeof a)return a;if(!1===a)return null;if(!0===a&&!n)throw new Error('Docusaurus error: The <Tabs> component groupId prop is required if queryString=true, because this value is used as the search param name. You can also provide an explicit value such as queryString="my-search-param".');return n??null}({queryString:a,groupId:n});return[(0,u.aZ)(t),(0,r.useCallback)((e=>{if(!t)return;const a=new URLSearchParams(l.location.search);a.set(t,e),l.replace({...l.location,search:a.toString()})}),[t,l])]}function f(e){const{defaultValue:a,queryString:n=!1,groupId:l}=e,t=m(e),[i,s]=(0,r.useState)((()=>function(e){let{defaultValue:a,tabValues:n}=e;if(0===n.length)throw new Error("Docusaurus error: the <Tabs> component requires at least one <TabItem> children component");if(a){if(!b({value:a,tabValues:n}))throw new Error(`Docusaurus error: The <Tabs> has a defaultValue "${a}" but none of its children has the corresponding value. Available values are: ${n.map((e=>e.value)).join(", ")}. If you intend to show no default tab, use defaultValue={null} instead.`);return a}const l=n.find((e=>e.default))??n[0];if(!l)throw new Error("Unexpected error: 0 tabValues");return l.value}({defaultValue:a,tabValues:t}))),[o,u]=g({queryString:n,groupId:l}),[d,p]=function(e){let{groupId:a}=e;const n=function(e){return e?`docusaurus.tab.${e}`:null}(a),[l,t]=(0,h.Dv)(n);return[l,(0,r.useCallback)((e=>{n&&t.set(e)}),[n,t])]}({groupId:l}),f=(()=>{const e=o??d;return b({value:e,tabValues:t})?e:null})();(0,c.A)((()=>{f&&s(f)}),[f]);return{selectedValue:i,selectValue:(0,r.useCallback)((e=>{if(!b({value:e,tabValues:t}))throw new Error(`Can't select invalid tab value=${e}`);s(e),u(e),p(e)}),[u,p,t]),tabValues:t}}var j=n(2303);const x={tabList:"tabList__CuJ",tabItem:"tabItem_LNqP"};function v(e){let{className:a,block:n,selectedValue:t,selectValue:r,tabValues:o}=e;const c=[],{blockElementScrollPositionUntilNextRender:u}=(0,s.a_)(),d=e=>{const a=e.currentTarget,n=c.indexOf(a),l=o[n].value;l!==t&&(u(a),r(l))},h=e=>{let a=null;switch(e.key){case"Enter":d(e);break;case"ArrowRight":{const n=c.indexOf(e.currentTarget)+1;a=c[n]??c[0];break}case"ArrowLeft":{const n=c.indexOf(e.currentTarget)-1;a=c[n]??c[c.length-1];break}}a?.focus()};return(0,l.jsx)("ul",{role:"tablist","aria-orientation":"horizontal",className:(0,i.A)("tabs",{"tabs--block":n},a),children:o.map((e=>{let{value:a,label:n,attributes:r}=e;return(0,l.jsx)("li",{role:"tab",tabIndex:t===a?0:-1,"aria-selected":t===a,ref:e=>c.push(e),onKeyDown:h,onClick:d,...r,className:(0,i.A)("tabs__item",x.tabItem,r?.className,{"tabs__item--active":t===a}),children:n??a},a)}))})}function w(e){let{lazy:a,children:n,selectedValue:t}=e;const i=(Array.isArray(n)?n:[n]).filter(Boolean);if(a){const e=i.find((e=>e.props.value===t));return e?(0,r.cloneElement)(e,{className:"margin-top--md"}):null}return(0,l.jsx)("div",{className:"margin-top--md",children:i.map(((e,a)=>(0,r.cloneElement)(e,{key:a,hidden:e.props.value!==t})))})}function k(e){const a=f(e);return(0,l.jsxs)("div",{className:(0,i.A)("tabs-container",x.tabList),children:[(0,l.jsx)(v,{...a,...e}),(0,l.jsx)(w,{...a,...e})]})}function y(e){const a=(0,j.A)();return(0,l.jsx)(k,{...e,children:p(e.children)},String(a))}const I={tabItem:"tabItem_Ymn6"};function N(e){let{children:a,hidden:n,className:t}=e;return(0,l.jsx)("div",{role:"tabpanel",className:(0,i.A)(I.tabItem,t),hidden:n,children:a})}const S={sidebar_position:2,slug:"/installation/",title:"Installation"},V="Installation",T={id:"installation",title:"Installation",description:"Using a package manager",source:"@site/docs/installation.mdx",sourceDirName:".",slug:"/installation/",permalink:"/jalapeno/installation/",draft:!1,unlisted:!1,editUrl:"https://github.com/futurice/jalapeno/tree/main/docs/site/docs/installation.mdx",tags:[],version:"current",sidebarPosition:2,frontMatter:{sidebar_position:2,slug:"/installation/",title:"Installation"},sidebar:"tutorialSidebar",previous:{title:"Home",permalink:"/jalapeno/"},next:{title:"Usage",permalink:"/jalapeno/usage/"}},A={},D=[{value:"Using a package manager",id:"using-a-package-manager",level:2},{value:"Updating with package manager",id:"updating-with-package-manager",level:2},{value:"From the Binary Releases",id:"from-the-binary-releases",level:2},{value:"Use via Docker",id:"use-via-docker",level:2},{value:"Build From Source",id:"build-from-source",level:2}];function E(e){const a={a:"a",code:"code",h1:"h1",h2:"h2",li:"li",ol:"ol",p:"p",pre:"pre",...(0,t.R)(),...e.components};return(0,l.jsxs)(l.Fragment,{children:[(0,l.jsx)(a.h1,{id:"installation",children:"Installation"}),"\n",(0,l.jsx)(a.h2,{id:"using-a-package-manager",children:"Using a package manager"}),"\n",(0,l.jsxs)(y,{groupId:"code-examples",children:[(0,l.jsx)(N,{value:"homebrew",label:"Homebrew",children:(0,l.jsx)(a.pre,{children:(0,l.jsx)(a.code,{className:"language-bash",metastring:'tab={"label":"Homebrew"}',children:"brew tap futurice/jalapeno\nbrew install jalapeno\n"})})}),(0,l.jsx)(N,{value:"winget",label:"Winget",children:(0,l.jsx)(a.pre,{children:(0,l.jsx)(a.code,{className:"language-powershell",metastring:'tab={"label":"Winget"}',children:"winget install jalapeno\n"})})})]}),"\n",(0,l.jsx)(a.h2,{id:"updating-with-package-manager",children:"Updating with package manager"}),"\n",(0,l.jsxs)(y,{groupId:"code-examples",children:[(0,l.jsx)(N,{value:"homebrew",label:"Homebrew",children:(0,l.jsx)(a.pre,{children:(0,l.jsx)(a.code,{className:"language-bash",metastring:'tab={"label":"Homebrew"}',children:"brew update\nbrew upgrade jalapeno\n"})})}),(0,l.jsx)(N,{value:"winget",label:"Winget",children:(0,l.jsx)(a.pre,{children:(0,l.jsx)(a.code,{className:"language-powershell",metastring:'tab={"label":"Winget"}',children:"winget upgrade jalapeno\n"})})})]}),"\n",(0,l.jsx)(a.h2,{id:"from-the-binary-releases",children:"From the Binary Releases"}),"\n",(0,l.jsxs)(a.p,{children:["Cross-platform binaries are provided with each release of Jalapeno. These can manually be\ndownloaded and installed from ",(0,l.jsx)(a.a,{href:"https://github.com/futurice/jalapeno/releases/",children:"GitHub releases"}),"."]}),"\n",(0,l.jsx)(a.p,{children:"In short the process is:"}),"\n",(0,l.jsxs)(a.ol,{children:["\n",(0,l.jsxs)(a.li,{children:["Download ",(0,l.jsx)(a.a,{href:"https://github.com/futurice/jalapeno/releases/latest",children:"the latest version of Jalapeno"}),"\nfor your platform"]}),"\n",(0,l.jsx)(a.li,{children:"Extract the archive"}),"\n",(0,l.jsx)(a.li,{children:"Make the binary executable"}),"\n",(0,l.jsx)(a.li,{children:"Rename and move the binary to proper location"}),"\n"]}),"\n",(0,l.jsx)(a.p,{children:"For example on MacOS (running on Apple Silicon) this can be done with:"}),"\n",(0,l.jsx)(a.pre,{children:(0,l.jsx)(a.code,{className:"language-bash",children:"curl -L https://github.com/futurice/jalapeno/releases/latest/download/jalapeno-darwin-arm64.tar.gz -o jalapeno.tar.gz\ntar -xvf jalapeno.tar.gz\nchmod +x jalapeno\nmv jalapeno /usr/local/bin/jalapeno\n"})}),"\n",(0,l.jsx)(a.h2,{id:"use-via-docker",children:"Use via Docker"}),"\n",(0,l.jsxs)(a.p,{children:["It is possible to use Jalapeno via Docker without installing it locally. Jalapeno image is available\nfrom the GitHub Container Registry, and thus the following command on a *nix system is equivalent\nto running ",(0,l.jsx)(a.code,{children:"jalapeno"})," locally:"]}),"\n",(0,l.jsxs)(y,{groupId:"code-examples",children:[(0,l.jsx)(N,{value:"macos-and linux",label:"MacOS and Linux",children:(0,l.jsx)(a.pre,{children:(0,l.jsx)(a.code,{className:"language-bash",metastring:'tab={"label":"MacOS and Linux"}',children:"docker run -it --rm -v $(pwd):/workdir ghcr.io/futurice/jalapeno:latest\n"})})}),(0,l.jsx)(N,{value:"windows-command line",label:"Windows Command Line",children:(0,l.jsx)(a.pre,{children:(0,l.jsx)(a.code,{className:"language-batch",metastring:'tab={"label":"Windows Command Line"}',children:"docker run -it --rm -v %cd%:/workdir ghcr.io/futurice/jalapeno:latest\n"})})}),(0,l.jsx)(N,{value:"powershell",label:"PowerShell",children:(0,l.jsx)(a.pre,{children:(0,l.jsx)(a.code,{className:"language-powershell",metastring:'tab={"label":"PowerShell"}',children:"docker run -it --rm -v ${PWD}:/workdir ghcr.io/futurice/jalapeno:latest\n"})})})]}),"\n",(0,l.jsx)(a.h2,{id:"build-from-source",children:"Build From Source"}),"\n",(0,l.jsxs)(a.p,{children:["First, make sure you have ",(0,l.jsx)(a.a,{href:"https://go.dev/doc/install",children:"Go"})," and\n",(0,l.jsx)(a.a,{href:"https://taskfile.dev/installation",children:"Task"})," installed."]}),"\n",(0,l.jsx)(a.p,{children:"Then you can compile the binary with following commands:"}),"\n",(0,l.jsx)(a.pre,{children:(0,l.jsx)(a.code,{className:"language-bash",children:"git clone https://github.com/futurice/jalapeno.git\ncd jalapeno\ntask build\n"})}),"\n",(0,l.jsxs)(a.p,{children:["After this the binary is available on path ",(0,l.jsx)(a.code,{children:"./bin/jalapeno"}),"."]})]})}function q(e={}){const{wrapper:a}={...(0,t.R)(),...e.components};return a?(0,l.jsx)(a,{...e,children:(0,l.jsx)(E,{...e})}):E(e)}}}]);