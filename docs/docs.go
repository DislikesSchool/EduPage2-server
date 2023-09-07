// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

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
        "/api/timeline": {
            "get": {
                "description": "Returns the user's timeline from any date to any other date or today.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "timeline"
                ],
                "summary": "Get the user's timeline",
                "parameters": [
                    {
                        "type": "string",
                        "description": "JWT token",
                        "name": "token",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "example": "2022-01-01T00:00:00Z2022-01-01T00:00:00Z",
                        "name": "from",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "default": "time.Now()",
                        "example": "2022-01-01T00:00:00Z",
                        "name": "to",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.Timeline"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/main.TimelineUnauthorizedResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/main.TimelineInternalErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/timeline/recent": {
            "get": {
                "description": "Returns the user's timeline from today to 30 days in the past.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "timeline"
                ],
                "summary": "Get the user's recent timeline",
                "parameters": [
                    {
                        "type": "string",
                        "description": "JWT token",
                        "name": "token",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.Timeline"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/main.RecentTimelineUnauthorizedResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/main.RecentTimelineInternalErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/timetable": {
            "get": {
                "description": "Returns the user's timetable from date specified to date specified or today.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "timetable"
                ],
                "summary": "Get the user's  timetable",
                "parameters": [
                    {
                        "type": "string",
                        "description": "JWT token",
                        "name": "token",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "example": "2022-01-01T00:00:00Z2022-01-01T00:00:00Z",
                        "name": "from",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "default": "time.Now()",
                        "example": "2022-01-01T00:00:00Z",
                        "name": "to",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Timetable"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/main.TimetableUnauthorizedResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/main.TimetableInternalErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/timetable/recent": {
            "get": {
                "description": "Returns the user's timetable from before yesterday to 7 days in the future.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "timetable"
                ],
                "summary": "Get the user's recent timetable",
                "parameters": [
                    {
                        "type": "string",
                        "description": "JWT token",
                        "name": "token",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Timetable"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/main.TimetableUnauthorizedResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/main.TimetableInternalErrorResponse"
                        }
                    }
                }
            }
        },
        "/login": {
            "post": {
                "description": "Logs in to your Edupage account using the provided credentials.",
                "consumes": [
                    "application/json",
                    "multipart/form-data",
                    "application/x-www-form-urlencoded"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Login to your Edupage account",
                "parameters": [
                    {
                        "description": "Login using username and password",
                        "name": "login",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/main.LoginRequestUsernamePassword"
                        }
                    },
                    {
                        "description": "Login using username, password and server",
                        "name": "loginServer",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/main.LoginRequestUsernamePasswordServer"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.LoginSuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/main.LoginBadRequestResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/main.LoginUnauthorizedResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/main.LoginInternalErrorResponse"
                        }
                    }
                }
            }
        },
        "/validate-token": {
            "get": {
                "description": "Validates your token and returns a 200 OK if it's valid.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Validate your token",
                "parameters": [
                    {
                        "type": "string",
                        "description": "JWT token",
                        "name": "token",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.ValidateTokenSuccessResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/main.ValidateTokenUnauthorizedResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "main.HomeworkReduced": {
            "type": "object",
            "properties": {
                "attachements": {},
                "data": {
                    "$ref": "#/definitions/model.TimelineItemData"
                },
                "datecreated": {
                    "type": "string"
                },
                "details": {
                    "type": "string"
                },
                "homeworkid": {
                    "type": "string"
                },
                "hwkid": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "period": {},
                "pocet_done": {
                    "type": "string"
                },
                "pocet_like": {
                    "type": "string"
                },
                "pocet_reakcii": {
                    "type": "string"
                },
                "posledny_vysledok": {
                    "type": "string"
                },
                "predmetid": {
                    "type": "string"
                },
                "stav": {
                    "type": "string"
                },
                "stavhodnotetimelinePathd": {
                    "type": "string"
                },
                "students_hidden": {
                    "type": "string"
                },
                "testid": {
                    "type": "string"
                },
                "timestamp": {
                    "type": "string"
                },
                "typ": {
                    "$ref": "#/definitions/model.TimelineItemType"
                },
                "userid": {
                    "type": "string"
                },
                "znamky_udalostid": {}
            }
        },
        "main.LoginBadRequestResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string",
                    "example": "Username and Password are required"
                },
                "success": {
                    "type": "boolean",
                    "example": false
                }
            }
        },
        "main.LoginInternalErrorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string",
                    "example": "failed to login: Post https://example.edupage.org/login/edubarLogin.php: dial tcp: lookup example.edupage.org: no such host"
                },
                "success": {
                    "type": "boolean",
                    "example": false
                }
            }
        },
        "main.LoginRequestUsernamePassword": {
            "type": "object",
            "required": [
                "password",
                "username"
            ],
            "properties": {
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "main.LoginRequestUsernamePasswordServer": {
            "type": "object",
            "required": [
                "password",
                "server",
                "username"
            ],
            "properties": {
                "password": {
                    "type": "string"
                },
                "server": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "main.LoginSuccessResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string",
                    "example": ""
                },
                "firstname": {
                    "type": "string",
                    "example": "John"
                },
                "lastname": {
                    "type": "string",
                    "example": "Doe"
                },
                "success": {
                    "type": "boolean",
                    "example": true
                },
                "token": {
                    "type": "string",
                    "example": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM"
                }
            }
        },
        "main.LoginUnauthorizedResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string",
                    "example": "Unexpected response from server, make sure credentials are specified correctly"
                },
                "success": {
                    "type": "boolean",
                    "example": false
                }
            }
        },
        "main.RecentTimelineInternalErrorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string",
                    "example": "failed to create payload"
                }
            }
        },
        "main.RecentTimelineUnauthorizedResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string",
                    "example": "Unauthorized"
                }
            }
        },
        "main.Timeline": {
            "type": "object",
            "properties": {
                "homeworks": {
                    "type": "object",
                    "additionalProperties": {
                        "$ref": "#/definitions/main.HomeworkReduced"
                    }
                },
                "items": {
                    "type": "object",
                    "additionalProperties": {
                        "$ref": "#/definitions/main.TimelineItemReduced"
                    }
                }
            }
        },
        "main.TimelineInternalErrorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string",
                    "example": "failed to create payload"
                }
            }
        },
        "main.TimelineItemReduced": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/model.TimelineItemData"
                },
                "poct_reakcii": {
                    "type": "integer"
                },
                "reakcia_na": {
                    "type": "string"
                },
                "removed": {
                    "type": "string"
                },
                "target_user": {
                    "type": "string"
                },
                "text": {
                    "type": "string"
                },
                "timelineid": {
                    "type": "string"
                },
                "timestamp": {
                    "type": "string"
                },
                "typ": {
                    "$ref": "#/definitions/model.TimelineItemType"
                },
                "user": {
                    "type": "string"
                },
                "vlastnik": {
                    "type": "string"
                }
            }
        },
        "main.TimelineUnauthorizedResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string",
                    "example": "Unauthorized"
                }
            }
        },
        "main.TimetableInternalErrorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string",
                    "example": "failed to create payload"
                }
            }
        },
        "main.TimetableUnauthorizedResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string",
                    "example": "Unauthorized"
                }
            }
        },
        "main.ValidateTokenSuccessResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string",
                    "example": ""
                },
                "expires": {
                    "type": "string",
                    "example": "1620000000"
                },
                "success": {
                    "type": "boolean",
                    "example": true
                }
            }
        },
        "main.ValidateTokenUnauthorizedResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string",
                    "example": "Unauthorized"
                },
                "success": {
                    "type": "boolean",
                    "example": false
                }
            }
        },
        "model.TimelineItemData": {
            "type": "object",
            "properties": {
                "value": {
                    "type": "object",
                    "additionalProperties": true
                }
            }
        },
        "model.TimelineItemType": {
            "type": "object",
            "properties": {
                "uint8": {
                    "type": "integer"
                }
            }
        },
        "model.Timetable": {
            "type": "object",
            "properties": {
                "days": {
                    "description": "key format is YYYY-MM-dd or 2006-01-02",
                    "type": "object",
                    "additionalProperties": {
                        "type": "array",
                        "items": {
                            "$ref": "#/definitions/model.TimetableItem"
                        }
                    }
                }
            }
        },
        "model.TimetableItem": {
            "type": "object",
            "properties": {
                "classids": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "classroomids": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "colors": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "date": {
                    "type": "string"
                },
                "endtime": {
                    "type": "string"
                },
                "groupnames": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "igroupid": {
                    "type": "string"
                },
                "starttime": {
                    "type": "string"
                },
                "studentids": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "subjectid": {
                    "type": "string"
                },
                "teacherids": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "type": {
                    "type": "string"
                },
                "uniperiod": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "EduPage2 API",
	Description:      "This is the backend for the EduPage2 app.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
