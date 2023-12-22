# Changelog

## [1.0.0](https://github.com/futurice/jalapeno/compare/v0.1.30...v1.0.0) (2023-12-22)


### ⚠ BREAKING CHANGES

* breaks all existing sauces.yml files, migrate them by creating `apiVersion: "v1"` property and moving recipe metadata under `recipe` property

### Features

* **cli:** If execute fails, output command line with all vars as --set args ([#57](https://github.com/futurice/jalapeno/issues/57)) ([4400305](https://github.com/futurice/jalapeno/commit/44003053dc2027a69c9d17c1175739a4c9bcbce7))
* move recipe struct to dedicated property in sauces.yml ([2544cb5](https://github.com/futurice/jalapeno/commit/2544cb5db59a79834ad2fdea4847cc9937e4cc7f))


### Bug Fixes

* check manually created new files when upgrading ([33b7d2e](https://github.com/futurice/jalapeno/commit/33b7d2e7c4ad4ec8c5738e4ccd9b39624045eecc))
* **docs:** push terraform bootstrap to examples OCI repo ([#61](https://github.com/futurice/jalapeno/issues/61)) ([07b6e2d](https://github.com/futurice/jalapeno/commit/07b6e2d59516c32b37da2ac270073f0324312a62))

## [0.1.30](https://github.com/futurice/jalapeno/compare/v0.1.29...v0.1.30) (2023-12-18)


### Miscellaneous Chores

* release 0.1.30 ([f7ffca3](https://github.com/futurice/jalapeno/commit/f7ffca3873002527a24b818088fef549eee0a7e4))
