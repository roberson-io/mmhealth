# Healthcheck tool

Simple Mattermost health check tool. This tool accepts a support packet and generates a markdown file containing the results of the investigation.

## Supported Mattermost Versions

This tool supports Mattermost support packets from the following versions:

**Currently Supported (as of October 2025):**
- **11.0.x** - Full support (until 2026-01-15)
- **10.12.x** - Full support (until 2025-12-15)
- **10.11.x** - Extended Support Release (ESR) (until 2026-08-15)

**Legacy Support:**
- **v9.10 - v10.10** - Older versions are supported via V1 format compatibility

The tool automatically detects and handles both V1 (pre-v10.6.0) and V2 (v10.6.0+) support packet formats.

> **Note:** Mattermost releases new versions monthly. The supported versions in `mmhealth/files/config.yaml` should be updated regularly by checking the [Mattermost Server Releases](https://docs.mattermost.com/product-overview/mattermost-server-releases.html) documentation.

## Getting Started

1. Ensure you have `go` installed on your OS. If you do not, follow the guide [here](https://go.dev/doc/install)

2. Download the `mmhealth` tool with `go install`

    ```bash
    go install github.com/coltoneshaw/mmhealth@latest
    ```

3. Ensure you have Docker setup and running on your computer.

## How to use

The easiest way to use this tool is to just point it at a support packet and generate the report, like below.

```bash
mmhealth generate -p ./mattermost_support_packet_2023-09-21-11-55.zip
```

This will output a `healthcheck-report.pdf` into your current directory.

You can modify the report name by passing `--outputName` or `-o`.

## Legend

### Statuses

- `pass` - The Check passed
- `warn` - The check did not pass, but it is just a warning right now. All adoption / proactive checks that fail are a warn.
- `fail` - The check failed and it should be addressed
- `ignore` - The check can be ignored.

### Severity

- `urgent` - Highest priority check and should be addressed immediately.
- `high` - The issue should be scheduled in the next change window.
- `medium` - Can be addressed later, also informational.
- `low` - Can be ignored, but fixing could provide benefits.

### Types

- `Proactive` - A proactive measure to increase the health, reliability, or otherwise inside of Mattermost.
- `Adoption` - Better configuration of Mattermost for optimal usage and adoption.
- `Health` - Environment health checks that should be remediated if failed.

## Contributing

### Setup the dev environment

1. Clone the repo
2. Start docker desktop
3. Run `make buildDocker` to generate the local docker image
4. Run `make build` to build the binary
5. Run `make test` to confirm everything is working

Now to use the tool you can run `./bin/mmhealth` and use the newly updated code. 

### Updating / adding the checks

> **Note:** This process is a bit clunky right now and needs to be improved. 
> The easiest way is to use `./bin/mmhealth add`, go through the automated process, 
> and then add it to the sheet in the correct format. This way gives it a unique ID that's next in line and properly keeps it in sync.


The checks are provided from the Findings Report, parsed via CSV and updated in the `./mmhealth/files/checks.yaml` file.

To parse a csv:

1. Add a check on the [Findings report](https://docs.google.com/spreadsheets/d/1biFuKKgjhAYi7bKyknoo3h4bz9jpb44oWZ1bRPjwicI/edit#gid=0)
   - Note that you need to give the check a unique ID for it to be parsed later on.
2. Download the finding report from Google as a csv
3. Copy to the root dir here and name it `healthcheck.csv`.
4. Run `make build`
5. Run `./bin/mmhealth parsecsv`
6. Check the git diff for exactly what's changed. You should see the new checks here. 
7. Run `make test` to confirm no checks have broke
8. Update or add the checks broken checks

### Modifying the template

When you modify the template it requires a rebuild of docker. 

1. Make changes to the template
2. run `make buildDocker` to rebuild the docker image
3. run `./bin/mmhealth generate` to use the new template

Once you are done, be sure to follow the docker release info below. 

### Making Docker Release

When you need to adjust anything inside the `./docker/dockerfile` or `./template`, you'll need to manually adjust the `DOCKER_VERSION` to trigger a rebuild. 

### Making mmhealth Release

Just release by adding a tag on github and publishing, everything else is handled.

