"use strict";(self.webpackChunkjalapeno=self.webpackChunkjalapeno||[]).push([[827],{351:(e,i,n)=>{n.r(i),n.d(i,{assets:()=>c,contentTitle:()=>s,default:()=>h,frontMatter:()=>r,metadata:()=>l,toc:()=>o});var t=n(5893),a=n(1151);const r={sidebar_position:3,slug:"/usage/",title:"Usage"},s="Usage [WIP]",l={id:"usage",title:"Usage",description:"Getting started",source:"@site/docs/usage.md",sourceDirName:".",slug:"/usage/",permalink:"/jalapeno/usage/",draft:!1,unlisted:!1,editUrl:"https://github.com/futurice/jalapeno/tree/main/docs/site/docs/usage.md",tags:[],version:"current",sidebarPosition:3,frontMatter:{sidebar_position:3,slug:"/usage/",title:"Usage"},sidebar:"tutorialSidebar",previous:{title:"Installation",permalink:"/jalapeno/installation/"},next:{title:"API Reference",permalink:"/jalapeno/api/"}},c={},o=[{value:"Getting started",id:"getting-started",level:2},{value:"Templating",id:"templating",level:2},{value:"Variables",id:"variables",level:2},{value:"Validation",id:"validation",level:3},{value:"Conditional variables",id:"conditional-variables",level:3},{value:"Pushing recipe to Container Registry",id:"pushing-recipe-to-container-registry",level:2},{value:"Executing a recipe from local path",id:"executing-a-recipe-from-local-path",level:3},{value:"Executing a recipe from Container registry",id:"executing-a-recipe-from-container-registry",level:3},{value:"Checking updates for a recipe",id:"checking-updates-for-a-recipe",level:3}];function d(e){const i={a:"a",code:"code",em:"em",h1:"h1",h2:"h2",h3:"h3",li:"li",p:"p",pre:"pre",ul:"ul",...(0,a.a)(),...e.components};return(0,t.jsxs)(t.Fragment,{children:[(0,t.jsx)(i.h1,{id:"usage-wip",children:"Usage [WIP]"}),"\n",(0,t.jsx)(i.h2,{id:"getting-started",children:"Getting started"}),"\n",(0,t.jsxs)(i.p,{children:["To get started with Jalapeno, you need to install the CLI tool. You can find the installation instructions ",(0,t.jsx)(i.a,{href:"/installation",children:"here"}),"."]}),"\n",(0,t.jsxs)(i.p,{children:["Then you can bootstrap a new ",(0,t.jsx)(i.em,{children:"recipe"})," (aka template) by running:"]}),"\n",(0,t.jsx)(i.pre,{children:(0,t.jsx)(i.code,{className:"language-bash",children:"jalapeno create my-recipe\n"})}),"\n",(0,t.jsxs)(i.p,{children:["After this you should have a new folder called ",(0,t.jsx)(i.code,{children:"my-recipe"})," with the following structure:"]}),"\n",(0,t.jsx)(i.pre,{children:(0,t.jsx)(i.code,{children:"my-recipe\n\u251c\u2500\u2500 recipe.yml\n\u251c\u2500\u2500 templates\n\u2502   \u2514\u2500\u2500 README.md\n\u2514\u2500\u2500 tests\n    \u2514\u2500\u2500 defaults\n        \u251c\u2500\u2500 test.yml\n        \u2514\u2500\u2500 files\n            \u2514\u2500\u2500 README.md\n"})}),"\n",(0,t.jsxs)(i.p,{children:["The ",(0,t.jsx)(i.code,{children:"templates"})," directory contains the templates which will be rendered to the project directory. You can add and edit files there or you can already execute this ",(0,t.jsx)(i.em,{children:"recipe"})," (render the templates) to your project directory by running:"]}),"\n",(0,t.jsx)(i.pre,{children:(0,t.jsx)(i.code,{className:"language-bash",children:"mkdir my-project\njalapeno execute my-recipe -d my-project\n"})}),"\n",(0,t.jsx)(i.p,{children:"After this, the project directory should have the following files:"}),"\n",(0,t.jsx)(i.pre,{children:(0,t.jsx)(i.code,{children:"my-project\n\u251c\u2500\u2500 .jalapeno\n\u2502   \u2514\u2500\u2500 sauces.yml\n\u2514\u2500\u2500 README.md\n"})}),"\n",(0,t.jsxs)(i.p,{children:["The ",(0,t.jsx)(i.code,{children:".jalapeno"})," directory contains files which Jalapeno uses internally. For example ",(0,t.jsx)(i.code,{children:"sauces.yml"})," is Jalapeno metadata file which contains information about the ",(0,t.jsx)(i.em,{children:"sauces"})," (aka executed recipes). This file is used to check for updates for the recipes later."]}),"\n",(0,t.jsx)(i.p,{children:"The rest of the files are rendered from the templates. You can edit the templates and execute the recipe again to update the files."}),"\n",(0,t.jsx)(i.h2,{id:"templating",children:"Templating"}),"\n",(0,t.jsxs)(i.p,{children:["Templates are done by using ",(0,t.jsx)(i.a,{href:"https://pkg.go.dev/text/template",children:"Go templates"}),"."]}),"\n",(0,t.jsx)(i.p,{children:"The following context is available on the templates:"}),"\n",(0,t.jsxs)(i.ul,{children:["\n",(0,t.jsxs)(i.li,{children:[(0,t.jsx)(i.code,{children:"Recipe"}),": Metadata object of the recipe","\n",(0,t.jsxs)(i.ul,{children:["\n",(0,t.jsxs)(i.li,{children:[(0,t.jsx)(i.code,{children:"Recipe.Name"}),": The name of the recipe"]}),"\n",(0,t.jsxs)(i.li,{children:[(0,t.jsx)(i.code,{children:"Recipe.Version"}),": The current version of the recipe"]}),"\n"]}),"\n"]}),"\n",(0,t.jsxs)(i.li,{children:[(0,t.jsx)(i.code,{children:"ID"}),": UUID which is generated after the first execution of the recipe. It will keep its value over upgrades. Can be used to generate unique pseudo-random values which stays the same over the upgrades, for example ",(0,t.jsx)(i.code,{children:"my-resource-{{ sha1sum .ID | trunc 5 }}"})]}),"\n",(0,t.jsxs)(i.li,{children:[(0,t.jsx)(i.code,{children:"Variables"}),": Object which contains the values of the variables defined for the recipe. Example: ",(0,t.jsx)(i.code,{children:"{{ .Variables.FOO }}"})]}),"\n"]}),"\n",(0,t.jsx)(i.h2,{id:"variables",children:"Variables"}),"\n",(0,t.jsx)(i.h3,{id:"validation",children:"Validation"}),"\n",(0,t.jsx)(i.h3,{id:"conditional-variables",children:"Conditional variables"}),"\n",(0,t.jsx)(i.h2,{id:"pushing-recipe-to-container-registry",children:"Pushing recipe to Container Registry"}),"\n",(0,t.jsx)(i.h3,{id:"executing-a-recipe-from-local-path",children:"Executing a recipe from local path"}),"\n",(0,t.jsx)(i.h3,{id:"executing-a-recipe-from-container-registry",children:"Executing a recipe from Container registry"}),"\n",(0,t.jsx)(i.h3,{id:"checking-updates-for-a-recipe",children:"Checking updates for a recipe"}),"\n",(0,t.jsx)(i.p,{children:(0,t.jsx)(i.code,{children:"jalapeno check"})})]})}function h(e={}){const{wrapper:i}={...(0,a.a)(),...e.components};return i?(0,t.jsx)(i,{...e,children:(0,t.jsx)(d,{...e})}):d(e)}}}]);