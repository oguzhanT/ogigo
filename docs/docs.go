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
        "/": {
            "get": {
                "description": "This endpoint retrieves the current wallet state.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Wallet"
                ],
                "summary": "Retrieve the current wallet state",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.WalletResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/main.ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "This endpoint processes wallet update events and places them into a Kafka queue.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Events"
                ],
                "summary": "Process incoming wallet update events",
                "parameters": [
                    {
                        "description": "Events to be processed",
                        "name": "events",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.EventRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/main.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/main.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "main.Balance": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "number"
                },
                "currency": {
                    "type": "string"
                }
            }
        },
        "main.ErrorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        },
        "main.Event": {
            "type": "object",
            "properties": {
                "app": {
                    "type": "string"
                },
                "attributes": {
                    "type": "object",
                    "properties": {
                        "amount": {
                            "type": "string"
                        },
                        "currency": {
                            "type": "string"
                        }
                    }
                },
                "meta": {
                    "type": "object",
                    "properties": {
                        "user": {
                            "type": "string"
                        }
                    }
                },
                "time": {
                    "type": "string"
                },
                "type": {
                    "type": "string"
                },
                "wallet": {
                    "type": "string"
                }
            }
        },
        "main.EventRequest": {
            "type": "object",
            "properties": {
                "events": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/main.Event"
                    }
                }
            }
        },
        "main.Wallet": {
            "type": "object",
            "properties": {
                "balances": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/main.Balance"
                    }
                },
                "id": {
                    "type": "string"
                }
            }
        },
        "main.WalletResponse": {
            "type": "object",
            "properties": {
                "wallets": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/main.Wallet"
                    }
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:80",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Ogigo API",
	Description:      "This is a test scenario server.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
