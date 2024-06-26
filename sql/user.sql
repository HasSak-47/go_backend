-- @CREATE-USER
CREATE TABLE IF NOT EXISTS User(
	ID       INTEGER PRIMARY KEY AUTOINCREMENT,
	Name     VARCHAR(40) NOT NULL,
	Role     VARCHAR(40) NOT NULL,
	Password VARCHAR(32) NOT NULL
);

-- @ADD-USER
INSERT INTO User(Name, Role, Password) VALUES (?, ?, ?);

-- @VALIDATE-USER
SELECT * FROM User WHERE Name == ? and Password == ?;

-- @FIND-USERNAME
SELECT Name FROM User WHERE Name == ?;

-- @FIND-USER
SELECT * FROM User WHERE Name == ?;

-- @GET-USERS
SELECT * FROM User;
