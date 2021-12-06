-- +goose Up
create table user
(
    id               int unsigned not null auto_increment,
    username         varchar(50)  not null,
    user_type        varchar(10)  not null,
    telegram_chat_id int unsigned,
    primary key (id),
    unique key (username)
);

create table client_user
(
    id          int unsigned not null auto_increment,
    client_guid varchar(36)  not null,
    user_fk     int unsigned not null,
    primary key (id),
    foreign key (user_fk) references user (id)
);

DROP TABLE user, client_user;
