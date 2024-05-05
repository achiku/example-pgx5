# example-pgx5

- login as database superuser

```
create role pgx_root;
alter role pgx_root with login;
alter user pgx_root with superuser;
create database pgx owner pgx_root;
```

- login as `pgx_root` 

```
create role pgx_api;
alter role pgx_api with login;
create schema pgx_api authorization pgx_api;
alter role pgx_api set search_path = pgx_api;

create role pgx_api_test;
alter role pgx_api_test with login;
alter role pgx_api_test set search_path = pgx_api_test;
```
