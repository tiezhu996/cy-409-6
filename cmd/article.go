package cmd

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"lawsearch/internal/store"
)

var articleLaw string
var articleNumber int
var articleFrom int
var articleTo int
var articleFavorite bool
var articleUnfavorite bool

var articleCmd = &cobra.Command{
	Use:   "article",
	Short: "按条款号查询",
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := store.Open(dbPath)
		if err != nil {
			return err
		}
		defer s.Close()
		law, err := s.LawByName(articleLaw)
		if err != nil {
			return err
		}

		var targetNumbers []int
		rangeMode := false
		if articleFrom > 0 && articleTo >= articleFrom {
			rangeMode = true
			for i := articleFrom; i <= articleTo; i++ {
				targetNumbers = append(targetNumbers, i)
			}
		} else {
			targetNumbers = append(targetNumbers, articleNumber)
		}

		if articleFavorite {
			for _, n := range targetNumbers {
				if err := s.AddFavorite(law.ID, n); err != nil {
					if rangeMode {
						continue
					}
					return err
				}
			}
		} else if articleUnfavorite {
			for _, n := range targetNumbers {
				if err := s.RemoveFavorite(law.ID, n); err != nil {
					if rangeMode {
						continue
					}
					return err
				}
			}
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.Header("收藏", "条款", "内容")
		for _, n := range targetNumbers {
			article, err := s.ArticleByNumber(law.ID, n)
			if err != nil {
				continue
			}
			fav, _ := s.IsFavorite(law.ID, article.Number)
			favMark := " "
			if fav {
				favMark = "★"
			}
			_ = table.Append(favMark, fmt.Sprintf("第%d条", article.Number), article.Content)
		}
		return table.Render()
	},
}

func init() {
	rootCmd.AddCommand(articleCmd)
	articleCmd.Flags().StringVar(&articleLaw, "law", "", "法规名称或 ID")
	articleCmd.Flags().IntVar(&articleNumber, "number", 0, "条款号")
	articleCmd.Flags().IntVar(&articleFrom, "from", 0, "范围起始条款号")
	articleCmd.Flags().IntVar(&articleTo, "to", 0, "范围结束条款号")
	articleCmd.Flags().BoolVar(&articleFavorite, "favorite", false, "标记收藏")
	articleCmd.Flags().BoolVar(&articleUnfavorite, "unfavorite", false, "取消收藏")
	_ = articleCmd.MarkFlagRequired("law")
}
