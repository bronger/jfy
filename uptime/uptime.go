package uptime

import (
	"errors"
	"regexp"
	"runtime"
	"strconv"

	"github.com/bronger/jfy/lib"
	"github.com/pborman/getopt/v2"
)

func MustAtoi(s string) int {
	if i, err := strconv.Atoi(s); err != nil {
		panic(err)
	} else {
		return i
	}
}

func Handle(settings lib.SettingsType, stdout, stderr []byte, args ...string) (any, any, error) {
	if runtime.GOOS != "linux" {
		return nil, nil, errors.New("Your OS is not yet supported by jfy wrapper")
	}
	commandLine := getopt.New()
	flags := []*bool{commandLine.BoolLong("pretty", 'c'), commandLine.BoolLong("help", 'h'),
		commandLine.BoolLong("since", 's'), commandLine.BoolLong("version", 'V')}
	if err := commandLine.Getopt(args, nil); err != nil {
		return nil, nil, err
	}
	for _, flag := range flags {
		if *flag {
			return nil, nil, errors.New("Command line option to uptime is not yet supported by jfy wrapper")
		}
	}

	regex := regexp.MustCompile(
		`^ (\d+):(\d+):(\d+) up (\d+):(\d+),  (\d+) users,  load average: ([0-9.]+), ([0-9.]+), ([0-9.]+)\n$`)
	submatches := regex.FindStringSubmatch(string(stdout))
	if len(submatches) != 10 {
		return nil, nil, errors.New("Could not parse output of uptime.  Wrong locale?")
	}
	var load1, load5, load15 float64
	var err error
	if load1, err = strconv.ParseFloat(submatches[7], 64); err != nil {
		return nil, nil, errors.New("Could not parse output of uptime.")
	}
	if load5, err = strconv.ParseFloat(submatches[8], 64); err != nil {
		return nil, nil, errors.New("Could not parse output of uptime.")
	}
	if load15, err = strconv.ParseFloat(submatches[9], 64); err != nil {
		return nil, nil, errors.New("Could not parse output of uptime.")
	}
	return map[string]any{"hour": MustAtoi(submatches[1]), "minute": MustAtoi(submatches[2]),
		"second": MustAtoi(submatches[3]), "hours": MustAtoi(submatches[4]), "minutes": MustAtoi(submatches[5]),
		"users": MustAtoi(submatches[6]), "load1": load1, "load5": load5, "load15": load15}, nil, nil
}
