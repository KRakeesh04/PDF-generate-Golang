package main

import (
	// "fmt"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"strconv"

	"github.com/signintech/gopdf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type District struct {
	ID           uint          `json:"id" gorm:"primaryKey"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	Name         string        `json:"name"`
	Coordinators []Coordinator `json:"coordinators"`
	ExamCentres  []ExamCentre  `json:"exam_centres"`
}
type Coordinator struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Name        string    `json:"name"`
	TelephoneNo string    `json:"telephone_no"`
	DistrictID  uint      `json:"district_id"`
}
type Stream struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`

	Subject1   *Subject `json:"subject1" gorm:"foreignKey:Subject1ID;references:ID"`
	Subject1ID *uint    `json:"subject1_id" gorm:"index"`
	Subject2   *Subject `json:"subject2" gorm:"foreignKey:Subject2ID;references:ID"`
	Subject2ID *uint    `json:"subject2_id" gorm:"index"`
	Subject3   *Subject `json:"subject3" gorm:"foreignKey:Subject3ID;references:ID"`
	Subject3ID *uint    `json:"subject3_id" gorm:"index"`
}
type Subject struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Code string `json:"code" gorm:"index,unique"`
	Name string `json:"name"`
}
type ExamCentre struct {
	ID               uint      `json:"id" gorm:"primaryKey"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	Name             string    `json:"name"`
	DistrictID       uint      `json:"district_id"`
	District         District  `json:"district" gorm:"foreignKey:DistrictID"`
	Place            string    `json:"place"`
	Gender           string    `json:"gender" gorm:"type:enum('Male', 'Female', 'Mixed'); not null" validate:"oneof=Male Female Mixed"`
	BusRoute         *string   `json:"bus_route"`
	BusDepartureTime *string   `json:"bus_departure_time"`
	BusArrivalTime   *string   `json:"bus_arrival_time"`
	TravelDuration   *string   `json:"travel_duration"`
	SubsituteTimes   *string   `json:"substitude_times"`
	Counts           []Count   `json:"counts" gorm:"-"`
}
type Count struct {
	SubjectName  string `json:"subject"`
	SubjectCode  string `json:"code"`
	Medium       string `json:"medium"`
	StudentCount int    `json:"count"`
}

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Role     string `json:"role"`
	Username string `json:"username" gorm:"unique"`
	Password string `json:"-"`
	Approved bool   `json:"approved"`
}
type Student struct {
	IndexNo        int        `json:"index_no" gorm:"primaryKey;unique"` // Primary key and unique
	Name           string     `json:"name" gorm:"not null"`              // Name cannot be null
	StreamID       uint       `json:"stream_id" gorm:"not null"`         // Foreign key reference
	Stream         Stream     `json:"stream" gorm:"foreignKey:StreamID"` // Reference to Stream
	Medium         string     `json:"medium" gorm:"type:enum('Tamil', 'English');not null" validate:"oneof=Tamil English"`
	RankDistrictID uint       `json:"rank_district_id"`                               // Foreign key reference
	RankDistrict   District   `json:"rank_district" gorm:"foreignKey:RankDistrictID"` // Reference to District
	ExamDistrictID uint       `json:"exam_district_id"`                               // Foreign key reference
	ExamDistrict   District   `json:"exam_district" gorm:"foreignKey:ExamDistrictID"` // Reference to District
	ExamCentreID   uint       `json:"exam_centre_id"`                                 // Foreign key reference
	ExamCentre     ExamCentre `json:"exam_centre" gorm:"foreignKey:ExamCentreID"`     // Reference to ExamCentre
	NIC            string     `json:"nic" gorm:"not null"`                            // NIC cannot be null
	Gender         string     `json:"gender" gorm:"type:enum('Male', 'Female');not null" validate:"oneof=Male Female"`
	School         *string    `json:"school"`
	Address        *string    `json:"address"`
	Email          *string    `json:"email"`
	TelephoneNo    *string    `json:"telephone_no"`
	RegisteredByID uint       `json:"registered_by_id" gorm:"not null"`               // Foreign key reference
	RegisteredBy   User       `json:"registered_by" gorm:"foreignKey:RegisteredByID"` // Reference to User
	CheckedByID    *uint      `json:"checked_by_id"`                                  // Nullable foreign key reference
	CheckedBy      *User      `json:"checked_by" gorm:"foreignKey:CheckedByID"`       // Nullable reference to User
	CheckedAt      *time.Time `json:"checked_at"`                                     // Nullable timestamp
}

type StudentAdmission struct {
	IndexNo int    `json:"index_no"`
	Name    string `json:"name"`
	Stream  string `json:"stream"`
	NIC     string `json:"nic"`
}

func getStudentsForAdmission(db *gorm.DB) ([]StudentAdmission, error) {
	var students []*Student
	// Fetch all students from the database

	// SELECT * FROM students ORDER BY LENGTH(name) DESC;
	result := db.Model(&Student{}).Find(&students)
	// sort students by name length in descending order in golang
	sort.Slice(students, func(i, j int) bool {
		return len(students[i].Name) > len(students[j].Name)
	})

	if result.Error != nil {
		return nil, result.Error
	}

	// Helper function to map stream_id to stream name
	getStreamName := func(streamID int) string {
		switch streamID {
		case 1:
			return "ICT"
		case 2:
			return "MATHS"
		case 3:
			return "OTHER"
		case 4:
			return "BIO"
		default:
			return "UNKNOWN"
		}
	}

	var studentsAdmission []StudentAdmission

	for _, student := range students {
		studentAdmission := StudentAdmission{
			IndexNo: student.IndexNo,
			Name:    student.Name,
			Stream:  getStreamName(int(student.StreamID)), // Use actual stream
			NIC:     student.NIC,
		}
		studentsAdmission = append(studentsAdmission, studentAdmission)
	}

	return studentsAdmission, nil
}

func ConnectDB() (*gorm.DB, error) {
	dsn := "moraexams:moraexams-testing@tcp(127.0.0.1:3306)/moraexams?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	// Initialize the database connection
	db, err := ConnectDB()
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	students, err := getStudentsForAdmission(db)
	if err != nil {
		log.Fatal("Error fetching students:", err)
	}

	if len(students) == 0 {
		log.Println("No students found for admission.")
		return
	}

	log.Printf("Found %d students for admission.\n", len(students))

	for i := 0; i < len(students); i += 25 {
		end := i + 25
		if end > len(students) {
			end = len(students)
		}
		GenerateAdmissionSheetForAllStudents(students[i:end], i/25+1)
	}

}

func GenerateAdmissionSheetForAllStudents(students []StudentAdmission, num int) {

	district := "JAFFNA"
	center := "Jaffna Central College"
	subject := "ICT"
	sub_number := "10"
	part := "II"
	medium := "EM"
	date := "2025-07-21"

	// Init PDF
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4}) // A4
	pdf.AddPage()
	pdf.SetLineWidth(0.5)

	// Add header image
	SetImgPositionAndStyle(pdf, 25, 15, "admission_src/Header.png", 545, 120.5)
	// Add footer image
	SetImgPositionAndStyle(pdf, 25, 823, "admission_src/Footer_black.png", 545, 15)
	// Add background img
	SetImgPositionAndStyle(pdf, 120, 370, "admission_src/logo_light.png", 350, 195)


	// Add normal font
	err := pdf.AddTTFFont("segoeuithis", "admission_src/segoe-ui-this/segoeuithis.ttf")
	if err != nil {
		log.Fatal("AddTTFFont normal failed:", err)
	}

	// Add bold font
	err = pdf.AddTTFFont("segoeuithis_bold", "admission_src/segoe-ui-this/segoeuithibd.ttf")
	if err != nil {
		log.Fatal("AddTTFFont bold failed:", err)
	}

	
	fontName := "segoeuithis"
	align := gopdf.CellOption{Align: gopdf.Center | gopdf.Middle}
	SetTextPositionAndAlign(pdf, 160, 103, align, sub_number, fontName, 12)
	SetTextPositionAndAlign(pdf, 200, 103, align, medium, fontName, 12)
	SetTextPositionAndAlign(pdf, 235, 103, align, part, fontName, 12)
	align = gopdf.CellOption{Align: gopdf.Left | gopdf.Middle}
	SetTextPositionAndAlign(pdf, 410, 107, align, district, fontName, 11)
	SetTextPositionAndAlign(pdf, 410, 82, align, subject, fontName, 11)
	SetTextPositionAndAlign(pdf, 140, 107, align, date, fontName, 11)
	
	SetTextPositionAndAlign(pdf, 40, 145, align,  "Center : ", fontName, 12)
	SetTextPositionAndAlign(pdf, 88, 145, align,  center, fontName, 11)


	// set bold font for the headings
	err = pdf.SetFont("segoeuithis_bold", "", 11.5)
	if err != nil {
		log.Fatal(err)
	}

	// Table settings
	startX := 25.0  
	startY := 175.0 
	rowHeight := 24.5
	numRows := 25 + 1 

	// Column widths
	colWidths := []float64{20, 70, 260, 80, 50, 65}
	headers := []string{"No", "IndexNo", "Name", "Signature", "Marks", "Checked"}

	x := startX
	y := startY

	pdf.SetFillColor(240, 240, 240) // light gray
	for _, w := range colWidths {
		pdf.RectFromUpperLeftWithStyle(x, y, w, rowHeight, "F") // fill rectangle
		x += w
	}

	pdf.Line(startX, y, startX+sum(colWidths), y)

	// Reset fill & text color for text
	pdf.SetFillColor(255, 255, 255) // no fill behind text
	pdf.SetTextColor(0, 0, 0)       // black text

	x = startX
	for i, h := range headers {
		pdf.SetX(x)
		pdf.SetY(y)
		pdf.CellWithOption(&gopdf.Rect{W: colWidths[i], H: rowHeight}, h, gopdf.CellOption{Align: gopdf.Center | gopdf.Middle})
		x += colWidths[i]
	}
	y += rowHeight

	err = pdf.SetFont("segoeuithis", "", 10)
	if err != nil {
		log.Fatal(err)
	}

	// Draw rows
	for i, s := range students {
		x = startX

		// Draw the top border for a row
		pdf.Line(startX, y, startX+sum(colWidths), y)

		cols := []string{
			formatNumber(i + 1),
			strconv.Itoa(s.IndexNo) + "000",
			FormatNameInitials(s.Name),
			"",
			"",
			"",
		}

		for j, text := range cols {
			// Draw the left border for a cell
			pdf.Line(x, startY, x, startY+float64(numRows)*rowHeight)

			pdf.SetX(x)
			pdf.SetY(y)
			if j != 2 {
				pdf.CellWithOption(&gopdf.Rect{W: colWidths[j], H: rowHeight}, text, gopdf.CellOption{Align: gopdf.Center | gopdf.Middle})

			} else {

				pdf.CellWithOption(&gopdf.Rect{W: colWidths[j], H: rowHeight}, "  "+text, gopdf.CellOption{Align: gopdf.Left | gopdf.Middle})
			}

			x += colWidths[j]
		}
		// Draw the right border for a cell
		pdf.Line(x, startY, x, startY+float64(numRows)*rowHeight)

		y += rowHeight
	}

	// Draw the bottom border for the last row
	pdf.Line(startX, y, startX+sum(colWidths), y)

	// Save pdf
	err = pdf.WritePdf(fmt.Sprintf("generated/admission_sheet_no_%d.pdf", num))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Generated: admission_sheet_no_%d.pdf\n", num)
}

func formatNumber(n int) string {
	if n < 10 {
		return "0" + fmt.Sprint('0'+n)
	}
	return fmt.Sprint('0'+n/10) + fmt.Sprint('0'+n%10)
}

func sum(arr []float64) float64 {
	total := 0.0
	for _, v := range arr {
		total += v
	}
	return total
}

func FormatNameInitials(name string) string {
	if len(name) == 0 {
		return ""
	}

	parts := strings.Fields(name)
	if len(parts) <= 4 {
		if len(name) > 34 {
			formattedName := ""
			for i, part := range parts {
				if i < 1 {
					formattedName += string(part[0]) + ". "
				} else {
					formattedName += part + " "
				}
			}
			return formattedName
		}
		return name
	}
	initials := ""
	for i, part := range parts {
		runes := []rune(part)
		if len(runes) == 0 {
			continue
		}
		if i > 3 {
			initials += string(runes) + " "
		} else {
			initials += string(runes[0]) + ". "
		}

	}
	return strings.TrimSpace(initials)
}

func SetTextPositionAndAlign(pdf *gopdf.GoPdf, x, y float64, align gopdf.CellOption, text string, fontName string, size int) {
	err := pdf.SetFont(fontName, "", size)
	if err != nil {
		log.Fatal(err)
	}
	pdf.SetX(x)
	pdf.SetY(y)
	// Use a fixed width and height for the cell
	pdf.CellWithOption(&gopdf.Rect{W: 200, H: 20}, text, align)
	pdf.SetTextColor(0, 0, 0)       // Reset text color to black
	pdf.SetFillColor(255, 255, 255) // Reset fill color
}

func SetImgPositionAndStyle(pdf *gopdf.GoPdf, x, y float64, imgPath string, w, h float64) {
	rect := gopdf.Rect{
		W: w, // width in points
		H: h,  // height in points
	}
	err := pdf.Image(imgPath, x, y, &rect)
	if err != nil {
		log.Fatal(err)
	}
}

