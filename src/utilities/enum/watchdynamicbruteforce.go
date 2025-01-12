package watch

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/bahramiofficial/watchtower/src/api/model"
	"github.com/bahramiofficial/watchtower/src/database"
	"github.com/bahramiofficial/watchtower/src/utilities"
)

func DynamicBrute(domain string) {
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
		fmt.Printf("%v running DynamicBrute on %v\n", utilities.GetFormattedTime(), domain)
		results, err = RunDynamicBrute(domain)
		if err != nil {
			fmt.Printf("Failed to run DynamicBrute: %v", err)
		}
		// Check if there are any results
		if results == nil || len(*results) == 0 {
			fmt.Println("No results found.")
			return
		}

		for _, subdomain := range *results {
			if subdomain != "" {
				model.UpsertSubdomain(db, program.ProgramName, subdomain, "dynamicbrute")

				model.UpsertLiveSubdomain(db, program.ProgramName, subdomain, domain, nil, "")
				// upsert_lives(domain=domain, subdomain=sub, ips=[], tag="")
			}
		}
	} else {

	}
}
func RunDynamicBrute(domain string) (*[]string, error) {
	// Create temporary directory for storing files
	tempDir, err := os.MkdirTemp("", "dynamic_brute")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tempDir) // Clean up temp directory after function returns

	// Paths for temporary files
	dnsBrute := fmt.Sprintf("%s/%s.dns_brute", tempDir, domain)
	dnsGenWords := fmt.Sprintf("%s/dnsgen-words.tx", tempDir)
	altDnsWords := fmt.Sprintf("%s/altdns-words.txt", tempDir)
	mergedPath := fmt.Sprintf("%s/words-merged.tx", tempDir)
	domainDnsGen := fmt.Sprintf("%s/%s.dns_gen", tempDir, domain)

	// Step 1: Prepare wordlist for dynamic brute
	commands := []string{
		fmt.Sprintf("curl -s https://raw.githubusercontent.com/AlephNullSK/dnsgen/master/dnsgen/words.txt -o %s", dnsGenWords),
		fmt.Sprintf("curl -s https://raw.githubusercontent.com/infosec-au/altdns/master/words.txt -o %s", altDnsWords),
		fmt.Sprintf("cat %s %s | sort -u > %s", dnsGenWords, altDnsWords, mergedPath),
	}
	for _, cmd := range commands {
		fmt.Printf("Executing command: %s\n", cmd)
		_, err := runCommandInZsh(cmd)
		if err != nil {
			return nil, err
		}
	}

	db, closeDb, err := database.GetDbAfterInit()
	if err != nil {
		return nil, err
	}
	defer closeDb()

	// Step 2: Get subdomains for dynamic brute
	subdomains, err := model.GetAllSubdomainWithScope(db, domain)
	if err != nil {
		return nil, err
	}

	// Write subdomains to dns_brute file
	dnsBruteFile, err := os.Create(dnsBrute)
	if err != nil {
		fmt.Printf("filad create dns_brute file")

		return nil, err
	}
	defer dnsBruteFile.Close()

	for _, sub := range subdomains {
		_, err := dnsBruteFile.WriteString(fmt.Sprintf("%s\n", sub.SubDomain))
		if err != nil {
			fmt.Printf("Failed to write subdomain to file")
			return nil, err
		}
	}

	// Step 3: Run dnsgen command
	command := fmt.Sprintf("cat %s | dnsgen -w %s - | tee %s", dnsBrute, mergedPath, domainDnsGen)
	dnsGenResult, err := runCommandInZsh(command)
	if err != nil {
		fmt.Printf("Error running dnsgen command")
		return nil, err
	}
	// Check if the file is too large and re-run with live subdomains
	if len(dnsGenResult) > 30000000 {

		// Here you would fetch live subdomains (this part is simplified)
		// For now, just simulate it with the getSubdomains function again
		liveSubdomains, err := model.GetAllLiveSubdomainWithScope(db, domain)
		if err != nil {
			return nil, err
		}

		// Write live subdomains to dns_brute file again
		dnsBruteFile, err = os.Create(dnsBrute)
		if err != nil {
			fmt.Printf("Failed to create file dnsBrute")
			return nil, err
		}
		defer dnsBruteFile.Close()

		for _, sub := range liveSubdomains {
			_, err := dnsBruteFile.WriteString(fmt.Sprintf("%s\n", sub.SubDomain))
			if err != nil {

				fmt.Printf("Failed to write subdomain to file")
				return nil, err
			}
		}

		// Re-run the dnsgen command
		_, err = runCommandInZsh(command)
		if err != nil {

			fmt.Printf("Error running dnsgen command")
			return nil, err

		}
	}

	// Step 4: Run shuffledns command
	shufflednsCommand := fmt.Sprintf(
		"shuffledns -list %s -d %s -r ~/.resolvers -m $(which massdns) -mode resolve -t 100 -silent", domainDnsGen, domain)
	fmt.Printf("Executing command: %s\n", shufflednsCommand)
	output, err := runCommandInZsh(shufflednsCommand)
	if err != nil {
		return nil, err
	}

	results := strings.Split(string(output), "\n")
	return &results, nil

}

func runCommandInZsh(command string) (string, error) {
	cmd := exec.Command("zsh", "-c", command)
	output, err := cmd.CombinedOutput()
	return string(output), err
}
