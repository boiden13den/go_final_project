package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "20060102"

// nextDayHandler parses the query parameters and writes the result of NextDate to the response.
func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
		return
	}
	nowStr := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	now := time.Now()
	if nowStr != "" {
		var err error
		now, err = time.Parse(DateFormat, nowStr)
		if err != nil {
			http.Error(w, "invalid now date", http.StatusBadRequest)
			return
		}
	}

	result, err := NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprint(w, result)
}

// NextDate calculates the next date after `now` based on the initial `date` and `repeat` rules.
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("repeat is empty")
	}
	date, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("invalid date: %w", err)
	}

	parts := strings.SplitN(repeat, " ", 2)
	switch parts[0] {
	case "d":
		if len(parts) < 2 {
			return "", errors.New("d: missing interval")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil || days < 1 || days > 400 {
			return "", errors.New("d: invalid interval")
		}
		for {
			date = date.AddDate(0, 0, days)
			if date.After(now) {
				break
			}
		}

	case "y":
		for {
			date = date.AddDate(1, 0, 0)
			if date.After(now) {
				break
			}
		}

	case "w":
		if len(parts) < 2 {
			return "", errors.New("w: missing days")
		}
		var allowed [8]bool
		for _, s := range strings.Split(parts[1], ",") {
			d, err := strconv.Atoi(s)
			if err != nil || d < 1 || d > 7 {
				return "", errors.New("w: invalid day")
			}
			allowed[d] = true
		}
		for {
			date = date.AddDate(0, 0, 1)
			wd := int(date.Weekday())
			if wd == 0 {
				wd = 7
			}
			if allowed[wd] && date.After(now) {
				break
			}
		}

	case "m":
		if len(parts) < 2 {
			return "", errors.New("m: missing days")
		}
		monthParts := strings.SplitN(parts[1], " ", 2)
		var dayAllowed [32]bool
		for _, s := range strings.Split(monthParts[0], ",") {
			d, err := strconv.Atoi(s)
			if err != nil || d < -2 || d > 31 || d == 0 {
				return "", errors.New("m: invalid day")
			}
			if d > 0 {
				dayAllowed[d] = true
			}
		}
		var monthAllowed [13]bool
		if len(monthParts) > 1 {
			for _, s := range strings.Split(monthParts[1], ",") {
				m, err := strconv.Atoi(s)
				if err != nil || m < 1 || m > 12 {
					return "", errors.New("m: invalid month")
				}
				monthAllowed[m] = true
			}
		} else {
			for i := 1; i <= 12; i++ {
				monthAllowed[i] = true
			}
		}
		for {
			date = date.AddDate(0, 0, 1)
			if !monthAllowed[int(date.Month())] {
				continue
			}
			last := lastDay(date)
			day := date.Day()
			if dayAllowed[day] ||
				(containsNeg(parts[1], -1) && day == last) ||
				(containsNeg(parts[1], -2) && day == last-1) {
				if date.After(now) {
					break
				}
			}
		}

	default:
		return "", fmt.Errorf("unsupported repeat format: %s", parts[0])
	}

	return date.Format(DateFormat), nil
}

// lastDay returns the number of the last day of the month for the given time
func lastDay(t time.Time) int {
	return time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, t.Location()).Day()
}

// containsNeg reports whether s contains the specified negative number.
func containsNeg(s string, n int) bool {
	return strings.Contains(s, strconv.Itoa(n))
}
