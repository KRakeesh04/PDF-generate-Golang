package main

import (
	// "fmt"
	"log"
	"github.com/signintech/gopdf"
)



type Student struct {
	Stream   string
	IndexNo  string
	Name     string
}

var students = []Student{
	{"Bio", "1001", "Alice Fernando"},
	{"Maths", "1002", "Bandara Silva"},
	{"Commerce", "1003", "Chamara Perera"},
	{"Bio", "1004", "Dilani Kumari"},
	{"Maths", "1005", "Eshan Jayawardena"},
	{"Commerce", "1006", "Farhan Ismail"},
	{"Bio", "1007", "Gihan Wickramasinghe"},
	{"Maths", "1008", "Harsha Rajapaksha"},
	{"Commerce", "1009", "Ishara Ranasinghe"},
	{"Bio", "1010", "Janani Gunasekara"},
	{"Maths", "1011", "Kavindu Hettiarachchi"},
	{"Commerce", "1012", "Lahiru Rathnayake"},
	{"Bio", "1013", "Mihiri Weerasinghe"},
	{"Maths", "1014", "Nadeesha Madushani"},
	{"Commerce", "1015", "Oshadi Dilrukshi"},
	{"Bio", "1016", "Pasindu Madushan"},
	{"Maths", "1017", "Ravindu Fernando"},
	{"Commerce", "1018", "Sajini Samarasinghe"},
	{"Bio", "1019", "Thilina Karunaratne"},
	{"Maths", "1020", "Udara Wickrama"},
	{"Commerce", "1021", "Vindya Abeywardana"},
	{"Bio", "1022", "Wasana Perera"},
	{"Maths", "1023", "Yasiru Senanayake"},
	{"Commerce", "1024", "Zuhair Ameer"},
	{"Bio", "1025", "Ashan Bandara"},
}


func main() {
	
	// for i, student := range students {
	// 	fmt.Printf("%02d. Stream: %s | IndexNo: %s | Name: %s\n", i+1, student.Stream, student.IndexNo, student.Name)
	// }

	// Init PDF
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4}) // A4
	pdf.AddPage()
	pdf.SetLineWidth(0.5)
	

	// Add normal font
	err := pdf.AddTTFFont("dejavu", "/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf")
	if err != nil {
		log.Fatal("AddTTFFont normal failed:", err)
	}

	// Add bold font
	err = pdf.AddTTFFont("dejavu_bold", "/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf")
	if err != nil {
		log.Fatal("AddTTFFont bold failed:", err)
	}


	// set bold font for the headings
	err = pdf.SetFont("dejavu_bold", "", 11.5)
	if err != nil {
		log.Fatal(err)
	}

	// Table settings
	startX := 35.0		// starting position on the page where the table begins.
	startY := 135.0		//
	rowHeight := 24.5
	numRows := len(students) + 1 // +1 for header row

	// Column widths
	colWidths := []float64{20, 75, 80, 160, 80, 50, 65}


	// Draw header
	headers := []string{"No", "Stream", "IndexNo", "Name", "Signature", "Marks", "Checked"}

	x := startX
	y := startY

	pdf.SetFillColor(220, 220, 220) // light gray
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

	err = pdf.SetFont("dejavu", "", 12)
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
			s.Stream,
			s.IndexNo,
			s.Name,
			"",
			"",
			"",
		}

		for j, text := range cols {
			// Draw the left border for a cell
			pdf.Line(x, startY, x, startY+float64(numRows)*rowHeight)

			pdf.SetX(x)
			pdf.SetY(y)
			if j == 2 || j == 0{
				pdf.CellWithOption(&gopdf.Rect{W: colWidths[j], H: rowHeight}, text, gopdf.CellOption{Align: gopdf.Center | gopdf.Middle})
				
			} else {

				pdf.CellWithOption(&gopdf.Rect{W: colWidths[j], H: rowHeight}, " "+text, gopdf.CellOption{Align: gopdf.Left | gopdf.Middle})
			}
			

			x += colWidths[j]
		}

		// Draw the right border for a cell
		pdf.Line(x, startY, x, startY+float64(numRows)*rowHeight)

		y += rowHeight
	}

	// Draw the bottom border for the last row
	pdf.Line(startX, y, startX+sum(colWidths), y)



	// Save
	err = pdf.WritePdf("admission_sheet.pdf")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Generated: admission_sheet.pdf")
}



func formatNumber(n int) string {
	if n < 10 {
		return "0" + string('0'+n)
	}
	return string('0'+n/10) + string('0'+n%10)
}


func sum(arr []float64) float64 {
    total := 0.0
    for _, v := range arr {
        total += v
    }
    return total
}
