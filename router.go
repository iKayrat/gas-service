package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/iKayrat/gas-service/gas"
)

func Handlers() *gin.Engine {
	router := gin.Default()

	router.GET("/transactions", GetTransactions)
	router.GET("/monthly", Monthly)
	router.GET("/daily", Daily)
	router.GET("/wholeperiod", WholePeriod)
	router.GET("/hourly", Hourly)

	return router
}

func GetTransactions(ctx *gin.Context) {
	Req := new(gas.Body)

	err := ctx.ShouldBindWith(&Req, binding.JSON)
	if err != nil {
		log.Printf("Binding json err: %v", err)
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(200, Req)
}

func Monthly(ctx *gin.Context) {
	res, err := http.Get("https://raw.githubusercontent.com/CryptoRStar/GasPriceTestTask/main/gas_price.json")
	if err != nil {
		log.Println("getfrom github err:", err)
		return
	}
	defer res.Body.Close()

	Req := new(gas.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("readAll err:", err)
		return
	}

	err = json.Unmarshal(body, &Req)
	if err != nil {
		log.Println("unmarshal err:", err)
		return
	}

	perMonth, err := Req.SpentPerMonth()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": err,
		})
		return
	}

	ctx.JSON(200, perMonth)
}

func Daily(ctx *gin.Context) {
	res, err := http.Get("https://raw.githubusercontent.com/CryptoRStar/GasPriceTestTask/main/gas_price.json")
	if err != nil {
		log.Println("getfrom github err:", err)
		return
	}
	defer res.Body.Close()

	Req := new(gas.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("readAll err:", err)
		return
	}

	err = json.Unmarshal(body, &Req)
	if err != nil {
		log.Println("unmarshal err:", err)
		return
	}

	perDay, err := Req.AveragePerDay()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": err,
		})
		return
	}

	ctx.JSON(200, perDay)

}

func Hourly(ctx *gin.Context) {
	res, err := http.Get("https://raw.githubusercontent.com/CryptoRStar/GasPriceTestTask/main/gas_price.json")
	if err != nil {
		log.Println("getfrom github err:", err)
		return
	}
	defer res.Body.Close()

	Req := gas.Body{}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("readAll err:", err)
		return
	}

	err = json.Unmarshal(body, &Req)
	if err != nil {
		log.Println("unmarshal err:", err)
		return
	}
	if err := ctx.ShouldBindJSON(&Req); err != nil {
		log.Println(err)
		ctx.JSON(400, gin.H{
			"err": err,
		})
		return
	}

	ctx.JSON(200, Req.PerHour())
}

func WholePeriod(ctx *gin.Context) {
	now := time.Now()
	res, err := http.Get("https://raw.githubusercontent.com/CryptoRStar/GasPriceTestTask/main/gas_price.json")
	if err != nil {
		log.Println("getfrom github err:", err)
		return
	}
	defer res.Body.Close()

	Req := gas.Body{}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("readAll err:", err)
		return
	}

	err = json.Unmarshal(body, &Req)
	if err != nil {
		log.Println("unmarshal err:", err)
		return
	}

	ctx.JSON(200, Req.WholePeriod())
	end := time.Now()
	log.Println("time per request:", end.Sub(now))
}
