package watch

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strings"

	_ "github.com/lib/pq" // PostgreSQL driver

	"github.com/bahramiofficial/watchtower/src/api/model"
	"github.com/bahramiofficial/watchtower/src/database"
	"github.com/bahramiofficial/watchtower/src/utilities"
	// Import SQLite3 driver
)

func Crtsh(domain string) {

	// Get database connection and the deferred CloseDb function
	db, closeDb, err := database.GetDbAfterInit()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer closeDb() // Ensure that the connection will be closed when the function exits

	var program *model.Program
	var results *[]string

	// Define the domain
	program, _ = model.FindDomainWithProgramName(db, domain)
	if program != nil {
		fmt.Printf("%v running crtsh on %v\n", utilities.GetFormattedTime(), domain)
		results, err = RunCrtsh(domain)
		if err != nil {
			fmt.Printf("Failed to run crtsh: %v", err)
		}
		// Check if there are any results
		if results == nil || len(*results) == 0 {
			fmt.Println("No results found.")
			return
		}

		for _, subdomain := range *results {
			if subdomain != "" {
				model.UpsertSubdomain(db, program.ProgramName, subdomain, "crtsh")
			}
		}
	} else {

	}

}

// func RunCrtsh(domain string, outputFile string) (*[]string, error) {

func RunCrtsh(domain string) (*[]string, error) {
	// Connect to PostgreSQL
	connStr := "user=guest host=crt.sh port=5432 dbname=certwatch sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Prepare the query
	query := fmt.Sprintf(`
		SELECT ci.NAME_VALUE
		FROM certificate_and_identities ci
		WHERE plainto_tsquery('certwatch', '%s') @@ identities(ci.CERTIFICATE)
	`, domain)

	// Execute the query
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Process results
	var nameValue string
	var results []string
	for rows.Next() {
		if err := rows.Scan(&nameValue); err != nil {
			return nil, err
		}

		// Clean the result

		nameValue = strings.ReplaceAll(nameValue, " ", "")
		nameValue = strings.ToLower(nameValue)
		nameValue = strings.TrimPrefix(nameValue, "*.")

		// Domain matching
		matchPattern := fmt.Sprintf(".*\\.%s$", domain)

		if match, _ := regexp.MatchString(matchPattern, nameValue); match {
			results = append(results, nameValue)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	// Remove duplicates and sort the results
	uniqueResults := uniqueSort(results)

	return &uniqueResults, nil
}

// uniqueSort removes duplicates and sorts the results
func uniqueSort(input []string) []string {
	// Use a map to remove duplicates
	unique := make(map[string]struct{})
	for _, item := range input {
		unique[item] = struct{}{}
	}

	// Convert map keys to a slice
	var sortedResults []string
	for key := range unique {
		sortedResults = append(sortedResults, key)
	}

	// Sort the results
	sort.Strings(sortedResults)
	return sortedResults
}
