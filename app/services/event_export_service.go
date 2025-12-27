package services

import (
	"bytes"
	"fmt"
	"time"

	"github.com/followCode/djjs-event-reporting-backend/app/models"
	"github.com/xuri/excelize/v2"
)

// ExportEventsToExcel exports events to an Excel file and returns the buffer
func ExportEventsToExcel(events []models.EventDetails) (*bytes.Buffer, error) {
	f := excelize.NewFile()

	sheetName := "Events"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, fmt.Errorf("failed to create sheet: %v", err)
	}

	// Set active sheet
	f.SetActiveSheet(index)

	// Delete default Sheet1
	f.DeleteSheet("Sheet1")

	// Define headers
	headers := []string{
		"ID",
		"Event Type",
		"Event Category",
		"Scale",
		"Theme",
		"Start Date",
		"End Date",
		"Daily Start Time",
		"Daily End Time",
		"Spiritual Orator",
		"Language",
		"Country",
		"State",
		"City",
		"District",
		"Post Office",
		"Pincode",
		"Address",
		"Beneficiaries - Men",
		"Beneficiaries - Women",
		"Beneficiaries - Children",
		"Initiation - Men",
		"Initiation - Women",
		"Initiation - Children",
		"Branch",
		"Status",
		"Created On",
		"Updated On",
		"Created By",
		"Updated By",
	}

	// Write headers
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
	}

	// Style header row
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:   true,
			Size:   12,
			Color:  "#FFFFFF",
			Family: "Arial",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#4472C4"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err == nil {
		f.SetCellStyle(sheetName, "A1", fmt.Sprintf("%c1", 'A'+len(headers)-1), headerStyle)
	}

	// Write event data
	for rowIndex, event := range events {
		row := rowIndex + 2 // Start from row 2 (row 1 is headers)

		// Helper function to safely get string values
		getString := func(s string) interface{} {
			if s == "" {
				return "-"
			}
			return s
		}

		// Helper function to format time
		formatTime := func(t time.Time) string {
			if t.IsZero() {
				return "-"
			}
			return t.Format("2006-01-02")
		}

		// Helper function to format time only
		formatTimeOnly := func(t *models.TimeOnly) string {
			if t == nil || t.Time.IsZero() {
				return "-"
			}
			return t.Time.Format("15:04")
		}

		// Write data
		values := []interface{}{
			event.ID,
			getString(event.EventType.Name),
			getString(event.EventCategory.Name),
			getString(event.Scale),
			getString(event.Theme),
			formatTime(event.StartDate),
			formatTime(event.EndDate),
			formatTimeOnly(event.DailyStartTime),
			formatTimeOnly(event.DailyEndTime),
			getString(event.SpiritualOrator),
			getString(event.Language),
			getString(event.Country),
			getString(event.State),
			getString(event.City),
			getString(event.District),
			getString(event.PostOffice),
			getString(event.Pincode),
			getString(event.Address),
			event.BeneficiaryMen,
			event.BeneficiaryWomen,
			event.BeneficiaryChild,
			event.InitiationMen,
			event.InitiationWomen,
			event.InitiationChild,
			func() interface{} {
				if event.Branch != nil {
					return event.Branch.Name
				}
				return "-"
			}(),
			getString(event.Status),
			formatTime(event.CreatedOn),
			func() interface{} {
				if event.UpdatedOn != nil {
					return formatTime(*event.UpdatedOn)
				}
				return "-"
			}(),
			getString(event.CreatedBy),
			getString(event.UpdatedBy),
		}

		for colIndex, value := range values {
			cell := fmt.Sprintf("%c%d", 'A'+colIndex, row)
			f.SetCellValue(sheetName, cell, value)
		}
	}

	// Set column widths (using helper function to convert index to column name)
	getColumnName := func(index int) string {
		result := ""
		for index >= 0 {
			result = string(rune('A'+(index%26))) + result
			index = index/26 - 1
		}
		return result
	}

	// Auto-size columns
	for i := 0; i < len(headers); i++ {
		colName := getColumnName(i)
		f.SetColWidth(sheetName, colName, colName, 15)
	}

	// Write to buffer
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, fmt.Errorf("failed to write Excel to buffer: %v", err)
	}

	// Close the file
	if err := f.Close(); err != nil {
		return nil, fmt.Errorf("failed to close Excel file: %v", err)
	}

	return &buf, nil
}

// ExportVolunteersToExcel exports volunteers to an Excel file and returns the buffer
func ExportVolunteersToExcel(volunteers []models.Volunteer, eventName string) (*bytes.Buffer, error) {
	f := excelize.NewFile()

	sheetName := "Volunteers"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, fmt.Errorf("failed to create sheet: %v", err)
	}

	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	// Define headers
	headers := []string{
		"ID",
		"Branch",
		"Volunteer Name",
		"Contact",
		"Number of Days",
		"Seva Involved",
		"Mention Seva",
		"Event ID",
		"Created On",
		"Updated On",
		"Created By",
		"Updated By",
	}

	// Write headers
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
	}

	// Style header row
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:   true,
			Size:   12,
			Color:  "#FFFFFF",
			Family: "Arial",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#4472C4"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err == nil {
		getColumnName := func(index int) string {
			result := ""
			for index >= 0 {
				result = string(rune('A'+(index%26))) + result
				index = index/26 - 1
			}
			return result
		}
		lastCol := getColumnName(len(headers) - 1)
		f.SetCellStyle(sheetName, "A1", lastCol+"1", headerStyle)
	}

	// Write volunteer data
	for rowIndex, vol := range volunteers {
		row := rowIndex + 2

		getString := func(s string) interface{} {
			if s == "" {
				return "-"
			}
			return s
		}

		formatTime := func(t time.Time) string {
			if t.IsZero() {
				return "-"
			}
			return t.Format("2006-01-02 15:04:05")
		}

		values := []interface{}{
			vol.ID,
			func() interface{} {
				if vol.Branch.ID > 0 {
					return vol.Branch.Name
				}
				return "-"
			}(),
			getString(vol.VolunteerName),
			getString(vol.Contact),
			vol.NumberOfDays,
			getString(vol.SevaInvolved),
			getString(vol.MentionSeva),
			vol.EventID,
			formatTime(vol.CreatedOn),
			func() interface{} {
				if vol.UpdatedOn != nil {
					return formatTime(*vol.UpdatedOn)
				}
				return "-"
			}(),
			getString(vol.CreatedBy),
			getString(vol.UpdatedBy),
		}

		for colIndex, value := range values {
			cell := fmt.Sprintf("%c%d", 'A'+colIndex, row)
			f.SetCellValue(sheetName, cell, value)
		}
	}

	// Set column widths
	getColumnName := func(index int) string {
		result := ""
		for index >= 0 {
			result = string(rune('A'+(index%26))) + result
			index = index/26 - 1
		}
		return result
	}

	for i := 0; i < len(headers); i++ {
		colName := getColumnName(i)
		f.SetColWidth(sheetName, colName, colName, 18)
	}

	// Write to buffer
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, fmt.Errorf("failed to write Excel to buffer: %v", err)
	}

	if err := f.Close(); err != nil {
		return nil, fmt.Errorf("failed to close Excel file: %v", err)
	}

	return &buf, nil
}

// ExportSpecialGuestsToExcel exports special guests to an Excel file and returns the buffer
func ExportSpecialGuestsToExcel(guests []models.SpecialGuest, eventName string) (*bytes.Buffer, error) {
	f := excelize.NewFile()

	sheetName := "Special Guests"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, fmt.Errorf("failed to create sheet: %v", err)
	}

	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	// Define headers
	headers := []string{
		"ID",
		"Gender",
		"Prefix",
		"First Name",
		"Middle Name",
		"Last Name",
		"Full Name",
		"Designation",
		"Organization",
		"Email",
		"City",
		"State",
		"Personal Number",
		"Contact Person",
		"Contact Person Number",
		"Reference Branch ID",
		"Reference Volunteer ID",
		"Reference Person Name",
		"Event ID",
		"Created On",
		"Updated On",
		"Created By",
		"Updated By",
	}

	// Write headers
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
	}

	// Style header row
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:   true,
			Size:   12,
			Color:  "#FFFFFF",
			Family: "Arial",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#4472C4"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err == nil {
		getColumnName := func(index int) string {
			result := ""
			for index >= 0 {
				result = string(rune('A'+(index%26))) + result
				index = index/26 - 1
			}
			return result
		}
		lastCol := getColumnName(len(headers) - 1)
		f.SetCellStyle(sheetName, "A1", lastCol+"1", headerStyle)
	}

	// Write guest data
	for rowIndex, guest := range guests {
		row := rowIndex + 2

		getString := func(s string) interface{} {
			if s == "" {
				return "-"
			}
			return s
		}

		formatTime := func(t time.Time) string {
			if t.IsZero() {
				return "-"
			}
			return t.Format("2006-01-02 15:04:05")
		}

		fullName := fmt.Sprintf("%s %s %s %s",
			getString(guest.Prefix),
			getString(guest.FirstName),
			getString(guest.MiddleName),
			getString(guest.LastName))

		values := []interface{}{
			guest.ID,
			getString(guest.Gender),
			getString(guest.Prefix),
			getString(guest.FirstName),
			getString(guest.MiddleName),
			getString(guest.LastName),
			fullName,
			getString(guest.Designation),
			getString(guest.Organization),
			getString(guest.Email),
			getString(guest.City),
			getString(guest.State),
			getString(guest.PersonalNumber),
			getString(guest.ContactPerson),
			getString(guest.ContactPersonNumber),
			getString(guest.ReferenceBranchID),
			getString(guest.ReferenceVolunteerID),
			getString(guest.ReferencePersonName),
			guest.EventID,
			formatTime(guest.CreatedOn),
			func() interface{} {
				if guest.UpdatedOn != nil {
					return formatTime(*guest.UpdatedOn)
				}
				return "-"
			}(),
			getString(guest.CreatedBy),
			getString(guest.UpdatedBy),
		}

		for colIndex, value := range values {
			cell := fmt.Sprintf("%c%d", 'A'+colIndex, row)
			f.SetCellValue(sheetName, cell, value)
		}
	}

	// Set column widths
	getColumnName := func(index int) string {
		result := ""
		for index >= 0 {
			result = string(rune('A'+(index%26))) + result
			index = index/26 - 1
		}
		return result
	}

	for i := 0; i < len(headers); i++ {
		colName := getColumnName(i)
		f.SetColWidth(sheetName, colName, colName, 18)
	}

	// Write to buffer
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, fmt.Errorf("failed to write Excel to buffer: %v", err)
	}

	if err := f.Close(); err != nil {
		return nil, fmt.Errorf("failed to close Excel file: %v", err)
	}

	return &buf, nil
}

// ExportEventMediaToExcel exports event media to an Excel file and returns the buffer
func ExportEventMediaToExcel(media []models.EventMedia, eventName string) (*bytes.Buffer, error) {
	f := excelize.NewFile()

	sheetName := "Event Media"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, fmt.Errorf("failed to create sheet: %v", err)
	}

	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	// Define headers
	headers := []string{
		"ID",
		"Event ID",
		"Media Coverage Type",
		"Company Name",
		"Company Email",
		"Company Website",
		"Gender",
		"Prefix",
		"First Name",
		"Middle Name",
		"Last Name",
		"Full Name",
		"Designation",
		"Contact",
		"Email",
		"File Type",
		"Original Filename",
		"Created On",
		"Updated On",
		"Created By",
		"Updated By",
	}

	// Write headers
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
	}

	// Style header row
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:   true,
			Size:   12,
			Color:  "#FFFFFF",
			Family: "Arial",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#4472C4"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err == nil {
		getColumnName := func(index int) string {
			result := ""
			for index >= 0 {
				result = string(rune('A'+(index%26))) + result
				index = index/26 - 1
			}
			return result
		}
		lastCol := getColumnName(len(headers) - 1)
		f.SetCellStyle(sheetName, "A1", lastCol+"1", headerStyle)
	}

	// Write media data
	for rowIndex, m := range media {
		row := rowIndex + 2

		getString := func(s string) interface{} {
			if s == "" {
				return "-"
			}
			return s
		}

		formatTime := func(t time.Time) string {
			if t.IsZero() {
				return "-"
			}
			return t.Format("2006-01-02 15:04:05")
		}

		fullName := fmt.Sprintf("%s %s %s %s",
			getString(m.Prefix),
			getString(m.FirstName),
			getString(m.MiddleName),
			getString(m.LastName))

		values := []interface{}{
			m.ID,
			m.EventID,
			func() interface{} {
				if m.MediaCoverageType.ID > 0 {
					return m.MediaCoverageType.MediaType
				}
				return "-"
			}(),
			getString(m.CompanyName),
			getString(m.CompanyEmail),
			getString(m.CompanyWebsite),
			getString(m.Gender),
			getString(m.Prefix),
			getString(m.FirstName),
			getString(m.MiddleName),
			getString(m.LastName),
			fullName,
			getString(m.Designation),
			getString(m.Contact),
			getString(m.Email),
			getString(m.FileType),
			getString(m.OriginalFilename),
			formatTime(m.CreatedOn),
			formatTime(m.UpdatedOn),
			getString(m.CreatedBy),
			getString(m.UpdatedBy),
		}

		for colIndex, value := range values {
			cell := fmt.Sprintf("%c%d", 'A'+colIndex, row)
			f.SetCellValue(sheetName, cell, value)
		}
	}

	// Set column widths
	getColumnName := func(index int) string {
		result := ""
		for index >= 0 {
			result = string(rune('A'+(index%26))) + result
			index = index/26 - 1
		}
		return result
	}

	for i := 0; i < len(headers); i++ {
		colName := getColumnName(i)
		f.SetColWidth(sheetName, colName, colName, 18)
	}

	// Write to buffer
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, fmt.Errorf("failed to write Excel to buffer: %v", err)
	}

	if err := f.Close(); err != nil {
		return nil, fmt.Errorf("failed to close Excel file: %v", err)
	}

	return &buf, nil
}
