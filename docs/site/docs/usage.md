---
sidebar_position: 3
slug: /usage/
title: Usage
---

# Usage

## Getting started

`jalapeno create my-recipe`

- Add variable
- Use that variable in templates

`jalapeno execute my-recipe`

## Variables

### Validation

### Conditional variables

## Templating

Templates are done by using [Go templates](https://pkg.go.dev/text/template).

The following context is available on the templates:

- `Recipe`: Metadata object of the recipe
  - `Recipe.Name`: The name of the recipe
  - `Recipe.Version`: The current version of the recipe
- `ID`: UUID which is generated after the first execution of the recipe. It will keep its value over upgrades. Can be used to generate unique pseudo-random values which stays the same over the upgrades, for example `my-resource-{{ sha1sum .ID | trunc 5 }}`
- `Variables`: Object which contains the values of the variables defined for the recipe. Example: `{{ .Variables.FOO }}`

## Pushing recipe to Container Registry

### Executing a recipe from local path

### Executing a recipe from Container registry

### Checking updates for a recipe

`jalapeno check`
