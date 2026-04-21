CREATE TABLE IF NOT EXISTS `orcamentos` (
    `id` INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `usuario_id` INT NOT NULL,
    `categoria_id` INT NOT NULL,
    `limite` DECIMAL(10, 2) NOT NULL,
    `mes` DATE NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT NOW(),
    `updated_at` TIMESTAMP NOT NULL DEFAULT NOW() ON UPDATE NOW(),
    CONSTRAINT `fk_orcamentos_categoria_id`
        FOREIGN KEY (`categoria_id`) REFERENCES `categorias` (`id`),
    CONSTRAINT `fk_orcamentos_usuario_id`
        FOREIGN KEY (`usuario_id`) REFERENCES `usuarios` (`id`)
);
