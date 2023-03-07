package main

// utf8mb4 - emoji表情

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strconv"
	"time"
)

func main() {
	// 连接数据库
	dsn := "root:1234@tcp(127.0.0.1:3306)/go-crud-list?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	fmt.Println(db)
	fmt.Print(err)

	sqlDB, err := db.DB()
	// 连接池中最大连接数
	sqlDB.SetMaxIdleConns(10)
	// 打开数据库连接的最大数量
	sqlDB.SetMaxOpenConns(100)
	// 连接可复用的最大时间
	sqlDB.SetConnMaxLifetime(10 * time.Second)

	// 结构体
	type List struct {
		gorm.Model
		Name    string `gorm:"type:varchar(20); not null" json:"name" binding:"required"`
		Sate    string `gorm:"type:varchar(20); not null" json:"state" binding:"required"`
		Phone   string `gorm:"type:varchar(20); not null" json:"phone" binding:"required"`
		Email   string `gorm:"type:varchar(40); not null" json:"email" binding:"required"`
		Address string `gorm:"type:varchar(20); not null" json:"address" binding:"required"`
	}

	db.AutoMigrate(&List{})

	// 接口
	r := gin.Default()

	// 获取
	r.POST("/user/add", func(c *gin.Context) {
		var data List

		err := c.ShouldBindJSON(&data)

		if err != nil {
			c.JSON(200, gin.H{
				"msg":  "添加失败",
				"data": gin.H{},
				"code": 400,
			})
		} else {
			// 数据库操作
			db.Create(&data) // 创建一条数据

			c.JSON(200, gin.H{
				"mes":  "添加成功",
				"data": data,
				"code": 200,
			})
		}
	})

	// 删除
	r.DELETE("/user/delete/:id", func(c *gin.Context) {
		var data []List

		id := c.Param("id")

		db.Where("id = ?", id).Find(&data)

		if len(data) == 0 {
			c.JSON(200, gin.H{
				"mes":  "ID没有找到，删除失败",
				"code": 400,
			})
		} else {

			db.Where("id = ?", id).Delete(&data)

			c.JSON(200, gin.H{
				"msg":  "删除成功",
				"code": 200,
			})
		}
	})

	// 更新
	r.PUT("/user/update/:id", func(c *gin.Context) {
		var data List
		id := c.Param("id")

		db.Select("id").Where("id = ?", id).Find(&data)

		if data.ID == 0 {
			c.JSON(200, gin.H{
				"msg":  "用户id没有找到",
				"code": 400,
			})
		} else {
			err := c.ShouldBindJSON(&data) // 接收并验证

			if err != nil {
				c.JSON(200, gin.H{
					"msg":  "修改失败",
					"code": 400,
				})
			} else {

				db.Where("id = ?", id).Updates(&data)

				c.JSON(200, gin.H{
					"msg":  "修改成功",
					"code": 200,
				})
			}
		}

	})

	// 查询
	r.GET("/user/list/:name", func(c *gin.Context) {
		var dataList []List
		name := c.Param("name")
		db.Where("name = ?", name).Find(&dataList)
		if len(dataList) == 0 {
			c.JSON(200, gin.H{
				"msg":  "没有查询到数据",
				"code": 400,
				"data": gin.H{},
			})
		} else {
			c.JSON(200, gin.H{
				"msg":  "查询成功",
				"code": 200,
				"data": dataList,
			})
		}
	})

	// 全部查询
	r.GET("/user/list", func(c *gin.Context) {

		var dataList []List

		var total int64

		// 查询分页
		pageSize, _ := strconv.Atoi(c.Query("pageSize"))
		pageNum, _ := strconv.Atoi(c.Query("pageNum"))

		if pageSize == 0 {
			pageSize = -1
		}
		if pageNum == 0 {
			pageSize = -1
		}
		offsetVal := (pageNum - 1) * pageSize
		if pageNum == -1 && pageSize == -1 {
			offsetVal = 1
		}

		db.Model(dataList).Count(&total).Limit(pageSize).Offset(offsetVal).Find(&dataList)

		if len(dataList) == 0 {
			c.JSON(200, gin.H{
				"msg":  "没有查询到数据",
				"code": 400,
				"data": gin.H{},
			})
		} else {
			c.JSON(200, gin.H{
				"msg":  "查询成功",
				"code": 200,
				"data": gin.H{
					"list":     dataList,
					"total":    total,
					"pageNum":  pageNum,
					"pageSize": pageSize,
				},
			})
		}

	})

	// 端口号
	PORT := "3000"
	r.Run(":" + PORT)
}
