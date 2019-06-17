INSERT INTO trader(email) VALUES ('admin@mockex.com');

INSERT INTO company(symbol) VALUES ('HFZ'), ('IFO');

INSERT INTO holding(trader_id, company_id, shares) VALUES (1, 1, 100);

INSERT INTO trade(trader_id, company_id, shares, price) VALUES (1, 1, 100, 134.66);

SELECT t.email, c.symbol, h.shares
FROM holding h
LEFT JOIN company c
ON h.company_id = c.id
LEFT JOIN trader t 
ON h.trader_id = t.id;