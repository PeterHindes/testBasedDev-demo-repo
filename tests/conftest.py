import pytest
import httpx
from datetime import datetime, timedelta
from pytest_mock import MockerFixture

@pytest.fixture(scope="session")
def anyio_backend():
    return "asyncio"

@pytest.fixture(scope="session")
async def api_client():
    async with httpx.AsyncClient(base_url="http://api-under-test:8000") as client:
        yield client

@pytest.fixture
def reference_time():
    return datetime(2024, 3, 20, 12, 0, 0)

@pytest.fixture
def frozen_time(reference_time, mocker: MockerFixture):
    mock_datetime = mocker.patch("datetime.datetime")
    mock_datetime.now.return_value = reference_time
    return reference_time