---
slug: /
sidebar_position: 1
title: Home
---

# Jalapeno

Jalapeno is **a project templating system** which support complex templating while staying developer friendly. Project templates used by Jalapeno are called _recipes_.

## Features

- **Modular**: Project can contain multiple recipes, so you can compose your projects from multiple smaller templates modules.
- **Tool agnostic**: Recipes can define whatever tools to be used in the project, and Jalapeno related files can be ejected from the project at any time.
- **Continous integration**: Recipes are versioned, and new versions of the recipe can be merged to existing projects. Recipes can be shared via OCI compatible (aka container) registries so project CI/CD pipelines which already utilize the registry can easily check and notify the developers if there are new versions available for the recipes.
- **Snapshots tests**: Recipes can define snapshot tests, which reduces regression by ensuring that the templates produces expected outputs.
