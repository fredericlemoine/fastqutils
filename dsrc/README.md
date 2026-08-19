# DSRC support

This package lets fastqutils read and write `.dsrc`-compressed FASTQ
files, alongside the plain and `.gz` files it already supports.

## How it works

fastqutils does **not** link against DSRC. DSRC
(https://github.com/refresh-bio/DSRC) is a separate C++ program licensed
under GPLv2, while fastqutils is Apache-2.0; statically embedding DSRC's
source (or a `libdsrc.a` built from it) into the fastqutils binary would
create a real license entanglement between the two, not just added build
complexity. So instead, this package shells out to the `dsrc`
command-line executable at run time, using DSRC's own `-s` streaming
flag:

- Reading a `.dsrc` file runs `dsrc d -s <file>` and reads the
  decompressed FASTQ data from the subprocess's stdout.
- Writing a `.dsrc` file runs `dsrc c -s <file>` and writes raw FASTQ
  data to the subprocess's stdin; dsrc compresses it and writes the
  archive to `<file>` itself (DSRC archives are not stream-friendly on
  the compressed side, so the archive is always a real file, never
  stdout).

**Known caveat:** DSRC's `-s` streaming mode has an open upstream
reliability report (https://github.com/refresh-bio/DSRC/issues/1), and
building DSRC from source surfaces real mismatched `new[]`/`delete` bugs
in its buffer-handling code — undefined behavior that different
platforms' allocators tolerate differently. It has been observed to
segfault on macOS (while working fine on Linux). If you hit
`dsrc: signal: segmentation fault`, the more robust alternative is to
stop passing `-s` and instead round-trip through real temp files using
DSRC's plain file-to-file interface (`dsrc c in out.dsrc` /
`dsrc d in.dsrc out`) — DSRC's primary, best-tested code path. That
version trades away true streaming (needs scratch disk space, and
`GetReader`/`GetWriter` are no longer zero-copy) for reliability; ask
for it if `-s` proves unworkable in your environment.

This keeps fastqutils itself free of cgo and C++ toolchain requirements,
so ordinary `go build` / cross-compilation (see the `deploy` target in
the top-level `Makefile`) keeps working unmodified for anyone who
doesn't touch `.dsrc` files.

## Requirements

Install DSRC separately and make sure the `dsrc` executable is on your
`$PATH`:

```
git clone https://github.com/refresh-bio/DSRC.git
cd DSRC
make -f Makefile.c++11   # or Makefile / Makefile.osx, depending on platform
# copy or symlink the resulting `dsrc` binary somewhere on $PATH
```

If you'd rather not modify `$PATH`, set the `DSRC_BIN` environment
variable to the full path of the `dsrc` executable instead; fastqutils
checks it before falling back to `$PATH`.

## Usage from fastqutils

Input files ending in `.dsrc` are auto-detected, exactly like `.gz`
files are today — no flag needed:

```
fastqutils stats -1 reads.dsrc
```

For commands that write FASTQ output (`sample`, `deinterlace`,
`filter length`, `mask quality`, `generate`), pass `--dsrc` instead of
`--gz` to compress the output with DSRC; the `.dsrc` extension is added
automatically, the same way `--gz` adds `.gz`:

```
fastqutils sample -1 reads.fastq -n 1000 --output1 sample --dsrc
# writes sample.dsrc
```

`--gz` and `--dsrc` are mutually exclusive. Because DSRC archives must
be written to a real file, `--dsrc` cannot be combined with
`--output1 stdout` (or `-`) — pass an actual output file name.

## Troubleshooting

If fastqutils reports that `dsrc` was not found on `$PATH`, install it
as above or set `DSRC_BIN`. If a `.dsrc` file fails to read or write,
fastqutils surfaces whatever `dsrc` printed to stderr in the error
message.
