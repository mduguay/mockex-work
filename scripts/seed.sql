INSERT INTO trader(email) VALUES ('trader@mockex.com');

INSERT INTO company(symbol) VALUES ('HFZ'), ('IFO');

INSERT INTO holding(trader_id, company_id, shares) VALUES (1, 1, 100);

INSERT INTO trade(trader_id, company_id, shares, price) VALUES (1, 1, 100, 134.66);

INSERT INTO price(company_id, price, stamp) VALUES (1, 11.34, NOW()), (2, 6.66, NOW());

INSERT INTO stock(company_id, price, vol, minchange, maxchange) 
VALUES (1, 11.34, 0.02, 0.0, 1.0), (2, 6.67, 0.03, 0.0, 1.2);

SELECT t.email, c.symbol, h.shares
FROM holding h
LEFT JOIN company c
ON h.company_id = c.id
LEFT JOIN trader t 
ON h.trader_id = t.id;