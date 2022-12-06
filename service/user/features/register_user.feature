Feature: Create new user

    In order to use Arjuna
    I need to register user

    Scenario: Invalid json request body (string)
        Given the user is empty
        When I register user with body
            | string |
        Then response status code must be 400
        And response must match json
            """
            {
                "code": 3,
                "message": "invalid character 's' looking for beginning of value",
                "details": []
            }
            """
    
    Scenario: Invalid json request body (integer)
        Given the user is empty
        When I register user with body
            | integer |
        Then response status code must be 400
        And response must match json
            """
            {
                "code": 3,
                "message": "invalid character 'i' looking for beginning of value",
                "details": []
            }
            """

    Scenario: Invalid name
        Given the user is empty
        When I register user with body
            | {"email": "a@a.com", "password": "a", "name": "a1"} |
        Then response status code must be 400
        And response must match json
            """
            {
                "code": 3,
                "message": "",
                "details": [
                    {
                        "@type": "type.googleapis.com/google.rpc.BadRequest",
                        "fieldViolations": [
                            {
                                "field": "name",
                                "description": "contain character outside of alphabet"
                            }
                        ]
                    },
                    {
                        "@type": "type.googleapis.com/api.v1.UserError",
                        "errorCode": "USER_ERROR_CODE_INVALID_NAME"
                    }
                ]
            }
            """

    Scenario: Invalid email
        Given the user is empty
        When I register user with body
            | {"email": "a@.com", "password": "a", "name": "a"} |
        Then response status code must be 400
        And response must match json
            """
            {
                "code": 3,
                "message": "",
                "details": [
                    {
                        "@type": "type.googleapis.com/google.rpc.BadRequest",
                        "fieldViolations": [
                            {
                                "field": "email",
                                "description": ""
                            }
                        ]
                    },
                    {
                        "@type": "type.googleapis.com/api.v1.UserError",
                        "errorCode": "USER_ERROR_CODE_INVALID_EMAIL"
                    }
                ]
            }
            """

    Scenario: Valid json request body
        Given the user is empty
        When I register user with body
            | {"email": "a@a.com", "password": "a", "name": "a"} |
            | {"email": "b@b.com", "password": "b", "name": "b"} |
            | {"email": "c@c.com", "password": "c", "name": "c"} |
        Then response status code must be 200

    Scenario: Create new users but already exists
        Given there are users with
            | {"email": "a@a.com", "password": "a", "name": "a"} |
            | {"email": "b@b.com", "password": "b", "name": "b"} |
            | {"email": "c@c.com", "password": "c", "name": "c"} |
        When I register user with body
            | {"email": "a@a.com", "password": "a", "name": "a"} |
            | {"email": "b@b.com", "password": "b", "name": "b"} |
            | {"email": "c@c.com", "password": "c", "name": "c"} |
        Then response status code must be 409
        And response must match json
            """
            {
                "code": 6,
                "message": "",
                "details": [
                    {
                        "@type": "type.googleapis.com/api.v1.UserError",
                        "errorCode": "USER_ERROR_CODE_ALREADY_EXISTS"
                    }
                ]
            }
            """
