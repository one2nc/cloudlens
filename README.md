<p align="center">
      <img src="assets/cloudlens.png" alt="Cloudlens" width="225" height="150" >
</p>

## Cloudlens - A Sophisticated Command Line Interface for Effortless AWS Service Management!

https://user-images.githubusercontent.com/87517248/220539848-f20b9791-29f9-4c09-b733-3fdf779d4911.mov

Your One-Stop Terminal Solution for Seamless AWS Service Management and Monitoring! With its intuitive UI and advanced features, it empowers you to effortlessly navigate, observe, and optimize your AWS environment, giving you more time to focus on your core business goals. Whether you're a seasoned cloud expert or just getting started, cloudlens will elevate your AWS experience and take your service management to new heights.

## Note
**Cloudlens reads your ~/.aws/config file, but it does not store or send your access and secret key anywhere. The access and secret key is used only to securely connect to AWS API via AWS SDK.**

**Since cloudlens is in readonly mode, we recommend you create an access and secret key that only has readonly permissions to the AWS services.**

## Features

Our terminal application supports various features, allowing users to view EC2 instances, S3 buckets, EBS volumes, VPCs, SQS queues, Lambda functions, subnets, security groups, and IAM roles, making their work faster and more efficient.

## Documentation

Please refer to our [Cloudlens documentation](https://one2n.gitbook.io/docs/) to know more.

## We appreciate [localstack's](https://localstack.cloud/) assistance with the development

<img src="assets/localstack.jpeg" alt="k9s">


## Building From Source

 Cloudlens is currently using go v0.0.1 . In order to build cloudlens from source you must:

 1. Clone the repo
 2. Build and run the executable

      ```shell
      make run
      ```
## Installation

Cloudlens is available on Linux and macOS.

* Via [Homebrew](https://brew.sh/) for macOS or Linux

   ```shell
   brew install one2nc/cloudlens/cloudlens
   ```

## Key Bindings

Cloudlens uses aliases to navigate most AWS Services.

| Action                                                         | Command                       | Comment                                                                |
|----------------------------------------------------------------|-------------------------------|------------------------------------------------------------------------|
| Show active keyboard mnemonics and help                        | `?`                           |                                                                        |                                                                      |
| To bail out of cloudlens                                             | `:q`, `ctrl-c`                |                                                                        |
| Bails out of view/command/filter mode                          | `<esc>`                       |                                                                        |
| To view and switch to another AWS Service               | `:`ctx⏎                       |                                                                        |


## Screenshots

1. EC2
      <img src="assets/ec2.png"/>
1. EC2 Details
      <img src="assets/Ec2Json.png"/>

2. S3
      <img src="assets/s3.png"/>
2. S3 Details
      <img src="assets/s3Details.png"/>


## Acknowledgements

We would like to express our sincere appreciation to `K9s` as it has been a invaluable source of reference for this project.


All materials licensed under [Apache v2.0](http://www.apache.org/licenses/LICENSE-2.0)
