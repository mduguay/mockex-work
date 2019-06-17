CREATE TABLE trader (
  id SERIAL PRIMARY KEY,
  email TEXT UNIQUE NOT NULL
);

CREATE TABLE company (
  id SERIAL PRIMARY KEY,
  symbol TEXT UNIQUE NOT NULL
);


CREATE TABLE holding (
  trader_id INT REFERENCES trader(id),
  company_id INT REFERENCES company(id),
  shares INT
);

CREATE INDEX idx_holding_trader_id ON holding(trader_id);

CREATE TABLE trade (
  id SERIAL PRIMARY KEY,
  trader_id INT REFERENCES trader(id),
  company_id INT REFERENCES company(id),
  direction BOOLEAN,
  shares INT,
  price NUMERIC(9, 2)
);