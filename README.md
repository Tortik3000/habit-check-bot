## Habit check bot

Данный бот умеет добавлить и удалять привыки,
выводить check-list сегодняшних привычек. Также может вывести календарь, где помечены дни в которых выполнены все привычки.


## Run

```shell
  make built
   ./bin/bot 
```
###
Переменные окружения по умолчанию

.env:
```
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_DB=bot
POSTGRES_USER=ed
POSTGRES_PASSWORD=1234567

TG_HOST=api.telegram.org
TELEGRAM_BOT_TOKEN = "8485788686:AAFDEl7cUKS6GHVw5_0Xee0zuEoo6APbY8k"
```

## ToDo
- [ ] Написать нормальный код(calendar, time_events)
- [ ] Привычки с определенными днями недели
- [ ] Добавить секреты
- [ ] Получение check по определенному дню
- [ ] Deploy with ci/cd

