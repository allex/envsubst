# Changelog

All notable changes to this project will be documented in this file. See [standard-version](https://github.com/conventional-changelog/standard-version) for commit guidelines.

### [1.0.5](https://github.com/allex/envsubst/compare/v1.0.4...v1.0.5) (2025-07-01)


### Features

* enhance nested variable substitution handling in lexer and parser ([c5656f8](https://github.com/allex/envsubst/commit/c5656f8ee5051300e7bb8131d763ebc7feb37394))

## [1.0.4](https://github.com/allex/envsubst/compare/v1.0.3...v1.0.4) (2025-06-30)


### Features

* add comprehensive API documentation for envsubst package ([2d048c2](https://github.com/allex/envsubst/commit/2d048c29ad8517f7ad86dfd99152ca206122ba53))
* add KeepUnset option to preserve undefined variables in envsubst ([b87f06c](https://github.com/allex/envsubst/commit/b87f06c65beb6f1e85cb271ffd137344086d2d0e))
* implement custom variable matcher functionality in lexer tests ([b77c7b3](https://github.com/allex/envsubst/commit/b77c7b34d06c14f16c5b8bff4686382296179c37))


### Bug Fixes

* replace fmt.Fprintf with fmt.Fprint for error message output in envsubst ([11481e7](https://github.com/allex/envsubst/commit/11481e736c98db0cdabcde88dfd5051c1c31d600))

### [1.0.3](https://github.com/allex/envsubst/compare/v1.0.2...v1.0.3) (2023-05-12)


### Bug Fixes

* cleanup mismatched text-variable subsDepth ([1b21529](https://github.com/allex/envsubst/commit/1b21529eec0213b814a2954308e3acd726aff37e))

### 1.0.1 (2023-03-13)


### Features

* add default behaviour of running through all the errors ([#21](https://github.com/allex/envsubst/issues/21)) ([7bc4df4](https://github.com/allex/envsubst/commit/7bc4df48ec6140d0c5670a3eb2d484bed0f6c284))
* add varMatch predicate function to filter the valid variable names ([7df40d4](https://github.com/allex/envsubst/commit/7df40d43549259d06a50815c28dffc83f798994a))


### Bug Fixes

* example in readme ([de1f237](https://github.com/allex/envsubst/commit/de1f237918b5935914667cdee93129d5aed87eaa))
* **lex:** change lastPos assignment ([49bbf0c](https://github.com/allex/envsubst/commit/49bbf0c66cede47052267e1b0cd9a2af7059dae3))
* test if something being pipe to stdin ([569d354](https://github.com/allex/envsubst/commit/569d3548760e3d1c52c6c5bbd40c8ecf2fa3495d))
