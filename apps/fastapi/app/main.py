from fastapi import FastAPI
from fastapi.responses import PlainTextResponse, StreamingResponse

import hashlib, asyncio
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
    await asyncio.sleep(0.05)
    return {'status': 'ok'}


@app.get('/json')
async def json_test():
    obj = {'name': 'test', 'value': 123, 'items': list(range(5000))}
    return obj


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

