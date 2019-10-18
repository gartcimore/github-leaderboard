package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/google/go-github/v28/github"
	"golang.org/x/oauth2"
	"log"
	"os"
)

func readParticipants() []string {
	file, err := os.Open("members")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var participants []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		participants = append(participants, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return participants
}
func main() {

	var participants = readParticipants()
	ctx := context.Background()
	token := os.Getenv("github_token")
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	// list all repositories for each user

	for i, user := range participants {
		fmt.Printf("Reading for user %s (%v)\n", user, i)
		repositories, _, err := client.Repositories.List(ctx, user, nil)
		if err != nil {
			fmt.Println("Error while reading repositories")
		}

		for i, repository := range repositories {
			fmt.Printf("%v. %v\n", i+1, repository.GetFullName())
		}

		fmt.Printf("user %v has %v repos\n", user, len(repositories))

		//issuesOptions := github.IssueListOptions{Filter: user}
		issues, _, err := client.Issues.List(ctx, true, nil)
		if err != nil {
			fmt.Println("Error while reading issues")
		}

		for i, issue := range issues {
			fmt.Printf("%v. %v\n", i+1, issue.Title)
		}
	}

}
