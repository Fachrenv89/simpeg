package api

import (
	"fmt"

	db "github.com/SemmiDev/simpeg/db/mysql"
	"github.com/SemmiDev/simpeg/util"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker util.Maker
	router     *gin.Engine
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(
	config util.Config,
	store db.Store,
	tokenMaker util.Maker,
) (*Server, error) {
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.Use(cors.Default())

	router.POST("/api/users/login", server.loginUser)
	router.POST("/api/users/register", server.createUser)
	router.PUT("/api/users/change-password", server.changePassword)
	router.PUT("/api/users/profile", server.updateUserProfile)

	router.POST("/api/lembur", server.createLembur)
	router.GET("/api/lembur", server.listLembur)
	router.GET("/api/lembur/:id", server.lemburDetails)
	router.DELETE("/api/lembur/:id", server.deleteLembur)

	router.GET("/api/non-staff", server.listNonStaff)
	router.GET("/api/non-staff/details", server.detailNonStaff)
	router.GET("/api/non-staff/slip-gaji", server.generateSlipGaji)
	router.GET("/api/non-staff/total", server.countNonStaff)
	router.POST("/api/non-staff", server.createNonStaff)
	router.DELETE("/api/non-staff/:id", server.deleteNonStaff)

	router.GET("/api/non-staff/edit", server.detailNonStaff)
	router.PUT("/api/non-staff/edit", server.editNonStaff)

	router.POST("/api/tanggal_lembur", server.createTanggalLembur)
	router.GET("/api/generate-report", func(c *gin.Context) {
		date := c.Query("date")
		nonStaff, _ := server.RekapitulasiNonStaff(date)
		fileName := fmt.Sprintf("Rekapitulasi | Non Staff | %s", date)

		rekapitulasi := make([]Rekapitulasi, 0, len(nonStaff))
		for _, v := range nonStaff {
			rekapitulasi = append(rekapitulasi, v.Rekapitulasi)
		}

		err := GenerateExcelReport(c.Writer, fileName, rekapitulasi)
		if err != nil {
			c.JSON(500, errorResponse(err))
		}
	})

	server.router = router
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(":" + address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
