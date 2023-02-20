GET /api/getUserInfo?id=user_id

user_id - Telegram ID получаемого пользователя

При успехе возвращает объект User с полями:    

    -id - id пользователя в Telegram    
    -name - username пользователя в Telegram     
    -balance - Баланс пользователя   
    -hasKey - Имеет ли пользователь ключ доступа 


Если пользователя нет ничего не возвращает 

---

GET /api/getUserPhoto?id=user_id

    -id - Telegram ID получаемого пользователя

При успехе возвращает последнюю фотографию пользователя в Telegram


---

GET /api/getKeyInfo?by=method&owner=user_id

method - Метод получения информации о ключе     

    -ownerId - Требует указание Telegram ID владельца ключа в поле owner

owner - Telegram ID владельца ключа