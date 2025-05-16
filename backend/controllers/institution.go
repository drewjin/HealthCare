package controllers

import (
	"errors"
	"fmt"
	"healthcare/controllers/utils"
	"healthcare/global"
	"healthcare/models"
	"net/http"
	"strconv"

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
		PlanName        *string `json:"plan_name"` //允许没有套餐名称
		HealthItem      string  `json:"health_item"`
		ItemDescription *string `json:"item_description"` //允许没有描述
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var newPlan models.Plan
	newPlan.RelationInstitutionID = institution.ID
	if planID == "" {
		newPlan.PlanName = *input.PlanName

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

	// 只允许查看已批准的机构的套餐
	if institution.Status != 1 {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "This institution is not approved",
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
	if err := global.DB.Where("plan_id IN ?", planIDs).Find(&planItems).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"institution": institution,
		"plans":       plans,
		"items":       planItems,
	})
}

func UpdateInsistutionPlanorItem(ctx *gin.Context) {
	institutionID := ctx.Param("id")

	var input struct {
		PlanID                   uint    `json:"plan_id"`
		ItemID                   uint    `json:"item_id"`
		ItemName                 *string `json:"item_name"`
		ItemDescription          *string `json:"item_description"`
		PlanName                 *string `json:"plan_name"`
		InstitutionName          *string `json:"institution_name"`
		InstitutionPhone         *string `json:"institution_phone"`
		InstitutionAddress       *string `json:"institution_address"`
		InstitutionQualification *string `json:"institution_qualification"`
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

	// Update plan name
	if input.PlanName != nil {
		utils.UpdateIt(&models.Plan{}, input.PlanID, "plan_name", *input.PlanName)
	}

	// Update item name
	if input.ItemName != nil {
		utils.UpdateIt(&models.HealthItem{}, input.ItemID, "item_name", *input.ItemName)
	}

	// Update item description
	if input.ItemDescription != nil {
		utils.UpdateIt(&models.PlanHeathItem{}, input.ItemID, "item_description", *input.ItemDescription)
	}

	// Update institution Name
	if input.InstitutionName != nil {
		utils.UpdateIt(&models.Institution{}, institution.ID, "institution_name", *input.InstitutionName)
	}

	// Update institution Phone
	if input.InstitutionPhone != nil {
		utils.UpdateIt(&models.Institution{}, institution.ID, "institution_phone", *input.InstitutionPhone)
	}

	// Update institution Address
	if input.InstitutionAddress != nil {
		utils.UpdateIt(&models.Institution{}, institution.ID, "institution_address", *input.InstitutionAddress)
	}

	// Update institution Qualification
	if input.InstitutionQualification != nil {
		utils.UpdateIt(&models.Institution{}, institution.ID, "institution_qualification", *input.InstitutionQualification)
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
func DeleteInsistutionPlan(ctx *gin.Context) {
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

// 删除机构
func DeleteInsistution(ctx *gin.Context) {}
