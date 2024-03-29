# Changelog

## [1.7.2](https://github.com/futurice/jalapeno/compare/v1.7.1...v1.7.2) (2024-03-13)


### Bug Fixes

* show prompt help also when using default value ([dd56128](https://github.com/futurice/jalapeno/commit/dd5612810154cd043d67d3e9d849a16a38093ee0))

## [1.7.1](https://github.com/futurice/jalapeno/compare/v1.7.0...v1.7.1) (2024-03-13)


### Bug Fixes

* remove all tags from URL ([e94a331](https://github.com/futurice/jalapeno/commit/e94a3313fe7cb526a9ffa74efbc4928bbf582f22))

## [1.7.0](https://github.com/futurice/jalapeno/compare/v1.6.2...v1.7.0) (2024-03-11)


### Features

* optionally push recipe to latest tag ([62fa60f](https://github.com/futurice/jalapeno/commit/62fa60f4d07a2bc6645c6767d8c867c32f10989a))

## [1.6.2](https://github.com/futurice/jalapeno/compare/v1.6.1...v1.6.2) (2024-02-29)


### Bug Fixes

* make the order of upgradeable sauces deterministic ([3994757](https://github.com/futurice/jalapeno/commit/399475757095967febc4dd66caf59dbeff0e7435))

## [1.6.1](https://github.com/futurice/jalapeno/compare/v1.6.0...v1.6.1) (2024-02-29)


### Bug Fixes

* honor user abort when upgrading recipes with check command ([7486129](https://github.com/futurice/jalapeno/commit/74861291c920ec422b979cb483e88caddaf09b2b))

## [1.6.0](https://github.com/futurice/jalapeno/compare/v1.5.0...v1.6.0) (2024-02-28)


### Features

* support directly upgrading recipes with check command ([f789c64](https://github.com/futurice/jalapeno/commit/f789c6456a112e59f04eb1ddad13ed388928067a))


### Bug Fixes

* when upgrading to same version, do not use previous values automatically ([9c32487](https://github.com/futurice/jalapeno/commit/9c3248735f2e3b112c911f9e367555af4b58ba03))

## [1.5.0](https://github.com/futurice/jalapeno/compare/v1.4.0...v1.5.0) (2024-02-19)


### Features

* add statuses to printed file trees ([b891c48](https://github.com/futurice/jalapeno/commit/b891c48a13e87545b09e6ed350463387d9c30783))
* print file tree after upgrade ([be1d443](https://github.com/futurice/jalapeno/commit/be1d443b7cc9afd38f2b5d6ed858b01e2fe00e14))


### Bug Fixes

* remove old files after upgrading ([6296c77](https://github.com/futurice/jalapeno/commit/6296c776a6051ddceb2697488bffb3c18b801f17))
* save files correctly after upgrading ([3c565c1](https://github.com/futurice/jalapeno/commit/3c565c1bd3fbd4a0402dfff08a002c40300d9a9d))

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
