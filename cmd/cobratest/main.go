package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:  "example",
	Long: "コマンドの説明",
	// メイン処理
	RunE: func(cmd *cobra.Command, args []string) error {
		// フラグを取得
		t, err := cmd.Flags().GetBool("toggle")
		if err != nil {
			return err
		}

		fmt.Printf("toggle: %t\n", t)
		return nil
	},
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// フラグの追加
	// bool 型の --toggle, -t フラグを設定
	rootCmd.Flags().BoolP("toggle", "t", false, "フラグの説明")
}
