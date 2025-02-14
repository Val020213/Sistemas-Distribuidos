import requests
import concurrent.futures
import time

URL = "http://10.0.10.2:8080/api/test"  # Ajusta la URL a tu endpoint


def send_request(i):
    """Función para enviar una petición al servidor."""
    try:
        response = requests.get(URL, timeout=5)  # Cambia a POST si es necesario
        print(f"Request {i}: {response.status_code} - {response.text[:50]}")
    except requests.exceptions.RequestException as e:
        print(f"Request {i} failed: {e}")


def stress_test(num_requests=100, max_workers=10):
    """Ejecuta múltiples peticiones concurrentes al servidor."""
    start_time = time.time()

    with concurrent.futures.ThreadPoolExecutor(max_workers=max_workers) as executor:
        futures = {executor.submit(send_request, i): i for i in range(num_requests)}
        for future in concurrent.futures.as_completed(futures):
            try:
                future.result()
            except Exception as e:
                print(f"Error en la petición {futures[future]}: {e}")

    elapsed_time = time.time() - start_time
    print(f"\n🏁 Test completado en {elapsed_time:.2f} segundos")


# Ejecutar la prueba con 100 peticiones y 10 hilos en paralelo
stress_test(num_requests=100, max_workers=10)
