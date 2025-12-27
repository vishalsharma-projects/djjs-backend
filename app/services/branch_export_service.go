package services

import (
	"bytes"
	"fmt"
	"time"

	"github.com/followCode/djjs-event-reporting-backend/app/models"
	"github.com/xuri/excelize/v2"
)

// ExportBranchesToExcel exports branches to an Excel file and returns the buffer
func ExportBranchesToExcel(branches []models.Branch) (*bytes.Buffer, error) {
	f := excelize.NewFile()

	sheetName := "Branches"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, fmt.Errorf("failed to create sheet: %v", err)
	}

	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	// Define headers
	headers := []string{
		"ID",
		"Name",
		"Email",
		"Coordinator Name",
		"Contact Number",
		"Established On",
		"Aashram Area",
		"Country",
		"State",
		"City",
		"District",
		"Address",
		"Pincode",
		"Post Office",
		"Police Station",
		"Open Days",
		"Daily Start Time",
		"Daily End Time",
		"Branch Code",
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

	// Helper functions
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
		return t.Format("2006-01-02")
	}

	formatTimePtr := func(t *time.Time) string {
		if t == nil || t.IsZero() {
			return "-"
		}
		return t.Format("2006-01-02")
	}

	// Write branch data
	for rowIndex, branch := range branches {
		row := rowIndex + 2

		values := []interface{}{
			branch.ID,
			getString(branch.Name),
			getString(branch.Email),
			getString(branch.CoordinatorName),
			getString(branch.ContactNumber),
			formatTimePtr(branch.EstablishedOn),
			branch.AashramArea,
			func() interface{} {
				if branch.Country.ID > 0 {
					return branch.Country.Name
				}
				return "-"
			}(),
			func() interface{} {
				if branch.State.ID > 0 {
					return branch.State.Name
				}
				return "-"
			}(),
			func() interface{} {
				if branch.City.ID > 0 {
					return branch.City.Name
				}
				return "-"
			}(),
			func() interface{} {
				if branch.District.ID > 0 {
					return branch.District.Name
				}
				return "-"
			}(),
			getString(branch.Address),
			getString(branch.Pincode),
			getString(branch.PostOffice),
			getString(branch.PoliceStation),
			getString(branch.OpenDays),
			getString(branch.DailyStartTime),
			getString(branch.DailyEndTime),
			getString(branch.BranchCode),
			func() string {
				if branch.Status {
					return "Active"
				}
				return "Inactive"
			}(),
			formatTime(branch.CreatedOn),
			formatTimePtr(branch.UpdatedOn),
			getString(branch.CreatedBy),
			getString(branch.UpdatedBy),
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
		f.SetColWidth(sheetName, colName, colName, 15)
	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, fmt.Errorf("failed to write Excel to buffer: %v", err)
	}

	if err := f.Close(); err != nil {
		return nil, fmt.Errorf("failed to close Excel file: %v", err)
	}

	return &buf, nil
}

// ExportMembersToExcel exports branch members to an Excel file and returns the buffer
func ExportMembersToExcel(members []models.BranchMember) (*bytes.Buffer, error) {
	f := excelize.NewFile()

	sheetName := "Members"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, fmt.Errorf("failed to create sheet: %v", err)
	}

	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	// Define headers
	headers := []string{
		"ID",
		"Name",
		"Member Type",
		"Branch Role",
		"Responsibility",
		"Age",
		"Qualification",
		"Date of Birth",
		"Date of Samarpan",
		"Branch Name",
		"Branch ID",
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

	// Helper functions
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
		return t.Format("2006-01-02")
	}

	formatTimePtr := func(t *time.Time) string {
		if t == nil || t.IsZero() {
			return "-"
		}
		return t.Format("2006-01-02")
	}

	// Write member data
	for rowIndex, member := range members {
		row := rowIndex + 2

		values := []interface{}{
			member.ID,
			getString(member.Name),
			getString(member.MemberType),
			getString(member.BranchRole),
			getString(member.Responsibility),
			func() interface{} {
				if member.Age > 0 {
					return member.Age
				}
				return "-"
			}(),
			getString(member.Qualification),
			formatTimePtr(member.DateOfBirth),
			formatTimePtr(member.DateOfSamarpan),
			func() interface{} {
				if member.Branch.ID > 0 {
					return member.Branch.Name
				}
				return "-"
			}(),
			member.BranchID,
			formatTime(member.CreatedOn),
			formatTimePtr(member.UpdatedOn),
			getString(member.CreatedBy),
			getString(member.UpdatedBy),
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
		f.SetColWidth(sheetName, colName, colName, 15)
	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, fmt.Errorf("failed to write Excel to buffer: %v", err)
	}

	if err := f.Close(); err != nil {
		return nil, fmt.Errorf("failed to close Excel file: %v", err)
	}

	return &buf, nil
}

