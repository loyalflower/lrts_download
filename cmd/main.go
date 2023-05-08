package cmd

import (
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "download",
		Short: "lrts 音频下载",
		RunE:  download,
	}
	rootCmd.Execute()
}

func download(cmd *cobra.Command, args []string) error {
	return nil
}
