package api

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"net/http"
)

func GenerateExcelReport(w http.ResponseWriter, fileName string, data []Rekapitulasi) error {
	xlsx := excelize.NewFile()
	defer func() {
		if err := xlsx.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	sheet1Name := "Non Staff"
	index, err := xlsx.NewSheet(sheet1Name)
	if err != nil {
		return err
	}

	err = xlsx.AutoFilter(sheet1Name, "A1:AA1", nil)
	if err != nil {
		return err
	}

	// add width for each column
	xlsx.SetColWidth(sheet1Name, "A", "AA", 23)

	xlsx.SetCellValue(sheet1Name, "A1", "Kode")
	xlsx.SetCellValue(sheet1Name, "B1", "Nama Karyawan")
	xlsx.SetCellValue(sheet1Name, "C1", "Jabatan")
	xlsx.SetCellValue(sheet1Name, "D1", "Golongan")
	xlsx.SetCellValue(sheet1Name, "E1", "Kode TNG")
	xlsx.SetCellValue(sheet1Name, "F1", "Gaji Pokok")
	xlsx.SetCellValue(sheet1Name, "G1", "Tunjangan Kemahalan")
	xlsx.SetCellValue(sheet1Name, "H1", "Tunjangan Perumahan")
	xlsx.SetCellValue(sheet1Name, "I1", "Tunjangan Jabatan")
	xlsx.SetCellValue(sheet1Name, "J1", "Tunjangan Lain Pph 21")
	xlsx.SetCellValue(sheet1Name, "K1", "TKJkk")
	xlsx.SetCellValue(sheet1Name, "L1", "TKJkm")
	xlsx.SetCellValue(sheet1Name, "M1", "Kesehatan P")
	xlsx.SetCellValue(sheet1Name, "N1", "Lembur")
	xlsx.SetCellValue(sheet1Name, "O1", "Gaji Bruto")
	xlsx.SetCellValue(sheet1Name, "P1", "PPH21 Per Bulan")
	xlsx.SetCellValue(sheet1Name, "Q1", "Pinjaman Tunai")
	xlsx.SetCellValue(sheet1Name, "R1", "Pinjaman Koperasi")
	xlsx.SetCellValue(sheet1Name, "S1", "Pinjaman LainLain")
	xlsx.SetCellValue(sheet1Name, "T1", "BPJS TKJHT")
	xlsx.SetCellValue(sheet1Name, "U1", "BPJS TKJPN")
	xlsx.SetCellValue(sheet1Name, "V1", "BPJS KesehatanK")
	xlsx.SetCellValue(sheet1Name, "W1", "BPJS KesehatanP")
	xlsx.SetCellValue(sheet1Name, "X1", "BPJS TKJKK")
	xlsx.SetCellValue(sheet1Name, "Y1", "BPJS TKJKM")
	xlsx.SetCellValue(sheet1Name, "Z1", "Jumlah Potongan")
	xlsx.SetCellValue(sheet1Name, "AA1", "Gaji Netto")

	for i, row := range data {
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("A%d", i+2), row.Kode)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("B%d", i+2), row.NamaKaryawan)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("C%d", i+2), row.Jabatan)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("D%d", i+2), row.Golongan)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("E%d", i+2), row.KodeTNG)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("F%d", i+2), row.GajiPokok)

		if row.TunjanganKemahalan != nil {
			xlsx.SetCellValue(sheet1Name, fmt.Sprintf("G%d", i+2), *row.TunjanganKemahalan)
		} else {
			xlsx.SetCellValue(sheet1Name, fmt.Sprintf("G%d", i+2), "-")
		}

		if row.TunjanganPerumahan != nil {
			xlsx.SetCellValue(sheet1Name, fmt.Sprintf("H%d", i+2), *row.TunjanganPerumahan)
		} else {
			xlsx.SetCellValue(sheet1Name, fmt.Sprintf("H%d", i+2), "-")
		}

		if row.TunjanganJabatan != nil {
			xlsx.SetCellValue(sheet1Name, fmt.Sprintf("I%d", i+2), *row.TunjanganJabatan)
		} else {
			xlsx.SetCellValue(sheet1Name, fmt.Sprintf("I%d", i+2), "-")
		}

		if row.TunjanganLainPph21 != nil {
			xlsx.SetCellValue(sheet1Name, fmt.Sprintf("J%d", i+2), *row.TunjanganLainPph21)
		} else {
			xlsx.SetCellValue(sheet1Name, fmt.Sprintf("J%d", i+2), "-")
		}

		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("K%d", i+2), row.TKJkk)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("L%d", i+2), row.TKJkm)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("M%d", i+2), row.KesehatanP)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("N%d", i+2), row.Lembur)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("O%d", i+2), row.GajiBruto)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("P%d", i+2), row.PPH21PerBulan)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("Q%d", i+2), row.PinjamanTunai)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("R%d", i+2), row.PinjamanKoperasi)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("S%d", i+2), row.PinjamanLainLain)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("T%d", i+2), row.BPJSTKJHT)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("U%d", i+2), row.BPJSTKJPN)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("V%d", i+2), row.BPJSKesehatanK)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("W%d", i+2), row.BPJSKesehatanP)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("X%d", i+2), row.BPJSTKJKK)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("Y%d", i+2), row.BPJSTKJKM)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("Z%d", i+2), row.JumlahPotongan)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("AA%d", i+2), row.GajiNetto)

		fmt.Printf("Non Staff with %s Payrol DONE \n", row.Kode)
	}

	xlsx.SetActiveSheet(index)

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName+".xlsx")
	w.Header().Set("Content-Transfer-Encoding", "binary")

	if err := xlsx.Write(w); err != nil {
		return err
	}

	return nil
}
