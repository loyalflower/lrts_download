package main

import (
	"github.com/cockroachdb/errors"
	"github.com/loyalflower/lrts_download"
	"github.com/spf13/cobra"
	"strconv"
)

var (
	rootCmd = &cobra.Command{
		Use: "lrts",
	}

	token = ""
)

func main() {
	downloadCmd := &cobra.Command{
		Use:   "download",
		Short: "lrts 音频下载",
		Example: `  lrts download book_id
`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("requires at least one arg")
			}
			v, e := strconv.Atoi(args[0])
			if e != nil || v < 1 {
				return errors.New("requires a number")
			}
			return nil
		},
		RunE: download,
	}
	downloadCmd.Flags().StringVarP(&token, "token", "t", "", "use special token")

	rootCmd.AddCommand(downloadCmd)
	rootCmd.Execute()
}

func download(cmd *cobra.Command, args []string) error {
	options := lrts_download.DownloadOptions{
		Token: token,
	}
	return lrts_download.Download(args[0], options)
}
