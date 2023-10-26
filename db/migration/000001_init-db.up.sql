CREATE TABLE clients (
    id BIGSERIAL PRIMARY KEY,
    fio VARCHAR(50) NOT NULL,
    phone VARCHAR(20) NOT NULL
);

CREATE TABLE balances(
    id BIGSERIAL PRIMARY KEY,
    client_id BIGINT NOT NULL,
    balance INT NOT NULL
);

ALTER TABLE ONLY balances
    ADD CONSTRAINT fk_clients FOREIGN KEY (client_id) REFERENCES clients(id);