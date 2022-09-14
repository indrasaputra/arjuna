## A system can get all users


### Retrieve users without limit

#### Given

- Default limit is 10.
- There are infinite users in system.
- An internal system wants to retrieve all users.
- The request doesn't define limit.

#### When

- An internal system call get all users endpoint.

#### Then

- System must return 10 users.


### Retrieve users with limit less than default limit

#### Given

- Default limit is 10.
- There are infinite users in system.
- An internal system wants to retrieve all users.
- The request defines limit as 2.

#### When

- An internal system call get all users endpoint.

#### Then

- System must return just 2 users.


### Retrieve users with limit more than default limit

#### Given

- Default limit is 10.
- There are infinite users in system.
- An internal system wants to retrieve all users.
- The request defines limit as 20.

#### When

- An internal system call get all users endpoint.

#### Then

- System must return 10 users.
