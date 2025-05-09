CREATE TABLE users (
    id BIGINT NOT NULL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    roles JSON NOT NULL,
    password_hash VARCHAR(60) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE posts (
    id BIGINT NOT NULL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    title VARCHAR(255) NOT NULL,
    description VARCHAR(300) NOT NULL,
    front_image VARCHAR(512) NOT NULL DEFAULT "",
    content_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX (user_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE comments (
    id         BIGINT PRIMARY KEY,
    user_id    BIGINT NOT NULL,
    post_id    BIGINT NOT NULL,
    parent_id  BIGINT DEFAULT NULL,
    content    VARCHAR(10000) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_comments_user FOREIGN KEY (user_id) 
        REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,

    CONSTRAINT fk_comments_post FOREIGN KEY (post_id) 
        REFERENCES posts(id) ON DELETE CASCADE ON UPDATE CASCADE,

    CONSTRAINT fk_comments_parent FOREIGN KEY (parent_id) 
        REFERENCES comments(id) ON DELETE CASCADE ON UPDATE CASCADE,

    INDEX idx_comments_post_id (post_id)
);



