SET foreign_key_checks=0;

DROP TABLE IF EXISTS `user`;

CREATE TABLE `user` (
    `id` BIGINT unsigned NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(191) NOT NULL,
    `token` VARCHAR(191) NOT NULL,
    `last_login_at` DATETIME NOT NULL,
    `rank` INTEGER unsigned NOT NULL DEFAULT 1,
    `money` BIGINT NOT NULL DEFAULT 0,
    `stamina` INTEGER unsigned NOT NULL DEFAULT 0,
    `created_at` DATETIME NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARACTER SET utf8mb4;


DROP TABLE IF EXISTS `user_business`;

CREATE TABLE `user_business` (
    `id` BIGINT unsigned NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT unsigned NOT NULL,
    `business_id` BIGINT unsigned NOT NULL,
    `level` INTEGER unsigned NOT NULL,
    `last_buy_at` DATETIME NOT NULL,
    UNIQUE `user_business_idx` (`user_id`, `business_id`),
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARACTER SET utf8mb4;


DROP TABLE IF EXISTS `user_rank`;

CREATE TABLE `user_rank` (
    `id` INTEGER unsigned NOT NULL AUTO_INCREMENT,
    `rank` INTEGER unsigned NOT NULL,
    `assets` BIGINT unsigned NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARACTER SET utf8mb4;


DROP TABLE IF EXISTS `business`;

CREATE TABLE `business` (
    `id` INTEGER unsigned NOT NULL,
    `prefecture` INTEGER unsigned NOT NULL,
    `name` VARCHAR(191) NOT NULL,
    `price_base` BIGINT unsigned NOT NULL,
    `price_level2` BIGINT unsigned NOT NULL,
    `price_level3` BIGINT unsigned NOT NULL,
    `return_rate_base` INTEGER unsigned NOT NULL,
    `return_rate_level2` INTEGER unsigned NOT NULL,
    `return_rate_level3` INTEGER unsigned NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARACTER SET utf8mb4;

SET foreign_key_checks=1;