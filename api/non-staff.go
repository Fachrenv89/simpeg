package api

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/SemmiDev/simpeg/util"
	"github.com/jung-kurt/gofpdf"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/Rhymond/go-money"
	db "github.com/SemmiDev/simpeg/db/mysql"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Server) deleteNonStaff(ctx *gin.Context) {
	id := ctx.Param("id")
	err := s.store.DeleteNonStaff(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (s *Server) createNonStaff(ctx *gin.Context) {
	var request db.CreateNonStaffParams
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Get data from db
	data, err := s.store.GetNonStaffByNoPayroll(ctx, request.NoPayroll)
	if err != nil {
		if err != sql.ErrNoRows {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	request.ID = uuid.New().String()

	// Check if data is empty
	if data.ID != "" {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Insert data to db
	err = s.store.CreateNonStaff(ctx, request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, request)
}

type Rekapitulasi struct {
	Kode         string `json:"kode"`
	NamaKaryawan string `json:"nama_karyawan"`
	Jabatan      string `json:"jabatan"`
	Golongan     string `json:"golongan"`
	KodeTNG      string `json:"kode_tng"` // PTKP
	GajiPokok    string `json:"gaji_pokok"`

	TunjanganKemahalan *float32 `json:"tunjangan_kemahalan"`
	TunjanganPerumahan *float32 `json:"tunjangan_perumahan"`
	TunjanganJabatan   *float32 `json:"tunjangan_jabatan"`
	TunjanganLainPph21 *float32 `json:"tunjangan_lain_pph21"`

	TKJkk      string `json:"jkk"`
	TKJkm      string `json:"jkm"`
	KesehatanP string `json:"kesehatan_p"`
	Lembur     string `json:"lembur"`
	GajiBruto  string `json:"gaji_bruto"`

	PPH21PerBulan    string `json:"pph_21_perbulan"`
	PinjamanTunai    string `json:"pinjaman_tunai"`
	PinjamanKoperasi string `json:"pinjaman_koperasi"`
	PinjamanLainLain string `json:"pinjaman_lain_lain"`

	BPJSTKJHT      string `json:"bpjs_tk_jht"`
	BPJSTKJPN      string `json:"bpjs_tk_jpn"`
	BPJSKesehatanK string `json:"bpjs_kesehatan_k"`
	BPJSKesehatanP string `json:"bpjs_kesehatan_p"`
	BPJSTKJKK      string `json:"bpjs_tk_jkk"`
	BPJSTKJKM      string `json:"bpjs_tk_jkm"`
	JumlahPotongan string `json:"jumlah_potongan"`
	GajiNetto      string `json:"gaji_netto"`
}

type NonStaffResponse struct {
	NonStaff     db.NonStaff `json:"non_staff"`
	Rekapitulasi `json:"rekapitulasi"`

	MasaKerja                  int32  `json:"masa_kerja"`
	Lembur                     string `json:"lembur"`
	Gaji                       string `json:"gaji"`
	GajiLemburTunjangan        string `json:"gaji_lembur_tunjangan"`
	Jkk                        string `json:"jkk"`
	Jkm                        string `json:"jkm"`
	Bpjs_4persen               string `json:"bpjs_4persen"`
	Penghasilan_bruto_perbulan string `json:"penghasilan_bruto_perbulan"`
	Penghasilan_bruto_pertahun string `json:"penghasilan_bruto_pertahun"`
	Biaya_jabatan              string `json:"biaya_jabatan"`
	Jht                        string `json:"jht"`
	Jpn                        string `json:"jpn"`
	Bpjs_kesehatan             string `json:"bpjs_kesehatan"`
	Gaji_netto_pertahun        string `json:"gaji_netto_pertahun"`
	Ptkp                       string `json:"ptkp"`
	Pkp                        string `json:"pkp"`
	Pph_21_pertahun            string `json:"pph_21_pertahun"`
	Pph_21_pertahun_tanpa_npwp string `json:"pph_21_pertahun_tanpa_npwp"`
	Pph_21_perbulan            string `json:"pph_21_perbulan"`
}

type listStaffResponse struct {
	db.NonStaff
	GajiPokokStr string `json:"gaji_pokok_str"`
}

func (server *Server) listNonStaff(ctx *gin.Context) {
	data, err := server.store.ListNonStaff(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := make([]listStaffResponse, 0, len(data))
	for _, v := range data {
		response = append(response, listStaffResponse{
			NonStaff:     v,
			GajiPokokStr: money.NewFromFloat(float64(v.GajiPokok), money.IDR).Display(),
		})
	}

	ctx.JSON(http.StatusOK, response)
}

func (s *Server) countNonStaff(ctx *gin.Context) {
	totalNonStaff := s.store.CountNonStaff(ctx)
	ctx.JSON(http.StatusOK, gin.H{
		"total_non_staff": totalNonStaff,
	})
}

func (s *Server) editNonStaff(ctx *gin.Context) {
	id := ctx.Query("id")

	var request db.UpdateNonStaffParams
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	request.ID = id

	err := s.store.UpdateNonStaff(ctx, request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func (s *Server) detailNonStaff(ctx *gin.Context) {
	id := ctx.Query("id")
	noPayroll := ctx.Query("no_payroll")
	date := ctx.Query("date")

	dataNonStaff, err := s.store.GetNonStaffByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	_, totalGajiLembur := s.lemburDetailsByID(noPayroll, date)

	var response NonStaffResponse

	// masa kerja
	var masaKerja int32
	if dataNonStaff.AwalKerja == 1 && dataNonStaff.StatusKerja == "aktif" {
		masaKerja = 12
	} else {
		masaKerja = dataNonStaff.AkhirKerja - dataNonStaff.AwalKerja + 1
	}

	// masa kerja
	response.MasaKerja = masaKerja

	// data non staff
	response.NonStaff = dataNonStaff

	// rekap lembur
	// make sure lembur exists
	var gaji float32
	if totalGajiLembur != 0 {
		response.Lembur = money.NewFromFloat(totalGajiLembur, money.IDR).Display()

		gaji = float32(dataNonStaff.GajiPokok)
		if dataNonStaff.TunjanganKemahalan != nil {
			gaji += *dataNonStaff.TunjanganKemahalan
		}
		if dataNonStaff.TunjanganPerumahan != nil {
			gaji += *dataNonStaff.TunjanganPerumahan
		}
		if dataNonStaff.TunjanganJabatan != nil {
			gaji += *dataNonStaff.TunjanganJabatan
		}

		// gaji
		response.Gaji = money.NewFromFloat(float64(gaji), money.IDR).Display()
	} else {
		gaji = float32(dataNonStaff.GajiPokok)
		response.Gaji = money.NewFromFloat(float64(gaji), money.IDR).Display()
	}

	// GajiLemburTunjangan
	var gajiLemburTunjangan float32
	gajiLemburTunjangan += dataNonStaff.GajiPokok
	if dataNonStaff.TunjanganKemahalan != nil {
		gajiLemburTunjangan += *dataNonStaff.TunjanganKemahalan
	}
	if dataNonStaff.TunjanganPerumahan != nil {
		gajiLemburTunjangan += *dataNonStaff.TunjanganPerumahan
	}
	if dataNonStaff.TunjanganJabatan != nil {
		gajiLemburTunjangan += *dataNonStaff.TunjanganJabatan
	}
	if dataNonStaff.TunjanganLainPph21 != nil {
		gajiLemburTunjangan += *dataNonStaff.TunjanganLainPph21
	}

	gajiLemburTunjangan += float32(totalGajiLembur)
	response.GajiLemburTunjangan = money.NewFromFloat(float64(gajiLemburTunjangan), money.IDR).Display()

	jkk := gaji * 0.54 / 100
	response.Jkk = money.NewFromFloat(float64(jkk), money.IDR).Display()

	jkm := gaji * 0.3 / 100
	response.Jkm = money.NewFromFloat(float64(jkm), money.IDR).Display()

	var bpjs4persen float32
	if gaji > 8000000 {
		bpjs4persen = 8000000 * 0.04
	} else {
		bpjs4persen = gaji * 0.04
	}
	response.Bpjs_4persen = money.NewFromFloat(float64(bpjs4persen), money.IDR).Display()

	var gajiBrutoPerBulan float32
	gajiBrutoPerBulan += gajiLemburTunjangan
	gajiBrutoPerBulan += jkk
	gajiBrutoPerBulan += jkm
	gajiBrutoPerBulan += bpjs4persen

	response.Penghasilan_bruto_perbulan = money.NewFromFloat(float64(gajiBrutoPerBulan), money.IDR).Display()

	var gajiBrutoPerTahun float32 = (gajiBrutoPerBulan * 11) + gajiBrutoPerBulan
	response.Penghasilan_bruto_pertahun = money.NewFromFloat(float64(gajiBrutoPerTahun), money.IDR).Display()

	var biayaJabatan float32
	if gajiBrutoPerTahun*0.05 > float32(500_000*masaKerja) {
		biayaJabatan = float32(500_000 * masaKerja)
	} else {
		biayaJabatan = 0.05 * gajiBrutoPerTahun
	}
	response.Biaya_jabatan = money.NewFromFloat(float64(biayaJabatan), money.IDR).Display()

	var jht float32 = (0.02 * gaji * 11) + (0.02 * gaji)
	response.Jht = money.NewFromFloat(math.Ceil(float64(jht)), money.IDR).Display()

	var jpn float32 = (0.01 * gaji * 11) + (0.01 * gaji)
	response.Jpn = money.NewFromFloat(math.Ceil(float64(jpn)), money.IDR).Display()

	var bpjsKesehatan float32 = (0.01 * gaji * 11) + (0.01 * gaji)
	response.Bpjs_kesehatan = money.NewFromFloat(math.Ceil(float64(bpjsKesehatan)), money.IDR).Display()

	var gajiNettoPertahun float32 = gajiBrutoPerTahun - biayaJabatan - jht - jpn
	response.Gaji_netto_pertahun = money.NewFromFloat(math.Ceil(float64(gajiNettoPertahun)), money.IDR).Display()

	// ptkp
	var ptkp float32

	ptkpStr := strings.ToLower(dataNonStaff.Ptkp)
	if ptkpStr == "k/3" {
		ptkp = 72_000_000
	} else if ptkpStr == "k/2" {
		ptkp = 67_500_000
	} else if ptkpStr == "k/1" {
		ptkp = 63_000_000
	} else if ptkpStr == "k/0" {
		ptkp = 58_500_000
	} else if ptkpStr == "tk/3" {
		ptkp = 67_500_000
	} else if ptkpStr == "tk/2" {
		ptkp = 63_000_000
	} else if ptkpStr == "tk/1" {
		ptkp = 58_500_000
	} else {
		ptkp = 54_000_000
	}

	response.Ptkp = money.NewFromFloat(float64(ptkp), money.IDR).Display()

	var pkp float32
	if gajiNettoPertahun-ptkp > 0 {
		pkp = gajiNettoPertahun - ptkp
	} else {
		pkp = 0
	}

	response.Pkp = money.NewFromFloat(float64(pkp), money.IDR).Display()

	var pph21Pertahun float32

	if pkp <= 50_000_000 {
		pph21Pertahun = 0.05 * pkp
	} else if pkp <= 250_000_000 {
		pph21Pertahun = 25_000_000 + (pkp-50_000_000)*0.15
	} else if pkp <= 500_000_000 {
		pph21Pertahun = 32_500_000 + (pkp-250_000_000)*0.25
	} else {
		pph21Pertahun = 95_000_000 + (pkp-500_000_000)*0.3
	}

	response.Pph_21_pertahun = money.NewFromFloat(float64(pph21Pertahun), money.IDR).Display()

	var pph21PertahunTanpaNpwp float32

	// 8 karakter
	npwp := dataNonStaff.Npwp[:9]
	npwpInt, _ := strconv.Atoi(npwp)
	if npwpInt*1 > 0 {
		pph21PertahunTanpaNpwp = float32(math.Round(float64(pph21Pertahun)))
	} else {
		pph21PertahunTanpaNpwp = (120 / 100) * pph21Pertahun
	}

	response.Pph_21_pertahun_tanpa_npwp = money.NewFromFloat(float64(pph21PertahunTanpaNpwp), money.IDR).Display()

	var pph21Perbulan float32 = pph21PertahunTanpaNpwp / float32(masaKerja)
	response.Pph_21_perbulan = money.NewFromFloat(float64(pph21Perbulan), money.IDR).Display()

	// rekapituasi
	response.Rekapitulasi.Kode = dataNonStaff.NoPayroll
	response.Rekapitulasi.NamaKaryawan = dataNonStaff.Nama
	response.Rekapitulasi.Jabatan = dataNonStaff.Jabatan
	response.Rekapitulasi.KodeTNG = dataNonStaff.Ptkp
	response.Rekapitulasi.Golongan = dataNonStaff.Golongan
	response.Rekapitulasi.GajiPokok = money.NewFromFloat(float64(dataNonStaff.GajiPokok), money.IDR).Display()
	response.Rekapitulasi.TunjanganKemahalan = dataNonStaff.TunjanganKemahalan
	response.Rekapitulasi.TunjanganPerumahan = dataNonStaff.TunjanganPerumahan
	response.Rekapitulasi.TunjanganJabatan = dataNonStaff.TunjanganJabatan
	response.Rekapitulasi.TunjanganLainPph21 = dataNonStaff.TunjanganLainPph21
	response.Rekapitulasi.TKJkk = response.Jkk
	response.Rekapitulasi.TKJkm = response.Jkm
	response.Rekapitulasi.KesehatanP = response.Bpjs_4persen
	response.Rekapitulasi.Lembur = response.Lembur

	var gajiBruto float32

	gajiBruto += dataNonStaff.GajiPokok

	if dataNonStaff.TunjanganKemahalan != nil {
		gajiBruto += *dataNonStaff.TunjanganKemahalan
	}
	if dataNonStaff.TunjanganPerumahan != nil {
		gajiBruto += *dataNonStaff.TunjanganPerumahan
	}
	if dataNonStaff.TunjanganJabatan != nil {
		gajiBruto += *dataNonStaff.TunjanganJabatan
	}
	if dataNonStaff.TunjanganLainPph21 != nil {
		gajiBruto += *dataNonStaff.TunjanganLainPph21
	}

	gajiBruto += jkk
	gajiBruto += jkm
	gajiBruto += bpjs4persen
	gajiBruto += float32(totalGajiLembur)
	response.Rekapitulasi.GajiBruto = money.NewFromFloat(float64(gajiBruto), money.IDR).Display()

	response.Rekapitulasi.PPH21PerBulan = response.Pph_21_perbulan

	var bpjstkjht float32
	bpjstkjht += dataNonStaff.GajiPokok

	if dataNonStaff.TunjanganKemahalan != nil {
		bpjstkjht += *dataNonStaff.TunjanganKemahalan
	}
	if dataNonStaff.TunjanganPerumahan != nil {
		bpjstkjht += *dataNonStaff.TunjanganPerumahan
	}
	if dataNonStaff.TunjanganJabatan != nil {
		bpjstkjht += *dataNonStaff.TunjanganJabatan
	}

	bpjstkjht *= 0.02
	bpjstkjht = float32(math.Ceil(float64(bpjstkjht)))
	response.Rekapitulasi.BPJSTKJHT = money.NewFromFloat(float64(bpjstkjht), money.IDR).Display()

	// bpjstkjpn
	var bpjstkjpn float32
	bpjstkjpn += dataNonStaff.GajiPokok

	if dataNonStaff.TunjanganKemahalan != nil {
		bpjstkjpn += *dataNonStaff.TunjanganKemahalan
	}
	if dataNonStaff.TunjanganPerumahan != nil {
		bpjstkjpn += *dataNonStaff.TunjanganPerumahan
	}
	if dataNonStaff.TunjanganJabatan != nil {
		bpjstkjpn += *dataNonStaff.TunjanganJabatan
	}

	bpjstkjpn *= 0.01
	bpjstkjpn = float32(math.Ceil(float64(bpjstkjpn)))
	response.Rekapitulasi.BPJSTKJPN = money.NewFromFloat(float64(bpjstkjpn), money.IDR).Display()
	response.Rekapitulasi.BPJSKesehatanK = response.Rekapitulasi.BPJSTKJPN
	response.Rekapitulasi.BPJSKesehatanP = money.NewFromFloat(float64(bpjs4persen), money.IDR).Display()

	response.Rekapitulasi.BPJSTKJKK = money.NewFromFloat(float64(jkk), money.IDR).Display()
	response.Rekapitulasi.BPJSTKJKM = money.NewFromFloat(float64(jkm), money.IDR).Display()

	var JumlahPotongan float32
	JumlahPotongan += pph21Perbulan

	// belum ditambahkan dengan pinjamna-pinjaman
	JumlahPotongan += bpjstkjht
	JumlahPotongan += bpjstkjpn
	JumlahPotongan += bpjstkjpn   // bpjs kesehatan k
	JumlahPotongan += bpjs4persen // bpjs kesehatan k
	JumlahPotongan += jkk
	JumlahPotongan += jkm

	response.Rekapitulasi.JumlahPotongan = money.NewFromFloat(float64(JumlahPotongan), money.IDR).Display()
	response.Rekapitulasi.GajiNetto = money.NewFromFloat(float64(math.Ceil(float64(gajiBruto-JumlahPotongan))), money.IDR).Display()
	ctx.JSON(http.StatusOK, response)
}

func (s *Server) generateSlipGaji(ctx *gin.Context) {
	id := ctx.Query("id")
	noPayroll := ctx.Query("no_payroll")
	date := ctx.Query("date")

	dataNonStaff, err := s.store.GetNonStaffByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	_, totalGajiLembur := s.lemburDetailsByID(noPayroll, date)

	var response NonStaffResponse

	// masa kerja
	var masaKerja int32
	if dataNonStaff.AwalKerja == 1 && dataNonStaff.StatusKerja == "aktif" {
		masaKerja = 12
	} else {
		masaKerja = dataNonStaff.AkhirKerja - dataNonStaff.AwalKerja + 1
	}

	// masa kerja
	response.MasaKerja = masaKerja

	// data non staff
	response.NonStaff = dataNonStaff

	// rekap lembur
	// make sure lembur exists
	var gaji float32
	if totalGajiLembur != 0 {
		response.Lembur = money.NewFromFloat(totalGajiLembur, money.IDR).Display()

		gaji = float32(dataNonStaff.GajiPokok)
		if dataNonStaff.TunjanganKemahalan != nil {
			gaji += *dataNonStaff.TunjanganKemahalan
		}
		if dataNonStaff.TunjanganPerumahan != nil {
			gaji += *dataNonStaff.TunjanganPerumahan
		}
		if dataNonStaff.TunjanganJabatan != nil {
			gaji += *dataNonStaff.TunjanganJabatan
		}

		// gaji
		response.Gaji = money.NewFromFloat(float64(gaji), money.IDR).Display()
	} else {
		gaji = float32(dataNonStaff.GajiPokok)
		response.Gaji = money.NewFromFloat(float64(gaji), money.IDR).Display()
	}

	// GajiLemburTunjangan
	var gajiLemburTunjangan float32
	gajiLemburTunjangan += dataNonStaff.GajiPokok
	if dataNonStaff.TunjanganKemahalan != nil {
		gajiLemburTunjangan += *dataNonStaff.TunjanganKemahalan
	}
	if dataNonStaff.TunjanganPerumahan != nil {
		gajiLemburTunjangan += *dataNonStaff.TunjanganPerumahan
	}
	if dataNonStaff.TunjanganJabatan != nil {
		gajiLemburTunjangan += *dataNonStaff.TunjanganJabatan
	}
	if dataNonStaff.TunjanganLainPph21 != nil {
		gajiLemburTunjangan += *dataNonStaff.TunjanganLainPph21
	}

	gajiLemburTunjangan += float32(totalGajiLembur)
	response.GajiLemburTunjangan = money.NewFromFloat(float64(gajiLemburTunjangan), money.IDR).Display()

	jkk := gaji * 0.54 / 100
	response.Jkk = money.NewFromFloat(float64(jkk), money.IDR).Display()

	jkm := gaji * 0.3 / 100
	response.Jkm = money.NewFromFloat(float64(jkm), money.IDR).Display()

	var bpjs4persen float32
	if gaji > 8000000 {
		bpjs4persen = 8000000 * 0.04
	} else {
		bpjs4persen = gaji * 0.04
	}
	response.Bpjs_4persen = money.NewFromFloat(float64(bpjs4persen), money.IDR).Display()

	var gajiBrutoPerBulan float32
	gajiBrutoPerBulan += gajiLemburTunjangan
	gajiBrutoPerBulan += jkk
	gajiBrutoPerBulan += jkm
	gajiBrutoPerBulan += bpjs4persen

	response.Penghasilan_bruto_perbulan = money.NewFromFloat(float64(gajiBrutoPerBulan), money.IDR).Display()

	var gajiBrutoPerTahun float32 = (gajiBrutoPerBulan * 11) + gajiBrutoPerBulan
	response.Penghasilan_bruto_pertahun = money.NewFromFloat(float64(gajiBrutoPerTahun), money.IDR).Display()

	var biayaJabatan float32
	if gajiBrutoPerTahun*0.05 > float32(500_000*masaKerja) {
		biayaJabatan = float32(500_000 * masaKerja)
	} else {
		biayaJabatan = 0.05 * gajiBrutoPerTahun
	}
	response.Biaya_jabatan = money.NewFromFloat(float64(biayaJabatan), money.IDR).Display()

	var jht float32 = (0.02 * gaji * 11) + (0.02 * gaji)
	response.Jht = money.NewFromFloat(math.Ceil(float64(jht)), money.IDR).Display()

	var jpn float32 = (0.01 * gaji * 11) + (0.01 * gaji)
	response.Jpn = money.NewFromFloat(math.Ceil(float64(jpn)), money.IDR).Display()

	var bpjsKesehatan float32 = (0.01 * gaji * 11) + (0.01 * gaji)
	response.Bpjs_kesehatan = money.NewFromFloat(math.Ceil(float64(bpjsKesehatan)), money.IDR).Display()

	var gajiNettoPertahun float32 = gajiBrutoPerTahun - biayaJabatan - jht - jpn
	response.Gaji_netto_pertahun = money.NewFromFloat(math.Ceil(float64(gajiNettoPertahun)), money.IDR).Display()

	// ptkp
	var ptkp float32

	ptkpStr := strings.ToLower(dataNonStaff.Ptkp)
	if ptkpStr == "k/3" {
		ptkp = 72_000_000
	} else if ptkpStr == "k/2" {
		ptkp = 67_500_000
	} else if ptkpStr == "k/1" {
		ptkp = 63_000_000
	} else if ptkpStr == "k/0" {
		ptkp = 58_500_000
	} else if ptkpStr == "tk/3" {
		ptkp = 67_500_000
	} else if ptkpStr == "tk/2" {
		ptkp = 63_000_000
	} else if ptkpStr == "tk/1" {
		ptkp = 58_500_000
	} else {
		ptkp = 54_000_000
	}

	response.Ptkp = money.NewFromFloat(float64(ptkp), money.IDR).Display()

	var pkp float32
	if gajiNettoPertahun-ptkp > 0 {
		pkp = gajiNettoPertahun - ptkp
	} else {
		pkp = 0
	}

	response.Pkp = money.NewFromFloat(float64(pkp), money.IDR).Display()

	var pph21Pertahun float32

	if pkp <= 50_000_000 {
		pph21Pertahun = 0.05 * pkp
	} else if pkp <= 250_000_000 {
		pph21Pertahun = 25_000_000 + (pkp-50_000_000)*0.15
	} else if pkp <= 500_000_000 {
		pph21Pertahun = 32_500_000 + (pkp-250_000_000)*0.25
	} else {
		pph21Pertahun = 95_000_000 + (pkp-500_000_000)*0.3
	}

	response.Pph_21_pertahun = money.NewFromFloat(float64(pph21Pertahun), money.IDR).Display()

	var pph21PertahunTanpaNpwp float32

	// 8 karakter
	npwp := dataNonStaff.Npwp[:9]
	npwpInt, _ := strconv.Atoi(npwp)
	if npwpInt*1 > 0 {
		pph21PertahunTanpaNpwp = float32(math.Round(float64(pph21Pertahun)))
	} else {
		pph21PertahunTanpaNpwp = (120 / 100) * pph21Pertahun
	}

	response.Pph_21_pertahun_tanpa_npwp = money.NewFromFloat(float64(pph21PertahunTanpaNpwp), money.IDR).Display()

	var pph21Perbulan float32 = pph21PertahunTanpaNpwp / float32(masaKerja)
	response.Pph_21_perbulan = money.NewFromFloat(float64(pph21Perbulan), money.IDR).Display()

	// rekapituasi
	response.Rekapitulasi.Kode = dataNonStaff.NoPayroll
	response.Rekapitulasi.NamaKaryawan = dataNonStaff.Nama
	response.Rekapitulasi.Jabatan = dataNonStaff.Jabatan
	response.Rekapitulasi.KodeTNG = dataNonStaff.Ptkp
	response.Rekapitulasi.Golongan = dataNonStaff.Golongan
	response.Rekapitulasi.GajiPokok = money.NewFromFloat(float64(dataNonStaff.GajiPokok), money.IDR).Display()
	response.Rekapitulasi.TunjanganKemahalan = dataNonStaff.TunjanganKemahalan
	response.Rekapitulasi.TunjanganPerumahan = dataNonStaff.TunjanganPerumahan
	response.Rekapitulasi.TunjanganJabatan = dataNonStaff.TunjanganJabatan
	response.Rekapitulasi.TunjanganLainPph21 = dataNonStaff.TunjanganLainPph21
	response.Rekapitulasi.TKJkk = response.Jkk
	response.Rekapitulasi.TKJkm = response.Jkm
	response.Rekapitulasi.KesehatanP = response.Bpjs_4persen
	response.Rekapitulasi.Lembur = response.Lembur

	var gajiBruto float32

	gajiBruto += dataNonStaff.GajiPokok

	if dataNonStaff.TunjanganKemahalan != nil {
		gajiBruto += *dataNonStaff.TunjanganKemahalan
	}
	if dataNonStaff.TunjanganPerumahan != nil {
		gajiBruto += *dataNonStaff.TunjanganPerumahan
	}
	if dataNonStaff.TunjanganJabatan != nil {
		gajiBruto += *dataNonStaff.TunjanganJabatan
	}
	if dataNonStaff.TunjanganLainPph21 != nil {
		gajiBruto += *dataNonStaff.TunjanganLainPph21
	}

	gajiBruto += jkk
	gajiBruto += jkm
	gajiBruto += bpjs4persen
	gajiBruto += float32(totalGajiLembur)
	response.Rekapitulasi.GajiBruto = money.NewFromFloat(float64(gajiBruto), money.IDR).Display()

	response.Rekapitulasi.PPH21PerBulan = response.Pph_21_perbulan

	var bpjstkjht float32
	bpjstkjht += dataNonStaff.GajiPokok

	if dataNonStaff.TunjanganKemahalan != nil {
		bpjstkjht += *dataNonStaff.TunjanganKemahalan
	}
	if dataNonStaff.TunjanganPerumahan != nil {
		bpjstkjht += *dataNonStaff.TunjanganPerumahan
	}
	if dataNonStaff.TunjanganJabatan != nil {
		bpjstkjht += *dataNonStaff.TunjanganJabatan
	}

	bpjstkjht *= 0.02
	bpjstkjht = float32(math.Ceil(float64(bpjstkjht)))
	response.Rekapitulasi.BPJSTKJHT = money.NewFromFloat(float64(bpjstkjht), money.IDR).Display()

	// bpjstkjpn
	var bpjstkjpn float32
	bpjstkjpn += dataNonStaff.GajiPokok

	if dataNonStaff.TunjanganKemahalan != nil {
		bpjstkjpn += *dataNonStaff.TunjanganKemahalan
	}
	if dataNonStaff.TunjanganPerumahan != nil {
		bpjstkjpn += *dataNonStaff.TunjanganPerumahan
	}
	if dataNonStaff.TunjanganJabatan != nil {
		bpjstkjpn += *dataNonStaff.TunjanganJabatan
	}

	bpjstkjpn *= 0.01
	bpjstkjpn = float32(math.Ceil(float64(bpjstkjpn)))
	response.Rekapitulasi.BPJSTKJPN = money.NewFromFloat(float64(bpjstkjpn), money.IDR).Display()
	response.Rekapitulasi.BPJSKesehatanK = response.Rekapitulasi.BPJSTKJPN
	response.Rekapitulasi.BPJSKesehatanP = money.NewFromFloat(float64(bpjs4persen), money.IDR).Display()

	response.Rekapitulasi.BPJSTKJKK = money.NewFromFloat(float64(jkk), money.IDR).Display()
	response.Rekapitulasi.BPJSTKJKM = money.NewFromFloat(float64(jkm), money.IDR).Display()

	var JumlahPotongan float32
	JumlahPotongan += pph21Perbulan

	// belum ditambahkan dengan pinjamna-pinjaman
	JumlahPotongan += bpjstkjht
	JumlahPotongan += bpjstkjpn
	JumlahPotongan += bpjstkjpn   // bpjs kesehatan k
	JumlahPotongan += bpjs4persen // bpjs kesehatan k
	JumlahPotongan += jkk
	JumlahPotongan += jkm

	response.Rekapitulasi.JumlahPotongan = money.NewFromFloat(float64(JumlahPotongan), money.IDR).Display()
	response.Rekapitulasi.GajiNetto = money.NewFromFloat(float64(math.Ceil(float64(gajiBruto-JumlahPotongan))), money.IDR).Display()

	// slip pdf
	pdf := gofpdf.New("P", "mm", "letter", "")
	pdf.AddPage()
	pdf.SetFont("helvetica", "B", 16)
	width, _ := pdf.GetPageSize()
	pdf.CellFormat(width, 5, "SLIP GAJI                 ", "", 1, "C", false, 0, "")
	pdf.SetFont("helvetica", "B", 12)
	pdf.CellFormat(width, 8, "PT. Hutahaean Group                        ", "", 1, "C", false, 0, "")
	pdf.Line(10, pdf.GetY()+5, width-10, pdf.GetY()+5)
	pdf.Ln(15)
	pdf.SetFont("Arial", "", 12)
	employeeDetails := [][]string{
		{"Tahun-Bulan", fmt.Sprintf(": %s", util.FormatDate(date)), "", "", "Tunjangan Kemahalan", "", fmt.Sprintf(": %s", isNilMoney(response.TunjanganKemahalan))},
		{"Nama", fmt.Sprintf(": %s", response.NamaKaryawan), "", "", "Tunjangan Perumahan", "", fmt.Sprintf(": %s", isNilMoney(response.TunjanganPerumahan))},
		{"Jabatan", fmt.Sprintf(": %s", response.Jabatan), "", "", "Tunjangan Jabatan", "", fmt.Sprintf(": %s", isNilMoney(response.TunjanganJabatan))},
		{"Status", fmt.Sprintf(": %s", response.NonStaff.StatusKerja), "", "", "Tunjangan Lain-lain", "", fmt.Sprintf(": %s", isNilMoney(response.TunjanganLainPph21))},
		{"PTKP", fmt.Sprintf(": %s", response.Ptkp)},
	}

	for _, row := range employeeDetails {
		for _, col := range row {
			pdf.CellFormat(25, 7, col, "", 0, "L", false, 0, "")
		}
		pdf.Ln(-1)
	}

	pdf.Line(10, pdf.GetY()+5, width-10, pdf.GetY()+5)
	pdf.Ln(15)

	employeeDetails2 := [][]string{
		{"", "", "", "", "Lembur", "", fmt.Sprintf(": %s", response.Lembur)},
		{"", "", "", "", "Jumlah Potongan", "", fmt.Sprintf(": %s", response.Rekapitulasi.JumlahPotongan)},
		{"", "", "", "", "Gaji Bruto", "", fmt.Sprintf(": %s", response.GajiBruto)},
		{"", "", "", "", "Gaji Netto", "", fmt.Sprintf(": %s", response.GajiNetto)},
	}

	// Write the table data.
	for i, row := range employeeDetails2 {
		if i == len(employeeDetails2)-1 {
			for _, col := range row {
				pdf.SetFont("Arial", "B", 12)
				pdf.CellFormat(25, 6, col, "0", 0, "L", false, 0, "")
			}
			continue
		}

		for _, col := range row {
			pdf.CellFormat(25, 6, col, "0", 0, "L", false, 0, "")
		}

		pdf.Ln(-1)
	}

	_ = pdf.Output(ctx.Writer)
}

func isNilMoney(data *float32) string {
	if data == nil {
		return "-"
	}

	formatData := money.NewFromFloat(float64(*data), money.IDR)
	return formatData.Display()
}

func (server *Server) RekapitulasiNonStaff(date string) ([]NonStaffResponse, error) {
	nonStaff, _ := server.store.ListNonStaff(context.Background())
	responses := make([]NonStaffResponse, 0, len(nonStaff))

	for _, v := range nonStaff {
		id := v.ID
		noPayroll := v.NoPayroll

		dataNonStaff, err := server.store.GetNonStaffByID(context.Background(), id)
		if err != nil {
			return nil, err
		}

		_, totalGajiLembur := server.lemburDetailsByID(noPayroll, date)

		var response NonStaffResponse

		// masa kerja
		var masaKerja int32
		if dataNonStaff.AwalKerja == 1 && dataNonStaff.StatusKerja == "aktif" {
			masaKerja = 12
		} else {
			masaKerja = dataNonStaff.AkhirKerja - dataNonStaff.AwalKerja + 1
		}

		// masa kerja
		response.MasaKerja = masaKerja

		// data non staff
		response.NonStaff = dataNonStaff

		// rekap lembur
		// make sure lembur exists
		var gaji float32
		if totalGajiLembur != 0 {
			response.Lembur = money.NewFromFloat(totalGajiLembur, money.IDR).Display()

			gaji = float32(dataNonStaff.GajiPokok)
			if dataNonStaff.TunjanganKemahalan != nil {
				gaji += *dataNonStaff.TunjanganKemahalan
			}
			if dataNonStaff.TunjanganPerumahan != nil {
				gaji += *dataNonStaff.TunjanganPerumahan
			}
			if dataNonStaff.TunjanganJabatan != nil {
				gaji += *dataNonStaff.TunjanganJabatan
			}

			// gaji
			response.Gaji = money.NewFromFloat(float64(gaji), money.IDR).Display()
		} else {
			gaji = float32(dataNonStaff.GajiPokok)
			response.Gaji = money.NewFromFloat(float64(gaji), money.IDR).Display()
		}

		// GajiLemburTunjangan
		var gajiLemburTunjangan float32
		gajiLemburTunjangan += dataNonStaff.GajiPokok
		if dataNonStaff.TunjanganKemahalan != nil {
			gajiLemburTunjangan += *dataNonStaff.TunjanganKemahalan
		}
		if dataNonStaff.TunjanganPerumahan != nil {
			gajiLemburTunjangan += *dataNonStaff.TunjanganPerumahan
		}
		if dataNonStaff.TunjanganJabatan != nil {
			gajiLemburTunjangan += *dataNonStaff.TunjanganJabatan
		}
		if dataNonStaff.TunjanganLainPph21 != nil {
			gajiLemburTunjangan += *dataNonStaff.TunjanganLainPph21
		}

		gajiLemburTunjangan += float32(totalGajiLembur)
		response.GajiLemburTunjangan = money.NewFromFloat(float64(gajiLemburTunjangan), money.IDR).Display()

		jkk := gaji * 0.54 / 100
		response.Jkk = money.NewFromFloat(float64(jkk), money.IDR).Display()

		jkm := gaji * 0.3 / 100
		response.Jkm = money.NewFromFloat(float64(jkm), money.IDR).Display()

		var bpjs4persen float32
		if gaji > 8000000 {
			bpjs4persen = 8000000 * 0.04
		} else {
			bpjs4persen = gaji * 0.04
		}
		response.Bpjs_4persen = money.NewFromFloat(float64(bpjs4persen), money.IDR).Display()

		var gajiBrutoPerBulan float32
		gajiBrutoPerBulan += gajiLemburTunjangan
		gajiBrutoPerBulan += jkk
		gajiBrutoPerBulan += jkm
		gajiBrutoPerBulan += bpjs4persen

		response.Penghasilan_bruto_perbulan = money.NewFromFloat(float64(gajiBrutoPerBulan), money.IDR).Display()

		var gajiBrutoPerTahun float32 = (gajiBrutoPerBulan * 11) + gajiBrutoPerBulan
		response.Penghasilan_bruto_pertahun = money.NewFromFloat(float64(gajiBrutoPerTahun), money.IDR).Display()

		var biayaJabatan float32
		if gajiBrutoPerTahun*0.05 > float32(500_000*masaKerja) {
			biayaJabatan = float32(500_000 * masaKerja)
		} else {
			biayaJabatan = 0.05 * gajiBrutoPerTahun
		}
		response.Biaya_jabatan = money.NewFromFloat(float64(biayaJabatan), money.IDR).Display()

		var jht float32 = (0.02 * gaji * 11) + (0.02 * gaji)
		response.Jht = money.NewFromFloat(math.Ceil(float64(jht)), money.IDR).Display()

		var jpn float32 = (0.01 * gaji * 11) + (0.01 * gaji)
		response.Jpn = money.NewFromFloat(math.Ceil(float64(jpn)), money.IDR).Display()

		var bpjsKesehatan float32 = (0.01 * gaji * 11) + (0.01 * gaji)
		response.Bpjs_kesehatan = money.NewFromFloat(math.Ceil(float64(bpjsKesehatan)), money.IDR).Display()

		var gajiNettoPertahun float32 = gajiBrutoPerTahun - biayaJabatan - jht - jpn
		response.Gaji_netto_pertahun = money.NewFromFloat(math.Ceil(float64(gajiNettoPertahun)), money.IDR).Display()

		// ptkp
		var ptkp float32

		ptkpStr := strings.ToLower(dataNonStaff.Ptkp)
		if ptkpStr == "k/3" {
			ptkp = 72_000_000
		} else if ptkpStr == "k/2" {
			ptkp = 67_500_000
		} else if ptkpStr == "k/1" {
			ptkp = 63_000_000
		} else if ptkpStr == "k/0" {
			ptkp = 58_500_000
		} else if ptkpStr == "tk/3" {
			ptkp = 67_500_000
		} else if ptkpStr == "tk/2" {
			ptkp = 63_000_000
		} else if ptkpStr == "tk/1" {
			ptkp = 58_500_000
		} else {
			ptkp = 54_000_000
		}

		response.Ptkp = money.NewFromFloat(float64(ptkp), money.IDR).Display()

		var pkp float32
		if gajiNettoPertahun-ptkp > 0 {
			pkp = gajiNettoPertahun - ptkp
		} else {
			pkp = 0
		}

		response.Pkp = money.NewFromFloat(float64(pkp), money.IDR).Display()

		var pph21Pertahun float32

		if pkp <= 50_000_000 {
			pph21Pertahun = 0.05 * pkp
		} else if pkp <= 250_000_000 {
			pph21Pertahun = 25_000_000 + (pkp-50_000_000)*0.15
		} else if pkp <= 500_000_000 {
			pph21Pertahun = 32_500_000 + (pkp-250_000_000)*0.25
		} else {
			pph21Pertahun = 95_000_000 + (pkp-500_000_000)*0.3
		}

		response.Pph_21_pertahun = money.NewFromFloat(float64(pph21Pertahun), money.IDR).Display()

		var pph21PertahunTanpaNpwp float32

		// 8 karakter
		npwp := dataNonStaff.Npwp[:9]
		npwpInt, _ := strconv.Atoi(npwp)
		if npwpInt*1 > 0 {
			pph21PertahunTanpaNpwp = float32(math.Round(float64(pph21Pertahun)))
		} else {
			pph21PertahunTanpaNpwp = (120 / 100) * pph21Pertahun
		}

		response.Pph_21_pertahun_tanpa_npwp = money.NewFromFloat(float64(pph21PertahunTanpaNpwp), money.IDR).Display()

		var pph21Perbulan float32 = pph21PertahunTanpaNpwp / float32(masaKerja)
		response.Pph_21_perbulan = money.NewFromFloat(float64(pph21Perbulan), money.IDR).Display()

		// rekapituasi
		response.Rekapitulasi.Kode = dataNonStaff.NoPayroll
		response.Rekapitulasi.NamaKaryawan = dataNonStaff.Nama
		response.Rekapitulasi.Jabatan = dataNonStaff.Jabatan
		response.Rekapitulasi.KodeTNG = dataNonStaff.Ptkp
		response.Rekapitulasi.Golongan = dataNonStaff.Golongan
		response.Rekapitulasi.GajiPokok = money.NewFromFloat(float64(dataNonStaff.GajiPokok), money.IDR).Display()
		response.Rekapitulasi.TunjanganKemahalan = dataNonStaff.TunjanganKemahalan
		response.Rekapitulasi.TunjanganPerumahan = dataNonStaff.TunjanganPerumahan
		response.Rekapitulasi.TunjanganJabatan = dataNonStaff.TunjanganJabatan
		response.Rekapitulasi.TunjanganLainPph21 = dataNonStaff.TunjanganLainPph21
		response.Rekapitulasi.TKJkk = response.Jkk
		response.Rekapitulasi.TKJkm = response.Jkm
		response.Rekapitulasi.KesehatanP = response.Bpjs_4persen
		response.Rekapitulasi.Lembur = response.Lembur

		var gajiBruto float32

		gajiBruto += dataNonStaff.GajiPokok

		if dataNonStaff.TunjanganKemahalan != nil {
			gajiBruto += *dataNonStaff.TunjanganKemahalan
		}
		if dataNonStaff.TunjanganPerumahan != nil {
			gajiBruto += *dataNonStaff.TunjanganPerumahan
		}
		if dataNonStaff.TunjanganJabatan != nil {
			gajiBruto += *dataNonStaff.TunjanganJabatan
		}
		if dataNonStaff.TunjanganLainPph21 != nil {
			gajiBruto += *dataNonStaff.TunjanganLainPph21
		}

		gajiBruto += jkk
		gajiBruto += jkm
		gajiBruto += bpjs4persen
		gajiBruto += float32(totalGajiLembur)
		response.Rekapitulasi.GajiBruto = money.NewFromFloat(float64(gajiBruto), money.IDR).Display()

		response.Rekapitulasi.PPH21PerBulan = response.Pph_21_perbulan

		var bpjstkjht float32
		bpjstkjht += dataNonStaff.GajiPokok

		if dataNonStaff.TunjanganKemahalan != nil {
			bpjstkjht += *dataNonStaff.TunjanganKemahalan
		}
		if dataNonStaff.TunjanganPerumahan != nil {
			bpjstkjht += *dataNonStaff.TunjanganPerumahan
		}
		if dataNonStaff.TunjanganJabatan != nil {
			bpjstkjht += *dataNonStaff.TunjanganJabatan
		}

		bpjstkjht *= 0.02
		bpjstkjht = float32(math.Ceil(float64(bpjstkjht)))
		response.Rekapitulasi.BPJSTKJHT = money.NewFromFloat(float64(bpjstkjht), money.IDR).Display()

		var bpjstkjpn float32
		bpjstkjpn += dataNonStaff.GajiPokok

		if dataNonStaff.TunjanganKemahalan != nil {
			bpjstkjpn += *dataNonStaff.TunjanganKemahalan
		}
		if dataNonStaff.TunjanganPerumahan != nil {
			bpjstkjpn += *dataNonStaff.TunjanganPerumahan
		}
		if dataNonStaff.TunjanganJabatan != nil {
			bpjstkjpn += *dataNonStaff.TunjanganJabatan
		}

		bpjstkjpn *= 0.01
		bpjstkjpn = float32(math.Ceil(float64(bpjstkjpn)))
		response.Rekapitulasi.BPJSTKJPN = money.NewFromFloat(float64(bpjstkjpn), money.IDR).Display()
		response.Rekapitulasi.BPJSKesehatanK = response.Rekapitulasi.BPJSTKJPN
		response.Rekapitulasi.BPJSKesehatanP = money.NewFromFloat(float64(bpjs4persen), money.IDR).Display()

		response.Rekapitulasi.BPJSTKJKK = money.NewFromFloat(float64(jkk), money.IDR).Display()
		response.Rekapitulasi.BPJSTKJKM = money.NewFromFloat(float64(jkm), money.IDR).Display()

		var JumlahPotongan float32
		JumlahPotongan += pph21Perbulan

		// belum ditambahkan dengan pinjaman-pinjaman
		JumlahPotongan += bpjstkjht
		JumlahPotongan += bpjstkjpn
		JumlahPotongan += bpjstkjpn   // bpjs kesehatan k
		JumlahPotongan += bpjs4persen // bpjs kesehatan k
		JumlahPotongan += jkk
		JumlahPotongan += jkm

		response.Rekapitulasi.JumlahPotongan = money.NewFromFloat(float64(JumlahPotongan), money.IDR).Display()
		response.Rekapitulasi.GajiNetto = money.NewFromFloat(math.Ceil(float64(gajiBruto-JumlahPotongan)), money.IDR).Display()
		responses = append(responses, response)
	}

	return responses, nil
}
