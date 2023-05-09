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

	token  = ""
	output = ""
)

func main() {
	cmdDownload()
	cmdSearch()
	cmdDetail()
	rootCmd.Execute()
}
func cmdDownload() {
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
	downloadCmd.Flags().StringVarP(&output, "output", "o", "", "output directory")

	rootCmd.AddCommand(downloadCmd)
}

func download(cmd *cobra.Command, args []string) error {
	options := lrts_download.DownloadOptions{
		Output: output,
	}
	options.Token = token
	return lrts_download.Download(args[0], options)
}

func cmdSearch() {
	searchCmd := &cobra.Command{
		Use:   "search",
		Short: "lrts 音频搜索",
		Example: `  lrts search keyword
`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("requires at least one arg")
			}
			return nil
		},
		RunE: search,
	}
	searchCmd.Flags().StringVarP(&token, "token", "t", "", "use special token")

	rootCmd.AddCommand(searchCmd)
}

func search(cmd *cobra.Command, args []string) error {
	options := lrts_download.SearchOptions{}
	return lrts_download.Search(args[0], options)
}

func cmdDetail() {
	detailCmd := &cobra.Command{
		Use:   "detail",
		Short: "lrts 音频详情",
		Example: `  lrts detail book_id
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
		RunE: detail,
	}
	detailCmd.Flags().StringVarP(&token, "token", "t", "", "use special token")

	rootCmd.AddCommand(detailCmd)
}

func detail(cmd *cobra.Command, args []string) error {
	options := lrts_download.DetailOptions{}
	return lrts_download.Detail(args[0], options)
}
