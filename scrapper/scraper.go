package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"database/sql"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/mattn/go-sqlite3" // Untuk SQLite
)

func main() {
	// Buka file HTML lokal
	f, err := os.Open("Elements (Little Alchemy 2) _ Little Alchemy Wiki _ Fandom.html")
	if err != nil {
		log.Fatalf("Gagal membuka file HTML: %v", err)
	}
	defer f.Close()

	// Buat dokument dari file reader
	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		log.Fatalf("Gagal membuat dokumen dari file: %v", err)
	}

	// Buat file CSV untuk gabungan tabel
	file, err := os.Create("alchemy.csv")
	if err != nil {
		log.Printf("Gagal membuat file alchemy.csv: %v", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Tulis header CSV
	writer.Write([]string{"Element", "Item1", "Item2"})

	// Setup SQLite database
	db, err := sql.Open("sqlite3", "./alchemy.db")
	if err != nil {
		log.Fatalf("Gagal membuka SQLite database: %v", err)
	}
	defer db.Close()

	// Buat tabel dalam database jika belum ada
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS elements (
			element TEXT,
			item1 TEXT,
			item2 TEXT
		);
	`)
	if err != nil {
		log.Fatalf("Gagal membuat tabel: %v", err)
	}

	// Menyimpan data ke database dan CSV
	doc.Find("h3").Each(func(i int, s *goquery.Selection) {
		headline := s.Find(".mw-headline")
		sectionID, exists := headline.Attr("id")
		if !exists {
			return
		}

		table := s.NextAllFiltered("table").First()
		if table.Length() == 0 {
			return
		}

		// Memasukkan data ke dalam SQLite dan CSV
		table.Find("tr").Each(func(i int, tr *goquery.Selection) {
			if i == 0 {
				return // skip header
			}
			tds := tr.Find("td")
			if tds.Length() < 2 {
				return
			}

			element := strings.TrimSpace(tds.Eq(0).Text())

			liFound := false
			tds.Eq(1).Find("li").Each(func(_ int, li *goquery.Selection) {
				text := strings.TrimSpace(li.Text())
				if strings.Contains(text, "+") {
					// Pisahkan berdasarkan '+'
					parts := strings.Split(text, "+")
					if len(parts) != 2 {
						return
					}
					item1 := strings.TrimSpace(parts[0])
					item2 := strings.TrimSpace(parts[1])

					// Menulis ke file CSV
					writer.Write([]string{element, item1, item2})

					// Menulis ke SQLite
					insertDataToSQLite(db, &element, &item1, &item2)

					liFound = true
				}
			})

			if !liFound {
				// Tidak ada '+' dalam <li>, tulis satu baris NULL ke CSV dan SQLite
				writer.Write([]string{element, "", ""})
				insertDataToSQLite(db, &element, nil, nil)
			}
		})

		fmt.Println("Saved:", sectionID)
	})

}

// Fungsi untuk memasukkan data ke SQLite
func insertDataToSQLite(db *sql.DB, element, item1, item2 *string) {

	// Menyisipkan data ke dalam database SQLite
	_, err := db.Exec("INSERT INTO elements (element, item1, item2) VALUES (?, ?, ?)", element, item1, item2)
	if err != nil {
		log.Printf("Gagal memasukkan data ke SQLite: %v", err)
	}
}
