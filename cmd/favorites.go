package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"lawsearch/internal/store"
)

var favoritesLaw string

var favoritesCmd = &cobra.Command{
	Use:   "favorites",
	Short: "查看收藏的条款",
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := store.Open(dbPath)
		if err != nil {
			return err
		}
		defer s.Close()

		if favoritesLaw != "" {
			law, err := s.LawByName(favoritesLaw)
			if err != nil {
				return err
			}
			articles, err := s.Favorites(law.ID)
			if err != nil {
				return err
			}
			sort.Slice(articles, func(i, j int) bool {
				return articles[i].Number < articles[j].Number
			})
			table := tablewriter.NewWriter(os.Stdout)
			table.Header("条款", "内容")
			for _, article := range articles {
				_ = table.Append(fmt.Sprintf("第%d条", article.Number), article.Content)
			}
			fmt.Printf("%s - 收藏条款（共 %d 条）\n", law.Name, len(articles))
			return table.Render()
		}

		laws, err := s.Laws()
		if err != nil {
			return err
		}
		table := tablewriter.NewWriter(os.Stdout)
		table.Header("法规", "收藏数")
		total := 0
		for _, law := range laws {
			count, err := s.FavoriteCount(law.ID)
			if err != nil {
				continue
			}
			if count > 0 {
				total += count
				_ = table.Append(law.Name, fmt.Sprintf("%d", count))
			}
		}
		fmt.Printf("收藏总览（共 %d 条）\n", total)
		return table.Render()
	},
}

func init() {
	rootCmd.AddCommand(favoritesCmd)
	favoritesCmd.Flags().StringVar(&favoritesLaw, "law", "", "法规名称或 ID，查看该法规的收藏条款")
}
