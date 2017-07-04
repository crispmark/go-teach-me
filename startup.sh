#! /bin/bash
# sudo sh -c 'echo "deb http://apt.postgresql.org/pub/repos/apt/ `lsb_release -cs`-pgdg main" >> /etc/apt/sources.list.d/pgdg.list'
# sudo wget -q https://www.postgresql.org/media/keys/ACCC4CF8.asc -O - | sudo apt-key add -
# sudo apt-get -y update
# sudo apt-get -y upgrade
# sudo apt-get -y install golang-1.8-go postgresql postgresql-contrib libpq-dev pgadmin3

# sudo iptables -I INPUT 5 -i eth0 -p tcp --dport 80 -m state --state NEW,ESTABLISHED -j ACCEPT
sudo rpm -Uvh https://yum.postgresql.org/9.6/redhat/rhel-7-x86_64/pgdg-centos96-9.6-3.noarch.rpm
sudo yum -n update
sudo yum install -y go git epel-release nginx policycoreutils-python postgresql96-server postgresql96

sudo systemctl start nginx
sudo systemctl enable nginx
sudo firewall-cmd --permanent --zone=public --add-service=http
sudo firewall-cmd --permanent --zone=public --add-service=https
sudo firewall-cmd --reload

sudo /usr/pgsql-9.6/bin/postgresql96-setup initdb
sudo systemctl start postgresql-9.6
sudo systemctl enable postgresql-9.6

sudo -u postgres createuser teachme
sudo -u postgres createdb teachme
sudo -u postgres psql << EOF
  \connect teachme
  REVOKE CONNECT ON DATABASE teachme FROM PUBLIC;
  GRANT CONNECT ON DATABASE teachme TO teachme;
  CREATE SEQUENCE IF NOT EXISTS user_group_user_group_id_seq;
  CREATE TABLE IF NOT EXISTS user_groups (
      user_group_id integer primary key not null default nextval('user_group_user_group_id_seq'::regclass),
      user_group_name varchar(20) NOT NULL UNIQUE,
      created_at timestamp NOT NULL default current_timestamp,
      updated_at timestamp NOT NULL default current_timestamp
  );
  ALTER SEQUENCE user_group_user_group_id_seq OWNED BY user_groups.user_group_id;

  INSERT INTO user_groups (user_group_name) VALUES ('ADMIN') ON CONFLICT DO NOTHING;
  INSERT INTO user_groups (user_group_name) VALUES ('STUDENT') ON CONFLICT DO NOTHING;

  CREATE SEQUENCE IF NOT EXISTS user_user_id_seq;
  CREATE TABLE IF NOT EXISTS users (
      user_id integer primary key not null default nextval('user_user_id_seq'::regclass),
      first_name varchar(60) NOT NULL,
      last_name varchar(60) NOT NULL,
      email varchar(256) NOT NULL UNIQUE,
      password varchar(120) NOT NULL,
      user_group_id integer NOT NULL REFERENCES user_groups ON DELETE RESTRICT,
      created_at timestamp NOT NULL default current_timestamp,
      updated_at timestamp NOT NULL default current_timestamp
  );
  ALTER SEQUENCE user_user_id_seq OWNED BY users.user_id;

  CREATE TABLE IF NOT EXISTS files (
      file_id varchar(36) primary key,
      filename varchar(60) NOT NULL,
      data bytea NOT NULL,
      created_at timestamp NOT NULL default current_timestamp
  );
  REVOKE ALL ON ALL TABLES IN SCHEMA public FROM PUBLIC;
  GRANT SELECT, INSERT, UPDATE ON ALL TABLES IN SCHEMA public TO teachme;
  GRANT SELECT, UPDATE ON ALL SEQUENCES IN SCHEMA public TO teachme;
EOF
rm -rf go-teach-me
git clone https://github.com/crispmark/go-teach-me.git
cd go-teach-me
export GOBIN="/go/bin"
export GOPATH="/go/src"
export GOPATH=$GOPATH:/go-teach-me
go get github.com/google/uuid
go get github.com/gorilla/mux
go get github.com/gorilla/sessions
go get github.com/lib/pq
go get golang.org/x/crypto/bcrypt
go build go-teach-me/app
./app
