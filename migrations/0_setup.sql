CREATE DATABASE IF NOT EXISTS `food_decider` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE admins
(
    id            INT AUTO_INCREMENT PRIMARY KEY,
    telegram_id   INT          NOT NULL,
    name          VARCHAR(255) NOT NULL,
    created_at    TIMESTAMP       DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP       DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_by    INT          NOT NULL,
    updated_by    INT          NULL,
    is_superadmin BOOLEAN         DEFAULT FALSE,
    status        ENUM ('A', 'I') DEFAULT 'A'
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

CREATE TABLE food
(
    id          VARCHAR(36) PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    description TEXT         NULL,
    created_at  TIMESTAMP       DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP       DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_by  INT          NOT NULL,
    updated_by  INT          NULL,
    status      ENUM ('A', 'D') DEFAULT 'A'
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

CREATE TABLE food_groups
(
    id          int AUTO_INCREMENT PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    created_at  TIMESTAMP       DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP       DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_by  INT          NOT NULL,
    updated_by  INT          NULL,
    status      ENUM ('A', 'D') DEFAULT 'A'
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

CREATE TABLE food_groups_link
(
    id          int AUTO_INCREMENT PRIMARY KEY,
    food_id     VARCHAR(36) NOT NULL,
    group_id    int NOT NULL,
    created_at  TIMESTAMP       DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP       DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_by  INT          NOT NULL,
    updated_by  INT          NULL,
    status      ENUM ('A', 'D') DEFAULT 'A',
    FOREIGN KEY (food_id) REFERENCES food (id),
    FOREIGN KEY (group_id) REFERENCES food_groups (id)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

CREATE TABLE locations
(
    id          int AUTO_INCREMENT PRIMARY KEY,
    food_id     VARCHAR(36) NOT NULL,
    name        VARCHAR(255) NOT NULL,
    latitude    DECIMAL(10, 8) NOT NULL,
    longitude   DECIMAL(11, 8) NOT NULL,
    created_at  TIMESTAMP       DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP       DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_by  INT          NOT NULL,
    updated_by  INT          NULL,
    status      ENUM ('A', 'D') DEFAULT 'A',
    FOREIGN KEY (food_id) REFERENCES food (id)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;