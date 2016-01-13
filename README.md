# sgctl
Add and remove security groups to an EC2 instance. This works offline by providing the EC2 instance ID and AWS_REGION and on the instance itself:

On the instance:

```
$ sgctl add sg-123456 sg-678901
$ sgctl del sg-123456 sg-678901
```

From outside the instance:

```
$ sgctl add sg-123456 sg-678901 -i i-asdfgh
$ sgctl del sg-123456 sg-678901 -i i-asdfgh
```

Make sure the EC2 instance has `ec2:ModifyInstanceAttribute` permission.

