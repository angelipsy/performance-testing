from fastapi import FastAPI
from fastapi.responses import PlainTextResponse, StreamingResponse, JSONResponse

import hashlib, asyncio, os, tempfile, json
from prometheus_fastapi_instrumentator import Instrumentator

app = FastAPI()

instrumentator = Instrumentator(
    should_group_status_codes=False,
    should_ignore_untemplated=False,
    should_respect_env_var=False, # Set to False
    should_instrument_requests_inprogress=True,
    inprogress_name="fastapi_inprogress",
    inprogress_labels=True,
)
instrumentator.instrument(app).expose(app)


@app.get("/health", response_class=PlainTextResponse)
async def health():
    return "ok"


@app.get("/cpu")
async def cpu_task(iterations: int = 100):
    data = b"benchmark test data" * 100
    for _ in range(iterations):
        hashlib.sha256(data).hexdigest()
    return {"message": f"Completed {iterations} SHA256 hashes"}


@app.get('/io')
async def io():
    # Write data to a temporary file
    temp_dir = tempfile.gettempdir()
    file_path = os.path.join(temp_dir, "io_test.txt")

    # Write 1000 lines to the file
    with open(file_path, 'w') as f:
        for i in range(1000):
            f.write(f"Line {i}: This is test data for I/O operations\n")

    # Read the file back to simulate I/O
    with open(file_path, 'r') as f:
        lines = f.readlines()

    # Clean up
    os.remove(file_path)

    return {'status': 'ok', 'lines_written': len(lines)}


@app.get('/json')
async def json_test():
    # Create a complex nested data structure
    data = {
        'status': 'success',
        'timestamp': '2024-01-01T00:00:00Z',
        'users': [
            {
                'id': i,
                'name': f'User {i}',
                'email': f'user{i}@example.com',
                'active': i % 2 == 0,
                'metadata': {
                    'role': 'admin' if i % 10 == 0 else 'user',
                    'created_at': '2024-01-01',
                    'preferences': {
                        'theme': 'dark',
                        'notifications': True
                    }
                }
            }
            for i in range(1000)
        ],
        'pagination': {
            'total': 1000,
            'page': 1,
            'per_page': 1000
        }
    }

    # Explicitly serialize to JSON and deserialize to demonstrate serialization
    json_string = json.dumps(data)
    result = json.loads(json_string)

    return JSONResponse(content=result)


async def stream_gen():
    while True:
        yield b'chunk\n'
        await asyncio.sleep(0.1)


@app.get('/stream')
async def stream():
    async def event_generator():
        chunks = 20
        for i in range(1, chunks + 1):
            yield f"chunk {i}\n"
            await asyncio.sleep(0.1)  # 100ms between chunks
    return StreamingResponse(event_generator(), media_type="text/plain")

