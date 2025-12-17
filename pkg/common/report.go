package common

import "fmt"

// ReportFilename builds the canonical filename for a deanonymization report.
func ReportFilename(appID ApplicationIdType, reportID RequestIdType) string {
	return fmt.Sprintf("%s_%s", appID.String(), reportID.String())
}
