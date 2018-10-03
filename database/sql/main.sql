SET foreign_key_checks=0;

DROP TABLE IF EXISTS `user`;

CREATE TABLE `user` (
    `id` BIGINT unsigned NOT NULL,
    `name` VARCHAR(191) NOT NULL,
    `token` VARCHAR(191) NOT NULL,
    `created_at` DATETIME NOT NULL,
    `updated_at` DATETIME NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARACTER SET utf8mb4;

SET foreign_key_checks=1;