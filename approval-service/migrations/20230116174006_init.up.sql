-- create "payments" table
CREATE TABLE `payments` (`id` bigint NOT NULL AUTO_INCREMENT, `request_id` char(36) NOT NULL, `status` enum('PENDING','APPROVED','REJECTED') NOT NULL DEFAULT 'PENDING', `remark` varchar(255) NULL, `create_time` datetime(6) NOT NULL, `update_time` datetime(6) NOT NULL, PRIMARY KEY (`id`), UNIQUE INDEX `payment_request_id` (`request_id`), INDEX `payment_create_time_id` (`create_time` DESC, `id` DESC)) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- create "reviews" table
CREATE TABLE `reviews` (`id` bigint NOT NULL AUTO_INCREMENT, `event` enum('APPROVE','REJECT') NOT NULL, `reviewer_id` varchar(255) NOT NULL, `comment` varchar(255) NULL, `create_time` datetime(6) NOT NULL, `update_time` datetime(6) NOT NULL, `payment_id` bigint NOT NULL, PRIMARY KEY (`id`), UNIQUE INDEX `review_payment_id_reviewer_id` (`payment_id`, `reviewer_id`), CONSTRAINT `reviews_payments_reviews` FOREIGN KEY (`payment_id`) REFERENCES `payments` (`id`) ON DELETE NO ACTION) CHARSET utf8mb4 COLLATE utf8mb4_bin;
