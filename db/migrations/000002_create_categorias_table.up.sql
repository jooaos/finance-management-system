CREATE TABLE IF NOT EXISTS `categorias` (
    `id` INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `nome` VARCHAR(120) NOT NULL,
    `usuario_id` INT NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT NOW(),
    `updated_at` TIMESTAMP NOT NULL DEFAULT NOW() ON UPDATE NOW(),
    CONSTRAINT `fk_categorias_usuario_id`
        FOREIGN KEY (`usuario_id`) REFERENCES `usuarios` (`id`)
);
