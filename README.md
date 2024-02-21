
# DeVerse Main Server



Require go version >= 1.14, if your go version is lower, please use legacy branches,
there are quite a lot of incompatible changes between master and legacy branches.

```
# start a web server listening on 0.0.0.0:8080
go run main.go
```

## Pre-required (Window User):
1. Install Go version 1.21
2. Install GCC, can follow [this](https://www.msys2.org/)
3. If there is any Error: 

   ```go build -o c:\Data\master-server__debug_bin3967034439.exe -gcflags all=-N -l . github.com/ethereum/go-ethereum/crypto/secp256k1: build constraints exclude all Go files in C:\Users\User\go\pkg\mod\github.com\ethereum\go-ethereum@v1.8.10\crypto\secp256k1 (exit status 1)```  
    Just run this command to update the  go-ethereum to the latest version`go get github.com/ethereum/go-ethereum@latest`

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

