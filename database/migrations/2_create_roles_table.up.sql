CREATE SEQUENCE roles_seq;

CREATE TABLE roles
(
	id INT NOT NULL DEFAULT NEXTVAL ('roles_seq'),
	uuid CHAR(36) NOT NULL UNIQUE,
	name VARCHAR(255) UNIQUE NOT NULL,
	created_at TIMESTAMP(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by INT DEFAULT NULL,
	updated_at TIMESTAMP(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by INT DEFAULT NULL,

	PRIMARY KEY (id)
);

CREATE TABLE role_permissions 
(
    role_id INT NOT NULL,
    permission_id INT NOT NULL,

    CONSTRAINT fk_role_permissions_role_id FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    CONSTRAINT fk_role_permissions_permission_id FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE,

    PRIMARY KEY (role_id, permission_id)
)