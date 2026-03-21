-- データベースがなければ作成する
CREATE DATABASE IF NOT EXISTS sysken_pay;

-- 'sysken_pay' データベースを使用する
USE sysken_pay;

-- ---------------------------------
-- 1. user テーブル (ユーザー)
-- ---------------------------------
CREATE TABLE `user` (
    id VARCHAR(20) NOT NULL,
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
    user_id VARCHAR(20) NOT NULL,
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
    user_id VARCHAR(20) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES `user`(id)
);

CREATE TABLE purchase_item (
    id INT NOT NULL AUTO_INCREMENT,
    purchase_id INT NOT NULL,  
    item_id INT NOT NULL,      
    quantity INT NOT NULL,     
    price_at_purchase INT NOT NULL, 
    
    PRIMARY KEY (id),
    FOREIGN KEY (purchase_id) REFERENCES purchase(id),
    FOREIGN KEY (item_id) REFERENCES item(id)
);

-- ---------------------------------
-- 5. balance テーブル (残高変動ログ)
-- ---------------------------------
-- ユーザーごとの残高レコードを更新(UPDATE)し続けると、アクセス集中時に行ロックの競合が発生しやすくなります。
-- そのため、入出金の履歴を追記(INSERT)していき、残高は SUM(amount) で算出する設計にします。
-- これにより、chargeとpurchaseの両方と紐づく中間テーブルのような役割も果たします。
CREATE TABLE balance (
    id INT NOT NULL AUTO_INCREMENT,
    user_id VARCHAR(20) NOT NULL,
    charge_id INT DEFAULT NULL,   -- チャージ由来の場合にIDが入る
    purchase_id INT DEFAULT NULL, -- 購入由来の場合にIDが入る
    amount INT NOT NULL,          -- 変動額（入金はプラス、出金はマイナス）
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES `user`(id),
    FOREIGN KEY (charge_id) REFERENCES charge(id),
    FOREIGN KEY (purchase_id) REFERENCES purchase(id),
    
    INDEX idx_user_id (user_id) -- 残高集計(SUM)を高速化するためのインデックス
);

