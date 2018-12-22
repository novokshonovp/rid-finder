package main

import (
  "bufio"
  "fmt"
  "log"
  "os"
  "time"
  "flag"
  "io"
  "regexp"
)

func timeTrack(start time.Time, name string) {
  elapsed := time.Since(start)
  log.Printf("%s took %s", name, elapsed)
}

func showKeys() {
  help := `usage: | rid-finder -f file_with_rids [-rse] [-regexp]
  # usage example
    $ cat path_to_file_with_logs | rid-finder -f path_to_file_with_rids -r RID > path_to_output_file
    $ cat path_to_file_with_logs | rid-finder -f path_to_file_with_rids -r RID -regexp > path_to_output_file
    $ cat path_to_file_with_logs | rid-finder -f path_to_file_with_jids -r JID -regexp > path_to_output_file
    `

  fmt.Println(help)
}

func readLines(path string) (map[string]int, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    lines := make(map[string]int)
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        lines[scanner.Text()] = 1
    }
    return lines, scanner.Err()
}

func get_shift_positions(startingPosition int, endingPosition int, formatPtr string) (int, int) {
    mask_starting_position := 0
    mask_ending_position := 0

    if endingPosition == 0 {
        mask_starting_position = 31
        mask_ending_position = 67

        if formatPtr == "JID" {
            mask_starting_position = 30
            mask_ending_position = 58
        }
    } else {
        mask_starting_position = startingPosition
        mask_ending_position = endingPosition
    }
    return mask_starting_position, mask_ending_position
}

func get_regexp(formatPtr string)(*regexp.Regexp) {
    reg, _ := regexp.Compile("RID-[0-9a-z]{32}")

    if formatPtr == "JID" {
        reg, _ = regexp.Compile("JID-[0-9a-z]{24}")
    }
    return reg
}

func main() {
    idsFilePtr := flag.String("f", "", "a string")
    formatPtr := flag.String("r", "RID", "a string")
    helpPtr := flag.Bool("h", false, "a bool")
    regexpPtr := flag.Bool("regexp", false, "a bool")
    startingPosition := flag.Int("s", 0, "an int")
    endingPosition := flag.Int("e", 0, "an int")

    flag.Parse()

    if len(os.Args) < 2 || *helpPtr || *idsFilePtr == "" { showKeys(); return }

    mask_starting_position, mask_ending_position := get_shift_positions(*startingPosition, *endingPosition, *formatPtr)
    reg := get_regexp(*formatPtr)

    defer timeTrack(time.Now(), "Execution")

    lines, err := readLines(*idsFilePtr)

    if err != nil {
      log.Fatalf("readLines: %s", err)
    }

    stdinReader := bufio.NewReader(os.Stdin)

    rid := ""
	for {
		line_to_analyze, err := stdinReader.ReadString('\n')
		if err != nil && err == io.EOF { break }

        if len(line_to_analyze) < 68 { continue }

        if *regexpPtr {
            rid = reg.FindString(line_to_analyze)
        } else {
            rid = line_to_analyze[mask_starting_position:mask_ending_position]
        }

        if lines[rid] == 1 {
            fmt.Print(line_to_analyze)
        }
	}
}