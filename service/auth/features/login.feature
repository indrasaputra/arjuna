Feature: Login

    In order to use Arjuna
    I need to login

    Scenario: Wrong client-id can't login
        Given there are users with
            | {"email": "a@a.com", "password": "a", "name": "a"} |
        When I login with user
            | {"email": "a@a.com", "password": "a", "clientId": "unknown-client"} |
        Then response status code must be 400
        And response must match json
            """
            {
                "code": 3,
                "message": "Invalid client credentials",
                "details": [
                    {
                        "@type": "type.googleapis.com/api.v1.AuthError",
                        "errorCode": "AUTH_ERROR_CODE_INVALID_ARGUMENT"
                    }
                ]
            }
            """

    Scenario: Non-exist user can't login
        Given there are users with
            | {"email": "a@a.com", "password": "a", "name": "a"} |
        When I login with user
            | {"email": "b@b.com", "password": "b", "clientId": "arjuna-client"} |
        Then response status code must be 401
        And response must match json
            """
            {
                "code": 16,
                "message": "",
                "details": [
                    {
                        "@type": "type.googleapis.com/api.v1.AuthError",
                        "errorCode": "AUTH_ERROR_CODE_UNAUTHORIZED"
                    }
                ]
            }
            """

    Scenario: Wrong password can't login
        Given there are users with
            | {"email": "a@a.com", "password": "a", "name": "a"} |
        When I login with user
            | {"email": "a@a.com", "password": "abc", "clientId": "arjuna-client"} |
        Then response status code must be 401
        And response must match json
            """
            {
                "code": 16,
                "message": "",
                "details": [
                    {
                        "@type": "type.googleapis.com/api.v1.AuthError",
                        "errorCode": "AUTH_ERROR_CODE_UNAUTHORIZED"
                    }
                ]
            }
            """

    Scenario: Registered user can login using the right credential
        Given there are users with
            | {"email": "a@a.com", "password": "a", "name": "a"} |
        When I login with user
            | {"email": "a@a.com", "password": "a", "clientId": "arjuna-client"} |
        Then response status code must be 200
