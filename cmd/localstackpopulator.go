package cmd

import (
	"log"

	"github.com/one2nc/cloud-lens/internal/config"
	pop "github.com/one2nc/cloud-lens/populator"
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

		errCB := pop.CreateBuckets(sess)
		if errCB != nil {
			log.Fatal("err: ", errCB)
		}
		errCEI := pop.CreateEC2Instances(sess)
		if errCEI != nil {
			log.Fatal("err: ", errCEI)
		}
	},
}

func init() {
	rootCmd.AddCommand(lspop)
}
