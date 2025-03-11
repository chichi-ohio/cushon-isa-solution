-- Insert sample funds
INSERT INTO funds (id, name, description, risk_level, min_investment, max_investment) VALUES
('CUSHON_EQ', 'Cushon Equities Fund', 'A diversified portfolio of global equities focused on long-term growth', 'High', 100.00, 1000000.00),
('CUSHON_BOND', 'Cushon Bond Fund', 'A conservative fund investing in high-quality government and corporate bonds', 'Low', 500.00, 500000.00),
('CUSHON_ESG', 'Cushon ESG Impact Fund', 'Sustainable investments in companies with strong environmental and social governance', 'Medium', 250.00, 750000.00),
('CUSHON_TECH', 'Cushon Technology Fund', 'Focused on innovative technology companies with high growth potential', 'High', 1000.00, 1000000.00);

-- Insert sample performance data
INSERT INTO performances (fund_id, date, value, change) VALUES
-- Cushon Equities Fund performance
('CUSHON_EQ', CURRENT_DATE - INTERVAL '30 days', 100.00, 0.00),
('CUSHON_EQ', CURRENT_DATE - INTERVAL '20 days', 102.50, 2.50),
('CUSHON_EQ', CURRENT_DATE - INTERVAL '10 days', 105.00, 2.44),
('CUSHON_EQ', CURRENT_DATE, 108.00, 2.86),

-- Cushon Bond Fund performance
('CUSHON_BOND', CURRENT_DATE - INTERVAL '30 days', 100.00, 0.00),
('CUSHON_BOND', CURRENT_DATE - INTERVAL '20 days', 100.80, 0.80),
('CUSHON_BOND', CURRENT_DATE - INTERVAL '10 days', 101.50, 0.69),
('CUSHON_BOND', CURRENT_DATE, 102.00, 0.49),

-- Cushon ESG Impact Fund performance
('CUSHON_ESG', CURRENT_DATE - INTERVAL '30 days', 100.00, 0.00),
('CUSHON_ESG', CURRENT_DATE - INTERVAL '20 days', 101.50, 1.50),
('CUSHON_ESG', CURRENT_DATE - INTERVAL '10 days', 103.00, 1.48),
('CUSHON_ESG', CURRENT_DATE, 105.00, 1.94),

-- Cushon Technology Fund performance
('CUSHON_TECH', CURRENT_DATE - INTERVAL '30 days', 100.00, 0.00),
('CUSHON_TECH', CURRENT_DATE - INTERVAL '20 days', 104.00, 4.00),
('CUSHON_TECH', CURRENT_DATE - INTERVAL '10 days', 107.00, 2.88),
('CUSHON_TECH', CURRENT_DATE, 112.00, 4.67); 