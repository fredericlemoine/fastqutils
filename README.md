# Fastqutils
[![build](https://github.com/fredericlemoine/fastqutils/actions/workflows/go.yml/badge.svg)](https://github.com/fredericlemoine/fastqutils/actions)

Tool to manipulate fastq and bam files.

Available commands :


-  bamtofasta  Converts the input bam file in fasta alignment
-  cap         Downsample reads at regions with too high coverage
-  deinterlace Place the first reads on file 1 and second reads on file 2
-  filter      Commands to filter reads
-  generate    Generates a random Fastq file
-  help        Help about any command
-  mask        Mask nucleotides from bam or fastq files
-  sample      Subsample a FastQ File
-  stats       Displays different statistics about fastq file(s)
-  tobam       Generates an unaligned bam file from FASTQ File(s)
-  tofasta     Converts input fastq file into fasta
-  varcap      Downsample reads at regions with too high coverage. Given maximum coverage can be variable along the genome.
-  version     Prints the version of fastqutils
