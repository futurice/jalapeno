# Changelog

## [1.14.4](https://github.com/futurice/jalapeno/compare/v1.14.3...v1.14.4) (2024-10-31)


### Bug Fixes

* show versions correctly when bumping recipe version ([78f325b](https://github.com/futurice/jalapeno/commit/78f325bf30491fd351adadda60a63412de00ff56))

## [1.14.3](https://github.com/futurice/jalapeno/compare/v1.14.2...v1.14.3) (2024-10-15)


### Bug Fixes

* text input width calculation ([#144](https://github.com/futurice/jalapeno/issues/144)) ([e7e2b9d](https://github.com/futurice/jalapeno/commit/e7e2b9d43f6f88d9de9a728ce4492d5936defa31))

## [1.14.2](https://github.com/futurice/jalapeno/compare/v1.14.1...v1.14.2) (2024-10-04)


### Bug Fixes

* validate existing table values correctly ([9309384](https://github.com/futurice/jalapeno/commit/93093845c6317e86c13710a8d20526f493889971))

## [1.14.1](https://github.com/futurice/jalapeno/compare/v1.14.0...v1.14.1) (2024-10-03)


### Bug Fixes

* allow updating values quickly with --set flag ([1875031](https://github.com/futurice/jalapeno/commit/1875031feb46c11e59e1a52c66942a2fa0487073))

## [1.14.0](https://github.com/futurice/jalapeno/compare/v1.13.0...v1.14.0) (2024-10-02)


### Features

* improved merge conflict solver ([#139](https://github.com/futurice/jalapeno/issues/139)) ([eb4b162](https://github.com/futurice/jalapeno/commit/eb4b162c4c8c844a9b186386e5de435d58916604))
* introduce `bumpver` command for bumping recipe version ([7f03a31](https://github.com/futurice/jalapeno/commit/7f03a31f100342858ed67df9a33e8e39f1cff2f1))


### Bug Fixes

* **ci:** provide whole commit history to commitlint ([efd35b7](https://github.com/futurice/jalapeno/commit/efd35b7231297493b5cffa7ab14b9e104d0db8af))
* select output correctly when using diff ([06e5f95](https://github.com/futurice/jalapeno/commit/06e5f953915ab6b09524243b328f8adb8ad2b42a))

## [1.13.0](https://github.com/futurice/jalapeno/compare/v1.12.1...v1.13.0) (2024-08-26)


### Features

* support optional multi select variables ([#136](https://github.com/futurice/jalapeno/issues/136)) ([0599255](https://github.com/futurice/jalapeno/commit/059925528382f2546f782f162a013f707f654b54))

## [1.12.1](https://github.com/futurice/jalapeno/compare/v1.12.0...v1.12.1) (2024-08-09)


### Bug Fixes

* remove redundant debug flag ([2cd5553](https://github.com/futurice/jalapeno/commit/2cd55531602628706b001debe9c0bcd6c513b036))

## [1.12.0](https://github.com/futurice/jalapeno/compare/v1.11.8...v1.12.0) (2024-08-01)


### Features

* implement multi-select variable type ([#132](https://github.com/futurice/jalapeno/issues/132)) ([cd0cc02](https://github.com/futurice/jalapeno/commit/cd0cc023a6f9280bc7ccfe345f9b0e93901f85b6))

## [1.11.8](https://github.com/futurice/jalapeno/compare/v1.11.7...v1.11.8) (2024-07-18)


### Bug Fixes

* check for extra variables defined in tests ([f51a639](https://github.com/futurice/jalapeno/commit/f51a639f04a2d05d0e8863e3a4f5acfc81d87915))
* update expectedInitHelp when updating snapshots ([01074df](https://github.com/futurice/jalapeno/commit/01074dfeb7e31d49b044ac37e898bcd75ed95a75))

## [1.11.7](https://github.com/futurice/jalapeno/compare/v1.11.6...v1.11.7) (2024-07-17)


### Bug Fixes

* skip prompting when no-input is set ([99e778f](https://github.com/futurice/jalapeno/commit/99e778f04b237ab66ede027178c42648c75d335a))

## [1.11.6](https://github.com/futurice/jalapeno/compare/v1.11.5...v1.11.6) (2024-07-16)


### Bug Fixes

* add new lines to check logs ([e412bf6](https://github.com/futurice/jalapeno/commit/e412bf60aee046d7870c1b300c1918a5b35eb6df))

## [1.11.5](https://github.com/futurice/jalapeno/compare/v1.11.4...v1.11.5) (2024-07-15)


### Miscellaneous Chores

* allow "create test" to be run multiple times ([641c7cc](https://github.com/futurice/jalapeno/commit/641c7cc3f69579456aca68fb8bc417dd89130f6f))

## [1.11.4](https://github.com/futurice/jalapeno/compare/v1.11.3...v1.11.4) (2024-07-11)


### Miscellaneous Chores

* try out publish workflow changes ([e1835a6](https://github.com/futurice/jalapeno/commit/e1835a6adb4a11d91475950e1da534487ea329a3))

## [1.11.3](https://github.com/futurice/jalapeno/compare/v1.11.2...v1.11.3) (2024-07-11)


### Bug Fixes

* allow same recipe to be executed twice ([330e6d0](https://github.com/futurice/jalapeno/commit/330e6d02a6195820abc248d930b5ed88490da592))

## [1.11.2](https://github.com/futurice/jalapeno/compare/v1.11.1...v1.11.2) (2024-07-10)


### Bug Fixes

* run rest of the tests when expectInitHelp is defined ([aab6c53](https://github.com/futurice/jalapeno/commit/aab6c53a04f2636bdd6cf07cdcd17c68e8e6f5ef))

## [1.11.1](https://github.com/futurice/jalapeno/compare/v1.11.0...v1.11.1) (2024-06-20)


### Features

* show which files were extra or missing when testing ([c5a4ddd](https://github.com/futurice/jalapeno/commit/c5a4ddd29a78d3077dce2fb5dd766b74ebaf2bb0))

## [1.11.0](https://github.com/futurice/jalapeno/compare/v1.10.4...v1.11.0) (2024-06-19)


### Features

* create test for initHelp ([f3839d6](https://github.com/futurice/jalapeno/commit/f3839d6ef9fecb2e30c02165fcc4a69050a269a7))

## [1.10.4](https://github.com/futurice/jalapeno/compare/v1.10.3...v1.10.4) (2024-06-18)


### Bug Fixes

* show "user aborted" correctly when executing a recipe ([74b9e25](https://github.com/futurice/jalapeno/commit/74b9e254e9b6c44a47c3585a8fcc9f21198f0d68))

## [1.10.3](https://github.com/futurice/jalapeno/compare/v1.10.2...v1.10.3) (2024-06-18)


### Bug Fixes

* set text input width dynamically ([41b1e75](https://github.com/futurice/jalapeno/commit/41b1e750b401702381002fc25f8356e6bdf5b5ed))

## [1.10.2](https://github.com/futurice/jalapeno/compare/v1.10.1...v1.10.2) (2024-06-13)


### Features

* allow forcing upgrade ([1492cb1](https://github.com/futurice/jalapeno/commit/1492cb105b5249c4c1e5b556e536bade2352a4be))

## [1.10.1](https://github.com/futurice/jalapeno/compare/v1.10.0...v1.10.1) (2024-06-12)


### Bug Fixes

* relative recipe paths starting with ".." did not work ([e8618a4](https://github.com/futurice/jalapeno/commit/e8618a46541193dca6ec61ad0ad717e3295cb6b6))
* validate old sauce values when upgrading a recipe ([#110](https://github.com/futurice/jalapeno/issues/110)) ([3a9590b](https://github.com/futurice/jalapeno/commit/3a9590b209b48ac48a6ab05ead8a3f47dfb50864))

## [1.10.0](https://github.com/futurice/jalapeno/compare/v1.9.2...v1.10.0) (2024-06-11)


### Features

* variables in template file names ([#105](https://github.com/futurice/jalapeno/issues/105)) ([1f45505](https://github.com/futurice/jalapeno/commit/1f4550509485a300f989e3228d5d9a0a0a5cdb34))

## [1.9.2](https://github.com/futurice/jalapeno/compare/v1.9.1...v1.9.2) (2024-06-06)


### Bug Fixes

* forward stderr to GH Action results output ([c79148e](https://github.com/futurice/jalapeno/commit/c79148e307c97d1246f055c3b978bb5e2384be7b))

## [1.9.1](https://github.com/futurice/jalapeno/compare/v1.9.0...v1.9.1) (2024-06-05)


### Bug Fixes

* move no-input flag as a common flag ([ff94f72](https://github.com/futurice/jalapeno/commit/ff94f72b8ed29769cd135664ed83f11b8743f72a))

## [1.9.0](https://github.com/futurice/jalapeno/compare/v1.8.2...v1.9.0) (2024-05-14)


### Features

* unique table column validator ([#99](https://github.com/futurice/jalapeno/issues/99)) ([d187ae9](https://github.com/futurice/jalapeno/commit/d187ae90a8a66cfcd9c7043c485695808cc7f077))

## [1.8.2](https://github.com/futurice/jalapeno/compare/v1.8.1...v1.8.2) (2024-05-07)


### Bug Fixes

* add helper texts for create command ([0e8653f](https://github.com/futurice/jalapeno/commit/0e8653ff74ef193a9b4b92d594990eddb7d22352))

## [1.8.1](https://github.com/futurice/jalapeno/compare/v1.8.0...v1.8.1) (2024-04-17)


### Bug Fixes

* skip already executed recipes when executing a manifest ([9463a4f](https://github.com/futurice/jalapeno/commit/9463a4fe3eae89955df7a0e04957ffcc26b88bec))

## [1.8.0](https://github.com/futurice/jalapeno/compare/v1.7.2...v1.8.0) (2024-04-15)


### Features

* introduce manifest file to execute multiple recipes ([#94](https://github.com/futurice/jalapeno/issues/94)) ([04cc757](https://github.com/futurice/jalapeno/commit/04cc75745344658c7df9f62f809a29b41d5bca89))
* show empty variable value when variable was optional ([de97ad9](https://github.com/futurice/jalapeno/commit/de97ad900feec1cfa0423006732f6c6a6f6a2ebc))
* show if the variable is optional ([dc265a4](https://github.com/futurice/jalapeno/commit/dc265a4d7aaffdbe2a4b5698f30e3b318f847bf0))

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
