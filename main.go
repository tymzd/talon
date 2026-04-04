/*
Command talon is a CLI tool for fetching and formatting workout data from the Hevy API.

This is meant to provide an OpenClaw-friendly Markdown digests of your workouts
and routines on Hevy.

You must set the HEVY_API_KEY environment variable before execution.

Usage:

	talon <command> [flags] [arguments]

The commands are:

	workouts    Fetch workout history (defaults to 10 most recent)
	routines    List all available routines in Hevy
	routine     Get detailed markdown for a specific routine

Examples:

Fetch the 5 most recent workouts:

	talon workouts -n 5

Fetch all workouts from a specific date range:

	talon workouts -from 2026-03-01 -to 2026-03-31

Fetch detailed markdown for a specific routine by name:

	talon routine "Core Progression"

Environment Variables:

	HEVY_API_KEY    Required. The API key for authenticating with Hevy.
*/
package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	flag.Usage = printGlobalHelp
	flag.Parse()

	// Flag sets per subcommand.
	workoutsCmd := flag.NewFlagSet("workouts", flag.ExitOnError)
	routineCmd := flag.NewFlagSet("routine", flag.ExitOnError)

	// Flags specific to workouts
	num := workoutsCmd.Int("n", 10, "Number of the most recent workouts to fetch")
	fromDate := workoutsCmd.String("from", "", "Start date (YYYY-MM-DD)")
	toDate := workoutsCmd.String("to", "", "End date (YYYY-MM-DD)")

	if len(os.Args) < 2 {
		printGlobalHelp()
		os.Exit(1)
	}

	apiKey := os.Getenv("HEVY_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "Error: HEVY_API_KEY environment variable is not set.")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "workouts":
		fmt.Println("Workouts")
		// TODO: Fetch n workouts.
		// TODO: Fetch from-to workouts.
		workoutsCmd.Parse(os.Args[2:])
		fmt.Printf("Num: %d\n", *num)
		fmt.Printf("From: %v\n", *fromDate)
		fmt.Printf("To: %v\n", *toDate)
	case "routines":
		fmt.Println("routines")
		// TODO: Fetch all routines and dump them.
	case "routine":
		fmt.Println("routine")
		// TODO: Fetch specific routine by name
		routineCmd.Parse(os.Args[2:])
		routineName := routineCmd.Arg(0)
		fmt.Printf("Routine: %s\n", routineName)
	default:
		printGlobalHelp()
	}
}

func printGlobalHelp() {
	helpText := `Talon: The OpenClaw Hevy API Fetcher

Usage:
  talon <command> [arguments] [flags]

Commands:
  workouts    Fetch workout history (defaults to 10 most recent)
  routines    List all available routines in Hevy
  routine     Get detailed markdown for a specific routine

Flags:
  -h, --help  Show this help message

Use "talon <command> -h" for more information about a given command.
`
	fmt.Fprint(os.Stderr, helpText)
}
