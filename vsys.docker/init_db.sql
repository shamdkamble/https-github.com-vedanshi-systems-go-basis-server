CREATE DATABASE IF NOT EXISTS `vsys-db`;
USE `vsys-db`;

CREATE TABLE IF NOT EXISTS `end_user` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `first_name` VARCHAR(255),
    `last_name` VARCHAR(255),
    `email` VARCHAR(255),
    `device_type` VARCHAR(255),
    `model_number` VARCHAR(255),
    `mobile_no` BIGINT UNSIGNED,
    `password` VARCHAR(255),
    `profile_pic_url` TEXT,
    `social_ident` VARCHAR(255),
    `status` VARCHAR(255),
    `activated` BIT(1),
    `reset_key` INT,
    `reset_date` DATETIME,
    `created_on` DATETIME,
    `created_by` BIGINT UNSIGNED,
    `modified_on` DATETIME,
    `modified_by` BIGINT UNSIGNED,
    PRIMARY KEY (`id`)
);

CREATE TABLE IF NOT EXISTS `end_user_role` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `role` VARCHAR(255),
    `status` VARCHAR(255),
    `created_on` DATETIME,
    `created_by` BIGINT UNSIGNED,
    `modified_on` DATETIME,
    `modified_by` BIGINT UNSIGNED,
    `end_user_id` BIGINT UNSIGNED,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`end_user_id`) REFERENCES `end_user`(`id`)
);

CREATE TABLE IF NOT EXISTS `otp_details` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `mobile_no` BIGINT UNSIGNED,
    `otp` INT,
    `created_on` DATETIME,
    `modified_on` DATETIME,
    PRIMARY KEY (`id`)
);
