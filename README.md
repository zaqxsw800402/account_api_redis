# account_api_golang

# 資料夾介紹

mux資料夾用的模組有mux跟sqlx

gin資料夾用gin跟gorm

兩份程式碼用一樣的架構及功能

# 架構圖

![image](https://github.com/zaqxsw800402/account_api_golang/blob/master/picture/golang_api.png?raw=true)

# Database Table

顧客表

查詢顧客的資料

帳戶表

查詢顧客開啟的帳戶

交易表

記錄所有交易紀錄

# URL

/customers

get: 查詢有哪些顧客

/custoemer/:id

get:查詢特製的顧客

/customer/:id/account

post: 創辦帳戶

/customer/:id/account/:account_id

post: 在該帳戶下存提款