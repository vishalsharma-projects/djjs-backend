package services

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/followCode/djjs-event-reporting-backend/app/models"
	"github.com/jung-kurt/gofpdf"
)

// GenerateEventPDF generates a PDF document for event details
func GenerateEventPDF(event *models.EventDetails, specialGuests []models.SpecialGuest, 
	volunteers []models.Volunteer, mediaList []models.EventMedia, 
	promotionMaterials []models.PromotionMaterialDetails, donations []models.Donation) ([]byte, error) {
	
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetAutoPageBreak(true, 25)
	pdf.SetMargins(10, 15, 10)
	pdf.AddPage()

	// Title
	pdf.SetFont("Arial", "B", 18)
	pdf.Cell(0, 12, "Event Details Report")
	pdf.Ln(8)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 6, fmt.Sprintf("Event ID: %d", event.ID))
	pdf.Ln(10)

	// Event Information Section
	pdf.SetFont("Arial", "B", 14)
	pdf.SetFillColor(240, 240, 240)
	pdf.CellFormat(0, 8, "Event Information", "", 1, "L", true, 0, "")
	pdf.SetFillColor(255, 255, 255)
	pdf.Ln(2)
	pdf.SetFont("Arial", "", 10)

	addField(pdf, "Event Type", event.EventType.Name, 45, 6)
	addField(pdf, "Event Category", event.EventCategory.Name, 45, 6)
	addField(pdf, "Scale", event.Scale, 45, 6)
	if event.Theme != "" {
		addField(pdf, "Theme", event.Theme, 45, 6)
	}
	addField(pdf, "Start Date", event.StartDate.Format("2006-01-02"), 45, 6)
	addField(pdf, "End Date", event.EndDate.Format("2006-01-02"), 45, 6)
	if event.DailyStartTime != nil {
		addField(pdf, "Daily Start Time", event.DailyStartTime.Format("15:04"), 45, 6)
	}
	if event.DailyEndTime != nil {
		addField(pdf, "Daily End Time", event.DailyEndTime.Format("15:04"), 45, 6)
	}
	if event.SpiritualOrator != "" {
		addField(pdf, "Spiritual Orator", event.SpiritualOrator, 45, 6)
	}
	pdf.Ln(3)

	// Venue Information
	pdf.SetFont("Arial", "B", 14)
	pdf.SetFillColor(240, 240, 240)
	pdf.CellFormat(0, 8, "Venue Information", "", 1, "L", true, 0, "")
	pdf.SetFillColor(255, 255, 255)
	pdf.Ln(2)
	pdf.SetFont("Arial", "", 10)
	if event.Country != "" {
		addField(pdf, "Country", event.Country, 45, 6)
	}
	if event.State != "" {
		addField(pdf, "State", event.State, 45, 6)
	}
	if event.City != "" {
		addField(pdf, "City", event.City, 45, 6)
	}
	if event.District != "" {
		addField(pdf, "District", event.District, 45, 6)
	}
	if event.PostOffice != "" {
		addField(pdf, "Post Office", event.PostOffice, 45, 6)
	}
	if event.Pincode != "" {
		addField(pdf, "Pincode", event.Pincode, 45, 6)
	}
	if event.Address != "" {
		addField(pdf, "Address", event.Address, 45, 6)
	}
	pdf.Ln(3)

	// Beneficiaries and Initiation in two columns
	pdf.SetFont("Arial", "B", 14)
	pdf.SetFillColor(240, 240, 240)
	pdf.CellFormat(95, 8, "Beneficiaries", "1", 0, "L", true, 0, "")
	pdf.CellFormat(95, 8, "Initiation", "1", 1, "L", true, 0, "")
	pdf.SetFillColor(255, 255, 255)
	pdf.SetFont("Arial", "", 9)
	pdf.Ln(2)

	// Beneficiaries column
	xStart := pdf.GetX()
	yStart := pdf.GetY()
	pdf.CellFormat(95, 6, fmt.Sprintf("Men: %d", event.BeneficiaryMen), "", 0, "L", false, 0, "")
	pdf.Ln(5)
	pdf.CellFormat(95, 6, fmt.Sprintf("Women: %d", event.BeneficiaryWomen), "", 0, "L", false, 0, "")
	pdf.Ln(5)
	pdf.CellFormat(95, 6, fmt.Sprintf("Children: %d", event.BeneficiaryChild), "", 0, "L", false, 0, "")
	pdf.Ln(5)
	pdf.SetFont("Arial", "B", 9)
	pdf.CellFormat(95, 6, fmt.Sprintf("Total: %d", event.BeneficiaryMen+event.BeneficiaryWomen+event.BeneficiaryChild), "", 0, "L", false, 0, "")
	pdf.Ln(6)

	// Initiation column
	pdf.SetXY(xStart+95, yStart)
	pdf.SetFont("Arial", "", 9)
	pdf.CellFormat(95, 6, fmt.Sprintf("Men: %d", event.InitiationMen), "", 0, "L", false, 0, "")
	pdf.Ln(5)
	pdf.SetXY(xStart+95, pdf.GetY())
	pdf.CellFormat(95, 6, fmt.Sprintf("Women: %d", event.InitiationWomen), "", 0, "L", false, 0, "")
	pdf.Ln(5)
	pdf.SetXY(xStart+95, pdf.GetY())
	pdf.CellFormat(95, 6, fmt.Sprintf("Children: %d", event.InitiationChild), "", 0, "L", false, 0, "")
	pdf.Ln(5)
	pdf.SetXY(xStart+95, pdf.GetY())
	pdf.SetFont("Arial", "B", 9)
	pdf.CellFormat(95, 6, fmt.Sprintf("Total: %d", event.InitiationMen+event.InitiationWomen+event.InitiationChild), "", 0, "L", false, 0, "")
	pdf.Ln(8)

	// Special Guests Table
	if len(specialGuests) > 0 {
		addTableSection(pdf, "Special Guests", len(specialGuests))
		headers := []string{"Name", "Designation", "Organization", "Email", "Contact", "City", "State"}
		colWidths := []float64{40, 35, 35, 40, 30, 25, 25}
		
		// Draw table header
		pdf.SetFont("Arial", "B", 8)
		pdf.SetFillColor(220, 220, 220)
		for i, header := range headers {
			pdf.CellFormat(colWidths[i], 7, header, "1", 0, "L", true, 0, "")
		}
		pdf.Ln(-1)
		
		// Draw table rows
		pdf.SetFont("Arial", "", 7)
		pdf.SetFillColor(255, 255, 255)
		for _, guest := range specialGuests {
			if pdf.GetY() > 270 {
				pdf.AddPage()
			}
			name := strings.TrimSpace(fmt.Sprintf("%s %s %s", guest.FirstName, guest.MiddleName, guest.LastName))
			if name == "" {
				name = "N/A"
			}
			rows := [][]string{
				{name, guest.Designation, guest.Organization, guest.Email, guest.PersonalNumber, guest.City, guest.State},
			}
			for _, row := range rows {
				for i, cell := range row {
					if cell == "" {
						cell = "N/A"
					}
					// Truncate long text
					if len(cell) > 20 {
						cell = cell[:17] + "..."
					}
					pdf.CellFormat(colWidths[i], 6, cell, "1", 0, "L", false, 0, "")
				}
				pdf.Ln(-1)
			}
		}
		pdf.Ln(5)
	}

	// Volunteers Table
	if len(volunteers) > 0 {
		addTableSection(pdf, "Volunteers", len(volunteers))
		headers := []string{"Name", "Branch", "Days", "Seva"}
		colWidths := []float64{60, 60, 25, 55}
		
		// Draw table header
		pdf.SetFont("Arial", "B", 8)
		pdf.SetFillColor(220, 220, 220)
		for i, header := range headers {
			pdf.CellFormat(colWidths[i], 7, header, "1", 0, "L", true, 0, "")
		}
		pdf.Ln(-1)
		
		// Draw table rows
		pdf.SetFont("Arial", "", 7)
		pdf.SetFillColor(255, 255, 255)
		for _, volunteer := range volunteers {
			if pdf.GetY() > 270 {
				pdf.AddPage()
			}
			branchName := "N/A"
			if volunteer.Branch.ID > 0 {
				branchName = volunteer.Branch.Name
				if len(branchName) > 25 {
					branchName = branchName[:22] + "..."
				}
			}
			seva := volunteer.SevaInvolved
			if len(seva) > 20 {
				seva = seva[:17] + "..."
			}
			rows := [][]string{
				{volunteer.VolunteerName, branchName, strconv.Itoa(volunteer.NumberOfDays), seva},
			}
			for _, row := range rows {
				for i, cell := range row {
					if cell == "" {
						cell = "N/A"
					}
					pdf.CellFormat(colWidths[i], 6, cell, "1", 0, "L", false, 0, "")
				}
				pdf.Ln(-1)
			}
		}
		pdf.Ln(5)
	}

	// Donations Table
	if len(donations) > 0 {
		addTableSection(pdf, "Donations", len(donations))
		headers := []string{"Type", "Details", "Amount (Rs.)"}
		colWidths := []float64{50, 80, 60}
		
		// Draw table header
		pdf.SetFont("Arial", "B", 8)
		pdf.SetFillColor(220, 220, 220)
		for i, header := range headers {
			pdf.CellFormat(colWidths[i], 7, header, "1", 0, "L", true, 0, "")
		}
		pdf.Ln(-1)
		
		// Draw table rows
		pdf.SetFont("Arial", "", 7)
		pdf.SetFillColor(255, 255, 255)
		totalAmount := 0.0
		for _, donation := range donations {
			if pdf.GetY() > 270 {
				pdf.AddPage()
			}
			details := donation.KindType
			if details == "" {
				details = "-"
			}
			if len(details) > 30 {
				details = details[:27] + "..."
			}
			amountStr := fmt.Sprintf("%.2f", donation.Amount)
			// Format donation type for display
			formattedType := formatDonationTypeForDisplay(donation.DonationType)
			rows := [][]string{
				{formattedType, details, amountStr},
			}
			for _, row := range rows {
				for i, cell := range row {
					if cell == "" {
						cell = "N/A"
					}
					align := "L"
					if i == 2 {
						align = "R"
					}
					pdf.CellFormat(colWidths[i], 6, cell, "1", 0, align, false, 0, "")
				}
				pdf.Ln(-1)
			}
			totalAmount += donation.Amount
		}
		// Total row
		if pdf.GetY() > 270 {
			pdf.AddPage()
		}
		pdf.SetFont("Arial", "B", 8)
		pdf.SetFillColor(240, 240, 240)
		pdf.CellFormat(colWidths[0]+colWidths[1], 7, "Total", "1", 0, "R", true, 0, "")
		pdf.CellFormat(colWidths[2], 7, fmt.Sprintf("%.2f", totalAmount), "1", 1, "R", true, 0, "")
		pdf.Ln(5)
	}

	// Promotion Materials Table
	if len(promotionMaterials) > 0 {
		addTableSection(pdf, "Promotion Materials", len(promotionMaterials))
		headers := []string{"Material Type", "Quantity", "Size"}
		colWidths := []float64{100, 40, 50}
		
		// Draw table header
		pdf.SetFont("Arial", "B", 8)
		pdf.SetFillColor(220, 220, 220)
		for i, header := range headers {
			pdf.CellFormat(colWidths[i], 7, header, "1", 0, "L", true, 0, "")
		}
		pdf.Ln(-1)
		
		// Draw table rows
		pdf.SetFont("Arial", "", 7)
		pdf.SetFillColor(255, 255, 255)
		for _, material := range promotionMaterials {
			if pdf.GetY() > 270 {
				pdf.AddPage()
			}
			materialType := "N/A"
			if material.PromotionMaterial.ID > 0 {
				materialType = material.PromotionMaterial.MaterialType
			}
			if len(materialType) > 35 {
				materialType = materialType[:32] + "..."
			}
			rows := [][]string{
				{materialType, strconv.Itoa(material.Quantity), material.Size},
			}
			for _, row := range rows {
				for i, cell := range row {
					if cell == "" {
						cell = "N/A"
					}
					pdf.CellFormat(colWidths[i], 6, cell, "1", 0, "L", false, 0, "")
				}
				pdf.Ln(-1)
			}
		}
		pdf.Ln(5)
	}

	// Media Coverage
	if len(mediaList) > 0 {
		addTableSection(pdf, "Media Coverage", len(mediaList))
		for i, media := range mediaList {
			if pdf.GetY() > 250 {
				pdf.AddPage()
			}
			pdf.SetFont("Arial", "B", 9)
			pdf.Cell(0, 6, fmt.Sprintf("Media %d", i+1))
			pdf.Ln(4)
			pdf.SetFont("Arial", "", 8)
			coverageType := "N/A"
			if media.MediaCoverageType.ID > 0 {
				coverageType = media.MediaCoverageType.MediaType
			}
			addFieldCompact(pdf, "Coverage Type", coverageType, 40, 5)
			if media.CompanyName != "" {
				addFieldCompact(pdf, "Company", media.CompanyName, 40, 5)
			}
			if media.CompanyEmail != "" {
				addFieldCompact(pdf, "Company Email", media.CompanyEmail, 40, 5)
			}
			if media.CompanyWebsite != "" {
				addFieldCompact(pdf, "Website", media.CompanyWebsite, 40, 5)
			}
			personName := strings.TrimSpace(fmt.Sprintf("%s %s %s %s", media.Prefix, media.FirstName, media.MiddleName, media.LastName))
			if personName != "" {
				addFieldCompact(pdf, "Contact Person", personName, 40, 5)
			}
			if media.Designation != "" {
				addFieldCompact(pdf, "Designation", media.Designation, 40, 5)
			}
			if media.Contact != "" {
				addFieldCompact(pdf, "Contact", media.Contact, 40, 5)
			}
			if media.Email != "" {
				addFieldCompact(pdf, "Email", media.Email, 40, 5)
			}
			pdf.Ln(3)
		}
	}

	// Footer
	pdf.SetY(-15)
	pdf.SetFont("Arial", "I", 7)
	pdf.CellFormat(0, 8, fmt.Sprintf("Generated on %s", time.Now().Format("2006-01-02 15:04:05")), "", 0, "C", false, 0, "")

	// Generate PDF bytes using buffer
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Helper function to add a field label and value
func addField(pdf *gofpdf.Fpdf, label, value string, labelWidth, lineHeight float64) {
	if value == "" {
		return
	}
	pdf.SetFont("Arial", "B", 9)
	pdf.Cell(labelWidth, lineHeight, label+":")
	pdf.SetFont("Arial", "", 9)
	// Wrap long values
	if len(value) > 60 {
		pdf.MultiCell(0, lineHeight, value, "", "L", false)
	} else {
		pdf.Cell(0, lineHeight, value)
		pdf.Ln(lineHeight)
	}
}

// Helper function for compact fields in media section
// formatDonationTypeForDisplay formats donation type with proper capitalization
func formatDonationTypeForDisplay(donationType string) string {
	if donationType == "" {
		return "N/A"
	}
	
	lowerType := strings.ToLower(strings.TrimSpace(donationType))
	switch lowerType {
	case "cash":
		return "Cash-Bank-Online"
	case "in-kind", "inkind":
		return "In-Kind"
	default:
		// Capitalize first letter of each word separated by hyphens
		parts := strings.Split(donationType, "-")
		formattedParts := make([]string, len(parts))
		for i, part := range parts {
			if len(part) > 0 {
				formattedParts[i] = strings.ToUpper(part[:1]) + strings.ToLower(part[1:])
			} else {
				formattedParts[i] = part
			}
		}
		return strings.Join(formattedParts, "-")
	}
}

func addFieldCompact(pdf *gofpdf.Fpdf, label, value string, labelWidth, lineHeight float64) {
	if value == "" {
		return
	}
	pdf.SetFont("Arial", "B", 7)
	pdf.Cell(labelWidth, lineHeight, label+":")
	pdf.SetFont("Arial", "", 7)
	if len(value) > 50 {
		value = value[:47] + "..."
	}
	pdf.Cell(0, lineHeight, value)
	pdf.Ln(lineHeight)
}

// Helper function to add table section header
func addTableSection(pdf *gofpdf.Fpdf, title string, count int) {
	if pdf.GetY() > 260 {
		pdf.AddPage()
	}
	pdf.SetFont("Arial", "B", 14)
	pdf.SetFillColor(240, 240, 240)
	pdf.CellFormat(0, 8, fmt.Sprintf("%s (%d)", title, count), "", 1, "L", true, 0, "")
	pdf.SetFillColor(255, 255, 255)
	pdf.Ln(2)
}
