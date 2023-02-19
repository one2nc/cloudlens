<img src="assets/cloud-lens.png" alt="k9s">

## Cloud-Lens - A Sophisticated Command Line Interface for Effortless AWS Service Management!

Your One-Stop Terminal Solution for Seamless AWS Service Management and Monitoring! With its intuitive UI and advanced features, it empowers you to effortlessly navigate, observe, and optimize your AWS environment, giving you more time to focus on your core business goals. Whether you're a seasoned cloud expert or just getting started, Cloud-Lens will elevate your AWS experience and take your service management to new heights.

## Building From Source

 Cloud-Lens is currently using go v0.0.1 . In order to build Cloud-Lens from source you must:

 1. Clone the repo
 2. Build and run the executable

      ```shell
      make run
      ```

## Key Bindings

Cloud-Lens uses aliases to navigate most AWS Services.

| Action                                                         | Command                       | Comment                                                                |
|----------------------------------------------------------------|-------------------------------|------------------------------------------------------------------------|
| Show active keyboard mnemonics and help                        | `?`                           |                                                                        |                                                                      |
| To bail out of Cloud-Lens                                             | `:q`, `ctrl-c`                |                                                                        |
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

3. IAM:G
      <img src="assets/iamg.png"/>
3. IAM:G Details
      <img src="assets/iamg-details.png"/>

4. IAM:U
      <img src="assets/iamu.png"/>
4. IAM:U Details
      <img src="assets/iamu-details.png"/>

5. IAM:R
      <img src="assets/iamr.png"/>
5. IAM:R Details
      <img src="assets/iamr-details.png"/>

6. SQS
      <img src="assets/sqs.png"/>
6. SQS Details:
       <img src="assets/sqs-details.png"/>

7. Lambda: 
      <img src="assets/lambda.png"/>

8. VPC:
      <img src="assets/vpc.png"/>
8. VPC Details:
      <img src="assets/vpc-details.png"/>


## Acknowledgements

We would like to express our sincere appreciation to `K9s` as it has been a invaluable source of reference for this project.