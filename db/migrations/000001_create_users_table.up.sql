CREATE TABLE
    IF NOT EXISTS `users` (
        `id` TINYINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
        `email` VARCHAR(255) UNIQUE NOT NULL,
        `name` VARCHAR(255) NOT NULL,
        `password` VARCHAR(255) NOT NULL,
        `refresh_token` CHAR(64)
    )