CREATE TABLE `users` ( \
    `id` int unsigned PRIMARY KEY AUTO_INCREMENT, \
    `session_id` varchar(255) NOT NULL COMMENT 'user session id', \
    `state` varchar(255) COMMENT 'user state', \
    `token` varchar(255) COMMENT 'user token', \
    `expiry` timestamp COMMENT 'user expiry', \
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP, \
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP \
) \
ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC;

CREATE TABLE `github_tokens` ( \
    `id` int unsigned PRIMARY KEY AUTO_INCREMENT, \
    `token` varchar(255) NOT NULL COMMENT 'user github token', \
    `type` varchar(255) NOT NULL COMMENT 'user github token type', \
    `refresh_token` varchar(255) NOT NULL COMMENT 'user refresh token', \
    `expiry` timestamp NOT NULL COMMENT 'user github token expiry', \
    `users_id` int unsigned NOT NULL COMMENT 'associated with users.id', \
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP, \
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, \
    FOREIGN KEY (`users_id`) \
    REFERENCES `users` (`id`) \
) \
ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC;

