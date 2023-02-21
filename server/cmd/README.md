GET /api/getUserInfo?id=user_id

user_id - Telegram ID получаемого пользователя

При успехе возвращает JSON объект с полями:    

    -id - id пользователя в Telegram (int)
    -name - username пользователя в Telegram (string)     
    -balance - Баланс пользователя  (float) 
    -hasKey - Имеет ли пользователь ключ доступа (boolean) 


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

При успехе возвращает JSON объект с полями

    -key - Ключ пользователя (string)
    -type - Тип подписки пользователя (string)   
    -ip - IP адрес к которому привязан ключ (string)
    -expire - Дата экспирации ключа (string)
    -is_expired - Истек ли ключ (boolean)
    -owner - Telegram ID владельца (int)