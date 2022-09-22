package cmd

import (
	"log"
	"os"

	"github.com/biogo/hts/bam"
	"github.com/biogo/hts/sam"
	"github.com/fredericlemoine/fastqutils/io"
	"github.com/fredericlemoine/fastqutils/stats"
	"github.com/spf13/cobra"
)

var output string

// tobamCmd represents the tobam command
var tobamCmd = &cobra.Command{
	Use:   "tobam",
	Short: "Generates an unaligned bam file from FASTQ File(s)",
	Long: `Generates an unaligned bam file
`,
	Run: func(cmd *cobra.Command, args []string) {
		var bamwriter *bam.Writer
		var header *sam.Header
		var r1, r2 *sam.Record
		var f *os.File
		var err error
		var parser *io.FastQParser
		var offset, bamoffset int
		var enc int
		//mdsum := md5.New()
		//io.WriteString(mdsum, "*")

		if enc, err = stats.EncodingFromString(encoding); err != nil {
			log.Fatal(err)
		}
		if offset, err = stats.EncodingOffset(enc); err != nil {
			log.Fatal(err)
		}
		if bamoffset, err = stats.EncodingOffset(stats.ILLUMINA_1_8); err != nil {
			log.Fatal(err)
		}

		if parser, err = openFastqParser(input1, input2); err != nil {
			log.Fatal(err)
		}

		if output == "stdout" || output == "-" {
			f = os.Stdout
		} else {
			if f, err = os.Create(output); err != nil {
				log.Fatal(err)
			}
		}
		if header, err = sam.NewHeader(nil, nil); err != nil {
			log.Fatal(err)
		}
		if bamwriter, err = bam.NewWriter(f, header, 1); err != nil {
			log.Fatal(err)
		}

		for {
			entry1, entry2, err := parser.NextEntry()
			if err != nil {
				if err.Error() != "EOF" {
					log.Fatal(err)
				}
				break
			}

			flag1 := sam.Read1 | sam.Unmapped
			if entry2 != nil {
				flag1 = flag1 | sam.Paired | sam.MateUnmapped
			}

			// We encode the quality with the right offset
			for i, q := range entry1.Quality {
				entry1.Quality[i] = byte(int(q) - offset + bamoffset)
			}

			if r1, err = sam.NewRecord(string(entry1.Name), nil, nil, -1, -1, 0, byte(0), []sam.CigarOp{}, entry1.Sequence, entry1.Quality, []sam.Aux{}); err != nil {
				log.Fatal(err)
			}
			r1.Flags = flag1
			if err = bamwriter.Write(r1); err != nil {
				log.Fatal(err)
			}

			if entry2 != nil {
				flag2 := sam.Read2 | sam.Unmapped | sam.Paired | sam.MateUnmapped
				// We encode the quality of read 2 with the right offset
				for i, q := range entry2.Quality {
					entry2.Quality[i] = byte(int(q) - offset + bamoffset)
				}
				if r2, err = sam.NewRecord(string(entry2.Name), nil, nil, -1, -1, 0, byte(0), []sam.CigarOp{}, entry2.Sequence, entry2.Quality, []sam.Aux{}); err != nil {
					log.Fatal(err)
				}
				r2.Flags = flag2
				if err = bamwriter.Write(r2); err != nil {
					log.Fatal(err)
				}
			}
		}

		bamwriter.Close()
		f.Close()
	},
}

func init() {
	RootCmd.AddCommand(tobamCmd)
	tobamCmd.PersistentFlags().StringVarP(&input1, "input1", "1", "stdin", "First read fastq file")
	tobamCmd.PersistentFlags().StringVarP(&input2, "input2", "2", "none", "Second read fastq file")
	tobamCmd.PersistentFlags().StringVarP(&output, "output", "o", "stdout", "Output unaligned BAM file")
	tobamCmd.PersistentFlags().StringVar(&encoding, "encoding", "illumina1.8", "Base quality encoding, possible values: sanger, solexa, illumina1.3, illumina1.5, illumina1.8")
}
