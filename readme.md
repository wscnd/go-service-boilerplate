# Go Service Boilerplate

Just very opinionated boilerplate, it has these layers:

## Structure
```
.
├── app
│   -> maintains all binaries, services, tooling & things
│      that need to be deployed together to work on this project.
│   -> job: startup, receive external input, process in business layer, send to output, ~clean shutdown~.
│
├── business
│   -> code that its related to solve the business logic/problelm & that can be reusable.
│   ├── core: core business problem
│   ├── data: data management, order, paging, db, transaction...
│   └── web: reusable stuff related to webapis.
│
├── foundation
│   -> is like the stdlib of this project.
│   -> can't have interdependency between them.
│   -> ideally can be deployed as a standalone pkg.
│
├── vendor
│   -> module for direct dependency, third party code.
│
├── zarf
│   -> configs
│   ├── docker
│   └── k8s: kustomize, cluster stuff...
│
├── makefile
│   -> scripts
```

## Goals

What we are trying to accomplish is the disctinction made by the go language data semantics: packages based on what it does (interfaces) instead of what it is (concrete types).