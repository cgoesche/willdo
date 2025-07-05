package edit

import (
	"fmt"
	"strings"

	"github.com/cgoesche/willdo/internal/database"
	"github.com/cgoesche/willdo/internal/modules/category"
	"github.com/spf13/cobra"
)

var (
	catID          int64
	catName        string
	currentName    string
	catDescription string

	categoryCmd = &cobra.Command{
		Use:   "category",
		Short: "Edit a category",
		RunE: func(cmd *cobra.Command, args []string) error {
			db := database.New(conf.Database)
			catService := category.NewService(db)

			cats, err := catService.GetAll()
			if err != nil {
				return fmt.Errorf("failed to find any categories in the database, %v", err)
			}

			if len(strings.TrimSpace(currentName)) != 0 {
				catID = category.GetCategoryIDFromName(cats, currentName)
				if catID == 0 {
					return fmt.Errorf("no category found with name '%s'", currentName)
				}
			}

			cat, err := catService.GetById(catID)
			if err != nil {
				return fmt.Errorf("could not find a category with ID %d", catID)
			}

			if len(strings.TrimSpace(catName)) != 0 {
				cat.Name = catName
			}
			if catDescription != "" {
				cat.Description = catDescription
			}

			id, err := catService.Update(cat)
			if err != nil {
				return fmt.Errorf("failed to update category %d, %v", id, err)
			}
			fmt.Printf("Category updated!\n")
			return nil
		},
	}
)

func init() {
	categoryCmd.Flags().StringVarP(&currentName, "current", "c", "", "current category name")
	categoryCmd.Flags().StringVarP(&catName, "name", "n", "", "new category name")
	categoryCmd.Flags().StringVarP(&catDescription, "description", "d", "", "new category description")
	categoryCmd.MarkFlagsOneRequired("current")
	categoryCmd.MarkFlagsRequiredTogether("current", "name", "description")
}
