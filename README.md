# Trello to Github

Import Trello cards into your GitHub projects.

## Installation

* With Docker, using the [GitHub Package](https://github.com/users/timtraversy/packages/container/package/trellotogithub)
```
docker run -it --rm -v "${PWD}:/local" ghcr.io/timtraversy/trellotogithub:latest \
    -c local/path/to/config.yaml
```

* Download binary for your OS from [Release page](https://github.com/timtraversy/trellotogithub/releases)

* Build from source
    * Dependencies: Go 1.15
```
git clone https://github.com/timtraversy/trellotogithub
cd trellotogithub
go run . 
```

## Usage

To Do

## Contributing

First, fork and clone this repo onto your machine.
   
(Option 1) If you have [VS Code](https://code.visualstudio.com) and [Docker](https://www.docker.com) installed, the simplest way to get going is to open this project in a [VS Code Development Container](https://code.visualstudio.com/docs/remote/containers). Just run `Remote-Containers: Open Folder in Container` in VS Code and select the cloned folder. This wil spin up a Go container with the proper Go version and settings, which you can use to edit and compile the source.

(Option 2) If you have Go installed on your machine, you can edit and compile the source directly.

## Support

Feel free to open issues in the GitHub issue tracker.

## Roadmap

To Do
