package cmd

import (
	"crypto/md5"
	"io"
	"os"

	"github.com/biogo/hts/bam"
	"github.com/biogo/hts/sam"
	errorp "github.com/fredericlemoine/fastqutils/error"
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
		mdsum := md5.New()
		io.WriteString(mdsum, "*")

		if output == "stdout" || output == "-" {
			f = os.Stdout
		} else {
			if f, err = os.Create(output); err != nil {
				errorp.ExitWithMessage(err)
			}
		}
		if header, err = sam.NewHeader(nil, nil); err != nil {
			errorp.ExitWithMessage(err)
		}
		bamwriter, err = bam.NewWriter(f, header, 1)

		for {
			entry1, entry2, err := parser.NextEntry()
			if err != nil {
				if err.Error() != "EOF" {
					errorp.WarnMessage(err)
				}
				break
			}

			flag1 := sam.Read1 | sam.Unmapped
			if entry2 != nil {
				flag1 = flag1 | sam.Paired | sam.MateUnmapped
			}
			if r1, err = sam.NewRecord(string(entry1.Name), nil, nil, -1, -1, 0, byte(0), []sam.CigarOp{}, entry1.Sequence, entry1.Quality, []sam.Aux{}); err != nil {
				errorp.ExitWithMessage(err)
			}
			r1.Flags = flag1
			if err = bamwriter.Write(r1); err != nil {
				errorp.ExitWithMessage(err)
			}

			if entry2 != nil {
				flag2 := sam.Read2 | sam.Unmapped | sam.Paired | sam.MateUnmapped
				if r2, err = sam.NewRecord(string(entry2.Name), nil, nil, -1, -1, 0, byte(0), []sam.CigarOp{}, entry2.Sequence, entry2.Quality, []sam.Aux{}); err != nil {
					errorp.ExitWithMessage(err)
				}
				r2.Flags = flag2
				if err = bamwriter.Write(r2); err != nil {
					errorp.ExitWithMessage(err)
				}
			}
		}

		bamwriter.Close()
		f.Close()
	},
}

func init() {
	RootCmd.AddCommand(tobamCmd)
	tobamCmd.PersistentFlags().StringVarP(&output, "output", "o", "stdout", "Output unaligned BAM file")
}
