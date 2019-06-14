CREATE TABLE trader (
  id SERIAL PRIMARY KEY,
  email TEXT UNIQUE NOT NULL
);

INSERT INTO trader(email) VALUES ('admin@mockex.com');

CREATE TABLE company (
  id SERIAL PRIMARY KEY,
  symbol TEXT UNIQUE NOT NULL
);

INSERT INTO company(symbol) VALUES ('HFZ'), ('IFO');

CREATE TABLE holding (
  trader_id INT REFERENCES trader(id),
  company_id INT REFERENCES company(id),
  shares INT
);

INSERT INTO holding(trader_id, company_id, shares) VALUES (1, 1, 100);

SELECT t.email, c.symbol, h.shares
FROM holding h
LEFT JOIN company c
ON h.company_id = c.id
LEFT JOIN trader t 
ON h.trader_id = t.id;