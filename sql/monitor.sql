create database if not exists monitor;

use monitor;

drop table account;
create table if not exists account (
    id int not null auto_increment,
    ip varchar(20) not null default '',
    created_time int unsigned not null default 0,
    primary key(id),
    unique key uk_ip(ip)
)engine=innodb default charset=utf8;

drop table stats;
create table if not exists stats (
	id bigint not null auto_increment,
	host varchar(20) not null default '',
	load1 decimal(5,2) not null default '0.00',
	load5 decimal(5,2) not null default '0.00',
	load15 decimal(5,2) not null default '0.00',
	buffers bigint not null default 0,
	cached bigint not null default 0,
	memtotal bigint not null default 0,
	memfree bigint not null default 0,
	swaptotal bigint not null default 0,
	swapused bigint not null default 0,
	swapfree bigint not null default 0,
	created_time int unsigned not null default 0,
	primary key(id),
	key idx_host(host)
)engine=innodb default charset=utf8;

drop table diskinfo;
create table if not exists diskinfo(
	id bigint not null auto_increment,
	host varchar(20) not null default '',
	mount varchar(20) not null default '',
	inodeused decimal(4,2) not null default '0.00',
	diskused decimal(4,2) not null default '0.00',
	created_time int unsigned not null default 0,
	primary key(id),
	key idx_host_mount(host,mount)
)engine=innodb default charset=utf8;
