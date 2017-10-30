# cde-client
This is a command line interface for cde.

## Compilation
```
make build
```

## Usage

### How to build a stack

```
cde register http://192.168.50.4:31088 --email stackmanager@tw.com --password admin  
cde whoami  
cde stacks:create javajersey-test <stack-definition-file.yml>  
```

### How to use a stack

```
cde register http://192.168.50.4:31088 --email=gm@tw.com * --password=111111  
cde keys:add ~/.ssh/id_rsa.pub  
cde stacks:list  
cde apps:create ketsu javajersey-test  
```

