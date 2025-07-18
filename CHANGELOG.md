# Changelog

All notable changes to this project will be documented in this file. See [standard-version](https://github.com/conventional-changelog/standard-version) for commit guidelines.

### [1.0.8](https://github.com/allex/envsubst/compare/v1.0.7...v1.0.8) (2025-07-14)


### Features

* add support for bash case conversion operators `${VAR^^}` and `${VAR,,}` ([d57f2df](https://github.com/allex/envsubst/commit/d57f2df5e40c4d4d14a1a81f88dcb7d477f92adb))

### [1.0.7](https://github.com/allex/envsubst/compare/v1.0.6...v1.0.7) (2025-07-10)


### Features

* add Strings method to Env for retrieving all environment variables as key-value strings ([67824dd](https://github.com/allex/envsubst/commit/67824ddc217d379dc52481de8362b064f6e4ff34))

### [1.0.6](https://github.com/allex/envsubst/compare/v1.0.5...v1.0.6) (2025-07-01)


### Features

* add env with lazy injection in envsubst ([5d37f59](https://github.com/allex/envsubst/commit/5d37f590590a0c3d55097236783039c7fd87b845))

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
