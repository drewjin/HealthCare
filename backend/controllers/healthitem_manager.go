package controllers

import (
	"fmt"
	"healthcare/controllers/utils"
	"healthcare/global"
	"healthcare/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// CreateHealthItemTemplate 创建检查项目模板
// 输入格式：项目1, 项目2, 项目3
// 存储格式：项目1:, 项目2:, 项目3:
func CreateHealthItemTemplate(ctx *gin.Context) {
	var input struct {
		HealthItems string `json:"health_items" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "输入格式不正确",
		})
		return
	}

	// 分割输入的项目
	items := strings.Split(input.HealthItems, ",")
	if len(items) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "至少需要输入一个检查项目",
		})
		return
	}

	// 构建检查项目模板字符串
	var templateItems []string
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item != "" {
			templateItems = append(templateItems, item+":")
		}
	}

	// 合并成最终的模板字符串
	templateString := strings.Join(templateItems, ", ")

	// 返回构建好的模板
	ctx.JSON(http.StatusOK, gin.H{
		"template": templateString,
	})
}

// UpdateHealthItemValues 更新健康项目值
// 输入格式：当前字符串和需要更新的键值对
// 返回格式："项目1:数值1, 项目2:数值2, 项目3:数值3"
func UpdateHealthItemValues(ctx *gin.Context) {
	var input struct {
		ItemID     uint              `json:"item_id" binding:"required"`
		ItemString string            `json:"item_string" binding:"required"` // 当前的健康项目字符串
		Updates    map[string]string `json:"updates" binding:"required"`     // 要更新的值
		DeleteKeys []string          `json:"delete_keys"`                    // 要删除的键
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "输入格式不正确",
		})
		return
	}

	// 使用工具函数解析当前项目字符串为键值对
	itemMap := utils.ParseHealthItemString(input.ItemString)

	// 更新键值对
	for key, value := range input.Updates {
		key = strings.TrimSpace(key)
		if key != "" {
			itemMap[key] = value
		}
	}

	// 删除指定的键
	for _, key := range input.DeleteKeys {
		key = strings.TrimSpace(key)
		if key != "" {
			delete(itemMap, key)
		}
	}

	// 重新构建项目字符串
	newItemString := utils.BuildHealthItemString(itemMap)

	// 查找HealthItem记录
	var healthItem models.HealthItem
	if err := global.DB.First(&healthItem, input.ItemID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "检查项目不存在",
		})
		return
	}

	// 更新健康项目模板字符串
	if err := global.DB.Model(&healthItem).Update("item_name", newItemString).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "更新检查项目失败: " + err.Error(),
		})
		return
	}

	// 将健康项目字符串转换为JSON并更新ItemMetrics字段
	jsonString, err := utils.ConvertToJSON(newItemString)
	if err == nil {
		// 查找相关的PlanHeathItem记录并更新ItemMetrics
		var planItems []models.PlanHeathItem
		if err := global.DB.Where("health_item_id = ?", input.ItemID).Find(&planItems).Error; err == nil {
			for _, planItem := range planItems {
				global.DB.Model(&planItem).Update("item_metrics", jsonString)
			}
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":     "检查项目更新成功",
		"item_string": newItemString,
	})
}

// GetHealthItemValues 获取健康项目值的解析结果
// 输入：健康项目ID
// 返回：解析后的键值对
func GetHealthItemValues(ctx *gin.Context) {
	itemID := ctx.Param("id")
	
	// 调试输出
	fmt.Printf("GetHealthItemValues - 接收到的项目ID参数: '%s'\n", itemID)
	
	// 检查ID是否是有效的数字
	if itemID == "undefined" || itemID == "" {
		fmt.Println("GetHealthItemValues - 无效的项目ID (undefined 或空)")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的检查项目ID",
		})
		return
	}

	// 尝试将ID转换为数值进行验证
	_, err := strconv.Atoi(itemID)
	if err != nil {
		fmt.Printf("GetHealthItemValues - 项目ID不是有效数字: %v\n", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的检查项目ID格式",
		})
		return
	}

	var healthItem models.HealthItem
	if err := global.DB.First(&healthItem, itemID).Error; err != nil {
		fmt.Printf("GetHealthItemValues - 数据库查询失败: %v\n", err)
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "检查项目不存在",
		})
		return
	}

	// 使用工具函数解析健康项目字符串为键值对
	itemMap := utils.ParseHealthItemString(healthItem.ItemName)

	ctx.JSON(http.StatusOK, gin.H{
		"item_id":     healthItem.ID,
		"item_string": healthItem.ItemName,
		"values":      itemMap,
	})
}

// SaveHealthItemTemplate 保存健康项目模板到数据库
func SaveHealthItemTemplate(ctx *gin.Context) {
	var input struct {
		ItemName string `json:"item_name" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "输入格式不正确",
		})
		return
	}

	// 创建新的健康项目记录
	healthItem := models.HealthItem{
		ItemName: input.ItemName,
	}

	if err := global.DB.Create(&healthItem).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "创建健康项目失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "健康项目保存成功",
		"item": healthItem,
	})
}