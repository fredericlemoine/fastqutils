/*
bammask : Filter reads in fastq files

# Copyright © 2022 Institut Pasteur, Paris

Author: Frederic Lemoine

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"bufio"
	"compress/gzip"
	"log"
	"os"

	"github.com/fredericlemoine/fastqutils/io"

	"github.com/spf13/cobra"
)

var bothReads bool
var minLength, maxLength int

// qualityCmd represents the quality command
var filterLengthCmd = &cobra.Command{
	Use:   "length",
	Short: "Remove reads that outside the given length interval",
	Long: `Remove reads that outside the given length interval
	
	fastqutils filter length --min-length <> --max-length <> 

	if --min-length -1 (default): then no minimal length
	if --max-length -1 (default): then no maximal length
	`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		var parser *io.FastQParser
		var w1, w2 *bufio.Writer
		var f1, f2 *os.File
		var g1, g2 *gzip.Writer
		var remove1, remove2, toWrite bool

		nbrecords := 0
		discarded := 0

		if parser, err = openFastqParser(input1, input2); err != nil {
			log.Fatal(err)
		}

		if w1, g1, f1, err = io.GetWriter(output1, gziped); err != nil {
			log.Fatal(err)
		}

		if input2 != "none" && output2 != "none" {
			if w2, g2, f2, err = io.GetWriter(output2, gziped); err != nil {
				log.Fatal(err)
			}
		}

		for {
			entry1, entry2, err := parser.NextEntry()
			if err != nil {
				if err.Error() != "EOF" {
					log.Fatal(err)
				}
				break
			}

			remove1 = (minLength != -1 && (len(entry1.Sequence) < minLength)) || (maxLength != -1 && (len(entry1.Sequence) > maxLength))
			remove2 = true
			toWrite = true

			if entry2 != nil {
				remove2 = (minLength != -1 && (len(entry2.Sequence) < minLength)) || (maxLength != -1 && (len(entry2.Sequence) > maxLength))
			}

			toWrite = (bothReads && !remove1 && !remove2) || (!bothReads && (!remove1 || !remove2))

			if toWrite {
				io.WriteEntry(w1, entry1)
				if w2 != nil {
					io.WriteEntry(w2, entry2)
				}
				nbrecords++
			} else {
				discarded++
			}
		}

		w1.Flush()
		if g1 != nil {
			g1.Flush()
			g1.Close()
		}
		f1.Close()
		if input2 != "none" && output2 != "none" {
			w2.Flush()
			if g2 != nil {
				g2.Flush()
				g2.Close()
			}
			f2.Close()
		}
		log.Printf("Wrote %d fastq records", nbrecords)
		log.Printf("Discarded %d fastq records", discarded)
	},
}

func init() {
	filterCmd.AddCommand(filterLengthCmd)
	filterLengthCmd.PersistentFlags().BoolVarP(&bothReads, "paired-both", "p", false, "Removes the two reads of a pair (if paired-end) if the two are outside length interval. Otherwise, removes the two reads if at least one read is outside of the range.")
	filterLengthCmd.PersistentFlags().IntVar(&minLength, "min-length", -1, "Minimum length to keep a read, default -1 (no cutoff)")
	filterLengthCmd.PersistentFlags().IntVar(&maxLength, "max-length", -1, "Maximum length to keep a read, default -1 (no cutoff)")
	filterLengthCmd.PersistentFlags().StringVar(&output1, "output1", "stdout", "Output file 1")
	filterLengthCmd.PersistentFlags().StringVar(&output2, "output2", "none", "Output file 2 (if paired)")
	filterLengthCmd.PersistentFlags().BoolVar(&gziped, "gz", false, "If true, will generate gziped file(s) : .gz extension is added automatically")
	filterLengthCmd.PersistentFlags().StringVarP(&input1, "input1", "1", "stdin", "First read fastq file")
	filterLengthCmd.PersistentFlags().StringVarP(&input2, "input2", "2", "none", "Second read fastq file")
}
