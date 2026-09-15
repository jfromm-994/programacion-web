CREATE TABLE juegos (
    id SERIAL PRIMARY KEY,
    titulo VARCHAR(255) UNIQUE NOT NULL,
    genero VARCHAR(100) NOT NULL,
    logros_totales INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE copias (
    id SERIAL PRIMARY KEY,
    juego_id INT NOT NULL REFERENCES juegos(id) ON DELETE CASCADE,
    plataforma VARCHAR(100) NOT NULL,
    horas_jugadas NUMERIC(10, 2) NOT NULL DEFAULT 0.0,
    logros_obtenidos INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_juego_plataforma UNIQUE (juego_id, plataforma)
);