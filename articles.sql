CREATE DATABASE IF NOT EXISTS `ptt`
USE `ptt`;

CREATE TABLE IF NOT EXISTS `articles` (
  `url` varchar(100) NOT NULL,
  `board` varchar(50) NOT NULL,
  `title` varchar(50) DEFAULT NULL,
  `author` varchar(50) DEFAULT NULL,
  `contains` longtext DEFAULT NULL,
  `reply` longtext DEFAULT NULL,
  `date` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`url`),
  UNIQUE KEY `url` (`url`),
  KEY `board` (`board`),
  KEY `author` (`author`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;