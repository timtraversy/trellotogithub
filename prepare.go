package main

import (
	"context"
	"fmt"
	"io"

	"github.com/adlio/trello"
	"github.com/google/go-github/v32/github"
)

func selectBoard(client *trello.Client, in io.Reader, out io.Writer) *trello.Board {
	boards, _ := client.GetMyBoards(trello.Defaults())
	fmt.Fprintln(out, "Select the Trello board you want to import from:")
	for i, board := range boards {
		fmt.Fprintf(out, "(%v) %v\n", i, board.Name)
	}

	var selection int
	fmt.Fscanf(in, "%d", selection)
	return boards[selection]
}

type repositoriesClient interface {
	List(ctx context.Context, user string, opts *github.RepositoryListOptions) ([]*github.Repository, *github.Response, error)
}

type usersClient interface {
	Get(ctx context.Context, user string) (*github.User, *github.Response, error)
	ListProjects(ctx context.Context, user string, opts *github.ProjectListOptions) ([]*github.Project, *github.Response, error)
}

func selectRepository(client repositoriesClient, in io.Reader, out io.Writer) *github.Repository {
	repos, _, _ := client.List(context.Background(), "", nil)
	fmt.Fprintln(out, "Select the GitHub repository whose projects you want to export to:")
	for i, repo := range repos {
		fmt.Fprintf(out, "(%v) %v\n", i, *repo.Name)
	}

	var selection int
	fmt.Fscanf(in, "%d", &selection)
	return repos[selection]
}

func selectProject(client usersClient, in io.Reader, out io.Writer) *github.Project {
	user, _, _ := client.Get(context.Background(), "")
	projects, _, _ := client.ListProjects(context.Background(), user.GetName(), nil)
	fmt.Fprintln(out, "Select the GitHub project you want to export to:")
	for i, project := range projects {
		fmt.Fprintf(out, "(%v) %v\n", i, *project.Name)
	}

	var selection int
	fmt.Fscanf(in, "%d\n", &selection)
	return projects[selection]
}
