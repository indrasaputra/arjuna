Feature: Get user

    In order to use Arjuna
    I need to get user

    Scenario: No user registered to the system
        Given the user is empty
        When I get all users
        Then response status code must be 200
        And response must match json
            """
            {
                "data": []
            }
            """

    Scenario: Users exist
        Given there are users with
            | {"email": "a@a.com", "password": "a", "name": "a"} |
            | {"email": "b@b.com", "password": "b", "name": "b"} |
            | {"email": "c@c.com", "password": "c", "name": "c"} |
        When I get all users
        Then response status code must be 200
        And email must match
            | a@a.com |
            | b@b.com |
            | c@c.com |
        And number of users must be 3

    Scenario: 15 users exist and default max limit is 10 
        Given there are users with
            | {"email": "a01@a.com", "password": "a", "name": "a"} |
            | {"email": "a02@a.com", "password": "a", "name": "a"} |
            | {"email": "a03@a.com", "password": "a", "name": "a"} |
            | {"email": "a04@a.com", "password": "a", "name": "a"} |
            | {"email": "a05@a.com", "password": "a", "name": "a"} |
            | {"email": "a06@a.com", "password": "a", "name": "a"} |
            | {"email": "a07@a.com", "password": "a", "name": "a"} |
            | {"email": "a08@a.com", "password": "a", "name": "a"} |
            | {"email": "a09@a.com", "password": "a", "name": "a"} |
            | {"email": "a10@a.com", "password": "a", "name": "a"} |
            | {"email": "a11@a.com", "password": "a", "name": "a"} |
            | {"email": "a12@a.com", "password": "a", "name": "a"} |
            | {"email": "a13@a.com", "password": "a", "name": "a"} |
            | {"email": "a14@a.com", "password": "a", "name": "a"} |
            | {"email": "a15@a.com", "password": "a", "name": "a"} |
        When I get all users
        Then response status code must be 200
        And number of users must be 10
