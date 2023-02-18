package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
	"log"
	"os"
	"time"
)

const (
	ISOLayout = "2006-01-02"
)

func readParticipants() []string {
	file, err := os.Open("participants")
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

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getFirstDayOfCurrentMonth() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
}

func main() {

	var participants = readParticipants()
	ctx := context.Background()

	token, ok := os.LookupEnv("github_token")
	if !ok {
		fmt.Println("no value for github_token, can not read info. Set env github_token")
		os.Exit(2)
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	// List all issues and PR since october for connected user
	october, _ := time.Parse(ISOLayout, "2019-10-01")
	issuesOptions := github.IssueListOptions{State: "all", Filter: "created", Since: october}
	issues, _, err := client.Issues.List(ctx, true, &issuesOptions)
	if err != nil {
		fmt.Println("Error while reading issues")
	}

	for i, issue := range issues {
		fmt.Printf("%v. %q %q \n", i+1, issue.GetTitle(), issue.GetState())
	}

	// list all repositories for each user

	for i, participant := range participants {
		fmt.Printf("Reading for participant %s (%v)\n", participant, i+1)
		// repositories, _, err := client.Repositories.List(ctx, participant, nil)
		// if err != nil {
		// 	fmt.Println("Error while reading repositories")
		// }
		user, _, err := client.Users.Get(ctx, participant)
		if err != nil {
			fmt.Println("Error while reading user")
		}
		fmt.Printf("%v login %s name %s repo %d\n", i+1, *user.Login, *user.Name, *user.PublicRepos)
		// for i, repository := range repositories {
		// 	fmt.Printf("%v. %v \n", i+1, repository.GetFullName())
		// 	if repository.GetHasIssues() {
		// 		fmt.Printf("Listing issues for repository %v \n", repository.GetFullName())
		// 		issuesOptions := github.IssueListByRepoOptions{State: "all", Assignee: "*"}
		// 		issues, _, err := client.Issues.ListByRepo(ctx, participant, repository.GetName(), &issuesOptions)
		// 		if err != nil {
		// 			fmt.Println("Error while reading issues")
		// 		}
		//
		// 		for i, issue := range issues {
		// 			fmt.Printf("%v. %q\n", i+1, issue.GetTitle())
		// 		}
		// 	}
		// }
		var allCommits []*github.Commit
		opt := &github.SearchOptions{
			ListOptions: github.ListOptions{PerPage: 10},
		}
		firstDay := getFirstDayOfCurrentMonth()
		for {
			commitResults, response, err := client.Search.Commits(ctx, "user:"+*user.Login+" committer-date:>"+firstDay.Format("2006-01-02"), opt)
			if err != nil {
				fmt.Println("Error while reading commits", err)
				break
			}
			if *commitResults.Total > 0 {
				fmt.Printf("found %d commits \n", *commitResults.Total)
				for _, commitResult := range commitResults.Commits {
					//fmt.Printf("CommitResult : %v \n",commitResult)
					fmt.Printf("CommitResult : %v %v %v\n", *commitResult.Author.Login, *commitResult.Repository.Name, *commitResult.Score)
					allCommits = append(allCommits, commitResult.Commit)
				}
			}

			if response.NextPage == 0 {
				break
			}
			opt.Page = response.NextPage
		}
		fmt.Printf("final list of commit for %v is %v \n", *user.Name, len(allCommits))
	}
}
