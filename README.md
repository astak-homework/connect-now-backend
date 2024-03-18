# connect-now-backend
API для создания и просмотра профилей в социальной сети

[<img src="https://run.pstmn.io/button.svg" alt="Run In Postman" style="width: 128px; height: 32px;">](https://app.getpostman.com/run-collection/5350211-23546d69-8bea-416f-9f6d-f044aeccfee7?action=collection%2Ffork&source=rip_markdown&collection-url=entityId%3D5350211-23546d69-8bea-416f-9f6d-f044aeccfee7%26entityType%3Dcollection%26workspaceId%3D92f3c2bc-d0cb-40cb-ab96-8694851be6bc)

## Инструкция по локальному запуску приложения

Приложение можно запустить локально с помощью утилиты [Docker Compose](https://docs.docker.com/compose/install/).

```
docker compose up -d
```

Для полноценной работы приложения необходимо указать учетные данные для базы данных, и ключ подписи JWT токенов. Это можно сделать через переменные окружения, добавив их непосредственно перед запуском команды `docker compose up`. Можно использовать следующие переменные:
 - `POSTGRES_USER` - определяет имя пользователя, который будет автоматически создан при первом запуске базы данных. 
 - `POSTGRES_PASSWORD` - определяет пароль пользователя, который будет автоматически создан при первом запуске базы данных.
 - `POSTGRES_DB` - определяет имя базы данных, которая будет автоматически создана при первом запуске.
 - `AUTH_SIGNING_KEY` - секретный ключ для подписи JWT токенов, выдаваемых приложением при аутентификации пользователя. Задается в формате Base64. Для подписи используется метод HMAC, поэтому в качестве ключа можно использовать случайный набор байт произвольной длины. На Linux можно использовать такую команду, чтобы получить случайный набор размером 256 байт в формате Base64: `dd if=/dev/urandom bs=256 count=1 2>/dev/null | base64`

Для удобства, рекомендуется добавить в корень проекта файл `.env`, где и хранить все эти переменные окружения в виде списка `ИМЯ=ЗНАЧЕНИЕ`. Docker Compose будет использовать этот файл при запуске команд в корне проекта.

```
POSTGRES_USER=connectnow
POSTGRES_PASSWORD=N6f5ygYPScYT
POSTGRES_DB=connectnow
AUTH_SIGNING_KEY=XKmXeEWjLZSKvE2MbJ2wPy4C7PVP1oFyDxHaUaaJkWE=
```

> [!NOTE]
> Настройки учетной записи пользователя базы данных читаются из переменных окружения только при первом запуске. При последующих запусках базы данных эти настройки игнорируются, т.к. учетные данные хранятся в томе `postgres_data`. Однако эти переменные окружения будет использоваться при перезапуске REST API сервиса, поэтому значения переменных окружения следует всегда задавать одни и те же. При необходимости изменить значения переменных `POSTGRES_USER`, `POSTGRES_PASSWORD`, и `POSTGRES_DB`, нужно вручную синхронизировать учетную запись в базе данных либо выполнив соотвествующие SQL команды, либо через интерфейс управления `adminer` (доступен через браузер по адресу http://localhost:8080).

## Postman-коллекции
### Тест "Создание учетной записи и анкеты"
```
{
	"info": {
		"_postman_id": "23546d69-8bea-416f-9f6d-f044aeccfee7",
		"name": "Creating Profile",
		"schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
		"_exporter_id": "5350211",
		"_collection_link": "https://www.postman.com/uriahas/workspace/highload-architect/collection/5350211-23546d69-8bea-416f-9f6d-f044aeccfee7?action=share&source=collection_link&creator=5350211"
	},
	"item": [
		{
			"name": "Creating a User",
			"event": [
				{
					"listen": "prerequest",
					"script": {
						"exec": [
							"const first_name = pm.variables.replaceIn('{{$randomFirstName}}');\r",
							"pm.variables.set(\"first_name\", first_name);\r",
							"\r",
							"const last_name = pm.variables.replaceIn('{{$randomLastName}}');\r",
							"pm.variables.set(\"last_name\", last_name);\r",
							"\r",
							"const datePastStr = pm.variables.replaceIn('{{$randomDatePast}}');\r",
							"const birthdate = new Date(Date.parse(datePastStr)).toISOString().split('T')[0];\r",
							"pm.variables.set(\"birthdate\", birthdate);\r",
							"\r",
							"const genders = ['male', 'female'];\r",
							"const gender = genders[Math.floor(Math.random() * genders.length)];\r",
							"pm.variables.set(\"gender\", gender);\r",
							"\r",
							"const biography = pm.variables.replaceIn('{{randomPhrase}}');\r",
							"pm.variables.set(\"biography\", biography);\r",
							"\r",
							"const city = pm.variables.replaceIn('{{$randomCity}}');\r",
							"pm.variables.set('city', city);\r",
							"\r",
							"const password = pm.variables.replaceIn('{{$randomPassword}}');\r",
							"pm.variables.set('password', password);\r",
							""
						],
						"type": "text/javascript",
						"packages": {}
					}
				},
				{
					"listen": "test",
					"script": {
						"exec": [
							"const schema = {\r",
							"    \"properties\": {\r",
							"        \"user_id\": {\r",
							"            \"type\": \"string\"\r",
							"        }\r",
							"    }\r",
							"};\r",
							"\r",
							"pm.test(\"response is ok\", function() {\r",
							"    pm.response.to.have.status(200);\r",
							"});\r",
							"\r",
							"pm.test(\"response body has valid schema\", function() {\r",
							"    pm.response.to.have.jsonSchema(schema);\r",
							"});\r",
							"\r",
							"const responseJson = pm.response.json();\r",
							"pm.variables.set(\"user_id\", responseJson.user_id);"
						],
						"type": "text/javascript",
						"packages": {}
					}
				}
			],
			"request": {
				"method": "POST",
				"header": [],
				"body": {
					"mode": "raw",
					"raw": "{\r\n    \"first_name\": \"{{first_name}}\",\r\n    \"second_name\": \"{{last_name}}\",\r\n    \"birthdate\": \"{{birthdate}}\",\r\n    \"gender\": \"{{gender}}\",\r\n    \"biography\": \"{{biography}}\",\r\n    \"city\": \"{{city}}\",\r\n    \"password\": \"{{password}}\"\r\n}",
					"options": {
						"raw": {
							"language": "json"
						}
					}
				},
				"url": {
					"raw": "http://localhost/user/register",
					"protocol": "http",
					"host": [
						"localhost"
					],
					"path": [
						"user",
						"register"
					]
				}
			},
			"response": []
		},
		{
			"name": "Authenticating a User",
			"event": [
				{
					"listen": "test",
					"script": {
						"exec": [
							"const schema = {\r",
							"    \"properties\": {\r",
							"        \"token\": {\r",
							"            \"type\": \"string\"\r",
							"        }\r",
							"    }\r",
							"};\r",
							"\r",
							"pm.test(\"response is ok\", function() {\r",
							"    pm.response.to.have.status(200);\r",
							"});\r",
							"\r",
							"pm.test(\"response body has valid schema\", function() {\r",
							"    pm.response.to.have.jsonSchema(schema);\r",
							"});\r",
							"\r",
							"const responseJson = pm.response.json();\r",
							"pm.variables.set(\"token\", responseJson.token);\r",
							""
						],
						"type": "text/javascript",
						"packages": {}
					}
				}
			],
			"request": {
				"method": "POST",
				"header": [],
				"body": {
					"mode": "raw",
					"raw": "{\r\n    \"id\": \"{{user_id}}\",\r\n    \"password\": \"{{password}}\"\r\n}",
					"options": {
						"raw": {
							"language": "json"
						}
					}
				},
				"url": {
					"raw": "http://localhost/login",
					"protocol": "http",
					"host": [
						"localhost"
					],
					"path": [
						"login"
					]
				}
			},
			"response": []
		},
		{
			"name": "Getting a User Profile",
			"event": [
				{
					"listen": "test",
					"script": {
						"exec": [
							"const schema = {\r",
							"    \"properties\": {\r",
							"        \"first_name\": {\r",
							"            \"type\": \"string\"\r",
							"        },\r",
							"        \"second_name\": {\r",
							"            \"type\": \"string\"\r",
							"        },\r",
							"        \"birthdate\": {\r",
							"            \"type\": \"string\",\r",
							"            \"format\": \"date\"\r",
							"        },\r",
							"        \"gender\": {\r",
							"            \"type\": \"string\",\r",
							"            \"enum\": [\"male\", \"female\"]\r",
							"        },\r",
							"        \"biography\": {\r",
							"            \"type\": \"string\"\r",
							"        },\r",
							"        \"city\": {\r",
							"            \"type\": \"string\"\r",
							"        }\r",
							"    }\r",
							"};\r",
							"\r",
							"pm.test(\"response is ok\", function() {\r",
							"    pm.response.to.have.status(200);\r",
							"});\r",
							"\r",
							"pm.test(\"response body has valid schema\", function() {\r",
							"    pm.response.to.have.jsonSchema(schema);\r",
							"});\r",
							"\r",
							"pm.test(\"response body has expected values\", function() {\r",
							"    pm.response.to.have.jsonBody('first_name', pm.variables.get('first_name'))\r",
							"        .and.have.jsonBody('second_name', pm.variables.get('last_name'))\r",
							"        .and.have.jsonBody('birthdate', pm.variables.get('birthdate'))\r",
							"        .and.have.jsonBody('gender', pm.variables.get('gender'))\r",
							"        .and.have.jsonBody('biography', pm.variables.get('biography'))\r",
							"        .and.have.jsonBody('city', pm.variables.get('city'));\r",
							"});\r",
							"\r",
							""
						],
						"type": "text/javascript",
						"packages": {}
					}
				}
			],
			"request": {
				"method": "GET",
				"header": [],
				"url": {
					"raw": "http://localhost/user/get/{{user_id}}",
					"protocol": "http",
					"host": [
						"localhost"
					],
					"path": [
						"user",
						"get",
						"{{user_id}}"
					]
				}
			},
			"response": []
		}
	],
	"event": [
		{
			"listen": "prerequest",
			"script": {
				"type": "text/javascript",
				"exec": [
					"pm.request.headers.add({ \r",
					"    key: \"Accept-Language\",\r",
					"    value: \"ru\" \r",
					"});\r",
					""
				]
			}
		},
		{
			"listen": "test",
			"script": {
				"type": "text/javascript",
				"exec": [
					""
				]
			}
		}
	],
	"variable": [
		{
			"key": "user_id",
			"value": "1",
			"type": "string"
		},
		{
			"key": "base_url",
			"value": "arch.homework",
			"type": "string"
		}
	]
}
```