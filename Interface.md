## 云中台用户管理接口

### 用户注册

`url`  

/cp-ua/v1/user/register

`method`  

post

`formdata`

- phone 用户手机号, 必须
- name 用户姓名，与gogs的用户名一致， 必须
- email 用户邮箱，与gogs的用户邮箱一致， 必须 
- password 用户密码，与gogs的密码一致，必须

`response`

```
{
  "code": 0,
  "info": "Success"
}
```

### 用户登录

`url`  

/cp-ua/v1/user/login

`method`  

post

`formdata`

- name 用户名，必须
- password 用户密码，必须

`response`

```
{
  "code": 0,
  "user": {
    "_id": "58d0998cc3666e21441349d9",
    "phone": "18099999999",
    "name": "patrick1",
    "reg_date": "2017-03-21T11:10:04.117+08:00",
    "no_enc_pwd": "123456",
    "email": "patrick1@126.com"
  }
}
```

### 用户登出

`url`  

/cp-ua/v1/user/logout

`method`  

post

`formdata`

- name 用户名，必须

`response`

```
{
  "code": 0,
  "info": "Success"
}
```

### 用户信息获取

`url`  

/cp-ua/v1/user/{name}

`method`  

get


`response`

```
{
  "code": 0,
  "user": {
    "_id": "58d0998cc3666e21441349d9",
    "phone": "18099999999",
    "name": "patrick1",
    "password": "6f841ebcce09ea0a372c0bddefad0d976f45118ceae734ad4a903a544946364ac3a8efd0155208c36b9eff5cdf2048d58049e2983172b697b78b3ab137093069",
    "salt": "d7a2f1153c5f2424e3bbc5f7e605873308cdd8db78d2c4403ef846dc92df41f0431dfb9478b67080eb0a679d928fd7c65df4f7d32a1f7e28616a0f54dc8f53b2",
    "reg_date": "2017-03-21T11:10:04.117+08:00",
    "no_enc_pwd": "123456",
    "email": "patrick1@126.com"
  }
}
```




