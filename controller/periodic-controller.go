package controller

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/jasonlvhit/gocron"
	"github.com/joho/godotenv"
	"github.com/olivere/elastic/v7"
	"github.com/tuanhnguyen888/server/entity"
	"github.com/tuanhnguyen888/server/service"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
)

type PeriodicController interface {
	Cron()
	UpdateServerPeriodic()
	SendEmailDaily()
	SaveDataByRedis()
	CustomSendEmail(ctx *gin.Context)
}
type periodicController struct {
	service       service.ServerService
	DB            *gorm.DB
	LogRorate     service.LogRorate
	RedisClient   *redis.Client
	ElasticClient *elastic.Client
}

type checkServer struct {
	Name   string `json:"name"`
	Status bool   `json:"status" gorm:"default:false"`
	Time   string `json:"time"`
}

func NewPeriodicController(service service.ServerService, logRorate service.LogRorate, db *gorm.DB, redis *redis.Client, elas *elastic.Client) PeriodicController {
	return &periodicController{
		service:       service,
		LogRorate:     logRorate,
		DB:            db,
		RedisClient:   redis,
		ElasticClient: elas,
	}
}

// Cron implements PeriodicController
func (c *periodicController) Cron() {
	c.LogRorate.Info("...")
	// gocron.Every(5).Second().Do(c.UpdateServerPeriodic)
	// gocron.Every(5).Second().Do(c.SaveDataByRedis)
	gocron.Every(5).Second().Do(c.SendEmailDaily)

	<-gocron.Start()
}

// CustomSendEmail implements PeriodicController
func (*periodicController) CustomSendEmail(ctx *gin.Context) {
	panic("unimplemented")
}

// SaveDataByRedis implements PeriodicController
func (c *periodicController) SaveDataByRedis() {
	dbServers := []entity.Server{}
	c.DB.Find(&dbServers)

	now := time.Now().Add(-24 * time.Hour)
	date := strconv.Itoa(now.Day()) + "/" + strconv.Itoa(int(now.Month())) + "/" + strconv.Itoa(now.Year())

	cachedServers, err := json.Marshal(dbServers)
	if err != nil {
		c.LogRorate.Error("Can not save date day %s", now.Format("01-02-2006"))
		return
	}

	err = c.RedisClient.Set(date, cachedServers, 60*24*time.Hour).Err()
	if err != nil {
		c.LogRorate.Error("Can not cache data day %s", now.Format("01-02-2006"))
		return
	}

	c.LogRorate.Info("Cache success data day %s", now.Format("01-02-2006"))
}

// SendEmailDaily implements PeriodicController
func (c *periodicController) SendEmailDaily() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	mail := os.Getenv("EMAIL_ACCOUNT")
	// pwd := os.Getenv("EMAIL_PASSWPRD")

	servers := []entity.Server{}

	c.DB.Find(&servers)
	serverOn := 0
	serverOff := 0
	msg1 := ""
	for _, server := range servers {
		if server.Status {
			serverOn++
		} else {
			serverOff++
		}
		// uptime
		// ctx := context.Background()
		CheckServer := []checkServer{}

		// searchSource := elastic.NewSearchSource()
		// searchSource.Query(elastic.NewMatchQuery("name", *server.Name))
		// searchSource.Query(elastic.NewMatchQuery("time", time.Now().Add(12*time.Hour).Format("02-03-2006")))

		// searchService := elasc.Search().Index("server").SearchSource(searchSource)
		// searchResult, err := searchService.Do(ctx)

		termQuery := elastic.NewTermQuery("name", server.Name)
		scroller := c.ElasticClient.Scroll().
			Index("server").
			Query(termQuery).
			Size(1)

		for {

			res, err := scroller.Do(context.TODO())
			if err == io.EOF {
				log.Println(err)
				break
			}
			for _, hit := range res.Hits.Hits {
				serverEmp := checkServer{}
				err := json.Unmarshal(hit.Source, &serverEmp)
				if err != nil {
					c.LogRorate.Error("[Getting Students][Unmarshal] Err=", err)
				}
				CheckServer = append(CheckServer, serverEmp)
			}
		}

		statusOn := 0
		for _, check := range CheckServer {
			if check.Status {
				statusOn++
			}
		}
		rateUptime := fmt.Sprintf("%.2f %s", 100*(float64(statusOn)/float64(len(CheckServer))), "%")
		msg1 = msg1 + fmt.Sprintf("\n '%s' rate uptime: %s ", server.Name, rateUptime)

	}

	msg2 := fmt.Sprintf("Total number of server : %s \nSERVERS ON : %s \nSERVERS OFF : %s \n\n", strconv.Itoa(len(servers)), strconv.Itoa(serverOn), strconv.Itoa(serverOff))
	msg := msg2 + msg1
	m := gomail.NewMessage()
	m.SetHeader("From", mail)
	m.SetHeader("To", "nguyentuanh5527@gmail.com")

	m.SetHeader("Subject", "Report Servers "+time.Now().Format("01-02-2006"))
	m.SetBody("text/plain", msg)

	d := gomail.NewDialer("smtp.gmail.com", 587, mail, "ikvjpolypjwerykg")
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// send
	time.Sleep(time.Second * 10)
	if err := d.DialAndSend(m); err != nil {
		// TODO: this function should return an error: sendEmail(receivers []string) error
		// panic here will make program/service stop, which is an unexpected behavior.
		log.Fatal(err)
	}

	c.LogRorate.Info("......done email........")
}

// UpdateServerPeriodic implements PeriodicController
func (c *periodicController) UpdateServerPeriodic() {
	servers := []entity.Server{}
	c.DB.Find(&servers)
	// c.LogRorate.Info(servers)
	checkServer := checkServer{}
	for _, server := range servers {

		_, err := exec.Command("ping", server.Ipv4).Output()

		if (err != nil) && (server.Status) {
			server.Status = false
			err = c.DB.Where("name = ? ", server.Name).Updates(&server).Error
			if err != nil {
				// ErrorLogger.Println("message : could not update Server " + server.Ipv4)
				c.LogRorate.Error("message : could not update Server " + server.Ipv4)
				continue
			}

			checkServer.Name = server.Name
			checkServer.Status = server.Status
			checkServer.Time = time.Now().Format("02-06-2006")

			dataJSON, _ := json.Marshal(checkServer)
			js := string(dataJSON)
			ind, err := c.ElasticClient.Index().
				Index("server").
				BodyJson(js).
				Do(context.Background())

			if err != nil {
				c.LogRorate.Error(err, ind.Index)
				continue
			}
			c.LogRorate.Info(server.Ipv4 + " has been update ON -> OFF")
			continue
		} else {
			if (err == nil) && (!server.Status) {
				server.Status = true
				err = c.DB.Where("name = ? ", server.Name).Updates(&server).Error
				if err != nil {
					c.LogRorate.Error("message : could not update Server " + server.Ipv4)
					continue
				}

				checkServer.Name = server.Name
				checkServer.Status = server.Status
				checkServer.Time = time.Now().Format("02-06-2006")

				dataJSON, _ := json.Marshal(checkServer)

				js := string(dataJSON)
				ind, err := c.ElasticClient.Index().
					Index("server").
					BodyJson(js).
					Do(context.Background())

				if err != nil {
					c.LogRorate.Error(err, ind.Index)
					continue
				}

				c.LogRorate.Info(server.Ipv4 + " has been update OFF -> ON")
				continue
			}
		}
	}
}
