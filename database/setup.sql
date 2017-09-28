-- MySQL Workbench Forward Engineering

SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='TRADITIONAL,ALLOW_INVALID_DATES';

-- -----------------------------------------------------
-- Schema app_mvp_dating
-- -----------------------------------------------------

-- -----------------------------------------------------
-- Schema app_mvp_dating
-- -----------------------------------------------------
CREATE SCHEMA IF NOT EXISTS `app_mvp_dating` DEFAULT CHARACTER SET utf8 ;
USE `app_mvp_dating` ;

-- -----------------------------------------------------
-- Table `app_mvp_dating`.`conversation`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `app_mvp_dating`.`conversation` (
  `create_time` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` TIMESTAMP NULL DEFAULT NULL,
  `conversation_id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id_one` INT(10) UNSIGNED NULL DEFAULT NULL,
  `user_id_two` INT(10) UNSIGNED NULL DEFAULT NULL,
  `ip` VARCHAR(30) NULL DEFAULT NULL COMMENT 'ip address',
  `deleted` BIT(1) NULL DEFAULT b'0',
  PRIMARY KEY (`conversation_id`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8;


-- -----------------------------------------------------
-- Table `app_mvp_dating`.`user`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `app_mvp_dating`.`user` (
  `user_id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `email` VARCHAR(255) NOT NULL,
  `password` VARCHAR(32) NOT NULL,
  `fullname` VARCHAR(45) NOT NULL,
  `nickname` VARCHAR(16) NULL DEFAULT NULL,
  `avatar_uri` VARCHAR(45) NULL DEFAULT NULL,
  `phone_number` VARCHAR(12) NULL DEFAULT NULL,
  `gender` ENUM('X', 'M', 'F') NULL DEFAULT 'X' COMMENT 'X => secret\nM => Male\nF => Female',
  `date_of_birth` DATE NULL DEFAULT '1900-01-01',
  `job` VARCHAR(45) NULL DEFAULT NULL,
  `living_at` VARCHAR(45) NULL DEFAULT NULL,
  `home_town` VARCHAR(45) NULL DEFAULT NULL,
  `working_at` VARCHAR(45) NULL DEFAULT NULL,
  `short_about` VARCHAR(255) NULL DEFAULT NULL,
  `height` TINYINT(3) UNSIGNED NULL DEFAULT NULL COMMENT 'unit: centimet\nrange: [139, 219]',
  `weight` TINYINT(3) UNSIGNED NULL DEFAULT NULL COMMENT 'unit: kilogram\nrange: [39, 149]',
  `status` ENUM('single', 'divorce', 'complicate') NULL DEFAULT NULL,
  `deleted` BIT(1) NULL DEFAULT b'0' COMMENT 'Whether this user is deleted.',
  `create_time` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` TIMESTAMP NULL DEFAULT NULL,
  PRIMARY KEY (`user_id`))
ENGINE = InnoDB
AUTO_INCREMENT = 4
DEFAULT CHARACTER SET = utf8;


-- -----------------------------------------------------
-- Table `app_mvp_dating`.`conversation_message`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `app_mvp_dating`.`conversation_message` (
  `create_time` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` TIMESTAMP NULL DEFAULT NULL,
  `conversation_message_id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` INT(10) UNSIGNED NOT NULL,
  `conversation_id` INT(10) UNSIGNED NOT NULL,
  `message` TEXT NULL DEFAULT NULL,
  `ip` VARCHAR(30) NULL DEFAULT NULL,
  PRIMARY KEY (`conversation_message_id`),
  INDEX `fk_conversation_message_1_idx` (`conversation_id` ASC),
  INDEX `fk_conversation_message_2_idx` (`user_id` ASC),
  CONSTRAINT `fk_conversation_message_1`
    FOREIGN KEY (`conversation_id`)
    REFERENCES `app_mvp_dating`.`conversation` (`conversation_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_conversation_message_2`
    FOREIGN KEY (`user_id`)
    REFERENCES `app_mvp_dating`.`user` (`user_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8;


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
-- Table `app_mvp_dating`.`subject`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `app_mvp_dating`.`subject` (
  `create_time` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` TIMESTAMP NULL DEFAULT NULL,
  `subject_id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `subject_alias` VARCHAR(45) NOT NULL,
  `deleted` BIT(1) NULL DEFAULT b'0',
  PRIMARY KEY (`subject_id`))
ENGINE = InnoDB
AUTO_INCREMENT = 13
DEFAULT CHARACTER SET = utf8;


-- -----------------------------------------------------
-- Table `app_mvp_dating`.`user_acts`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `app_mvp_dating`.`user_acts` (
  `create_time` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` TIMESTAMP NULL DEFAULT NULL,
  `user_id` INT(10) UNSIGNED NOT NULL,
  `target_user_id` INT(10) UNSIGNED NOT NULL,
  `likes` BIT(1) NULL DEFAULT b'0',
  `visit_date` DATE NULL DEFAULT NULL,
  `seen_count` TINYINT(3) NULL DEFAULT '0' COMMENT 'The number of times the user_id has seen the target user',
  `deleted` BIT(1) NULL DEFAULT b'0',
  PRIMARY KEY (`user_id`, `target_user_id`),
  INDEX `fk_user_likes_2_idx` (`target_user_id` ASC),
  CONSTRAINT `fk_user_likes_1`
    FOREIGN KEY (`user_id`)
    REFERENCES `app_mvp_dating`.`user` (`user_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_user_likes_2`
    FOREIGN KEY (`target_user_id`)
    REFERENCES `app_mvp_dating`.`user` (`user_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB
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
  INDEX `fk_user_answer_question_1_idx` (`user_id` ASC),
  INDEX `fk_user_answer_question_2_idx` (`question_id` ASC),
  CONSTRAINT `fk_user_answer_question_1`
    FOREIGN KEY (`user_id`)
    REFERENCES `app_mvp_dating`.`user` (`user_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_user_answer_question_2`
    FOREIGN KEY (`question_id`)
    REFERENCES `app_mvp_dating`.`question` (`question_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8;


-- -----------------------------------------------------
-- Table `app_mvp_dating`.`user_has_interests`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `app_mvp_dating`.`user_has_interests` (
  `create_time` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` TIMESTAMP NULL DEFAULT NULL,
  `deleted` BIT(1) NULL DEFAULT b'0',
  `user_id` INT(10) UNSIGNED NOT NULL,
  `interest_id` INT(10) UNSIGNED NOT NULL,
  PRIMARY KEY (`user_id`, `interest_id`),
  INDEX `fk_user_has_interests_1_idx` (`user_id` ASC),
  INDEX `fk_user_has_interests_2_idx` (`interest_id` ASC),
  CONSTRAINT `fk_user_has_interests_1`
    FOREIGN KEY (`user_id`)
    REFERENCES `app_mvp_dating`.`user` (`user_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_user_has_interests_2`
    FOREIGN KEY (`interest_id`)
    REFERENCES `app_mvp_dating`.`subject` (`subject_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8;


-- -----------------------------------------------------
-- Table `app_mvp_dating`.`user_has_strongness`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `app_mvp_dating`.`user_has_strongness` (
  `create_time` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` TIMESTAMP NULL DEFAULT NULL,
  `user_id` INT(10) UNSIGNED NOT NULL,
  `subject_id` INT(10) UNSIGNED NOT NULL,
  `deleted` BIT(1) NULL DEFAULT b'0',
  PRIMARY KEY (`user_id`, `subject_id`),
  INDEX `fk_user_has_strongness_2_idx` (`subject_id` ASC),
  CONSTRAINT `fk_user_has_strongness_1`
    FOREIGN KEY (`user_id`)
    REFERENCES `app_mvp_dating`.`user` (`user_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_user_has_strongness_2`
    FOREIGN KEY (`subject_id`)
    REFERENCES `app_mvp_dating`.`subject` (`subject_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8;


-- -----------------------------------------------------
-- Table `app_mvp_dating`.`user_has_weakness`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `app_mvp_dating`.`user_has_weakness` (
  `create_time` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` TIMESTAMP NULL DEFAULT NULL,
  `user_id` INT(10) UNSIGNED NOT NULL,
  `subject_id` INT(10) UNSIGNED NOT NULL,
  `deleted` BIT(1) NULL DEFAULT b'0',
  PRIMARY KEY (`user_id`, `subject_id`),
  INDEX `fk_user_has_weakness_2_idx` (`subject_id` ASC),
  CONSTRAINT `fk_user_has_weakness_1`
    FOREIGN KEY (`user_id`)
    REFERENCES `app_mvp_dating`.`user` (`user_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_user_has_weakness_2`
    FOREIGN KEY (`subject_id`)
    REFERENCES `app_mvp_dating`.`subject` (`subject_id`)
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

USE `app_mvp_dating` ;

-- -----------------------------------------------------
-- procedure conversation_check
-- -----------------------------------------------------

DELIMITER $$
USE `app_mvp_dating`$$
CREATE DEFINER=`root`@`localhost` PROCEDURE `conversation_check`(pCId1 INT, pCId2 INT)
BEGIN
	SELECT `conversation_id`
	FROM `conversation`
	WHERE
	(`user_id_one` = pUserId1 AND `user_id_two` = pUserId2)
	OR
	(`user_id_one` = pUserId2 AND `user_id_two` = pUserId1);
END$$

DELIMITER ;

-- -----------------------------------------------------
-- procedure get_conversation_list_for_user_id
-- -----------------------------------------------------

DELIMITER $$
USE `app_mvp_dating`$$
CREATE DEFINER=`root`@`localhost` PROCEDURE `get_conversation_list_for_user_id`(pUserId INT)
BEGIN
	SELECT U.user_id, C.conversation_id, U.fullname, U.email
	FROM user U, conversation C, conversation_message R
	WHERE 
	CASE

	WHEN C.user_id_one = pUserId
	THEN C.user_id_two = U.user_id
	WHEN C.user_id_two = pUserId
	THEN C.user_id_one= U.user_id
	END

	AND
	C.conversation_id = R.conversation_id
	AND
	(C.user_id_one = pUserId OR C.user_id_two = pUserId)
    ORDER BY C.conversation_id DESC;
END$$

DELIMITER ;

-- -----------------------------------------------------
-- procedure get_last_message_for_conversation_id
-- -----------------------------------------------------

DELIMITER $$
USE `app_mvp_dating`$$
CREATE DEFINER=`root`@`localhost` PROCEDURE `get_last_message_for_conversation_id`(pCId INT)
BEGIN
	SELECT `conversation_message_id`, `create_time`, `message`
	FROM `conversation_message` R
	WHERE `conversation_id` = pCId
	ORDER BY `conversation_message_id` DESC LIMIT 1;
END$$

DELIMITER ;

SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;
