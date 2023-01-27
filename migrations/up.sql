CREATE TABLE IF NOT EXISTS user(
		userId  INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE,
		password text, 	
		email text,
		posts INT DEFAULT 0,
		creationDate DATE DEFAULT (datetime('now')),
		token TEXT DEFAULT NULL,
		expiresAt DATETIME DEFAULT NULL 
);

CREATE TABLE IF NOT EXISTS posts(
		postId INTEGER PRIMARY KEY AUTOINCREMENT,
		author text ,
		title text,
		content text,
		creationDate DATE DEFAULT (datetime('now')),
    	likes INT DEFAULT 0,
    	dislikes INT DEFAULT 0,
		FOREIGN KEY (author) REFERENCES user(username)
);

CREATE TABLE IF NOT EXISTS posts_category(
    	postCategoryId INTEGER,
		category TEXT,
		FOREIGN KEY (postCategoryId) REFERENCES posts(postId) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS comments(
    	commentsId INTEGER PRIMARY KEY AUTOINCREMENT,
    	postId INTEGER,
    	author TEXT,
    	content TEXT,
    	likes INT DEFAULT 0,
    	dislikes INT DEFAULT 0,
    	FOREIGN KEY (postId)  REFERENCES posts(postId)
);

CREATE TABLE IF NOT EXISTS likes(
    	likeId INTEGER PRIMARY KEY AUTOINCREMENT,
    	username TEXT,
    	postId INTEGER DEFAULT NULL,
    	commentsId INTEGER DEFAULT NULL,
    	FOREIGN KEY (postId) REFERENCES posts(postId) ON DELETE CASCADE,
    	FOREIGN KEY (commentsId) REFERENCES comments(commentsId) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS dislikes(
    	dislikeId INTEGER PRIMARY KEY AUTOINCREMENT,
    	username TEXT,
    	postId INTEGER DEFAULT NULL,
    	commentsId INTEGER DEFAULT NULL,
    	FOREIGN KEY (postId) REFERENCES posts(postId) ON DELETE CASCADE,
		FOREIGN KEY (commentsId) REFERENCES comments(commentsId) ON DELETE CASCADE
);