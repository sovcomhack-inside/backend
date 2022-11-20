# Sovcombank Team Challenge 2022 | Inside

## Состав команды

- Артемий Звонарев | [@artemiyzvonarev](https://t.me/artemiyzvonarev) | **Product Manager** **BizDev**
- Никита Биченов | [@weij33t](https://t.me/weij33t) | **Frontend Engineer**
- Кирилл Смирнов | [@yabhzett](https://t.me/yabhzett) | **Designer**
- Илья Новиков | [@ougirez](https://t.me/ougirez) | **Backend Engineer**
- Оганес Мирзоян | [@senago](https://t.me/senago) | **Fullstack Engineer**

## Ссылки

- [Figma](https://www.figma.com/file/wg5W6zodk9CvxH4yV887PX/SovComHack?node-id=0%3A1&t=f1DZOc9YqKpub7IL-1)
- [Customer Journey Map](https://www.figma.com/file/DKiZ8ZQDANAlJBGOugYrEi/Inside-Job?t=XlC4L3wtSGN1YSp6-0)
- Репозиторий фронтенда: [github.com/sovcomhack-inside/frontend](https://github.com/sovcomhack-inside/frontend)

## Запуск

1. Заполнить конфиг [resources/config/config_default.yaml](resources/config/config_default.yaml)
2. Запустить `make`

## API

Документация описана в [api/swagger.yaml](api/swagger.yaml)

## Архитектура приложения

1. Обработчики ручек - [internal/api/controller](internal/api/controller)
2. Сервисы с бизнес логикой [internal/pkg/service](internal/pkg/service)
3. Запросы к БД - [internal/pkg/store](internal/pkg/store)
