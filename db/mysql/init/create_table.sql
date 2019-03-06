use test;

CREATE TABLE companies (
  id int(11) unsigned not null auto_increment,
  primary key (id)
);

CREATE TABLE offices (
  id int(11) unsigned not null auto_increment,
  company_id int(11) unsigned not null,
  primary key (id)
);
