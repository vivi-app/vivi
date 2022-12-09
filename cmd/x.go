package cmd

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/vivi-app/vivi/color"
	"github.com/vivi-app/vivi/extensions/extension"
	"github.com/vivi-app/vivi/extensions/manager"
	"github.com/vivi-app/vivi/icon"
	"github.com/vivi-app/vivi/style"
	"github.com/vivi-app/vivi/text"
	"regexp"
)

func init() {
	rootCmd.AddCommand(xCmd)
}

var xCmd = &cobra.Command{
	Use:   "x",
	Short: "Manage extensions",
	Args:  cobra.NoArgs,
}

func init() {
	xCmd.AddCommand(xNewCmd)

	xNewCmd.Flags().BoolP("print-path", "p", false, "print path")
}

var xNewCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new extension",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		ext, err := extension.GenerateInteractive()
		handleErr(err)

		handleErr(ext.LoadPassport())

		fmt.Printf(
			"%s Created %s extension\n",
			style.Fg(color.Green)(icon.Check),
			style.Fg(color.Purple)(ext.String()),
		)

		if printPath := lo.Must(cmd.Flags().GetBool("print-path")); printPath {
			fmt.Println(ext.Path())
		}
	},
}

func init() {
	xCmd.AddCommand(xListCmd)
}

var xListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List installed extensions",
	Aliases: []string{"ls"},
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		extensions, err := manager.InstalledExtensions()
		handleErr(err)

		var (
			byAuthor = make(map[string][]*extension.Extension)
			authors  = make(map[string]any, 0)
		)

		for _, ext := range extensions {
			byAuthor[ext.Author()] = append(byAuthor[ext.Author()], ext)
			authors[ext.Author()] = nil
		}

		printForDomain := func(author string) {
			fmt.Println(
				style.
					New().
					Foreground(color.Yellow).
					Bold(true).
					Render(author),
			)

			for _, e := range byAuthor[author] {
				fmt.Printf(
					"%s %s %s\n",
					style.Fg(color.Purple)(e.Passport().Name),
					style.Bold(e.Passport().Version().String()),
					style.Faint(e.Passport().About),
				)
			}
		}

		for author, _ := range authors {
			if _, ok := byAuthor[author]; ok {
				printForDomain(author)
				fmt.Println()
			}
		}

		fmt.Printf(
			"%s %s installed\n",
			style.Fg(color.Pink)(icon.Heart),
			text.Quantify(len(extensions), "extension", "extensions"),
		)
	},
}

func init() {
	xCmd.AddCommand(xUninstallCmd)

	xUninstallCmd.Flags().StringP("id", "i", "", "id of the extension to remove")
}

var xUninstallCmd = &cobra.Command{
	Use:     "del",
	Short:   "Uninstall an extension",
	Aliases: []string{"rm", "remove", "uninstall"},
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		extensions, err := manager.InstalledExtensions()
		handleErr(err)

		if id := lo.Must(cmd.Flags().GetString("id")); id != "" {
			toRemove, ok := lo.Find(extensions, func(e *extension.Extension) bool {
				return e.String() == id
			})

			if !ok {
				handleErr(fmt.Errorf(
					"extension %s not found",
					style.Fg(color.Purple)(id),
				))
			}

			handleErr(manager.UninstallExtension(toRemove))

			return
		}

		var nameExtensionMap = make(map[string]*extension.Extension)

		for _, ext := range extensions {
			nameExtensionMap[ext.String()] = ext
		}

		var selected []string
		err = survey.AskOne(&survey.MultiSelect{
			Message: "Select extensions to remove",
			Options: lo.Keys(nameExtensionMap),
		}, &selected)
		handleErr(err)

		var confirm bool
		err = survey.AskOne(&survey.Confirm{
			Message: fmt.Sprintf("Remove %s?", text.Quantify(len(selected), "extension", "extensions")),
			Default: false,
		}, &confirm)
		handleErr(err)

		if !confirm {
			fmt.Printf("%s OK, aborted\n", style.Fg(color.Green)(icon.Check))
			return
		}

		for _, s := range selected {
			err = manager.UninstallExtension(nameExtensionMap[s])
			handleErr(err)
		}

		fmt.Printf(
			"%s Successfully removed %s\n",
			style.Fg(color.Green)(icon.Check),
			text.Quantify(len(selected), "extension", "extensions"),
		)
	},
}

func init() {
	xCmd.AddCommand(xSelCmd)

	xSelCmd.Flags().StringP("path", "p", "", "path to the extension")
	xSelCmd.Flags().StringP("id", "i", "", "id of the extension")
	xSelCmd.MarkFlagsMutuallyExclusive("path", "id")
	xSelCmd.MarkFlagDirname("path")
	xSelCmd.RegisterFlagCompletionFunc("id", completionExtensionsIDs)

	xSelCmd.Flags().Bool("run", false, "run the selected extension")
	xSelCmd.Flags().Bool("test", false, "test the selected extension")
	xSelCmd.Flags().Bool("info", false, "show info about the extension")
	xSelCmd.Flags().BoolP("json", "j", false, "output in json format")

	xSelCmd.MarkFlagsMutuallyExclusive("run", "test", "info")
}

var xSelCmd = &cobra.Command{
	Use:   "sel",
	Short: "Select an extension to perform an action on",
	Args:  cobra.NoArgs,
	PreRunE: preRunERequiredMutuallyExclusiveFlags(
		[]string{"path", "id"},
		[]string{"run", "test", "info"},
	),
	Run: func(cmd *cobra.Command, args []string) {
		ext := loadExtension(cmd.Flag("path"), cmd.Flag("id"))
		ext.Init()

		switch {
		case lo.Must(cmd.Flags().GetBool("run")):
			handleErr(ext.LoadScraper())
		case lo.Must(cmd.Flags().GetBool("test")):
			handleErr(ext.LoadTester())
			handleErr(ext.Tester().Test())
		case lo.Must(cmd.Flags().GetBool("info")):
			fmt.Println(ext.Passport().Info())
		}
	},
}

func init() {
	xCmd.AddCommand(xAddCmd)

	xAddCmd.Flags().BoolP("yes", "y", false, "skip install confirmation")
	xAddCmd.Flags().BoolP("force", "f", false, "skip passport validation")

	xAddCmd.MarkFlagsMutuallyExclusive("yes", "force")
}

var xAddCmd = &cobra.Command{
	Use:     "add",
	Short:   "Install an extension",
	Aliases: []string{"install"},
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		arg := args[0]

		var url string

		if text.IsURL(arg) {
			url = arg
		} else if matched, _ := regexp.MatchString(`^[\w-]+/[\w-]+$`, arg); matched {
			url = fmt.Sprintf("https://github.com/%s", arg)
		}

		ext, err := manager.InstallExtension(&manager.InstallOptions{
			URL:          url,
			SkipConfirm:  lo.Must(cmd.Flags().GetBool("yes")),
			SkipValidate: lo.Must(cmd.Flags().GetBool("force")),
		})
		handleErr(err)

		fmt.Printf("%s Successfully installed %s\n", style.Fg(color.Green)(icon.Check), style.Fg(color.Purple)(ext.String()))
	},
}

func init() {
	xCmd.AddCommand(xUpCmd)
}

var xUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Update extensions",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		extensions, err := manager.InstalledExtensions()
		handleErr(err)

		var updated int
		for _, ext := range extensions {
			err = manager.UpdateExtension(ext)
			if err != nil {
				fmt.Printf("%s %s %s\n", style.Fg(color.Red)(icon.Cross), style.Fg(color.Purple)(ext.String()), err)
			} else {
				fmt.Printf("%s %s %s\n", style.Fg(color.Green)(icon.Check), style.Fg(color.Purple)(ext.String()), "updated")
				updated++
			}
		}

		fmt.Printf("%s Successfully updated %s\n", style.Fg(color.Green)(icon.Check), text.Quantify(updated, "extension", "extensions"))
	},
}