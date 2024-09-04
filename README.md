# Practical DDD Example #

[It is based on the ddd-gin-admin template](https://github.com/linzhengen/ddd-gin-admin)

- The entrypoint is: internal/app/app.go
- I remooved the "wire" framework and do dependency injection manually
- I added missing casbin related config files
- I added postgresql/mongodb to the infrastructure layer (dependency injection, repository pattern)
- There is some attempt on how to add different login providers (dependency injection)