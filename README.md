# Remote Exec
Simple go package to exec command remotely.

### Getting Started
```bash
go get -u github.com/ahmadmuzakki29/rexec
```

### Usage
```bash
rexec [-h <hosts>|-e] <command>

# example of executing ls command remotely

# single server
rexec -h root@192.168.100.160 ls

# multiple server
rexec -h root@192.168.100.160,root@192.168.100.161 ls

# using file config
rexec ls

# edit file config
rexec -e ls
```

### Thanks
- https://github.com/alileza/rtail