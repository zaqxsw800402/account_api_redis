# account_api_golang

運用gin、gorm、redis來架設api

# 架構圖

![image](https://github.com/zaqxsw800402/account_api_redis/blob/master/picture/redis.png?raw=true)

# Database Table

顧客表

&ensp;&ensp;顧客的資料

帳戶表

&ensp;&ensp;顧客創立哪些帳戶

交易表

&ensp;&ensp;記錄所有交易紀錄

# Redis
在一分鐘內只能申請五次帳號 <BR>
在一分鐘內只能申請五次交易  <BR>
提供GET的緩存<BR>

# URL
/customers

&ensp;&ensp;&ensp;get: 查詢全部顧客<BR>
&ensp;&ensp;&ensp;post: 增加顧客<BR>

/custoemer/:id

&ensp;&ensp;&ensp;get:查詢特定的顧客

/customer/:id/account

&ensp;&ensp;&ensp;post: 在特定的顧客下創辦帳戶

/customer/:id/account/:account_id

&ensp;&ensp;&ensp;post: 在該帳戶下存提款 <BR>
&ensp;&ensp;&ensp;get: 查詢該帳戶的交易紀錄<BR>