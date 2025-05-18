package controllers

import (
	"errors"
	"fmt"
	"healthcare/controllers/utils"
	"healthcare/global"
	"healthcare/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateInstitution(ctx *gin.Context) {
	var institution models.Institution
	if err := ctx.ShouldBindJSON(&institution); err != nil {
		fmt.Printf("Failed to bind JSON: %v\n", err)
		fmt.Printf("Request body: %v\n", ctx.Request.Body)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "请求格式错误: " + err.Error(),
		})
		return
	}

	// Get user ID from the path
	userID := ctx.Param("id")
	uid, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的用户ID",
		})
		return
	}

	// Verify that the user is an institution user
	var user models.User
	if err := global.DB.First(&user, uid).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "用户不存在",
		})
		return
	}

	if user.UserType != 3 {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "只有机构用户可以创建机构信息",
		})
		return
	}

	// Check if institution already exists for this user
	var existingInstitution models.Institution
	err = global.DB.Where("user_id = ?", uid).First(&existingInstitution).Error
	if err == nil {
		// 找到了已存在的机构记录
		ctx.JSON(http.StatusConflict, gin.H{
			"error": "该用户已经创建过机构信息",
		})
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// 发生了除了"记录未找到"之外的错误
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "数据库错误: " + err.Error(),
		})
		return
	}
	// 如果是 ErrRecordNotFound 错误，说明可以创建新机构，继续执行

	institution.UserID = uint(uid)
	institution.Status = 0 // Set status as pending

	fmt.Printf("Creating institution: %+v\n", institution)

	if err := global.DB.Create(&institution).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "创建机构信息失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":     "机构信息提交成功，等待管理员审核",
		"institution": institution,
	})
}

func GetPendingInstitutions(ctx *gin.Context) {
	// Verify that the user is an admin
	username := ctx.GetString("username")
	var user models.User
	if err := global.DB.Where("username = ?", username).First(&user).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	if user.UserType != 2 {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "Only admin users can view pending institutions",
		})
		return
	}

	var institutions []models.Institution
	if err := global.DB.Where("status = ?", 0).Find(&institutions).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, institutions)
}

func ReviewInstitution(ctx *gin.Context) {
	// Verify that the user is an admin
	username := ctx.GetString("username")
	var user models.User
	if err := global.DB.Where("username = ?", username).First(&user).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	if user.UserType != 2 {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "Only admin users can review institutions",
		})
		return
	}

	institutionID := ctx.Param("id")
	var input struct {
		Approved bool `json:"approved"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var institution models.Institution
	if err := global.DB.First(&institution, institutionID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Institution not found",
		})
		return
	}

	institution.Status = map[bool]uint8{true: 1, false: 2}[input.Approved]
	if err := global.DB.Save(&institution).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": map[bool]string{
			true:  "Institution approved successfully",
			false: "Institution rejected",
		}[input.Approved],
	})
}

func GetInstitutions(ctx *gin.Context) {
	var institutions models.Institution
	if err := global.DB.Where("status = ?", 1).Find(&institutions).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, institutions)
}

func GetInstitutionDetail(ctx *gin.Context) {
	institutionID := ctx.Param("id")
	var institution models.Institution

	if err := global.DB.First(&institution, institutionID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Institution not found",
		})
		return
	}

	// 检查请求是否来自管理员
	username := ctx.GetString("username")
	var user models.User
	if err := global.DB.Where("username = ?", username).First(&user).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	// 如果不是管理员，只能看到已批准的机构
	if user.UserType != 2 && institution.Status != 1 {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "You don't have permission to view this institution",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"institution": institution,
		"isAdmin":     user.UserType == 2,
	})
}

// CreateInstitutionPlans 创建机构套餐&体检项目或者创建目标套餐的新增体检项目
func CreateInstitutionPlans(ctx *gin.Context) {
	institutionID := ctx.Param("id")
	planID := ctx.Param("plan_id")
	var institution models.Institution

	if err := global.DB.First(&institution, institutionID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Institution not found",
		})
		return
	}
	// 验证操作者权限（仅机构所有者或管理员可以更新套餐）
	username := ctx.GetString("username")
	var user models.User
	if err := global.DB.Where("username = ?", username).First(&user).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	if user.ID != institution.UserID && user.UserType != 2 {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "You don't have permission to create plans for this institution",
		})
		return
	}

	var input struct {
		PlanName        *string  `json:"plan_name"` //允许没有套餐名称
		HealthItem      string   `json:"health_item"`
		ItemDescription *string  `json:"item_description"` //允许没有描述
		PlanPrice       *float64 `json:"plan_price"`       // 套餐价格
		Description     *string  `json:"description"`      // 套餐描述
		SuitableFor     *string  `json:"suitable_for"`     // 适用人群
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	fmt.Printf("Received input: %+v\n", input)

	var newPlan models.Plan
	newPlan.RelationInstitutionID = institution.ID
	if planID == "" {
		newPlan.PlanName = *input.PlanName

		// 设置可选的套餐字段
		if input.PlanPrice != nil {
			newPlan.PlanPrice = *input.PlanPrice
		}
		if input.Description != nil {
			newPlan.Description = *input.Description
		}
		if input.SuitableFor != nil {
			newPlan.SuitableFor = *input.SuitableFor
		}

		// 检查套餐名称是否已存在
		exists, err := utils.CheckExists(&models.Plan{}, "plan_name", *input.PlanName)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "数据库错误: " + err.Error(),
			})
			return
		}
		if exists {
			ctx.JSON(http.StatusConflict, gin.H{
				"error": fmt.Sprintf("%v 已经存在，请更换名称或者是更新已有内容", input.PlanName)})
			return
		}

		if err := global.DB.Create(&newPlan).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	} else {
		// 检查planid对应的套餐是否存在
		if err := global.DB.First(&newPlan, planID).Error; err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "套餐不存在",
			})
			return
		}
	}

	var newHealthItem models.HealthItem
	newHealthItem.ItemName = input.HealthItem
	// 检查指标名称是否已存在
	exists, err := utils.CheckExists(&models.HealthItem{}, "item_name", input.HealthItem)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "数据库错误: " + err.Error(),
		})
		return
	}
	if !exists {
		if err := global.DB.Create(&newHealthItem).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	} else {
		// 已存在则查出ID
		if err := global.DB.Where("item_name = ?", input.HealthItem).First(&newHealthItem).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "查询已存在指标失败: " + err.Error(),
			})
			return
		}
	}

	var newPlanHealthItem models.PlanHeathItem
	newPlanHealthItem.RelationPlanId = newPlan.ID
	newPlanHealthItem.RelationHealthItemId = newHealthItem.ID
	if input.ItemDescription != nil {
		newPlanHealthItem.ItemDescription = *input.ItemDescription
	}

	if err := global.DB.Create(&newPlanHealthItem).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": map[bool]string{
			true:  "套餐创建成功",
			false: "体检项目创建成功"}[planID == ""],

		"plan":      newPlan,
		"item":      newHealthItem,
		"plan_item": newPlanHealthItem,
	})
}

func GetInstitutionPlans(ctx *gin.Context) {
	institutionID := ctx.Param("id")
	var institution models.Institution

	if err := global.DB.First(&institution, institutionID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Institution not found",
		})
		return
	}

	// 获取用户身份
	username := ctx.GetString("username")
	var user models.User
	if err := global.DB.Where("username = ?", username).First(&user).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	// 如果机构未批准，只有机构所有者和管理员才能查看套餐
	// 但如果机构已批准，任何用户都可以查看套餐
	if institution.Status != 1 && user.ID != institution.UserID && user.UserType != 2 {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "This institution is not approved",
		})
		return
	}

	// 获取机构所属下的套餐名称（可能有多个套餐）
	var plans []models.Plan
	if err := global.DB.Where("institution_id = ?", institution.ID).Find(&plans).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 获取套餐对应的id切片
	planIDs := make([]uint, len(plans))
	for i, p := range plans {
		planIDs[i] = p.ID
	}

	// 获取套餐对应的指标信息
	var planItems []models.PlanHeathItem
	if err := global.DB.Where("plan_id IN ?", planIDs).Preload("ThisHeathItem").Find(&planItems).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 组织套餐详细信息以满足前端需要
	type PlanDetail struct {
		ID          uint    `json:"id"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		SuitableFor string  `json:"suitable_for"`
		Items       string  `json:"items"`
		Price       float64 `json:"price"`
	}

	planDetails := make([]PlanDetail, 0, len(plans))

	// 为每个套餐收集指标项
	planItemsMap := make(map[uint][]models.PlanHeathItem)
	for _, item := range planItems {
		if _, exists := planItemsMap[item.RelationPlanId]; !exists {
			planItemsMap[item.RelationPlanId] = make([]models.PlanHeathItem, 0)
		}
		planItemsMap[item.RelationPlanId] = append(planItemsMap[item.RelationPlanId], item)
	}

	// 构建每个套餐的详细信息
	for _, plan := range plans {
		items := planItemsMap[plan.ID]
		itemNames := make([]string, 0, len(items))

		// 收集所有指标项的名称
		for _, item := range items {
			if item.ThisHeathItem.ItemName != "" {
				itemNames = append(itemNames, item.ThisHeathItem.ItemName)
			}
		}

		// 如果没有描述，使用第一个项目的描述作为套餐描述
		description := plan.Description
		if description == "" && len(items) > 0 {
			description = items[0].ItemDescription
		}

		planDetail := PlanDetail{
			ID:          plan.ID,
			Name:        plan.PlanName,
			Description: description,
			SuitableFor: plan.SuitableFor,
			Items:       strings.Join(itemNames, ", "),
			Price:       plan.PlanPrice,
		}

		fmt.Printf("Plan Detail: %+v\n", planDetail)
		planDetails = append(planDetails, planDetail)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"institution": institution,
		"plans":       planDetails,
		"items":       planItems,
	})
}

func UpdateInstitutionPlanorItem(ctx *gin.Context) {
	institutionID := ctx.Param("id")

	var input struct {
		PlanID          uint     `json:"plan_id"`
		ItemID          *uint    `json:"item_id"`
		ItemName        *string  `json:"item_name"`
		ItemDescription *string  `json:"item_description"`
		PlanName        *string  `json:"plan_name"`
		PlanPrice       *float64 `json:"plan_price"`
		PlanDescription *string  `json:"description"`
		PlanSuitableFor *string  `json:"suitable_for"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var institution models.Institution
	if err := global.DB.First(&institution, institutionID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Institution not found",
		})
		return
	}
	if institution.Status != 1 {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "This institution is not approved",
		})
		return
	}

	// Update plan name and other fields
	if input.PlanName != nil && *input.PlanName != "" {
		utils.UpdateIt(&models.Plan{}, input.PlanID, "plan_name", *input.PlanName)
	}

	// Update plan price
	if input.PlanPrice != nil {
		utils.UpdateIt(&models.Plan{}, input.PlanID, "plan_price", *input.PlanPrice)
	}

	// Update plan description
	if input.PlanDescription != nil && *input.PlanDescription != "" {
		utils.UpdateIt(&models.Plan{}, input.PlanID, "description", *input.PlanDescription)
	}

	// Update plan suitable for
	if input.PlanSuitableFor != nil && *input.PlanSuitableFor != "" {
		utils.UpdateIt(&models.Plan{}, input.PlanID, "suitable_for", *input.PlanSuitableFor)
	}
	if input.ItemID != nil {
		// Update item name
		if input.ItemName != nil && *input.ItemName != "" {
			utils.UpdateIt(&models.HealthItem{}, input.ItemID, "item_name", *input.ItemName)
		}

		// Update item description
		if input.ItemDescription != nil {
			utils.UpdateIt(&models.PlanHeathItem{}, input.ItemID, "item_description", *input.ItemDescription)
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "更新成功"})

}

// 删除套餐内一个体检项目
func DeleteInsistutionPlanonItem(ctx *gin.Context) {
	var input struct {
		PlanID uint `json:"plan_id"`
		ItemID uint `json:"item_id"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	fmt.Println("input:", input)

	// 根据planid，itemid查找删除
	var planheathitem models.PlanHeathItem
	if err := global.DB.Where("plan_id = ? AND item_id = ?", input.PlanID, input.ItemID).First(&planheathitem).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "套餐内体检项目不存在",
		})
		return
	}

	// 删除套餐内体检项目
	if err := global.DB.Unscoped().Delete(&planheathitem).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "删除失败: " + err.Error(),
		})
		return
	}

	var remainitem []models.PlanHeathItem
	if err := global.DB.Where("item_id = ?", input.ItemID).Find(&remainitem).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "数据库错误" + err.Error(),
		})
		return
	}

	// 如果没有其他套餐引用该体检项目，则删除该体检项目
	if len(remainitem) == 0 {
		var healthitem models.HealthItem
		if err := global.DB.Where("id = ?", input.ItemID).First(&healthitem).Error; err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "体检项目不存在",
			})
			return
		}
		// 删除体检项目
		if err := global.DB.Unscoped().Delete(&healthitem).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "删除失败: " + err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"message": "体检项目删除成功",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "套餐内体检项目删除成功",
	})

}

// 删除套餐
func DeleteInstitutionPlan(ctx *gin.Context) {
	var input struct {
		PlanID uint `json:"plan_id"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	fmt.Println("input:", input)

	// 删除套餐plan
	var plan models.Plan
	if err := global.DB.Where("id = ?", input.PlanID).First(&plan).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "套餐不存在",
		})
		return
	}

	// 删除套餐相关评论
	if err := global.DB.Where("plan_id = ?", input.PlanID).Unscoped().Delete(&models.Commentary{}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "删除套餐评论失败: " + err.Error(),
		})
		return
	}

	// 删除套餐相关体检项目
	var planheathitems []models.PlanHeathItem
	if err := global.DB.Where("plan_id = ?", input.PlanID).Find(&planheathitems).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "套餐内体检项目不存在",
		})
		return
	}

	for _, item := range planheathitems {
		if err := global.DB.Where("id = ?", item.ID).Unscoped().Delete(&models.PlanHeathItem{}).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "删除套餐内体检项目失败: " + err.Error(),
			})
			return
		}

		// 查找相同体检项目数目
		var count int64
		if err := global.DB.Model(&models.PlanHeathItem{}).Where("item_id = ?", item.RelationHealthItemId).Count(&count).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"errors": "统计删除套餐内体相同检项目失败" + err.Error(),
			})
		}

		if count == 0 {
			if err := global.DB.Where("id = ?", item.RelationHealthItemId).Unscoped().Delete(&models.HealthItem{}).Error; err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": "删除套餐内体检项目失败: " + err.Error(),
				})
				return
			}
		}
	}

	if err := global.DB.Unscoped().Delete(&plan).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "删除失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "套餐删除成功",
	})

}
