package cmd

import (
	"fmt"

	list "github.com/jedib0t/go-pretty/v6/list"
	table "github.com/jedib0t/go-pretty/v6/table"
	text "github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
	pkg "github.com/vapusdata-oss/aistudio/cli/pkgs"
)

func NewExplainOps() *cobra.Command {

	cmd := &cobra.Command{
		Use:   pkg.ExplainOps,
		Short: "This command is to list all agents. You will have to provide the resource name as an argument to perform certain actions on that resource",
		Long: `
	You can use this command to list all the agents or get a specific resource.
			`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				cobra.CheckErr(fmt.Errorf("no resource provided for this command"))
			}
			getValidResources()
		},
	}

	return cmd
}

func getValidResources() {
	var xx string

	xx = text.FormatUpper.Apply("Goal Based Agents with list of their actions: ")
	xx = text.Underline.Sprintf(xx)
	vapusGlobals.logger.Info().Msgf("\n%v", xx)
	tw := pkg.NewTableWritter()
	tw.AppendHeader(table.Row{"Agents", "Actions", "Version", "Chaining Support"})
	for resource, operations := range vapusGlobals.AgentsActions {
		tw.AppendRow(table.Row{resource, pkg.NewListWritter(operations, list.StyleMarkdown).Render(), "v1alpha1", true})
		tw.AppendSeparator()
	}
	tw.Render()

	// xx = text.FormatUpper.Apply("Model Based Agents with list of their reflexes:")
	// xx = text.Underline.Sprintf(xx)
	// vapusGlobals.logger.Info().Msgf("\n%v", xx)
	// tw1 := pkg.NewTableWritter()
	// tw1.AppendHeader(table.Row{"Agents", "Reflexes", "Version", "Chaining Support"})
	// for resource, operations := range vapusGlobals.AgentsReflexes {
	// 	tw1.AppendRow(table.Row{resource, pkg.NewListWritter(operations, list.StyleMarkdown).Render(), "v1alpha1", true})
	// 	tw1.AppendSeparator()
	// }
	// tw1.Render()

	// xx = text.FormatUpper.Apply("Utility Agents with list of their utilities: ")
	// xx = text.Underline.Sprintf(xx)
	// vapusGlobals.logger.Info().Msgf("\n%v", xx)
	// tw2 := pkg.NewTableWritter()
	// tw2.AppendHeader(table.Row{"Agents", "Utilities", "Version", "Chaining Support"})
	// for resource, operations := range vapusGlobals.AgentsUtilities {
	// 	tw2.AppendRow(table.Row{resource, pkg.NewListWritter(operations, list.StyleMarkdown).Render(), "v1alpha1", true})
	// 	tw2.AppendSeparator()
	// }
	// tw2.Render()
}
