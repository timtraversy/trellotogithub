module github.com/timtraversy/trellotogithub

go 1.15

require (
	github.com/google/go-github/v32 v32.1.0
	github.com/timtraversy/trellotogithub/githubdeviceauth v0.0.0-20201126204849-5ac8fb4bd0e6
)

replace github.com/timtraversy/trellotogithub/githubdeviceauth => ./githubdeviceauth
