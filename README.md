# Remote Exec
Simple go package to exec command remotely.

### Getting Started
```bash
go get -u github.com/ahmadmuzakki29/rexec
```

### Usage
```bash
rexec [-h <hosts>|-e] <command>

# example of executing command remotely

# single server
rexec -h root@192.168.100.160 tail -f /var/log/nginx/access.log

# multiple server
rexec -h root@192.168.100.160,root@192.168.100.161 zgrep 500 /var/log/nginx/access.log.1.gz

# using file config
rexec grep 500 /var/log/nginx/access.log

# edit file config
rexec -e
```

### Thanks
- https://github.com/alileza/rtail