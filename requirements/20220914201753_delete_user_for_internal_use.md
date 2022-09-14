## A user can be deleted


### Delete non-exist user

#### Given

- An internal system wants to delete user.
- An internal system provides user's email.
- User doesn't exist in system.

#### When

- An internal system deletes user.

#### Then

- System returns user not found error.


### Delete existing user

#### Given

- An internal system wants to delete user.
- An internal system provides user's email.
- User exists in system.

#### When

- An internal system deletes user.

#### Then

- System hard-deletes the user.
