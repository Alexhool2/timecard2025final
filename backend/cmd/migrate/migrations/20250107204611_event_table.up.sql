CREATE TABLE IF NOT EXISTS event (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `userId` INT UNSIGNED NOT NULL,
    `date` DATE NOT NULL,
    `start_time` DATETIME,
    `start_lunch` DATETIME,
    `end_lunch` DATETIME,
    `end_time` DATETIME,
    `createdAt` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`userId`) REFERENCES users(`id`) ON DELETE CASCADE
);