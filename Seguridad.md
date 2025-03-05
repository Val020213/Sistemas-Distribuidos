### **1. Generar Certificados y CA Privada**

Crea una **Autoridad de Certificación (CA)** autofirmada y emite certificados para cada componente:

```bash
# Generar CA
openssl req -x509 -newkey rsa:4096 -keyout ca.key -out ca.crt -days 365 -nodes

# Generar certificados para:
# - Clientes (Next.js)
openssl req -newkey rsa:2048 -nodes -keyout client.key -out client.csr -subj "/CN=client"
openssl x509 -req -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out client.crt -days 365

# - Router (Go)
openssl req -newkey rsa:2048 -nodes -keyout router.key -out router.csr -subj "/CN=router"
openssl x509 -req -in router.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out router.crt -days 365

# - Servidores (Chord)
openssl req -newkey rsa:2048 -nodes -keyout server.key -out server.csr -subj "/CN=server"
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 365
```

---

### **2. Configuración de Seguridad por Componente**

#### **Clientes (Next.js)**

- **Objetivo:** Autenticar al cliente frente al router usando mTLS.
- **Implementación:**

  - Usa `https` con certificados en las peticiones al router.

  ```javascript
  // Ejemplo con axios
  import https from "https";
  import fs from "fs";

  const agent = new https.Agent({
    cert: fs.readFileSync("client.crt"),
    key: fs.readFileSync("client.key"),
    ca: fs.readFileSync("ca.crt"),
  });

  axios.get("https://router:3000/api/data", { httpsAgent: agent });
  ```

#### **Router (Go)**

- **Objetivo:**
  - Validar certificados de los clientes (mTLS).
  - Comunicarse con los servidores usando mTLS.
- **Implementación:**

  - Configura el servidor HTTP para requerir mTLS:

  ```go
  // Servidor del router (recibe peticiones de clientes)
  cert, _ := tls.LoadX509KeyPair("router.crt", "router.key")
  caCert, _ := ioutil.ReadFile("ca.crt")
  caCertPool := x509.NewCertPool()
  caCertPool.AppendCertsFromPEM(caCert)

  tlsConfig := &tls.Config{
    Certificates: []tls.Certificate{cert},
    ClientAuth:   tls.RequireAndVerifyClientCert,
    ClientCAs:    caCertPool,
  }

  server := &http.Server{
    Addr:      ":3000",
    TLSConfig: tlsConfig,
  }
  ```

  - Configura el cliente HTTP para hablar con los servidores:

  ```go
  // Cliente del router (envía peticiones a servidores)
  cert, _ := tls.LoadX509KeyPair("router.crt", "router.key")
  caCert, _ := ioutil.ReadFile("ca.crt")
  caCertPool := x509.NewCertPool()
  caCertPool.AppendCertsFromPEM(caCert)

  transport := &http.Transport{
    TLSClientConfig: &tls.Config{
      Certificates: []tls.Certificate{cert},
      RootCAs:      caCertPool,
    },
  }
  client := &http.Client{Transport: transport}
  ```

#### **Servidores (Chord en Go)**

- **Objetivo:**
  - Validar certificados del router (mTLS en REST).
  - Usar mTLS en las comunicaciones gRPC entre nodos.
- **Implementación:**

  - **REST Server (recibe peticiones del router):**

  ```go
  cert, _ := tls.LoadX509KeyPair("server.crt", "server.key")
  caCert, _ := ioutil.ReadFile("ca.crt")
  caCertPool := x509.NewCertPool()
  caCertPool.AppendCertsFromPEM(caCert)

  tlsConfig := &tls.Config{
    Certificates: []tls.Certificate{cert},
    ClientAuth:   tls.RequireAndVerifyClientCert,
    ClientCAs:    caCertPool,
  }

  server := &http.Server{
    Addr:      ":8080",
    TLSConfig: tlsConfig,
  }
  ```

  - **gRPC Server/Client (comunicación entre nodos):**

  ```go
  // Servidor gRPC
  creds := credentials.NewTLS(&tls.Config{
    Certificates: []tls.Certificate{cert},
    ClientAuth:   tls.RequireAndVerifyClientCert,
    ClientCAs:    caCertPool,
  })
  grpcServer := grpc.NewServer(grpc.Creds(creds))

  // Cliente gRPC
  creds := credentials.NewTLS(&tls.Config{
    Certificates: []tls.Certificate{cert},
    RootCAs:      caCertPool,
  })
  conn, _ := grpc.Dial("server-node:50051", grpc.WithTransportCredentials(creds))
  ```

---

### **3. Configuración de Docker**

- **Monta los certificados en los contenedores:**

  ```yaml
  # docker-compose.yml
  services:
    client:
      build: ./client
      volumes:
        - ./certs/client.crt:/app/client.crt
        - ./certs/client.key:/app/client.key
        - ./certs/ca.crt:/app/ca.crt

    router:
      build: ./router
      volumes:
        - ./certs/router.crt:/app/router.crt
        - ./certs/router.key:/app/router.key
        - ./certs/ca.crt:/app/ca.crt

    server:
      build: ./server
      volumes:
        - ./certs/server.crt:/app/server.crt
        - ./certs/server.key:/app/server.key
        - ./certs/ca.crt:/app/ca.crt
  ```

---

### **4. Validaciones Clave**

1. **Clientes → Router:** Solo clientes con `client.crt` pueden conectarse.
2. **Router → Servidores:** Solo el router con `router.crt` puede enviar peticiones.
3. **Servidores entre sí:** Solo nodos con `server.crt` pueden unirse al anillo Chord.

---

### **5. Consideraciones Adicionales**

- **Renovación de Certificados:** Para un proyecto escolar, los certificados con vigencia de 1 año son suficientes.
- **Almacenamiento Seguro:** En producción, usa Docker Secrets o Vault, pero para tu caso, los volúmenes son aceptables.
- **Firewall:** Asegura que los puertos expuestos en Docker estén restringidos (ej: solo el router expone el puerto 3000).

Con esto, garantizas que todas las comunicaciones sean cifradas y autenticadas sin necesidad de usuarios o servicios externos.
