CREATE TABLE trader (
  id SERIAL PRIMARY KEY,
  email TEXT UNIQUE NOT NULL
);

CREATE TABLE company (
  id SERIAL PRIMARY KEY,
  symbol TEXT UNIQUE NOT NULL
);

CREATE TABLE holding (
  id SERIAL PRIMARY KEY
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

CREATE TABLE quote (
  id SERIAL PRIMARY KEY,
  company_id INT REFERENCES company(id),
  price NUMERIC(9, 2),
  stamp TIMESTAMP
);

CREATE INDEX idx_quote_company ON quote(company_id)

CREATE TABLE stock (
  id SERIAL PRIMARY KEY,
  company_id INT REFERENCES company(id),
  price NUMERIC(9, 2),
  vol NUMERIC(6, 4),
  minchange NUMERIC(6, 4),
  maxchange NUMERIC(6, 4)
);

CREATE TABLE cash (
  id SERIAL PRIMARY KEY,
  trader_id INT REFERENCES trader(id),
  amount NUMERIC(14, 4)
);