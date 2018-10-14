SET foreign_key_checks=0;

DROP TABLE IF EXISTS `user`;

CREATE TABLE `user` (
    `id` BIGINT unsigned NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(191) NOT NULL,
    `token` VARCHAR(191) NOT NULL,
    `last_login_at` DATETIME NOT NULL,
    `money` BIGINT unsigned NOT NULL DEFAULT 0,
    `stamina` INTEGER unsigned NOT NULL DEFAULT 0,
    `created_at` DATETIME NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARACTER SET utf8mb4;

SET foreign_key_checks=1;