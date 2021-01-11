package cmd

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/biogo/hts/bam"
	"github.com/biogo/hts/sam"
	errorp "github.com/fredericlemoine/fastqutils/error"
	"github.com/spf13/cobra"
)

var capOutput string
var capInput string
var capReadLength int
var capWindowSize int
var capCoverage int

// capCmd represents the tobam command
var capCmd = &cobra.Command{
	Use:   "cap",
	Short: "Downsample reads at regions with too high coverage",
	Long: `Produces a bam file with lower coverage

	Be careful: 
	1) Input bam file must be sorted
	2) Output bam file is not sorted anymore
`,
	Run: func(cmd *cobra.Command, args []string) {
		var bamwriter *bam.Writer
		var bamreader *bam.Reader
		var header *sam.Header
		var prev, rec *sam.Record
		var outfile *os.File
		var infile *os.File
		var reservoir []*sam.Record
		var curRef int
		var err error
		var windowStart, windowElements int
		var maxNReadsPerWindow int = int(float64(capCoverage*capWindowSize) / float64(capReadLength))

		if capInput == "stdin" || capInput == "-" {
			infile = os.Stdin
		} else {
			if infile, err = os.Open(capInput); err != nil {
				errorp.ExitWithMessage(err)
			}
		}
		if bamreader, err = bam.NewReader(infile, 1); err != nil {
			errorp.ExitWithMessage(err)
		}
		header = bamreader.Header()

		if capOutput == "stdout" || capOutput == "-" {
			outfile = os.Stdout
		} else {
			if outfile, err = os.Create(capOutput); err != nil {
				errorp.ExitWithMessage(err)
			}
		}
		bamwriter, err = bam.NewWriter(outfile, header, 1)
		curRef = -1
		reservoir = make([]*sam.Record, maxNReadsPerWindow)
		windowStart = 0
		windowElements = 0
		for {
			if rec, err = bamreader.Read(); err != nil {
				if err.Error() != "EOF" {
					errorp.WarnMessage(err)
				}
				break
			}

			if prev != nil && prev.Ref.ID() == rec.Ref.ID() && prev.Start() > rec.Start() {
				errorp.ExitWithMessage(fmt.Errorf("Bam file is not sorted by coordinate, please consider using samtools sort"))
			}

			if curRef != rec.Ref.ID() {
				fmt.Printf("WriteReservoir ChangeChr: %d - %d - %d\n", maxNReadsPerWindow, windowStart, windowElements)
				if err = writeReservoir(bamwriter, reservoir, min(windowElements, maxNReadsPerWindow)); err != nil {
					errorp.ExitWithMessage(err)
				}
				// Then change the ref
				curRef = rec.Ref.ID()
				reservoir = make([]*sam.Record, maxNReadsPerWindow)
				windowStart = 0
				windowElements = 0
			}

			if rec.Flags&sam.Unmapped == sam.Unmapped ||
				rec.Flags&sam.Secondary == sam.Secondary ||
				rec.Flags&sam.Supplementary == sam.Supplementary ||
				rec.Flags&sam.QCFail == sam.QCFail {
				continue
			}

			if rec.Start() >= windowStart+capWindowSize {
				fmt.Printf("WriteReservoir ChangeWindows: %d - %d - %d\n", maxNReadsPerWindow, windowStart, windowElements)
				if err = writeReservoir(bamwriter, reservoir, min(windowElements, maxNReadsPerWindow)); err != nil {
					errorp.ExitWithMessage(err)
				}
				reservoir = make([]*sam.Record, maxNReadsPerWindow)
				windowStart = rec.Start()
				windowElements = 0
			} else {
				// Reservoir Sampling in the window
				if windowElements < maxNReadsPerWindow {
					reservoir[windowElements] = rec
				} else {
					random := rand.Intn(windowElements)
					if random < maxNReadsPerWindow {
						reservoir[random] = rec
					}
				}
				windowElements++
			}
			prev = rec
		}
		if err = writeReservoir(bamwriter, reservoir, min(windowElements, maxNReadsPerWindow)); err != nil {
			errorp.ExitWithMessage(err)
		}

		bamwriter.Close()
		outfile.Close()
	},
}

func writeReservoir(bamwriter *bam.Writer, reservoir []*sam.Record, windowElements int) (err error) {
	for i := 0; i < windowElements; i++ {
		if err = bamwriter.Write(reservoir[i]); err != nil {
			return
		}
	}
	return
}

func init() {
	RootCmd.AddCommand(capCmd)
	capCmd.PersistentFlags().StringVarP(&capInput, "input", "i", "stdin", "Input bam file")
	capCmd.PersistentFlags().IntVarP(&capReadLength, "length", "l", 200, "Read Length")
	capCmd.PersistentFlags().IntVarP(&capWindowSize, "window", "w", 200, "Size of the coverage windows")
	capCmd.PersistentFlags().StringVarP(&capOutput, "output", "o", "stdout", "Output unaligned BAM file")
	capCmd.PersistentFlags().IntVarP(&capCoverage, "coverage", "c", 100, "Max desired output coverage")
}
