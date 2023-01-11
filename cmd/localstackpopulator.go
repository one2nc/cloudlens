package cmd

import (
	"fmt"
	"log"

	"github.com/one2nc/cloud-lens/internal/config"
	"github.com/spf13/cobra"
)

var lspop = &cobra.Command{
	Use:   `lspop`,
	Short: ``,
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Get()
		if err != nil {
			log.Fatal("err: ", err)
		}
		sess, err := config.GetSession("test", "ap-south-1", cfg.AwsConfig)
		if err != nil {
			log.Fatal("err: ", err)
		}

		bucketList := CreateAndListS3Buckets(sess)
		fmt.Println("Bucket list is:", bucketList)

		ec2InstanceInfo := CreateAndListEC2Instances(sess)
		fmt.Println("Ec2 instance info is:", ec2InstanceInfo)
	},
}

func init() {
	rootCmd.AddCommand(lspop)
}
