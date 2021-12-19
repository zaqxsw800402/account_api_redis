# Bank

## 架構圖
![image](https://github.com/zaqxsw800402/account_api_redis/blob/master/picture/redis.png?raw=true)

架設模擬能夠存錢、提錢、轉帳的網頁(僅限本網站) <br>
### api(backend)
主要使用gin、gorm、redis來架設
### web(frontend)
主要使用html、js
### cronjob
固定時間來清理一些軟刪除的資料，目前設置每兩分鐘清一次
### mailer
負責提供寄信的微服務(開辦帳號，密碼重設)，目前收發信件只能在[mailtrap](https://mailtrap.io/) 裡面查看
### redis
用在 all customer, all account 提供資料的緩存，及middleware 驗證的緩存 
### nsq
用來傳遞資料給mailer發送信件
### travis.ci
練習ci cd到aws ebs上


