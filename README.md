This tool will automatically delete all CloudFormation stacks given a filter.


Example:

```shell
aws-delete-stacks --force mystack
```

This will delete all stacks that has a name that contains the string "mystack".
You can have a look at the stacks yourself at
[AWS CloudFormation](https://eu-central-1.console.aws.amazon.com/cloudformation/home?region=eu-central-1#/stacks) 
(remember to switch to your zone).

# Install

```shell
go install
```

# Run

Log in to AWS

* If you use [saml2aws](https://github.com/Versent/saml2aws)

```shell
saml2aws login --profile default
```

* Or

```shell
aws configure
```

Now, you can run 

```shell
aws-delete-stacks --help
```

# Details

Running

```shell
aws-delete-stacks somestack
```

is totally safe. Nothing is deleted until you specify the `--force` flag.
