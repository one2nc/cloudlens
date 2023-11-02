<p align="center">
      <img src="assets/cloudlens.png" alt="Cloudlens" width="225" height="150" >
</p>

## Cloudlens - k9s like CLI for AWS. 

![](./assets/cloudlensdemo.gif)

AWS Console in your terminal! well, almost. Explore AWS services like EC2, S3, IAM, VPC, etc. from your terminal. If you like k9s for Kubernetes, you'll love cloudlens.

## Installation

* Via [Homebrew](https://brew.sh/) for macOS or Linux

   ```shell
   brew install one2nc/cloudlens/cloudlens
   ```
* Using go install
      *cloudlens* requires go1.19 to install successfully. Run the following command to install the latest version -
   ```shell
       go install -v github.com/one2nc/cloudlens@latest
   ```

* Building from source code
      Cloudlens is currently in active development. We use Go 1.19. Follow these steps to build cloudlens locally:

       1. Clone the repo
       2. Build and run the executable

  To Run:
  ```shell
  make run
  ```

## Prerequisite
1. Docker installed on your local. Refer this [documentation](https://docs.docker.com/engine/install/)
2. If you want to use localstack for populating dummy data, use our repo [cloud-lens-populator](https://github.com/one2nc/cloud-lens-populator) 

## Usage

For the simple usage, just run the command without any options.

```shell
cloudlens
```

For knowing all the options available, use:
```shell
cloudlens help
```

### Self Update
For updating to latest version, use:
```console
cloudlens update
```
### Using Localstack
- Configure localstack server to listen on port `4566`. 
-  Use our repo [cloud-lens-populator](https://github.com/one2nc/cloud-lens-populator) to setup and populate dummy data.
- To run cloudlens with localstack use aws sub-command with `-l` or `--localstack` flag 
```console
cloudlens aws --localstack
```

## Features

Cloudlens supports viewing EC2 instances, S3 buckets, EBS volumes, VPCs, SQS queues, Lambda functions, Subnets, Security Groups, and IAM roles. Read the [cloudlens documentation](https://one2n.gitbook.io/docs/) to know more.

## Screenshots

1. EC2
      <img src="assets/ec2.png"/>
1. EC2 Details
      <img src="assets/ec2Details.png"/>

2. S3
      <img src="assets/s3.png"/>
2. S3 Details
      <img src="assets/s3Details.png"/>

## Documentation

Please refer to our [cloudlens documentation](https://one2n.gitbook.io/docs/) to know more.


## Key Bindings

Cloudlens uses k9s like shortcuts for navigation. Listed below are few of the shortcuts:

| **Action**                                | **Command**   |
|-------------------------------------------|---------------|
| Show active keyboard mnemonics and help   | ?             |
| To bail out of cloudlens                  | :q ,   ctrl-c |
| Bails out of view/command/filter mode     | esc         |
| To view and switch to another AWS Service | :S3/EC2/VPC⏎  |

## Note
**Cloudlens reads your ~/.aws/config file, but it does not store or send your access and secret key anywhere. The access and secret key is used only to securely connect to AWS API via AWS SDK.**

**Since cloudlens is in readonly mode, we recommend you create an access and secret key that only has readonly permissions to the AWS services.**

## Acknowledgements

We would like to express our sincere appreciation to `K9s` as it has been a invaluable source of reference for this project.

All materials licensed under [Apache v2.0](http://www.apache.org/licenses/LICENSE-2.0)

## We appreciate [localstack's](https://localstack.cloud/) assistance with the development

<img src="assets/localstack.jpeg" alt="localstack">
