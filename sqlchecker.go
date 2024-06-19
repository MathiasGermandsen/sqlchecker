package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// Structure to hold the connection string from JSON
type Config struct {
	ConnectionString string `json:"connectionString"`
}

// Function to read connection string from JSON file
func readConnectionString(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	var config Config
	if err := json.Unmarshal(byteValue, &config); err != nil {
		return "", err
	}

	return config.ConnectionString, nil
}

func main() {
	// Read connection string from JSON file
	connStr, err := readConnectionString("conString.json")
	if err != nil {
		log.Fatal(err)
	}

	// Open a connection to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Check if the connection is successful
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully connected to the database")

	// Name of the table to check
	tableName := "employees"

	// Check if table exists
	tableExists, err := checkTableExists(db, tableName)
	if err != nil {
		log.Fatal(err)
	}

	// If table does not exist, create it
	if !tableExists {
		err := createTable(db, tableName)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Table %s created successfully\n", tableName)

		// Insert 50 employees
		err = insertEmployees(db)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("50 employees inserted successfully")
	} else {
		fmt.Printf("Table %s already exists\n", tableName)
	}
}

// Function to check if a table exists
func checkTableExists(db *sql.DB, tableName string) (bool, error) {
	query := fmt.Sprintf("SELECT to_regclass('%s');", tableName)
	var result sql.NullString
	err := db.QueryRow(query).Scan(&result)
	if err != nil {
		return false, err
	}
	return result.Valid, nil
}

// Function to create a table
func createTable(db *sql.DB, tableName string) error {
	createTableSQL := fmt.Sprintf(`CREATE TABLE %s (
		emp_id SERIAL PRIMARY KEY,
		firstname VARCHAR(50),
		lastname VARCHAR(50),
		email VARCHAR(100),
		news BOOLEAN,
		role VARCHAR(50)
	);`, tableName)

	_, err := db.Exec(createTableSQL)
	if err != nil {
		return err
	}
	return nil
}

// Function to insert 50 employees
func insertEmployees(db *sql.DB) error {
	firstnames := []string{"John", "Jane", "Alice", "Bob", "Charlie"}
	lastnames := []string{"Doe", "Smith", "Johnson", "Williams", "Brown"}
	roles := []string{"Janitor", "SoftwareDev", "AiEthics", "ItConsultant", "CleaningLady"}

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 50; i++ {
		firstname := firstnames[rand.Intn(len(firstnames))]
		lastname := lastnames[rand.Intn(len(lastnames))]
		email := fmt.Sprintf("%s%s@techtech.com", firstname, lastname)
		role := roles[rand.Intn(len(roles))]

		insertSQL := `INSERT INTO employees (firstname, lastname, email, news, role) VALUES ($1, $2, $3, $4, $5)`
		_, err := db.Exec(insertSQL, firstname, lastname, email, false, role)
		if err != nil {
			return err
		}
	}

	return nil
}
