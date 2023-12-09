# Raccoon!

A stupid simple download manager that I wrote to learn Go concepts better.

## Usage

Build: `make build`

Download a single file: `raccoon url <file-url>`

Batch downloads: `raccoon readfile <path-to-file>`

Help: `raccoon help`

Help for subcommands: `raccoon <subcommand> --help`

Shell completion: `raccoon completion`

Number of connections and download path can be configured with `--connections` or `-c` and `--directory` or `-d` respectively.