CREATE DATABASE IF NOT EXISTS `torch-database`;
USE `torch-database`;

-- Country table
CREATE TABLE IF NOT EXISTS `countries` (
  `country_id` TINYINT UNSIGNED PRIMARY KEY,
  `country_code` char(2) NOT NULL,
  `country` varchar(60) NOT NULL
) 
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_unicode_ci
;

INSERT INTO `countries` (`country_id`, `country_code`, `country`) VALUES
(1, 'AF', 'Afghanistan'),
(2, 'AX', 'Aland Islands'),
(3, 'AL', 'Albania'),
(4, 'DZ', 'Algeria'),
(5, 'AS', 'American Samoa'),
(6, 'AD', 'Andorra'),
(7, 'AO', 'Angola'),
(8, 'AI', 'Anguilla'),
(9, 'AQ', 'Antarctica'),
(10, 'AG', 'Antigua and Barbuda'),
(11, 'AR', 'Argentina'),
(12, 'AM', 'Armenia'),
(13, 'AW', 'Aruba'),
(14, 'AU', 'Australia'),
(15, 'AT', 'Austria'),
(16, 'AZ', 'Azerbaijan'),
(17, 'BS', 'Bahamas'),
(18, 'BH', 'Bahrain'),
(19, 'BD', 'Bangladesh'),
(20, 'BB', 'Barbados'),
(21, 'BY', 'Belarus'),
(22, 'BE', 'Belgium'),
(23, 'BZ', 'Belize'),
(24, 'BJ', 'Benin'),
(25, 'BM', 'Bermuda'),
(26, 'BT', 'Bhutan'),
(27, 'BO', 'Bolivia'),
(28, 'BQ', 'Bonaire, Sint Eustatius and Saba'),
(29, 'BA', 'Bosnia and Herzegovina'),
(30, 'BW', 'Botswana'),
(31, 'BV', 'Bouvet Island'),
(32, 'BR', 'Brazil'),
(33, 'IO', 'British Indian Ocean Territory'),
(34, 'BN', 'Brunei Darussalam'),
(35, 'BG', 'Bulgaria'),
(36, 'BF', 'Burkina Faso'),
(37, 'BI', 'Burundi'),
(38, 'KH', 'Cambodia'),
(39, 'CM', 'Cameroon'),
(40, 'CA', 'Canada'),
(41, 'CV', 'Cape Verde'),
(42, 'KY', 'Cayman Islands'),
(43, 'CF', 'Central African Republic'),
(44, 'TD', 'Chad'),
(45, 'CL', 'Chile'),
(46, 'CN', 'China'),
(47, 'CX', 'Christmas Island'),
(48, 'CC', 'Cocos (Keeling) Islands'),
(49, 'CO', 'Colombia'),
(50, 'KM', 'Comoros'),
(51, 'CG', 'Congo'),
(52, 'CD', 'Congo, Democratic Republic of the Congo'),
(53, 'CK', 'Cook Islands'),
(54, 'CR', 'Costa Rica'),
(55, 'CI', "Cote D'Ivoire"),
(56, 'HR', 'Croatia'),
(57, 'CU', 'Cuba'),
(58, 'CW', 'Curacao'),
(59, 'CY', 'Cyprus'),
(60, 'CZ', 'Czech Republic'),
(61, 'DK', 'Denmark'),
(62, 'DJ', 'Djibouti'),
(63, 'DM', 'Dominica'),
(64, 'DO', 'Dominican Republic'),
(65, 'EC', 'Ecuador'),
(66, 'EG', 'Egypt'),
(67, 'SV', 'El Salvador'),
(68, 'GQ', 'Equatorial Guinea'),
(69, 'ER', 'Eritrea'),
(70, 'EE', 'Estonia'),
(71, 'ET', 'Ethiopia'),
(72, 'FK', 'Falkland Islands (Malvinas)'),
(73, 'FO', 'Faroe Islands'),
(74, 'FJ', 'Fiji'),
(75, 'FI', 'Finland'),
(76, 'FR', 'France'),
(77, 'GF', 'French Guiana'),
(78, 'PF', 'French Polynesia'),
(79, 'TF', 'French Southern Territories'),
(80, 'GA', 'Gabon'),
(81, 'GM', 'Gambia'),
(82, 'GE', 'Georgia'),
(83, 'DE', 'Germany'),
(84, 'GH', 'Ghana'),
(85, 'GI', 'Gibraltar'),
(86, 'GR', 'Greece'),
(87, 'GL', 'Greenland'),
(88, 'GD', 'Grenada'),
(89, 'GP', 'Guadeloupe'),
(90, 'GU', 'Guam'),
(91, 'GT', 'Guatemala'),
(92, 'GG', 'Guernsey'),
(93, 'GN', 'Guinea'),
(94, 'GW', 'Guinea-Bissau'),
(95, 'GY', 'Guyana'),
(96, 'HT', 'Haiti'),
(97, 'HM', 'Heard Island and Mcdonald Islands'),
(98, 'VA', 'Holy See (Vatican City State)'),
(99, 'HN', 'Honduras'),
(100, 'HK', 'Hong Kong'),
(101, 'HU', 'Hungary'),
(102, 'IS', 'Iceland'),
(103, 'IN', 'India'),
(104, 'ID', 'Indonesia'),
(105, 'IR', 'Iran, Islamic Republic of'),
(106, 'IQ', 'Iraq'),
(107, 'IE', 'Ireland'),
(108, 'IM', 'Isle of Man'),
(109, 'IL', 'Israel'),
(110, 'IT', 'Italy'),
(111, 'JM', 'Jamaica'),
(112, 'JP', 'Japan'),
(113, 'JE', 'Jersey'),
(114, 'JO', 'Jordan'),
(115, 'KZ', 'Kazakhstan'),
(116, 'KE', 'Kenya'),
(117, 'KI', 'Kiribati'),
(118, 'KP', "Korea, Democratic People's Republic of"),
(119, 'KR', 'Korea, Republic of'),
(120, 'XK', 'Kosovo'),
(121, 'KW', 'Kuwait'),
(122, 'KG', 'Kyrgyzstan'),
(123, 'LA', "Lao People's Democratic Republic"),
(124, 'LV', 'Latvia'),
(125, 'LB', 'Lebanon'),
(126, 'LS', 'Lesotho'),
(127, 'LR', 'Liberia'),
(128, 'LY', 'Libyan Arab Jamahiriya'),
(129, 'LI', 'Liechtenstein'),
(130, 'LT', 'Lithuania'),
(131, 'LU', 'Luxembourg'),
(132, 'MO', 'Macao'),
(133, 'MK', 'Macedonia, the Former Yugoslav Republic of'),
(134, 'MG', 'Madagascar'),
(135, 'MW', 'Malawi'),
(136, 'MY', 'Malaysia'),
(137, 'MV', 'Maldives'),
(138, 'ML', 'Mali'),
(139, 'MT', 'Malta'),
(140, 'MH', 'Marshall Islands'),
(141, 'MQ', 'Martinique'),
(142, 'MR', 'Mauritania'),
(143, 'MU', 'Mauritius'),
(144, 'YT', 'Mayotte'),
(145, 'MX', 'Mexico'),
(146, 'FM', 'Micronesia, Federated States of'),
(147, 'MD', 'Moldova, Republic of'),
(148, 'MC', 'Monaco'),
(149, 'MN', 'Mongolia'),
(150, 'ME', 'Montenegro'),
(151, 'MS', 'Montserrat'),
(152, 'MA', 'Morocco'),
(153, 'MZ', 'Mozambique'),
(154, 'MM', 'Myanmar'),
(155, 'NA', 'Namibia'),
(156, 'NR', 'Nauru'),
(157, 'NP', 'Nepal'),
(158, 'NL', 'Netherlands'),
(159, 'AN', 'Netherlands Antilles'),
(160, 'NC', 'New Caledonia'),
(161, 'NZ', 'New Zealand'),
(162, 'NI', 'Nicaragua'),
(163, 'NE', 'Niger'),
(164, 'NG', 'Nigeria'),
(165, 'NU', 'Niue'),
(166, 'NF', 'Norfolk Island'),
(167, 'MP', 'Northern Mariana Islands'),
(168, 'NO', 'Norway'),
(169, 'OM', 'Oman'),
(170, 'PK', 'Pakistan'),
(171, 'PW', 'Palau'),
(172, 'PS', 'Palestinian Territory, Occupied'),
(173, 'PA', 'Panama'),
(174, 'PG', 'Papua New Guinea'),
(175, 'PY', 'Paraguay'),
(176, 'PE', 'Peru'),
(177, 'PH', 'Philippines'),
(178, 'PN', 'Pitcairn'),
(179, 'PL', 'Poland'),
(180, 'PT', 'Portugal'),
(181, 'PR', 'Puerto Rico'),
(182, 'QA', 'Qatar'),
(183, 'RE', 'Reunion'),
(184, 'RO', 'Romania'),
(185, 'RU', 'Russian Federation'),
(186, 'RW', 'Rwanda'),
(187, 'BL', 'Saint Barthelemy'),
(188, 'SH', 'Saint Helena'),
(189, 'KN', 'Saint Kitts and Nevis'),
(190, 'LC', 'Saint Lucia'),
(191, 'MF', 'Saint Martin'),
(192, 'PM', 'Saint Pierre and Miquelon'),
(193, 'VC', 'Saint Vincent and the Grenadines'),
(194, 'WS', 'Samoa'),
(195, 'SM', 'San Marino'),
(196, 'ST', 'Sao Tome and Principe'),
(197, 'SA', 'Saudi Arabia'),
(198, 'SN', 'Senegal'),
(199, 'RS', 'Serbia'),
(200, 'CS', 'Serbia and Montenegro'),
(201, 'SC', 'Seychelles'),
(202, 'SL', 'Sierra Leone'),
(203, 'SG', 'Singapore'),
(204, 'SX', 'Sint Maarten'),
(205, 'SK', 'Slovakia'),
(206, 'SI', 'Slovenia'),
(207, 'SB', 'Solomon Islands'),
(208, 'SO', 'Somalia'),
(209, 'ZA', 'South Africa'),
(210, 'GS', 'South Georgia and the South Sandwich Islands'),
(211, 'SS', 'South Sudan'),
(212, 'ES', 'Spain'),
(213, 'LK', 'Sri Lanka'),
(214, 'SD', 'Sudan'),
(215, 'SR', 'Suriname'),
(216, 'SJ', 'Svalbard and Jan Mayen'),
(217, 'SZ', 'Swaziland'),
(218, 'SE', 'Sweden'),
(219, 'CH', 'Switzerland'),
(220, 'SY', 'Syrian Arab Republic'),
(221, 'TW', 'Taiwan, Province of China'),
(222, 'TJ', 'Tajikistan'),
(223, 'TZ', 'Tanzania, United Republic of'),
(224, 'TH', 'Thailand'),
(225, 'TL', 'Timor-Leste'),
(226, 'TG', 'Togo'),
(227, 'TK', 'Tokelau'),
(228, 'TO', 'Tonga'),
(229, 'TT', 'Trinidad and Tobago'),
(230, 'TN', 'Tunisia'),
(231, 'TR', 'Turkey'),
(232, 'TM', 'Turkmenistan'),
(233, 'TC', 'Turks and Caicos Islands'),
(234, 'TV', 'Tuvalu'),
(235, 'UG', 'Uganda'),
(236, 'UA', 'Ukraine'),
(237, 'AE', 'United Arab Emirates'),
(238, 'GB', 'United Kingdom'),
(239, 'US', 'United States'),
(240, 'UM', 'United States Minor Outlying Islands'),
(241, 'UY', 'Uruguay'),
(242, 'UZ', 'Uzbekistan'),
(243, 'VU', 'Vanuatu'),
(244, 'VE', 'Venezuela'),
(245, 'VN', 'Viet Nam'),
(246, 'VG', 'Virgin Islands, British'),
(247, 'VI', 'Virgin Islands, U.s.'),
(248, 'WF', 'Wallis and Futuna'),
(249, 'EH', 'Western Sahara'),
(250, 'YE', 'Yemen'),
(251, 'ZM', 'Zambia'),
(252, 'ZW', 'Zimbabwe');

-- Users table
CREATE TABLE IF NOT EXISTS `users` (
  `user_id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `clerk_id` VARCHAR(50) NOT NULL,
  `username` VARCHAR(30) NOT NULL,
  `email` VARCHAR(255) NOT NULL,
  `birthday` DATE,
  `gender` ENUM('MALE', 'FEMALE', 'OTHER'),
  `country_id` TINYINT UNSIGNED,
  `city` VARCHAR(255),
  `description` VARCHAR(300),
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  INDEX `idx_clerk_id` (`clerk_id`)
)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_unicode_ci
;

--  Items table
CREATE TABLE IF NOT EXISTS `items` (
  `item_id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `user_id` BIGINT UNSIGNED NOT NULL,
  `title` VARCHAR(255) NOT NULL,
  `type` ENUM('DREAM', 'GOAL', 'TASK') NOT NULL,
  `target_date` DATE,
  `priority` ENUM('LOW', 'MEDIUM', 'HIGH'),
  `duration` INT UNSIGNED,
  `time_spent` INT UNSIGNED DEFAULT 0,
  `rec_times` INT UNSIGNED,
  `rec_period` ENUM('WEEK', 'DAY', 'MONTH'),
  `rec_progress` INT UNSIGNED,
  `parent_id` BIGINT UNSIGNED,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  INDEX `idx_user_id` (`user_id`)
)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_unicode_ci
;

--  Relationship table
CREATE TABLE IF NOT EXISTS `item_relations` (
  `user_id` BIGINT UNSIGNED NOT NULL,
  `parent_id` BIGINT UNSIGNED,
  `child_id` BIGINT UNSIGNED,
  PRIMARY KEY (`user_id`, `parent_id`, `child_id`)
)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_unicode_ci
;

--  Timer History table
CREATE TABLE IF NOT EXISTS `timer_history` (
  `user_id` BIGINT UNSIGNED NOT NULL,
  `start_time` TIMESTAMP,
  `end_time` TIMESTAMP,
  `item_id` BIGINT UNSIGNED,
  PRIMARY KEY (`user_id`, `end_time`)
)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_unicode_ci
;


-- ======= Sample data ==========

-- Users table

INSERT INTO `users` (`username`, `email`, `clerk_id`) VALUES ('kaspis245', 'kaspis245@gmail.com', 'user_2Ui3QUrAix2lslvsWBmvGeNBKtu');

-- Items table

INSERT INTO `items` (`user_id`, `title`, `type`, `priority`) VALUES (1, 'Learn Spanish', 'DREAM', 'HIGH');
INSERT INTO `items` (`user_id`, `title`, `type`) VALUES (1, 'Get fit', 'DREAM');
INSERT INTO `items` (`user_id`, `title`, `type`) VALUES (1, 'Get good at math', 'DREAM');

INSERT INTO `items` (`user_id`, `title`, `type`, `target_date`, `priority`) VALUES (1, 'Make a todo/timer app', 'GOAL', '2023-12-01', 'HIGH');
INSERT INTO `items` (`user_id`, `title`, `type`, `priority`) VALUES (1, 'Learn chess', 'GOAL', 'LOW');
INSERT INTO `items` (`user_id`, `title`, `type`, `parent_id`) VALUES (1, 'Learn Spanish vocabulary', 'GOAL', 1);
INSERT INTO `item_relations` (`user_id`, `parent_id`, `child_id`) VALUES (1, 1, 6);
INSERT INTO `items` (`user_id`, `title`, `type`, `parent_id`) VALUES (1, 'Learn Spanish grammar', 'GOAL', 1);
INSERT INTO `item_relations` (`user_id`, `parent_id`, `child_id`) VALUES (1, 1, 7);
INSERT INTO `items` (`user_id`, `title`, `type`, `parent_id`) VALUES (1, 'Spanish language comprehension', 'GOAL', 1);
INSERT INTO `item_relations` (`user_id`, `parent_id`, `child_id`) VALUES (1, 1, 8);
INSERT INTO `items` (`user_id`, `title`, `type`, `parent_id`) VALUES (1, 'Spanish writing', 'GOAL', 1);
INSERT INTO `item_relations` (`user_id`, `parent_id`, `child_id`) VALUES (1, 1, 9);
INSERT INTO `items` (`user_id`, `title`, `type`, `parent_id`) VALUES (1, 'Build muscle', 'GOAL', 2);
INSERT INTO `item_relations` (`user_id`, `parent_id`, `child_id`) VALUES (1, 2, 10);
INSERT INTO `items` (`user_id`, `title`, `type`, `parent_id`, `target_date`) VALUES (1, 'Learn Linear Algebra', 'GOAL', 3, '2023-12-01');
INSERT INTO `item_relations` (`user_id`, `parent_id`, `child_id`) VALUES (1, 3, 11);
INSERT INTO `items` (`user_id`, `title`, `type`, `parent_id`) VALUES (1, 'Learn Calculus', 'GOAL', 3);
INSERT INTO `item_relations` (`user_id`, `parent_id`, `child_id`) VALUES (1, 3, 12);
INSERT INTO `items` (`user_id`, `title`, `type`) VALUES (1, 'Read "Demons" by Dostoevsky', 'GOAL');
INSERT INTO `items` (`user_id`, `title`, `type`) VALUES (1, 'Read "The Shape of Space"', 'GOAL');

INSERT INTO `items` (`user_id`, `title`, `type`, `duration`, `target_date`, `priority`, `parent_id`) VALUES (1, 'Make a Figma design sketch', 'TASK', 100800, '2023-10-30', 'MEDIUM', 4);
INSERT INTO `item_relations` (`user_id`, `parent_id`, `child_id`) VALUES (1, 4, 15);
INSERT INTO `items` (`user_id`, `title`, `type`, `duration`, `target_date`, `priority`, `parent_id`) VALUES (1, 'Code MVP frontend', 'TASK', 144000, '2023-10-30', 'HIGH', 4);
INSERT INTO `item_relations` (`user_id`, `parent_id`, `child_id`) VALUES (1, 4, 16);
INSERT INTO `items` (`user_id`, `title`, `type`, `duration`, `target_date`, `priority`, `parent_id`) VALUES (1, 'Make MVP backend', 'TASK', 108000, '2023-10-30', 'HIGH', 4);
INSERT INTO `item_relations` (`user_id`, `parent_id`, `child_id`) VALUES (1, 4, 17);
INSERT INTO `items` (`user_id`, `title`, `type`, `duration`, `parent_id`) VALUES (1, 'Learn common Spanish greeting phrases', 'TASK', 36000, 6);
INSERT INTO `item_relations` (`user_id`, `parent_id`, `child_id`) VALUES (1, 6, 18);
INSERT INTO `items` (`user_id`, `title`, `type`, `duration`, `parent_id`) VALUES (1, 'Memorize a list of essential words', 'TASK', 36000, 6);
INSERT INTO `item_relations` (`user_id`, `parent_id`, `child_id`) VALUES (1, 6, 19);
INSERT INTO `items` (`user_id`, `title`, `type`, `duration`, `parent_id`) VALUES (1, 'Learn Spanish pronunciation', 'TASK', 36000, 6);
INSERT INTO `item_relations` (`user_id`, `parent_id`, `child_id`) VALUES (1, 6, 20);
INSERT INTO `items` (`user_id`, `title`, `type`, `parent_id`, `rec_times`, `rec_period`, `rec_progress`) VALUES (1, 'Do weight lifting', 'TASK', 10, 3, 'WEEK', 2);
INSERT INTO `item_relations` (`user_id`, `parent_id`, `child_id`) VALUES (1, 10, 21);

UPDATE items
SET
  time_spent = time_spent + 1000
WHERE
  user_id = 1 AND item_id = 1;