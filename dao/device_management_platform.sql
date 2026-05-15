create database device_management_platform;
use device_management_platform;

create table device (
    id bigint primary key auto_increment,
    device_no varchar(255) not null,
    register_time datetime not null,
    create_time datetime not null, 
    update_time datetime not null,
    uk_device_no unique (device_no)
);