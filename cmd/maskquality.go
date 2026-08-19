/*
bammask : Masking nucleotides in bam files

Copyright © 2022 Institut Pasteur, Paris

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
	stdio "io"
	"log"
	"os"

	"github.com/biogo/hts/bam"
	"github.com/biogo/hts/sam"
	"github.com/fredericlemoine/fastqutils/fastq"
	"github.com/fredericlemoine/fastqutils/io"
	"github.com/fredericlemoine/fastqutils/stats"

	"github.com/spf13/cobra"
)

var inbam, outbam string
var bamformat bool
var qual int

// qualityCmd represents the quality command
var qualityCmd = &cobra.Command{
	Use:   "quality",
	Short: "Mask read bases in a bam file, based on base quality",
	Long:  `Mask read nucleotides in a bam file, based on base quality`,
	Run: func(cmd *cobra.Command, args []string) {
		if bamformat {
			if err := maskQUalityBam(inbam, outbam, qual); err != nil {
				log.Fatal(err)
			}
		} else {
			if err := maskQualityFastq(input1, input2, encoding, output1, output2, gziped, dsrcOut, qual); err != nil {
				log.Fatal(err)
			}
		}
	},
}

func init() {
	maskCmd.AddCommand(qualityCmd)

	qualityCmd.PersistentFlags().BoolVarP(&bamformat, "bam", "b", false, "Whether the input is bam or fastq format")
	qualityCmd.PersistentFlags().StringVarP(&inbam, "input-bam", "i", "stdin", "Input bam file")
	qualityCmd.PersistentFlags().StringVarP(&outbam, "out-bam", "o", "stdout", "Output bam file")
	qualityCmd.PersistentFlags().StringVar(&output1, "output1", "stdout", "Output file 1")
	qualityCmd.PersistentFlags().StringVar(&output2, "output2", "none", "Output file 2 (if paired)")
	qualityCmd.PersistentFlags().StringVarP(&input1, "input1", "1", "stdin", "First read fastq file")
	qualityCmd.PersistentFlags().StringVarP(&input2, "input2", "2", "none", "Second read fastq file")
	qualityCmd.PersistentFlags().StringVar(&encoding, "encoding", "illumina1.8", "Base quality encoding, possible values: sanger, solexa, illumina1.3, illumina1.5, illumina1.8")
	qualityCmd.PersistentFlags().IntVarP(&qual, "quality", "q", 20, "Quality cutoff below which bases are masked")
	qualityCmd.PersistentFlags().BoolVar(&gziped, "gz", false, "If true, will generate gziped file(s) : .gz extension is added automatically")
	qualityCmd.PersistentFlags().BoolVar(&dsrcOut, "dsrc", false, "If true, will generate dsrc-compressed file(s) : .dsrc extension is added automatically (requires the 'dsrc' executable, see dsrc/README.md)")
}

func maskQUalityBam(inbam, outbam string, qual int) (err error) {
	var bamwriter *bam.Writer
	var bamreader *bam.Reader
	var header *sam.Header
	var rec *sam.Record
	var outfile *os.File
	var infile *os.File

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
		s := rec.Seq.Expand()
		modified := false
		for i, q := range rec.Qual {
			// If the base at the current index has bad quality, we replace it with a N
			if int(q) < qual {
				modified = true
				s[i] = 'N'
			}
		}
		// If the sequence has been modified, we update it
		if modified {
			rec.Seq = sam.NewSeq(s)
		}

		// We write the record in the output
		if err = bamwriter.Write(rec); err != nil {
			return
		}
	}

	bamwriter.Close()
	outfile.Close()
	return
}

func maskQualityFastq(input1, input2, encoding, output1, output2 string, gziped, dsrced bool, qual int) (err error) {
	var w1, w2 *bufio.Writer
	var closer1, closer2 stdio.Closer
	var entry1, entry2 *fastq.FastqEntry
	var parser *io.FastQParser
	var enc int
	var offset int

	if enc, err = stats.EncodingFromString(encoding); err != nil {
		return
	}
	if offset, err = stats.EncodingOffset(enc); err != nil {
		return
	}

	if parser, err = openFastqParser(input1, input2); err != nil {
		return
	}
	defer parser.Close()

	if w1, closer1, err = io.GetWriter(output1, gziped, dsrced); err != nil {
		return
	}
	if input2 != "none" && output2 != "none" {
		if w2, closer2, err = io.GetWriter(output2, gziped, dsrced); err != nil {
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
		for i, q := range entry1.Quality {
			// If the base at the current index has bad quality, we replace it with a N
			if int(q)-offset < qual {
				entry1.Sequence[i] = 'N'
			}
		}
		// NOTE: this used to call io.WriteEntryFasta here, which silently
		// wrote 2-line FASTA records (name+sequence, no quality) even
		// though this command reads and is meant to write FASTQ. Fixed to
		// io.WriteEntry so quality-masked reads keep their quality scores
		// and the output is valid FASTQ -- which --gz/--dsrc compression
		// also depends on, since DSRC specifically expects 4-line FASTQ
		// records.
		io.WriteEntry(w1, entry1)
		if w2 != nil {
			for i, q := range entry2.Quality {
				// If the base at the current index has bad quality, we replace it with a N
				if int(q)-offset < qual {
					entry2.Sequence[i] = 'N'
				}
			}
			io.WriteEntry(w2, entry2)
		}
	}

	if err = closer1.Close(); err != nil {
		return
	}

	if input2 != "none" && output2 != "none" {
		if err = closer2.Close(); err != nil {
			return
		}
	}
	return
}
