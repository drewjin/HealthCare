### API修改介绍

## PATCH /api/user/:id/profile **用户更新个人信息**

- **User**需要*任意*权限
- 路由样式：
  [127.0.0.1:xxxx/api/user/:id/profile](https://127.0.0.1:xxxx/api/user/:id/profile)
    (*:id*应为当前用户的User_id)
- 请求体：
    **最好都不是空值**

  ```json
        {
            //均设置允许空值
            "username":"",
            "name":"",
            "gender":"",
            "birthday":"",
            "phone":"",
            "email":"",
            "address":""
        }
  ```

-返回信息：
*成功*

```json
        {
            "message": "User profile updated successfully"
        }
```

*失败*
请查看具体错误error信息，比对纠错

---

## DELETE /api/user/:id **删除用户（物理删除）**

-User需要管理员或本人权限

-路由样式：
    [127.0.0.1:xxxx/api/user/:id](https://127.0.0.1:xxxx/api/user/:id)
    (*:id*应为当前用户的User_id)

-请求体：无

-返回信息：
**成功**

```json
  {
      "message": "User deleted successfully",
      "logout": true
  }
```

**失败**  
  请查看具体错误error信息，比对纠错

---

## PATCH /api/user/:id/permission **管理员更改用户权限**

- **用户注销后应跳转登录界面**
- User需要**管理员**权限

- 路由样式：  
  [127.0.0.1:xxxx/api/user/:id/permission](https://127.0.0.1:xxxx/api/user/:id/permission)  
  (*:id*应为目标用户的User_id)

- 请求体：

  ```json
  {
      "user_type": 2
  }
  ```

- 返回信息：  
  **成功**

  ```json
  {
      "message": "User permission updated successfully",
      "user_id": "3",
      "user_type": 2,
      "logout": true
  }
  ```

  **失败**  
  请查看具体错误error信息，比对纠错

  ---

## POST /api/commentary/add **发布评论**

- User需要**任意**权限

- 路由样式：  
  [127.0.0.1:xxxx/api/commentary/add](https://127.0.0.1:xxxx/api/commentary/add)

- 请求体：

  ```json
  {
      "user_id": 3,
      "plan_id": 2,
      "commentary": ""
  }
  ```

- 返回信息：  
  **成功**

  ```json
  {
      "message": "Commentary added successfully"
  }
  ```

  **失败**  
  请查看具体错误error信息，比对纠错

---

## DELETE /api/commentary/delete/:id **删除评论（物理删除）**

- User需要登录权限，且为**评论创建者**或**管理员**

- 路由样式：  
  [127.0.0.1:xxxx/api/commentary/delete/:id](https://127.0.0.1:xxxx/api/commentary/delete/:id)  
  (*:id*为评论的ID)

- 请求体：无

- 返回信息：  
  **成功**

  ```json
  {
      "message": "Commentary deleted successfully"
  }
  ```

  **失败**  
  请查看具体错误error信息，比对纠错

---

## GET /api/commentary/get/plan/:id **查看评论（按套餐ID）**

- User需要**任意**权限

- 路由样式：  
  [127.0.0.1:xxxx/api/commentary/get/plan/:id](https://127.0.0.1:xxxx/api/commentary/get/plan/:id)  
  (*:id*为套餐的ID)

- 请求体：无

- 返回信息：  
  **成功**

  ```json
  {
      "commentaries": [
          {
              "id": null,
              "user_id": null,
              "plan_id": null,
              "commentary": "",
              "created_at": ""
          }
      ]
  }
  ```

  **失败**  
  请查看具体错误error信息，比对纠错

---

## GET /api/commentary/get/user **查看评论（按用户ID）**

- User需要**登录**权限

- 路由样式：  
  [127.0.0.1:xxxx/api/commentary/get/user](https://127.0.0.1:xxxx/api/commentary/get/user)

- 请求体：无

- 返回信息：  
  **成功**

  ```json
  {
      "commentaries": [
          {
              "id": null,
              "user_id": null,
              "plan_id": null,
              "commentary": "",
              "created_at": ""
          }
      ]
  }
  ```

  **失败**  
  请查看具体错误error信息，比对纠错

---

## DELETE /api/institutions/plan/item **删除套餐内一个体检项目**

- User需要**机构用户**或**管理员**权限

- 路由样式：  
  [127.0.0.1:xxxx/api/institutions/plan/item](https://127.0.0.1:xxxx/api/institutions/plan/item)

- 请求体：

  ```json
  {
      "plan_id": null,
      "item_id": null
  }
  ```

- 返回信息：  
  **成功**

  ```json
  {
      "message": "套餐内体检项目删除成功"
  }
  ```

  **失败**  
  请查看具体错误error信息，比对纠错

---

## DELETE /api/institutions/plan **删除套餐**

- User需要**机构用户**或**管理员**权限

- 路由样式：  
  [127.0.0.1:xxxx/api/institutions/plan](https://127.0.0.1:xxxx/api/institutions/plan)

- 请求体：

  ```json
  {
      "plan_id": null
  }
  ```

- 返回信息：  
  **成功**

  ```json
  {
      "message": "套餐删除成功"
  }
  ```

  **失败**  
  请查看具体错误error信息，比对纠错

---

## GET /api/userview/ **查看所有体检项目（包含所有体检套餐）**

- User需要**任意**权限

- 路由样式：  
  [127.0.0.1:xxxx/api/userview/](https://127.0.0.1:xxxx/api/userview/)

- 请求体：无

- 返回信息：  
  **成功**

  ```json
  {
      "items": [
          {
              "item_name": "",
              "item_value": ""
          },
          {
              "item_name": "",
              "item_value": ""
          }
      ]
  }
  ```

  **失败**  
  请查看具体错误error信息，比对纠错

---

## GET /api/userview/plan **查看指定套餐体检项目**

- User需要登录权限

- 路由样式：  
  [127.0.0.1:xxxx/api/userview/plan](https://127.0.0.1:xxxx/api/userview/plan)

- 请求体：

  ```json
  {
      "plan_id": null
  }
  ```

- 返回信息：  
  **成功**

  ```json
  {
      "plan_items": [
          {
              "plan_name": "",
              "item_name": "",
              "item_value": ""
          },
          {
              "plan_name": "",
              "item_name": "",
              "item_value": ""
          }
      ]
  }
  ```

  **失败**  
  请查看具体错误error信息，比对纠错

---

## POST /api/adduserdata/:customer_id/:plan_id **为用户添加体检数据**

- User需要**机构用户**权限

- 路由样式：  
  [127.0.0.1:xxxx/api/adduserdata/:customer_id/:plan_id](https://127.0.0.1:xxxx/api/adduserdata/:customer_id/:plan_id)  
  (*:customer_id*为目标用户ID，*:plan_id*为套餐ID)

- 请求体：

  ```json
  [
      { "item_id": null, "item_value": "" },
      { "item_id": null, "item_value": "" }
  ]
  ```

- 返回信息：  
  **成功**

  ```json
  {
      "message": "User health items updated successfully"
  }
  ```

  **失败**  
  请查看具体错误error信息，比对纠错

---