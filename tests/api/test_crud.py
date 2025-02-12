import pytest
from httpx import AsyncClient
from tests.models.task_models import TaskCreate, TaskResponse

@pytest.mark.asyncio
class TestTaskCRUD:
    async def test_create_valid_task(self, api_client: AsyncClient, frozen_time):
        response = await api_client.post("/tasks", json={
            "name": "Valid Task",
            "start_date": frozen_time.isoformat(),
            "due_date": (frozen_time + timedelta(days=5)).isoformat(),
            "priority": 3
        })
        assert response.status_code == 201
        TaskResponse(**response.json())

    @pytest.mark.parametrize("invalid_task", [
        {"name": "AB", "start_date": "2024-01-01", "due_date": "2024-01-02", "priority": 3},  # Name too short
        {"name": "Invalid Date", "start_date": "2024-01-01", "due_date": "2023-01-01", "priority": 2},  # Due before start
        {"name": "Bad Priority", "start_date": "2024-01-01", "due_date": "2024-01-02", "priority": 0}
    ])
    async def test_create_invalid_tasks(self, api_client: AsyncClient, invalid_task):
        response = await api_client.post("/tasks", json=invalid_task)
        assert response.status_code == 400

    async def test_complete_task_workflow(self, api_client: AsyncClient, frozen_time):
        # Create
        create_res = await api_client.post("/tasks", json={
            "name": "Completable Task",
            "start_date": frozen_time.isoformat(),
            "due_date": (frozen_time + timedelta(hours=1)).isoformat(),
            "priority": 1
        })
        task_id = create_res.json()["id"]
        
        # Complete
        complete_res = await api_client.put(f"/tasks/{task_id}/complete")
        assert complete_res.status_code == 200
        assert TaskResponse(**complete_res.json()).status == "completed"