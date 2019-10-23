create schema matrix default charset utf8 COLLATE utf8_general_ci;
use matrix;
create table application
(
  id          varchar(34)   not null comment '唯一ID'
    primary key,
  name        varchar(128)  not null comment '应用名称',
  description varchar(1024) not null comment '详情描述',
  section     varchar(128)  not null comment '数据维护部门',
  sources     text          not null comment '数据来源项目IDs',
  owner       varchar(32)   not null comment '所属者',
  secret_key  varchar(128)  not null comment 'secret key',
  type        varchar(128)  null comment '应用类型'
);

create table source
(
  id          varchar(34)   not null comment '唯一ID'
    primary key,
  name        varchar(128)  not null comment '名称',
  description varchar(1024) not null comment '详情描述',
  section     varchar(128)  not null comment '数据维护部门',
  owner       varchar(32)   not null comment '所属者',
  service     text          not null comment '服务URI地址',
  api_doc     text          not null comment 'API开发文档URI',
  type        varchar(128)  not null comment '连通方式',
  `check`     int           null comment '检测方式',
  secret      varchar(128)  not null comment '加密字符串',
  create_at   bigint        not null comment '创建时间'
);

create table user
(
  id        varchar(32)  not null comment 'ID',
  username  varchar(128) not null comment '用户名',
  password  varchar(128) not null comment '密码',
  salt      varchar(128) not null comment '加密盐',
  create_at bigint       not null comment '创建时间'
)comment '用户名';
INSERT INTO matrix.user (id, username, password, salt, create_at) VALUES ('b6c837870b0a4fcca737e9459867500e', 'admin', '8e4f00f95901e923569a38207c2c5870', 'pDT9fFDSVrSet1VbKFacFO5Q4DUWijYp', 0);