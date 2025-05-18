# Institution API修改

-**统一路由组：/api/institution**
-**所有数据都用""代替**

## POST /:id | 创建机构

- **User**需要*机构方*权限:3级
- 路由样式：
  [127.0.0.1:xxxx/api/institutions/:id](https://127.0.0.1:xxxx/api/institutions/:id)
    (*:id*应为当前用户的User_id)
- 请求体：
    **最好都不是空值**

  ```json
  {
  "institution_name": "",
  "institution_address": "",
  "institution_phone":"",
  "institution_qualification": ""
  }
  ```

-返回信息：
*成功*

```json
    {
    "institution": {
        "ID": "",
        "CreatedAt": "",
        "UpdatedAt": "",
        "DeletedAt": ,
        "institution_name": "",
        "institution_address": "",
        "institution_phone": "",
        "institution_qualification": "",
        "user_id": "",
        "status": "",
        "user": {
            "ID": "",
            "CreatedAt": "",
            "UpdatedAt": "",
            "DeletedAt": "",
            "Username": "",
            "Password": "",
            "Name": "",
            "Gender": "",
            "Birthday": "",
            "Phone": "",
            "Email": "",
            "Address": "",
            "user_type": ""
        }
    },
    "message": "机构信息提交成功，等待管理员审核"
}
```

*失败*
请查看具体错误error信息，比对纠错

---

## GET /pending | 查看机构创建请求

- **User**需要*管理者*权限:2级
- 路由样式：[127.0.0.1:xxxx/api/institutions/pending](http://127.0.0.1:xxxx/api/institutions/pending)
- 请求体：无
- 返回信息：
  *成功*：返回创建成功的institution

  ```json
    {
        "ID": "",
        "CreatedAt": "",
        "UpdatedAt": "",
        "DeletedAt": "",
        "institution_name": "",
        "institution_address": "",
        "institution_phone": "",
        "institution_qualification": "",
        "user_id": "",
        "status": "",
        "user": {
            "ID": "",
            "CreatedAt": "",
            "UpdatedAt": "",
            "DeletedAt": "",
            "Username": "",
            "Password": "",
            "Name": "",
            "Gender": "",
            "Birthday": "",
            "Phone": "",
            "Email": "",
            "Address": "",
            "user_type": ""
        }
    }
  ```

    *失败*
请查看具体错误error信息，比对纠错

---

## POST /:id/review | 审核机构创建请求

- **User**需要*管理者*权限:2级
- 路由样式：[127.0.0.1:xxxx/api/institutions/:id/review](http://127.0.0.1:xxxx/api/institutions/:id/review)
  (*:id*为待审核机构的ID)
- 请求体：

    ```json
    {
        "approved":true/false //true为通过，false为拒绝申请
    }
    ```

- 返回信息：
  *成功*

    ```json
    {
        "message": "Institution approved successfully"
    }
    ```

  *失败*

    ```json
    {
        "message": "Institution rejected"
    }
    ```

  请查看具体错误error信息，比对纠错

---

## GET / | User查看所有管理机构

- **User**需要*机构方*或*管理者*权限:3级或2级
- 路由样式：[127.0.0.1:xxxx/api/institutions](http://127.0.0.1:xxxx/api/institutions)
- 请求体：无
- 返回信息：
  *成功*：返回所有管理机构的列表

    ```json
    [
      {
        "ID": "",
        "CreatedAt": "",
        "UpdatedAt": "",
        "DeletedAt": "",
        "institution_name": "",
        "institution_address": "",
        "institution_phone": "",
        "institution_qualification": "",
        "user_id": "",
        "status": ""
      }
    ]
    ```

  *失败*
  请查看具体错误error信息，比对纠错

---

## GET /:id | User查看指定机构

- **User**需要*机构方*或*管理者*权限:3级或2级
- 路由样式：[127.0.0.1:xxxx/api/institutions/:id](http://127.0.0.1:xxxx/api/institutions/:id)
  (*:id*为指定机构的ID)
- 请求体：无
- 返回信息：
  *成功*

    ```json
    {
        "institution": {
            "ID": "",
            "CreatedAt": "",
            "UpdatedAt": "",
            "DeletedAt": "",
            "institution_name": "",
            "institution_address": "",
            "institution_phone": "",
            "institution_qualification": "",
            "user_id": "",
            "status": "",
            "user": {
                "ID": "",
                "CreatedAt": "",
                "UpdatedAt": "",
                "DeletedAt": "",
                "Username": "",
                "Password": "",
                "Name": "",
                "Gender": "",
                "Birthday": "",
                "Phone": "",
                "Email": "",
                "Address": "",
                "user_type": ""
            }
        },
        "isAdmin": true/false //根据用户是否拥有管理权限决定
    }
    ```

  *失败*
  请查看具体错误error信息，比对纠错

---

## POST /:id/plans | 为指定机构创建套餐(创建时必须有一个体检项目)

- **User**需要*机构方*:3级,或*管理者*权限:2级
- 路由样式：[127.0.0.1:xxxx/api/institutions/:id/plans](http://127.0.0.1:xxxx/api/institutions/:id/plans)
  (*:id*为机构ID)
-**特别注意**
该API与*POST /:id/:plan_id/item*使用同一个控制器实现
- 请求体：
  
    ```json
    {
      "plan_name": "", // 不可为空
      "health_item": "", // 不可为空
      "item_description": "" // 可以为空
    }
    ```

- 返回信息：
  *成功*

    ```json
    {
        "item": {
            "ID": "",
            "CreatedAt": "",
            "UpdatedAt": "",
            "DeletedAt": "",
            "ItemName": ""
        },
        "message": "套餐创建成功",
        "plan": {
            "ID": "",
            "CreatedAt": "",
            "UpdatedAt": "",
            "DeletedAt": "",
            "PlanName": "",
            "RelationInstitutionID": "",
            "ThisInstitution": {
                "ID": "",
                "CreatedAt": "",
                "UpdatedAt": "",
                "DeletedAt": "",
                "institution_name": "",
                "institution_address": "",
                "institution_phone": "",
                "institution_qualification": "",
                "user_id": "",
                "status": "",
                "user": {
                    "ID": "",
                    "CreatedAt": "",
                    "UpdatedAt": "",
                    "DeletedAt": "",
                    "Username": "",
                    "Password": "",
                    "Name": "",
                    "Gender": "",
                    "Birthday": "",
                    "Phone": "",
                    "Email": "",
                    "Address": "",
                    "user_type": ""
                }
            }
        },
        "plan_item": {
            "ID": "",
            "CreatedAt": "",
            "UpdatedAt": "",
            "DeletedAt": "",
            "RelationPlanId": "",
            "RelationHealthItemId": "",
            "ItemDescription": "",
            "ThisPlan": {
                "ID": "",
                "CreatedAt": "",
                "UpdatedAt": "",
                "DeletedAt": "",
                "PlanName": "",
                "RelationInstitutionID": "",
                "ThisInstitution": {
                    "ID": "",
                    "CreatedAt": "",
                    "UpdatedAt": "",
                    "DeletedAt": "",
                    "institution_name": "",
                    "institution_address": "",
                    "institution_phone": "",
                    "institution_qualification": "",
                    "user_id": "",
                    "status": "",
                    "user": {
                        "ID": "",
                        "CreatedAt": "",
                        "UpdatedAt": "",
                        "DeletedAt": "",
                        "Username": "",
                        "Password": "",
                        "Name": "",
                        "Gender": "",
                        "Birthday": "",
                        "Phone": "",
                        "Email": "",
                        "Address": "",
                        "user_type": ""
                    }
                }
            },
            "ThisHeathItem": {
                "ID": "",
                "CreatedAt": "",
                "UpdatedAt": "",
                "DeletedAt": "",
                "ItemName": ""
            }
        }
    }
    ```

  *失败*
  请查看具体错误error信息，比对纠错

---

## GET /:id/plans | 查看指定机构的所有套餐信息

- **User**无权限要求
- 路由样式：[127.0.0.1:xxxx/api/institutions/:id/plans](http://127.0.0.1:xxxx/api/institutions/:id/plans)
  (*:id*为机构ID)
- 请求体：无
- 返回信息：
  *成功*

    ```json
    {
       "" "institution": {
            "ID": "",
            "CreatedAt": "",
            "UpdatedAt": "",
            "DeletedAt": "",
            "institution_name": "",
            "institution_address": "",
            "institution_phone": "",
            "institution_qualification": "",
            "user_id": "",
            "status": "",
            "user": {
                "ID": "",
                "CreatedAt": "",
                "UpdatedAt": "",
                "DeletedAt": "",
                "Username": "",
                "Password": "",
                "Name": "",
                "Gender": "",
                "Birthday": "",
                "Phone": "",
                "Email": "",
                "Address": "",
                "user_type": ""
            }
        },
        "items": [
            {
                "ID": "",
                "CreatedAt": "",
                "UpdatedAt": "",
                "DeletedAt": "",
                "RelationPlanId": "",
                "RelationHealthItemId": "",
                "ItemDescription": "",
                "ThisPlan": {
                    "ID": "",
                    "CreatedAt": "",
                    "UpdatedAt": "",
                    "DeletedAt": "",
                    "PlanName": "",
                    "RelationInstitutionID": "",
                    "ThisInstitution": {
                        "ID": "",
                        "CreatedAt": "",
                        "UpdatedAt": "",
                        "DeletedAt": "",
                        "institution_name": "",
                        "institution_address": "",
                        "institution_phone": "",
                        "institution_qualification": "",
                        "user_id": "",
                        "status": "",
                        "user": {
                            "ID": "",
                            "CreatedAt": "",
                            "UpdatedAt": "",
                            "DeletedAt": "",
                            "Username": "",
                            "Password": "",
                            "Name": "",
                            "Gender": "",
                            "Birthday": "",
                            "Phone": "",
                            "Email": "",
                            "Address": "",
                            "user_type": ""
                        }
                    }
                },
                "ThisHeathItem": {
                    "ID": "",
                    "CreatedAt": "",
                    "UpdatedAt": "",
                    "DeletedAt": "",
                    "ItemName": ""
                }
            }
        ],
        "plans": [
            {
                "ID": "",
                "CreatedAt": "",
                "UpdatedAt": "",
                "DeletedAt": "",
                "PlanName": "",
                "RelationInstitutionID":"",
                "ThisInstitution": {
                    "ID": "",
                    "CreatedAt": "",
                    "UpdatedAt": "",
                    "DeletedAt": "",
                    "institution_name": "",
                    "institution_address": "",
                    "institution_phone": "",
                    "institution_qualification": "",
                    "user_id":"",
                    "status": "",
                    "user": {
                        "ID": "",
                        "CreatedAt": "",
                        "UpdatedAt": "",
                        "DeletedAt": "",
                        "Username": "",
                        "Password": "",
                        "Name": "",
                        "Gender": "",
                        "Birthday": "",
                        "Phone": "",
                        "Email": "",
                        "Address": "",
                        "user_type": ""
                    }
                }
            }
        ]
    }
    ```

  *失败*
  请查看具体错误error信息，比对纠错

---

## POST /:id/:plan_id/item | 为指定机构的指定套餐添加体检项目

- **User**需要*机构方*权限:3级，或*管理员*:2级
- 路由样式：[127.0.0.1:xxxx/api/institutions/:id/:plan_id/item](http://127.0.0.1:xxxx/api/institutions/:id/:plan_id/item)
  (*:id*为机构ID, *:plan_id*为套餐ID)
-**特别注意**
该API与*POST /:id/plans*使用同一个控制器实现
- 请求体：
  
    ```json
    {
      "health_item": "",
      "item_description": ""
    }
    ```

- 返回信息：
  *成功*

    ```json
    {
        "item": {
            "ID": "",
            "CreatedAt": "",
            "UpdatedAt": "",
            "DeletedAt": "",
            "ItemName": ""
        },
        "message": "套餐创建成功",
        "plan": {
            "ID": "",
            "CreatedAt": "",
            "UpdatedAt": "",
            "DeletedAt": "",
            "PlanName": "",
            "RelationInstitutionID": "",
            "ThisInstitution": {
                "ID": "",
                "CreatedAt": "",
                "UpdatedAt": "",
                "DeletedAt": "",
                "institution_name": "",
                "institution_address": "",
                "institution_phone": "",
                "institution_qualification": "",
                "user_id": "",
                "status": "",
                "user": {
                    "ID": "",
                    "CreatedAt": "",
                    "UpdatedAt": "",
                    "DeletedAt": "",
                    "Username": "",
                    "Password": "",
                    "Name": "",
                    "Gender": "",
                    "Birthday": "",
                    "Phone": "",
                    "Email": "",
                    "Address": "",
                    "user_type": ""
                }
            }
        },
        "plan_item": {
            "ID": "",
            "CreatedAt": "",
            "UpdatedAt": "",
            "DeletedAt": "",
            "RelationPlanId": "",
            "RelationHealthItemId": "",
            "ItemDescription": "",
            "ThisPlan": {
                "ID": "",
                "CreatedAt": "",
                "UpdatedAt": "",
                "DeletedAt": "",
                "PlanName": "",
                "RelationInstitutionID": "",
                "ThisInstitution": {
                    "ID": "",
                    "CreatedAt": "",
                    "UpdatedAt": "",
                    "DeletedAt": "",
                    "institution_name": "",
                    "institution_address": "",
                    "institution_phone": "",
                    "institution_qualification": "",
                    "user_id": "",
                    "status": "",
                    "user": {
                        "ID": "",
                        "CreatedAt": "",
                        "UpdatedAt": "",
                        "DeletedAt": "",
                        "Username": "",
                        "Password": "",
                        "Name": "",
                        "Gender": "",
                        "Birthday": "",
                        "Phone": "",
                        "Email": "",
                        "Address": "",
                        "user_type": ""
                    }
                }
            },
            "ThisHeathItem": {
                "ID": "",
                "CreatedAt": "",
                "UpdatedAt": "",
                "DeletedAt": "",
                "ItemName": ""
            }
        }
    }
    ```

  *失败*
  请查看具体错误error信息，比对纠错

---

## PATCH /:id/update | 更新机构相关信息

- **User**需要*机构方*权限:3级，或*管理员*权限:2级
- 路由样式：[127.0.0.1:xxxx/api/institutions/:id/update](http://127.0.0.1:xxxx/api/institutions/:id/update)
  (*:id*为机构ID)
- 请求体（可选字段，按需填写）：
  控制器内哪个键值不为空，则更新哪个属性值

    ```json
    {
        "plan_id":"", // 不可为空
        "item_id":"", // 不可为空
        "item_name":"", // 可为空
        "item_description":"", // 可为空
        "plan_name":"", // 可为空
        "institution_name":"", // 可为空
        "institution_phone":"", // 可为空
        "institution_address":"", // 可为空
        "institution_qualification":"" // 可为空
    }
    ```

- 返回信息：
  *成功*

    ```json
    {
      "message": "机构信息更新成功"
    }
    ```

  *失败*
  请查看具体错误error信息，比对纠错
