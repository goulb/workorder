create table departments (name);
create table providers (name);
create table users (name,password,department_id,privileges);
create table sessions (uuid,user_id,created_at);
create table cartypes (weight,type_name);
create table orders(department_id,date_begin,date_end,provider_id,cartype_id,
	carnum,usefor,submit,locked,careat_at);