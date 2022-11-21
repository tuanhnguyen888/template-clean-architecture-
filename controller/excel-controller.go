package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
	"github.com/tuanhnguyen888/server/entity"
	"github.com/tuanhnguyen888/server/helper"
	"github.com/tuanhnguyen888/server/service"
	"gorm.io/gorm"
)

type ExcelController interface {
	ExportServer(ctx *gin.Context)
	ImportServer(ctx *gin.Context)
}

type excelController struct {
	DB      *gorm.DB
	service service.ServerService
}

func NewExcelController(db *gorm.DB, service service.ServerService) ExcelController {
	return &excelController{
		DB:      db,
		service: service,
	}
}

// ExportServer implements ExcelController
func (c *excelController) ExportServer(ctx *gin.Context) {
	servers := []entity.Server{}
	payload := struct {
		From  int    `json:"from"`
		To    int    `json:"to"`
		Field string `json:"field"`
		Kind  string `json:"kind"`
	}{}

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
	}

	if payload.From > payload.To {
		response := helper.BuildErrorResponse("Failed to process request", "invalid page number", helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	if payload.From <= 0 {
		payload.From = 1
	}

	offset := payload.From - 1
	limit := (payload.To - payload.From + 1)
	if payload.Field != "" {
		c.DB.Order(payload.Field + " " + payload.Kind).Offset(offset * 10).Limit(limit * 10).Find(&servers)
	} else {
		c.DB.Offset((payload.From - 1) * 10).Limit(limit).Find(&servers)
	}

	f := excelize.NewFile()

	index := f.NewSheet("Sheet1")
	f.SetCellValue("Sheet1", "A1", "ServerName")
	f.SetCellValue("Sheet1", "B1", "Status")
	f.SetCellValue("Sheet1", "C1", "Ipv4")
	f.SetCellValue("Sheet1", "D1", "CreateTime")
	f.SetCellValue("Sheet1", "E1", "UpdateTime")
	// set trang hoat donog
	f.SetActiveSheet(index)

	for i, server := range servers {
		SNByte, err := json.Marshal(server.Name)
		if err != nil {
			continue
		}
		StatusByte, err := json.Marshal(server.Status)
		if err != nil {
			continue
		}
		ipv4Byte, err := json.Marshal(server.Ipv4)
		if err != nil {
			continue
		}
		createTime, err := json.Marshal(time.UnixMilli(server.CreatedAt))
		if err != nil {
			continue
		}
		updateTime, err := json.Marshal(time.UnixMilli(server.UpdatedAt))
		if err != nil {
			continue
		}
		f.SetCellValue("Sheet1", "A"+strconv.Itoa(i+2), string(SNByte))
		f.SetCellValue("Sheet1", "B"+strconv.Itoa(i+2), string(StatusByte))
		f.SetCellValue("Sheet1", "C"+strconv.Itoa(i+2), string(ipv4Byte))
		f.SetCellValue("Sheet1", "D"+strconv.Itoa(i+2), string(createTime))
		f.SetCellValue("Sheet1", "E"+strconv.Itoa(i+2), string(updateTime))
	}

	if err := f.SaveAs("server.xlsx"); err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "okk",
		"file":    "server.xlsx",
	})

	// ctx.FileAttachment("data/excel/server.xlsx")
}

// ImportServer implements ExcelController
func (c *excelController) ImportServer(ctx *gin.Context) {
	// file, err := ctx.FormFile("file")
	// if err != nil {
	// 	response := helper.BuildErrorResponse("Failed to process ", err.Error(), helper.EmptyObj{})
	// 	ctx.JSON(http.StatusBadRequest, response)
	// 	return
	// }

	xlsx, err := excelize.OpenFile("listOfServers.xlsx")
	if err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	rows := xlsx.GetRows("servers")

	var strGoodImport []string

	var errImports []string

	// connect db

	allServers := []entity.Server{}
	c.DB.Find(&allServers)

	for i := 1; i < (len(rows)); i++ {
		server := entity.Server{}
		server.Name = rows[i][0]
		server.Ipv4 = rows[i][1]

		// if err != nil {
		// 	// ErrorLogger.Printf(" %s - %s", *server.Name, *server.Ipv4)
		// 	errImports = append(errImports, fmt.Sprintf(" %s - %s", server.Name, server.Ipv4))
		// 	continue
		// }

		// _, err = exec.Command("ping", *server.Ipv4).Output()
		// if err != nil {
		// 	server.Status = false
		// } else {
		// 	server.Status = true
		// }

		server.CreatedAt = time.Now().UnixMilli()
		server.UpdatedAt = time.Now().UnixMilli()

		err = c.DB.Create(&server).Error
		if err != nil {
			// ErrorLogger.Printf(" %s - %s", *server.Name, *server.Ipv4)
			errImports = append(errImports, fmt.Sprintf(" %s - %s", server.Name, server.Ipv4))
		} else {
			strGoodImport = append(strGoodImport, fmt.Sprintf("%s - %s ", server.Name, server.Ipv4))
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":                   "servers has been added by Excel",
		"numbers of success server": len(strGoodImport),
		"servers added success":     strGoodImport,
		"numbers of error groups":   len(errImports),
		"servers added error":       errImports,
	})
}
