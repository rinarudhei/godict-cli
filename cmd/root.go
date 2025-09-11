/*
Copyright © 2025 Rinaldi Adrian Mohammad <rinaldiadrian5@gmail.com>
*/
package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var verbose bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:          "godict-cli",
	Short:        "English dictionary CLI application",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		apiURL := viper.GetString("api-url")

		return lookAction(os.Stdout, verbose, apiURL, args[0])
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().String("api-url", "https://api.dictionaryapi.dev/api/v2/entries/en/", "Free Dictionary API")
	rootCmd.PersistentFlags().Duration("api-timeout", 5, "Set API timeout in duration second")
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix("GODICT")

	viper.BindPFlag("api-url", rootCmd.PersistentFlags().Lookup("api-url"))
	viper.BindPFlag("api-timeout", rootCmd.PersistentFlags().Lookup("api-timeout"))
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose response")
}

func lookAction(out io.Writer, verbose bool, apiURL, word string) error {
	resps, err := lookDictionary(apiURL, word)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return printEntryNotFoundResponse(out)
		}

		return err
	}

	if verbose {
		return printVerboseResponse(out, resps[0])
	}

	return printSimpleResponse(out, resps[0])
}

// printEntryNotFoundResponse prints not found message
func printEntryNotFoundResponse(w io.Writer) error {
	_, err := fmt.Fprintln(w, "No Definitions Found.\nSorry pal, we couldn't find definitions for the word you were looking for.\nYou can try the search again at later time or head to the web instead.\n")
	return err
}

// printSimpleResponse prints minimal output
func printSimpleResponse(w io.Writer, entry DictionaryEntry) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	fmt.Fprintf(tw, "Word:\t%s\n", entry.Word)
	if entry.Phonetic != "" {
		fmt.Fprintf(tw, "Phonetic:\t%s\n", entry.Phonetic)
	}

	// print first meaning and first definition if available
	if len(entry.Meanings) > 0 && len(entry.Meanings[0].Definitions) > 0 {
		def := entry.Meanings[0].Definitions[0]
		fmt.Fprintf(tw, "Meaning:\t%s\n", def.Definition)
		if def.Example != "" {
			fmt.Fprintf(tw, "Example:\t%s\n", def.Example)
		}
	}

	return tw.Flush()
}

// printVerboseResponse prints detailed output
func printVerboseResponse(w io.Writer, entry DictionaryEntry) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	fmt.Fprintf(tw, "Word:\t%s\n", entry.Word)
	if entry.Phonetic != "" {
		fmt.Fprintf(tw, "Phonetic:\t%s\n", entry.Phonetic)
	}

	// Print phonetics
	if len(entry.Phonetics) > 0 {
		fmt.Fprintln(tw, "Phonetics:")
		for _, p := range entry.Phonetics {
			if p.Audio != "" {
				fmt.Fprintf(tw, "\t- %s (audio: %s)\n", p.Text, p.Audio)
			} else {
				fmt.Fprintf(tw, "\t- %s\n", p.Text)
			}
		}
	}

	// Print origin
	if entry.Origin != "" {
		fmt.Fprintf(tw, "Origin:\t%s\n", entry.Origin)
	}

	// Print meanings
	if len(entry.Meanings) > 0 {
		fmt.Fprintln(tw, "Meanings:")
		for _, m := range entry.Meanings {
			fmt.Fprintf(tw, "\tPart of Speech:\t%s\n", m.PartOfSpeech)
			for _, d := range m.Definitions {
				fmt.Fprintf(tw, "\tDefinition:\t%s\n", d.Definition)
				if d.Example != "" {
					fmt.Fprintf(tw, "\tExample:\t%s\n", d.Example)
				}
				if len(d.Synonyms) > 0 {
					fmt.Fprintf(tw, "\tSynonyms:\t%v\n", d.Synonyms)
				}
				if len(d.Antonyms) > 0 {
					fmt.Fprintf(tw, "\tAntonyms:\t%v\n", d.Antonyms)
				}
			}
		}
	}

	return tw.Flush()
}
