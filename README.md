# GitStore

A simple git server written in golang

## Design

1. Keep the design simple initially don't add too many dependencies. For now not using any framework this might change in future.
2. Since we don't even have the business logic, don't even think about writing production ready code.
3. Use as little files as needed.

## Tasks

- [x] Add a route to be able to create a repository.
- [ ] Add a route to pull the repository.
- [ ] Add a route to push the code to remote repository.
