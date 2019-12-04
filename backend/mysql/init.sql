CREATE TABLE `github_tokens` ( \
    `id` int unsigned PRIMARY KEY AUTO_INCREMENT, \
    `token` varchar(255) NOT NULL COMMENT 'user github token', \
    `type` varchar(255) NOT NULL COMMENT 'user github token type', \
    `refresh_token` varchar(255) NOT NULL COMMENT 'user refresh token', \
    `expiry` timestamp NOT NULL COMMENT 'user github token expiry' \
) \
ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC;

CREATE TABLE `users` ( \
    `id` int unsigned PRIMARY KEY AUTO_INCREMENT, \
    `session_id` varchar(255) NOT NULL COMMENT 'user session id', \
    `state` varchar(255) COMMENT 'user state', \
    `token` varchar(255) COMMENT 'user token', \
    `expiry` timestamp NOT NULL COMMENT 'user expiry', \
    `github_tokens_id` int unsigned NOT NULL, \
    FOREIGN KEY (`github_tokens_id`) \
    REFERENCES `github_tokens` (`id`) \
) \
ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC;
