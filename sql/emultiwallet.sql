/*
Navicat MySQL Data Transfer

Source Server         : 192.168.1.107
Source Server Version : 50629
Source Host           : 192.168.1.107:3306
Source Database       : emultiwallet

Target Server Type    : MYSQL
Target Server Version : 50629
File Encoding         : 65001

Date: 2018-12-27 11:45:49
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for tbl_acct_config
-- ----------------------------
DROP TABLE IF EXISTS `tbl_acct_config`;
CREATE TABLE `tbl_acct_config` (
  `acctid` int(11) NOT NULL AUTO_INCREMENT,
  `cellphone` varchar(64) NOT NULL,
  `realname` varchar(64) NOT NULL,
  `idcard` varchar(64) NOT NULL,
  `pubkey` varchar(512) DEFAULT NULL,
  `role` int(11) NOT NULL,
  `state` int(11) NOT NULL,
  `createtime` datetime DEFAULT NULL,
  `updatetime` datetime DEFAULT NULL,
  PRIMARY KEY (`acctid`),
  UNIQUE KEY `cellphone` (`cellphone`),
  UNIQUE KEY `idcard` (`idcard`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for tbl_acct_wallet_relation
-- ----------------------------
DROP TABLE IF EXISTS `tbl_acct_wallet_relation`;
CREATE TABLE `tbl_acct_wallet_relation` (
  `relationid` int(11) NOT NULL AUTO_INCREMENT,
  `acctid` int(11) NOT NULL,
  `walletid` int(11) NOT NULL,
  `createtime` datetime DEFAULT NULL,
  PRIMARY KEY (`relationid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for tbl_coin_config
-- ----------------------------
DROP TABLE IF EXISTS `tbl_coin_config`;
CREATE TABLE `tbl_coin_config` (
  `coinid` int(11) NOT NULL,
  `coinsymbol` varchar(16) NOT NULL,
  `ip` varchar(64) NOT NULL,
  `rpcport` int(11) NOT NULL,
  `rpcuser` varchar(64) DEFAULT NULL,
  `rpcpass` varchar(64) DEFAULT NULL,
  `state` int(11) NOT NULL,
  `createtime` datetime DEFAULT NULL,
  `updatetime` datetime DEFAULT NULL,
  PRIMARY KEY (`coinid`),
  UNIQUE KEY `coinsymbol` (`coinsymbol`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for tbl_notification
-- ----------------------------
DROP TABLE IF EXISTS `tbl_notification`;
CREATE TABLE `tbl_notification` (
  `notifyid` int(11) NOT NULL AUTO_INCREMENT,
  `acctid` int(11) NOT NULL,
  `wallettid` int(11) DEFAULT NULL,
  `trxid` int(11) DEFAULT NULL,
  `notifytype` int(11) NOT NULL,
  `notification` text,
  `state` int(11) NOT NULL,
  `reserved1` text,
  `reserved2` text,
  `createtime` datetime DEFAULT NULL,
  `updatetime` datetime DEFAULT NULL,
  PRIMARY KEY (`notifyid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for tbl_operator_log
-- ----------------------------
DROP TABLE IF EXISTS `tbl_operator_log`;
CREATE TABLE `tbl_operator_log` (
  `logid` int(11) NOT NULL AUTO_INCREMENT,
  `acctid` int(11) NOT NULL,
  `optype` int(11) NOT NULL,
  `content` text,
  `createtime` datetime DEFAULT NULL,
  PRIMARY KEY (`logid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for tbl_pending_transaction
-- ----------------------------
DROP TABLE IF EXISTS `tbl_pending_transaction`;
CREATE TABLE `tbl_pending_transaction` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `trxuuid` varchar(64) NOT NULL,
  `coinid` int(11) NOT NULL,
  `vintrxid` varchar(128) DEFAULT NULL,
  `vinvout` int(11) DEFAULT NULL,
  `fromaddress` varchar(128) DEFAULT NULL,
  `balance` varchar(64) DEFAULT NULL,
  `createtime` datetime DEFAULT NULL,
  `updatetime` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for tbl_pubkey_pool
-- ----------------------------
DROP TABLE IF EXISTS `tbl_pubkey_pool`;
CREATE TABLE `tbl_pubkey_pool` (
  `serverid` int(11) NOT NULL DEFAULT '0',
  `keyindex` int(11) NOT NULL DEFAULT '0',
  `pubkey` varchar(512) NOT NULL,
  `isused` tinyint(1) NOT NULL,
  `createtime` datetime DEFAULT NULL,
  `usedtime` datetime DEFAULT NULL,
  UNIQUE KEY `idx_serverid_keyindex` (`serverid`,`keyindex`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for tbl_sequence
-- ----------------------------
DROP TABLE IF EXISTS `tbl_sequence`;
CREATE TABLE `tbl_sequence` (
  `seqvalue` int(11) NOT NULL AUTO_INCREMENT,
  `idtype` int(11) NOT NULL,
  `state` int(11) NOT NULL,
  `createtime` datetime NOT NULL,
  `updatetime` datetime DEFAULT NULL,
  PRIMARY KEY (`seqvalue`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for tbl_server_info
-- ----------------------------
DROP TABLE IF EXISTS `tbl_server_info`;
CREATE TABLE `tbl_server_info` (
  `serverid` int(11) unsigned NOT NULL,
  `servername` varchar(64) NOT NULL,
  `islocalserver` tinyint(1) NOT NULL,
  `serverpubkey` varchar(256) DEFAULT NULL,
  `serverstartindex` int(11) NOT NULL,
  `serverstatus` int(11) NOT NULL,
  `createtime` datetime DEFAULT NULL,
  `updatetime` datetime DEFAULT NULL,
  PRIMARY KEY (`serverid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for tbl_session_persistence
-- ----------------------------
DROP TABLE IF EXISTS `tbl_session_persistence`;
CREATE TABLE `tbl_session_persistence` (
  `sessionid` varchar(64) NOT NULL,
  `sessionvalue` text NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for tbl_task_persistence
-- ----------------------------
DROP TABLE IF EXISTS `tbl_task_persistence`;
CREATE TABLE `tbl_task_persistence` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `taskuuid` varchar(64) NOT NULL,
  `walletuuid` varchar(64) DEFAULT NULL,
  `trxuuid` varchar(64) DEFAULT NULL,
  `pushtype` int(11) NOT NULL,
  `state` int(11) NOT NULL,
  `createtime` datetime DEFAULT NULL,
  `updatetime` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for tbl_transaction
-- ----------------------------
DROP TABLE IF EXISTS `tbl_transaction`;
CREATE TABLE `tbl_transaction` (
  `trxid` int(11) NOT NULL AUTO_INCREMENT,
  `trxuuid` varchar(64) NOT NULL,
  `rawtrxid` varchar(128) DEFAULT NULL,
  `walletid` int(11) NOT NULL,
  `coinid` int(11) NOT NULL,
  `contractaddr` varchar(128) DEFAULT NULL,
  `acctid` int(11) NOT NULL,
  `serverid` int(11) NOT NULL,
  `fromaddr` varchar(128) NOT NULL,
  `todetails` text NOT NULL,
  `feecost` varchar(128) DEFAULT NULL,
  `trxtime` datetime DEFAULT NULL,
  `needconfirm` int(11) NOT NULL,
  `confirmed` int(11) NOT NULL,
  `acctconfirmed` text NOT NULL,
  `signedtrxs` text,
  `signedserverids` text,
  `fee` varchar(128) DEFAULT NULL,
  `gasprice` varchar(128) DEFAULT NULL,
  `gaslimit` varchar(128) DEFAULT NULL,
  `state` int(11) NOT NULL,
  `signature` varchar(256) DEFAULT NULL,
  PRIMARY KEY (`trxid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for tbl_wallet_config
-- ----------------------------
DROP TABLE IF EXISTS `tbl_wallet_config`;
CREATE TABLE `tbl_wallet_config` (
  `walletid` int(11) NOT NULL AUTO_INCREMENT,
  `walletuuid` varchar(64) NOT NULL,
  `coinid` int(16) NOT NULL,
  `walletname` varchar(64) NOT NULL,
  `serverkeys` varchar(128) NOT NULL,
  `createserver` int(11) NOT NULL,
  `keycount` int(11) NOT NULL,
  `needkeysigcount` int(11) NOT NULL,
  `address` varchar(64) NOT NULL,
  `destaddress` text,
  `needsigcount` int(11) NOT NULL,
  `fee` varchar(64) DEFAULT NULL,
  `gasprice` varchar(64) DEFAULT NULL,
  `gaslimit` varchar(64) DEFAULT NULL,
  `state` int(11) NOT NULL,
  `createtime` datetime DEFAULT NULL,
  `updatetime` datetime DEFAULT NULL,
  PRIMARY KEY (`walletid`),
  UNIQUE KEY `walletname` (`walletname`),
  UNIQUE KEY `keyindex` (`serverkeys`),
  UNIQUE KEY `address` (`address`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
