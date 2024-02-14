# Changelog

## [1.4.0](https://github.com/futurice/jalapeno/compare/v1.3.5...v1.4.0) (2024-02-14)


### Features

* support executing recipe to a subpath ([c3777fc](https://github.com/futurice/jalapeno/commit/c3777fc383e4b6fdc158494ab03b13ab54bb4073))
* support templates in initHelp ([75da7fa](https://github.com/futurice/jalapeno/commit/75da7fa6e75f51d9e9dc0a04617beec56c401e10))

## [1.3.5](https://github.com/futurice/jalapeno/compare/v1.3.4...v1.3.5) (2024-02-07)


### Bug Fixes

* make the list of conflicted files visible after upgrading ([b50bdb4](https://github.com/futurice/jalapeno/commit/b50bdb47db9efd955cf3c749a697016bbd5a765c))

## [1.3.4](https://github.com/futurice/jalapeno/compare/v1.3.3...v1.3.4) (2024-02-06)


### Bug Fixes

* use single source instead of a list ([40d453b](https://github.com/futurice/jalapeno/commit/40d453bed93b0a9dfabadcb8e70ae3d982f03f49))

## [1.3.3](https://github.com/futurice/jalapeno/compare/v1.3.2...v1.3.3) (2024-02-05)


### Bug Fixes

* retain CheckFrom URL from old sauce ([b1eef93](https://github.com/futurice/jalapeno/commit/b1eef93dfaece6f7f3dfb38711b5f2bb9608b4cf))

## [1.3.2](https://github.com/futurice/jalapeno/compare/v1.3.1...v1.3.2) (2024-02-05)


### Continuous Integration

* add published tag after publishing ([5bfd1e8](https://github.com/futurice/jalapeno/commit/5bfd1e85fd3205ec3fd82db1f55ee3a63adf0102))

## [1.3.1](https://github.com/futurice/jalapeno/compare/v1.3.0...v1.3.1) (2024-02-05)


### Bug Fixes

* **survey:** avoid infinite loop if first variable is skipped ([aeadec7](https://github.com/futurice/jalapeno/commit/aeadec7239b0464af1c3c5df9334648db457cdc6))

## [1.3.0](https://github.com/futurice/jalapeno/compare/v1.2.1...v1.3.0) (2024-02-02)


### Features

* show proper diff for failing snapshot tests ([e61eb0c](https://github.com/futurice/jalapeno/commit/e61eb0c09a78860ff34f160ef1f59c3340b5b229))
* try variable default values when using no-input mode ([7030c7c](https://github.com/futurice/jalapeno/commit/7030c7c97c75278d40d987cd47267434ce3525a7))


### Bug Fixes

* handle files manually deleted by the user ([c4ecb80](https://github.com/futurice/jalapeno/commit/c4ecb80017a55a7cd006245674e49b9fce8046f8))
* **survey:** reflect actual confirm value ([ecf631f](https://github.com/futurice/jalapeno/commit/ecf631fac2ea7891d9decdf28ea3d97597ee7609))

## [1.2.1](https://github.com/futurice/jalapeno/compare/v1.2.0...v1.2.1) (2024-01-26)


### Bug Fixes

* remove tag from check URL ([d009a89](https://github.com/futurice/jalapeno/commit/d009a89cfdfcb3023e89c6cf79a743c5efeaf3b8))
* **survey:** do not validate cell if the table variable is optional and the table is empty ([40ff5e6](https://github.com/futurice/jalapeno/commit/40ff5e69574375817e218e0a73ef8644fca42768))
* **survey:** handle window resizing ([6fd6ed5](https://github.com/futurice/jalapeno/commit/6fd6ed5acb6f0f04bbcfd9e81950adb1ae33d120))

## [1.2.0](https://github.com/futurice/jalapeno/compare/v1.1.1...v1.2.0) (2024-01-22)


### Features

* allow upgrades to same version to redefine variable values ([f84b14d](https://github.com/futurice/jalapeno/commit/f84b14db98d355f6f72c7cd74fe8391ddadbb0fd))


### Bug Fixes

* unmarshal TableValue correctly ([a05f7e5](https://github.com/futurice/jalapeno/commit/a05f7e52607a87b737710b832d3264b57f395894))

## [1.1.1](https://github.com/futurice/jalapeno/compare/v1.1.0...v1.1.1) (2024-01-18)


### Bug Fixes

* check if table variable is optional ([4860bfa](https://github.com/futurice/jalapeno/commit/4860bfa7fada9fca898894ad8a572cf6bf16620a))
* handle boolean values correctly in retry message ([9432685](https://github.com/futurice/jalapeno/commit/9432685a3c4cfac8d3746f2181cf921f54bd8e82))
* make sure that table column order is consistent ([799c1d0](https://github.com/futurice/jalapeno/commit/799c1d0fff79eceb5bd24fa181faf870765b9580))
* save recipe correctly when directory name does not match recipe name ([d4a23bd](https://github.com/futurice/jalapeno/commit/d4a23bd6c191e9f798b3597797b72ac6a55486eb))
* validate table cells when value is given with flags ([65677bb](https://github.com/futurice/jalapeno/commit/65677bbf2b142395832878c918373786eb4c6162))

## [1.1.0](https://github.com/futurice/jalapeno/compare/v1.0.0...v1.1.0) (2024-01-17)


### Features

* upgrade sprig to v3 ([882f10e](https://github.com/futurice/jalapeno/commit/882f10ec2754d6b6dc413f6fb417eaa0470e8018))

## 1.0.0 (2023-12-22)

* Initial release
