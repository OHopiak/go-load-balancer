# go-load-balancer
Створити веб застосування, backend якого здатний коректно обслуговувати трудомісткі запити користувачів.

Прикладами таких запитів можуть бути задачі: розв"язування системи рівнянь, метод скінченних елементів, розпізнавання образів, обробка зображень та інші, на вибір студента. 

Особливості:
1) Повинна бути обмежена максимальна трудомісткість одної задачі (кількість невідомих, час обчислення).У разі перевищення - помилка або відмова у виконанні. (Перевірку вхідних даних задачі можна виконати ще на клієнті (frontend)).
2) Оскільки запити трудомісткі, клієнт повинен інформуватись про хід виконання задачі (відсоток/етап виконання).
3) Клієнт повинен мати доступ до історії результатів виконання задач, переглядати стан поточної задачі (див. пункт 1.), скасувати виконання задачі, запустити виконання нової задачі. (максимальна кількість задач повинна бути обмежена, необхідні дані зберігатись в базі даних)
4) Авторизація клієнта (з використанням https (+бали))

# Найважливіша частина завдання
5) Веб застосування повинне забезпечувати т.з. балансування навантаження як мінімум для 2-ох серверів. 
(для цього повинен бути застосований так званий load balancer, в його ролі може виступати і сам веб сервер, а задачі виконуватись на кількох серверах застосувань (application server)).


## Додаткові бали
6) (+бали) Якщо всі ресурси системи вичерпані, клієнт повинен бути повідомлений про те, що його запит поставлено у чергу та інформуватись про приблизний час очікування.
7) (+бали) Горизонтальне розширення за рахунок автоматичного запуску нових серверів (віртуальні машини, контейнери), максимальна кількість повинна бути обмежена.
8) (+бали) Панель адміністратора для моніторингу всіх запущених процесів та навантаження на систему серверів.

Кількість балів: 
До 2-рьох тижнів - 35 балів.
До 4-рьох тижнів включно - 30балів
5-тижнів 25 балів
6 тижнів 20 балів.
