CREATE DATABASE IF NOT EXISTS balances;
USE balances;
CREATE TABLE IF NOT EXISTS user_balances(
                              id INT(12) unsigned NOT NULL AUTO_INCREMENT,
                              balance FLOAT(6) unsigned DEFAULT 0,
                              reserved FLOAT(6) unsigned DEFAULT 0,
                              PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS orders(
                       order_id INT(12) unsigned NOT NULL,
                       service_id INT(12) unsigned NOT NULL,
                       user_id INT(12) unsigned NOT NULL,
                       cost FLOAT(6) unsigned DEFAULT 0,
                       PRIMARY KEY (order_id),
                       FOREIGN KEY (user_id)
                           REFERENCES user_balances(id)
                           ON DELETE CASCADE
);
