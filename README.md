# czar

## Setup:

* Must have AWS environment variables set or configuration file.

## Install Options:

**MAC**

`curl -O --insecure -o /usr/local/bin/czar https://github.com/feelobot/czar/releases/download/0.0.3/czar-macosx && sudo chmod +x /usr/local/bin/czar`


## LS

```
➜  czar git:(master) ./czar ls --tag "application" --value "external-ingress"
external-ingress: i-e53a6166 ec2-52-87-193-89.compute-1.amazonaws.com 172.21.11.148
external-ingress: i-ecde1176 ec2-52-91-238-92.compute-1.amazonaws.com 172.16.13.218
external-ingress: i-02144685 ec2-52-201-237-57.compute-1.amazonaws.com 172.16.12.173
```

**wildcards** 

```
➜  czar git:(master) ✗ ./czar ls -t "application" -v "*ingress*"
external-ingress-prod: i-f2a5df75 ec2-107-23-130-139.compute-1.amazonaws.com 172.16.12.245
external-ingress-prod: i-f0a5df77 ec2-107-23-130-198.compute-1.amazonaws.com 172.16.12.246
external-ingress-prod: i-dc9d6046 ec2-54-175-100-186.compute-1.amazonaws.com 172.16.13.235
external-ingress-prod: i-dd9d6047 ec2-54-172-82-223.compute-1.amazonaws.com 172.16.13.234
external-ingress-prod: i-fa9d6060 ec2-52-90-137-214.compute-1.amazonaws.com 172.16.13.233
```

## ssh

```
./czar ssh -u "core" -t "Name" -v "prod_kube*" "df -h"
```
