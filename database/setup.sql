-- -----------------------------------------------------
-- Schema app_mvp_dating
-- -----------------------------------------------------
CREATE SCHEMA IF NOT EXISTS `app_mvp_dating` DEFAULT CHARACTER SET utf8 ;
USE `app_mvp_dating` ;

-- -----------------------------------------------------
-- Table `app_mvp_dating`.`user`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `app_mvp_dating`.`user` (
  `user_id` INT(11) NOT NULL AUTO_INCREMENT,
  `email` VARCHAR(255) NOT NULL,
  `password` VARCHAR(32) NOT NULL,
  `username` VARCHAR(16) NULL DEFAULT NULL,
  `fullname` VARCHAR(45) NOT NULL,
  `avatar_uri` VARCHAR(45) NULL DEFAULT NULL,
  `phone_number` VARCHAR(12) NULL DEFAULT NULL,
  `gender` ENUM('secret', 'male', 'female') NULL DEFAULT 'secret',
  `date_of_birth` DATE NULL DEFAULT '1900-01-01',
  `live_city` VARCHAR(45) NULL DEFAULT NULL,
  `create_time` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` TIMESTAMP NULL DEFAULT NULL,
  PRIMARY KEY (`user_id`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8;
