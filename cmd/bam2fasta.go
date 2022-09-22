package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/biogo/hts/bam"
	"github.com/biogo/hts/sam"
	"github.com/spf13/cobra"
)

var toFastaOutput string
var toFastaInput string

// bam2FastaCmd represents the bam2FastaCmd command
var bam2FastaCmd = &cobra.Command{
	Use:   "bamtofasta",
	Short: "Converts the input bam file in fasta alignment",
	Long: `Converts the input bam file in fasta alignment
`,
	Run: func(cmd *cobra.Command, args []string) {
		var writer *bufio.Writer
		var bamreader *bam.Reader
		var prev, rec *sam.Record
		var header *sam.Header
		var outfile *os.File
		var infile *os.File
		var totalRefLength int
		var err error
		var offset int = 0
		if toFastaInput == "stdin" || toFastaInput == "-" {
			infile = os.Stdin
		} else {
			if infile, err = os.Open(toFastaInput); err != nil {
				log.Fatal(err)
			}
		}

		if bamreader, err = bam.NewReader(infile, 1); err != nil {
			log.Fatal(err)
		}
		header = bamreader.Header()
		for _, r := range header.Refs() {
			totalRefLength += r.Len()
		}

		if toFastaOutput == "stdout" || toFastaOutput == "-" {
			outfile = os.Stdout
		} else {
			if outfile, err = os.Create(toFastaOutput); err != nil {
				log.Fatal(err)
			}
		}
		writer = bufio.NewWriter(outfile)
		for {
			if rec, err = bamreader.Read(); err != nil {
				if err.Error() != "EOF" {
					log.Fatal(err)
				}
				break
			}

			if prev != nil && prev.Ref.ID() == rec.Ref.ID() && prev.Start() > rec.Start() {
				log.Fatal(fmt.Errorf("bam file is not sorted by coordinate, please consider using samtools sort"))
			}
			if prev != nil && prev.Ref.ID() != rec.Ref.ID() {
				offset += prev.Ref.Len()
			}

			if rec.Flags&sam.Unmapped == sam.Unmapped ||
				rec.Flags&sam.Secondary == sam.Secondary ||
				rec.Flags&sam.Supplementary == sam.Supplementary ||
				rec.Flags&sam.QCFail == sam.QCFail {
				continue
			}
			fmt.Fprintf(writer, ">%s\n", rec.Name)
			fmt.Fprintf(writer, "%s", strings.Repeat("-", offset+rec.Start()))
			seq := rec.Seq.Expand()
			pos := 0
			for _, op := range rec.Cigar {
				t := op.Type()
				c := t.Consumes()
				if t == sam.CigarMatch {
					fmt.Fprintf(writer, "%s", string(seq[pos:pos+op.Len()]))
				} else if t == sam.CigarDeletion {
					fmt.Fprintf(writer, "%s", strings.Repeat("-", op.Len()))
				}
				pos += c.Query * op.Len()
			}

			if (offset + rec.End()) < (totalRefLength) {
				fmt.Fprintf(writer, "%s", strings.Repeat("-", totalRefLength-(offset+rec.End())))
			}
			prev = rec
			fmt.Fprint(writer, "\n")
		}
		writer.Flush()
		outfile.Close()
	},
}

func init() {
	RootCmd.AddCommand(bam2FastaCmd)
	bam2FastaCmd.PersistentFlags().StringVarP(&toFastaInput, "input", "i", "stdin", "Input bam file")
	bam2FastaCmd.PersistentFlags().StringVarP(&toFastaOutput, "output", "o", "stdout", "Output aligned Fasta")
}
