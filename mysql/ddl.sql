CREATE DATABASE IF NOT EXISTS  minimum_wages DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci;

CREATE TABLE IF NOT EXISTS `checksums` (
  `name` varchar(64) NOT NULL,
  `checksum` text,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `wage_logs` (
  `id` char(26) NOT NULL COMMENT '最低賃金ID',
  `prefecture_id` int(11) NOT NULL COMMENT '都道府県ID',
  `hourly` int(11) NOT NULL COMMENT '最低時給額',
  `daily` int(11) DEFAULT NULL COMMENT '最低日給額',
  `name` varchar(1000) NOT NULL COMMENT '対象産業名',
  `regional` tinyint(1) NOT NULL DEFAULT '0' COMMENT '地域最低賃金（0:特定最低賃金, 1:地域最低賃金）',
  `implemented_at` timestamp NOT NULL COMMENT '施行日',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '作成日時',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新日時',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '削除日時',
  PRIMARY KEY (`id`,`implemented_at`),
  KEY `idx_id` (`id`),
  KEY `idx_prefecture_id` (`prefecture_id`),
  KEY `idx_regional` (`regional`),
  KEY `idx_implementation_at` (`implemented_at`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_updated_at` (`updated_at`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='logs for japanese minimum wages';

CREATE TABLE IF NOT EXISTS `prefectures` (
  `id` int(11) NOT NULL,
  `name` varchar(64) NOT NULL DEFAULT '',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='prefectures';

INSERT INTO `prefectures` (`id`, `name`)
VALUES
	(1,'北海道'),
	(2,'青森県'),
	(3,'岩手県'),
	(4,'宮城県'),
	(5,'秋田県'),
	(6,'山形県'),
	(7,'福島県'),
	(8,'茨城県'),
	(9,'栃木県'),
	(10,'群馬県'),
	(11,'埼玉県'),
	(12,'千葉県'),
	(13,'東京都'),
	(14,'神奈川県'),
	(15,'新潟県'),
	(16,'富山県'),
	(17,'石川県'),
	(18,'福井県'),
	(19,'山梨県'),
	(20,'長野県'),
	(21,'岐阜県'),
	(22,'静岡県'),
	(23,'愛知県'),
	(24,'三重県'),
	(25,'滋賀県'),
	(26,'京都府'),
	(27,'大阪府'),
	(28,'兵庫県'),
	(29,'奈良県'),
	(30,'和歌山県'),
	(31,'鳥取県'),
	(32,'島根県'),
	(33,'岡山県'),
	(34,'広島県'),
	(35,'山口県'),
	(36,'徳島県'),
	(37,'香川県'),
	(38,'愛媛県'),
	(39,'高知県'),
	(40,'福岡県'),
	(41,'佐賀県'),
	(42,'長崎県'),
	(43,'熊本県'),
	(44,'大分県'),
	(45,'宮崎県'),
	(46,'鹿児島県'),
	(47,'沖縄県');