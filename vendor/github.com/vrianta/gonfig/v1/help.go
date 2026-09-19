package structparser

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

/**
 * Code here will help to do auto print of arguments the application supports
 */

type help_record struct {
	field       string
	args        string
	env         string
	description string
	def         string
	required    bool
}

// help formats and prints the configuration help metadata to os.Stdout.
func help(help_records []help_record) {
	if len(help_records) == 0 {
		return
	}

	fmt.Println("\nUsage & Options:")
	fmt.Println(strings.Repeat("-", 80))

	// Use tabwriter to keep columns aligned dynamically
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "FIELD\tFLAG\tENV VAR\tDEFAULT\tREQUIRED\tDESCRIPTION")
	fmt.Fprintln(w, "-----\t----\t-------\t-------\t--------\t-----------")

	for _, rec := range help_records {
		argStr := "-"
		if rec.args != "" {
			argStr = "--" + rec.args
		}

		envStr := "-"
		if rec.env != "" {
			envStr = rec.env
		}

		defStr := "-"
		if rec.def != "" {
			defStr = rec.def
		}

		reqStr := "false"
		if rec.required {
			reqStr = "true"
		}

		descStr := "-"
		if rec.description != "" {
			descStr = rec.description
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
			rec.field,
			argStr,
			envStr,
			defStr,
			reqStr,
			descStr,
		)
	}

	w.Flush()
	fmt.Println()
}
