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

	"github.com/biogo/hts/bam"
	"github.com/biogo/hts/sam"
	"github.com/fredericlemoine/fastqutils/fastq"
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

	To filter reads in bam files:
	fastqutils filter length --min-length <> --max-length <> -b  -i <inbam> -o <outbam>

	To filter reads in fastq files:
	fastqutils filter length --min-length <> --max-length <> -b  -1 <fastq1> -2 <fastq2> --output1 <outfastq1> --output2 <outfastq2>

	Option --paired-both is only functionnal for fastqfiles.
	For bam files, each record is kept or discarded based on its length, independently of its mate read.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		if bamformat {
			if err = filterLengthBam(inbam, outbam, minLength, maxLength); err != nil {
				log.Fatal(err)
			}
		} else {
			if err = filterLengthFastq(input1, input2, output1, output2, gziped, minLength, maxLength); err != nil {
				log.Fatal(err)
			}
		}
	},
}

func init() {
	filterCmd.AddCommand(filterLengthCmd)
	filterLengthCmd.PersistentFlags().BoolVarP(&bothReads, "paired-both", "p", false, "Removes the two reads of a pair (if paired-end) if the two are outside length interval. Otherwise, removes the two reads if at least one read is outside of the range.")
	filterLengthCmd.PersistentFlags().IntVar(&minLength, "min-length", -1, "Minimum length to keep a read, default -1 (no cutoff)")
	filterLengthCmd.PersistentFlags().IntVar(&maxLength, "max-length", -1, "Maximum length to keep a read, default -1 (no cutoff)")
	filterLengthCmd.PersistentFlags().StringVarP(&inbam, "input-bam", "i", "stdin", "Input bam file")
	filterLengthCmd.PersistentFlags().StringVarP(&outbam, "out-bam", "o", "stdout", "Output bam file")
	filterLengthCmd.PersistentFlags().BoolVarP(&bamformat, "bam", "b", false, "Whether the input is bam or fastq format")
	filterLengthCmd.PersistentFlags().StringVar(&output1, "output1", "stdout", "Output file 1")
	filterLengthCmd.PersistentFlags().StringVar(&output2, "output2", "none", "Output file 2 (if paired)")
	filterLengthCmd.PersistentFlags().BoolVar(&gziped, "gz", false, "If true, will generate gziped file(s) : .gz extension is added automatically")
	filterLengthCmd.PersistentFlags().StringVarP(&input1, "input1", "1", "stdin", "First read fastq file")
	filterLengthCmd.PersistentFlags().StringVarP(&input2, "input2", "2", "none", "Second read fastq file")
}

func filterLengthFastq(input1, input2, output1, output2 string, gziped bool, minLength, maxLength int) (err error) {
	var parser *io.FastQParser
	var w1, w2 *bufio.Writer
	var f1, f2 *os.File
	var g1, g2 *gzip.Writer
	var remove1, remove2, toWrite bool
	var entry1, entry2 *fastq.FastqEntry

	nbrecords := 0
	discarded := 0

	if parser, err = openFastqParser(input1, input2); err != nil {
		return
	}

	if w1, g1, f1, err = io.GetWriter(output1, gziped); err != nil {
		return
	}

	if input2 != "none" && output2 != "none" {
		if w2, g2, f2, err = io.GetWriter(output2, gziped); err != nil {
			return
		}
	}

	for {
		entry1, entry2, err = parser.NextEntry()
		if err != nil {
			if err.Error() != "EOF" {
				return
			}
			err = nil
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

	return
}

func filterLengthBam(inbam, outbam string, minLength, maxLength int) (err error) {
	var bamwriter *bam.Writer
	var bamreader *bam.Reader
	var header *sam.Header
	var rec *sam.Record
	var outfile *os.File
	var infile *os.File
	var nbrecords, discarded int

	// Opening new bam reader
	if inbam == "stdin" || inbam == "-" {
		infile = os.Stdin
	} else {
		if infile, err = os.Open(inbam); err != nil {
			return
		}
	}
	if bamreader, err = bam.NewReader(infile, 1); err != nil {
		return
	}
	header = bamreader.Header()

	// Opening new bam writer
	if outbam == "stdout" || outbam == "-" {
		outfile = os.Stdout
	} else {
		if outfile, err = os.Create(outbam); err != nil {
			return
		}
	}
	if bamwriter, err = bam.NewWriter(outfile, header, 1); err != nil {
		return
	}

	// Reading bam file, record by record
	for {
		if rec, err = bamreader.Read(); err != nil {
			if err.Error() != "EOF" {
				return
			}
			err = nil
			break
		}

		// Converting sequence from doublets to one byte per nucleotide
		reclen := rec.Seq.Length
		remove := (minLength != -1 && (reclen < minLength)) || (maxLength != -1 && (reclen > maxLength))

		if !remove {
			// We write the record in the output
			if err = bamwriter.Write(rec); err != nil {
				return
			}
			nbrecords++
		} else {
			discarded++
		}
	}

	bamwriter.Close()
	outfile.Close()

	log.Printf("Wrote %d fastq records", nbrecords)
	log.Printf("Discarded %d fastq records", discarded)

	return
}
