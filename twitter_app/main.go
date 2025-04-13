package main

import (
	"MessingWithGo/Learning/learning_go_repo/twitter_app/twitter"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {
	var (
		keyFile    string
		usersFile  string
		tweetId    string
		numWinners int
	)

	flag.StringVar(&keyFile, "key", "keys.json", "the file where you store the api key and secret for X api")
	flag.StringVar(&usersFile, "users", "users.csv", "the file where users who have retweeted the tweet are stored")
	flag.StringVar(&tweetId, "tweet", "375655219813560320", "the id of the tweet you wish to find the retweeters of")
	flag.IntVar(&numWinners, "winners", 3, "how many winners to pick")
	flag.Parse()

	key, secret, err := keys(keyFile)
	if err != nil {
		panic(err)
	}

	client, err := twitter.New(key, secret)
	if err != nil {
		panic(err)
	}

	newUsernames, err := client.Retweeters(tweetId)
	if err != nil {
		panic(err)
	}

	exisitngUsernames := exisiting(usersFile)

	allUsernames := merge(newUsernames, exisitngUsernames)

	err = writeUsers(usersFile, allUsernames)
	if err != nil {
		panic(err)
	}

	if numWinners < 1 {
		return
	}

	exisitngUsernames = exisiting(usersFile)

	winners := pickWinners(allUsernames, numWinners)
	fmt.Println("The winners are:")
	for _, winner := range winners {
		fmt.Printf("\t%s\n", winner)
	}

}

func keys(keyFile string) (key, secret string, err error) {
	var keys struct {
		Key    string `json:"consumer_key"`
		Secret string `json:"consumer_secret"`
	}

	f, err := os.Open(keyFile)
	if err != nil {
		return "", "", err
	}

	defer f.Close()
	dec := json.NewDecoder(f)
	dec.Decode(&keys)

	return keys.Key, keys.Secret, nil
}

func exisiting(usersFile string) []string {
	f, err := os.Open(usersFile)
	if err != nil {
		return []string{}
	}

	defer f.Close()

	r := csv.NewReader(f)
	lines, err := r.ReadAll()
	users := make([]string, 0, len(lines))
	for _, line := range lines {
		users = append(users, line[0])
	}

	return users
}

func merge(a, b []string) []string {
	uniq := make(map[string]struct{}, 0)

	for _, user := range a {
		uniq[user] = struct{}{}
	}

	for _, user := range b {
		uniq[user] = struct{}{}
	}

	ret := make([]string, 0, len(uniq))
	for user := range uniq {
		ret = append(ret, user)
	}

	return ret
}

func writeUsers(userFile string, users []string) error {
	f, err := os.OpenFile(userFile, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	defer f.Close()

	w := csv.NewWriter(f)
	for _, username := range users {
		if err := w.Write([]string{username}); err != nil {
			return err
		}
	}

	w.Flush()

	if err := w.Error(); err != nil {
		return err
	}

	return nil
}

func pickWinners(users []string, numWinners int) []string {
	r := rand.New(rand.NewSource(time.Now().Unix()))

	perm := r.Perm(len(users))
	winners := perm[:numWinners]

	var ret = make([]string, 0, numWinners)
	for _, idx := range winners {
		ret = append(ret, users[idx])
	}

	return ret
}
