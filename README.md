# DeVerse Main Server



Require go version >= 1.14, if your go version is lower, please use legacy branches,
there are quite a lot of incompatible changes between master and legacy branches.

```
# start a web server listening on 0.0.0.0:8080
go run main.go
```

## Components

- Framework: [gin-gonic/gin](https://github.com/gin-gonic/gin)
- Database ORM: [go-gorm/gorm](https://github.com/go-gorm/gorm)
- Database migration: [rubenv/sql-migrate](https://github.com/rubenv/sql-migrate)
- Zero Allocation JSON Logger: [rs/zerolog](https://github.com/rs/zerolog)
- YAML support: [go-yaml/yaml](https://github.com/go-yaml/yaml)
- Testing toolkit: [stretchr/testify](https://github.com/stretchr/testify)

## Configuration

Edit the `config.yml` with your own config

## Main DB Schema:
https://dbdiagram.io/d/619a007e02cf5d186b606f44

## Database Migration

**Setup Database**
1/ install mysql https://flaviocopes.com/mysql-how-to-install/
2/ Change in config.yml the db password
3/ Install TablePlus as GUI for db

**Create the database first**

```
CREATE DATABASE IF NOT EXISTS `gin` DEFAULT CHARACTER SET utf8mb4;
```

**Migrates the database to the most recent version available**
```
# Include go path to main path to access go libs executables
e.g export PATH="$PATH:/Users/hieuletrung/go/bin"
```

```
./migrate.sh up
```

**Undo a database migration**

```
./migrate.sh down
```

**Show migration status**

```
./migrate.sh status
```

**Create a new migration**

```
./migrate.sh new a_new_migration
```
# main-server
