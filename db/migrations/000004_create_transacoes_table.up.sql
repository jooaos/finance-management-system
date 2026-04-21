CREATE TABLE IF NOT EXISTS `transacoes` (
    `id` INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `usuario_id` INT NOT NULL,
    `categoria_id` INT NOT NULL,
    `valor` DECIMAL(10, 2) NOT NULL,
    `data` DATE NOT NULL,
    `descricao` TEXT NULL,
    `tipo` VARCHAR(50) NOT NULL,
    `parcelas` INT NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT NOW(),
    `updated_at` TIMESTAMP NOT NULL DEFAULT NOW() ON UPDATE NOW(),
    CONSTRAINT `fk_transacoes_categoria_id`
        FOREIGN KEY (`categoria_id`) REFERENCES `categorias` (`id`),
    CONSTRAINT `fk_transacoes_usuario_id`
        FOREIGN KEY (`usuario_id`) REFERENCES `usuarios` (`id`)
);
