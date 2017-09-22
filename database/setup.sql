-- -----------------------------------------------------
-- Schema app_mvp_dating
-- -----------------------------------------------------
CREATE SCHEMA IF NOT EXISTS `app_mvp_dating` DEFAULT CHARACTER SET utf8 ;
USE `app_mvp_dating` ;

-- -----------------------------------------------------
-- Table `app_mvp_dating`.`question`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `app_mvp_dating`.`question` (
  `question_id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `content` VARCHAR(255) NOT NULL,
  `create_time` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` TIMESTAMP NULL DEFAULT NULL,
  `deleted` BIT(1) NULL DEFAULT b'0',
  PRIMARY KEY (`question_id`))
ENGINE = InnoDB
AUTO_INCREMENT = 3
DEFAULT CHARACTER SET = utf8;


-- -----------------------------------------------------
-- Table `app_mvp_dating`.`user`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `app_mvp_dating`.`user` (
  `user_id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `email` VARCHAR(255) NOT NULL,
  `password` VARCHAR(32) NOT NULL,
  `username` VARCHAR(16) NULL DEFAULT NULL,
  `fullname` VARCHAR(45) NOT NULL,
  `avatar_uri` VARCHAR(45) NULL DEFAULT NULL,
  `phone_number` VARCHAR(12) NULL DEFAULT NULL,
  `gender` ENUM('secret', 'male', 'female') NULL DEFAULT 'secret',
  `date_of_birth` DATE NULL DEFAULT '1900-01-01',
  `live_city` VARCHAR(45) NULL DEFAULT NULL,
  `deleted` BIT(1) NULL DEFAULT b'0' COMMENT 'Whether this user is deleted.',
  `create_time` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` TIMESTAMP NULL DEFAULT NULL,
  PRIMARY KEY (`user_id`))
ENGINE = InnoDB
AUTO_INCREMENT = 4
DEFAULT CHARACTER SET = utf8;


-- -----------------------------------------------------
-- Table `app_mvp_dating`.`user_answer_question`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `app_mvp_dating`.`user_answer_question` (
  `user_id` INT(10) UNSIGNED NOT NULL,
  `question_id` INT(10) UNSIGNED NOT NULL,
  `self_answer` VARCHAR(255) NULL DEFAULT NULL COMMENT 'Answer which the user_id provides to others who ask this kind of question.',
  `expect_answer` VARCHAR(255) NULL DEFAULT NULL COMMENT 'Answer which user_id expects from others for this kind of question_id',
  `hate_answer` VARCHAR(255) NULL DEFAULT NULL COMMENT 'Anwer which the user_id doesn\'t expect from others for this kind of question_id',
  `also_ask` BIT(1) NULL DEFAULT b'0' COMMENT 'If true, the question_id is also a question which the user_id cares and wants to ask others,\nelse the question_id is just referring to a question which the user_id has answered.',
  `deleted` BIT(1) NULL DEFAULT b'0',
  `create_time` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` TIMESTAMP NULL DEFAULT NULL,
  PRIMARY KEY (`user_id`, `question_id`),
  INDEX `fk_answer_question_idx` (`question_id` ASC),
  CONSTRAINT `fk_answer_question`
    FOREIGN KEY (`question_id`)
    REFERENCES `app_mvp_dating`.`question` (`question_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_user_answer`
    FOREIGN KEY (`user_id`)
    REFERENCES `app_mvp_dating`.`user` (`user_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8;


-- -----------------------------------------------------
-- Table `app_mvp_dating`.`user_photo`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `app_mvp_dating`.`user_photo` (
  `photo_id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `photo_uri` VARCHAR(255) NOT NULL,
  `user_id` INT(10) UNSIGNED NOT NULL,
  `deleted` BIT(1) NULL DEFAULT b'0',
  `create_time` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` TIMESTAMP NULL DEFAULT NULL,
  PRIMARY KEY (`photo_id`),
  INDEX `fk_user_photo` (`user_id` ASC),
  CONSTRAINT `fk_user_photo`
    FOREIGN KEY (`user_id`)
    REFERENCES `app_mvp_dating`.`user` (`user_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8;
