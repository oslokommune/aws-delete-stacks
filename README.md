This tool will automatically delete all CloudFormation stacks given a filter.

Motivation: Instead of using the AWS web console to click, refresh, delete, rinse and repeat, you can just use this
tool.

Example:

```shell
$ aws-delete-stacks mystack

- Would delete 3 stack(s)
eksctl-mystack-dev-cluster (19 Mar 21 12:17 UTC)
myproject-vpc-mystack-dev (19 Mar 21 12:13 UTC)
myproject-domain-mystack-dev (19 Mar 21 12:11 UTC)

No CloudFormation stacks were deleted, as you didn't specify the --force flag.
```

This will delete all stacks that has a name that contains the string "mystack".
You can have a look at the stacks yourself at
[AWS CloudFormation](https://eu-central-1.console.aws.amazon.com/cloudformation/home?region=eu-central-1#/stacks) 
(remember to switch to your zone).

See more examples below.

# Install

```shell
curl --silent --location "https://github.com/oslokommune/aws-delete-stacks/releases/latest/download/aws-delete-stacks_$(uname -s)_x86_64.tar.gz" | tar xz -C /tmp
sudo mv /tmp/aws-delete-stacks /usr/local/bin
```

or

```shell
git clone git@github.com:oslokommune/aws-delete-stacks.git
go install
```

or download directly binary from https://github.com/oslokommune/aws-delete-stacks/releases.

# Run

## Login to AWS

There are numerous ways to that is documented
[here](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html), but here is a straight forward way
assuming you use [saml2aws](https://github.com/Versent/saml2aws) and the AWS profile called 'default'.

1. Setup AWS profile defaults
  
Create a `~/.aws/config` that contains the following contens: (you can do this this by running `aws configure`)

```
[default]
region = eu-west-1 # or whatever region you like
```

2. Log in to AWS

```shell
saml2aws login --profile default
```

## Run 

```shell
aws-delete-stacks --help
```

Examples:

Delete all stacks that contains the name 'mystack':

```shell
aws-delete-stacks mystack
```

Delete all stacks that contains the name 'mystack', minus the stacks that contains 'mystack-hostedzone':

```shell
aws-delete-stacks mystack --exclude mystack-hostedzone
```

This is totally safe. Nothing is deleted until you specify the `--force` flag.

When you do, know that stuff in AWS will be deleted you do this. I.e. know what you are doing.

Disclaimer: I don't take responsibility for anything this tool does or your usage of it.

# Syntax

```shell
Delete AWS cloudformation stacks with names containing the string INCLUDE FILTER (minus stack containing exclude filter).

Usage:
  delete <INCLUDE FILTER> [flags]

Flags:
  -e, --exclude string   Set filter for which stacks to subtract from included results (filter method: string contains).
  -f, --force            Use this flag to actually delete stacks. Otherwise nothing is deleted.
  -h, --help             help for delete
```