---
description: >-
  Your One-Stop Terminal Solution for Seamless AWS Service Management and
  Monitoring!
---

# 👋 Cloudlens

![](<.gitbook/assets/cloudlens (2).png>)

## Terminal-ate Your Cloud Worries!

Cloudlens provides a terminal UI for AWS Service Management and Monitoring! With its easy to use UI and advanced features, you can easily navigate, observe, and optimize your AWS environment. Whether you're a cloud expert or just getting started, cloudlens will elevate your AWS experience and take your service management to new heights.

### Got 2 minutes? Check out a video overview of cloudlens:

Add video

### Getting set up

<details>

<summary>Install</summary>

```
brew install one2nc/cloudlens/cloudlens
```

</details>

<details>

<summary>Run</summary>

```
cloudlens
```

</details>

### Features

| Services                             | Description                                                           | Alias   |
| ------------------------------------ | --------------------------------------------------------------------- | ------- |
| [S3](./#s3)                          | View all S3 buckets and their contents                                | `s3`    |
| [EC2](./#ec2)                        | view all instances and their associated metadata, including JSON data | `ec2`   |
| [EC2 Snapshot](./#ec2-snapshot)      | view a list of all EC2 snapshots                                      | `ec2:s` |
| [EC2 Image](./#ec2-image)            | See a list of all EC2 images                                          | `ec2:i` |
| [VPC](./#vpc)                        | view all  VPC's and their associated metadata, including JSON data    | `vpc`   |
| [Security Group](./#security-group)  | Security Groups and their associated metadata                         | `sg`    |
| [IAM users](./#iam-users)            | view all IAM users and their associated metadata                      | `iam:u` |
| [IAM user groups](./#iam-user-group) | view all IAM users groups                                             | `iam:g` |
| [IAM Roles](./#iam-roles)            | view all IAM Roles                                                    | `iam:r` |
| [EBS](./#ebs)                        | View all available EBS volumes                                        | `ebs`   |
| [SQS](./#sqs)                        | view a list of all SQS queues                                         | `sqs`   |
| [Lambda](./#lambda)                  | view a list of all Lambda functions                                   | `lamda` |

### Prompt

To view the input prompt, press the  :  key. From there, you can try different commands to access and view various services.You can use the tab key or the right arrow key for autocomplete to make entering commands faster and easier.

<figure><img src=".gitbook/assets/image (4).png" alt="Prompt"><figcaption><p>Prompt</p></figcaption></figure>

### Help

To access the Help Page, press the "?" key while on the prompt. This will display all available commands and services for you to view.

<figure><img src=".gitbook/assets/image (1).png" alt="Help Page"><figcaption><p>Help Page</p></figcaption></figure>

### Resources Tab

The resources window located in the top right corner of the terminal provides you with all the necessary commands to browse for a specific service.

<figure><img src=".gitbook/assets/image (10).png" alt="Resources Tab"><figcaption><p>Resources Tab</p></figcaption></figure>

### Dropdowns

You can switch between the dropdown options using the tab button.

<figure><img src=".gitbook/assets/image (28).png" alt="Dropdowns"><figcaption><p>Dropdowns</p></figcaption></figure>

### S3

To view the <mark style="color:orange;">S3</mark> page, use the command `s3` in your prompt. Pressing enter will display all the available buckets, folders, and files. You can use the escape key to go back to the previous page.. Additionally, you can download a CSV file using the `z` command.

<figure><img src=".gitbook/assets/image (30).png" alt="S3 Page"><figcaption><p>S3 Page</p></figcaption></figure>

<figure><img src=".gitbook/assets/image (22).png" alt="S3 Details Page"><figcaption><p>S3 Details Page</p></figcaption></figure>

{% content-ref url="overview/features.md" %}
[features.md](overview/features.md)
{% endcontent-ref %}

## Get Started

We've put together some helpful guides for you to get setup with cloudlens quickly and easily.

{% content-ref url="overview/getting-set-up.md" %}
[getting-set-up.md](overview/getting-set-up.md)
{% endcontent-ref %}
