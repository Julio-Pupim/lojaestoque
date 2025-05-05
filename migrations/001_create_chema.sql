-- Enable foreign key support
PRAGMA foreign_keys = ON;

-- 1. Table: fornecedores (suppliers)
CREATE TABLE IF NOT EXISTS fornecedores (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nome TEXT NOT NULL
);

-- 2. Table: clientes (customers)
CREATE TABLE IF NOT EXISTS clientes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nome TEXT NOT NULL,
    telefone TEXT,
    data_cadastro DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(nome)
);

-- 3. Table: produtos (products)
CREATE TABLE IF NOT EXISTS produtos (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nome TEXT NOT NULL,
    fornecedor_id INTEGER NOT NULL,
    codigo_fornecedor TEXT NOT NULL,
    codigo_barras TEXT,
    qtd_estoque INTEGER NOT NULL DEFAULT 0,
    preco_unitario REAL NOT NULL DEFAULT 0.0,
    qtd_minima INTEGER ,
    UNIQUE(fornecedor_id, codigo_fornecedor),
    FOREIGN KEY(fornecedor_id) REFERENCES fornecedores(id)
);

-- 4. Table: compras (purchases)
CREATE TABLE IF NOT EXISTS compras (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    fornecedor_id INTEGER NOT NULL,
    data DATETIME NOT NULL,
    total REAL NOT NULL,
    FOREIGN KEY(fornecedor_id) REFERENCES fornecedores(id)
);

-- 5. Table: compras_produtos (purchase items)
CREATE TABLE IF NOT EXISTS compras_produtos (
    compra_id INTEGER NOT NULL,
    produto_id INTEGER NOT NULL,
    quantidade INTEGER NOT NULL,
    preco_unitario REAL NOT NULL,
    PRIMARY KEY(compra_id, produto_id),
    FOREIGN KEY(compra_id) REFERENCES compras(id),
    FOREIGN KEY(produto_id) REFERENCES produtos(id)
);

-- 6. Table: vendas (sales)
CREATE TABLE IF NOT EXISTS vendas (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    cliente_id INTEGER NOT NULL,
    data_venda DATETIME NOT NULL,
    total REAL NOT NULL,
    data_pagamento DATETIME,
    status_pagamento TEXT NOT NULL DEFAULT 'PENDENTE',
    FOREIGN KEY(cliente_id) REFERENCES clientes(id)
);

-- 7. Table: vendas_produtos (sale items)
CREATE TABLE IF NOT EXISTS vendas_produtos (
    venda_id INTEGER NOT NULL,
    produto_id INTEGER NOT NULL,
    quantidade INTEGER NOT NULL,
    preco_unitario REAL NOT NULL,
    total REAL NOT NULL,
    PRIMARY KEY(venda_id, produto_id),
    FOREIGN KEY(venda_id) REFERENCES vendas(id),
    FOREIGN KEY(produto_id) REFERENCES produtos(id)
);

-- 8. Table: pagamentos (optional: payment installments)
-- CREATE TABLE IF NOT EXISTS pagamentos (
--     id INTEGER PRIMARY KEY AUTOINCREMENT,
--     venda_id INTEGER NOT NULL,
--     valor_pago REAL NOT NULL,
--     data_pagamento DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
--     FOREIGN KEY(venda_id) REFERENCES vendas(id)
-- );
