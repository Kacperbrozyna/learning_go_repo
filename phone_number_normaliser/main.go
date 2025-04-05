package main

import (
	"fmt"
	"regexp"

	phone_db "github.com/Kacperbrozyna/learning_go_repo/phone_number_normaliser/database"
)

const (
	host         = "localhost"
	port         = 5432
	user         = "postgres"
	password     = "password"
	databaseName = "phone_database"
	driverName   = "postgres"
)

type phone struct {
	id     int
	number string
}

func normalise(phone_number string) string {
	re := regexp.MustCompile("\\D")
	return re.ReplaceAllString(phone_number, "")
}

func main() {
	psql_info := fmt.Sprintf("host=%s port=%d  user=%s sslmode=disable", host, port, user)

	phone_db.Reset(driverName, psql_info, databaseName)

	psql_info = fmt.Sprintf("%s dbname = %s", psql_info, databaseName)

	err := phone_db.Migrate(driverName, psql_info)
	if err != nil {
		panic(err)
	}

	database, err := phone_db.Open(driverName, psql_info)
	if err != nil {
		panic(err)
	}
	defer database.Close()

	err = database.Seed()
	if err != nil {
		panic(err)
	}

	phones, err := database.AllPhones()
	if err != nil {
		panic(err)
	}

	for _, phone := range phones {
		fmt.Printf("Working on... %+v\n", phone)
		normalised_number := normalise(phone.Number)
		if normalised_number != phone.Number {
			existing, err := database.FindPhone(normalised_number)
			if err != nil {
				panic(err)
			}

			if existing != nil {
				fmt.Println("Removing..", normalised_number)
				err = database.DeletePhone(phone.ID)
				if err != nil {
					panic(err)
				}
			} else {
				fmt.Println("Updating..", normalised_number)
				phone.Number = normalised_number
				err = database.UpdatePhone(&phone)
				if err != nil {
					panic(err)
				}
			}

		} else {
			fmt.Println("No changes needed")
		}
	}
}
