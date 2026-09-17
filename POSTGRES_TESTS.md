# Pruebas de concurrencia de evaluaciones

Usar PostgreSQL descartable local, nunca una base productiva. El usuario necesita permiso CREATE SCHEMA. Cada ejecución crea un esquema único con tablas mínimas de prueba y lo elimina al terminar; no aplica migraciones a tablas existentes.

```sh
CAMPUS_TEST_DATABASE_URL='postgres://usuario@127.0.0.1:5432/campus_test?sslmode=disable' \
  go test ./internal/adapters/datasources/repositories/quiz_attempt -run TestPostgresLifecycle -v -count=1
```

Sin variable, tests PostgreSQL se omiten explícitamente. La suite normal continúa funcionando con `go test ./...`.

Cobertura real de repositorios:

- Ocho inicios concurrentes recuperan el mismo intento.
- El inicio toma el advisory lock del par (quiz, alumno): comprobado de forma determinista sosteniendo ese lock desde fuera y exigiendo que el inicio espere. La prueba de los ocho inicios por sí sola no detecta que se quite la serialización — en una base local el primero confirma antes de que los demás lean, así que reanudan en vez de colisionar.
- Un lote inválido no guarda respuestas parcialmente. Usa una pregunta de **otro quiz**, no una inexistente: con un id inexistente la que rechaza es la foreign key, no el chequeo de pertenencia.
- Finalización espera lock de un guardado y califica su última respuesta, no una lectura previa.
- Otra institución no puede finalizar intento ajeno.
- Entregas concurrentes cierran una sola vez; borradores posteriores se rechazan; intento vivo no se finaliza.

Cada caso se verificó por mutación: al romper deliberadamente la lógica que protege (advisory lock, lectura del borrador dentro del lock, chequeo de pertenencia), el test correspondiente falla. La guarda `submitted_at IS NULL` del UPDATE de entrega es la excepción: no es observable porque la relectura bajo el lock de fila ya rechaza antes de llegar ahí — es defensa en profundidad, correcta pero no cubierta.

Ejecutadas satisfactoriamente con PostgreSQL 16 y 18.3 locales. No sustituyen pruebas HTTP autenticadas ni pruebas de migración del esquema completo. Las carreras de corrección de tareas/rúbricas requieren tests adicionales de sus propios repositorios.
