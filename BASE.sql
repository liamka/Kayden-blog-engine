# ************************************************************
# SQL dump for Kayden
# ************************************************************

CREATE TABLE `kayden_blog_posts` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `title` varchar(1000) NOT NULL,
  `body` varchar(5000) NOT NULL,
  `tags` varchar(1000) NOT NULL DEFAULT '',
  `time` int(100) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;

INSERT INTO `kayden_blog_posts` (`id`, `title`, `body`, `tags`, `time`)
VALUES (1,'Hello golang!','This is test message','hello, world, ',1397480139);



CREATE TABLE `kayden_blog_drafts` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `title` varchar(1000) NOT NULL,
  `body` varchar(5000) NOT NULL,
  `tags` varchar(1000) NOT NULL DEFAULT '',
  `time` int(100) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;

INSERT INTO `kayden_blog_drafts` (`id`, `title`, `body`, `tags`, `time`)
VALUES (1,'Test draft','Example draft','draft, tags, ',1397480139);