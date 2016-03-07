create database if not exists monitor;

use monitor;

create table account if not exists (
    id int not null auto_increment,
    ip varchar(20) not null default '',
    created_time int unsigned not null default 0,
    primary key(id),
    unique key uk_ip(ip)
)engine=innodb default charset=utf8;

create table avgload (
	id bigint not null auto_increment,
	host varchar(20) not null default '',
	load1 decimal(5,2) not null default '0.00',
	load5 decimal(5,2) not null default '0.00',
	load15 decimal(5,2) not null default '0.00',
	created_time int unsigned not null default 0,
	primary key(id),
	key idx_host(host)
)engine=innodb default charset=utf8;
