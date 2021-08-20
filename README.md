## 架构设计

![无标题](D:\GoProjects\src\netdisk\无标题.png)

## 已实现的功能

文件上传 

文件路径修改、重命名、删除

文件下载

登录注册

文件分享、权限管理（获得链接的人可下载、需要提取码下载、仅自己见）

断点续传

加密分享链接

二维码分享链接

下载限速

## 接口介绍

### 登录

**简要描述:**

- 登录使用

**请求方式以及URL:**

- GET http://localhost:8080/user/login

**参数:**

| 参数名   | 必选 | 类型   | 说明     |
| -------- | ---- | ------ | -------- |
| username | 是   | string | 用户名   |
| password | 是   | string | 用户密码 |
|          |      |        |          |

**返回:**

header中的token字段以及成功与否



### 注册

**简要描述:**

- 注册使用

**请求方式以及URL:**

- POST	http://localhost:8080/user/register

**参数:**

| 参数名   | 必选 | 类型   | 说明     |
| -------- | ---- | ------ | -------- |
| username | 是   | string | 用户名   |
| password | 是   | string | 用户密码 |
|          |      |        |          |

**返回:**

header中的token字段以及成功与否



### 文件上传

**简要描述:**

- 上传文件使用

**请求方式以及URL:**

- POST	http://localhost:8080/file/upload

**参数:**

| 参数名 | 必选 | 类型 | 说明           |
| ------ | ---- | ---- | -------------- |
| upload | 是   | file | 需要上传的文件 |

**header中需要携带:**

| 参数名 | 必选 | 类型   | 说明                                                   |
| ------ | ---- | ------ | :----------------------------------------------------- |
| token  | 是   | string | 识别用户的凭证，需在登录或者注册时获取。有效时常30分钟 |

**返回:**

成功是否的提示string



### 文件上传

**简要描述:**

- 查询当前用户拥有的所有文件

**请求方式以及URL:**

- GET http://localhost:8080/file/query

**参数:**

无

**header中需要携带:**

| 参数名 | 必选 | 类型   | 说明                                                   |
| ------ | ---- | ------ | :----------------------------------------------------- |
| token  | 是   | string | 识别用户的凭证，需在登录或者注册时获取。有效时常30分钟 |

**返回示例:**

```json
[
    {
        "FileName": "begonia源码.md",
        "FileSize": 13187,
        "Md5hash": "3f78d4cb488fa7f52f5db302e38c8ca7",
        "Path": "file\\222\\begonia源码.md",
        "Username": "222",
        "Power": 0,
        "Secret": "",
        "UploadTime": "2021-08-20T16:13:57+08:00"
    }
]
```



### 文件下载

**简要描述:**

- 下载文件或获得别人分享的文件

**请求方式以及URL:**

- GET http://localhost:8080/file/download

**参数:**

| 参数名  | 必选 | 类型   | 说明                   |
| ------- | ---- | ------ | ---------------------- |
| md5hash | 是   | string | 文件的唯一标识         |
| secret  | 否   | string | 当文件需要提取码时携带 |

**header中需要携带:**

| 参数名 | 必选 | 类型   | 说明                                                   |
| ------ | ---- | ------ | :----------------------------------------------------- |
| token  | 否   | string | 识别用户的凭证，需在登录或者注册时获取。有效时常30分钟 |

**返回:**

成功是否的提示string



### 文件分享

**简要描述:**

- 更改文件的分享权限等级并获取下载链接和二维码

**请求方式以及URL:**

- GET http://localhost:8080/file/share

**参数:**

| 参数名  | 必选 | 类型   | 说明                                          |
| ------- | ---- | ------ | --------------------------------------------- |
| power   | 是   | int    | 0为仅自己可见，1为需要提取码，2为所有人可下载 |
| md5hash | 是   | string | 文件的唯一标识                                |
| sercret | 否   | string | 当power要改为1时需要填写                      |

**header中需要携带:**

| 参数名 | 必选 | 类型   | 说明                                                   |
| ------ | ---- | ------ | :----------------------------------------------------- |
| token  | 是   | string | 识别用户的凭证，需在登录或者注册时获取。有效时常30分钟 |

**返回:**

```json
{
    "p": "iVBORw0KGgoAAAANSUhEUgAAAQAAAAEAAQMAAABmvDolAAAABlBMVEX///8AAABVwtN+AAADGklEQVR42uyYMY77LBTEx6KgfDdYLhKZa32FJSNR/K9lKxchN6CksDyfhk2yq+03oViUIlJ+RYCZ9+aBv/W3xltG8ogVgCtY4DNJYg4k2zhAApaJx4e+rGwLEOm3AOB1wD/WZWK+MReyeJJb8FvwrIMBsPSBteg8tYvizwGByBTqyvu6Fo4GSA+z5dNysYTmCiKBH4J5NyDjzJZvdQmWS1uCv5b201lvBvqaKqbq2HfBFunPH3XolwHd8mZpYq8tTGC3Nrf7SQ4BeDJFpg/LlMFz4QHP0jASsG6qk5nyjlMJAgJ3YiIGAvbqdtVJsq6lOZ236uRexwGsRR6RmTLOgbbqkP0W2mzjAFhPS1NFkFy7o/0RMMPvD3e/AADsiDwmpmAp8OhSjGwwjAP4zbLcXR3h6DNVBiO/djECsG7dO7AjVKAtaBc06LoHAu6nrSRWF3U9pcFYuNeBgPWUHnLRLlyRYiWGgO/GeTdgCV0P+tVykUCubJN28TJANiYrPlSi2XtK/588H9c9AqDwMMOpUCs/9DaoEvRszSMAfq/LjPWEK/fWzB5Zd2IYwNpUEet63wVT8FclMc86EJD63R/3JKaud0HXAwYCuFneq6OAA/psas2YHhPK7wMS36LB3Eg7Pkc5cgPPr673dsBUazTKMZfqChxxCfJOfE4o7wfwXxeku1UES/D9qFWOntscALAWqzvhbnVlT4NUnZwK5se8OQLAHh54qjUrV6vptBn+Sw8DAClKD0eQJNL9RUtF8imYFwD/FLRkmVXxXky/7udcPATQAHf2NwTAPcJ/LO1LD+8H1JqXqa5FxWdVp2uu/xqJcQBZ+zTeVMZZZJxrkcGnZ/h/PyDRrqf1yd36GwguoX1LYiMAfc2Wpm4uDcXKYBpI68uAz1c1uBupXM0UmiPm4DfDOEB/n+yC7EGL9EqD34PWAEB/6a2OdenPC4u6nmrl9v1hfwiAmWrNR/C9Tn5+MBggg2dFa59Lu4Q24/nsNgQgPcS6fD5gUkVSooVmgpcBMs6koOVYHdVNEPQPZxsH+Ft/65Xr/wAAAP//icWwmod9ASwAAAAASUVORK5CYII=",
    "url": "http://localhost:8080/file/download?md5hash=6f1da8a762bb7b2f950a913a56f0c839"
}
```

p为二维码图片的[]byte，传入前端后转换为图片即可

url为下载的链接





### 改变存放路径

**简要描述:**

- 改变文件的存放位置

**请求方式以及URL:**

- PUT	http://localhost:8080/file/change_path

**参数:**

| 参数名   | 必选 | 类型   | 说明           |
| -------- | ---- | ------ | -------------- |
| new_path | 是   | string | 比如/abc/      |
| md5hash  | 是   | string | 需要操作的文件 |

**header中需要携带:**

| 参数名 | 必选 | 类型   | 说明                                                   |
| ------ | ---- | ------ | :----------------------------------------------------- |
| token  | 是   | string | 识别用户的凭证，需在登录或者注册时获取。有效时常30分钟 |

**返回:**

成功是否的提示string



### 文件重命名

**简要描述:**

- 文件重命名

**请求方式以及URL:**

- PUT	http://localhost:8080/file/rename

**参数:**

| 参数名   | 必选 | 类型   | 说明           |
| -------- | ---- | ------ | -------------- |
| new_name | 是   | string | 新的文件名     |
| md5hash  | 是   | string | 需要操作的文件 |

**header中需要携带:**

| 参数名 | 必选 | 类型   | 说明                                                   |
| ------ | ---- | ------ | :----------------------------------------------------- |
| token  | 是   | string | 识别用户的凭证，需在登录或者注册时获取。有效时常30分钟 |

**返回:**

成功是否的提示string



### 文件删除

**简要描述:**

- 删除文件本身及数据库记录

**请求方式以及URL:**

- DELETE	http://localhost:8080/file/delete

**参数:**

| 参数名  | 必选 | 类型   | 说明           |
| ------- | ---- | ------ | -------------- |
| md5hash | 是   | string | 需要操作的文件 |

**header中需要携带:**

| 参数名 | 必选 | 类型   | 说明                                                   |
| ------ | ---- | ------ | :----------------------------------------------------- |
| token  | 是   | string | 识别用户的凭证，需在登录或者注册时获取。有效时常30分钟 |

**返回:**

成功是否的提示string



