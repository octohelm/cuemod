# Change Log

All notable changes to this project will be documented in this file.
See [Conventional Commits](https://conventionalcommits.org) for commit guidelines.



# [0.3.2](https://github.com/octohelm/cuemod/compare/v0.3.1...v0.3.2)

### Bug Fixes

* **fix** flags fix cuem fmt ([ceba578](https://github.com/octohelm/cuemod/commit/ceba57888cad488c8396ae1e4686d4c2a280dbb0))



# [0.3.1](https://github.com/octohelm/cuemod/compare/v0.3.0...v0.3.1)

### Bug Fixes

* **fix(cuemod):** local path replace only work for root mod ([3e8cb5a](https://github.com/octohelm/cuemod/commit/3e8cb5a9b214a874961b0f7200444fa6e9950e22))



# [0.3.0](https://github.com/octohelm/cuemod/compare/v0.2.2...v0.3.0)

### Features

* **feat** remove auto delect for extracting to cue ([a0c88d2](https://github.com/octohelm/cuemod/commit/a0c88d2abda93056393307cb4f3ddc9f30829bbc))



# [0.2.2](https://github.com/octohelm/cuemod/compare/v0.2.1...v0.2.2)



# [0.2.1](https://github.com/octohelm/cuemod/compare/v0.2.0...v0.2.1)



# [0.2.0](https://github.com/octohelm/cuemod/compare/v0.1.1...v0.2.0)

### Features

* **feat** added support to overwrites from cli ([31ff716](https://github.com/octohelm/cuemod/commit/31ff716c606f30a09a8bc4e83b43656e2f3ca880))
* **feat** using server-apply ([3870d92](https://github.com/octohelm/cuemod/commit/3870d92ba5dfdfa1a721c0a3540b843774059d62))



# [0.1.1](https://github.com/octohelm/cuemod/compare/v0.1.0...v0.1.1)



# [0.1.0](https://github.com/octohelm/cuemod/compare/v0.0.4...v0.1.0)

### Features

* **feat** show -o with .yaml ext will generate to one single manifest file ([5b47606](https://github.com/octohelm/cuemod/commit/5b476060deb1c4035050b02f068b9326ed014d54))



# [0.0.4](https://github.com/octohelm/cuemod/compare/v0.0.3...v0.0.4)



# [0.0.3](https://github.com/octohelm/cuemod/compare/v0.0.2...v0.0.3)



# [0.0.2](https://github.com/octohelm/cuemod/compare/v0.0.1...v0.0.2)



# [0.0.1](https://github.com/octohelm/cuemod/compare/v0.0.0...v0.0.1)



# 0.0.0

### Bug Fixes

* **fix** double require statements in module.cue like golang 1.17 did ([5d9a8f3](https://github.com/octohelm/cuemod/commit/5d9a8f32c4d2a87c4e518e09ebc27eb4d92118d0))
* **fix** enhance apply ([018e0ee](https://github.com/octohelm/cuemod/commit/018e0eef62f1844333bb15d809d27355afae3fc6))
* **fix** v1.PersistentVolumeClaim should be merge ([073414d](https://github.com/octohelm/cuemod/commit/073414d8a05c624fc61f4b3a39992e3df2453d06))
* **fix** always set namespace ([d33bfeb](https://github.com/octohelm/cuemod/commit/d33bfebcf66b108673e6690cadbf255beb25710a))
* **fix** jsonnet.alias/xxx hook for replace jsonnet ([5e46c25](https://github.com/octohelm/cuemod/commit/5e46c25b0f6975e1daa103a5e96e1ef1d682fdaf))
* **fix(extractor/crd):** handle additionalProps ([79f8b30](https://github.com/octohelm/cuemod/commit/79f8b30663717f422f9e68d06f7f778ac3e55c61))
* **fix(extractor/crd):** fix boolean ([df542c6](https://github.com/octohelm/cuemod/commit/df542c64553e5cd8ef7b01597609d39443d4e5b8))
* **fix(extractor/crd):** should gen const apiVersion & kind ([e8829b1](https://github.com/octohelm/cuemod/commit/e8829b19c1b38717e9bb89a1d57dae2992beb341))
* **fix** translate should support nested ([f853a99](https://github.com/octohelm/cuemod/commit/f853a992b883d51e534cd3ab29c70dbfbece2efd))
* **fix(translator/toml):** if number can be int, to convert it to int ([b8e7967](https://github.com/octohelm/cuemod/commit/b8e79675cac70760cf246b26b6720340e0abd913))
* **fix** don't link root mod to cue.mod/usr ([a1e7415](https://github.com/octohelm/cuemod/commit/a1e741563c5d27bb4cb10f429b60d538264eed84))
* **fix(translator/helm):** should trim template before yaml decode ([4a014e3](https://github.com/octohelm/cuemod/commit/4a014e3d4a19ebbc750d26a94cd0e6666266402f))
* **fix(tanka):** should inject namespace ([06387c9](https://github.com/octohelm/cuemod/commit/06387c9dceda186e7c27f900a4356961da9173a0))
* **fix** print cue validate error stacks ([80dde6d](https://github.com/octohelm/cuemod/commit/80dde6d10d31cd8e5d4dd3dc624b727e71f1914d))
* **fix** enhance k plugin ([415d601](https://github.com/octohelm/cuemod/commit/415d601e33453f12beb137b0f8bfe6a02dfe2be7))
* **fix(extractor/helm):** file content should be simple bytes to avoid ident issue in yaml ([64722e5](https://github.com/octohelm/cuemod/commit/64722e5991b3078e723d0ca3d0b806a0c7ac91df))
* **fix(extractor/core):** should avoid symlink for hash dir ([9042c3d](https://github.com/octohelm/cuemod/commit/9042c3db39983d5ad2b3619a69c8ca3a2c7dac5b))
* **fix(extractor/helm):** values should all optional ([52b4772](https://github.com/octohelm/cuemod/commit/52b477267fb0155926c9fc23f8833d2901f49495))
* **fix** init ([267b96a](https://github.com/octohelm/cuemod/commit/267b96a51dfe9a1c7e9d669021b96897ad901ae3))


### Features

* **feat** cuem-operator ([5ac9092](https://github.com/octohelm/cuemod/commit/5ac9092a78d694f700848b006ab097b76baf9aa5))
* **feat** eval support bundle to single file ([e150778](https://github.com/octohelm/cuemod/commit/e15077809917b71e07fe2923359e30e201be7e3d))
* **feat** extractor for crd ([cad3b38](https://github.com/octohelm/cuemod/commit/cad3b385d032f5eb8f25ebc7ef20848388274abe))
* **feat** add translator yaml ([c5f36b3](https://github.com/octohelm/cuemod/commit/c5f36b324497fcb998ff177c4cfa3909414ec573))
* **feat** sum file & download only if needs ([478519c](https://github.com/octohelm/cuemod/commit/478519ce3fb25647a91779c34f0895086cf29a71))
