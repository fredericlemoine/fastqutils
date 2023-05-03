# Fastqutils
[![build](https://github.com/fredericlemoine/fastqutils/actions/workflows/go.yml/badge.svg)](https://github.com/fredericlemoine/fastqutils/actions)

Tool to manipulate fastq and bam files.

Available commands :

- bamtofasta  Converts the input bam file in fasta alignment
- cap         Downsample reads in regions with too much coverage
- completion  Generate the autocompletion script for the specified shell
- filter      Removes reads having length outside a given range
- generate    Generates a random Fastq file (for test purpose)
- help        Help about any command
- mask        Mask nucleotides from bam or fastq files
- sample      Subsample a FastQ File
- stats       Displays different statistics about fastq file
- tobam       Generates an unaligned bam file from FASTQ file
- tofasta     Converts input fastq file into fasta
- version     Prints the version of fastqutils
