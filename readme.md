# Go Service Boilerplate

## Goals

It is a very opinionated golang service boilerplate, I use it as a study practice and build many features to it as I progress my knowledge.

What we are trying to accomplish is the disctinction made by the go language data semantics: packages based on what it does (interfaces) instead of what it is (concrete types).

## Structure
```
.
├── apis
│   -> code that its related to solve the business logic/problem & that can be reusable.
│   -> code here is not very opinionated, we can define policy of usage as needed.
│   ├── auth: related to OPA.
│   ├── debug/metrics: custom metrics collector server, exposed with statsviz.
│   ├── errs: custom error handling.
│   ├── middleware: middleware handling logic.
│   └── mux: our custom mux that constructs & ties with all server routes.
│
├── apps
│   -> maintains the servers and tools built from this project, they output binaries.
│
├── libs
│   -> is like the stdlib of this project.
│   -> a restriction is that we can't have interdependency between them bc ideally they
│       can be deployed as a standalone pkg.
│
├── vendor
│   -> module for direct dependency, third party code.
│   -> it just gives me security of the code that it's being used.
│
├── zarf
│   -> these are the infrastructure or security related configuration files.
│   ├── docker
│   ├── k8s: kustomize, kind cluster management, deploys...
│   └── keys: pem files generated for authn/authz
│
├── makefile
│   -> scripts to facilitate the development of this project.
│   -> currently features starting up/down/restart, stress testing & binaries building.
```

## Features

1. Custom metrics server with Statsviz (apis/debug & apis/metrics).
2. AuthoriZation, AutheNtication and Token Generation with Open Policy Agent/Rego & JWT (apis/auth)
3. Clean startup/shutdown.
4. Flexible way of configuring the launch of application, facilitating the development/testing/debugging.
5. Clear (or the attempt to) separation of concerns, leveraging packages to communicate intent.