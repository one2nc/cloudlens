---
description: k9s like CLI for AWS
---

# 👋 Cloudlens

![](<.gitbook/assets/cloudlens (2).png>)

## Terminal-ate Your Cloud Worries!

AWS Console in your terminal! well, almost. Explore AWS services like EC2, S3, IAM, VPC, etc. from your terminal. If you like k9s for Kubernetes, you'll love cloudlens.

### Got 2 minutes? Check out a video overview of cloudlens:

{% embed url="https://drive.google.com/file/d/1Nc-b4g8F9F7V1ARi6LFo0vglbqKLOGZQ/view?usp=share_link" %}

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

| Services                            | Description                                                           | Alias   |
| ----------------------------------- | --------------------------------------------------------------------- | ------- |
| [S3](./#s3)                         | View all S3 buckets and their contents                                | `s3`    |
| [EC2](./#ec2)                       | view all instances and their associated metadata, including JSON data | `ec2`   |
| [EC2 Snapshot](./#ec2-snapshot)     | view a list of all EC2 snapshots                                      | `ec2:s` |
| [EC2 Image](./#ec2-image)           | See a list of all EC2 images                                          | `ec2:i` |
| [ECS Clusters](./#ecs-clusters)     | View all ECS Clusters                                                 | `ecs:c` |
| [VPC](./#vpc)                       | view all  VPC's and their associated metadata, including JSON data    | `vpc`   |
| [Security Group](./#security-group) | Security Groups and their associated metadata                         | `sg`    |
| [IAM users](./#iam-users)           | view all IAM users and their associated metadata                      | `iam:u` |
| [IAM Groups](./#iam-group)          | view all IAM users groups                                             | `iam:g` |
| [IAM Roles](./#iam-roles)           | view all IAM Roles                                                    | `iam:r` |
| [EBS](./#ebs)                       | View all available EBS volumes                                        | `ebs`   |
| [SQS](./#sqs)                       | view a list of all SQS queues                                         | `sqs`   |
| [Lambda](./#lambda)                 | view a list of all Lambda functions                                   | `lamda` |

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

### EC2

To view the <mark style="color:purple;">EC2</mark> page, enter the command `ec2` in your prompt. Pressing enter will display all available instances. To view specific details about an instance, select it and press enter to display a JSON with the EC2 information. You can use the escape key to navigate back to the previous page. Additionally, you can download a CSV file by using the `z` command.

<figure><img src=".gitbook/assets/image (7).png" alt="EC2 Page"><figcaption><p>EC2 Page</p></figcaption></figure>

<figure><img src=".gitbook/assets/image (32).png" alt="EC2 Details Page"><figcaption><p>EC2 Details Page</p></figcaption></figure>

### EC2 Snapshot

To access the <mark style="color:blue;">EC2 snapshots</mark> page, enter `ec2:s` in your prompt and press Enter. This will display a list of all available snapshots. Select a snapshot to view its details, then press Enter again. To return to the previous page, press the Escape key.

<figure><img src=".gitbook/assets/image (24).png" alt="EC2 Snapshot Page"><figcaption><p>EC2 Snapshot Page</p></figcaption></figure>

<figure><img src=".gitbook/assets/image (5).png" alt="EC2 Snapshot Details Page"><figcaption><p>EC2 Snapshot Details Page</p></figcaption></figure>

### EC2 Image

To access the <mark style="color:blue;">EC2 image</mark> page, enter `ec2:i` in your prompt and press Enter. This will display a list of all available images. Select an image to view its details, then press Enter again. To return to the previous page, press the Escape key.

<figure><img src=".gitbook/assets/image (3).png" alt="EC2 Image Page"><figcaption><p>EC2 Image Page</p></figcaption></figure>

<figure><img src=".gitbook/assets/image (21).png" alt="EC2 Image Details Page"><figcaption><p>EC2 Image Details Page</p></figcaption></figure>

### ECS Clusters

To acess the <mark style="color:blue;">ECS Clusters</mark> page, enter `ecs:c`  in your prompt and press Enter. This will display a list of all deployed Clusetrs.

### VPC

To access the <mark style="color:green;">VPC</mark> management functionality, type `vpc` in the command prompt to display a list of available VPCs. Selecting a specific VPC and pressing enter will show a JSON file with its information. You can view the VPC's subnets by using the  `s` command and navigate back to the previous page by pressing the escape key. Additionally, you can download a CSV file using the `z` command.

<figure><img src=".gitbook/assets/image (27).png" alt="VPC Page"><figcaption><p>VPC Page</p></figcaption></figure>

<figure><img src=".gitbook/assets/image (13).png" alt="VPC Details Page"><figcaption><p>VPC Details Page</p></figcaption></figure>

### Security Group

To access the <mark style="color:blue;">security group</mark> management functionality, enter `sg` in the command prompt to display a list of available security groups. Selecting a specific security group and pressing enter will display a JSON file with its information. To go back to the previous page, press the escape key. You can also download a CSV file using the `z` command.

<figure><img src=".gitbook/assets/image (16).png" alt="Security Group Page"><figcaption><p>Security Group Page</p></figcaption></figure>

<figure><img src=".gitbook/assets/image (2).png" alt="Security Group Details Page"><figcaption><p>Security Group Details Page</p></figcaption></figure>

### IAM Users

To view <mark style="color:yellow;">IAM users</mark>, enter `iam:u` in the command prompt to display a list of users. Selecting a specific user and pressing `Shift+P` will display their policy. To go back to the previous page, press the escape key. You can also download a CSV file using the `z` command.

<figure><img src=".gitbook/assets/image (14).png" alt="IAM Users"><figcaption><p>IAM Users</p></figcaption></figure>

<figure><img src=".gitbook/assets/image (19).png" alt="IAM User Policy"><figcaption><p>IAM User Policy</p></figcaption></figure>

### IAM Group

To view <mark style="color:purple;">IAM groups</mark>, enter `iam:g` in the command prompt to display a list of groups. Selecting a specific group and pressing Shift+P will display its users. To go back to the previous page, press the escape key. You can also download a CSV file using the `z` command.

<figure><img src=".gitbook/assets/image (20).png" alt="IAM Group"><figcaption><p>IAM Group</p></figcaption></figure>

<figure><img src=".gitbook/assets/image (11).png" alt="Group Policy"><figcaption><p>Group Policy</p></figcaption></figure>

### IAM Roles

To view <mark style="color:purple;">IAM roles</mark>, enter `iam:r` in the command prompt and press Enter to display a list of user roles. To view the policy of a specific role, select the role and press `Shift+P`. To return to the previous page, press the Escape key. To download a CSV file, use the `z` command

<figure><img src=".gitbook/assets/image (8).png" alt="IAM Role Page"><figcaption><p>IAM Role Page</p></figcaption></figure>

<figure><img src=".gitbook/assets/image (18).png" alt="IAM Role Policy Page"><figcaption><p>IAM Role Policy Page</p></figcaption></figure>

### EBS

To view <mark style="color:green;">EBS volumes</mark>, enter `ebs` in the command prompt and press Enter to show a list of volumes. Select a volume and press Enter to view its details. To return to the previous page, press the Escape key.

<figure><img src=".gitbook/assets/image (29).png" alt="EBS Page"><figcaption><p>EBS Page</p></figcaption></figure>

<figure><img src=".gitbook/assets/image (15).png" alt="EBS Details Page"><figcaption><p>EBS Details Page</p></figcaption></figure>

### SQS

To view <mark style="color:orange;">SQS queues</mark>, enter `sqs` in the command prompt and press Enter to display a list of queues. Select a queue and press Enter to view its details. To go back to the previous page, press the Escape key.

<figure><img src=".gitbook/assets/image.png" alt="SQS Page"><figcaption><p>SQS Page</p></figcaption></figure>

<figure><img src=".gitbook/assets/image (23).png" alt="SQS Details Page"><figcaption><p>SQS Details Page</p></figcaption></figure>

### Lambda

View all your Lambda functions easily by entering 'lambda' command in your terminal prompt.

<figure><img src=".gitbook/assets/image (12).png" alt=""><figcaption></figcaption></figure>
