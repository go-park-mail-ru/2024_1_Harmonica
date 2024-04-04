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
        "/is_auth": {
            "get": {
                "description": "Get user by request cookie",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Authorization"
                ],
                "summary": "Get auth status",
                "parameters": [
                    {
                        "type": "string",
                        "default": "session-token=",
                        "description": "session-token",
                        "name": "Cookie",
                        "in": "header"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entity.UserResponse"
                        }
                    },
                    "400": {
                        "description": "Possible code responses: 3.",
                        "schema": {
                            "$ref": "#/definitions/errs.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Possible code responses: 2.",
                        "schema": {
                            "$ref": "#/definitions/errs.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Possible code responses: 11.",
                        "schema": {
                            "$ref": "#/definitions/errs.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/likes/{pin_id}/users": {
            "get": {
                "description": "Get users that liked pin by pin ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Likes"
                ],
                "summary": "Get last 20 users that liked pin",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Pin ID",
                        "name": "pin_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "session-token=",
                        "description": "session-token",
                        "name": "Cookie",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entity.UserList"
                        }
                    },
                    "400": {
                        "description": "Possible code responses: 12.",
                        "schema": {
                            "$ref": "#/definitions/errs.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Possible code responses: 11.",
                        "schema": {
                            "$ref": "#/definitions/errs.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/login": {
            "post": {
                "description": "Login user by request.body json",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Authorization"
                ],
                "summary": "Login user",
                "parameters": [
                    {
                        "type": "string",
                        "default": "session-token=",
                        "description": "session-token",
                        "name": "Cookie",
                        "in": "header"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entity.User"
                        }
                    },
                    "400": {
                        "description": "Possible code responses: 3, 4, 5.",
                        "schema": {
                            "$ref": "#/definitions/errs.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Possible code responses: 7, 8.",
                        "schema": {
                            "$ref": "#/definitions/errs.ErrorResponse"
                        }
                    },
                    "403": {
                        "description": "Possible code responses: 1.",
                        "schema": {
                            "$ref": "#/definitions/errs.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Possible code responses: 11.",
                        "schema": {
                            "$ref": "#/definitions/errs.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/logout": {
            "get": {
                "description": "Logout user by their session cookie",
                "tags": [
                    "Authorization"
                ],
                "summary": "Logout user",
                "parameters": [
                    {
                        "type": "string",
                        "default": "session-token=",
                        "description": "session-token",
                        "name": "Cookie",
                        "in": "header"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Possible code responses: 3.",
                        "schema": {
                            "$ref": "#/definitions/errs.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/pins": {
            "get": {
                "description": "Get pins by page",
                "tags": [
                    "Pins"
                ],
                "summary": "Pins list",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Page num from 0",
                        "name": "page",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entity.FeedPins"
                        }
                    },
                    "400": {
                        "description": "Possible code responses: 4.",
                        "schema": {
                            "$ref": "#/definitions/errs.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Possible code responses: 11.",
                        "schema": {
                            "$ref": "#/definitions/errs.ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "Create pin by description",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Pins"
                ],
                "summary": "Create pin",
                "parameters": [
                    {
                        "type": "string",
                        "default": "session-token=",
                        "description": "session-token",
                        "name": "Cookie",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Pin information in json",
                        "name": "pin",
                        "in": "formData"
                    },
                    {
                        "type": "file",
                        "description": "Pin image",
                        "name": "image",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entity.PinPageResponse"
                        }
                    },
                    "400": {
                        "description": "Possible code responses: 3, 4, 15, 18, 19.",
                        "schema": {
                            "$ref": "#/definitions/errs.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Possible code responses: 2.",
                        "schema": {
                            "$ref": "#/definitions/errs.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Possible code responses: 11.",
                        "schema": {
                            "$ref": "#/definitions/errs.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/pins/created/{nickname}": {
            "get": {
                "description": "Get pins of user by page",
                "tags": [
                    "Pins"
                ],
                "summary": "Get pins that created by user id",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Page num from 0",
                        "name": "page",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entity.FeedPins"
                        }
                    },
                    "400": {
                        "description": "Possible code responses: 4, 12.",
                        "schema": {
                            "$ref": "#/definitions/errs.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Possible code responses: 11.",
                        "schema": {
                            "$ref": "#/definitions/errs.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/pins/{pin_id}": {
            "get": {
                "description": "Get pin by id in the slug",
                "tags": [
                    "Pins"
                ],
                "summary": "Get pin by id",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Pin ID",
                        "name": "pin_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entity.PinPageResponse"
                        }
                    },
                    "400": {
                        "description": "Possible code responses: 4, 12.",
                        "schema": {
                            "$ref": "#/definitions/errs.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Possible code responses: 11.",
                        "schema": {
                            "$ref": "#/definitions/errs.ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "Update pin by description",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Pins"
                ],
                "summary": "Update pin",
                "parameters": [
                    {
                        "description": "Pin information",
                        "name": "pin",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/entity.Pin"
                        }
                    },
                    {
                        "type": "string",
                        "default": "session-token=",
                        "description": "session-token",
                        "name": "Cookie",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entity.PinPageResponse"
                        }
                    },
                    "400": {
                        "description": "Possible code responses: 3, 4, 12",
                        "schema": {
                            "$ref": "#/definitions/errs.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Possible code responses: 2.",
                        "schema": {
                            "$ref": "#/definitions/errs.ErrorResponse"
                        }
                    },
                    "403": {
                        "description": "Possible code responses: 14.",
                        "schema": {
                            "$ref": "#/definitions/errs.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Possible code responses: 11.",
                        "schema": {
                            "$ref": "#/definitions/errs.ErrorResponse"
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete pin by id (allowed only for pin's author)",
                "tags": [
                    "Pins"
                ],
                "summary": "Delete pin",
                "parameters": [
                    {
                        "type": "string",
                        "default": "session-token=",
                        "description": "session-token",
                        "name": "Cookie",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Pin ID",
                        "name": "pin_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entity.PinPageResponse"
                        }
                    },
                    "400": {
                        "description": "Possible code responses: 3, 4 12",
                        "schema": {
                            "$ref": "#/definitions/errs.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Possible code responses: 2.",
                        "schema": {
                            "$ref": "#/definitions/errs.ErrorResponse"
                        }
                    },
                    "403": {
                        "description": "Possible code responses: 14.",
                        "schema": {
                            "$ref": "#/definitions/errs.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Possible code responses: 11.",
                        "schema": {
                            "$ref": "#/definitions/errs.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/pins/{pin_id}/like": {
            "post": {
                "description": "Sets like by pin id and auth token",
                "tags": [
                    "Likes"
                ],
                "summary": "Set like on the pin",
                "parameters": [
                    {
                        "type": "string",
                        "default": "session-token=",
                        "description": "session-token",
                        "name": "Cookie",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Pin ID",
                        "name": "pin_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Possible code responses: 3, 12.",
                        "schema": {
                            "$ref": "#/definitions/errs.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Possible code responses: 2.",
                        "schema": {
                            "$ref": "#/definitions/errs.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Possible code responses: 11.",
                        "schema": {
                            "$ref": "#/definitions/errs.ErrorResponse"
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete like by pin id and auth token",
                "tags": [
                    "Likes"
                ],
                "summary": "Delete like on the pin",
                "parameters": [
                    {
                        "type": "string",
                        "default": "session-token=",
                        "description": "session-token",
                        "name": "Cookie",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Pin ID",
                        "name": "pin_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Possible code responses: 3, 12.",
                        "schema": {
                            "$ref": "#/definitions/errs.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Possible code responses: 2.",
                        "schema": {
                            "$ref": "#/definitions/errs.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Possible code responses: 11.",
                        "schema": {
                            "$ref": "#/definitions/errs.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/users": {
            "post": {
                "description": "Register user by POST request and add them to DB",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Authorization"
                ],
                "summary": "Register user",
                "parameters": [
                    {
                        "description": "json",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/entity.User"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entity.UserResponse"
                        }
                    },
                    "400": {
                        "description": "Possible code responses: 3, 4, 5.",
                        "schema": {
                            "$ref": "#/definitions/errs.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Possible code responses: 7, 8.",
                        "schema": {
                            "$ref": "#/definitions/errs.ErrorResponse"
                        }
                    },
                    "403": {
                        "description": "Possible code responses: 1.",
                        "schema": {
                            "$ref": "#/definitions/errs.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Possible code responses: 11.",
                        "schema": {
                            "$ref": "#/definitions/errs.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/users/{user_id}": {
            "post": {
                "description": "Update user by description and user id.",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Update user",
                "parameters": [
                    {
                        "type": "string",
                        "default": "session-token=",
                        "description": "session-token",
                        "name": "Cookie",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "User information in json",
                        "name": "user",
                        "in": "formData"
                    },
                    {
                        "type": "file",
                        "description": "User avatar",
                        "name": "image",
                        "in": "formData"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entity.PinPageResponse"
                        }
                    },
                    "400": {
                        "description": "Possible code responses: 3, 4, 5, 12, 13, 18",
                        "schema": {
                            "$ref": "#/definitions/errs.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Possible code responses: 2.",
                        "schema": {
                            "$ref": "#/definitions/errs.ErrorResponse"
                        }
                    },
                    "403": {
                        "description": "Possible code responses: 14.",
                        "schema": {
                            "$ref": "#/definitions/errs.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Possible code responses: 6, 11.",
                        "schema": {
                            "$ref": "#/definitions/errs.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "entity.FeedPinResponse": {
            "description": "PinResponse information with author, pin id and content URL.",
            "type": "object",
            "properties": {
                "author": {
                    "$ref": "#/definitions/entity.PinAuthor"
                },
                "content_url": {
                    "type": "string"
                },
                "pin_id": {
                    "type": "integer"
                }
            }
        },
        "entity.FeedPins": {
            "description": "Pins array of FeedPinResponse",
            "type": "object",
            "properties": {
                "pins": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/entity.FeedPinResponse"
                    }
                }
            }
        },
        "entity.Pin": {
            "description": "Full pin information",
            "type": "object",
            "properties": {
                "allow_comments": {
                    "type": "boolean"
                },
                "click_url": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "entity.PinAuthor": {
            "description": "User-author information with user id, nickname and avatar",
            "type": "object",
            "properties": {
                "avatar_url": {
                    "type": "string"
                },
                "nickname": {
                    "type": "string"
                },
                "user_id": {
                    "type": "integer"
                }
            }
        },
        "entity.PinPageResponse": {
            "description": "Full pin information",
            "type": "object",
            "properties": {
                "allow_comments": {
                    "type": "boolean"
                },
                "author": {
                    "$ref": "#/definitions/entity.PinAuthor"
                },
                "click_url": {
                    "type": "string"
                },
                "content_url": {
                    "type": "string"
                },
                "created_at": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "is_owner": {
                    "type": "boolean"
                },
                "likes_count": {
                    "type": "integer"
                },
                "pin_id": {
                    "type": "integer"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "entity.User": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "nickname": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "entity.UserList": {
            "description": "User information with user id, email, nickname and avatar_url",
            "type": "object",
            "properties": {
                "users": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/entity.UserResponse"
                    }
                }
            }
        },
        "entity.UserResponse": {
            "description": "User information with user id, email, nickname and avatar_url",
            "type": "object",
            "properties": {
                "avatar_url": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "nickname": {
                    "type": "string"
                },
                "user_id": {
                    "type": "integer"
                }
            }
        },
        "errs.ErrorResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "message": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "https://harmoniums.ru",
	BasePath:         "api/v1",
	Schemes:          []string{},
	Title:            "Harmonium backend API",
	Description:      "This is API-docs of backend server of Harmonica team.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
