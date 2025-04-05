package phone_db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Phone struct {
	ID     int
	Number string
}

type DB struct {
	db *sql.DB
}

func Open(driver_name, data_source string) (*DB, error) {
	database, err := sql.Open(driver_name, data_source)
	if err != nil {
		return nil, err
	}

	return &DB{database}, nil
}

func (database *DB) Close() error {
	return database.db.Close()
}
func Reset(driver_name, data_source, database_name string) error {
	database, err := sql.Open(driver_name, data_source)
	if err != nil {
		return err
	}

	err = resetDB(database, database_name)
	if err != nil {
		return err
	}

	return database.Close()
}

func Migrate(driver_name, data_source string) error {
	database, err := sql.Open(driver_name, data_source)
	if err != nil {
		return err
	}

	err = createPhoneNumbersTable(database)
	if err != nil {
		return err
	}

	return database.Close()
}

func createPhoneNumbersTable(database *sql.DB) error {
	statement := `
		CREATE TABLE IF NOT EXISTS phone_numbers (
		id SERIAL,
		value VARCHAR(255)
		)`

	_, err := database.Exec(statement)

	return err
}

func createDB(db *sql.DB, name string) error {
	_, err := db.Exec("CREATE DATABASE " + name)
	if err != nil {
		return err
	}
	return nil
}

func resetDB(db *sql.DB, name string) error {
	_, err := db.Exec("DROP DATABASE IF EXISTS " + name)
	if err != nil {
		return err
	}

	return createDB(db, name)
}

func (database *DB) Seed() error {
	data := []string{
		"1234567890",
		"123 456 7891",
		"(123) 456 7892",
		"(123) 456-7893",
		"123-456-7890",
		"1234567892",
		"(123)456-7892",
	}

	for _, number := range data {
		if _, err := insertPhone(database.db, number); err != nil {
			return err
		}
	}
	return nil
}

func insertPhone(database *sql.DB, phone string) (int, error) {
	statement := `INSERT INTO phone_numbers(value) VALUES($1) RETURNING id`

	var id int
	err := database.QueryRow(statement, phone).Scan(&id)
	if err != nil {
		return -1, err
	}

	return id, nil
}

func getPhone(database *sql.DB, id int) (string, error) {

	var number string
	err := database.QueryRow("SELECT * FROM phone_numbers where id=$1", id).Scan(&id, &number)
	if err != nil {
		return "", err
	}

	return number, nil
}

func (database *DB) AllPhones() ([]Phone, error) {
	rows, err := database.db.Query("SELECT id, value from phone_numbers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var return_value []Phone
	for rows.Next() {
		var p Phone
		if err := rows.Scan(&p.ID, &p.Number); err != nil {
			return nil, err
		}

		return_value = append(return_value, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return return_value, nil
}

func (database *DB) FindPhone(number string) (*Phone, error) {

	var phone Phone
	err := database.db.QueryRow("SELECT * FROM phone_numbers where value=$1", number).Scan(&phone.ID, &phone.Number)
	if err != nil {

		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &phone, nil
}

func (database *DB) UpdatePhone(p *Phone) error {
	statement := `UPDATE phone_numbers SET value=$2 where id=$1`
	_, err := database.db.Exec(statement, p.ID, p.Number)
	return err
}

func (database *DB) DeletePhone(id int) error {
	statement := `DELETE FROM phone_numbers WHERE id=$1`
	_, err := database.db.Exec(statement, id)
	return err
}
