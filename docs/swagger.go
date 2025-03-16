package docs

import "github.com/swaggo/swag"

// docTemplate содержит JSON-шаблон Swagger-документации для API.
const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/info": {
            "get": {
                "description": "Retrieve detailed information about a song, add to database if not present",
                "produces": [
                    "application/json"
                ],
                "summary": "Get song details",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Group",
                        "name": "group",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Song",
                        "name": "song",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.SongDetail"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/songs": {
            "get": {
                "description": "Retrieve all songs with optional filtering and pagination",
                "produces": [
                    "application/json"
                ],
                "summary": "Get all songs",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Group",
                        "name": "group",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Song",
                        "name": "song",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Release Date",
                        "name": "release_date",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Text",
                        "name": "text",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Link",
                        "name": "link",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 1,
                        "description": "Page number",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 10,
                        "description": "Results per page",
                        "name": "limit",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Song"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.Song": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "deleted_at": {
                    "type": "string"
                },
                "group": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "link": {
                    "type": "string"
                },
                "release_date": {
                    "type": "string"
                },
                "song": {
                    "type": "string"
                },
                "text": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        },
        "models.SongDetail": {
            "type": "object",
            "properties": {
                "link": {
                    "type": "string"
                },
                "release_date": {
                    "type": "string"
                },
                "text": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo содержит информацию о Swagger-документации API.
// Этот объект регистрируется в Swag и позволяет клиентам изменять параметры API-документации.
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",                                   // Версия API
	Host:             "localhost:8080",                        // Хост API (можно изменить в продакшене)
	BasePath:         "/",                                     // Базовый путь API
	Schemes:          []string{},                              // Схемы запроса (http, https)
	Title:            "Music Library API",                     // Заголовок API
	Description:      "API для управления библиотекой песен.", // Описание API
	InfoInstanceName: "swagger",                               // Имя инстанса
	SwaggerTemplate:  docTemplate,                             // JSON-шаблон Swagger
	LeftDelim:        "{{",                                    // Левый разделитель для шаблонизации
	RightDelim:       "}}",                                    // Правый разделитель для шаблонизации
}

// init регистрирует Swagger-документацию при загрузке пакета.
func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
