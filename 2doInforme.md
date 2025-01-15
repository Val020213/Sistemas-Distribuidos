# 2do Informe de Distribuidos

---

## 1. Arquitectura

El sistema está diseñado como una red distribuida utilizando el protocolo Chord para organizar los nodos en un anillo lógico. Este enfoque permite distribuir eficientemente las tareas de scraping y garantizar la redundancia de los datos. Los clientes interactúan con el sistema a través de una **API REST**, mientras que los servidores se comunican entre sí utilizando **RPC**, lo que mejora la eficiencia y la sincronización de los datos. El sistema está distribuido en dos redes de Docker independientes: una para los clientes, que se encargan de la interacción con el usuario, y otra para los servidores, que gestionan las operaciones de scraping, almacenamiento de datos y coordinación entre nodos. Esta separación garantiza escalabilidad, aislamiento de procesos y facilita la implementación de políticas de seguridad y balanceo de carga.

---

## 2. Procesos y Servicios

El sistema está compuesto por dos tipos de procesos:

### Clientes:

Los clientes tienen la responsabilidad de gestionar la interfaz gráfica y permitir a los usuarios interactuar con el sistema. Sus funciones principales incluyen agregar nuevas URLs para scraping, listar las URLs con su estado (que puede ser `success`, `pending` o `error`), y permitir la descarga del contenido de las URLs que ya han sido procesadas con éxito. Además, los usuarios pueden reintentar las solicitudes fallidas (con estado `error`).

### Servidores:

Los servidores reciben las solicitudes de scraping de los clientes y se encargan de gestionar las tareas de scraping. Cada servidor tiene la responsabilidad de coordinar la distribución de tareas y la replicación de datos entre nodos del sistema. Los servidores también gestionan la redistribución de claves y datos en caso de que la topología del anillo cambie debido a la adición o eliminación de nodos.

---

## 3. Comunicación

### Cliente - Servidor:

La comunicación entre los clientes y los servidores se realiza a través de una **API REST** utilizando métodos HTTP estándar. El cliente envía solicitudes a los servidores utilizando **POST** para agregar URLs y **GET** para consultar el estado de las URLs o recuperar su contenido.

### Servidor - Servidor:

La comunicación entre los servidores se realiza a través de **RPC**. Este protocolo es más adecuado para la comunicación entre nodos de un sistema distribuido, ya que ofrece baja latencia y alto rendimiento. Además, se utilizan **canales** y **semáforos** para garantizar la sincronización y la consistencia de los datos replicados.

### Entre Procesos:

Dentro de cada servidor, la concurrencia se maneja utilizando las **goroutines** y **canales** de Go, lo que permite una ejecución eficiente y sincronizada de múltiples tareas simultáneas, como la gestión de solicitudes de scraping y la replicación de datos.

---

## 4. Coordinación

La coordinación entre nodos es crucial en un sistema distribuido, y en este caso, se realiza mediante el protocolo Chord. Este protocolo asegura que los nodos del anillo se mantengan sincronizados y que la replicación de datos se realice de forma coherente. Las tareas de scraping y la distribución de claves se coordinan entre nodos para garantizar que los datos estén disponibles incluso si un nodo falla. La sincronización periódica de datos entre nodos es fundamental para asegurar la consistencia.

---

## 5. Nombrado y Localización

El sistema asigna cada URL a un nodo específico del anillo Chord mediante una clave hash. Esto se hace usando una función hash consistente, lo que permite determinar rápidamente a qué nodo le corresponde cada tarea de scraping. La **tabla de enrutamiento** o **finger table** de Chord facilita la localización rápida de los datos, permitiendo que los nodos encuentren rápidamente las URLs y su contenido, incluso a gran escala.

---

## 6. Consistencia y Replicación

Los datos críticos son replicados en al menos **tres nodos** para garantizar la disponibilidad y la tolerancia a fallos. Además, se implementan mecanismos de sincronización de actualizaciones entre los nodos para asegurar que las réplicas sean consistentes. Cada vez que un nodo se une al anillo, se redistribuyen las claves y se sincronizan las réplicas para mantener la coherencia de los datos en todo momento.

---

## 7. Tolerancia a Fallos

El sistema está diseñado con tolerancia a fallos de nivel 2, lo que significa que puede continuar funcionando incluso si algunos nodos fallan temporalmente. La replicación de datos en múltiples nodos garantiza que los datos sigan estando disponibles incluso si un nodo no está disponible. Cuando un nodo cae, se reasignan sus responsabilidades a otros nodos, y cuando se reintegra al sistema, se sincroniza automáticamente con los otros nodos para recuperar la consistencia de los datos. Este mecanismo garantiza la disponibilidad continua sin afectar la experiencia del usuario.

---

## 8. Seguridad

Para proteger la comunicación entre los componentes del sistema, tanto las conexiones **API REST** como las **RPC** están cifradas mediante **TLS**, lo que asegura la confidencialidad e integridad de los datos en tránsito. Sin embargo, el sistema no implementa actualmente un mecanismo de autenticación de usuarios. Este podría ser añadido utilizando métodos como **tokens JWT** o **OAuth** para garantizar la autenticación de las solicitudes. Debido a que el sistema no almacena datos sensibles de los usuarios, el riesgo de exposición de información es mínimo.

---
