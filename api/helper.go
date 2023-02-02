package api

import (
	"context"
	"strings"

	db "github.com/SemmiDev/simpeg/db/mysql"
)

func (server *Server) lemburDetailsByID(noPayroll, date string) (ListLemburResponse, float64) {
	filterByDate := date

	lemburDetails, err := server.store.GetLemburByNoPayroll(context.Background(), noPayroll)
	if err != nil {
		return ListLemburResponse{}, 0
	}

	listTanggalLembur, err := server.store.GetTanggalLemburByIdLembur(context.Background(), lemburDetails.ID)
	if err != nil {
		return ListLemburResponse{}, 0
	}

	lemburAggregate := ListLemburResponse{Lembur: db.LemburResponse{
		ID:        lemburDetails.ID,
		NoPayroll: lemburDetails.NoPayroll,
		Nama:      lemburDetails.Nama,
		Jabatan:   lemburDetails.Jabatan,
	}}

	var filtered []db.GetTanggalLemburByIdLemburRow

	for _, v := range listTanggalLembur {
		if filterByDate != "" {
			splitDate := strings.Split(v.Tanggal, "-")
			getYear := splitDate[0]
			getMonth := splitDate[1]
			getKeyword := getYear + "-" + getMonth

			// if filter equals to tanggal from database
			if getKeyword == filterByDate {
				// so we will appen to lembur aggregate
				lemburAggregate.ListTanggalLembur = append(lemburAggregate.ListTanggalLembur, db.ListTanggalLemburRow{
					ID:       v.ID,
					IDLembur: lemburDetails.ID,
					Tanggal:  v.Tanggal,
					TotalJam: v.TotalJam,
				})

				// and dont forget to also append to filtered slice
				filtered = append(filtered, v)
			}
		} else {
			// if filter doesn't exists
			// we will add each tanggal lembur to lemburAggregate
			lemburAggregate.ListTanggalLembur = append(lemburAggregate.ListTanggalLembur, db.ListTanggalLemburRow{
				ID:       v.ID,
				IDLembur: lemburDetails.ID,
				Tanggal:  v.Tanggal,
				TotalJam: v.TotalJam,
			})

			// and also add to filtered
			filtered = append(filtered, v)
		}
	}

	// so we will replace value from listTanggalLembur
	// to  filtered data
	listTanggalLembur = filtered

	// and then we will process the gaji lembur
	lemburAggregateResult, gajiLemburInFloat64 := calculateGajiLembur2(
		lemburDetails,
		listTanggalLembur,
		lemburAggregate,
	)

	return lemburAggregateResult, gajiLemburInFloat64
}
