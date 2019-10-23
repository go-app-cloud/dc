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
  `check`     int           null comment '检测方式'
);

create table application
(
  id          varchar(34)   not null comment '唯一ID'
    primary key,
  name        varchar(128)  not null comment '应用名称',
  description varchar(1024) not null comment '详情描述',
  section     varchar(128)  not null comment '数据维护部门',
  sources     TEXT          not null comment '数据来源项目IDs',
  owner       varchar(32)   not null comment '所属者',
  app_key     varchar(128)  not null comment 'app key',
  type        varchar(128)  not null comment '应用类型',
  secret_key  varchar(128)  not null comment 'secret key'
);