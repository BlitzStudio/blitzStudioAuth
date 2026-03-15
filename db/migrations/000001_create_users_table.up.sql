CREATE TABLE
    IF NOT EXISTS `users` (
        `id` INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
        `email` VARCHAR(255) UNIQUE NOT NULL,
        `name` VARCHAR(255) NOT NULL,
        `password` VARCHAR(255) NOT NULL
    );

CREATE TABLE
    IF NOT EXISTS `jwts` (
        `id` CHAR(26) PRIMARY KEY,
        `userId` INT UNSIGNED,
        -- `deviceId` CHAR(36) NOT NULL,
        `tokenFamily` CHAR(36) NOT NULL,
        -- `tokenHash` VARCHAR(255),
        `expiresAt` DATETIME,
        `isRevoked` BOOLEAN DEFAULT FALSE,
        `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        CONSTRAINT `fk_jwts_user` FOREIGN KEY (`userId`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
    );