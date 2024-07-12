package commands

import (
	"fmt"
	"os"
	"testing"
	"time"
)


func TestExportToMDTable(t *testing.T) {
  //given
  engine := Engine{}
  engine.Setup("./testResources/sessions.json")
  tempDir, err := os.MkdirTemp("", "")
  if err != nil {
    t.Errorf("Testsetup failed")
  }

  testExportPath := fmt.Sprintf("%s/%s", tempDir, "testExport.md")

  //when
  err = engine.ExportToMDTable(testExportPath, buildDate("09-07-2024"))

  //then
  if err != nil {
    t.Errorf("Expected no err got %s", err)
  }
  actualContent, err := os.ReadFile(testExportPath)
  if err != nil {
    t.Errorf("%s", err)
  }
  expected := `

|No|Start|Dauer|Beschreibung|Anmerkungen|
|-----|-----|-----|------|------|
|1 |2024-07-09 06:44:20| 00:25| warmup| |
|2 |2024-07-09 07:17:13| 00:25| A| |
|3 |2024-07-09 08:12:07| 00:25| A| |
|4 |2024-07-09 08:36:42| 00:25| A| |
|5 |2024-07-09 09:13:27| 00:25| A| |
||Total||------|------|
<!-- TBLFM: @>$3=sum(@I..@-1);hm -->
`
  if string(actualContent) != expected {
    t.Errorf("<%s>\n to be equal\n<%s>", actualContent, expected)
  }
} 

func buildDate(dateStr string) time.Time {
  layout := "02-01-2006"
  date, _ := time.Parse(layout, dateStr)
  return date
}
