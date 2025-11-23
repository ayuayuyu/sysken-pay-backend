-- データベースがなければ作成する
CREATE DATABASE IF NOT EXISTS app;

-- 'app' データベースを使用する
USE app;

-- ---------------------------------
-- 1. user テーブル (ユーザー)
-- ---------------------------------
CREATE TABLE `user` (
    id CHAR(36) NOT NULL ,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    
    PRIMARY KEY (id)
);

-- ---------------------------------
-- 2. item テーブル (商品)
-- ---------------------------------
CREATE TABLE item (
    id INT NOT NULL AUTO_INCREMENT,
    jan_code VARCHAR(13) NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    price INT NOT NULL CHECK (price >= 0),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    
    PRIMARY KEY (id)
);

-- ---------------------------------
-- 3. charge テーブル (チャージ履歴)
-- ---------------------------------
CREATE TABLE charge (
    id INT NOT NULL AUTO_INCREMENT,
    user_id CHAR(36) NOT NULL,
    amount INT NOT NULL CHECK (amount > 0),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES `user`(id)
);

-- ---------------------------------
-- 4. purchase テーブル (購入履歴)
-- ---------------------------------
CREATE TABLE purchase (
    id INT NOT NULL AUTO_INCREMENT,
    user_id CHAR(36) NOT NULL,
    item_id INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES `user`(id),
    FOREIGN KEY (item_id) REFERENCES item(id)
);
