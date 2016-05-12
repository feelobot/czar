# czar

## Setup:

* Must have AWS environment variables set or configuration file.

## Install Options:

**MAC**

`curl -O --insecure -o /usr/local/bin/czar https://github.com/feelobot/czar/releases/download/0.0.3/czar-macosx && sudo chmod +x /usr/local/bin/czar`


## LS

```
➜ czar ls --tag "application" --value "external-ingress"
```


**wildcards** 

```
➜ czar ls -t "application" -v "*ingress*"
```

## ssh

```
➜ czar ssh -u "core" -t "Name" -v "prod_kube*" "df -h"
```
