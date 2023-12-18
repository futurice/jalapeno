---
slug: /
sidebar_position: 1
title: Home
---

# Jalapeno

Jalapeno is **a project templating system** which support complex templating while staying developer friendly. Project templates used by Jalapeno are called _recipes_.

## Features

- **Easy to use**: The CLI guides the user as much as possible when using recipes. This is done by prompting the user [interactively](https://github.com/charmbracelet/bubbletea) for required template values, provides validation for the values, giving hints about what to do after executing the recipe etc.
- **Modular**: Project can contain multiple recipes, so you can compose your projects from multiple smaller templates modules.
- **Tool agnostic**: Recipes can define whatever tools to be used in the project, and Jalapeno related files can be ejected from the project at any time.
- **Continous integration**: Recipes are versioned, and new versions of the recipe can be merged to existing projects. Recipes can be shared via OCI compatible registries (aka container registries) so project CI/CD pipelines which already utilize the registry can easily check and notify the developers if there are new versions available for the recipes.
- **Snapshots tests**: Recipes can define snapshot tests, which reduces regression when developing the templates by ensuring that the templates produces expected outputs.
