## A user can register to Arjuna

### Register already exist user
#### Given

- A user wants to register him/herself to Arjuna.
- A user must provide these data:
    - Full name
    - Email
    - Password
- The user already exists in the system.

#### When

- The user agrees to register.

#### Then

- System returns already exist error.


### Register non-exist user
#### Given

- A user wants to register him/herself to system.
- A user must provide these data:
    - Full name
    - Email
    - Password
- The user doesn't exist in the system.

#### When

- The user agrees to register.

#### Then

- A record of user is being stored in system.
