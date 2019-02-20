// Package main provides ...
package main

import (
	"fmt"
	//"bytes"
    "strconv"
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

const BASE_SQL string =  "select audio_id, audio_name, audio_url, audio_mp3_url from gcore_audio"
var db, err = sql.Open("mysql", "root@tcp(127.0.0.1:3306)/gcore")

type Audio struct {
	Audio_id      int    `json:"audio_id" form:"audio_id"`
	Audio_name    string `json:"audio_name" form:"audio_name"`
	Audio_url     string `json:"audio_url" form:"audio_url"`
	Audio_mp3_url string `json:"audio_mp3_url" form:"audio_mp3_url"`
}

type Playinfo struct {
    Audio_flow_info string `json:"audio_flow_info" form:"audio_flow_info"`
    Audio_djs string `json:"audio_djs" form:"audio_djs"`
}



func AudioPlayinfoGet(c *gin.Context) {
	var (
		playinfo  Playinfo
		result gin.H
	)
	audio_id := c.Param("audio_id")
	row := db.QueryRow("select audio_flow_info, audio_djs from gcore_audio where audio_id = ?", audio_id)
	err := row.Scan(&playinfo.Audio_flow_info, &playinfo.Audio_djs)
	if err != nil {
		fmt.Println(err.Error())

		result = gin.H{
			"result": nil,
			"count":  0,
		}
	} else {
		result = gin.H{
			"result": playinfo,
			"count":  1,
		}
	}
	c.JSON(http.StatusOK, result)
}

func RecentAudiosGet(c *gin.Context) {
	var (
		audio  Audio
		audios []Audio
	)
	page, err := strconv.Atoi(c.DefaultQuery("page", "0"))
	offset := page * 10
	rows, err := db.Query("select audio_id, audio_name, audio_url, audio_mp3_url from gcore_audio order by audio_date desc limit 10 offset ?", offset)
    if err != nil {
        fmt.Println(err.Error())
    }
	for rows.Next() {
		err := rows.Scan(&audio.Audio_id, &audio.Audio_name, &audio.Audio_url, &audio.Audio_mp3_url)
		audios = append(audios, audio)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	defer rows.Close()
	c.JSON(http.StatusOK, gin.H{
		"result": audios,
		"count":  len(audios),
	})
}

func GaycoreHandler(sql string) gin.HandlerFunc {
    fn := func(c *gin.Context){
	var (
		audio  Audio
		audios []Audio
	)
	page, err := strconv.Atoi(c.DefaultQuery("page", "0"))
	offset := page * 10
	rows, err := db.Query(sql, offset)
    if err != nil {
        fmt.Println(err.Error())
    }
	for rows.Next() {
		err := rows.Scan(&audio.Audio_id, &audio.Audio_name, &audio.Audio_url, &audio.Audio_mp3_url)
		audios = append(audios, audio)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	defer rows.Close()
	c.JSON(http.StatusOK, gin.H{
		"result": audios,
		"count":  len(audios),
	})}
    return gin.HandlerFunc(fn)
}

func GaycoreParamHandler(param, sql, sql_count string) gin.HandlerFunc {
    fn := func(c *gin.Context){
	var (
        count int
		audio  Audio
		audios []Audio
	)
	page, err := strconv.Atoi(c.DefaultQuery("page", "0"))
    if err != nil {
        fmt.Println(err.Error())
    }
	audio_param := c.Param(param)
	offset := page * 10
    row := db.QueryRow(sql_count, audio_param)
    err = row.Scan(&count)
    if err != nil {
        fmt.Println(err.Error())
    }
	rows, err := db.Query(sql, audio_param, offset)
    if err != nil {
        fmt.Println(err.Error())
    }
	for rows.Next() {
		err := rows.Scan(&audio.Audio_id, &audio.Audio_name, &audio.Audio_url, &audio.Audio_mp3_url)
		audios = append(audios, audio)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	defer rows.Close()
	c.JSON(http.StatusOK, gin.H{
		"result": audios,
		"count":  count,
	})}
    return gin.HandlerFunc(fn)
}


func main() {
	if err != nil {
		fmt.Println(err.Error())
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		fmt.Println(err.Error())
	}
    RECENT_SQL := BASE_SQL + " order by audio_date desc limit 10 offset ?"
    HOT_COMMENT_SQL := BASE_SQL + " order by audio_comment desc limit 10 offset ?"
    HOT_LIKE_SQL := BASE_SQL + " order by audio_like desc limit 10 offset ?"
    CATEGORY_SQL := BASE_SQL + " where substring_index(audio_cate, '/', -1) = ? order by audio_date desc limit 10 offset ?"
    CATEGORY_COUNT_SQL := "select count(*) from gcore_audio where audio_cate = ?"
    DJS_SQL := BASE_SQL + " where audio_djs like concat('%', ?, '%') order by audio_like desc limit 10 offset ?"
    DJS_COUNT_SQL := "select count(*) from gcore_audio where audio_djs like concat('%', ?, '%')"
    TOPIC_SQL := BASE_SQL + " where audio_name like concat('%', ?, '%') order by audio_date limit 10 offset ?"
    TOPIC_COUNT_SQL := "select count(*) from gcore_audio where audio_name like concat('%', ?, '%')"

	route := gin.Default()
	route.GET("/audio/:audio_id", AudioPlayinfoGet)
    route.GET("/audios/recent", GaycoreHandler(RECENT_SQL))
    route.GET("audios/hot/comment", GaycoreHandler(HOT_COMMENT_SQL))
    route.GET("audios/hot/like", GaycoreHandler(HOT_LIKE_SQL))
    route.GET("audios/category/:audio_cate", GaycoreParamHandler("audio_cate", CATEGORY_SQL, CATEGORY_COUNT_SQL))
    route.GET("audios/djs/:djs_name", GaycoreParamHandler("djs_name", DJS_SQL, DJS_COUNT_SQL))
    route.GET("audios/topic/:audio_topic", GaycoreParamHandler("audio_topic", TOPIC_SQL, TOPIC_COUNT_SQL))
	route.Run(":3000")
}
