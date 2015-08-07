# czar

## Install Options

* Add binary file to your /usr/local/bin
* Run build-all script to generate a binary for your machine (requires docker)

## Usage
![help](http://i.imgur.com/CBOYuhm.png)

* exec

**Example:**

```
czar exec -p "/Users/felix.rodriguez/.ssh/aws-key.pem" -u ubuntu --tag Name --value "*spark*" "echo hello"
```

**Output:**
```
> Number of reservation sets:  2
  > Number of instances:  1
    - Instance ID:  i-521e7282
    - DNS Name:  ec2-xx-xxx-xxx-xxx.compute-1.amazonaws.com
    - Tag:  spark-datacrunch-process-worker03
[golang-sh]$ eval `ssh-agent`
[golang-sh]$ ssh-add /Users/felix.rodriguez/.ssh/aws-key.pem
Identity added: /Users/felix.rodriguez/.ssh/aws-key.pem (/Users/felix.rodriguez/.ssh/aws-key.pem)
[golang-sh]$ ssh -o StrictHostKeyChecking=no -i /Users/felix.rodriguez/.ssh/aws-key.pem ubuntu@ec2-xx-xxx-xxx-xxx.compute-1.amazonaws.com echo hello
hello
  > Number of instances:  1
    - Instance ID:  i-7f4abad6
    - DNS Name:  ec2-xx-xx-xxx-xx.compute-1.amazonaws.com
    - Tag:  spark-datacrunch-process-worker04
[golang-sh]$ eval `ssh-agent`
[golang-sh]$ ssh-add /Users/felix.rodriguez/.ssh/aws-key.pem
Identity added: /Users/felix.rodriguez/.ssh/aws-key.pem (/Users/felix.rodriguez/.ssh/aws-key.pem)
[golang-sh]$ ssh -o StrictHostKeyChecking=no -i /Users/felix.rodriguez/.ssh/aws-key.pem ubuntu@ec2-xx-xxx-xxx-xxx.compute-1.amazonaws.com echo hello
hello
```

* list
todo

* ssh
todo
