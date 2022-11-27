package client

import (
	"SDR_Labo1/src/communication/protocol"
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/term"
)

var messagingProtocol = &protocol.SDRProtocol{}
var consoleIn = bufio.NewReader(os.Stdin)

// userInterface is the main function for the user interface,
// the client can go through each different functionality
func userInterface(c *Connection) {
	fmt.Println("Welcome!")

	var choice int
	for {
		menu := "Choose one of the following functionality\n[1] Create a new event\n[2] Register to an event as a volunteer\n[3] List all current events\n[4] List all jobs for a specific event\n[5] List the volunteer repartition for a specific event\n[6] Close an event\n[7] Quit\n"

		choice = c.integerReader(menu)
		switch choice {
		case 1:
			createEvent(c)
		case 2:
			volunteerRegistration(c)
		case 3:
			printEvents(c)
		case 4:
			listJobs(c)
		case 5:
			volunteerRepartition(c)
		case 6:
			closeEvent(c)
		case 7:
			c.Close()
			return
		default:
			fmt.Println("You have entered a bad request")
		}
		fmt.Println()
	}
}

// loginClient allow the user to log in
// If the login is invalid, the user has to try again
func loginClient(c *Connection) {
	for {
		username := stringReader("Enter your username : ")
		password := passwordReader("Enter your password : ")

		if c.LoginClient(username, password) {
			break
		}
	}
}

// createEvent creates a new event made by an organizer
// The user has to log in and must be an organizer
func createEvent(c *Connection) {
	loginClient(c)

	eventName := stringReader("\nEnter the event name : ")
	fmt.Println("List all job's name followed by the number of volunteers needed\n" +
		"(tap STOP when ended) : ")

	var jobList []string
	jobList = append(jobList, eventName)
	var i = 0
	for {
		i++
		jobName := stringReader("Insert a name for Job " + strconv.Itoa(i) + ": ")
		if strings.Compare(jobName, "STOP") == 0 {
			break
		}
		nbVolunteers := c.integerReader("Number of volunteers needed for this job : ")
		jobList = append(jobList, jobName, fmt.Sprint(nbVolunteers))
	}
	c.CreateEvent(jobList)
}

// printEvents prints all the events
func printEvents(c *Connection) {
	c.PrintEvents()
}

// volunteerRegistration allows a volunteer to register to an event
// The user has to log in
func volunteerRegistration(c *Connection) {
	loginClient(c)

	var eventId int
	var jobId int
	input := stringReader("Enter [event id] [job id] : ")
	_, err := fmt.Sscan(input, &eventId, &jobId)
	if err != nil {
		log.Fatal(err)
	}
	c.VolunteerRegistration(eventId, jobId)
}

// listJobs prints all the jobs for a specific event
func listJobs(c *Connection) {
	eventId := c.integerReader("Enter event id : ")
	c.ListJobs(eventId)
}

// volunteerRepartition prints the volunteer repartition for a specific event
func volunteerRepartition(c *Connection) {
	eventId := c.integerReader("Enter event id : ")
	c.VolunteerRepartition(eventId)
}

// closeEvent closes an event
// The user has to log in and must be the organizer of the event
func closeEvent(c *Connection) {
	loginClient(c)
	eventId := c.integerReader("Enter event id: ")
	c.CloseEvent(eventId)
}

// stringReader reads a string from the console
func stringReader(optionalMessage string) string {
	fmt.Print(optionalMessage)

	message, err := consoleIn.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimRight(message, "\r\n")
}

func passwordReader(optionalMessage string) string {
	fmt.Print(optionalMessage)

	message, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimRight(string(message), "\r\n")
}

// integerReader reads an integer from the console and returns it
func (c *Connection) integerReader(optionalMessage string) int {
	for {
		fmt.Print(optionalMessage)
		n, err := strconv.ParseInt(stringReader(""), 10, 32)
		if errors.Is(err, strconv.ErrSyntax) {
			fmt.Println()
			fmt.Println("Not a number !")
		} else if errors.Is(err, strconv.ErrRange) {
			fmt.Println()
			fmt.Println("Number too big !")
		} else if err != nil {
			log.Fatal(err)
		} else {
			return int(n)
		}
	}
}

// printDataPacket prints the content of a data packet
func printDataPacket(data protocol.DataPacket) {
	for i := 0; i < len(data.Data); i++ {
		fmt.Println(data.Data[i])
	}
}

func printEventTable(events map[string]interface{}) {
	fmt.Printf("%s | %-20s | %-15s | %-6s |\n", "Id", "Name", "Organizer", "Status")
	for _, event := range events {
		eventMap := event.(map[string]interface{})
		openStatus := "open"
		if eventMap["IsOpen"] == "false" {
			openStatus = "closed"
		}
		fmt.Printf("%.0f | %-20s | %-15s | %-6s |\n", eventMap["ID"], eventMap["Name"], eventMap["Organizer"], openStatus)
	}
}

func printJobTAble(jobs map[string]interface{}) {
	fmt.Printf("%s | %-20s | %s |\n", "Id", "Name", "Required")
	for _, job := range jobs {
		jobMap := job.(map[string]interface{})
		fmt.Printf("%.0f | %-20s | %.0f |\n", jobMap["id"], jobMap["name"], jobMap["required"])
	}
}

// GetJobsRepartitionTable returns a table with the jobs and which volunteers are assigned to them
func printJobPrepartitionTable(jobs map[string]interface{}) {
	head := "| Volunteers     | "
	// Sort jobs by id
	var keys []int
	for _, job := range jobs {
		jobMap := job.(map[string]interface{})
		keys = append(keys, int(jobMap["id"].(float64)))
	}
	sort.Ints(keys)

	jobsName := make([]string, 0)
	for _, job := range jobs {
		jobMap := job.(map[string]interface{})
		s := fmt.Sprintf("%0.f: %-10s /%0.f", jobMap["id"], jobMap["name"], (jobMap["required"])) + " | "
		jobsName = append(jobsName, jobMap["name"].(string))
		head += s
	}

	var tab []string
	tab = append(tab, head)
	// Repartition of volunteers
	for k := range keys {
		i := 0
		jobI := jobs[strconv.Itoa(keys[k])].(map[string]interface{})
		fmt.Println("jobI : ", jobI)
		volunteers := jobI["volunteers"].([]interface{})
		for _, volunteer := range volunteers {
			fmt.Println("volunteer: ", volunteer)
			line := fmt.Sprintf("| %-15s", volunteer)
			for i := range keys {
				jobJ := jobs[strconv.Itoa(keys[k])].(map[string]interface{})
				if jobJ["name"] == jobsName[i] {
					line += "|" + fmt.Sprintf("%-9s", "") + "X" + fmt.Sprintf("%-8s", "")
				} else {
					line += "|" + fmt.Sprintf("%-18s", " ")
				}
			}
			i++
			tab = append(tab, line+"|")
		}
	}

	//for _, job := range jobs {
	//	fmt.Println(jobs)
	//	jobMap := job.(map[string]interface{})
	//	volunteers := jobMap["volunteers"].([]interface{})
	//	for _, volunteer := range volunteers {
	//		fmt.Println("volunteer : ", volunteer)
	//		line := fmt.Sprintf("| %-15s", volunteer)
	//		for _, job := range jobs {
	//			jobMap := job.(map[string]interface{})
	//			volunteers := jobMap["volunteers"].([]interface{})
	//			for _, volunteer2 := range volunteers {
	//				fmt.Println("volunteer2 : ", volunteer)
	//				fmt.Println(jobMap)
	//				if volunteer2 == volunteer {
	//					line += "|" + fmt.Sprintf("%-9s", "") + "X" + fmt.Sprintf("%-8s", "")
	//				} else {
	//					line += "|" + fmt.Sprintf("%-18s", " ")
	//				}
	//			}
	//		}
	//		tab = append(tab, line+"|")
	//	}
	//}
	for i := 0; i < len(tab); i++ {
		fmt.Println(tab[i])
	}
}
