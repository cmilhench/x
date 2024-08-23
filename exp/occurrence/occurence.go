package occurrence

import (
	"slices"
	"time"

	"github.com/cmilhench/x/exp/ptr"
)

func NextMinutelyOccurrence(current time.Time, interval int, daysOfWeek []time.Weekday, between []time.Time) time.Time {
	// every (N) minutes (between A and B) on Monday, Tuesday and Wednesday until ...
	// next occurrence is in interval minutes as long as within range and on next weekday
	found := current.Add(time.Duration(interval) * time.Minute)
	if len(between) > 0 {
		if !timeIsBetween(found, between[0], between[1]) {
			found = found.AddDate(0, 0, 1)
			found = time.Date(found.Year(), found.Month(), found.Day(), between[0].Hour(), between[0].Minute(), 0, 0, time.UTC)
		}
	}
	if len(daysOfWeek) > 0 {
		for slices.Index(daysOfWeek, found.Weekday()) == -1 {
			found = found.AddDate(0, 0, 1)
			if len(between) > 0 {
				found = time.Date(found.Year(), found.Month(), found.Day(), between[0].Hour(), between[0].Minute(), 0, 0, time.UTC)
			} else {
				found = time.Date(found.Year(), found.Month(), found.Day(), 0, 0, 0, 0, time.UTC)
			}
		}
	}
	current = found
	return current
}

func timeIsBetween(t, min, max time.Time) bool {
	t = time.Date(0, time.January, 1, t.Hour(), t.Minute(), 0, 0, time.UTC)
	min = time.Date(0, time.January, 1, min.Hour(), min.Minute(), 0, 0, time.UTC)
	max = time.Date(0, time.January, 1, max.Hour(), max.Minute(), 0, 0, time.UTC)

	if min.After(max) {
		min, max = max, min
	}
	return (t.Equal(min) || t.After(min)) && (t.Equal(max) || t.Before(max))
}

func NextDailyOccurrence(current time.Time, interval int) time.Time {
	// at (T) every (N) days until ...
	// next occurrence is in interval days
	current = current.AddDate(0, 0, interval)
	return current
}

func NextWeeklyOccurrence(current time.Time, interval int, daysOfWeek []time.Weekday) time.Time {
	if len(daysOfWeek) > 0 {
		// at (T) every (N) week(s) on Thursday, Friday and Saturday until ...
		// next occurrence is in interval weeks on next weekday (maybe less then n*7 if there is a day this week)
		found := time.Time{}
		// floor the date to the start of the week (Sunday)
		first := current.AddDate(0, 0, (int(time.Sunday)-int(current.Weekday())-7)%7)
		// work out the dates for each weekday in week 0
		// find the first date this week that is after current
		for _, d := range daysOfWeek {
			candidate := getNextWeekday(first, d)
			if candidate.After(current) {
				found = candidate
				break
			}
		}
		// if there are no dates found this week add interval to the first date
		if found.IsZero() {
			first = first.AddDate(0, 0, interval*7)
			found = getNextWeekday(first, daysOfWeek[0])
		}
		// set the result to current
		current = found
	} else {
		// at (T) every (N) week(s) until ...
		// next occurrence is in interval weeks
		current = current.AddDate(0, 0, 7*interval)
	}
	return current
}

func getNextWeekday(current time.Time, weekday time.Weekday) time.Time {
	// thursday to sunday = (0 - 4 + 7) % 7 = 3
	// wednesday to monday = (1 - 3 + 7) % 7 = 5
	offset := (int(weekday) - int(current.Weekday()) + 7) % 7
	current = current.AddDate(0, 0, offset)
	return current
}

func getNthWeekday(current time.Time, nth int, weekday time.Weekday) time.Time {
	offset := (int(weekday) - int(current.Weekday()) + 7) % 7
	offset += (nth - 1) * 7
	current = current.AddDate(0, 0, offset)
	return current
}

func NextMonthlyOccurrence(current time.Time, interval int, daysOfWeek []time.Weekday, weekOfMonth, dayOfMonth *int32) time.Time {
	if weekOfMonth == nil {
		weekOfMonth = ptr.Int32(1)
	}
	switch {
	case len(daysOfWeek) == 1:
		// at (T) every (N) month(s) on the (Nth) Sunday until ...
		// next occurrence is on nth weekday in interval months (or this month if date not passed)
		found := getNthWeekdayOfMonth(current, current.Month(), daysOfWeek[0], int(*weekOfMonth))
		if found.After(current) {
			current = found
		} else {
			current = getNthWeekdayOfMonth(current, time.Month(int(current.Month())+interval), daysOfWeek[0], int(*weekOfMonth))
		}
	case dayOfMonth != nil && len(daysOfWeek) == 0:
		// at (T) every (N) month(s) on the (Nth) until ...
		// next occurrence is on the same date in interval months (or this month if date not passed)
		found := time.Date(current.Year(), current.Month(), int(*dayOfMonth), current.Hour(), current.Minute(), current.Second(), 0, time.UTC)
		if found.After(current) {
			current = found
		} else {
			current = current.AddDate(0, interval, 0)
		}
	default:
		// at (T) every (N) month(s)
		// next occurrence is in interval days
		current = current.AddDate(0, interval, 0)
	}
	// Invalid/Unimplemented cases;
	// - dayOfMonth == nil && daysOfWeek != 1
	// at (T) every (N) month(s) on the Nth weekday that's a ?
	// at (T) every (N) month(s) on the Nth weekday that's a Monday or Tuesday
	// - dayOfMonth != nil && daysOfWeek != 0
	// at (T) every (N) month(s) on the Nth day, that's also a Monday or Tuesday?
	return current
}

func getNthWeekdayOfMonth(date time.Time, month time.Month, weekday time.Weekday, nth int) time.Time {
	firstOfMonth := time.Date(date.Year(), month, 1, 0, 0, 0, 0, date.Location())
	found := getNthWeekday(firstOfMonth, nth, weekday)
	return time.Date(found.Year(), found.Month(), found.Day(), date.Hour(), date.Minute(), date.Second(), 0, date.Location())
}
