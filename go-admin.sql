-- phpMyAdmin SQL Dump
-- version 5.1.1
-- https://www.phpmyadmin.net/
--
-- 主机： localhost
-- 生成日期： 2025-11-28 00:48:25
-- 服务器版本： 9.3.0
-- PHP 版本： 8.0.30

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- 数据库： `go-admin`
--

-- --------------------------------------------------------

--
-- 表的结构 `admin`
--

CREATE TABLE `admin` (
  `id` int NOT NULL,
  `pid` int NOT NULL DEFAULT '0' COMMENT '上级管理员',
  `username` varchar(255) COLLATE utf8mb4_general_ci NOT NULL COMMENT '管理员账号',
  `password` varchar(255) COLLATE utf8mb4_general_ci NOT NULL COMMENT '管理员密码',
  `salt` varchar(64) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '密码盐',
  `avatar` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '管理员头像',
  `email` varchar(128) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '邮箱',
  `failures` int NOT NULL DEFAULT '0' COMMENT '登陆失败次数',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态：0=禁用，1=启用',
  `token` longtext COLLATE utf8mb4_general_ci COMMENT 'TOKEN',
  `realname` varchar(255) COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '真实姓名',
  `mobile` varchar(13) COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '手机号',
  `last_login_time` int NOT NULL DEFAULT '0' COMMENT '最后登录时间',
  `created_at` int NOT NULL COMMENT '创建时间',
  `updated_at` int NOT NULL COMMENT '更新时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- 转存表中的数据 `admin`
--

INSERT INTO `admin` (`id`, `pid`, `username`, `password`, `salt`, `avatar`, `email`, `failures`, `status`, `token`, `realname`, `mobile`, `last_login_time`, `created_at`, `updated_at`) VALUES
(1, 0, 'admin', '$2a$10$DiwUJcVoG1svEvhkEmZQPOnJXKYLwFPg4zp8WSsTlPl.7VF0F2GLu', 'a44b73a6c777626de536b4a26e08aec5', '/logo.png', 'admin@test.test', 0, 1, 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwidXNlcm5hbWUiOiJhZG1pbiIsImV4cCI6MTc2NDg5NTU2NywibmJmIjoxNzY0MjkwNzY3LCJpYXQiOjE3NjQyOTA3Njd9.uBaSb3DdIuBDjHsSAG1eWgBXaXINFpvG8M9_qYgRmbc', '', '', 1746333107, 1733297562, 1764290848);

-- --------------------------------------------------------

--
-- 表的结构 `admin_login_log`
--

CREATE TABLE `admin_login_log` (
  `id` int NOT NULL,
  `username` varchar(255) COLLATE utf8mb4_general_ci NOT NULL COMMENT '账号',
  `ip` varchar(255) COLLATE utf8mb4_general_ci NOT NULL COMMENT 'IP',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '状态',
  `created_at` int NOT NULL COMMENT '时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- 表的结构 `admin_menu`
--

CREATE TABLE `admin_menu` (
  `id` int UNSIGNED NOT NULL COMMENT '主键ID',
  `pid` int UNSIGNED NOT NULL DEFAULT '0' COMMENT '父级ID',
  `name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '菜单名称',
  `enname` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '英文名称',
  `route` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '路由地址',
  `component` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '组件路径',
  `path` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '路由路径',
  `icon` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '图标',
  `sort` int NOT NULL DEFAULT '0' COMMENT '排序',
  `visible` tinyint(1) NOT NULL DEFAULT '1' COMMENT '显示/隐藏：0=隐藏，1=显示',
  `fixed_tag` tinyint(1) NOT NULL DEFAULT '0' COMMENT '固定标签：0=不固定，1=固定',
  `show_tag` tinyint(1) NOT NULL DEFAULT '1' COMMENT '标签是否显示：0=不显示，1=显示',
  `iframe` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT 'iframe链接地址',
  `external` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '外部连接地址',
  `type` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '类型：menu=菜单，button=按钮，iframe=内嵌页面，link=外部链接',
  `permission` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '权限标识',
  `created_at` int UNSIGNED NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` int UNSIGNED NOT NULL DEFAULT '0' COMMENT '更新时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='菜单表';

--
-- 转存表中的数据 `admin_menu`
--

INSERT INTO `admin_menu` (`id`, `pid`, `name`, `enname`, `route`, `component`, `path`, `icon`, `sort`, `visible`, `fixed_tag`, `show_tag`, `iframe`, `external`, `type`, `permission`, `created_at`, `updated_at`) VALUES
(1, 0, '面板', 'Dashboard', '', '/home/index', '/home', 'ri:dashboard-3-line', 100, 1, 0, 0, '', '', 'menu', '', 0, 1764254528),
(2, 0, '权限', 'Admins', '', '/admin', '/admin', 'fa-solid:assistive-listening-systems', 80, 1, 0, 0, '', '', 'menu', '', 0, 1764254536),
(3, 2, '管理账号', 'Admin', '', '/admin/admin/list', '/admin/admin', 'ri:admin-line', 3, 1, 0, 0, '', '', 'menu', '', 1746330668, 1754379801),
(4, 2, '角色管理', 'Role', '', '/admin/role/list', '/admin/role', 'fa-solid:users-cog', 2, 1, 0, 0, '', '', 'menu', '', 0, 1759125251),
(5, 2, '菜单管理', 'Menu', NULL, '/admin/menu/list', '/admin/menu', 'ri:menu-line', 1, 1, 0, 0, '', '', 'menu', '', 0, 1759048594),
(6, 3, '管理员列表', 'AdminList', '/backend/admin/admin/list', '', '', '', 6, 1, 0, 0, '', '', 'button', 'admin:admin:list', 1746330668, 1754168163),
(7, 3, '添加/编辑管理员', 'AdminSave', '/backend/admin/admin/save', '', '', '', 7, 1, 0, 0, '', '', 'button', 'admin:admin:save', 1746330668, 1754168219),
(11, 5, '菜单列表', 'MenuList', '/backend/admin/menu/list', '', '', '', 11, 1, 0, 0, '', '', 'button', 'admin:menu:list', 1754149159, 1754149159),
(12, 4, '角色列表', 'RoleList', '/backend/admin/role/list', '', '', '', 5, 1, 0, 0, '', '', 'button', 'admin:role:list', 1754210565, 1754210565),
(13, 4, '添加/编辑角色', 'RoleSave', '/backend/admin/role/save', '', '', '', 4, 1, 0, 0, '', '', 'button', 'admin:role:save', 1754210677, 1754210677),
(14, 4, '删除角色', 'RoleDelete', '/backend/admin/role/delete', '', '', '', 3, 1, 0, 0, '', '', 'button', 'admin:role:delete', 1754377594, 1754377594),
(16, 4, '权限读取', 'RoleMyMenus', '/backend/admin/role/my-menus', '', '', '', 1, 1, 0, 0, '', '', 'button', 'admin:role:my-menus', 1754379101, 1754379101),
(17, 5, '添加/编辑菜单', 'MenuSave', '/backend/admin/menu/save', '', '', '', 17, 1, 0, 0, '', '', 'button', 'admin:menu:save', 1754379147, 1754379147),
(18, 5, '菜单删除', 'MenuDelete', '/backend/admin/menu/delete', '', '', '', 18, 1, 0, 0, '', '', 'button', 'admin:menu:delete', 1754379175, 1754379175),
(19, 0, '系统', 'Systems', '', '/system', '/system', 'carbon:settings', 90, 1, 0, 0, '', '', 'menu', '', 0, 1764254532),
(20, 19, '附件管理', 'Attachments', '', '/system/attachment/index', '/system/attachment', 'ri:attachment-2', 40, 1, 0, 0, '', '', 'menu', '', 0, 1764253574),
(22, 1, '帐号信息', 'AccountInfo', '', '/home/userinfo', '/userinfo', 'ri:dashboard-3-line', 2, 0, 0, 0, '', '', 'menu', '', 1746330668, 1754168134),
(30, 1, '数据统计', 'HomeStatic', '', '/home/index', '/home', 'ri:dashboard-3-line', 100, 0, 0, 0, '', '', 'menu', '', 1746330668, 1754168134),
(31, 19, '配置管理', 'Configs', '', '/system/config/index', '/system/config', 'carbon:cloud-satellite-config', 50, 1, 0, 0, '', '', 'menu', '', 0, 1764253553),
(32, 0, '用户', 'User', '', '/user/list', '/user', 'carbon:user', 70, 1, 0, 0, '', '', 'menu', '', 0, 1764254994),
(33, 32, '用户列表', 'Users', '', '/user/list', '/user/list', 'carbon:user-multiple', 1, 1, 0, 0, '', '', 'menu', '', 0, 1764254984),
(34, 19, '代码生成', 'AutoCrud', '', '/system/gen/list', '/system/gen', 'carbon:code', 10, 1, 0, 0, '', '', 'menu', '', 0, 1764253645),
(35, 31, '配置列表', 'ConfigList', '/backend/system/config/list', '', '', '', 1, 1, 0, 0, '', '', 'button', 'system:config:list', 1764253696, 1764253696),
(36, 31, '配置创建', 'ConfigCreate', '/backend/system/config/create', '', '', '', 2, 1, 0, 0, '', '', 'button', 'system:config:create', 1764253729, 1764253729),
(37, 31, '配置更新', 'ConfigUpdate', '/backend/system/config/update', '', '', '', 3, 1, 0, 0, '', '', 'button', 'system:config:update', 1764253832, 1764253832),
(38, 31, '配置删除', 'ConfigDelete', '/backend/system/config/delete', '', '', '', 4, 1, 0, 0, '', '', 'button', 'system:config:delete', 1764253862, 1764253862),
(39, 20, '附件文件夹', 'AttachmentDirectories', '/backend/system/attachment/directories', '', '', '', 1, 1, 0, 0, '', '', 'button', 'system:attachment:directories', 1764253917, 1764253917),
(40, 20, '附件列表', 'AttachmentList', '/backend/system/attachment/list', '', '', '', 2, 1, 0, 0, '', '', 'button', 'system:attachment:list', 1764253950, 1764253950),
(41, 20, '附件上传', 'AttachmentUpload', '/backend/system/attachment/upload', '', '', '', 3, 1, 0, 0, '', '', 'button', 'system:attachment:upload', 1764253983, 1764253983),
(42, 20, '附件删除', 'AttachmentDelete', '/backend/system/attachment/delete', '', '', '', 4, 1, 0, 0, '', '', 'button', 'system:attachment:delete', 1764254020, 1764254020),
(43, 34, '表列表', 'GenTableList', '/backend/system/gen/table-list', '', '', '', 1, 1, 0, 0, '', '', 'button', 'system:gen:table-list', 1764254071, 1764254071),
(44, 34, '表详情', 'GenTableInfo', '/backend/system/gen/table-info', '', '', '', 2, 1, 0, 0, '', '', 'button', 'system:gen:table-info', 1764254141, 1764254141),
(45, 34, '表配置', 'GenTableConfig', '/backend/system/gen/table-config', '', '', '', 3, 1, 0, 0, '', '', 'button', 'system:gen:table-config', 1764254169, 1764254169),
(46, 34, '预览', 'GenPreview', '/backend/system/gen/preview', '', '', '', 4, 1, 0, 0, '', '', 'button', 'system:gen:preview', 1764254202, 1764254202),
(47, 34, '生成', 'GenGenerate', '/backend/system/gen/generate', '', '', '', 5, 1, 0, 0, '', '', 'button', 'system:gen:generate', 1764254281, 1764254281),
(48, 34, '历史', 'GenHistory', '/backend/system/gen/history', '', '', '', 6, 1, 0, 0, '', '', 'button', 'system:gen:history', 1764254324, 1764254324),
(49, 34, '删除', 'GenDelete', '/backend/system/gen/delete', '', '', '', 7, 1, 0, 0, '', '', 'button', 'system:gen:delete', 1764254354, 1764254354),
(50, 34, '下载', 'GenDownload', '/backend/system/gen/download', '', '', '', 8, 1, 0, 0, '', '', 'button', 'system:gen:download', 1764254382, 1764254382),
(51, 3, '管理员删除', 'AdminDelete', '/backend/admin/admin/delete', '', '', '', 8, 1, 0, 0, '', '', 'button', 'admin:admin:delete', 1764254434, 1764254434),
(52, 33, '用户读取', 'UserList', '/backend/user/list', '', '', '', 1, 1, 0, 0, '', '', 'button', 'user:list', 1764254690, 1764254690),
(53, 33, '用户创建', 'UserCreate', '/backend/user/create', '', '', '', 2, 1, 0, 0, '', '', 'button', 'user:create', 1764254710, 1764254710),
(54, 33, '用户更新', 'UserUpdate', '/backend/user/update', '', '', '', 3, 1, 0, 0, '', '', 'button', 'user:update', 1764254727, 1764254727),
(55, 33, '用户删除', 'UserDelete', '/backend/user/delete', '', '', '', 4, 1, 0, 0, '', '', 'button', 'user:delete', 1764254768, 1764254768),
(56, 33, '用户操作', 'UserOperate', '/backend/user/operate', '', '', '', 5, 1, 0, 0, '', '', 'button', 'user:operate', 1764254787, 1764254787),
(57, 33, '用户余额操作', 'UserMoney', '/backend/user/update-money', '', '', '', 6, 1, 0, 0, '', '', 'button', 'user:money', 1764254815, 1764254815),
(58, 33, '用户积分操作', 'UserScore', '/backend/user/update-score', '', '', '', 7, 1, 0, 0, '', '', 'button', 'user:score', 1764254833, 1764254833),
(59, 30, '统计接口', 'HomeIndex', '/backend/home/index', '', '', '', 1, 1, 0, 0, '', '', 'button', 'home:index', 1764255590, 1764255590);

-- --------------------------------------------------------

--
-- 表的结构 `admin_role`
--

CREATE TABLE `admin_role` (
  `id` int UNSIGNED NOT NULL COMMENT '主键ID',
  `pid` int UNSIGNED NOT NULL DEFAULT '0' COMMENT '父级ID',
  `name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '角色名称',
  `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '角色描述',
  `is_super` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否超级管理员：0=否，1=是',
  `sort` int NOT NULL DEFAULT '50' COMMENT '排序',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态',
  `created_at` int UNSIGNED NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` int UNSIGNED NOT NULL DEFAULT '0' COMMENT '更新时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='角色表';

--
-- 转存表中的数据 `admin_role`
--

INSERT INTO `admin_role` (`id`, `pid`, `name`, `description`, `is_super`, `sort`, `status`, `created_at`, `updated_at`) VALUES
(1, 0, '超级管理员', '拥有所有权限', 1, 100, 1, 1746330668, 1754143818);

-- --------------------------------------------------------

--
-- 表的结构 `admin_role_admin`
--

CREATE TABLE `admin_role_admin` (
  `id` int UNSIGNED NOT NULL COMMENT '主键ID',
  `admin_id` int UNSIGNED NOT NULL COMMENT '管理员ID',
  `role_id` int UNSIGNED NOT NULL COMMENT '角色ID',
  `created_at` int UNSIGNED NOT NULL DEFAULT '0' COMMENT '创建时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='管理员-角色关联表';

--
-- 转存表中的数据 `admin_role_admin`
--

INSERT INTO `admin_role_admin` (`id`, `admin_id`, `role_id`, `created_at`) VALUES
(1, 1, 1, 0);

-- --------------------------------------------------------

--
-- 表的结构 `admin_role_menu`
--

CREATE TABLE `admin_role_menu` (
  `id` int UNSIGNED NOT NULL COMMENT '主键ID',
  `role_id` int UNSIGNED NOT NULL COMMENT '角色ID',
  `menu_id` int UNSIGNED NOT NULL COMMENT '菜单ID',
  `created_at` int UNSIGNED NOT NULL DEFAULT '0' COMMENT '创建时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='角色-菜单关联表';

-- --------------------------------------------------------

--
-- 表的结构 `attachment`
--

CREATE TABLE `attachment` (
  `id` int UNSIGNED NOT NULL,
  `type` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'local' COMMENT '存储类型',
  `admin_id` int UNSIGNED NOT NULL DEFAULT '0' COMMENT '管理员ID',
  `user_id` int UNSIGNED NOT NULL DEFAULT '0' COMMENT '用户ID',
  `path` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '存储路径',
  `parent` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '父级',
  `url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '在线地址',
  `filename` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '文件名称',
  `size` int UNSIGNED NOT NULL COMMENT '文件大小',
  `media_type` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '文件类型',
  `extension` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '文件后缀',
  `created_at` int UNSIGNED NOT NULL COMMENT '创建时间',
  `updated_at` int UNSIGNED DEFAULT '0' COMMENT '更新时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='附件列表';

-- --------------------------------------------------------

--
-- 表的结构 `config`
--

CREATE TABLE `config` (
  `id` int UNSIGNED NOT NULL,
  `dir` varchar(255) COLLATE utf8mb4_general_ci NOT NULL COMMENT '配置组',
  `key` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '配置键',
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '配置名称',
  `tip` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '提示说明',
  `type` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT 'string' COMMENT '配置类型',
  `value` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '配置值',
  `variable` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '配置变量',
  `created_at` int NOT NULL DEFAULT '0' COMMENT '创建时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='配置';

--
-- 转存表中的数据 `config`
--

INSERT INTO `config` (`id`, `dir`, `key`, `name`, `tip`, `type`, `value`, `variable`, `created_at`) VALUES
(1, '网站设置', 'site_name', '网站名称', '', 'input', '吊炸天', '', 1764256035),
(2, '短信配置', 'sms_type', '短信服务商', '', 'radio', 'smsbao', '[\n  {\"label\": \"阿里云\", \"value\": \"aliyun\"},\n  {\"label\": \"腾讯云\", \"value\": \"tencent\"},\n  {\"label\": \"云片\", \"value\": \"yunpian\"},\n{\"label\": \"短信宝\", \"value\": \"smsbao\"}\n]', 1764329758),
(3, '短信配置', 'sms_id', '短信配置ID', '阿里云：assessKeyId，腾讯云：APPID，云片：留空，短信宝：用户名', 'input', '', '', 1764329908),
(4, '短信配置', 'sms_key', '短信配置KEY', '阿里云：assessKeySecret，腾讯云：APPKEY，云片：ApiKey，短信宝：用户密码', 'input', '', '', 1764329943),
(5, '短信配置', 'sms_token', '短信配置Token', '平台Token/短信签名', 'input', '', '', 1764329986),
(6, '短信配置', 'sms_template', '短信通用模版ID', '', 'input', '', '', 1764330007);

-- --------------------------------------------------------

--
-- 表的结构 `gen_history`
--

CREATE TABLE `gen_history` (
  `id` int NOT NULL COMMENT '主键ID',
  `table_name` varchar(255) NOT NULL COMMENT '表名',
  `table_comment` varchar(255) DEFAULT NULL COMMENT '表注释',
  `module_name` varchar(255) NOT NULL COMMENT '模块名',
  `struct_name` varchar(255) NOT NULL COMMENT '结构体名',
  `package_name` varchar(255) NOT NULL COMMENT '包名',
  `frontend_src_path` varchar(500) DEFAULT NULL COMMENT '前端 src 路径',
  `config` text COMMENT '完整配置（JSON）',
  `created_at` int DEFAULT NULL COMMENT '创建时间',
  `updated_at` int DEFAULT NULL COMMENT '更新时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='代码生成历史';

-- --------------------------------------------------------

--
-- 表的结构 `user`
--

CREATE TABLE `user` (
  `id` int UNSIGNED NOT NULL,
  `pid` int NOT NULL DEFAULT '0' COMMENT '上级',
  `tid` int NOT NULL DEFAULT '0' COMMENT '顶级',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态',
  `status_text` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '状态信息',
  `code` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '邀请码',
  `wechat_unionid` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'UNIONID',
  `wechat_openid` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'OPENID',
  `version` int UNSIGNED DEFAULT '1' COMMENT '版本',
  `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '头像',
  `username` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '账号',
  `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '密码',
  `salt` varchar(255) COLLATE utf8mb4_general_ci NOT NULL COMMENT '密码盐',
  `email` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '邮箱',
  `mobile` varchar(11) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '手机号码',
  `score` decimal(10,2) DEFAULT '0.00' COMMENT '积分',
  `money` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '余额',
  `token` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT 'TOKEN',
  `created_at` int NOT NULL COMMENT '创建时间',
  `updated_at` int NOT NULL DEFAULT '0' COMMENT '更新时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='用户表';

-- --------------------------------------------------------

--
-- 表的结构 `user_money_log`
--

CREATE TABLE `user_money_log` (
  `id` int UNSIGNED NOT NULL,
  `user_id` int UNSIGNED NOT NULL COMMENT '用户',
  `type` int UNSIGNED DEFAULT '1' COMMENT '类型',
  `money` decimal(10,2) DEFAULT NULL COMMENT '金额',
  `note` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '备注',
  `source` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '来源',
  `created_at` int DEFAULT NULL COMMENT '创建时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- 表的结构 `user_score_log`
--

CREATE TABLE `user_score_log` (
  `id` int UNSIGNED NOT NULL,
  `user_id` int UNSIGNED NOT NULL COMMENT '用户',
  `type` int UNSIGNED DEFAULT '1' COMMENT '类型',
  `score` decimal(10,2) DEFAULT NULL COMMENT '数量',
  `note` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '备注',
  `source` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '来源',
  `created_at` int DEFAULT NULL COMMENT '创建时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- 转储表的索引
--

--
-- 表的索引 `admin`
--
ALTER TABLE `admin`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_admin_username` (`username`),
  ADD KEY `idx_admin_status` (`status`),
  ADD KEY `idx_admin_pid` (`pid`),
  ADD KEY `idx_admin_username_status` (`username`,`status`);

--
-- 表的索引 `admin_login_log`
--
ALTER TABLE `admin_login_log`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_login_log_ip` (`ip`),
  ADD KEY `idx_login_log_created_at` (`created_at`),
  ADD KEY `username` (`username`);

--
-- 表的索引 `admin_menu`
--
ALTER TABLE `admin_menu`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_pid` (`pid`),
  ADD KEY `idx_menu_pid` (`pid`),
  ADD KEY `idx_menu_type` (`type`),
  ADD KEY `idx_menu_route` (`route`),
  ADD KEY `idx_menu_type_route` (`type`,`route`),
  ADD KEY `visible` (`visible`);

--
-- 表的索引 `admin_role`
--
ALTER TABLE `admin_role`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_pid` (`pid`),
  ADD KEY `idx_role_pid` (`pid`),
  ADD KEY `idx_role_is_super` (`is_super`),
  ADD KEY `idx_role_status` (`status`);

--
-- 表的索引 `admin_role_admin`
--
ALTER TABLE `admin_role_admin`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `idx_admin_role` (`admin_id`,`role_id`),
  ADD KEY `idx_role_admin_admin_id` (`admin_id`),
  ADD KEY `idx_role_admin_role_id` (`role_id`),
  ADD KEY `idx_role_admin_admin_role` (`admin_id`,`role_id`);

--
-- 表的索引 `admin_role_menu`
--
ALTER TABLE `admin_role_menu`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `idx_role_menu` (`role_id`,`menu_id`),
  ADD KEY `idx_role_menu_role_id` (`role_id`),
  ADD KEY `idx_role_menu_menu_id` (`menu_id`),
  ADD KEY `idx_role_menu_role_menu` (`role_id`,`menu_id`);

--
-- 表的索引 `attachment`
--
ALTER TABLE `attachment`
  ADD PRIMARY KEY (`id`),
  ADD KEY `filename` (`filename`),
  ADD KEY `url` (`url`),
  ADD KEY `path` (`path`),
  ADD KEY `admin_id` (`admin_id`),
  ADD KEY `user_id` (`user_id`),
  ADD KEY `parent` (`parent`),
  ADD KEY `updated_at` (`updated_at`),
  ADD KEY `idx_attachment_type` (`type`),
  ADD KEY `idx_attachment_path` (`path`),
  ADD KEY `idx_attachment_created_at` (`created_at`);

--
-- 表的索引 `config`
--
ALTER TABLE `config`
  ADD PRIMARY KEY (`id`),
  ADD KEY `dir` (`dir`),
  ADD KEY `key` (`key`),
  ADD KEY `idx_config_key` (`key`),
  ADD KEY `idx_config_group` (`dir`),
  ADD KEY `idx_config_group_key` (`dir`,`key`);

--
-- 表的索引 `gen_history`
--
ALTER TABLE `gen_history`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_table_name` (`table_name`),
  ADD KEY `idx_created_at` (`created_at`);

--
-- 表的索引 `user`
--
ALTER TABLE `user`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `mobile` (`mobile`) USING BTREE,
  ADD KEY `code` (`code`),
  ADD KEY `version` (`id`,`version`),
  ADD KEY `idx_user_username` (`username`),
  ADD KEY `idx_user_mobile` (`mobile`),
  ADD KEY `idx_user_email` (`email`),
  ADD KEY `idx_user_code` (`code`),
  ADD KEY `idx_user_status` (`status`),
  ADD KEY `idx_user_created_at` (`created_at`);

--
-- 表的索引 `user_money_log`
--
ALTER TABLE `user_money_log`
  ADD PRIMARY KEY (`id`),
  ADD KEY `1` (`user_id`,`type`,`created_at`),
  ADD KEY `2` (`source`),
  ADD KEY `idx_money_log_user_id` (`user_id`),
  ADD KEY `idx_money_log_type` (`type`),
  ADD KEY `idx_money_log_created_at` (`created_at`),
  ADD KEY `idx_money_log_user_created` (`user_id`,`created_at`);

--
-- 表的索引 `user_score_log`
--
ALTER TABLE `user_score_log`
  ADD PRIMARY KEY (`id`),
  ADD KEY `1` (`user_id`,`type`,`created_at`),
  ADD KEY `2` (`source`),
  ADD KEY `idx_score_log_user_id` (`user_id`),
  ADD KEY `idx_score_log_type` (`type`),
  ADD KEY `idx_score_log_created_at` (`created_at`),
  ADD KEY `idx_score_log_user_created` (`user_id`,`created_at`);

--
-- 在导出的表使用AUTO_INCREMENT
--

--
-- 使用表AUTO_INCREMENT `admin`
--
ALTER TABLE `admin`
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=2;

--
-- 使用表AUTO_INCREMENT `admin_login_log`
--
ALTER TABLE `admin_login_log`
  MODIFY `id` int NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `admin_menu`
--
ALTER TABLE `admin_menu`
  MODIFY `id` int UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID', AUTO_INCREMENT=60;

--
-- 使用表AUTO_INCREMENT `admin_role`
--
ALTER TABLE `admin_role`
  MODIFY `id` int UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID', AUTO_INCREMENT=2;

--
-- 使用表AUTO_INCREMENT `admin_role_admin`
--
ALTER TABLE `admin_role_admin`
  MODIFY `id` int UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID', AUTO_INCREMENT=2;

--
-- 使用表AUTO_INCREMENT `admin_role_menu`
--
ALTER TABLE `admin_role_menu`
  MODIFY `id` int UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID';

--
-- 使用表AUTO_INCREMENT `attachment`
--
ALTER TABLE `attachment`
  MODIFY `id` int UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `config`
--
ALTER TABLE `config`
  MODIFY `id` int UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=2;

--
-- 使用表AUTO_INCREMENT `gen_history`
--
ALTER TABLE `gen_history`
  MODIFY `id` int NOT NULL AUTO_INCREMENT COMMENT '主键ID';

--
-- 使用表AUTO_INCREMENT `user`
--
ALTER TABLE `user`
  MODIFY `id` int UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `user_money_log`
--
ALTER TABLE `user_money_log`
  MODIFY `id` int UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `user_score_log`
--
ALTER TABLE `user_score_log`
  MODIFY `id` int UNSIGNED NOT NULL AUTO_INCREMENT;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
