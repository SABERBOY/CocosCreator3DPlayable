package main

import (
	"context"
	"fmt"
	"os"

	"CocosCreator3DPlayable/internal/packager"

	"github.com/spf13/cobra"
)

func main() {
	if err := newRootCommand().ExecuteContext(context.Background()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "cocos-playable",
		Short:         "Pack a Cocos Creator web-mobile export into one playable HTML file",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	cmd.AddCommand(newPackCommand())
	return cmd
}

func newPackCommand() *cobra.Command {
	var opts packager.Options

	cmd := &cobra.Command{
		Use:   "pack --input <web-mobile-dir> --output <playable.html>",
		Short: "Pack a Cocos Creator 3.x web-mobile export",
		RunE: func(cmd *cobra.Command, args []string) error {
			report, err := packager.Run(cmd.Context(), opts)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "wrote %s\n", report.OutputFile)
			fmt.Fprintf(cmd.OutOrStdout(), "embedded %d resources and %d scripts\n", report.ResourceCount, report.ScriptCount)
			if report.SkippedSymlinks > 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "skipped %d symlinks\n", report.SkippedSymlinks)
			}
			for _, warning := range report.Warnings {
				fmt.Fprintf(cmd.ErrOrStderr(), "warning: %s\n", warning)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&opts.InputDir, "input", "", "Cocos Creator web-mobile export directory")
	cmd.Flags().StringVar(&opts.OutputFile, "output", "", "single playable HTML output file")
	cmd.Flags().BoolVar(&opts.Minify, "minify", false, "minify generated HTML")
	cmd.Flags().BoolVar(&opts.Force, "force", false, "overwrite an existing output file")
	_ = cmd.MarkFlagRequired("input")
	_ = cmd.MarkFlagRequired("output")

	return cmd
}
