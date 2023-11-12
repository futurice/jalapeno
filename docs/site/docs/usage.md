---
sidebar_position: 3
slug: /usage/
title: Usage
---

# Usage [WIP]

## Getting started

To get started with Jalapeno, you need to install the CLI tool. You can find the installation instructions [here](/installation).

Then you can bootstrap a new _recipe_ (aka template) by running:

```bash
jalapeno create my-recipe
```

After this you should have a new folder called `my-recipe` with the following structure:

```
my-recipe
├── recipe.yml
├── templates
│   └── README.md
└── tests
    └── defaults
        ├── test.yml
        └── files
            └── README.md
```

The `templates` directory contains the templates which will be rendered to the project directory. You can add and edit files there or you can already execute this _recipe_ (render the templates) to your project directory by running:

```bash
mkdir my-project
jalapeno execute my-recipe -d my-project
```

You can also execute any of the examples from the Jalapeno repository, for example:

```bash
jalapeno execute oci://ghcr.io/futurice/jalapeno/examples/variable-types:0.0.0 -d my-project
```

After this, the project directory should have the following files:

```
my-project
├── .jalapeno
│   └── sauces.yml
└── README.md
```

The `.jalapeno` directory contains files which Jalapeno uses internally. For example `sauces.yml` is Jalapeno metadata file which contains information about the _sauces_ (aka executed recipes). This file is used to check for updates for the recipes later.

The rest of the files are rendered from the templates. You can edit the templates and execute the recipe again to update the files.

## Templating

Templates are done by using [Go templates](https://pkg.go.dev/text/template).

The following context is available on the templates:

- `Recipe`: Metadata object of the recipe
  - `Recipe.APIVersion`: The API version which the recipe file uses
  - `Recipe.Name`: The name of the recipe
  - `Recipe.Version`: The current version of the recipe
- `ID`: UUID which is generated after the first execution of the recipe. It will keep its value over upgrades. Can be used to generate unique pseudo-random values which stays the same over the upgrades, for example `my-resource-{{ sha1sum .ID | trunc 5 }}`
- `Variables`: Object which contains the values of the variables defined for the recipe. Example: `{{ .Variables.FOO }}`

## Variables

### Validation

### Conditional variables

## Pushing recipe to Container Registry

### Executing a recipe from local path

### Executing a recipe from Container registry

### Checking updates for a recipe

`jalapeno check`
