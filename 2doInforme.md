## 1. Arquitectura

El sistema está diseñado para ejecutar scrapers distribuidos organizados en un anillo lógico basado en el protocolo Chord, con las siguientes características:

- Clientes: Interactúan con los servidores a través de una API-REST para enviar solicitudes y recibir datos.
- Servidores: Forman un sistema distribuido basado en Chord, que asigna tareas y almacena datos de forma balanceada y redundante.
- Redes de Docker:
  - Red de clientes: Aloja instancias que interactúan con servidores mediante API-REST.
  - Red de servidores: Los servidores del anillo Chord se comunican entre sí usando RPC.

---

## 2. Organización de su sistema distribuido

### Roles del sistema

1. Clientes:
   - Envían solicitudes para iniciar scrapers.
   - Reciben los resultados procesados.
2. Servidores:
   - Distribuyen las tareas de scraping.
   - Almacenan y replican los resultados.
   - Mantienen la sincronización y el balance del sistema.

### Distribución de servicios

- Clientes en una red de Docker independiente para facilitar la escalabilidad y evitar interferencias con la red de servidores.
- Servidores en otra red de Docker, organizados en un anillo Chord, con roles de procesamiento, almacenamiento, y coordinación.

---

## 3. Procesos

### Tipos de procesos

1. Procesos cliente:
   - Gestionan la interfaz gráfica para enviar solicitudes.
   - Procesan y almacenan resultados localmente.
2. Procesos servidor:
   - Manejan solicitudes de scraping desde clientes.
   - Coordina tareas y replicación entre nodos usando RPC.
   - Redistribuyen claves y datos en caso de cambios en el anillo.

### Organización de procesos

- Cada instancia de servidor ejecuta hilos concurrentes para manejar múltiples solicitudes simultáneamente.

### Patrón de diseño

- Modelo basado en hilos para procesar solicitudes concurrentes.
- Uso de asincronía en la API-REST para una experiencia fluida entre cliente y servidor.

---

## 4. Comunicación

### Cliente-Servidor

- Protocolo: API-REST (HTTP).
- Métodos:
  - POST: Para iniciar tareas de scraping.
  - GET: Para recuperar resultados.

### Servidor-Servidor

- Protocolo: RPC.
- Uso:
  - Sincronización de datos entre nodos.
  - Asignación de tareas y redistribución de claves en el anillo Chord.

### Comunicación entre procesos

- Uso de canales (channels) y sincronización mediante semáforos en Go para evitar condiciones de carrera.

---

## 5. Coordinación

### Sincronización

- Los nodos en el anillo Chord sincronizan claves y tareas pendientes periódicamente.
- Uso de mensajes de confirmación para garantizar que las réplicas estén actualizadas.

### Acceso exclusivo

- Implementación de bloqueos distribuidos para evitar condiciones de carrera al modificar datos replicados.

### Toma de decisiones

- Completamente descentralizada: Los nodos toman decisiones localmente en función de su rango de responsabilidad.

---

## 6. Nombrado y Localización

### Identificación

- Las URLs y datos de scraping se identifican mediante claves hash.

### Ubicación

- Los servidores son responsables de un rango de claves en el anillo.
- Las tablas de enrutamiento de Chord permiten localizar datos con un número mínimo de saltos.

---

## 7. Consistencia y Replicación

### Distribución de datos

- Cada nodo maneja datos primarios y almacena réplicas de otros nodos.

### Cantidad de réplicas

- Datos críticos replicados en al menos 3 nodos para tolerancia a fallos.

### Confiabilidad

- Actualizaciones coordinadas garantizan que las réplicas sean consistentes en todo momento.

---

## 8. Tolerancia a Fallos

### Respuesta a errores

- Si un nodo falla, sus réplicas son promovidas a datos principales.

### Nivel esperado

- El sistema tolera hasta 2 fallos simultáneos.

Toledo Daniel 46, [1/14/2025 9:34 PM]

### Fallos parciales

- Los nodos caídos temporalmente se sincronizan al reintegrarse al sistema.

---

## 9. Seguridad

### Comunicación

- Cifrado mediante TLS para proteger las conexiones API-REST y RPC.

### Diseño

- Los datos sensibles en los servidores y clientes se cifran localmente.

### Autorización y autenticación

- Uso de tokens JWT para autenticar clientes.
- Verificación de permisos en cada solicitud API-REST.
