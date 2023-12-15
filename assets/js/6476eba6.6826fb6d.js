"use strict";(self.webpackChunkjalapeno=self.webpackChunkjalapeno||[]).push([[827],{351:(e,i,n)=>{n.r(i),n.d(i,{assets:()=>l,contentTitle:()=>s,default:()=>d,frontMatter:()=>t,metadata:()=>c,toc:()=>o});var r=n(5893),a=n(1151);const t={sidebar_position:3,slug:"/usage/",title:"Usage"},s="Usage",c={id:"usage",title:"Usage",description:"Getting started",source:"@site/docs/usage.md",sourceDirName:".",slug:"/usage/",permalink:"/jalapeno/usage/",draft:!1,unlisted:!1,editUrl:"https://github.com/futurice/jalapeno/tree/main/docs/site/docs/usage.md",tags:[],version:"current",sidebarPosition:3,frontMatter:{sidebar_position:3,slug:"/usage/",title:"Usage"},sidebar:"tutorialSidebar",previous:{title:"Installation",permalink:"/jalapeno/installation/"},next:{title:"API Reference",permalink:"/jalapeno/api/"}},l={},o=[{value:"Getting started",id:"getting-started",level:2},{value:"Templating",id:"templating",level:2},{value:"Template only specific files",id:"template-only-specific-files",level:3},{value:"Variables",id:"variables",level:2},{value:"Variable types",id:"variable-types",level:3},{value:"Validation",id:"validation",level:3},{value:"Publishing recipes",id:"publishing-recipes",level:2},{value:"Pushing a recipe to Container registry",id:"pushing-a-recipe-to-container-registry",level:3},{value:"Executing a recipe from Container registry",id:"executing-a-recipe-from-container-registry",level:3},{value:"Checking updates for a recipe",id:"checking-updates-for-a-recipe",level:3}];function h(e){const i={a:"a",admonition:"admonition",code:"code",em:"em",h1:"h1",h2:"h2",h3:"h3",li:"li",p:"p",pre:"pre",ul:"ul",...(0,a.a)(),...e.components};return(0,r.jsxs)(r.Fragment,{children:[(0,r.jsx)(i.h1,{id:"usage",children:"Usage"}),"\n",(0,r.jsx)(i.h2,{id:"getting-started",children:"Getting started"}),"\n",(0,r.jsxs)(i.p,{children:["To get started with Jalapeno, you need to install the CLI tool. You can find the installation instructions ",(0,r.jsx)(i.a,{href:"/installation",children:"here"}),"."]}),"\n",(0,r.jsxs)(i.p,{children:["Then you can bootstrap a new ",(0,r.jsx)(i.em,{children:"recipe"})," (aka template) by running:"]}),"\n",(0,r.jsx)(i.pre,{children:(0,r.jsx)(i.code,{className:"language-bash",children:"jalapeno create my-recipe\n"})}),"\n",(0,r.jsxs)(i.p,{children:["After this you should have a new folder called ",(0,r.jsx)(i.code,{children:"my-recipe"})," with the following structure:"]}),"\n",(0,r.jsx)(i.pre,{children:(0,r.jsx)(i.code,{children:"my-recipe\n\u251c\u2500\u2500 recipe.yml\n\u251c\u2500\u2500 templates\n\u2502   \u2514\u2500\u2500 README.md\n\u2514\u2500\u2500 tests\n    \u2514\u2500\u2500 defaults\n        \u251c\u2500\u2500 test.yml\n        \u2514\u2500\u2500 files\n            \u2514\u2500\u2500 README.md\n"})}),"\n",(0,r.jsxs)(i.p,{children:["The ",(0,r.jsx)(i.code,{children:"templates"})," directory contains the templates which will be rendered to the project directory. You can add and edit files there or you can already ",(0,r.jsx)(i.em,{children:"execute the recipe"})," (render the templates) to your project directory by running:"]}),"\n",(0,r.jsx)(i.pre,{children:(0,r.jsx)(i.code,{className:"language-bash",children:"mkdir my-project\njalapeno execute my-recipe -d my-project\n"})}),"\n",(0,r.jsxs)(i.admonition,{type:"tip",children:[(0,r.jsx)(i.p,{children:"You can also execute any of the examples from the Jalapeno repository, for example:"}),(0,r.jsx)(i.pre,{children:(0,r.jsx)(i.code,{className:"language-bash",children:"git clone git@github.com:futurice/jalapeno.git\njalapeno execute ./jalapeno/examples/variable-types -d my-project\n"})}),(0,r.jsx)(i.p,{children:"Or execute them directly from GitHub Container Registry:"}),(0,r.jsx)(i.pre,{children:(0,r.jsx)(i.code,{className:"language-bash",children:"jalapeno execute oci://ghcr.io/futurice/jalapeno/examples/variable-types:v0.0.0 -d my-project\n"})})]}),"\n",(0,r.jsx)(i.p,{children:"After this, the project directory should have the following files:"}),"\n",(0,r.jsx)(i.pre,{children:(0,r.jsx)(i.code,{children:"my-project\n\u251c\u2500\u2500 .jalapeno\n\u2502   \u2514\u2500\u2500 sauces.yml\n\u2514\u2500\u2500 README.md\n"})}),"\n",(0,r.jsxs)(i.p,{children:["The ",(0,r.jsx)(i.code,{children:".jalapeno"})," directory contains files which Jalapeno uses internally. For example ",(0,r.jsx)(i.code,{children:"sauces.yml"})," is Jalapeno metadata file which contains information about the ",(0,r.jsx)(i.em,{children:"sauces"})," (aka executed recipes). This file is used to check for updates for the recipes later."]}),"\n",(0,r.jsx)(i.p,{children:"The rest of the files are rendered from the templates. You can edit the templates and execute the recipe again to update the files."}),"\n",(0,r.jsx)(i.h2,{id:"templating",children:"Templating"}),"\n",(0,r.jsxs)(i.p,{children:["Templates are done by using ",(0,r.jsx)(i.a,{href:"https://pkg.go.dev/text/template",children:"Go templates"}),"."]}),"\n",(0,r.jsx)(i.p,{children:"The following context is available on the templates:"}),"\n",(0,r.jsxs)(i.ul,{children:["\n",(0,r.jsxs)(i.li,{children:[(0,r.jsx)(i.code,{children:".Recipe"}),": Metadata object of the recipe","\n",(0,r.jsxs)(i.ul,{children:["\n",(0,r.jsxs)(i.li,{children:[(0,r.jsx)(i.code,{children:".Recipe.APIVersion"}),": The API version which the recipe file uses"]}),"\n",(0,r.jsxs)(i.li,{children:[(0,r.jsx)(i.code,{children:".Recipe.Name"}),": The name of the recipe"]}),"\n",(0,r.jsxs)(i.li,{children:[(0,r.jsx)(i.code,{children:".Recipe.Version"}),": The current version of the recipe"]}),"\n"]}),"\n"]}),"\n",(0,r.jsxs)(i.li,{children:[(0,r.jsx)(i.code,{children:".ID"}),": UUID which is generated after the first execution of the recipe. It will keep its value over upgrades. Can be used to generate unique pseudo-random values which stays the same over the upgrades, for example ",(0,r.jsx)(i.code,{children:"my-resource-{{ sha1sum .ID | trunc 5 }}"})]}),"\n",(0,r.jsxs)(i.li,{children:[(0,r.jsx)(i.code,{children:".Variables"}),": Object which contains the values of the variables defined for the recipe. Example: ",(0,r.jsx)(i.code,{children:"{{ .Variables.FOO }}"})]}),"\n"]}),"\n",(0,r.jsx)(i.h3,{id:"template-only-specific-files",children:"Template only specific files"}),"\n",(0,r.jsxs)(i.p,{children:["By defining ",(0,r.jsx)(i.code,{children:"templateExtension"})," property in the ",(0,r.jsx)(i.code,{children:"recipe.yml"})," file, you can define that only the files with the given file extensions should be rendered from the ",(0,r.jsx)(i.code,{children:"templates"})," directory. The rest of the files will be copied as is. This is useful when there are files that do not need templating, but you would still need to escape the ",(0,r.jsx)(i.code,{children:"{{"})," and ",(0,r.jsx)(i.code,{children:"}}"})," characters (for example ",(0,r.jsx)(i.a,{href:"https://taskfile.dev/usage/",children:"Taskfiles"}),")."]}),"\n",(0,r.jsx)(i.h2,{id:"variables",children:"Variables"}),"\n",(0,r.jsxs)(i.p,{children:["Recipe variables let you define values which user need to provide to be able to render the tempaltes. Variables are defined in the ",(0,r.jsx)(i.code,{children:"recipe.yml"})," file. You can check schema ",(0,r.jsx)(i.a,{href:"/api#variable",children:"here"}),"."]}),"\n",(0,r.jsx)(i.h3,{id:"variable-types",children:"Variable types"}),"\n",(0,r.jsx)(i.p,{children:"Recipe variables supports the following types:"}),"\n",(0,r.jsxs)(i.ul,{children:["\n",(0,r.jsx)(i.li,{children:(0,r.jsx)(i.a,{href:"https://github.com/futurice/jalapeno/blob/main/examples/variable-types/recipe.yml#L9-L11",children:"String"})}),"\n",(0,r.jsx)(i.li,{children:(0,r.jsx)(i.a,{href:"https://github.com/futurice/jalapeno/blob/main/examples/variable-types/recipe.yml#L13-L15",children:"Boolean"})}),"\n",(0,r.jsx)(i.li,{children:(0,r.jsx)(i.a,{href:"https://github.com/futurice/jalapeno/blob/main/examples/variable-types/recipe.yml#L20-L22",children:"Select (predefined options)"})}),"\n",(0,r.jsx)(i.li,{children:(0,r.jsx)(i.a,{href:"https://github.com/futurice/jalapeno/blob/main/examples/variable-types/recipe.yml#L29-L38",children:"Table"})}),"\n"]}),"\n",(0,r.jsxs)(i.p,{children:["You can see examples of all the possible variables in the ",(0,r.jsx)(i.a,{href:"https://github.com/futurice/jalapeno/blob/main/examples/variable-types/recipe.yml",children:"example recipe"}),"."]}),"\n",(0,r.jsx)(i.admonition,{type:"note",children:(0,r.jsxs)(i.p,{children:["If you need to use numbers in the templates, you can use the ",(0,r.jsx)(i.code,{children:"atoi"})," function to convert a string variable to integer: ",(0,r.jsx)(i.code,{children:"{{ .Variables.FOO | atoi }}"})]})}),"\n",(0,r.jsx)(i.h3,{id:"validation",children:"Validation"}),"\n",(0,r.jsxs)(i.p,{children:["Variables can be validated by defining ",(0,r.jsx)(i.a,{href:"/api#variable",children:(0,r.jsx)(i.code,{children:"validators"})})," property for the variable. Validators support regular expression pattern matching."]}),"\n",(0,r.jsx)(i.h2,{id:"publishing-recipes",children:"Publishing recipes"}),"\n",(0,r.jsx)(i.p,{children:"Jalapeno supports publishing and storing recipes in OCI compatible registries. This means that versioned recipes can be pushed and pulled to/from ordinary Container Registries. This is useful when you want to make your recipes available or you want to check for updates for the recipe manually or programatically from CI/CD pipeline."}),"\n",(0,r.jsx)(i.h3,{id:"pushing-a-recipe-to-container-registry",children:"Pushing a recipe to Container registry"}),"\n",(0,r.jsxs)(i.p,{children:["You can push a recipe to a Container registry by using ",(0,r.jsx)(i.code,{children:"jalapeno push"})," command. For authentication you can use ",(0,r.jsx)(i.code,{children:"docker login"})," before pushing the recipe or provide credentials directly to ",(0,r.jsx)(i.code,{children:"jalapeno push"})," command by using flags. For example:"]}),"\n",(0,r.jsx)(i.pre,{children:(0,r.jsx)(i.code,{className:"language-bash",children:"jalapeno create my-recipe\n\n# You can find possible authentication methods to Github Container Registry at https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry#authenticating-to-the-container-registry\ndocker login ghcr.io -u my-user\n\njalapeno push my-recipe ghcr.io/my-user/my-recipe\n"})}),"\n",(0,r.jsx)(i.p,{children:"After this you should be able to see the recipe in the Container registry from the UI or by running:"}),"\n",(0,r.jsx)(i.pre,{children:(0,r.jsx)(i.code,{className:"language-bash",children:"docker inspect ghcr.io/my-user/my-recipe:v0.0.0\n"})}),"\n",(0,r.jsx)(i.admonition,{type:"note",children:(0,r.jsxs)(i.p,{children:["The tag of the recipe image in Container Registry is determined by the version in ",(0,r.jsx)(i.code,{children:"recipe.yml"})," file. So if you want to push a new version of the recipe, you need to update the version in ",(0,r.jsx)(i.code,{children:"recipe.yml"})," file first."]})}),"\n",(0,r.jsx)(i.h3,{id:"executing-a-recipe-from-container-registry",children:"Executing a recipe from Container registry"}),"\n",(0,r.jsxs)(i.p,{children:["You can execute a recipe directly from Container registry by using ",(0,r.jsx)(i.code,{children:"jalapeno execute"})," command. For example:"]}),"\n",(0,r.jsx)(i.pre,{children:(0,r.jsx)(i.code,{className:"language-bash",children:"mkdir my-project && cd my-project\njalapeno execute oci://ghcr.io/my-user/my-recipe:v0.0.0\n"})}),"\n",(0,r.jsx)(i.p,{children:"Another way is to pull the recipe first on your local machine and then execute it:"}),"\n",(0,r.jsx)(i.pre,{children:(0,r.jsx)(i.code,{className:"language-bash",children:"mkdir my-project && cd my-project\njalapeno pull oci://ghcr.io/my-user/my-recipe:v0.0.0\njalapeno execute my-recipe\n"})}),"\n",(0,r.jsx)(i.h3,{id:"checking-updates-for-a-recipe",children:"Checking updates for a recipe"}),"\n",(0,r.jsxs)(i.p,{children:["After you have executed a recipe, Jalapeno will create a ",(0,r.jsx)(i.code,{children:".jalapeno/sauces.yml"})," file to the project directory. This file contains information about the executed recipes and their versions. If the recipe was executed directly from a Container Registry, the registry URL is also stored in the file for checking for new versions. To check the new versions for the recipes, you can run:"]}),"\n",(0,r.jsx)(i.pre,{children:(0,r.jsx)(i.code,{className:"language-bash",children:"jalapeno check\n"})}),"\n",(0,r.jsx)(i.p,{children:"This will check for updates for all the recipes in the project directory. You can also check for updates for a specific recipe by running:"}),"\n",(0,r.jsx)(i.pre,{children:(0,r.jsx)(i.code,{className:"language-bash",children:"jalapeno check -r my-recipe\n"})}),"\n",(0,r.jsx)(i.p,{children:"If you've executed the recipe from local directory and the registry URL is still unknown, you can set the registry URL manually by running:"}),"\n",(0,r.jsx)(i.pre,{children:(0,r.jsx)(i.code,{className:"language-bash",children:"jalapeno check -r my-recipe --from oci://ghcr.io/my-user/my-recipe\n"})}),"\n",(0,r.jsx)(i.admonition,{type:"tip",children:(0,r.jsxs)(i.p,{children:["If you want to run the check in a CI/CD pipeline (like Github Actions), you can check the ",(0,r.jsx)(i.a,{href:"https://github.com/futurice/jalapeno/tree/main/examples/github-action",children:(0,r.jsx)(i.code,{children:"examples/github-action"})})," recipe how to do it or you can execute it to your project with ",(0,r.jsx)(i.code,{children:"jalapeno execute oci://ghcr.io/futurice/jalapeno/examples/github-action:v0.0.0"}),"."]})})]})}function d(e={}){const{wrapper:i}={...(0,a.a)(),...e.components};return i?(0,r.jsx)(i,{...e,children:(0,r.jsx)(h,{...e})}):h(e)}}}]);