package cmd

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/biogo/hts/bam"
	"github.com/biogo/hts/sam"
	"github.com/fredericlemoine/fastqutils/io"
	"github.com/spf13/cobra"
)

var varCapOutput string
var varCapInput string
var varCapReadLength int
var varCapWindowSize int
var varCapCoverageFile string

// capCmd represents the tobam command
var varCapCmd = &cobra.Command{
	Use:   "varcap",
	Short: "Downsample reads at regions with too high coverage. Given maximum coverage can be variable along the genome.",
	Long: `Produces a bam file with lower coverage, having a pattern as given in an input coverage file.

	The coverage input file (-c) must be tab separated with the following fields:
	- chromosome
	- start
	- end
	- desired coverage

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
		var maxNReadsPerWindow int
		var incov map[string][]int
		var varCapCoverage int

		if incov, err = parseInputCoverageFile(varCapCoverageFile); err != nil {
			log.Fatal(err)
		}

		if varCapInput == "stdin" || varCapInput == "-" {
			infile = os.Stdin
		} else {
			if infile, err = os.Open(varCapInput); err != nil {
				log.Fatal(err)
			}
		}
		if bamreader, err = bam.NewReader(infile, 1); err != nil {
			log.Fatal(err)
		}
		header = bamreader.Header()

		if varCapOutput == "stdout" || varCapOutput == "-" {
			outfile = os.Stdout
		} else {
			if outfile, err = os.Create(varCapOutput); err != nil {
				log.Fatal(err)
			}
		}
		if bamwriter, err = bam.NewWriter(outfile, header, 1); err != nil {
			log.Fatal(err)
		}
		curRef = -1
		windowStart = 0
		windowElements = 0
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

			if curRef != rec.Ref.ID() {
				if reservoir != nil {
					fmt.Printf("WriteReservoir ChangeChr: %d - %d - %d\n", maxNReadsPerWindow, windowStart, windowElements)
					if err = writeReservoir(bamwriter, reservoir, min(windowElements, maxNReadsPerWindow)); err != nil {
						log.Fatal(err)
					}
				}
				// Then change the ref
				windowStart = 0
				windowElements = 0
				curRef = rec.Ref.ID()
				covslice, ok := incov[rec.Ref.Name()]
				if !ok {
					varCapCoverage = 10000000
				} else if len(covslice) < windowStart || covslice[windowStart] <= 0 {
					varCapCoverage = 10000000
				} else {
					varCapCoverage = covslice[windowStart]
				}
				maxNReadsPerWindow = int(float64(varCapCoverage*varCapWindowSize) / float64(varCapReadLength))
				reservoir = make([]*sam.Record, maxNReadsPerWindow)
			}

			if rec.Flags&sam.Unmapped == sam.Unmapped ||
				rec.Flags&sam.Secondary == sam.Secondary ||
				rec.Flags&sam.Supplementary == sam.Supplementary ||
				rec.Flags&sam.QCFail == sam.QCFail {
				continue
			}

			if rec.Start() >= windowStart+varCapWindowSize {
				fmt.Printf("WriteReservoir ChangeWindows: %d - %d - %d\n", maxNReadsPerWindow, windowStart, windowElements)
				if err = writeReservoir(bamwriter, reservoir, min(windowElements, maxNReadsPerWindow)); err != nil {
					log.Fatal(err)
				}
				windowStart = rec.Start()
				windowElements = 0
				covslice, ok := incov[rec.Ref.Name()]
				if !ok {
					varCapCoverage = 10000000
				} else if len(covslice) < windowStart || covslice[windowStart] <= 0 {
					varCapCoverage = 10000000
				} else {
					varCapCoverage = covslice[windowStart]
				}
				maxNReadsPerWindow = int(float64(varCapCoverage*varCapWindowSize) / float64(varCapReadLength))
				reservoir = make([]*sam.Record, maxNReadsPerWindow)
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
			log.Fatal(err)
		}

		bamwriter.Close()
		outfile.Close()
	},
}

func init() {
	RootCmd.AddCommand(varCapCmd)
	varCapCmd.PersistentFlags().StringVarP(&varCapInput, "input", "i", "stdin", "Input bam file")
	varCapCmd.PersistentFlags().IntVarP(&varCapReadLength, "length", "l", 200, "Read Length")
	varCapCmd.PersistentFlags().IntVarP(&varCapWindowSize, "window", "w", 200, "Size of the coverage windows")
	varCapCmd.PersistentFlags().StringVarP(&varCapOutput, "output", "o", "stdout", "Output unaligned BAM file")
	varCapCmd.PersistentFlags().StringVarP(&varCapCoverageFile, "coverage", "c", "none", "Input file describing output variable genome coverage ")
}

// Input file:
// - 1 line per coverage
// - tab separated: chrom \t start \t end \t coverage
// Output structure:
// map:
// - key: chromosome
// - value: Array of coverages
//   - length: max start/end given in input for this chromosome
//   - values: desired coverage at each position. If not specified: -1 (means no constraint)
func parseInputCoverageFile(infile string) (cov map[string][]int, err error) {
	var ifilereader *bufio.Reader
	var line string
	var err2 error
	var cols []string
	var chrom string
	var start, end int
	var poscov int   // tmp variable for coverage at a single line of the input
	var chrcov []int // tmp coverage slice for a given chromosome
	var ok bool      // tmp variable to test existance of a key in the map

	cov = make(map[string][]int)
	if ifilereader, err = io.GetReader(infile); err == nil {
		line, err2 = Readln(ifilereader)
		for err2 == nil {
			cols = strings.Split(line, "\t")
			if len(cols) != 4 {
				err = fmt.Errorf("input coverage file does not contain the right columns")
				return
			}
			chrom = cols[0]
			if start, err = strconv.Atoi(cols[1]); err != nil {
				return
			}
			if end, err = strconv.Atoi(cols[2]); err != nil {
				return
			}
			if poscov, err = strconv.Atoi(cols[3]); err != nil {
				return
			}

			if start < 0 {
				err = fmt.Errorf("input start cannot be negative")
				return
			}

			if end < 0 {
				err = fmt.Errorf("input end cannot be negative")
				return
			}
			if end < start {
				err = fmt.Errorf("end cannot be < start")
				return
			}
			if poscov < 0 {
				err = fmt.Errorf("given coverage cannot be negative")
				return
			}

			if chrcov, ok = cov[chrom]; !ok {
				chrcov = make([]int, 0, 10000)
				cov[chrom] = chrcov
			}

			if len(chrcov) < end {
				tmpcov := make([]int, end*2)
				for i, v := range chrcov {
					tmpcov[i] = v
				}
				chrcov = tmpcov
				cov[chrom] = chrcov
			}

			for i := start; i < end; i++ {
				chrcov[i] = poscov
			}
			line, err2 = Readln(ifilereader)
		}
	}

	return
}

// Readln returns a single line (without the ending \n)
// from the input buffered reader.
// An error is returned iff there is an error with the
// buffered reader.
func Readln(r *bufio.Reader) (string, error) {
	var (
		isPrefix bool  = true
		err      error = nil
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln), err
}
