CREATE TABLE `profiles` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `first_name` varchar(50) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `last_name` varchar(30) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `city_id` int(4) DEFAULT '1',
  `phone` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `info` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `photoID` int(11) DEFAULT NULL,
  `status` tinyint(4) DEFAULT NULL,
  PRIMARY KEY (`id`),
  FULLTEXT KEY `first_name` (`first_name`,`last_name`),
  FULLTEXT KEY `first_name_2` (`first_name`,`last_name`),
  FULLTEXT KEY `first_name_3` (`first_name`,`last_name`),
  FULLTEXT KEY `first_name_4` (`first_name`,`last_name`),
  FULLTEXT KEY `first_name_5` (`first_name`,`last_name`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8_general_ci COLLATE=utf8mb4_general_ci;