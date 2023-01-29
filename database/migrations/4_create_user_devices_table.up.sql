CREATE SEQUENCE user_devices_seq;

CREATE TABLE user_devices
(
	id INT NOT NULL DEFAULT NEXTVAL ('user_devices_seq'),
	uuid CHAR(36) NOT NULL UNIQUE,
    user_id INT NOT NULL,
    token VARCHAR(255) NOT NULL,
    ip VARCHAR(255),
    location VARCHAR(255),
    platform VARCHAR(255),
    user_agent VARCHAR(255),
    app_version VARCHAR(255),
    vendor VARCHAR(255),
	created_at TIMESTAMP(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by INT DEFAULT NULL,
	updated_at TIMESTAMP(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by INT DEFAULT NULL,

    CONSTRAINT fk_user_devices_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,

	PRIMARY KEY (id)
);