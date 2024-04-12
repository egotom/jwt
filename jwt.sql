-- --------------------------------------------------------
-- 主机:                           127.0.0.1
-- 服务器版本:                        11.1.0-MariaDB - mariadb.org binary distribution
-- 服务器操作系统:                      Win64
-- HeidiSQL 版本:                  12.3.0.6589
-- --------------------------------------------------------

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8 */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


-- 导出 jwtgo 的数据库结构
CREATE DATABASE IF NOT EXISTS `jwtgo` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci */;
USE `jwtgo`;

-- 导出  表 jwtgo.chat_room_users 结构
CREATE TABLE IF NOT EXISTS `chat_room_users` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `wxid` varchar(50) NOT NULL,
  `custom_account` varchar(100) DEFAULT NULL,
  `name` varchar(100) DEFAULT NULL,
  `nickname` varchar(100) DEFAULT NULL,
  `invite_at` longtext DEFAULT NULL,
  `room` longtext DEFAULT NULL,
  `is_menber` tinyint(1) DEFAULT 0,
  `is_friend` tinyint(1) DEFAULT 0,
  `created_at` datetime(3) DEFAULT current_timestamp(3),
  `updated_at` datetime(3) DEFAULT current_timestamp(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `wxid` (`wxid`),
  KEY `idx_chat_room_users_wxid` (`wxid`)
) ENGINE=InnoDB AUTO_INCREMENT=13935 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- 数据导出被取消选择。

-- 导出  表 jwtgo.contexts 结构
CREATE TABLE IF NOT EXISTS `contexts` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `vendor` tinyint(3) unsigned DEFAULT NULL,
  `wxid` varchar(50) DEFAULT NULL,
  `role` enum('model','user','system') DEFAULT NULL,
  `msg` longtext DEFAULT NULL,
  `created_at` datetime(3) DEFAULT current_timestamp(3),
  PRIMARY KEY (`id`),
  KEY `idx_contexts_wxid` (`wxid`)
) ENGINE=InnoDB AUTO_INCREMENT=3244 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- 数据导出被取消选择。

-- 导出  表 jwtgo.promotions 结构
CREATE TABLE IF NOT EXISTS `promotions` (
  `id` varchar(191) NOT NULL,
  `publisher` varchar(50) NOT NULL,
  `consumer` text DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  KEY `idx_promotions_publisher` (`publisher`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- 数据导出被取消选择。

-- 导出  表 jwtgo.tickets 结构
CREATE TABLE IF NOT EXISTS `tickets` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `amount` bigint(20) DEFAULT NULL,
  `wxid` varchar(50) DEFAULT NULL,
  `commit` text DEFAULT NULL,
  `created_at` datetime(3) DEFAULT current_timestamp(3),
  PRIMARY KEY (`id`),
  KEY `idx_tickets_wxid` (`wxid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- 数据导出被取消选择。

-- 导出  表 jwtgo.users 结构
CREATE TABLE IF NOT EXISTS `users` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `ticket` bigint(20) DEFAULT 5000,
  `is_verified` tinyint(1) DEFAULT 0,
  `is_active` tinyint(1) DEFAULT 0,
  `account` varchar(50) DEFAULT NULL,
  `avatar` varchar(200) DEFAULT NULL,
  `name` varchar(50) DEFAULT NULL,
  `email` varchar(50) DEFAULT NULL,
  `password` varchar(50) DEFAULT NULL,
  `wxid` varchar(50) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT current_timestamp(3),
  `updated_at` datetime(3) DEFAULT current_timestamp(3),
  PRIMARY KEY (`id`),
  KEY `idx_users_wxid` (`wxid`)
) ENGINE=InnoDB AUTO_INCREMENT=1996 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- 数据导出被取消选择。

/*!40103 SET TIME_ZONE=IFNULL(@OLD_TIME_ZONE, 'system') */;
/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IFNULL(@OLD_FOREIGN_KEY_CHECKS, 1) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40111 SET SQL_NOTES=IFNULL(@OLD_SQL_NOTES, 1) */;
