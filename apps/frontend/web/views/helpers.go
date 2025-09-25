package views

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/a-h/templ"
)

type TripWizardStep string

const (
	TripWizardStepBasics    TripWizardStep = "basics"
	TripWizardStepLogistics TripWizardStep = "logistics"
	TripWizardStepReview    TripWizardStep = "review"
)

var wizardOrder = map[TripWizardStep]int{
	TripWizardStepBasics:    1,
	TripWizardStepLogistics: 2,
	TripWizardStepReview:    3,
}

type TripsWizardData struct {
	Step   TripWizardStep
	TripID string
	OrgID  string
	Form   map[string]string
}

type TripSummary struct {
	ID          string
	Name        string
	Location    string
	Status      string
	StatusClass string
	DateRange   string
	IsCurrent   bool
}

func makeWizardProgress(step TripWizardStep) templ.Component {
	return TripWizardProgress(step)
}

func normalizeStep(step TripWizardStep) TripWizardStep {
	if step == "" {
		return TripWizardStepBasics
	}
	return step
}

func isStepActive(current, target TripWizardStep) bool {
	return normalizeStep(current) == target
}

func isStepCompleted(current, target TripWizardStep) bool {
	return wizardOrder[normalizeStep(current)] > wizardOrder[target]
}

func wizardTextClass(current, target TripWizardStep) string {
	if isStepActive(current, target) {
		return "text-neutral-100"
	}
	if isStepCompleted(current, target) {
		return "text-neutral-200"
	}
	return "text-neutral-400"
}

func tripStatusBadge(status string) string {
	switch status {
	case "listed":
		return "text-lg leading-6 text-green-400"
	case "complete":
		return "text-lg leading-6 text-neutral-100"
	case "draft":
		return "text-lg leading-6 text-yellow-300"
	default:
		return "text-lg leading-6 text-neutral-400"
	}
}

func toInt64(v any) int64 {
	switch val := v.(type) {
	case int:
		return int64(val)
	case int64:
		return val
	case float64:
		return int64(val)
	case json.Number:
		n, _ := val.Int64()
		return n
	case string:
		if val == "" {
			return 0
		}
		if i, err := strconv.ParseInt(val, 10, 64); err == nil {
			return i
		}
		return 0
	default:
		return 0
	}
}

func floatString(v any) string {
	switch val := v.(type) {
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(val), 'f', -1, 64)
	case int:
		return strconv.Itoa(val)
	case int64:
		return strconv.FormatInt(val, 10)
	case json.Number:
		f, err := val.Float64()
		if err == nil {
			return strconv.FormatFloat(f, 'f', -1, 64)
		}
		return val.String()
	case string:
		return val
	default:
		return fmt.Sprint(val)
	}
}

func intString(v any) string {
	switch val := v.(type) {
	case int:
		return strconv.Itoa(val)
	case int64:
		return strconv.FormatInt(val, 10)
	case float64:
		return strconv.FormatInt(int64(val), 10)
	case json.Number:
		i, err := val.Int64()
		if err == nil {
			return strconv.FormatInt(i, 10)
		}
		return val.String()
	case string:
		return val
	default:
		return fmt.Sprint(val)
	}
}

func formatDateInput(ts int64) string {
	if ts == 0 {
		return ""
	}
	return time.Unix(ts, 0).UTC().Format("2006-01-02")
}

func FormatTripDateRange(start, end int64) string {
	if start == 0 && end == 0 {
		return "Dates TBD"
	}
	layout := "Jan 2"
	startStr := time.Unix(start, 0).UTC().Format(layout)
	if end == 0 {
		return startStr
	}
	endStr := time.Unix(end, 0).UTC().Format(layout)
	if start == end {
		return startStr
	}
	return fmt.Sprintf("%s – %s", startStr, endStr)
}

func FormatLocation(city, country string) string {
	switch {
	case city != "" && country != "":
		return fmt.Sprintf("%s, %s", city, country)
	case city != "":
		return city
	default:
		return country
	}
}

func formValue(data TripsWizardData, key string) string {
	if data.Form == nil {
		return ""
	}
	return data.Form[key]
}

func TripFormFromPayload(data map[string]any) map[string]string {
	form := make(map[string]string)
	copyString := func(key string) {
		if val, ok := data[key]; ok && val != nil {
			str := fmt.Sprint(val)
			if str != "" {
				form[key] = str
			}
		}
	}
	copyString("name")
	copyString("privacy_type")
	copyString("housing_type")
	copyString("trip_type")
	copyString("city")
	copyString("country")
	copyString("currency")
	copyString("description")
	copyString("mission")
	copyString("status")

	if val, ok := data["volunteer_limit"]; ok {
		form["volunteer_limit"] = intString(val)
	}
	if val, ok := data["price"]; ok {
		form["price"] = floatString(val)
	}
	if val, ok := data["latitude"]; ok {
		form["latitude"] = floatString(val)
	}
	if val, ok := data["longitude"]; ok {
		form["longitude"] = floatString(val)
	}

	if ts := toInt64(data["start_date"]); ts != 0 {
		form["start_date"] = formatDateInput(ts)
	}
	if ts := toInt64(data["end_date"]); ts != 0 {
		form["end_date"] = formatDateInput(ts)
	}

	return form
}

func NewTripSummaryFromPayload(data map[string]any) TripSummary {
	form := TripFormFromPayload(data)
	status := fmt.Sprint(data["status"])
	if status == "" {
		status = "draft"
	}
	start := toInt64(data["start_date"])
	end := toInt64(data["end_date"])
	summary := TripSummary{
		ID:          fmt.Sprint(data["id"]),
		Name:        form["name"],
		Location:    FormatLocation(form["city"], form["country"]),
		Status:      status,
		StatusClass: tripStatusBadge(status),
		DateRange:   FormatTripDateRange(start, end),
		IsCurrent:   false,
	}
	if summary.Name == "" {
		summary.Name = "Untitled Trip"
	}
	return summary
}
