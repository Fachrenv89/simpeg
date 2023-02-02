package api

import (
	"database/sql"
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/Rhymond/go-money"
	db "github.com/SemmiDev/simpeg/db/mysql"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ListLemburResponse struct {
	TotalJamLembur  string `json:"total_jam_lembur"`
	TotalGajiLembur string `json:"total_gaji_lembur"`

	Lembur            db.LemburResponse         `json:"lembur"`
	ListTanggalLembur []db.ListTanggalLemburRow `json:"list_tanggal_lembur"`
}

type tanggalLemburRequest struct {
	Tanggal  string  `json:"tanggal"`
	TotalJam float32 `json:"total_jam"`
}

type createLemburRequest struct {
	NoPayroll     string                 `json:"no_payroll"`
	Basis         float32                `json:"basis"`
	TanggalLembur []tanggalLemburRequest `json:"tanggal_lembur"` // optional
}

func (server *Server) createLembur(ctx *gin.Context) {
	var req createLemburRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	nonStaff, err := server.store.GetNonStaffByNoPayroll(ctx, req.NoPayroll)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	argLembur := db.CreateLemburParams{
		ID:        uuid.New().String(),
		NoPayroll: nonStaff.NoPayroll,
		Nama:      nonStaff.Nama,
		Jabatan:   nonStaff.Jabatan,
		GajiPokok: nonStaff.GajiPokok,
		Basis:     req.Basis,
	}

	err = server.store.CreateLembur(ctx, argLembur)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	for _, tanggalLembur := range req.TanggalLembur {
		arg := db.CreateTanggalLemburParams{
			ID:       uuid.New().String(),
			IDLembur: argLembur.ID,
			Tanggal:  tanggalLembur.Tanggal,
			TotalJam: tanggalLembur.TotalJam,
		}

		err := server.store.CreateTanggalLembur(ctx, arg)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Lembur berhasil dibuat",
	})
}

type createTanggalLemburRequest struct {
	Tanggal  string  `json:"tanggal"`
	TotalJam float32 `json:"total_jam"`
}

func (server *Server) createTanggalLembur(ctx *gin.Context) {
	// id of lembur table
	id := ctx.Query("id")

	// parse multiple tanggal lembur inserted by user
	var req []createTanggalLemburRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// each of requests payload, insert to tanggal_lembur table
	for _, v := range req {
		arg := db.CreateTanggalLemburParams{
			ID:       uuid.New().String(),
			IDLembur: id,
			Tanggal:  v.Tanggal,
			TotalJam: v.TotalJam,
		}

		err := server.store.CreateTanggalLembur(ctx, arg)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Tanggal lembur berhasil dibuat",
	})
}

type getLemburDetailsRequest struct {
	ID string `uri:"id"`
}

func (server *Server) lemburDetails(ctx *gin.Context) {
	// get filter by date by request
	filterByDate := ctx.Query("date")

	var req getLemburDetailsRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// first we will get the lembur details by lembur id
	lemburDetails, err := server.store.GetLemburByID(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// get all tanggal lembur by id lembur
	tanggalLemburDetails, err := server.store.GetTanggalLemburByIdLembur(ctx, lemburDetails.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// lemburAggregate is the the response payload which will server to client
	lemburAggregate := ListLemburResponse{Lembur: db.LemburResponse{
		ID:        lemburDetails.ID,
		NoPayroll: lemburDetails.NoPayroll,
		Nama:      lemburDetails.Nama,
		Jabatan:   lemburDetails.Jabatan,
	}}

	filtered := []db.GetTanggalLemburByIdLemburRow{}

	// tanggal lembur must filtered first
	// to make sure will append by filter
	for _, v := range tanggalLemburDetails {

		// if filter exists
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

	// so we will replace value from tanggalLemburDetails
	// to  filtered data
	tanggalLemburDetails = filtered

	// and then we will process the gaji lembur
	lemburAggregateResult := calculateGajiLembur(
		lemburDetails,
		tanggalLemburDetails,
		lemburAggregate,
	)

	ctx.JSON(http.StatusOK, lemburAggregateResult)
}

func (server *Server) deleteLembur(ctx *gin.Context) {
	id := ctx.Param("id")
	err := server.store.DeleteLembur(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Lembur berhasil dihapus",
	})
}

func (server *Server) listLembur(ctx *gin.Context) {
	filterByDate := ctx.Query("date")

	listLembur, err := server.store.ListLembur(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	listTanggalLembur, err := server.store.ListTanggalLembur(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	lemburAggregate := make([]ListLemburResponse, 0, len(listLembur))
	for _, lembur := range listLembur {
		temp := ListLemburResponse{Lembur: db.LemburResponse{
			ID:        lembur.ID,
			NoPayroll: lembur.NoPayroll,
			Nama:      lembur.Nama,
			Jabatan:   lembur.Jabatan,
		}}

		for _, tanggalLembur := range listTanggalLembur {
			if tanggalLembur.IDLembur == lembur.ID {
				temp.ListTanggalLembur = append(temp.ListTanggalLembur, tanggalLembur)
			}
		}

		tanggalLemburDetails := make([]db.GetTanggalLemburByIdLemburRow, 0, len(temp.ListTanggalLembur))

		// filter by date
		filteredLembur := []db.ListTanggalLemburRow{}
		for _, v := range temp.ListTanggalLembur {
			var getKeyword string

			if filterByDate == "" {
				year, month, _ := time.Now().Date()
				getKeyword = fmt.Sprintf("%d-%d", year, month)
			} else {
				splitDate := strings.Split(v.Tanggal, "-")
				getYear := splitDate[0]
				getMonth := splitDate[1]
				getKeyword = getYear + "-" + getMonth
			}

			if getKeyword == filterByDate {
				filteredLembur = append(filteredLembur, v)
			}
		}

		for _, v := range filteredLembur {
			tanggalLemburDetails = append(tanggalLemburDetails, db.GetTanggalLemburByIdLemburRow{
				ID:       v.ID,
				Tanggal:  v.Tanggal,
				TotalJam: v.TotalJam,
			})
		}

		if len(tanggalLemburDetails) == 0 {
			continue
		}

		temp.ListTanggalLembur = filteredLembur
		result := calculateGajiLembur(lembur, tanggalLemburDetails, temp)
		lemburAggregate = append(lemburAggregate, result)
	}

	ctx.JSON(http.StatusOK, lemburAggregate)
}

func calculateGajiLembur(
	lemburDetails db.Lembur,
	tanggalLemburDetails []db.GetTanggalLemburByIdLemburRow,
	lemburAggregate ListLemburResponse) ListLemburResponse {

	var (
		totalGajiLembur decimal.Decimal
		totalJamLembur  float32
	)

	// each tanggalLembur we iterate
	for _, tanggalLembur := range tanggalLemburDetails {
		var hasil = decimal.Zero
		var basis = decimal.NewFromFloat32(lemburDetails.Basis)

		tanggalLemburWithTotalJam := tanggalLembur.TotalJam
		if tanggalLemburWithTotalJam != 0 {
			if tanggalLemburWithTotalJam == 1 {
				hasil = basis.Mul(decimal.NewFromFloat32(1.5))
			} else {
				if tanggalLemburWithTotalJam > 2 {
					add := basis.Mul(decimal.NewFromFloat32(1.5))
					left := decimal.NewFromFloat32(tanggalLemburWithTotalJam - 1)
					right := basis.Mul(decimal.NewFromFloat32(2))
					hasil = left.Mul(right).Add(add)
				} else {
					add := basis.Mul(decimal.NewFromFloat32(1.5))
					hasil = decimal.NewFromFloat32(tanggalLemburWithTotalJam).Mul(basis).Add(add)
				}
			}

			// the the result we add to tao totalGajiLembur
			totalGajiLembur = totalGajiLembur.Add(hasil)

			// and add also to totalJamLembur
			totalJamLembur += tanggalLemburWithTotalJam
		}
	}

	// we convert to float64
	totalGajiLemburInFloat64, _ := totalGajiLembur.Float64()
	// and totalGajiLemburInFloat64 round up
	totalGajiLemburInFloat64 = math.Ceil(totalGajiLemburInFloat64)

	// display for Rp format
	lemburAggregate.TotalGajiLembur = money.NewFromFloat(totalGajiLemburInFloat64, money.IDR).Display()
	lemburAggregate.TotalJamLembur = fmt.Sprintf("%1.f", totalJamLembur)

	// display for Rp format
	lemburAggregate.Lembur.GajiPokok = money.NewFromFloat(float64(lemburDetails.GajiPokok), money.IDR).Display()
	lemburAggregate.Lembur.Basis = money.NewFromFloat(float64(lemburDetails.Basis), money.IDR).Display()

	return lemburAggregate
}

func calculateGajiLembur2(
	lemburDetails db.Lembur,
	tanggalLemburDetails []db.GetTanggalLemburByIdLemburRow,
	lemburAggregate ListLemburResponse) (ListLemburResponse, float64) {

	var (
		totalGajiLembur decimal.Decimal
		totalJamLembur  float32
	)

	// each tanggalLembur we iterate
	for _, tanggalLembur := range tanggalLemburDetails {
		var hasil = decimal.Zero
		var basis = decimal.NewFromFloat32(lemburDetails.Basis)

		tanggalLemburWithTotalJam := tanggalLembur.TotalJam
		if tanggalLemburWithTotalJam != 0 {
			if tanggalLemburWithTotalJam == 1 {
				hasil = basis.Mul(decimal.NewFromFloat32(1.5))
			} else {
				if tanggalLemburWithTotalJam > 2 {
					add := basis.Mul(decimal.NewFromFloat32(1.5))
					left := decimal.NewFromFloat32(tanggalLemburWithTotalJam - 1)
					right := basis.Mul(decimal.NewFromFloat32(2))
					hasil = left.Mul(right).Add(add)
				} else {
					add := basis.Mul(decimal.NewFromFloat32(1.5))
					hasil = decimal.NewFromFloat32(tanggalLemburWithTotalJam).Mul(basis).Add(add)
				}
			}

			// the the result we add to tao totalGajiLembur
			totalGajiLembur = totalGajiLembur.Add(hasil)

			// and add also to totalJamLembur
			totalJamLembur += tanggalLemburWithTotalJam
		}
	}

	// we convert to float64
	totalGajiLemburInFloat64, _ := totalGajiLembur.Float64()
	// and totalGajiLemburInFloat64 round up
	totalGajiLemburInFloat64 = math.Ceil(totalGajiLemburInFloat64)

	// display for Rp format
	lemburAggregate.TotalGajiLembur = money.NewFromFloat(totalGajiLemburInFloat64, money.IDR).Display()
	lemburAggregate.TotalJamLembur = fmt.Sprintf("%1.f", totalJamLembur)

	// display for Rp format
	lemburAggregate.Lembur.GajiPokok = money.NewFromFloat(float64(lemburDetails.GajiPokok), money.IDR).Display()
	lemburAggregate.Lembur.Basis = money.NewFromFloat(float64(lemburDetails.Basis), money.IDR).Display()

	return lemburAggregate, totalGajiLemburInFloat64
}
