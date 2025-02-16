import requests
import concurrent.futures
import time

# CONFIG

BASE_URL = "http://localhost:8080"  # Ajusta la URL a tu endpoint

ROUTES = {"check_health": "/health", "get_tasks": "/task"}

# UTILS


def send_request(request, i):
    """Envía una petición al servidor."""
    try:
        response = request()
        print(f"Petición {i}: {response.status_code}")
        return response
    except Exception as e:
        print(f"Error en la petición {i}: {e}")
    return False


def stress_test(
    num_requests=100, max_workers=10, build_request=lambda i: ({"status_code": 500}, i)
):
    """Ejecuta múltiples peticiones concurrentes al servidor."""
    start_time = time.time()

    with concurrent.futures.ThreadPoolExecutor(max_workers=max_workers) as executor:
        futures = {
            executor.submit(send_request, lambda: build_request(i), i): i
            for i in range(num_requests)
        }
        for future in concurrent.futures.as_completed(futures):
            try:
                future.result()
            except Exception as e:
                print(f"Error en la petición {futures[future]}: {e}")

    elapsed_time = time.time() - start_time
    print(f"\n🏁 Test completado en {elapsed_time:.2f} segundos")


# BUILDERS
check_health_builder = lambda _: requests.request(
    "GET", BASE_URL + ROUTES["check_health"], json={}
)

get_tasks_builder = lambda _: requests.request(
    "GET", BASE_URL + ROUTES["get_tasks"], json={}
)

# MAIN

if __name__ == "__main__":
    stress_test(100, 10, check_health_builder)
    time.sleep(5)
    stress_test(100, 10, get_tasks_builder)
